package data

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	err error
)

func ErrorFormat(str string, err error) {
	_, file, line, _ := runtime.Caller(1)
	prefix := "[xrsec]"
	date := time.Now().Format("15:04:05")
	fmt.Printf("%s %s%s %s:%d %s%s\n", prefix, viper.GetString("failed"), date, path.Base(file), line, str, err)
}

func SuccessFormat(str string) {
	_, file, line, _ := runtime.Caller(1)
	prefix := "[xrsec]"
	date := time.Now().Format("15:04:05")
	fmt.Printf("%s %s%s %s:%d %s\n", prefix, viper.GetString("success"), date, path.Base(file), line, str)
}

func ViperWrite() {
	if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		ErrorFormat("Viper Write file Error: ", err)
	} else {
		SuccessFormat("Viper Write file Success")
	}
}

func ViperRead() {
	if err = viper.ReadInConfig(); err != nil {
		ErrorFormat("Viper Read Config , Try", err)
		if _, err = os.Stat(viper.GetString("rootPath") + "/config.yaml"); err != nil {
			if _, err = os.Stat(viper.GetString("exePath") + "/config.yaml"); err != nil {
				log.Println("配置文件放在当前路劲即可, 注意检测配置是否正确")
			}
			ErrorFormat("config.yaml not found", err)
			viper.Set("wechaty.wechaty_puppet_endpoint", "Please Fill In Your Server Address")
			viper.Set("wechaty.wechaty_puppet_service_token", "Please Fill In Your Token")
			viper.Set("success", "[\033[01;32m✓\033[0m] ")
			viper.Set("failed", "[\033[01;31m✗\033[0m] ")
			viper.Set("info", "[\033[01;33m➜\033[0m] ")
			var f *os.File
			if f, err = os.Create(viper.GetString("exePath") + "/config.yaml"); err != nil {
				ErrorFormat("Create Config File", err)
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					ErrorFormat("Close Config File", err)
				}
			}(f)
		}
	}
}

func SayMsg(message *user.Message, msg string) {
	if _, err := message.Say("@" + message.From().Name() + msg); err != nil {
		ErrorFormat("SayMsg", err)
	}
}
