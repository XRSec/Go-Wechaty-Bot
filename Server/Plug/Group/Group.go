package Group

import (
	"github.com/spf13/viper"
	"strings"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func Group() *wechaty.Plugin {
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
	if !m.AtMe {
		log.Errorf("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
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
	if !m.Status {
		log.Errorf("Status: [%v] CoptRight: [%v]", m.Status, Copyright(make([]uintptr, 1)))
	}
	for _, v := range viper.GetStringMapString("grouppass") {
		if strings.Contains(message.Text(), v) {
			log.Printf("%v Pass, [%v] CoptRight: [%v]", v, message.Talker().Name(), Copyright(make([]uintptr, 1)))
			m.PassResult = v
			m.Pass = true
			context.SetData("msgInfo", m)
			return
		}
	}
}
