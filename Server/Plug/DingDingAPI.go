package Plug

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

/*
	DingMessage("string")
*/
func DingMessage(msg string) {
	type DingBotResult struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	var (
		dingBotResult DingBotResult
		resp          *http.Response
	)
	dingWebHook := viper.GetString("Ding.URL") + viper.GetString("Ding.TOKEN")
	content := fmt.Sprintf(" {\"msgtype\": \"text\",\"text\": {\"content\": \"%v %v\"}}", viper.GetString("Ding.KEYWORD"), msg)
	// 发送请求
	if resp, err = http.Post(dingWebHook, "application/json; charset=utf-8", strings.NewReader(content)); err != nil {
		log.Errorf("[Ding] 机器人请求错误: Error: [%v]", err)
		return
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("[Ding] 关闭请求错误: Error: [%v]", err)
		}
	}(resp.Body)
	if err = json.NewDecoder(resp.Body).Decode(&dingBotResult); err != nil {
		log.Errorf("[Ding] 机器人请求错误: Error: [%v]", err)
		return
	}
	if dingBotResult.Errcode == 0 {
		log.Println("[Ding] 消息发送成功!")
	} else {
		log.Errorf("[Ding] 消息发送失败: Error: [%v]", err)
	}
}
