package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	. "wechatBot/api"
	. "wechatBot/data"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/viper"
	. "github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	. "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func onScan(context *Context, qrCode string, status schemas.ScanStatus, data string) {
	fmt.Printf("\n\n")
	i := 0
	if status.String() == "ScanStatusWaiting" {
		i++
		if i > 5 {
			os.Exit(1)
		}
		qrterminal.GenerateWithConfig(qrCode, qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: " \u2005",
			WhiteChar: "💵",
			QuietZone: 1,
		})
		fmt.Printf("\n\n")
		log.Printf("%s[Scan] https://wechaty.js.org/qrcode/%s %s\n", viper.GetString("info"), qrCode, data)
		messages := fmt.Sprintf("账号未登录请扫码!\n\n---\n\n[qrCode](https://wechaty.js.org/qrcode/%v)", qrCode)
		DingMessagePic(messages, viper.GetString("bot.adminid"))
		time.Sleep(120 * time.Second)
	} else if status.String() == "ScanStatusScanned" {
		log.Printf("%v[Scan] Status: %v %v\n", viper.GetString("info"), status.String(), data)
	} else {
		log.Printf("%v[Scan] Status: %v %v\n", viper.GetString("info"), status.String(), data)
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
        ##MMMMMMMMMMMMMMMMMMMMMM##      微信机器人: [%s] 已经登录成功了。
        %s`, "\n", user.Name(), "\n")
	viper.Set("bot.name", user.Name())
}

/**
@method onLogout 当机器人检测到登出的时候，会触发事件，并会在事件中传递机器人的信息。
@param {*} user
*/
func onLogout(context *Context, user *user.ContactSelf, reason string) {
	log.Println("========================onLogout👇========================")
	messages := fmt.Sprintf("%v账号已退出登录, 请检查账号!\n\n---\n\n错误: %v", user.Name(), reason)
	fmt.Println(messages)
	DingMessagePic(messages, viper.GetString("bot.adminid"))
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
		log.Errorf("[RoomInvite] Error: %v", err)
		return
	}
	if inviter, err = roomInvitation.Inviter(); err != nil {
		log.Errorf("[RoomInvite] 获取用户属性失败, Error: %v", err)
		return
	}
	if roomName, err = roomInvitation.Topic(); err != nil {
		log.Errorf("[RoomInvite] 获取群聊名称失败, Error: %v", err)
		return
	}
	log.Printf("[RoomInvite] 通过群聊邀请, 群聊名称: [%v] 邀请人: [%v]", roomName, inviter.Name())
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
		log.Printf("机器人加入群聊, 群聊名称:[%v] ,邀请人: [%v], 时间: [%v]", room.Topic(), inviter.Name(), date)
		if _, err = room.Say(fmt.Sprintf("大家好呀.我是%v, 以后请多多关照!", newUser)); err != nil {
			log.Errorf("[onRoomJoin] 加入群聊自我介绍消息发送失败, Error: %v", err)
			return
		} else {
			log.Printf("[onRoomJoin] 加入群聊自我介绍消息发送成功")
			return
		}
	}
	log.Printf("群聊名称: [%v], 新人: [%v], 邀请人: [%v], 时间: [%v]", room.Topic(), newUser, inviter.Name(), date)
	if !Plug.NightMode(inviter.ID()) {
		return
	}
	if _, err = room.Say(fmt.Sprintf("@%v 欢迎新人!", newUser)); err != nil {
		log.Errorf("[onRoomJoin] 欢迎新人加入群聊消息发送失败, Error: %v", err)
	} else {
		log.Printf("[onRoomJoin] 欢迎新人加入群聊消息发送成功")
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
			log.Errorf("[onFriendship] 添加好友失败, 好友名称: [%v], Error: [%v]", friendship.Contact().Name(), err)
			return
		}
		log.Printf("[onFriendship] 添加好友成功, 好友名称:%v", friendship.Contact().Name())
	case 3: // FriendshipTypeVerify
		if err = friendship.GetWechaty().Friendship().Add(friendship.Contact(), fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("bot.name"))); err != nil {
			log.Errorf("[onFriendship] 添加好友失败, 好友名称: [%v], Error: [%v]", friendship.Contact().Name(), err)
			return
		}
		log.Printf("[onFriendship] 添加好友成功, 好友名称:%v", friendship.Contact().Name())

	default:
		//	NONE
	}
	log.Printf("[onFriendship] %v好友关系是: %v Hello: %v ", friendship.Contact().Name(), friendship.Type(), friendship.Hello())
}

/*
	@method onHeartbeat 获取机器人的心跳。
	@param {*} user
*/
func onHeartbeat(context *Context, data string) {
	fmt.Println("========================onHeartbeat👇========================")
	log.Printf("[onHeartbeat] 获取机器人的心跳: %v", data)
}

func groupChat(messages MessageInfo, message *user.Message) {
	if message.MentionSelf() { // @我 的我操作
		// 管理员操作台 以下操作都需要管理员权限
		if message.From().ID() == viper.GetString("bot.adminid") {
			log.Printf("MentionText: [%s]", message.MentionText())
			log.Printf("MentionSelf: [%v]", message.MentionSelf())
			//if strings.Contains(messages.Content, "add") { // 添加好友
			if message.MentionText() == "add" { // 添加好友
				addUserName := strings.Replace(strings.Replace(message.Text(), "\u2005", "", -1), fmt.Sprintf("@%sadd @", viper.GetString("bot.name")), "", 1) // 过滤用户名
				if member, err := message.Room().Member(addUserName); err != nil && member != nil {                                                            //查找添加用户的 ID
					ErrorFormat(fmt.Sprintf("搜索用户名ID失败, 用户名: [%s], 用户信息: [%s]", addUserName, member.String()), err)
				} else {
					SuccessFormat("搜索用户名ID成功!")
					if message.GetWechaty().Contact().Load(member.ID()).Friend() {
						if _, err := message.GetWechaty().Contact().Load(member.ID()).Say("已经是好友啦!"); err != nil {
							ErrorFormat("向还有发送已经是好友的消息发送失败, Error: ", err)
						}
						if _, err = message.Say("已经是好友了!"); err != nil {
							ErrorFormat("已经是好友的消息发送失败, Error: ", err)
						} else {
							SuccessFormat("已经是好友的消息发送成功!")
						}
					} else {
						SuccessFormat("您与对方不是好友, 正在尝试添加!")
						if err = message.GetWechaty().Friendship().Add(member, fmt.Sprintf("你好,我是%s,以后请多多关照!", viper.GetString("bot.name"))); err != nil {
							ErrorFormat(fmt.Sprintf("添加好友失败, 用户名: [%s]用户ID:[%s], Error: ", member.Name(), member.ID()), err)
							if _, err = message.Say("好友申请失败!"); err != nil {
								ErrorFormat("好友申请失败 通知失败, Error: ", err)
							} else {
								SuccessFormat("好友申请失败 通知成功!")
							}
						} else {
							if _, err = message.Say("好友申请发送成功!"); err != nil {
								ErrorFormat("好友申请发送成功 通知失败, Error: ", err)
							} else {
								SuccessFormat("好友申请发送成功 通知成功!")
							}
						}
					}
				}
				return
				//} else if strings.Contains(messages.Content, "del") { // 从群聊中移除用户
			} else if message.MentionText() == "del" { // 从群聊中移除用户
				deleteUserName := strings.Replace(strings.Replace(message.Text(), "\u2005", "", -1), fmt.Sprintf("@%sdel @", viper.GetString("bot.name")), "", 1) // 过滤用户名
				//deleteUserName := strings.Replace(strings.Replace(strings.Replace(messages.Content, "delete", "", 1), "@", "", 1), " ", "", -1)
				if member, err := message.Room().Member(deleteUserName); err != nil && member != nil {
					ErrorFormat(fmt.Sprintf("搜索用户名ID失败: [%s]", deleteUserName), err)
				} else {
					SuccessFormat(fmt.Sprintf("搜索用户名ID成功, 用户名: [%s]", deleteUserName))
					if err = message.Room().Del(member); err != nil {
						ErrorFormat(fmt.Sprintf("从群聊中删除成员失败: [%s]", deleteUserName), err)
					} else {
						SuccessFormat(fmt.Sprintf("从群聊中删除成员成功!, 用户名: [%s]", deleteUserName))
					}
				}
				return
				//} else if strings.Contains(messages.Content, "quit") { // 退群
			} else if message.MentionText() == "quit" { // 退群
				if _, err = message.Say("我走了, 拜拜👋🏻, 记得想我哦 [大哭]"); err != nil {
					ErrorFormat("退出群聊 告别语发送失败! ", err)
				} else {
					if err = message.Room().Quit(); err != nil {
						ErrorFormat(fmt.Sprintf("退出群聊失败, 群聊名称: [%s], Error: ", messages.RoomName), err)
					} else {
						SuccessFormat(fmt.Sprintf("退出群聊成功! 群聊名称: [%s]", messages.RoomName))
					}
				}
			} else if strings.Contains(messages.Content, "gmz") {
				newName := strings.Replace(message.MentionText(), "gmz ", "", 1)
				if err = message.GetPuppet().SetContactSelfName(newName); err != nil {
					if _, err = message.Say(fmt.Sprintf("修改用户名失败, Error: %v", err)); err != nil {
						ErrorFormat("发送修改用户名失败消息 失败, Error:", err)
					} else {
						SuccessFormat("发送修改用户名失败消息 成功!")
					}
				} else {
					if _, err = message.Say(fmt.Sprintf("修改用户名成功! 新的名称: %s", newName)); err != nil {
						ErrorFormat("发送修改用户名成功消息 失败, Error:", err)
					} else {
						SuccessFormat("发送修改用户名成功消息 成功!")
					}
				}
				return
			}
		}

		// 非管理员操作台
		if strings.Contains(messages.Content, "djs") {
			return
		} else if strings.Contains(messages.Content, "fdj") {
			log.Println("----------------")
			if _, err = message.Say(strings.Replace(messages.Content, "fdj ", "", 1)); err != nil {
				ErrorFormat("复读机消息发送失败, Error: ", err)
			} else {
				SuccessFormat("复读机消息发送成功")
			}
			return
		}
		//没有匹配指令,调用机器人回复 记得最后 return
		SayMessage(messages, message)
		DingMessageText(messages.AutoInfo, messages.UserID)
	} // 没有 @我 就老老实实的
}

/*
	messages MessageInfo, message *user.Message
*/
func privateChat(messages MessageInfo, message *user.Message) {
	if strings.Contains(messages.Content, "加群") || strings.Contains(messages.Content, "交流群") {
		keys := ""
		for k := range viper.GetStringMap("Group") {
			keys += "『" + k + "』"
		}
		reply := "现有如下交流群, 请问需要加入哪个呢? 请发交流群名字!\n" + keys
		if _, err = message.Say(reply); err != nil {
			ErrorFormat("群聊介绍发送失败, Error: ", err)
		} else {
			SuccessFormat("群聊介绍发送成功!")
		}
		return
	} else if strings.Contains(fmt.Sprintf("%s", viper.GetStringMap("Group")), messages.Content) {
		for i, v := range viper.GetStringMap("Group") {
			if strings.Contains(messages.Content, i) && v != "" {
				//	邀请好友进群
				if err = message.GetWechaty().Room().Load(v.(string)).Add(message.From()); err != nil {
					ErrorFormat("邀请好友进群失败, Error: ", err)
					return
				} else {
					SuccessFormat("邀请好友进群成功!")
					if _, err = message.Say("已经拉你啦! 等待管理员审核通过呀!"); err != nil {
						ErrorFormat("邀请好友成功提示信息发送失败, Error:", err)
					} else {
						SuccessFormat("邀请好友成功提示信息发送成功!")
					}
					return
				}
			}
			log.Printf("用户输入: [%s] i:[%v] i.key: [%s]", messages.Content, i, v)
		}
		//if _, err = message.Say("当前群聊我也没有权限,请重新输入!"); err != nil {
		//	ErrorFormat("群聊权限不足消息发送失败", err)
		//} else {
		//	SuccessFormat("群聊权限不足消息发送成功!")
		//}
		//return
	}
	SayMessage(messages, message)
}

func onMessage(context *Context, message *user.Message) {
	// 编码信息
	General.EncodeMessage(message) // map 加锁
	// Debug Model
	if message.Talker().ID() != viper.GetString("bot.adminid") {
		return
	}
	Plug.AdminManage(message)
	Plug.Manage(message)
	Plug.AutoReply(message)
	Plug.FileBox(message)
	if message.MentionSelf() {
		// 到这里的时候基本设置好了一些默认的值了
		Plug.DingMessage(fmt.Sprintf("%v@我了\n\n---\n\n### 用户属性\n\n用户名: [%v]\n\n用户ID: [%v]\n\n---\n\n### 群聊属性\n\n群聊名称: [%v]\n\n群聊ID: [%v]\n\n---\n\n**内容**: [%v]\n\n**回复**: [%v]", General.Messages.UserName, General.Messages.UserName, General.Messages.UserID, General.Messages.RoomName, General.Messages.RoomID, General.Messages.Content, General.Messages.Reply), General.Messages.UserID)
	}
	go General.ExportMessages()
}

func main() {
	i := 0
	// 重试次数 10
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	for i <= 10 {
		i++
		// 钉钉推送
		General.WechatBotInit()
		var bot = NewWechaty(WithPuppetOption(puppet.Option{
			Token:    viper.GetString("wechaty.token"),
			Endpoint: viper.GetString("wechaty.endpoint"),
		}))
		log.Printf("Token: %v", viper.GetString("wechaty.token"))
		log.Printf("Endpoint: %v", viper.GetString("wechaty.endpoint"))
		log.Printf("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT: [%v]", viper.GetString("wechaty.WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT"))

		bot.OnScan(onScan).
			OnLogin(onLogin).
			OnLogout(onLogout).
			OnMessage(onMessage).
			OnRoomInvite(onRoomInvite). // 有问题，暂时不用，等待修复
			OnRoomTopic(onRoomTopic).
			OnRoomJoin(onRoomJoin).
			OnRoomLeave(onRoomleave).
			OnFriendship(onFriendship).
			//OnHeartbeat(onHeartbeat).
			OnError(onError)
		//bot.DaemonStart()

		if err = bot.Start(); err != nil {
			// 重启Bot
			log.Printf("[main] Bot 错误: %v", err)
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
				General.ViperWrite()
				log.Fatal("程序退出!")
			}
		}
	}
}
