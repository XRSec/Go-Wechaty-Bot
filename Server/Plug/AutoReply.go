package Plug

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"time"
	"wechatBot/General"
)

func AutoReply(message *user.Message) {
	if message.Type() != schemas.MessageTypeText {
		return
	}
	if message.Self() {
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Println("消息已丢弃，因为它太旧（超过2分钟）")
		return
	}
	if !General.ChatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", message.From().ID()))) { // 消息频率限制，可能会存在 map问题
		return
	}
	if General.Messages.ReplyStatus {
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Errorf("消息已丢弃，因为它太旧（超过2分钟）")
		return
	}
	//if message.Room() != nil { // TODO 这是干嘛的？
	//	if !message.MentionSelf() {
	//		return
	//	}
	//}
	General.SayMessage(message, "")
}
