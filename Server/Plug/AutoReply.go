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
	if viper.GetString(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID())) == "true" {
		return
	}
	if message.Type() != schemas.MessageTypeText {
		return
	}
	if message.Self() {
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Errorf("消息已丢弃，因为它太旧（超过2分钟）")
		return
	}
	if message.Room() != nil {
		if !message.MentionSelf() {
			return
		}
	}
	General.SayMessage(message, "")
}
