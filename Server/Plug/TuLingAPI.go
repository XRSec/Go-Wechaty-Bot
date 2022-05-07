package Plug

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
	TulingMessage("msg")
	: 图灵机器人
	: 发送消息到图灵机器人
	: 并返回 回复
*/
func Tuling(msg string) string {
	type (
		TulingBotResult struct {
			Code int    `json:"code"`
			Text string `json:"text"`
		}
	)
	var (
		tulingBotResult TulingBotResult
		resp            *http.Response
	)
	// 发送请求
	tulingWebhook := viper.GetString("Tuling.URL") + viper.GetString("Tuling.TOKEN")
	if resp, err = http.Get(tulingWebhook + url.QueryEscape(msg)); err != nil {
		log.Errorf("[图灵] 机器人请求错误: [%v]", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println("[图灵] Close body error:", err)
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&tulingBotResult); err != nil {
		log.Errorf("[图灵] 机器人解析错误: [%v]", err)
		return ""
	}
	// 判断响应码
	// TODO 添加自动更换TOKEN
	if tulingBotResult.Code != 100000 {
		if strings.Contains(tulingBotResult.Text, "当天请求次数已用完") {
			if token2 := viper.GetString("tuling.token2"); token2 != "" {
				viper.Set("tuling.token2", viper.GetString("tuling.token"))
				viper.Set("tuling.token", token2)
				return ""
			}
		}
		if tulingBotResult.Text != "你想和我说什么呢?" {
			return ""
		}
		log.Errorf("[图灵] 机器人返回错误: [%v]", tulingBotResult.Text)
	}
	// 输出结果
	log.Printf("[图灵] 机器人 回复信息: %+v", tulingBotResult.Text)
	return tulingBotResult.Text
}
