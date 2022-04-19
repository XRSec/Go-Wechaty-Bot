package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/url"
	. "wechatBot/data"
)

type TulingBotResult struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

var (
	tulingBotResult TulingBotResult
)

func TulingMessage(messages MessageInfo) MessageInfo {
	if NightMode(messages.UserID) {
		log.Println("现在处于夜间模式，请在白天使用")
		return messages
	} else {
		// 发送请求
		tulingWebhook := viper.GetString("Tuling.URL") + viper.GetString("Tuling.TOKEN")
		if resp, err = http.Get(tulingWebhook + url.QueryEscape(messages.Content)); err != nil {
			ErrorFormat("图灵机器人请求错误: ", err)
		} else {
			if err = json.NewDecoder(resp.Body).Decode(&tulingBotResult); err != nil {
				return messages
			} else {
				if tulingBotResult.Code != 100000 {
					log.Println("图灵机器人返回错误: ", tulingBotResult.Text)
					return messages
				} else {
					messages.Reply = tulingBotResult.Text
					log.Printf("图灵机器人 回复信息: %+v", messages.Reply)
					messages.AutoInfo += " 回复: [" + messages.Reply + "]"
					return messages
				}
			}
		}
		// 关闭请求
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Println("Close body error:", err)
			}
		}(resp.Body)
		return messages
	}
}
