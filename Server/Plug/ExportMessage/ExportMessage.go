package ExportMessages

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	. "wechatBot/General"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	fp       *os.File
	filename = viper.GetString("rootPath") + "/data.json"
	result   []byte
	err      error
)

/*
	ExportMessages()
	对消息内容进行存储
*/
func ExportMessage() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(*MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed")
	}
	if m.Pass {
		m.AutoInfo += fmt.Sprintf(" Pass: [%v]", m.PassResult)
	}
	if m.Reply {
		m.AutoInfo += fmt.Sprintf(" 回复: [%v]", m.ReplyResult)
	}
	if message.Type() != schemas.MessageTypeText {
		log.Errorf("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
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

	if result, err = json.Marshal(m); err != nil {
		log.Errorf("[ExportMessages] Json 解析失败! Error: [%v]", err)
		return
	}
	if _, err = os.Stat(filename); err != nil {
		log.Errorf("[ExportMessages] 聊天备份文件不存在,正在创建! Error: [%v]", err)
	}
	if fp, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
		log.Errorf("[ExportMessages] 打开聊天备份文件失败! Error: [%v]", err)
	}
	defer func(fp *os.File) {
		if err = fp.Close(); err != nil {
			log.Errorf("[ExportMessages] 关闭聊天备份文件失败! Error: [%v]", err)
		}
	}(fp)
	if _, err = fp.Write(result); err != nil {
		log.Errorf("[ExportMessages] 写入聊天记录到聊天备份文件失败! Error: [%v]", err)
		return
	}
	if _, err = fp.WriteString("\n"); err != nil {
		log.Errorf("[ExportMessages] 写入换行符到聊天备份文件失败! Error: [%v]", err)
		return
	}
	log.Printf(m.AutoInfo)
	context.SetData("msgInfo", &m)
}
