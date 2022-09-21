package main

import (
	"encoding/json"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	. "github.com/wechaty/go-wechaty/wechaty"
	puppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/*

*****************************
	TODO
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build test.go

init().RootPath
main().Endpoint
main().Token
heartbeat().roomID

*****************************

*/

func init() {
	RootPath, _ := os.UserConfigDir()
	RootPath = RootPath + "/Go-Wechaty-Bot"
	viper.Set("RootPath", RootPath)
	viper.Set("LogPath", RootPath+"/logs")
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf(" %v:%v ", path.Base(frame.File), strconv.Itoa(frame.Line))
			function = strings.Replace(path.Ext(frame.Function), ".", "", -1)
			return function, fileName
		},
		//PrettyPrint: true,
	})
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   viper.GetString("LogPath") + "/Go-Wechaty-Bot.log", // 日志文件位置
		MaxSize:    50,                                                 // 单文件最大容量,单位是MB
		MaxBackups: 1,                                                  // 最大保留过期文件个数
		MaxAge:     365,                                                // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,                                               // 是否需要压缩滚动日志, 使用的 gzip 压缩
	}))
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error: %v", err)
		}
	}()
	Token, Endpoint := tokenGet()
	info(Token, Endpoint)
	for i := 0; i <= 10; i++ {
		var bot = NewWechaty(WithPuppetOption(puppet.Option{
			Endpoint: Endpoint,
			Token:    Token,
		}))
		botMain(bot)
		heartbeat(bot)
		daemonStart(bot, i)
	}
}

func botMain(bot *Wechaty) {
	bot.OnScan(func(context *Context, qrCode string, status schemas.ScanStatus, data string) {
		log.Infof("Scan QR Code: [ https://wechaty.js.org/qrcode/%s ] [%s] [%s]", qrCode, status, data)
	}).OnLogin(func(context *Context, user *user.ContactSelf) {
		log.Infof(user.Name() + " login")
	}).OnLogout(func(context *Context, user *user.ContactSelf, reason string) {
		log.Infof(user.Name() + " logout")
	}).OnError(func(context *Context, err error) {
		log.Error(err)
	}).OnMessage(func(context *Context, message *user.Message) {
		if message.Room() != nil {
			log.Infof("{User: %v, ID: %v}, {Group: %v, ID: %v}, Context: %v", message.Talker().Name(), message.Talker().ID(), message.Room().Topic(), message.Room().ID(), message.Text())
		} else {
			log.Infof("{User: %v, ID: %v}, Context: %v", message.Talker().Name(), message.Talker().ID(), message.Text())
		}

		if message.Type() != schemas.MessageTypeText || message.Age() > 2*60*time.Second || message.Text() != "ding" {
			return
		}
		if _, err := message.Say(result()); err != nil {
			log.Errorf("message.Say Error: [%v]", err)
		}
	})
}

func heartbeat(bot *Wechaty) {
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(nyc))
	if _, err := c.AddFunc("0 23 * * *", func() {
		var roomID _interface.IRoom
		//if roomID = bot.Room().Find(&schemas.RoomQueryFilter{Id: "roomID@chatroom"}); roomID == nil {
		if roomID = bot.Room().Find("Debug"); roomID == nil {
			dingSend("RoomID Find Error")
			log.Infof("RoomID Find Error")
			return
		}

		if _, err := roomID.Say(result()); err != nil {
			dingSend("failed to send messages")
			log.Errorf("onHeartbeat Say Error: [%v]", err)
			return
		}
		log.Infof("Heartbeat Say Success")
	}); err != nil {
		dingSend("Heartbeat Cron Add Error: " + err.Error())
		log.Errorf("Heartbeat Cron Add Error: [%v]", err)
	}
	c.Start()
}

func daemonStart(bot *Wechaty, i int) {
	if err := bot.Start(); err != nil {
		log.Infof("[main] Bot 错误: %v", err)
		if i > 10 {
			os.Exit(0)
		}
		log.Printf("正在重新启动程序, 当前重试次数: 第%v次", i)
		time.Sleep(10 * time.Second)
	} else {
		i = 0
		var quitSig = make(chan os.Signal)
		signal.Notify(quitSig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quitSig:
			log.Fatal("程序退出!")
		}
	}
}

/* SET YOUT TOKEN */
func tokenGet() (Token, Endpoint string) {
	/**************************/

	Token = "insecure_xxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	Endpoint = "127.0.0.1:25000"

	/**************************/

	if _, err := os.Stat(viper.GetString("RootPath") + "/config.yaml"); err == nil {
		log.Info("found config.yaml")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(viper.GetString("RootPath"))
		if viper.ReadInConfig() == nil {
			token := viper.GetString("WECHATY.WECHATY_TOKEN")
			endpoint := viper.GetString("WECHATY.WECHATY_ENDPOINT")
			if token == "" || endpoint == "" {
				return Token, Endpoint
			}
			return token, endpoint
		}
		log.Infof(err.Error())
	}
	return Token, Endpoint
}

func info(Token, Endpoint string) {
	fmt.Println("\n\n--------------------")
	fmt.Printf("LogPath: %v/Go-Wechaty-Bot.log \n", viper.GetString("LogPath"))
	fmt.Printf("Config: %v/config.yaml \n", viper.GetString("RootPath"))
	fmt.Printf("Token: %v \n", Token)
	fmt.Printf("Endpoint: %v \n", Endpoint)
	fmt.Println("--------------------\n\n")
}

func dingSend(msg string) {
	cli := dingtalk.InitDingTalkWithSecret(viper.GetString("Ding.TOKEN"), viper.GetString("Ding.SECRET"))
	if err := cli.SendMarkDownMessage(msg, msg); err != nil {
		log.Errorf("DingMessage Error: [%v]", err)
		return
	}
	log.Infof("DingTalk 通知成功!")
}

func result() string {
	type Jin struct {
		Content  string
		origin   string
		Author   string
		Category string
	}

	var (
		jin  Jin
		resp *http.Response
		err  error
	)
	// 发起请求
	resp, err = http.Get("https://v1.jinrishici.com/shuqing")
	if err != nil {
		log.Errorf("今日诗词接口请求错误: [%v] ", err)
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("今日诗词请求 Close Body Error: [%v]", err)
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&jin); err != nil {
		log.Errorf("今日诗词接口解析错误: [%v]", err)
	}
	if jin.Content == "" {
		jin.Content = "情如之何，暮涂为客，忍堪送君。"
	}
	return jin.Content
}
