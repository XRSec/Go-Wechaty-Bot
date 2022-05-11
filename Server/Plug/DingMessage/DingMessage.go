package DingMessage

import (
	"fmt"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"

	"github.com/blinkbean/dingtalk"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func DingMessage() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnScan(onScan).
		OnLogout(onLogout).
		OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.AtMe {
		log.Errorf("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Errorf("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Errorf("Self CoptRight: [%s] CoptRight: [%s]", Copyright(make([]uintptr, 1)), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Errorf("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}
	msg := fmt.Sprintf("%v@我了\n\n---\n\n### 用户属性\n\n用户名: [%v]\n\n用户ID: [%v]", message.Talker().Name(), message.Talker().Name(), message.Talker().ID())
	if message.Room() != nil {
		msg += fmt.Sprintf("\n\n---\n\n### 群聊属性\n\n群聊名称: [%v]\n\n群聊ID: [%v]", message.Room().Topic(), message.Room().ID())
	}
	msg += fmt.Sprintf("\n\n---\n\n**内容**: [%v]", message.Text())
	if m.Pass {
		msg += fmt.Sprintf("\n\n**Pass**: [%v]", m.PassResult)
	} else if m.Reply {
		msg += fmt.Sprintf("\n\n**回复**: [%v]", m.ReplyResult)
	} else {
		//
	}
	// 到这里的时候基本设置好了一些默认的值了
	DingSend(viper.GetString("Bot.AdminID"), msg)
}

func DingSend(userID, msg string) {
	if NightMode(userID) {
		cli := dingtalk.InitDingTalkWithSecret(viper.GetString("Ding.TOKEN"), viper.GetString("Ding.SECRET"))
		if err := cli.SendMarkDownMessage(msg, msg); err != nil {
			log.Errorf("DingMessage Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("DingTalk 通知成功! Copyright: [%s]", Copyright(make([]uintptr, 1)))
	} else {
		log.Errorf("现在处于夜间模式，请在白天使用 CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
}

func onScan(context *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
	if status.String() == "ScanStatusWaiting" {
		msg := fmt.Sprintf("账号未登录请扫码!\n\n---\n\n[qrCode](https://wechaty.js.org/qrcode/%v)", qrCode)
		DingSend(viper.GetString("Bot.AdminID"), msg)
	} else if status.String() == "ScanStatusScanned" {
		fmt.Printf("[Scan] Status: %v %v\n", status.String(), data)
	} else {
		fmt.Printf("[Scan] Status: %v %v\n", status.String(), data)
	}
	time.Sleep(120 * time.Second)
}

func onLogout(context *wechaty.Context, user *user.ContactSelf, reason string) {
	DingSend(viper.GetString("Bot.AdminID"), fmt.Sprintf("%s已登出!\n\n---\n\n**reason**\n\n%v", user.Name(), reason))
}
