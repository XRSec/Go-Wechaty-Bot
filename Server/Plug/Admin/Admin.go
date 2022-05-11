package Admin

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
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

/*
	Admin()
	管理员
*/
func Admin() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
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
	//if !m.AtMe {
	//	log.Errorf("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
	//	return
	//}
	if message.Type() != schemas.MessageTypeText {
		log.Errorf("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Errorf("Self CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Errorf("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}
	if m.UserID != viper.GetString("Bot.AdminID") {
		log.Errorf("UserID: [%s] CoptRight: [%s]", m.UserID, Copyright(make([]uintptr, 1)))
		return
	}
	if message.MentionText() == "add" || message.MentionText() == "加" { // 添加好友
		var (
			addUser = message.MentionList()[0]
			member  _interface.IContact
		)
		//if member, err = message.Room().Member(addUserName); err != nil && member != nil {
		//	log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%v], 用户信息: [%v]", addUserName, member.String()), err)
		//}
		//log.Printf("搜索用户名ID成功, 用户名: [%v]", addUser.Name())
		if message.GetWechaty().Contact().Load(addUser.ID()).Friend() {
			log.Infof("用户已经是好友, 用户名: [%v] CoptRight: [%s]", addUser.Name(), Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("用户: [%v] 已经是好友了", addUser))
			return
		}

		if err = message.GetWechaty().Friendship().Add(member, fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("Bot.Name"))); err != nil {
			log.Errorf("添加好友失败, 用户名: [%v], Error: [%v] CoptRight: [%s]", addUser, err, Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("添加好友失败, 用户: [%v]", addUser))
			return
		}

		SayMessage(context, message, fmt.Sprintf("好友申请发送成功, 用户: [%v]", addUser))
		return
	}

	if message.MentionText() == "del" || message.MentionText() == "踢" { // 从群聊中移除用户
		var (
			delUser = message.MentionList()[0]
		)
		if err = message.Room().Del(delUser); err != nil {
			log.Errorf("从群聊中移除用户失败, 用户名: [%v] Error: [%v] CoptRight: [%s]", delUser.Name(), err, Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("从群聊中移除用户失败, 用户: [%v]", delUser.Name()))
			return
		}
		m.PassResult = fmt.Sprintf("从群聊中移除用户: [%v]", delUser.Name())
		m.Pass = true
		context.SetData("msgInfo", m)
		return
	}

	if message.MentionText() == "quit" || message.MentionText() == "退" { // 退群
		SayMessage(context, message, "我走了, 拜拜👋🏻, 记得想我哦 [大哭]")
		if err = message.Room().Quit(); err != nil {
			log.Errorf("退出群聊失败, 群聊名称: [%v], Error: [%v] CoptRight: [%s]", message.Room().Topic(), err, Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("退出群聊失败, 群聊名称: [%v], Error: [%v]", message.Room().Topic(), err))
			return
		}

		m.PassResult = fmt.Sprintf("退出群聊成功! 群聊名称: [%v]", message.Room().Topic())
		m.Pass = true
		return
	}

	if strings.Contains(message.MentionText(), "gmz") {
		var (
			newName = strings.Replace(message.MentionText(), "gmz ", "", 1)
		)

		if err = message.GetPuppet().SetContactSelfName(newName); err != nil {
			log.Errorf("修改用户名失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("修改用户名失败, Error: [%v]", err))
			return
		}

		log.Infof("修改用户名成功! 新的名称: %v CoptRight: [%s]", newName, Copyright(make([]uintptr, 1)))
		m.PassResult = fmt.Sprintf("改名字: [%v]", newName)
		m.Pass = true
		context.SetData("msgInfo", m)
		return
	}
}
