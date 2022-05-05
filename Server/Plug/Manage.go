package Plug

import (
	"fmt"
	"strings"
	"time"
	"wechatBot/General"

	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

/*
	Manage()
	普通用户模式
*/
func Manage(message *user.Message) {
	var reply string
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
	if message.MentionSelf() && strings.EqualFold(message.Text(), "所有人") {
		log.Printf("Mention Self All Members Pass, [%v]", message.Talker().Name())
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
	if strings.EqualFold(message.MentionText(), "djs") {
		log.Printf("添加定时提醒成功! 任务详情: %v", "暂无")
		reply = "添加定时提醒成功! 任务详情: 暂无"
	}
	if strings.EqualFold(message.MentionText(), "fdj") {
		log.Printf("复读机模式, 复读内容: [%v]", message.MentionText())
		reply = strings.Replace(message.MentionText(), "fdj ", "", 1)
	}
	if strings.EqualFold(message.MentionText(), "print") {
		reply = strings.Replace(message.MentionText(), "print", "", 1)
	}
	if message.Room() == nil {
		if message.Text() == "加群" || message.Text() == "交流群" {
			keys := ""
			for k := range viper.GetStringMap("Group") {
				keys += "『" + k + "』"
			}
			reply = "现有如下交流群, 请问需要加入哪个呢? 请发交流群名字!\n" + keys
		}
		if strings.EqualFold(fmt.Sprintf("%v", viper.GetStringMap("Group")), message.Text()) {
			for i, v := range viper.GetStringMap("Group") {
				if strings.EqualFold(message.Text(), i) && v != "" {
					//	邀请好友进群
					if err = message.GetWechaty().Room().Load(v.(string)).Add(message.Talker()); err != nil {
						log.Errorf("邀请好友进群失败, Error: [%v]", err)
						return
					}
					log.Println("邀请好友进群成功!")
					reply = "已经拉你啦! 等待管理员审核通过呀!"
				}
			}

		}
	}
	if reply == "" {
		//reply = "有事请留言,可馨一定会给你回复的!" // TODO 是否使用机器人作为智能消息回复（私人）
		return
	}
	SayMessage(message, reply)
}