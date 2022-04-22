package Plug

import (
	"fmt"
	"strings"
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
	if message.Type() != schemas.MessageTypeText {
		return
	}
	if message.Self() {
		return
	}
	if !ChatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", message.From().ID()))) { // 消息频率限制，可能会存在 map问题
		return
	}
	if General.Messages.ReplyStatus {
		return
	}
	if !message.MentionSelf() {
		return
	}
	if strings.Contains(message.MentionText(), "djs") {
		log.Printf("添加定时提醒成功! 任务详情: %v", "暂无")
		return
	}
	if strings.Contains(message.MentionText(), "fdj") {
		log.Printf("复读机模式, 复读内容: [%v]", message.MentionText())
		SayMessage(message, strings.Replace(message.MentionText(), "fdj ", "", 1))
		return
	}
	if message.Room() == nil {
		if message.Text() == "加群" || message.Text() == "交流群" {
			keys := ""
			for k := range viper.GetStringMap("Group") {
				keys += "『" + k + "』"
			}
			reply := "现有如下交流群, 请问需要加入哪个呢? 请发交流群名字!\n" + keys
			SayMessage(message, reply)
			return
		}
		if strings.Contains(fmt.Sprintf("%v", viper.GetStringMap("Group")), message.Text()) {
			for i, v := range viper.GetStringMap("Group") {
				if strings.Contains(message.Text(), i) && v != "" {
					//	邀请好友进群
					if err = message.GetWechaty().Room().Load(v.(string)).Add(message.From()); err != nil {
						log.Errorf("邀请好友进群失败, Error: [%v]", err)
						return
					}
					log.Println("邀请好友进群成功!")
					SayMessage(message, "已经拉你啦! 等待管理员审核通过呀!")
					return
				}
			}
		}
		SayMessage(message, "有事请留言,可馨一定会给你回复的!") // TODO 是否使用机器人作为智能消息回复（私人）
		return
	}
	if strings.Contains(message.MentionText(), "print") {
		SayMessage(message, strings.Replace(message.MentionText(), "print", "", 1))
		return
	}
}
