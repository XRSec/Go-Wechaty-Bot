package General

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
	. "wechatBot/Plug"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

//var (
//	err error
//)

func Pretreatment() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	var m MessageInfo
	m.UserName = message.Talker().Name()
	m.UserID = message.Talker().ID()
	m.Date = message.Date().Format("2006-01-02 15:04:05")
	m.Reply = false
	m.Pass = false
	m.Status = false
	m.AtMe = false
	// 公众号消息
	if message.Type() == schemas.MessageTypeRecalled {
		log.Errorf("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		m.PassResult = "MessageTypeRecalled"
		m.Pass = true
		return
	}
	if message.Type() == schemas.MessageTypeUnknown && message.Talker().Name() == "微信团队" {
		log.Errorf("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		m.PassResult = "微信团队"
		m.Pass = true
		return
	}
	if message.Talker().Type().String() == "ContactTypeOfficial" {
		log.Errorf("Official Pass, [%v] CoptRight: [%s]", message.Talker().Name(), Copyright(make([]uintptr, 1)))
		m.PassResult = "Official"
		m.Pass = true
		return
	}
	if message.MentionSelf() || message.Room() == nil {
		if strings.Contains(message.Text(), "所有人") {
			m.PassResult = "所有人"
			m.Pass = true
			log.Errorf("At All Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
			return
		}
		m.AtMe = true
	}
	if message.Type() != schemas.MessageTypeText {
		m.Content = "[未知消息类型: " + message.Type().String() + "] " + message.Text()
	} else {
		m.Content = message.Text()
	}
	m.AutoInfo = fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v]", m.UserID, m.UserName, m.Content)
	//m.AutoInfo = fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v]", m.UserID, m.UserName, strings.Replace(m.Content, "\u2005", " ", -1))
	if message.Room() != nil {
		m.RoomName = message.Room().Topic()
		m.RoomID = message.Room().ID()
		m.Status = true
		m.AutoInfo = fmt.Sprintf("群聊ID: [%v] 群聊名称: [%v] %v", m.RoomID, m.RoomName, m.AutoInfo)
	}
	context.SetData("msgInfo", m)
	chatTimeLimit(context, message)
	//chatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", m.UserID)))
}

/*
	ChatTimeLimit(message.Date().Format("2006-01-02 15:04:05"))
		: 判断消息是否在规定时间内
		: 如果是，则返回true，否则返回false
*/
func chatTimeLimit(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	//当前时间
	var (
		now      time.Time
		loc      *time.Location
		lastDate time.Time
		date     = viper.GetString(fmt.Sprintf("Chat.%v.Date", m.UserID))
	)
	if date == "" {
		log.Errorf("Not Set Date, CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	if loc, err = time.LoadLocation("Local"); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Loc: [%v] CoptRight: [%s]", err, loc, Copyright(make([]uintptr, 1)))
		// [ChatTimeLimit] time.ParseInLocation, Error: [The system cannot find the path specified.], Loc: [UTC]
		return
	}
	if now, err = time.ParseInLocation("2006-01-02 15:04:05", timeNow, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Now: [%v] CoptRight: [%s]", err, now, Copyright(make([]uintptr, 1)))
		return
	}
	//当前时间转换为"年-月-日"的格式
	if lastDate, err = time.ParseInLocation("2006-01-02 15:04:05", date, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Lastdate: [%v] CoptRight: [%s]", err, lastDate, Copyright(make([]uintptr, 1)))
		return
	}
	//计算两个时间相差的秒数
	if second := int(now.Sub(lastDate).Seconds()); second < 30 {
		log.Errorf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d] CoptRight: [%s]", lastDate, now, second, Copyright(make([]uintptr, 1)))
		// Messages.Reply = fmt.Sprintf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		m.PassResult = "ChatTimeLimit"
		m.Pass = true
		context.SetData("msgInfo", m)
		return
	} else {
		log.Error("[ChatTimeLimit] 时间相差超过 30 秒")
	}
}

/*
	如果有自定义消息内容则填写，没有则为空
	SayMessage(context, message, "hello word")
	SayMessage(context, message, "")
	请确保你设置过了 ChatTimeLimit函数
*/
func SayMessage(context *wechaty.Context, message *user.Message, msg string) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !NightMode(m.UserID) { // 夜间模式
		return
	}

	//if m.Reply {
	//	return
	//}
	//if m.Pass {
	//	return
	//}

	if msg == "" {
		msg = "你想和我说什么呢?"
	}

	// TODO 0.79 私聊有问题
	//if _, err = message.Say(msg); err != nil {
	//	log.Errorf("[SayMessage] [%v], error: %v CoptRight: [%s]", msg, err, Copyright(make([]uintptr, 1)))
	//	return
	//}

	if _, err = message.Say(msg); err != nil {
		log.Errorf("SayMessage Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	}
	m.ReplyResult = msg
	m.Reply = true
	//m.AutoInfo += fmt.Sprintf(" 回复: [%v]", msg)
	context.SetData("msgInfo", m)
	viper.Set(fmt.Sprintf("Chat.%v.Date", m.UserID), m.Date)
}

/*
	func NightMode(message.Talker().ID())
	: 管理员 返回 true
	: 凌晨返回 false
*/
func NightMode(userID string) bool {
	//当前时间
	startTimeStr := "00:00:00"
	endTimeStr := "06:00:00"
	now := time.Now()
	//当前时间转换为"年-月-日"的格式
	format := now.Format("2006-01-02")
	//转换为time类型需要的格式
	layout := "2006-01-02 15:04:05"
	//将开始时间拼接“年-月-日 ”转换为time类型
	timeStart, _ := time.ParseInLocation(layout, format+" "+startTimeStr, time.Local)
	//将结束时间拼接“年-月-日 ”转换为time类型
	timeEnd, _ := time.ParseInLocation(layout, format+" "+endTimeStr, time.Local)
	//使用time的Before和After方法，判断当前时间是否在参数的时间范围
	if userID == viper.GetString("Bot.AdminId") {
		log.Infof("[NightMode] 管理员 Copyright: [%s]", Copyright(make([]uintptr, 1)))
		return true
	} else {
		return !(now.Before(timeEnd) && now.After(timeStart))
	}
}

/*
	use: Copyright(make([]uintptr, 1))
	return: main.xxx.xxx.xxx
*/
func Copyright(pc []uintptr) string {
	s := ""
	for i := 1; i < 4; i++ {
		runtime.Callers(i, pc)
		if i == 3 {
			s += strings.Replace(path.Ext(runtime.FuncForPC(pc[0]).Name()), ".", "", -1)
		} else {
			s += strings.Replace(path.Ext(runtime.FuncForPC(pc[0]).Name()), ".", "", -1) + " 》"
		}
	}
	return s
}
