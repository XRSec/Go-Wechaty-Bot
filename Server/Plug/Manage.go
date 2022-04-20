package Plug

import (
	"fmt"
	"strings"
	"wechatBot/General"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func Manage(message *user.Message) {
	// map 加锁
	if viper.GetString(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID())) == "true" {
		return
	}
	if strings.Contains(message.MentionText(), "djs") {
		log.Printf("添加定时提醒成功! 任务详情: %s", "暂无")
		return
	}
	if strings.Contains(message.MentionText(), "fdj") {
		log.Printf("复读机模式, 复读内容: [%s]", message.MentionText())
		General.SayMessage(message, strings.Replace(message.MentionText(), "fdj ", "", 1))
		return
	}
	if strings.Contains(message.Text(), "加群") || strings.Contains(message.Text(), "交流群") {
		keys := ""
		for k := range viper.GetStringMap("Group") {
			keys += "『" + k + "』"
		}
		reply := "现有如下交流群, 请问需要加入哪个呢? 请发交流群名字!\n" + keys
		General.SayMessage(message, reply)
		return
	}
	if strings.Contains(fmt.Sprintf("%s", viper.GetStringMap("Group")), message.Text()) {
		for i, v := range viper.GetStringMap("Group") {
			if strings.Contains(message.Text(), i) && v != "" {
				//	邀请好友进群
				if err = message.GetWechaty().Room().Load(v.(string)).Add(message.From()); err != nil {
					log.Errorf("邀请好友进群失败, Error: [%s]", err)
					return
				}
				log.Println("邀请好友进群成功!")
				General.SayMessage(message, "已经拉你啦! 等待管理员审核通过呀!")
				return
			}
		}
	}
}
