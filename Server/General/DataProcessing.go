package General

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	err  error
	Smap sync.Map
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

func viperInit() {
	// 初始化配置文件
	fileInit(false, viper.GetString("rootPath"))
	log.Printf("Viper Config Path: [%s]", viper.GetString("rootPath")+"/config.yaml")
	log.Printf("Logs Path: [%s]", viper.GetString("logPath"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("rootPath"))
	fileInit(false, viper.GetString("rootPath")+"/config.yaml")
}

func logInit() {
	// 设置日志格式
	// 创建日志文件夹
	fileInit(true, viper.GetString("rootPath")+"/logs")
	log.SetLevel(log.WarnLevel | log.InfoLevel | log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf(" %s:%s ", path.Base(frame.File), strconv.Itoa(frame.Line))
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

func fileInit(fileAttributes bool, fileName string) {
	if _, err = os.Stat(fileName); err != nil {
		if fileAttributes {
			if err = os.MkdirAll(fileName, os.ModePerm); err != nil {
				log.Printf("创建[%s] 目录失败, Error: [%s]", fileName, err)
			} else {
				log.Printf("创建[%s] 目录成功", fileName)
			}
		} else {
			var f *os.File
			if f, err = os.Create(fileName); err != nil {
				log.Printf("创建[%s] 文件失败, Error: [%s]", fileName, err)
			} else {
				log.Printf("创建[%s] 文件成功", fileName)
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					log.Printf("关闭[%s] 文件失败, Error: [%s]", fileName, err)
				}
			}(f)
		}
	}
}

func WechatBotInit() {
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("[viper] 读取配置文件失败, Error: [%s]", err)
	} else {
		log.Printf("[viper] 读取配置文件成功")
	}
	// WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT
	if viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT") == "" {
		viper.Set("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", "false")
	}
	if err = os.Setenv("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT")); err != nil {
		log.Errorf("设置环境变量: [WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT] 失败, Error: [%s]", err)
	}
	// ENDPOINT
	if viper.GetString("wechaty.endpoint") == "" {
		log.Errorf("请填写服务器地址, endpoint: [%s]", viper.GetString("wechaty.endpoint"))
		viper.Set("wechaty.endpoint", "Please Fill In Your Server Address")
		os.Exit(1)
	}
	// TOKEN
	if viper.GetString("wechaty.token") == "" {
		viper.Set("wechaty.token", "Please Fill In Your Token")
		log.Errorf("请填写服务器 Token: [%s]", viper.GetString("wechaty.token"))
		os.Exit(1)
	}
}

func ViperWrite() {
	viper.Set("chat", "")
	if err = viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		log.Errorf("Viper Write file Error: %s", err)
	} else {
		log.Printf("Viper Write file Success")
	}
}
