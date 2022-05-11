package General

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	err error
)

func init() {
	// 出厂设置
	RootPath, _ := os.Getwd()
	viper.Set("RootPath", RootPath)
	viper.Set("LogPath", RootPath+"/logs/")
	// 初始化日志
	logInit()
	// 初始化viper
	ViperInit()
}

/*
	初始化日志
*/
func ViperInit() {
	// 初始化配置文件
	//fileInit(false, viper.GetString("RootPath"))
	log.Printf("Viper Config Path: [%v]", viper.GetString("RootPath")+"/config.yaml")
	log.Printf("Logs Path: [%v]", viper.GetString("LogPath"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("RootPath"))
	fileInit(false, viper.GetString("RootPath")+"/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		log.Errorf("[viper] 读取配置文件失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		log.Infof("[viper] 读取配置文件成功")
	}
	if viper.GetString("WECHATY.TOKEN") == "" && viper.GetString("BOT.adminID") == "" {
		viper.Set("BOT.adminID", "wxid_xxxxx")
		viper.Set("BOT.NAME", "xxxxxxx")
		viper.Set("WECHATY.ENDPOINT", "127.0.0.1:25001")
		viper.Set("WECHATY.TOKEN", "insecure_xxxxxxxxxxxxxxxxxxxxxx")
		viper.Set("WECHATY.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", true)
		viper.Set("GROUP.XXXXX", "xxxxx@chatroom")
		viper.Set("GROUP.PASS.1", "xxxxx")
		viper.Set("GROUP.PASS.2", "xxxxx")
		viper.Set("DING.SECRET", "xxxxxxxxxxxxxxxxx")
		viper.Set("DING.TOKEN", "xxxxxxxxxxxxxxxxx")
		viper.Set("WXOPENAI.ENV", "online")
		viper.Set("WXOPENAI.SIGNURL", "https://openai.weixin.qq.com/openapi/sign/")
		viper.Set("WXOPENAI.URL", "https://openai.weixin.qq.com/openapi/aibot/")
		viper.Set("WXOPENAI.TOKEN", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
		viper.Set("TULING.TOKEN", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&info=")
		viper.Set("TULING.URL", "http://www.tuling123.com/openapi/api?key=")
		viper.Set("QINGYUNKE.AppID", 0)
		viper.Set("QINGYUNKE.Key", "free")
		viper.Set("QINGYUNKE.URL", "https://api.qingyunke.com/api.php?key=free&appid=0&msg=")
	}
	ViperWrite()
}

/*
	初始化日志
*/
func logInit() {
	// 设置日志格式
	// 创建日志文件夹
	fileInit(true, viper.GetString("RootPath")+"/logs")
	log.SetLevel(log.WarnLevel | log.InfoLevel | log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf(" %v:%v ", path.Base(frame.File), strconv.Itoa(frame.Line))
			function = strings.Replace(path.Ext(frame.Function), ".", "", -1)
			return function, fileName
		},
		PrettyPrint: true,
	})
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   viper.Get("LogPath").(string) + "/wechatBot.log", //日志文件位置
		MaxSize:    50,                                               // 单文件最大容量,单位是MB
		MaxBackups: 1,                                                // 最大保留过期文件个数
		MaxAge:     365,                                              // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,                                             // 是否需要压缩滚动日志, 使用的 gzip 压缩
	}))
}

/*
	初始化文件/夹
*/
func fileInit(fileAttributes bool, fileName string) {
	if _, err = os.Stat(fileName); err != nil {
		if fileAttributes {
			if err = os.MkdirAll(fileName, os.ModePerm); err != nil {
				log.Errorf("创建[%v] 目录失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
			} else {
				log.Infof("创建[%v] 目录成功", fileName)
			}
		} else {
			var f *os.File
			if f, err = os.Create(fileName); err != nil {
				log.Errorf("创建[%v] 文件失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
			} else {
				log.Infof("创建[%v] 文件成功", fileName)
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					log.Errorf("关闭[%v] 文件失败, Error: [%v] CoptRight: [%s]", fileName, err, Copyright(make([]uintptr, 1)))
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
	viper.Set("chat", "")
	if err = viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		log.Errorf("Viper Write file Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		log.Infof("Viper Write file Success")
	}
}
