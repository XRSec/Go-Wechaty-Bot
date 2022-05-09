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
	if message.Room() != nil && !message.MentionSelf() { // 不允许私聊使用
		log.Printf("Room Pass, [%v]", message.Talker().Name())
		return
	}

	// All Members Pass
	if message.MentionSelf() && strings.Contains(message.Text(), "所有人") {
		log.Printf("Mention Self All Members Pass, [%v]", message.Talker().Name())
		return
	}

	// if message.Room() == nil {
	// 	return
	// }
	// if !strings.Contains(message.Room().Topic(), "Debug") {
	// 	return
	// }

	// PassStatus
	if General.Messages.PassStatus {
		log.Printf("PassStatus Pass, [%v]", message.Talker().Name())
		return
	}

	// ReplyStatus
	if General.Messages.ReplyStatus {
		log.Printf("ReplyStatus Pass, [%v]", message.Talker().Name())
		return
	}

	// Processing message content
	var msg string
	if message.MentionText() == "" {
		msg = "你想表达什么[破涕为笑]?"
	}

	if msg = WXAPI(message); msg != "" {
		goto labelSay
	}

	if msg = Qingyunke(message.MentionText()); msg != "" {
		goto labelSay
	}

	if msg = Tuling(message.MentionText()); msg == "" {
		msg = "我又不会了?"
	}

labelSay:
	// time.Sleep(5 * time.Second)
	SayMessage(message, msg)
	//General.Messages.ReplyStatus = true
	//General.Messages.AutoInfo = General.Messages.AutoInfo + "[" + General.Messages.Reply + "]"
	//viper.Set(fmt.Sprintf("Chat.%v.Date", message.Talker().ID()), General.Messages.Date)
}
