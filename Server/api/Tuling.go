package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	. "wechatBot/data"
)

type TulingBotResult struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

var (
	tulingBotResult TulingBotResult
)

func TulingMessage(msg MessageInfo) MessageInfo {
	if NightMode() {
		log.Println("现在处于夜间模式，请在白天使用")
		return msg
	} else {
		// 发送请求
		tulingWebhook := viper.GetString("Tuling.URL") + viper.GetString("Tuling.TOKEN")
		if resp, err = http.Get(tulingWebhook + msg.Content); err != nil {
			log.Printf("%s机器人请求错误: %s\n", viper.GetString("faild"), err)
			ErrorFormat("图灵机器人请求错误: ", err)
		} else {
			if err = json.NewDecoder(resp.Body).Decode(&tulingBotResult); err != nil {
				return msg
			} else {
				if tulingBotResult.Code != 100000 {
					return msg
				} else {
					msg.Reply = tulingBotResult.Text
					log.Printf("图灵机器人 回复信息: %+v", msg.Reply)
					msg.AutoInfo += " 回复: [" + msg.Reply + "]"
					return msg
				}
			}
		}
		// 关闭请求
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Println("Close body error:", err)
			}
		}(resp.Body)
		return msg
	}
}
