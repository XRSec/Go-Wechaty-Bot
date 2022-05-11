package test

import (
	. "wechatBot/General"
	. "wechatBot/Plug"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func Test() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	return
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Talker().ID() != viper.GetString("Bot.AdminID") {
		m.Pass = true
		m.PassResult = "您不是管理员，无法使用"
		context.SetData("msgInfo", m)
	}
}
