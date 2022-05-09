package Plug

import (
	"strings"
	"time"
	"wechatBot/General"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func GroupPass(message *user.Message) {
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
	if message.Room() == nil {
		log.Printf("Room Pass, [%v]", message.Talker().Name())
		return
	}

	if !message.MentionSelf() {
		log.Printf("MentionSelf Pass, [%v]", message.Talker().Name())
		return
	}

	// All Members Pass
	if strings.Contains(message.Text(), "所有人") {
		log.Printf("MentionSelf All Pass, [%v]", message.Talker().Name())
		return
	}

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

	if strings.Contains(message.Room().Topic(), "nft") || strings.Contains(message.Room().Topic(), "NFT") {
		log.Printf("NFT Pass, [%v]", message.Talker().Name())
		General.Messages.Pass = "NFT"
		General.Messages.PassStatus = true
		return
	}
}
