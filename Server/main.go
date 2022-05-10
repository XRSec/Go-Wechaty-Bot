package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	. "wechatBot/General"

	. "wechatBot/Plug/Admin"
	. "wechatBot/Plug/AutoReply"
	. "wechatBot/Plug/Average"
	. "wechatBot/Plug/DingMessage"
	. "wechatBot/Plug/ExportMessage"
	. "wechatBot/Plug/FileBox"
	. "wechatBot/Plug/Group"

	"github.com/mdp/qrterminal/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	. "github.com/wechaty/go-wechaty/wechaty"
	puppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	//. "wechatBot/Plug/Test"
)

var (
	err error
)

func onScan(context *Context, qrCode string, status schemas.ScanStatus, data string) {
	fmt.Printf("\n\n")
	if status.String() == "ScanStatusWaiting" {
		qrterminal.GenerateWithConfig(qrCode, qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: " \u2005",
			WhiteChar: "💵",
			QuietZone: 1,
		})
		fmt.Printf("\n\n")
		log.Printf("%v[Scan] https://wechaty.js.org/qrcode/%v %v", viper.GetString("info"), qrCode, data)
	} else if status.String() == "ScanStatusScanned" {
		fmt.Printf("%v[Scan] Status: %v %v\n", viper.GetString("info"), status.String(), data)
	} else {
		fmt.Printf("%v[Scan] Status: %v %v\n", viper.GetString("info"), status.String(), data)
	}
}

/*
	@method onLogin 当机器人成功登陆后，会触发事件，并会在事件中传递当前登陆机器人的信息
	@param {*} user
*/
func onLogin(context *Context, user *user.ContactSelf) {
	fmt.Printf(`%v
                           //
               \\         //
                \\       //
        ## DDDDDDDDDDDDDDDDDDDD ##
        ## DDDDDDDDDDDDDDDDDDDD ##
        ## hh                hh ##      ##         ## ## ## ##   ## ## ## ###   ##    ####     ##
        ## hh    //    \\    hh ##      ##         ##       ##   ##             ##    ## ##    ##
        ## hh   //      \\   hh ##      ##         ##       ##   ##             ##    ##   ##  ##
        ## hh                hh ##      ##         ##       ##   ##     ##      ##    ##    ## ##
        ## hh      wwww      hh ##      ##         ##       ##   ##       ##    ##    ##     ####
        ## hh                hh ##      ## ## ##   ## ## ## ##   ## ## ## ###   ##    ##      ###
        ## MMMMMMMMMMMMMMMMMMMM ##
        ##MMMMMMMMMMMMMMMMMMMMMM##      微信机器人: [%v] 已经登录成功了。
        %v`, "\n", user.Name(), "\n")
	viper.Set("bot.name", user.Name())
}

/**
@method onLogout 当机器人检测到登出的时候，会触发事件，并会在事件中传递机器人的信息。
@param {*} user
*/
func onLogout(context *Context, user *user.ContactSelf, reason string) {
	fmt.Println("========================onLogout👇========================")
}

/*
  @method onRoomInvite 当收到群邀请的时候，会触发这个事件。
  @param {*} user
*/
func onRoomInvite(context *Context, roomInvitation *user.RoomInvitation) {
	fmt.Println("========================onRoomInvite👇========================")
	// TODO 自动通过群聊申请有问题 待解决(官方的问题)
	var (
		inviter  _interface.IContact
		roomName string
	)
	if err = roomInvitation.Accept(); err != nil {
		log.Errorf("[RoomInvite] Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
		return
	}
	if inviter, err = roomInvitation.Inviter(); err != nil {
		log.Errorf("[RoomInvite] 获取用户属性失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
		return
	}
	if roomName, err = roomInvitation.Topic(); err != nil {
		log.Errorf("[RoomInvite] 获取群聊名称失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
		return
	}
	log.Infof("[RoomInvite] 通过群聊邀请, 群聊名称: [%v] 邀请人: [%v]", roomName, inviter.Name())
	// 机器人进群自我介绍 onRoomJoin 已经实现
}

/*
	@method onRoomTopic 当有人修改群名称的时候会触发这个事件。
	@param {*} user
*/
func onRoomTopic(context *Context, room *user.Room, newTopic string, oldTopic string, changer _interface.IContact, date time.Time) {
	fmt.Println("========================onRoomTopic👇========================")
}

/*
	进入房间监听回调 room-群聊 inviteeList-受邀者名单 inviter-邀请者
	判断配置项群组id数组中是否存在该群聊id
*/
func onRoomJoin(context *Context, room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time) {
	fmt.Println("========================onRoomJoin👇========================")
	newUser := inviteeList[0].Name()
	if inviteeList[0].Self() {
		log.Infof("机器人加入群聊, 群聊名称:[%v] ,邀请人: [%v], 时间: [%v]", room.Topic(), inviter.Name(), date)
		if _, err = room.Say(fmt.Sprintf("大家好呀.我是%v, 以后请多多关照!", newUser)); err != nil {
			log.Errorf("[onRoomJoin] 加入群聊自我介绍消息发送失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		} else {
			log.Infof("[onRoomJoin] 加入群聊自我介绍消息发送成功")
			return
		}
	}
	log.Infof("群聊名称: [%v], 新人: [%v], 邀请人: [%v], 时间: [%v]", room.Topic(), newUser, inviter.Name(), date)
	//if !Plug.NightMode(inviter.ID()) {
	//	return
	//}
	if _, err = room.Say(fmt.Sprintf("@%v 欢迎新人!", newUser)); err != nil {
		log.Errorf("[onRoomJoin] 欢迎新人加入群聊消息发送失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		log.Infof("[onRoomJoin] 欢迎新人加入群聊消息发送成功")
	}
}

/*
	@method onRoomleave 当机器人把群里某个用户移出群聊的时候会触发这个时间。用户主动退群是无法检测到的。
	@param {*} user
*/
func onRoomleave(context *Context, room *user.Room, leaverList []_interface.IContact, remover _interface.IContact, date time.Time) {
	fmt.Println("========================onRoomleave👇========================")
	fmt.Printf("[onRoomleave] 群聊名称: [%v] 用户[%v] 被移出群聊", room.Topic(), leaverList[0].Name())
}

func onFriendship(context *Context, friendship *user.Friendship) {
	fmt.Println("========================onFriendship👇========================")
	switch friendship.Type() {
	case 0: // FriendshipTypeUnknown

	/*
		1. 新的好友请求
	*/
	case 1: // FriendshipTypeConfirm
		//log.Printf("friend ship confirmed with%v", friendship.Contact().Name())

	/*
		2. 通过好友申请
	*/
	case 2: // FriendshipTypeReceive
		if err = friendship.Accept(); err != nil {
			log.Errorf("[onFriendship] 添加好友失败, 好友名称: [%v], Error: [%v] CoptRight: [%s]", friendship.Contact().Name(), err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("[onFriendship] 添加好友成功, 好友名称:%v", friendship.Contact().Name())
	case 3: // FriendshipTypeVerify
		if err = friendship.GetWechaty().Friendship().Add(friendship.Contact(), fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("bot.name"))); err != nil {
			log.Errorf("[onFriendship] 添加好友失败, 好友名称: [%v], Error: [%v] CoptRight: [%s]", friendship.Contact().Name(), err, Copyright(make([]uintptr, 1)))
			return
		}
		log.Infof("[onFriendship] 添加好友成功, 好友名称:%v", friendship.Contact().Name())

	default:
		//	NONE
	}
	log.Infof("[onFriendship] %v好友关系是: %v Hello: %v ", friendship.Contact().Name(), friendship.Type(), friendship.Hello())
}

/*
	@method onHeartbeat 获取机器人的心跳。
	@param {*} user
*/
func onHeartbeat(context *Context, data string) {
	fmt.Println("========================onHeartbeat👇========================")
	log.Printf("[onHeartbeat] 获取机器人的心跳: %v", data)
}

func onError(context *Context, err error) {
	//log.Errorf("[onError] Error: [%v] 消息来自函数: [%v]", err, Plug.Copyright(make([]uintptr, 1)))
}

func main() {
	i := 0
	// 重试次数 10
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
		}
	}()
	for i <= 10 {
		i++
		// 钉钉推送
		WechatBotInit()
		var bot = NewWechaty(WithPuppetOption(puppet.Option{
			Token:    viper.GetString("Wechaty.Token"),
			Endpoint: viper.GetString("Wechaty.Endpoint"),
		}))
		log.Printf("Token: %v", viper.GetString("Wechaty.Token"))
		log.Printf("Endpoint: %v", viper.GetString("Wechaty.Endpoint"))
		log.Printf("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT: [%v]", viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT"))

		bot.OnScan(onScan).
			OnLogin(onLogin).
			OnLogout(onLogout).
			OnError(onError).
			OnRoomInvite(onRoomInvite). // 有问题，暂时不用，等待修复
			OnRoomTopic(onRoomTopic).
			OnRoomJoin(onRoomJoin).
			OnRoomLeave(onRoomleave).
			OnFriendship(onFriendship).
			Use(Pretreatment()).
			//Use(Test())
			Use(Group()).
			Use(Admin()).
			Use(Average()).
			Use(AutoReply()).
			Use(FileBox()).
			//OnMessage(onMessage).
			Use(ExportMessage()).
			//OnHeartbeat(onHeartbeat).
			Use(DingMessage())
			//bot.DaemonStart()

			/*
				Use(func() *Plugin {
					plug := NewPlugin()
					plug.OnMessage(func(context *Context, message *user.Message) {
					})
					return plug
				}()).
				Use(func() *Plugin {
					plug := NewPlugin()
					plug.OnMessage(func(context *Context, message *user.Message) {
					})
					return plug
				}())
			// */

		if err = bot.Start(); err != nil {
			// 重启Bot
			log.Infof("[main] Bot 错误: %v", err)
			if i > 10 {
				os.Exit(0)
			}
			log.Printf("正在重新启动程序, 当前重试次数: 第%v次", i)
			time.Sleep(10 * time.Second)
		} else {
			i = 0
			// Bot 守护程序
			var quitSig = make(chan os.Signal)
			signal.Notify(quitSig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-quitSig:
				ViperWrite()
				log.Fatal("程序退出!")
			}
		}
	}
}
