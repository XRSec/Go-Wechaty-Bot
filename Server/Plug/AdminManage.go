package Plug

import (
	"fmt"
	"strings"
	"wechatBot/General"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func AdminManage(message *user.Message) {
	// map 加锁
	if viper.GetString(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID())) == "true" {
		return
	}
	if message.Room() == nil {
		return
	}
	if !message.MentionSelf() {
		return
	}
	if message.From().ID() != viper.GetString("bot.adminid") {
		log.Printf("%s is not admin", message.From().ID())
		return
	}
	if message.MentionText() == "add" { // 添加好友
		var (
			addUserName = strings.Replace(strings.Replace(message.Text(), "\u2005", "", -1), fmt.Sprintf("@%sadd @", viper.GetString("bot.name")), "", 1) // 过滤用户名
			member      _interface.IContact
		)
		if member, err = message.Room().Member(addUserName); err != nil && member != nil {
			log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%s], 用户信息: [%s]", addUserName, member.String()), err)
		}
		log.Printf("搜索用户名ID成功, 用户名: [%s], 用户信息: [%s]", addUserName, member.String())
		if message.GetWechaty().Contact().Load(member.ID()).Friend() {
			log.Printf("用户已经是好友, 用户名: [%s], 用户信息: [%s]", addUserName, member.String())
			General.SayMessage(message, fmt.Sprintf("用户: [%s] 已经是好友了", addUserName))
			return
		}
		if err = message.GetWechaty().Friendship().Add(member, fmt.Sprintf("你好,我是%s,以后请多多关照!", viper.GetString("bot.name"))); err != nil {
			log.Errorf("添加好友失败, 用户名: [%s], 用户信息: [%s], Error: [%v]", addUserName, member.String(), err)
			General.SayMessage(message, fmt.Sprintf("添加好友失败, 用户: [%s]", addUserName))
			return
		}
		General.SayMessage(message, fmt.Sprintf("好友申请发送成功, 用户: [%s]", addUserName))
		return
	}
	if message.MentionText() == "del" { // 从群聊中移除用户
		var (
			deleteUserName = strings.Replace(strings.Replace(message.Text(), "\u2005", "", -1), fmt.Sprintf("@%sdel @", viper.GetString("bot.name")), "", 1) // 过滤用户名
			member         _interface.IContact
		)
		if member, err = message.Room().Member(deleteUserName); err != nil && member != nil {
			log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%s], 用户信息: [%s]", deleteUserName, member.String()), err)
			return
		}
		log.Printf("搜索用户名ID成功, 用户名: [%s], 用户信息: [%s]", deleteUserName, member.String())
		if err = message.Room().Del(member); err != nil {
			log.Errorf("从群聊中移除用户失败, 用户名: [%s], 用户信息: [%s], Error: [%v]", deleteUserName, member.String(), err)
			General.SayMessage(message, fmt.Sprintf("从群聊中移除用户失败, 用户: [%s]", deleteUserName))
			return
		}
		viper.Set(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID()), true)
		return
	}
	if message.MentionText() == "quit" { // 退群
		General.SayMessage(message, "我走了, 拜拜👋🏻, 记得想我哦 [大哭]")
		if err = message.Room().Quit(); err != nil {
			log.Errorf("退出群聊失败, 群聊名称: [%s], Error: [%v]", message.Room().Topic(), err)
			General.SayMessage(message, fmt.Sprintf("退出群聊失败, 群聊名称: [%s], Error: [%v]", message.Room().Topic(), err))
			return
		}
		viper.Set(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID()), true)
		log.Printf("退出群聊成功! 群聊名称: [%s]", message.Room().Topic())
		return
	}
	if strings.Contains(message.MentionText(), "gmz") {
		var (
			newName = strings.Replace(message.MentionText(), "gmz ", "", 1)
		)
		if err = message.GetPuppet().SetContactSelfName(newName); err != nil {
			log.Errorf("修改用户名失败, Error: [%v]", err)
			General.SayMessage(message, fmt.Sprintf("修改用户名失败, Error: [%v]", err))
			return
		}
		log.Printf("修改用户名成功! 新的名称: %s", newName)
		viper.Set(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID()), true)
		return
	}
}
