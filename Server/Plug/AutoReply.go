package Plug

import (
	"strings"
	"time"
	"wechatBot/General"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

/*
	AutoReply()
	自动回复
*/
func AutoReply(message *user.Message) {
	if message.Type() != schemas.MessageTypeText {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	if message.Self() {
		log.Printf("Self Pass, [%v]", message.Talker().Name())
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Println("消息已丢弃，因为它太旧（超过2分钟）")
		return
	}
	// If there is no @me in the group chat, I will not reply
	if message.Room() != nil && !message.MentionSelf() { // 不允许私聊使用
		log.Printf("Room Pass, [%v]", message.Talker().Name())
		return
	}
	// All Members Pass
	if message.MentionSelf() && strings.Contains(message.Text(), "全体成员") {
		log.Printf("Mention Self All Members Pass, [%v]", message.Talker().Name())
		return
	}
	if General.Messages.ReplyStatus {
		log.Printf("ReplyStatus Pass, [%v]", message.Talker().Name())
		return
	}
	// Processing message content
	var msg string
	if message.MentionText() == "" {
		msg = "你想和我说什么呢?"
	}
	if msg = WXAPI(message); msg != "" {
		goto labelSay
	}
	if msg = Tuling(message.MentionText()); msg == "" {
		msg = "你想和我说什么呢?"
	}
labelSay:
	SayMessage(message, msg)
	//General.Messages.ReplyStatus = true
	//General.Messages.AutoInfo = General.Messages.AutoInfo + "[" + General.Messages.Reply + "]"
	//viper.Set(fmt.Sprintf("Chat.%v.Date", message.Talker().ID()), General.Messages.Date)
}
