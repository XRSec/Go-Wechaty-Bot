package General

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"path"
	"runtime"
	"strings"
	"time"
)

//var (
//	err error
//)

type (
	MessageInfo struct {
		Date        string `json:"data"`         // 时间
		Status      bool   `json:"status"`       // 群聊属性
		AtMe        bool   `json:"atme"`         // 是否@我
		RoomName    string `json:"room_name"`    // 群聊名称
		RoomID      string `json:"room_id"`      // 群聊ID
		UserName    string `json:"user_name"`    // 用户名称
		UserID      string `json:"userid"`       // 用户ID
		Content     string `json:"content"`      // 聊天内容
		AutoInfo    string `json:"auto_info"`    // 信息一览
		ReplyResult string `json:"reply_result"` // 自动回复
		Reply       bool   `json:"reply"`        // 自动回复状态
		PassResult  string `json:"pass_result"`  // pass 原因
		Pass        bool   `json:"pass"`         // pass 状态
	}
)

func Pretreatment() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func encodeMessage(context *wechaty.Context, message *user.Message) {
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
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		m.PassResult = "MessageTypeRecalled"
		m.Pass = true
		return
	}
	if message.Type() == schemas.MessageTypeUnknown && message.Talker().Name() == "微信团队" {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		m.PassResult = "微信团队"
		m.Pass = true
		return
	}
	if message.Talker().Type().String() == "ContactTypeOfficial" {
		log.Printf("Official Pass, [%v]", message.Talker().Name())
		m.PassResult = "Official"
		m.Pass = true
		return
	}
	if message.MentionSelf() || message.Room() == nil {
		if strings.Contains(message.Text(), "所有人") {
			m.PassResult = "所有人"
			m.Pass = true
			return
		}
		m.AtMe = true
		fmt.Println("set atme ", m.AtMe)
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

	chatTimeLimit(m)
	//chatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", m.UserID)))
	context.SetData("msgInfo", &m)
}

func onMessage(context *wechaty.Context, message *user.Message) {
	encodeMessage(context, message)
}

/*
	ChatTimeLimit(message.Date().Format("2006-01-02 15:04:05"))
		: 判断消息是否在规定时间内
		: 如果是，则返回true，否则返回false
*/
func chatTimeLimit(m MessageInfo) {
	//当前时间
	var (
		now      time.Time
		loc      *time.Location
		lastDate time.Time
		date     = viper.GetString(fmt.Sprintf("Chat.%v.Date", m.UserID))
	)

	if date == "" {
		return
	}

	//if m.Status && !m.AtMe {
	//	log.Infof("Room !Atme [%v]", m.UserName)
	//	return
	//}
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	if loc, err = time.LoadLocation("Local"); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Loc: [%v]", err, loc)
		// [ChatTimeLimit] time.ParseInLocation, Error: [The system cannot find the path specified.], Loc: [UTC]
		return
	}
	if now, err = time.ParseInLocation("2006-01-02 15:04:05", timeNow, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Now: [%v]", err, now)
		return
	}
	//当前时间转换为"年-月-日"的格式
	if lastDate, err = time.ParseInLocation("2006-01-02 15:04:05", date, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Lastdate: [%v]", err, lastDate)
		return
	}
	//计算两个时间相差的秒数
	if second := int(now.Sub(lastDate).Seconds()); second < 30 {
		log.Errorf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		// Messages.Reply = fmt.Sprintf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		m.PassResult = "ChatTimeLimit"
		m.Pass = true
		return
	}
}

/*
	如果有自定义消息内容则填写，没有则为空
	SayMessage(context, message, "hello word")
	SayMessage(context, message, "")
	请确保你设置过了 ChatTimeLimit函数
*/
func SayMessage(context *wechaty.Context, message *user.Message, msg string) {
	m, ok := (context.GetData("msgInfo")).(*MessageInfo)
	if !ok {
		log.Errorf("Conversion failed")
	}
	log.Printf("消息来自函数: [%v]", Copyright(make([]uintptr, 1)))
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
	//	log.Errorf("[SayMessage] [%v], error: %v", msg, err)
	//	return
	//}

	if _, err = message.Say(msg); err != nil {
		log.Errorf("SayMessage Error: [%v]", err)
	}
	m.ReplyResult = msg
	m.Reply = true
	//m.AutoInfo += fmt.Sprintf(" 回复: [%v]", msg)
	context.SetData("msgInfo", &m)
	viper.Set(fmt.Sprintf("Chat.%v.Date", m.UserID), m.Date)
}

func DingSend(context *wechaty.Context, message *user.Message, msg string) {
	if NightMode(message.Talker().ID()) {
		cli := dingtalk.InitDingTalkWithSecret(viper.GetString("Ding.TOKEN"), viper.GetString("Ding.SECRET"))
		if err := cli.SendMarkDownMessage(msg, msg); err != nil {
			log.Errorf("DingMessage Error: %v", err)
			return
		}
		log.Println("DingTalk 通知成功! Copyright: ", Copyright(make([]uintptr, 1)))
	} else {
		log.Println("现在处于夜间模式，请在白天使用")
		return
	}
}

func DingMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(*MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
	}
	if m.Pass {
		log.Errorf("Pass CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	//if m.Reply {
	//	log.Errorf("Reply CoptRight: [%s]", Copyright(make([]uintptr, 1)))
	//	return
	//}
	if !m.AtMe {
		log.Errorf("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
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
	DingSend(context, message, msg)
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
	if userID == viper.GetString("bot.adminid") {
		log.Println("[NightMode] 管理员 Copyright:", Copyright(make([]uintptr, 1)))
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
