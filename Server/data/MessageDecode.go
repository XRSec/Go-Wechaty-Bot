package data

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"os"
	"strings"
	"time"
)

type (
	MessageInfo struct {
		Data     string
		Status   bool // 群聊属性
		AtMe     bool // 是否@我
		RoomName string
		RoomID   string
		UserName string
		UserID   string
		Content  string
		AutoInfo string
		Reply    string
	}
)

func EncodeMessage(message *user.Message) MessageInfo {
	//var result []byte
	UserName := message.From().Name()
	UserID := message.From().ID()
	messages := MessageInfo{
		Status:   false,
		AtMe:     false,
		UserName: UserName,
		UserID:   UserID,
		Content:  message.Text(),
		AutoInfo: "用户ID: [" + UserID + "] 用户名称: [" + UserName + "] 说: [" + message.Text() + "]",
	}
	if message.Room() != nil {
		messages.Status = true
		messages.RoomID = message.Room().ID()
		messages.RoomName = strings.Replace(strings.Replace(message.Room().String(), "Room<", "", 1), ">", "", 1)
		messages.AutoInfo = "群聊ID: [" + messages.RoomID + "] 群聊名称: [" + messages.RoomName + "] " + messages.AutoInfo
		if message.MentionSelf() || strings.Contains(message.Text(), "@"+viper.GetString("bot.name")) {
			messages.AtMe = true
			messages.Content = strings.Replace(message.Text(), "@"+viper.GetString("bot.name"), "", 1)
		}
	}
	return messages
}

func ExportMessages(context MessageInfo) {
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
	)
	context.Data = time.Now().Format("2006_01_02_15_04_05.00000")
	var result []byte
	if result, err = json.Marshal(context); err != nil {
		ErrorFormat("Json 解析失败!", err)
	}
	if _, err = os.Stat(filename); err != nil {
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
	if _, err = fp.Write(result); err != nil {
		ErrorFormat("写入聊天记录到聊天备份文件失败!", err)
	}
	if _, err = fp.WriteString("\n"); err != nil {
		ErrorFormat("写入换行符到聊天备份文件失败!", err)
	}
}
