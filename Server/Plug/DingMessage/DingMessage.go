package DingMessage

import (
	"fmt"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"os"
	"syscall"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"

	"github.com/blinkbean/dingtalk"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnScan(onScan).
		OnLogin(onLogin).
		OnLogout(onLogout).
		OnMessage(onMessage).
		OnError(onError)
	go onHeartbeat()
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.AtMe {
		log.Infof("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Infof("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Infof("Self CoptRight: [%s] CoptRight: [%s]", Copyright(make([]uintptr, 1)), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Infof("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
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
		if err = cli.SendMarkDownMessage(msg, msg); err != nil {
			log.Errorf("DingMessage Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("DingTalk 通知成功! Copyright: [%s]", Copyright(make([]uintptr, 1)))
	} else {
		log.Infof("现在处于夜间模式，请在白天使用 CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
}

func onScan(context *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
	if status.String() == "ScanStatusWaiting" {
		msg := fmt.Sprintf("账号未登录请扫码!\n\n---\n\n[qrCode](https://wechaty.js.org/qrcode/%v)", qrCode)
		DingSend(viper.GetString("Bot.AdminID"), msg)
	} else if status.String() == "ScanStatusCancel" {
		DingSend(viper.GetString("Bot.AdminID"), "扫码取消")
	} else if status.String() == "ScanStatusTimeout" {
		DingSend(viper.GetString("Bot.AdminID"), "扫码超时")
		// 超时退出
		os.Exit(1)
	}
}

func onLogout(context *wechaty.Context, user *user.ContactSelf, reason string) {
	DingSend(viper.GetString("Bot.AdminID"), fmt.Sprintf("%s已登出!\n\n---\n\n**reason**\n\n%v", user.Name(), reason))
}

func onLogin(context *wechaty.Context, user *user.ContactSelf) {
	DingSend(viper.GetString("Bot.AdminID"), fmt.Sprintf("微信机器人: [%v] 已经登录成功了。", user.Name()))
}

/*
	先启用 GatewayDemon
		func main() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			}
		}()
		go GatewayDemon()
		time.Sleep(time.Second * 30)
		// 重试次数 10
		wechatBotDaemon()
	再注册 onHeartbeatGateway by New()
}
*/
func onHeartbeatGateway() {
	c := cron.New()
	if _, err = c.AddFunc("0 * * * *", func() {
		if GatewayCmd == nil {
			return
		}
		fmt.Println("========================onHeartbeat👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		//if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: viper.GetString("BOT.GROUPID")}); roomID == nil {
		if roomID = Bot.Room().Find(viper.GetString("BOT.GROUP")); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			goto end
		}
		// 发送消息
		if _, err = roomID.Say("onHeartbeat: " + time.Now().Format("01-02 15:04:05")); err == nil {
			log.Errorf("onHeartbeat Error")
			return
		}
		// 消息异常啦！
	end:
		log.Errorf("Heartbeat Say Error: [%v]", err)
		DingSend(viper.GetString("Bot.AdminID"), "心跳发送失败 请检查账号状态")

		// Gateway Process Restart
		if err = syscall.Kill(-GatewayCmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Errorf("Gateway Process Kill Error: [%v]", err)
			return
		}
		go GatewayDemon()
		time.Sleep(30 * time.Second)
		log.Infof("Gateway Process Restarting...")
	}); err != nil {
		return
	}
	c.Start()
}

func onHeartbeat() {
	c := cron.New()
	if _, err = c.AddFunc("@hourly", func() {
		fmt.Println("========================onHeartbeat👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: viper.GetString("BOT.GROUP")}); roomID == nil {
			log.Infof("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			goto end
		}
		// 发送消息
		if _, err = roomID.Say("我还活着!"); err != nil {
			log.Errorf("onHeartbeat Say Error: [%v] Copyright: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("Heartbeat Say Success")
		// 消息异常啦！
	end:
		log.Errorf("Heartbeat Say Error: [%v]", err)
		DingSend(viper.GetString("Bot.AdminID"), "心跳发送失败 请检查账号状态")
		DockerRestart()
		time.Sleep(30 * time.Second)
		log.Infof("Gateway Process Restarting...")
	}); err != nil {
		return
	}
	c.Start()
}

func onError(context *wechaty.Context, err error) {
	// TODO 就怕遇到 BUG 一直发送消息
	// DingSend(viper.GetString("Bot.AdminID"), fmt.Sprintf("微信机器人: [%v] 出现错误了。\n\n---\n\n**Error**\n\n%v", viper.GetString("Bot.Name"), err))
}
