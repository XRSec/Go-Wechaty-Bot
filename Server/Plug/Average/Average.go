package Average

import (
	"fmt"
	"strings"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	reply := ""
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Pass {
		log.Infof("Pass CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Reply {
		log.Infof("Reply CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.AtMe {
		log.Infof("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Infof("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Infof("Self CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Infof("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}
	// 这里需要确定三件事,第一 admin 插件权限更高,所以先走admin插件
	// 第二,如果是管理员,则不需要走普通插件
	// 第三 Admin 插件没有 wxid phone 关键字, 匹配用的是 (MentionText=="add")
	if strings.Contains(message.MentionText(), "add") {
		var (
			wxQuery schemas.FriendshipSearchCondition
			search  _interface.IContact
		)
		if strings.Contains(message.MentionText(), "wxid") {
			wxQuery.WeiXin = strings.Replace(strings.Replace(strings.Replace(message.MentionText(), "add", "", 1), "wxid", "", 1), " ", "", -1)
		} else if strings.Contains(message.MentionText(), "phone") {
			wxQuery.Phone = strings.Replace(strings.Replace(strings.Replace(message.MentionText(), "add", "", 1), "phone", "", 1), " ", "", 1)
		} else {
			log.Infof("wxQuery: [%v] CoptRight: [%v]", message.Text(), Copyright(make([]uintptr, 1)))
			goto end
		}
		if search, err = message.GetWechaty().Friendship().Search(&wxQuery); err != nil || search == nil {
			SayMessage(context, message, "查询失败")
		}
		if message.GetWechaty().Contact().Load(search.ID()).Friend() {
			log.Infof("用户已经是好友, 用户名: [%v] CoptRight: [%s]", search.Name(), Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("用户: [%v] 已经是好友了", search.Name()))
			return
		}
		if err = message.GetWechaty().Friendship().Add(search, fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("Bot.Name"))); err != nil {
			log.Errorf("添加好友失败, 用户名: [%v], Error: [%v] CoptRight: [%s]", search.Name(), err, Copyright(make([]uintptr, 1)))
			SayMessage(context, message, fmt.Sprintf("添加好友失败, 用户: [%v]", search.Name()))
			return
		}
	end:
	}

	if strings.Contains(message.MentionText(), "djs") {
		log.Infof("添加定时提醒成功! 任务详情: %v CoptRight: [%s]", "暂无", Copyright(make([]uintptr, 1)))
		reply = "添加定时提醒成功! 任务详情: 暂无"
	}

	if strings.Contains(message.MentionText(), "fdj") {
		log.Infof("复读机模式, 复读内容: [%v] CoptRight: [%s]", message.MentionText(), Copyright(make([]uintptr, 1)))
		reply = strings.Replace(message.MentionText(), "fdj ", "", 1)
	}

	if strings.Contains(message.MentionText(), "print") {
		reply = strings.Replace(message.MentionText(), "print", "", 1)
	}

	if !m.Status {
		if message.Text() == "加群" || message.Text() == "交流群" {
			keys := ""
			/* TODO 临时方案 GetStringMap 异常 */
			out := GetStringMapFixed("group")
			//out := viper.GetStringMap("group")
			/* https://github.com/spf13/viper/issues/708 */

			for k := range out {
				keys += "『" + k + "』"
			}
			reply = "现有如下交流群, 请问需要加入哪个呢? 请发交流群名字!\n" + keys
		}
		/* TODO GetStringMap */
		if strings.Contains(fmt.Sprintf("%v", viper.GetStringMap("Group")), message.Text()) {
			for i, v := range viper.GetStringMap("Group") {
				if strings.Contains(message.Text(), i) && v != "" {
					//	邀请好友进群
					if err = message.GetWechaty().Room().Load(v.(string)).Add(message.Talker()); err != nil {
						log.Errorf("邀请好友进群失败, Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
						reply = "邀请好友进群失败, 群已满或者不存在!"
						goto end2
					}
					log.Infof("邀请好友进群成功! CoptRight: [%s]", Copyright(make([]uintptr, 1)))
					reply = "已经拉你啦! 等待管理员审核通过呀!"
				}
			}
		}
		/* end GetStringMap */
	}
end2:
	if reply == "" {
		return
	}
	SayMessage(context, message, reply)
}

func GetStringMapFixed(key string) map[string]interface{} {
	out := make(map[string]interface{})
	for _, k := range viper.AllKeys() {
		if !strings.Contains(k, ".") {
			continue
		}
		lastInd := strings.LastIndex(k, ".")
		if k[:lastInd] == key {
			//if strings.HasPrefix(, key) {
			out[strings.TrimPrefix(k, key+".")] = viper.Get(k)
		}
	}
	return out
}
