package Test

import (
	"fmt"
	. "wechatBot/General"
	. "wechatBot/Plug"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnLogin(onLogin)
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
	fmt.Println("MentionSelf", message.MentionSelf())
	fmt.Println("MentionText", message.MentionText())
	m.Pass = true
	m.PassResult = "您不是管理员，无法使用"
	context.SetData("msgInfo", m)
}

func onLogin(context *wechaty.Context, user *user.ContactSelf) {

}
