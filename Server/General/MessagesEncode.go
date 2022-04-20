package General

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type (
	MessageInfo struct {
		Data     string `json:"data"`     // 日期
		Status   bool   `json:"status"`   // 群聊属性
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

func EncodeMessage(message *user.Message) {
	viper.Set(fmt.Sprintf("Chat.%s.ReplyStatus", message.From().ID()), "false")
}

func ExportMessages(message *user.Message) {
	if message.Type() != schemas.MessageTypeText {
		return
	}
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
		result   []byte
		userName = message.From().Name()
		userID   = message.From().ID()
		content  = message.Text()
		reply    = viper.GetString(fmt.Sprintf("Chat.%s.Reply", userID))
		messages = MessageInfo{
			Data:     message.Date().Format("2006-01-02 15:04:05"),
			Status:   false,
			AtMe:     false,
			UserName: userName,
			UserID:   userID,
			Content:  content,
			Reply:    reply,
			AutoInfo: fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v] 回复: [%s]", userID, userName, content, reply),
		}
	)
	if message.Room() != nil {
		messages.RoomName = message.Room().Topic()
		messages.RoomID = message.Room().ID()
		messages.Status = true
		messages.AutoInfo = fmt.Sprintf("群聊ID: [%v] 群聊名称: [%v] %s", messages.RoomID, messages.RoomName, messages.AutoInfo)
	}
	if result, err = json.Marshal(messages); err != nil {
		log.Errorf("[ExportMessages] Json 解析失败! Error: [%s]", err)
		return
	}
	if _, err = os.Stat(filename); err != nil {
		log.Errorf("[ExportMessages] 聊天备份文件不存在,正在创建! Error: [%s]", err)
	}
	if fp, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
		log.Errorf("[ExportMessages] 打开聊天备份文件失败! Error: [%s]", err)
		//	TODO 这里好像找不到文件会创建文件
	}
	defer func(fp *os.File) {
		if err = fp.Close(); err != nil {
			log.Errorf("[ExportMessages] 关闭聊天备份文件失败! Error: [%s]", err)
		}
	}(fp)
	if _, err = fp.Write(result); err != nil {
		log.Errorf("[ExportMessages] 写入聊天记录到聊天备份文件失败! Error: [%s]", err)
		return
	}
	if _, err = fp.WriteString("\n"); err != nil {
		log.Errorf("[ExportMessages] 写入换行符到聊天备份文件失败! Error: [%s]", err)
		return
	}
	log.Printf(messages.AutoInfo)
}
