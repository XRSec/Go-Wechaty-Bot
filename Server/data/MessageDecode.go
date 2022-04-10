package data

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"os"
	"strings"
	"time"
)

var ()

type (
	MessageInfo struct {
		Data     string
		Status   bool
		RoomName string
		RoomID   string
		UserName string
		UserID   string
		AutoInfo string
	}
)

func EncodeMessage(message *user.Message) MessageInfo {
	//var result []byte
	UserName := message.From().Name()
	UserID := message.From().ID()
	messages := MessageInfo{
		Status:   false,
		UserName: UserName,
		UserID:   UserID,
		AutoInfo: "用户ID: [" + UserID + "] 用户名称: [" + UserName + "]",
	}
	if message.Room() != nil {
		messages.Status = true
		messages.RoomID = message.Room().ID()
		messages.RoomName = strings.Replace(strings.Replace(message.Room().String(), "Room<", "", 1), ">", "", 1)
		messages.AutoInfo = "群聊ID: [" + messages.RoomID + "] 群聊名称: [" + messages.RoomName + "] " + messages.AutoInfo
	}

	u := messages
	u.Data = time.Now().Format("2006_01_02_15_04_05.00000")
	u.AutoInfo = message.Text()
	var result []byte
	var err error
	if result, err = json.Marshal(u); err != nil {
		ErrorFormat("Json 解析失败!", err)
	}
	if message.Type() == schemas.MessageTypeText {
		go exportMessages(result)
	}
	return messages
}

func exportMessages(context []byte) {
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
	)
	if _, err := os.Stat(filename); err != nil {
		ErrorFormat("聊天备份文件不存在,正在创建!", err)
	}
	if fp, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
		ErrorFormat("打开聊天备份文件失败!", err)
	}
	defer func(fp *os.File) {
		if err = fp.Close(); err != nil {
			ErrorFormat("关闭聊天备份文件失败!", err)
		}
	}(fp)
	if _, err = fp.Write(context); err != nil {
		ErrorFormat("写入聊天记录到聊天备份文件失败!", err)
	}
	if _, err = fp.WriteString("\n"); err != nil {
		ErrorFormat("写入换行符到聊天备份文件失败!", err)
	}
}
