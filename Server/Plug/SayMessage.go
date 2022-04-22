package Plug

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"path"
	"runtime"
	"strings"
	"time"
	"wechatBot/General"
)

var (
	lastDate time.Time
)

/*
	func NightMode(message.From().ID())
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
		log.Println("[NightMode] 管理员")
		return true
	} else {
		return !(now.Before(timeEnd) && now.After(timeStart))
	}
}

/*
	ChatTimeLimit(message.Date().Format("2006-01-02 15:04:05"))
		: 判断消息是否在规定时间内
		: 如果是，则返回true，否则返回false
*/
func ChatTimeLimit(date string) bool {
	//当前时间
	var (
		now time.Time
		loc *time.Location
	)
	if date == "" {
		return true
	}
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	if loc, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Loc: [%v]", err, loc)
	}
	if now, err = time.ParseInLocation("2006-01-02 15:04:05", timeNow, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Now: [%v]", err, now)
	}
	//当前时间转换为"年-月-日"的格式
	if lastDate, err = time.ParseInLocation("2006-01-02 15:04:05", date, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Lastdate: [%v]", err, lastDate)
		return false
	}
	//计算两个时间相差的秒数
	if second := int(now.Sub(lastDate).Seconds()); second < 30 {
		log.Errorf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		General.Messages.Reply = fmt.Sprintf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		General.Messages.ReplyStatus = true
		General.Messages.AutoInfo = General.Messages.AutoInfo + "[" + General.Messages.Reply + "]"
		return false
	}
	return true
}

/*
	如果有自定义消息内容则填写，没有则为空
	SayMessage(message, "hello word")
	SayMessage(message, "")
	请确保你设置过了 ChatTimeLimit函数
*/
func SayMessage(message *user.Message, msg string) string {
	log.Printf("消息来自函数: [%v]", Copyright(make([]uintptr, 1)))
	if !NightMode(message.From().ID()) { // 夜间模式
		return "当前处于夜间模式!"
	}
	if msg == "" {
		msg = "你想和我说什么呢?"
	}
	// TODO 0.79 私聊有问题
	//if _, err = message.Say(msg); err != nil {
	//	log.Errorf("[SayMessage] [%v], error: %v", msg, err)
	//	return
	//}

	_, _ = message.Say(msg)
	return msg
}

/*
	Copyright(make([]uintptr, 1))
	log.Printf("消息来自: [%v]", Copyright(make([]uintptr, 1)))
	返回上一个函数的名称
*/
func Copyright(pc []uintptr) string {
	runtime.Callers(3, pc)
	return strings.Replace(path.Ext(runtime.FuncForPC(pc[0]).Name()), ".", "", -1)
}
