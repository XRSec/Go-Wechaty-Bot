package General

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type (
	MessageInfo struct {
		Date        string `json:"data"`        // 时间
		Status      bool   `json:"status"`      // 群聊属性
		AtMe        bool   `json:"atme"`        // 是否@我
		RoomName    string `json:"roomname"`    // 群聊名称
		RoomID      string `json:"roomid"`      // 群聊ID
		UserName    string `json:"username"`    // 用户名称
		UserID      string `json:"userid"`      // 用户ID
		Content     string `json:"content"`     // 聊天内容
		AutoInfo    string `json:"autoinfo"`    // 信息一览
		Reply       string `json:"reply"`       // 自动回复
		ReplyStatus bool   `json:"replystatus"` // 自动回复状态
		Pass        string `json:"pass"`        // pass 原因
		PassStatus  bool   `json:"passstatus"`  // pass 状态
	}
)

var Messages = MessageInfo{}

/*
	EncodeMessage()
	对消息内容进行编码
*/
func EncodeMessage(message *user.Message) {
	if message.Type() == schemas.MessageTypeRecalled {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		Messages.Pass = "MessageTypeRecalled"
		Messages.PassStatus = true
		return
	}
	if message.Type() == schemas.MessageTypeUnknown && message.Talker().Name() == "微信团队" {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		Messages.Pass = "微信团队"
		Messages.PassStatus = true
		return
	}
	if message.Type() != schemas.MessageTypeText {
		Messages.Content = "[未知消息类型: " + message.Type().String() + "] " + message.MentionText()
	} else {
		Messages.Content = message.Text()
	}
	Messages.Date = message.Date().Format("2006-01-02 15:04:05")
	Messages.Status = false
	Messages.AtMe = false
	Messages.UserName = message.Talker().Name()
	Messages.UserID = message.Talker().ID()
	Messages.AutoInfo = fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v]", Messages.UserID, Messages.UserName, strings.Replace(Messages.Content, "\u2005", " ", -1))
	Messages.RoomName = ""
	Messages.RoomID = ""
	Messages.Reply = ""
	Messages.ReplyStatus = false
	Messages.PassStatus = false
	if message.MentionSelf() {
		Messages.AtMe = true
	}
	if message.Room() != nil {
		Messages.RoomName = message.Room().Topic()
		Messages.RoomID = message.Room().ID()
		Messages.Status = true
		Messages.AutoInfo = fmt.Sprintf("群聊ID: [%v] 群聊名称: [%v] %v", Messages.RoomID, Messages.RoomName, Messages.AutoInfo)
	}
	// ChatTimeLimit(viper.GetString(fmt.Sprintf("Chat.%v.Date", Messages.UserID)))
}

/*
	ChatTimeLimit(message.Date().Format("2006-01-02 15:04:05"))
		: 判断消息是否在规定时间内
		: 如果是，则返回true，否则返回false
*/
func ChatTimeLimit(date string) {
	//当前时间
	var (
		now      time.Time
		loc      *time.Location
		lastDate time.Time
	)
	if date == "" {
		return
	}
	if Messages.Status && !Messages.AtMe {
		return
	}
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	if loc, err = time.LoadLocation("Local"); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Loc: [%v]", err, loc)
		// [ChatTimeLimit] time.ParseInLocation, Error: [The system cannot find the path specified.], Loc: [UTC]
		return
	}
	if now, err = time.ParseInLocation("2006-01-02 15:04:05", timeNow, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Now: [%v]", err, now)
		return
	}
	//当前时间转换为"年-月-日"的格式
	if lastDate, err = time.ParseInLocation("2006-01-02 15:04:05", date, loc); err != nil {
		log.Errorf("[ChatTimeLimit] time.ParseInLocation, Error: [%v], Lastdate: [%v]", err, lastDate)
		return
	}
	//计算两个时间相差的秒数
	if second := int(now.Sub(lastDate).Seconds()); second < 30 {
		log.Errorf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		// Messages.Reply = fmt.Sprintf("[ChatTimeLimit] 时间相差不足 开始时间: [%v], 结束时间: [%v], 相差秒数: [%d]", lastDate, now, second)
		Messages.Reply = ""
		Messages.ReplyStatus = false
		Messages.Pass = "ChatTimeLimit"
		Messages.PassStatus = true
		return
	}
}

/*
	ExportMessages()
	对消息内容进行存储
*/
func ExportMessages() {
	if Messages.ReplyStatus {
		Messages.AutoInfo += fmt.Sprintf(" 回复: [%v]", Messages.Reply)
	} else if Messages.PassStatus {
		Messages.AutoInfo += fmt.Sprintf(" Pass: [%v]", Messages.Pass)
	}
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
		result   []byte
	)
	if result, err = json.Marshal(Messages); err != nil {
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
	log.Printf(Messages.AutoInfo)
}
