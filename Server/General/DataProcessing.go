package General

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/wechaty/go-wechaty/wechaty"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// 出厂设置
	RootPath, _ := os.UserConfigDir()
	RootPath = RootPath + "/Go-Wechaty-Bot"
	viper.Set("RootPath", RootPath)
	viper.Set("LogPath", RootPath+"/logs")
	//logrus.Printf("Logs Path: [%v]", viper.GetString("LogPath"))
	// 初始化viper
	ViperInit()
	// 初始化日志
	logInit()
}

func NewGlobleService() *GlobleService {
	return globleServiceInstance
}

func (g *GlobleService) SetBot(bot *wechaty.Wechaty) {
	g.Bot = bot
}

func (g *GlobleService) GetBot() *wechaty.Wechaty {
	return g.Bot
}

/*
	初始化日志
*/
func ViperInit() {
	// 初始化配置文件
	//fileInit(false, viper.GetString("RootPath"))
	//logrus.Printf("Viper Config Path: [%v]", viper.GetString("RootPath")+"/config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("RootPath"))
	fileInit(true, viper.GetString("RootPath"))
	fileInit(false, viper.GetString("RootPath")+"/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		logrus.Errorf("[viper] 读取配置文件失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		logrus.Infof("[viper] 读取配置文件成功")
	}
	if viper.GetString("WECHATY.WECHATY_TOKEN") == "" && viper.GetString("BOT.ADMINID") == "" {
		var cmd *exec.Cmd
		if err = GetGatewayConfig("https://ghproxy.com/https://raw.githubusercontent.com/XRSec/Go-Wechaty-Bot/main/Server/example.yaml", viper.GetString("RootPath")+"/config.yaml"); err != nil {
			logrus.Errorf("[viper] 初始化配置文件失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		}
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "start", viper.GetString("RootPath")+"/config.yaml")
		} else {
			cmd = exec.Command("open", viper.GetString("RootPath")+"/config.yaml")
		}
		if err = cmd.Run(); err != nil {
			return
		}

		if err = viper.ReadInConfig(); err != nil {
			logrus.Errorf("[viper] 读取配置文件失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
		} else {
			logrus.Infof("[viper] 读取配置文件成功")
		}
		logrus.Infof("请检查配置文件! (CTRL + Z 撤销覆盖)")
		os.Exit(1)
	}
	ViperWrite()
}

/*
	初始化日志
*/
func logInit() {
	// 设置日志格式
	log.SetPrefix("\x1b[1;32m[Go-Wechaty-Bot] \x1b[0m")
	// \x1b[%dm%s\x1b[0m
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 创建日志文件夹
	fileInit(true, viper.GetString("LogPath"))

	// TODO Debug Model
	//logrus.SetLevel(logrus.PanicLevel | logrus.FatalLevel | logrus.ErrorLevel | logrus.WarnLevel | logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf(" %v:%v ", path.Base(frame.File), strconv.Itoa(frame.Line))
			function = strings.Replace(path.Ext(frame.Function), ".", "", -1)
			return function, fileName
		},
		//PrettyPrint: true,
	})

	//if viper.GetString("SYSLOG.STATUS") == "true" {
	//	fmt.Printf("[syslog] 启动日志记录\n")
	//	hook, err := lSyslog.NewSyslogHook(viper.GetString("SYSLOG.PROTOCOL"), viper.GetString("SYSLOG.REMOTE_ADDRESS"), syslog.LOG_INFO, "Go-Wechaty-Bot")
	//	if err != nil {
	//		fmt.Printf("SysLog Error: %v\n", err)
	//	}
	//	logrus.AddHook(hook)
	//	fmt.Printf("[syslog] 日志记录器初始化完成\n")
	//	return
	//}
	logrus.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   viper.GetString("LogPath") + "/Go-Wechaty-Bot.log", //日志文件位置
		MaxSize:    50,                                                 // 单文件最大容量,单位是MB
		MaxBackups: 1,                                                  // 最大保留过期文件个数
		MaxAge:     365,                                                // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,                                               // 是否需要压缩滚动日志, 使用的 gzip 压缩
	}))
}

/*
	初始化文件/夹
*/
func fileInit(fileAttributes bool, fileName string) {
	if _, err = os.Stat(fileName); err != nil {
		if fileAttributes {
			if err = os.MkdirAll(fileName, os.ModePerm); err != nil {
				logrus.Errorf("创建[%v] 目录失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
			} else {
				logrus.Infof("创建[%v] 目录成功", fileName)
			}
		} else {
			var f *os.File
			if f, err = os.Create(fileName); err != nil {
				logrus.Errorf("创建[%v] 文件失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
			} else {
				logrus.Infof("创建[%v] 文件成功", fileName)
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					logrus.Errorf("关闭[%v] 文件失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
				}
			}(f)
		}
	}
}

/*
	ViperWrite()
	写入配置文件
*/
func ViperWrite() {
	if err = viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		logrus.Errorf("Viper Write file Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		logrus.Infof("Viper Write file Success")
	}
}

/*
	DockerRestart() 重启容器
*/
func DockerRestart() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "chcp 437 && docker restart "+viper.GetString("BOT.DOCKER"))
	} else {
		cmd = exec.Command("/bin/sh", "-c", "docker restart "+viper.GetString("BOT.DOCKER"))
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		logrus.Errorf("Docker Restart Error: %v", err)
	} else {
		logrus.Infof("Docker Restart Success: %s", out)
	}
}

func DingSend(userID, msg string) {
	if NightMode(userID) {
		cli := dingtalk.InitDingTalkWithSecret(viper.GetString("Ding.TOKEN"), viper.GetString("Ding.SECRET"))
		if err = cli.SendMarkDownMessage(msg, msg); err != nil {
			logrus.Errorf("DingMessage Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		}
		logrus.Infof("DingTalk 通知成功! Copyright: [%s]", Copyright(make([]uintptr, 1)))
	} else {
		logrus.Infof("现在处于夜间模式，请在白天使用 CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
}
