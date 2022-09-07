package General

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

func GatewayStart() error {
	var (
		cmd      *exec.Cmd
		stderr   bytes.Buffer
		RootPath = viper.GetString("RootPath")
	)

	if _, err = os.Stat(viper.GetString("RootPath") + "/package.json"); err != nil {
		if err = GetGatewayConfig("https://ghproxy.com/https://raw.githubusercontent.com/XRSec/Go-Wechaty-Bot/main/padlocal/package.json", viper.GetString("RootPath")+"/package.json"); err != nil {
			return err
		}
	}

	log.Infof("[gateway] 发现 Node Gateway 配置信息, 存在启动 Node Gateway")

	// 安装依赖
	if _, err = os.Stat(RootPath + "/node_modules/.bin/wechaty"); err != nil {
		cmd = exec.Command("npm", "install", "--registry=https://registry.npmmirror.com")
		cmd.Dir = RootPath
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	log.Infof("[gateway] 安装依赖成功")

	// 启动 Node Gateway
	cmd = exec.Command("npm", "run", "serve")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = RootPath
	cmd.Env = os.Environ()
	// 设置环境变量
	cmd.Env = append(cmd.Env, "WECHATY_PUPPET=wechaty-puppet-padlocal",
		"WECHATY_TOKEN="+viper.GetString("WECHATY.WECHATY_TOKEN"),
		"WECHATY_PUPPET_PADLOCAL_TOKEN="+viper.GetString("WECHATY.WECHATY_PUPPET_PADLOCAL_TOKEN"),
		"WECHATY_PUPPET_SERVER_PORT="+viper.GetString("WECHATY.WECHATY_PUPPET_SERVER_PORT"),
		"WECHATY_LOG="+viper.GetString("WECHATY.WECHATY_LOG"),
		"WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT="+viper.GetString("WECHATY.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT"),
		"WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER="+viper.GetString("WECHATY.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER"),
	)
	_ = os.Setenv("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER", viper.GetString("WECHATY.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER"))
	_ = os.Setenv("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT", viper.GetString("WECHATY.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT"))
	//cmd.Stdout = io.MultiWriter(os.Stdout, &lumberjack.Logger{
	cmd.Stdout = io.MultiWriter(&lumberjack.Logger{
		Filename:   viper.Get("LogPath").(string) + "gateway.log", //日志文件位置
		MaxSize:    50,                                            // 单文件最大容量,单位是MB
		MaxBackups: 1,                                             // 最大保留过期文件个数
		MaxAge:     365,                                           // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,                                          // 是否需要压缩滚动日志, 使用的 gzip 压缩
	})

	cmd.Stderr = &stderr // 标准错误输出
	err = cmd.Start()

	log.Errorf(string(stderr.Bytes()))
	log.Infof("[gateway] 启动 Node Gateway 成功,ID: %v", cmd.Process.Pid)
	GatewayCmd = cmd
	go func() {
		if err = GatewayCmd.Wait(); err != nil {
			fmt.Printf("Child process %d exit with err: %v\n", GatewayCmd.Process.Pid, err)
		}
	}()
	return err
}

func GetGatewayConfig(configUrl string, file string) (err error) {
	var (
		resp *http.Response
		data []byte
	)

	if configUrl == "" {
		return errors.New("[gateway] 未配置 Node Gateway 的 NPM 下载地址")
	}

	if resp, err = http.Get(configUrl); err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	if err = ioutil.WriteFile(file, data, 0644); err != nil {
		return err
	}
	return err
}

func GatewayDemon() {
	// 这里就没有用for 循环，必须先正常连接一次 否则 server 循环10次连接不上直接退出
	if err = GatewayStart(); err != nil {
		log.Errorf("[gateway] Node Gateway 启动失败, Error: [%v]", err)
	}

}
