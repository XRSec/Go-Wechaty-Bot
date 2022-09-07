package Cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"time"
	. "wechatBot/General"
)

var err error

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	cronTab()
	return plug
}

func cronTab() {
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(nyc))
	if _, err = c.AddFunc("0 1 * * *", func() {
		fmt.Println("========================HackingClub技术交流群16签到👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: "25649532727@chatroom"}); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			return
		}
		// 发送消息
		if _, err = roomID.Say("签到：HackingClub祝大家天天得0day!"); err != nil {
			log.Errorf("HackingClub技术交流群16签到 Say Error [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("HackingClub技术交流群16签到 Say Success")
	}); err != nil {
		log.Errorf("cronTab Error: [%v]", err)
	}
	if _, err = c.AddFunc("0 2 * * *", func() {
		fmt.Println("========================别卷了提醒👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: "24068632533@chatroom"}); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			return
		}
		// 发送消息
		if _, err = roomID.Say("@\u2005Qian别卷了,赶紧睡觉!"); err != nil {
			log.Errorf("别卷了提醒 Say Error [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("别卷了提醒 Say Success")
	}); err != nil {
		log.Errorf("cronTab Error: [%v]", err)
	}
	if _, err = c.AddFunc("0 9 * * *", func() {
		fmt.Println("========================上班打卡提示👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: "25649532727@chatroom"}); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			return
		}
		// 发送消息
		if _, err = roomID.Say("打卡啦~打卡啦!"); err != nil {
			log.Errorf("上班打卡提示 Say Error [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("上班打卡提示 Say Success")
	}); err != nil {
		log.Errorf("cronTab Error: [%v]", err)
	}
	if _, err = c.AddFunc("0 16 * * *", func() {
		fmt.Println("========================下班打卡提示👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: "25649532727@chatroom"}); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			return
		}
		// 发送消息
		if _, err = roomID.Say("打卡啦~打卡啦!"); err != nil {
			log.Errorf("下班打卡提示 Say Error [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("下班打卡提示 Say Success")
	}); err != nil {
		log.Errorf("cronTab Error: [%v]", err)
	}
	if _, err = c.AddFunc("0 22 * * *", func() {
		fmt.Println("========================睡觉提示👇========================")
		var roomID _interface.IRoom
		// 得不到 roomID 所以出问题了
		if roomID = Bot.Room().Find(&schemas.RoomQueryFilter{Id: "24068632533@chatroom"}); roomID == nil {
			log.Errorf("RoomID Find Error CoptRight: [%s]", Copyright(make([]uintptr, 1)))
			return
		}
		// 发送消息
		if _, err = roomID.Say("睡觉啦~睡觉啦!"); err != nil {
			log.Errorf("睡觉提示 Say Error [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("睡觉提示 Say Success")
	}); err != nil {
		log.Errorf("cronTab Error: [%v]", err)
	}
	// 启动定时器
	c.Start()
}
