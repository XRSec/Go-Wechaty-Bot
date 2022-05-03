package Plug

import (
	"github.com/blinkbean/dingtalk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
	DingMessage("string",message.Talker().ID())
*/
func DingMessage(messages, UserID string) {
	if NightMode(UserID) {
		cli := dingtalk.InitDingTalk([]string{viper.GetString("Ding.TOKEN")}, viper.GetString("Ding.KEYWORD"))
		cli.SendMarkDownMessage("", messages)
	} else {
		log.Println("现在处于夜间模式，请在白天使用")
		return
	}
	// 	dingWebHook := viper.GetString("Ding.URL") + viper.GetString("Ding.TOKEN")
	// 	content := fmt.Sprintf(" {\"msgtype\": \"text\",\"text\": {\"content\": \"%s %s\"}}", viper.GetString("Ding.KEYWORD"), messages.AutoInfo)
	// 	// 发送请求
	// 	if resp, err = http.Post(dingWebHook, "application/json; charset=utf-8", strings.NewReader(content)); err != nil {
	// 		ErrorFormat("机器人请求错误: ", err)
	// 	} else {
	// 		if err = json.NewDecoder(resp.Body).Decode(&dingBotResult); err != nil {
	// 			ErrorFormat("机器人请求错误: ", err)
	// 		} else {
	// 			if dingBotResult.Errcode == 0 {
	// 				SuccessFormat("消息发送成功!")
	// 			} else {
	// 				ErrorFormat("消息发送失败: ", err)
	// 			}
	// 		}
	// 	}
	// 	// 关闭请求
	// 	defer func(Body io.ReadCloser) {
	// 		if err = Body.Close(); err != nil {
	// 			ErrorFormat("关闭请求错误: ", err)
	// 		}
	// 	}(resp.Body)
	// }
}
