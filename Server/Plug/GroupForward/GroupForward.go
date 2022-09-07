package GroupForward

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	. "github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io/ioutil"
	"os"
	"strings"
	"time"
	. "wechatBot/General"
)

/*
 目前已知问题：
 1. 微信默认的服务号（文件传输助手）会识别成好友
 2. 消息频率 风控问题
*/
var (
	lists []_interface.IContact
	err   error
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(func(context *Context, message *user.Message) {
		if message.Self() || message.Talker().ID() == viper.GetString("BOT.ADMINID") {
		} else {
			return
		}
		if message.Type() != schemas.MessageTypeText {
			return
		}
		if !strings.Contains(message.Text(), "节日祝福") {
			return
		}
		if message.Text()[0:13] == "节日祝福 " {
			if _, err = os.Stat("friend.json"); err != nil {
				getAllToFile(message.GetWechaty().Contact())
				SayMessage(context, message, fmt.Sprintf("群发开始,共计%v人", len(lists)))
			} else {
				readFromFile()
				SayMessage(context, message, fmt.Sprintf("群发继续, 剩余%v人", len(lists)))
			}
			//if msg == "" {
			//	if message.Text()[0:7] == "群发 " {
			//		msg = message.Text()[7:]
			//	}
			//	if message.Text()[0:8] == "forward " {
			//		msg = message.Text()[8:]
			//	}
			//}
			for i := 1; i < len(lists); i = 0 {
				if _, err = message.GetWechaty().Contact().Load(lists[i].ID()).Say(fmt.Sprintf("嗨,亲爱的%v, %v", lists[i].Name(), message.Text()[13:])); err != nil {
					_, _ = message.Say(fmt.Sprintf("群发失败, 剩余%v人未发送成功", len(lists)))
					writeToFile()
					return
				}
				lists = append(lists[:0], lists[(1):]...)
				writeToFile()
				time.Sleep(time.Second * 8)
			}
			if err := os.Remove("friend.json"); err != nil {
				log.Errorf("os.Remove Error: [%v]", err)
				return
			}
		}
		if message.Text()[0:19] == "节日祝福测试 " {
			SayMessage(context, message, fmt.Sprintf("嗨, 亲爱的%v, %v", message.Talker().Name(), message.Text()[19:]))
		}
	})
	//Do(plug.Wechaty)
	return plug
}

func getAllToFile(c _interface.IContactFactory) {
	var lists2 []_interface.IContact
	lists = c.FindAll(nil)
	log.Infoln("ContactList: 加载成功")
	for _, v := range lists {
		if v.Type() != schemas.ContactTypePersonal {
			continue
		}
		if !v.Friend() {
			continue
		}
		lists2 = append(lists2, v)
	}
	lists = lists2
	writeToFile()
}

func readFromFile() {
	result, err := ioutil.ReadFile("friend.json")
	if err != nil {
		log.Errorf("ioutil.ReadFile Error: [%v]", err)
		return
	}
	err = json.Unmarshal(result, &lists)
	if err != nil {
		log.Errorf("json.Unmarshal Error: [%v]", err)
		return
	}
}

func writeToFile() {
	result, err := json.Marshal(lists)
	if err != nil {
		log.Errorf("json.Marshal Error: [%v]", err)
		return
	}
	_ = ioutil.WriteFile("friend.json", result, 0644)
}
