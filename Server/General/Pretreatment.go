package General

import (
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

type (
	messageInfo struct {
		Date        string `json:"data"`        // 时间
		Status      bool   `json:"status"`      // 群聊属性
		AtMe        bool   `json:"atme"`        // 是否@我
		RoomName    string `json:"roomname"`    // 群聊名称
		RoomID      string `json:"roomid"`      // 群聊ID
		UserName    string `json:"username"`    // 用户名称
		UserID      string `json:"userid"`      // 用户ID
		Content     string `json:"content"`     // 聊天内容
		AutoInfo    string `json:"autoinfo"`    // 信息一览
		Reply       string `json:"reply"`       // 自动回复
		ReplyStatus bool   `json:"replystatus"` // 自动回复状态
		Pass        string `json:"pass"`        // pass 原因
		PassStatus  bool   `json:"passstatus"`  // pass 状态
	}
)

func Pretreatment() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func encodeMessage(ctx *wechaty.Context, msg *user.Message) {
	var m messageInfo
	m.UserName = msg.Talker().Name()
	m.UserID = msg.Talker().ID()
	if msg.Type() == schemas.MessageTypeRecalled {
		log.Printf("Type Pass, Type: [%v]:[%v]", msg.Type().String(), msg.Talker().Name())
		Messages.Pass = "MessageTypeRecalled"
		Messages.PassStatus = true
		return
	}
	ctx.SetData("msgInfo", &m)
}

func onMessage(context *wechaty.Context, message *user.Message) {
	encodeMessage(context, message)
}
