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

/*
	AutoReply()
	自动回复
*/
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
	// 群聊中没有@我则不回复
	if message.Room() != nil {
		if !message.MentionSelf() {
			return
		}
	}
	if General.Messages.ReplyStatus {
		return
	}
	// 注意 聊天速率影响 Dingding 正常推送消息
	if !ChatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", message.From().ID()))) { // 消息频率限制，可能会存在 map问题
		return
	}

	// 处理消息内容
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
	General.Messages.Reply = SayMessage(message, msg)
	General.Messages.ReplyStatus = true
	General.Messages.AutoInfo = General.Messages.AutoInfo + "[" + General.Messages.Reply + "]"
	viper.Set(fmt.Sprintf("Chat.%v.Date", message.From().ID()), General.Messages.Date)
}
