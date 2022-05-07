package Plug

import (
	"fmt"
	"strings"
	"time"
	"wechatBot/General"

	"github.com/blinkbean/dingtalk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

/*
	DingMessage("string",message.Talker().ID())
*/
func DingMessageSend(messages, UserID string) {
	if NightMode(UserID) {
		cli := dingtalk.InitDingTalkWithSecret(viper.GetString("Ding.TOKEN"), viper.GetString("Ding.SECRET"))
		if err := cli.SendMarkDownMessage(messages, messages); err != nil {
			log.Errorf("DingMessage Error: %v", err)
			return
		}
		log.Println("DingTalk 通知成功! Copyright: ", Copyright(make([]uintptr, 1)))
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

func DingMessage(message *user.Message) {
	// MessageTypeText
	if message.Type() != schemas.MessageTypeText {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	// Self
	if message.Self() {
		log.Printf("Self Pass, [%v]", message.Talker().Name())
		return
	}
	// TIMEOUT
	if message.Age() > 2*60*time.Second {
		log.Println("消息已丢弃，因为它太旧（超过2分钟）")
		return
	}
	// If there is no @me in the group chat, I will not reply
	if message.Room() != nil && !message.MentionSelf() {
		log.Printf("Room Pass, [%v]", message.Talker().Name())
		return
	}
	// All Members Pass
	if message.MentionSelf() && strings.Contains(message.Text(), "所有人") {
		log.Printf("Mention Self All Members Pass, [%v]", message.Talker().Name())
		return
	}
	// PassStatus
	if General.Messages.PassStatus {
		log.Printf("PassStatus Pass, [%v]", message.Talker().Name())
		return
	}
	msg := fmt.Sprintf("%v@我了\n\n---\n\n### 用户属性\n\n用户名: [%v]\n\n用户ID: [%v]", General.Messages.UserName, General.Messages.UserName, General.Messages.UserID)
	if General.Messages.Status {
		msg += fmt.Sprintf("\n\n---\n\n### 群聊属性\n\n群聊名称: [%v]\n\n群聊ID: [%v]", General.Messages.RoomName, General.Messages.RoomID)
	}
	msg += fmt.Sprintf("\n\n---\n\n**内容**: [%v]", General.Messages.Content)
	if General.Messages.PassStatus {
		msg += fmt.Sprintf("\n\n**Pass**: [%v]", General.Messages.Pass)
	} else if General.Messages.ReplyStatus {
		msg += fmt.Sprintf("\n\n**回复**: [%v]", General.Messages.Reply)
	} else {
		//
	}
	// 到这里的时候基本设置好了一些默认的值了
	DingMessageSend(msg, General.Messages.UserID)
}
