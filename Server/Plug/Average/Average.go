package Average

import (
	"fmt"
	"strings"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func Average() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	reply := ""
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Pass {
		log.Errorf("Pass CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Reply {
		log.Errorf("Reply CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.AtMe {
		log.Errorf("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Errorf("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Infof("Self CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Errorf("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}

	if strings.Contains(message.MentionText(), "djs") {
		log.Infof("添加定时提醒成功! 任务详情: %v CoptRight: [%s]", "暂无", Copyright(make([]uintptr, 1)))
		reply = "添加定时提醒成功! 任务详情: 暂无"
	}

	if strings.Contains(message.MentionText(), "fdj") {
		log.Infof("复读机模式, 复读内容: [%v] CoptRight: [%s]", message.MentionText(), Copyright(make([]uintptr, 1)))
		reply = strings.Replace(message.MentionText(), "fdj ", "", 1)
	}

	if strings.Contains(message.MentionText(), "print") {
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
		if strings.Contains(fmt.Sprintf("%v", viper.GetStringMap("Group")), message.Text()) {
			for i, v := range viper.GetStringMap("Group") {
				if strings.Contains(message.Text(), i) && v != "" {
					//	邀请好友进群
					if err = message.GetWechaty().Room().Load(v.(string)).Add(message.Talker()); err != nil {
						log.Errorf("邀请好友进群失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
						return
					}
					log.Infof("邀请好友进群成功! CoptRight: [%s]", Copyright(make([]uintptr, 1)))
					reply = "已经拉你啦! 等待管理员审核通过呀!"
				}
			}
		}
	}
	if reply == "" {
		return
	}
	SayMessage(context, message, reply)
}
