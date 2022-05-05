package Plug

import (
	"fmt"
	"strings"
	"time"
	"wechatBot/General"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

/*
	AdminManage(message)
	管理员管理
*/
func AdminManage(message *user.Message) {
	// MessageTypeText
	if message.Type() != schemas.MessageTypeText {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	// self
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
	// Admin
	if message.Talker().ID() != viper.GetString("bot.adminid") { // 以下功能仅对管理员开放
		log.Printf("%v is not admin", message.Talker().ID())
		return
	}
	if message.MentionText() == "add" { // 添加好友
		var (
			addUser = message.MentionList()[0]
			member  _interface.IContact
		)
		//if member, err = message.Room().Member(addUserName); err != nil && member != nil {
		//	log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%v], 用户信息: [%v]", addUserName, member.String()), err)
		//}
		//log.Printf("搜索用户名ID成功, 用户名: [%v]", addUser.Name())
		if message.GetWechaty().Contact().Load(addUser.ID()).Friend() {
			log.Printf("用户已经是好友, 用户名: [%v]", addUser.Name())
			SayMessage(message, fmt.Sprintf("用户: [%v] 已经是好友了", addUser))
			return
		}
		if err = message.GetWechaty().Friendship().Add(member, fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("bot.name"))); err != nil {
			log.Errorf("添加好友失败, 用户名: [%v], Error: [%v]", addUser, err)
			SayMessage(message, fmt.Sprintf("添加好友失败, 用户: [%v]", addUser))
			return
		}
		SayMessage(message, fmt.Sprintf("好友申请发送成功, 用户: [%v]", addUser))
		return
	}
	if message.MentionText() == "del" { // 从群聊中移除用户
		var (
			delUser = message.MentionList()[0]
		)
		log.Printf(message.MentionText())
		//if member, err = message.Room().Member(delUser.ID()); err != nil && member != nil {
		//	log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%v], 用户信息: [%v]", delUser.Name(), member.String()), err)
		//	return
		//}
		//log.Printf("搜索用户名ID成功, 用户名: [%v], 用户信息: [%v]", deleteUserName, member.String())
		if err = message.Room().Del(delUser); err != nil {
			log.Errorf("从群聊中移除用户失败, 用户名: [%v] Error: [%v]", delUser.Name(), err)
			SayMessage(message, fmt.Sprintf("从群聊中移除用户失败, 用户: [%v]", delUser.Name()))
			return
		}
		General.Messages.Reply = fmt.Sprintf("从群聊中移除用户: [%v]", delUser.Name())
		General.Messages.ReplyStatus = true
		return
	}
	if message.MentionText() == "quit" { // 退群
		SayMessage(message, "我走了, 拜拜👋🏻, 记得想我哦 [大哭]")
		if err = message.Room().Quit(); err != nil {
			log.Errorf("退出群聊失败, 群聊名称: [%v], Error: [%v]", message.Room().Topic(), err)
			SayMessage(message, fmt.Sprintf("退出群聊失败, 群聊名称: [%v], Error: [%v]", message.Room().Topic(), err))
			return
		}
		General.Messages.ReplyStatus = true
		log.Printf("退出群聊成功! 群聊名称: [%v]", message.Room().Topic())
		return
	}
	if strings.EqualFold(message.MentionText(), "gmz") {
		var (
			newName = strings.Replace(message.MentionText(), "gmz ", "", 1)
		)
		if err = message.GetPuppet().SetContactSelfName(newName); err != nil {
			log.Errorf("修改用户名失败, Error: [%v]", err)
			SayMessage(message, fmt.Sprintf("修改用户名失败, Error: [%v]", err))
			return
		}
		log.Printf("修改用户名成功! 新的名称: %v", newName)
		General.Messages.ReplyStatus = true
		General.Messages.Reply = fmt.Sprintf("改名字: [%v]", newName)
		return
	}
}
