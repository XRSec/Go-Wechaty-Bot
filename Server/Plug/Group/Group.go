package Group

import (
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"strings"
	"time"
	. "wechatBot/General"
)

var (
	err error
)

func Group() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(*MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed")
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
	if !m.Status {
		log.Errorf("Status: [%v] CoptRight: [%v]", m.Status, Copyright(make([]uintptr, 1)))
	}
	if strings.Contains(m.RoomName, "nft") || strings.Contains(m.RoomName, "NFT") {
		log.Printf("NFT Pass, [%v] CoptRight: [%v]", message.Talker().Name(), Copyright(make([]uintptr, 1)))
		m.PassResult = "NFT"
		m.Pass = true
		context.SetData("msgInfo", &m)
		return
	}
}
