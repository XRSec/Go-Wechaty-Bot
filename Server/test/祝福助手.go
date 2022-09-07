package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	. "github.com/wechaty/go-wechaty/wechaty"
	puppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
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

var (
	lists []_interface.IContact
	err   error
)

func init() {
	RootPath, _ := os.UserConfigDir()
	RootPath = RootPath + "/wechatBot"
	viper.Set("RootPath", RootPath)
	viper.Set("LogPath", RootPath+"/logs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("RootPath"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("viper.ReadInConfig() Error: [%v]", err)
	}
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf(" %v:%v ", path.Base(frame.File), strconv.Itoa(frame.Line))
			function = strings.Replace(path.Ext(frame.Function), ".", "", -1)
			return function, fileName
		},
		//PrettyPrint: true,
	})
}

func main() {
	var bot = NewWechaty(WithPuppetOption(puppet.Option{
		Token:    viper.GetString("WECHATY.WECHATY_TOKEN"),
		Endpoint: viper.GetString("WECHATY.WECHATY_ENDPOINT"),
	}))
	bot.OnScan(func(context *Context, qrCode string, status schemas.ScanStatus, data string) {
		log.Infof("Scan QR Code: [ https://wechaty.js.org/qrcode/%s ] [%s] [%s]", qrCode, status, data)
	}).OnLogin(func(context *Context, user *user.ContactSelf) {
		log.Infof(user.Name() + " login")
	}).OnLogout(func(context *Context, user *user.ContactSelf, reason string) {
		log.Infof(user.Name() + " logout")
	}).OnError(func(context *Context, err error) {
		log.Error(err)
	}).OnMessage(func(context *Context, message *user.Message) {
		if message.Self() || message.Talker().ID() == viper.GetString("BOT.ADMINID") {
		} else {
			return
		}
		if message.Type() != schemas.MessageTypeText {
			return
		}
		if !strings.Contains(message.Text(), "节日祝福") {
			return
		}
		if message.Text()[0:13] == "节日祝福 " {
			if _, err = os.Stat("friend.json"); err != nil {
				getAllToFile(message.GetWechaty().Contact())
				SayMessage(message, fmt.Sprintf("群发开始,共计%v人", len(lists)))
			} else {
				readFromFile()
				SayMessage(message, fmt.Sprintf("群发继续, 剩余%v人", len(lists)))
			}
			//if msg == "" {
			//	if message.Text()[0:7] == "群发 " {
			//		msg = message.Text()[7:]
			//	}
			//	if message.Text()[0:8] == "forward " {
			//		msg = message.Text()[8:]
			//	}
			//}
			for i := 1; i < len(lists); i = 0 {
				if _, err = message.GetWechaty().Contact().Load(lists[i].ID()).Say(fmt.Sprintf("嗨,亲爱的%v, %v", lists[i].Name(), message.Text()[13:])); err != nil {
					_, _ = message.Say(fmt.Sprintf("群发失败, 剩余%v人未发送成功", len(lists)))
					writeToFile()
					return
				}
				lists = append(lists[:0], lists[(1):]...)
				writeToFile()
				time.Sleep(time.Second * 8)
			}
			if err := os.Remove("friend.json"); err != nil {
				log.Errorf("os.Remove Error: [%v]", err)
				return
			}
		}
		if message.Text()[0:19] == "节日祝福测试 " {
			SayMessage(message, fmt.Sprintf("嗨, 亲爱的%v, %v", message.Talker().Name(), message.Text()[19:]))
		}
	})
	bot.DaemonStart()
}

func getAllToFile(c _interface.IContactFactory) {
	var lists2 []_interface.IContact
	lists = c.FindAll(nil)
	log.Infoln("ContactList: 加载成功")
	for _, v := range lists {
		if v.Type() != schemas.ContactTypePersonal {
			continue
		}
		if !v.Friend() {
			continue
		}
		lists2 = append(lists2, v)
	}
	lists = lists2
	writeToFile()
}

func readFromFile() {
	result, err := ioutil.ReadFile("friend.json")
	if err != nil {
		log.Errorf("ioutil.ReadFile Error: [%v]", err)
		return
	}
	err = json.Unmarshal(result, &lists)
	if err != nil {
		log.Errorf("json.Unmarshal Error: [%v]", err)
		return
	}
}

func writeToFile() {
	result, err := json.Marshal(lists)
	if err != nil {
		log.Errorf("json.Marshal Error: [%v]", err)
		return
	}
	_ = ioutil.WriteFile("friend.json", result, 0644)
}

func SayMessage(message *user.Message, msg string) {
	if _, err = message.Say(msg); err != nil {
		log.Errorf("[SayMessage] [%v], error: %v", msg, err)
	}
}
