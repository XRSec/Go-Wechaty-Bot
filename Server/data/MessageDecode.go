package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type (
	MessageInfo struct {
		Data     string `json:"data"`     // 日期
		Status   bool   `json:"status`    // 群聊属性
		AtMe     bool   `json:"atme"`     // 是否@我
		RoomName string `json:"roomname"` // 群聊名称
		RoomID   string `json:"roomid"`   // 群聊ID
		UserName string `json:"username"` // 用户名称
		UserID   string `json:"userid"`   // 用户ID
		Content  string `json:"content"`  // 聊天内容
		AutoInfo string `json:"autoinfo"` // 信息一览
		Reply    string `json:"reply"`    // 自动回复
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
		Content:  message.MentionText(),
		//AutoInfo: "用户ID: [" + UserID + "] 用户名称: [" + UserName + "] 说: [" +  + "]",
		AutoInfo: fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v]", UserID, UserName, strings.Replace(message.Text(), "\u2005", " ", -1)),
		//	 TODO Autoinfo 内容
	}
	if message.Room() != nil {
		messages.Status = true
		messages.RoomID = message.Room().ID()
		messages.RoomName = strings.Replace(strings.Replace(message.Room().String(), "Room<", "", 1), ">", "", 1)
		//messages.AutoInfo = "群聊ID: [" + messages.RoomID + "] 群聊名称: [" + messages.RoomName + "] " + messages.AutoInfo
		messages.AutoInfo = fmt.Sprintf("群聊ID: [%v] 群聊名称: [%v] %s", messages.RoomID, messages.RoomName, messages.AutoInfo)
		if message.MentionSelf() {
			//if message.MentionSelf() || strings.Contains(message.Text(), "@"+viper.GetString("bot.name")) {
			messages.AtMe = true
			//messages.Content = strings.Replace(strings.Replace(message.Text(), "\u2005", "", -1), "@"+viper.GetString("bot.name"), "", 1) // 去皮操作
		}
	}
	return messages
}

/*
	Json 日志格式化
	ExportMessages(messages MessageInfo)
*/
func ExportMessages(messages MessageInfo) {
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
	)
	messages.Data = time.Now().Format("2006-01-02T15:04:05.00000")
	var result []byte
	if result, err = json.Marshal(messages); err != nil {
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
