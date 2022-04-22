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
	rootPath, _ := os.Getwd()
	viper.Set("rootPath", rootPath)
	viper.Set("logPath", rootPath+"/logs/")
	// 初始化日志
	logInit()
	// 初始化viper
	viperInit()
}

/*
	初始化日志
*/
func viperInit() {
	// 初始化配置文件
	fileInit(false, viper.GetString("rootPath"))
	log.Printf("Viper Config Path: [%v]", viper.GetString("rootPath")+"/config.yaml")
	log.Printf("Logs Path: [%v]", viper.GetString("logPath"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("rootPath"))
	fileInit(false, viper.GetString("rootPath")+"/config.yaml")
}

/*
	初始化日志
*/
func logInit() {
	// 设置日志格式
	// 创建日志文件夹
	fileInit(true, viper.GetString("rootPath")+"/logs")
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
		Filename:   viper.Get("logPath").(string) + "/wechatBot.log", //日志文件位置
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
				log.Printf("创建[%v] 目录失败, Error: [%v]", fileName, err)
			} else {
				log.Printf("创建[%v] 目录成功", fileName)
			}
		} else {
			var f *os.File
			if f, err = os.Create(fileName); err != nil {
				log.Printf("创建[%v] 文件失败, Error: [%v]", fileName, err)
			} else {
				log.Printf("创建[%v] 文件成功", fileName)
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					log.Printf("关闭[%v] 文件失败, Error: [%v]", fileName, err)
				}
			}(f)
		}
	}
}

/*
	WechatBotInit()
	检查环境变量
*/
func WechatBotInit() {
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("[viper] 读取配置文件失败, Error: [%v]", err)
	} else {
		log.Printf("[viper] 读取配置文件成功")
	}
	// WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT
	if viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT") == "" {
		viper.Set("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", "false")
	}
	if err = os.Setenv("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT")); err != nil {
		log.Errorf("设置环境变量: [WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT] 失败, Error: [%v]", err)
	}
	// ENDPOINT
	if viper.GetString("wechaty.endpoint") == "" {
		log.Errorf("请填写服务器地址, endpoint: [%v]", viper.GetString("wechaty.endpoint"))
		viper.Set("wechaty.endpoint", "Please Fill In Your Server Address")
		os.Exit(1)
	}
	// TOKEN
	if viper.GetString("wechaty.token") == "" {
		viper.Set("wechaty.token", "Please Fill In Your Token")
		log.Errorf("请填写服务器 Token: [%v]", viper.GetString("wechaty.token"))
		os.Exit(1)
	}
}

/*
	ViperWrite()
	写入配置文件
*/
func ViperWrite() {
	viper.Set("chat", "")
	if err = viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		log.Errorf("Viper Write file Error: %v", err)
	} else {
		log.Printf("Viper Write file Success")
	}
}
