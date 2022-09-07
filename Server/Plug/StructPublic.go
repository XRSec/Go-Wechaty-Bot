package Plug

import (
	"time"

	"gorm.io/gorm"
)

type (
	MessageInfo struct {
		ID          string         `json:"ID" gorm:"column:id;comment:消息ID;primaryKey"`               // 消息ID
		Date        string         `json:"Date" gorm:"column:date;comment:日期"`                        // 时间
		Status      bool           `json:"Status" gorm:"column:status;comment:群聊?"`                   // 群聊属性
		AtMe        bool           `json:"AtMe" gorm:"column:atme;comment:提及?"`                       // 是否@我
		RoomName    string         `json:"RoomName" gorm:"column:room_name;comment:群聊名称"`             // 群聊名称
		RoomID      string         `json:"RoomID" gorm:"column:room_id;comment:群聊ID"`                 // 群聊ID
		UserName    string         `json:"UserName" gorm:"column:user_name;comment:用户名"`              // 用户名称
		UserID      string         `json:"UserID" gorm:"column:user_id;comment:用户ID"`                 // 用户ID
		Content     string         `json:"Content" gorm:"column:content;comment:内容;type:TEXT(10000)"` // 聊天内容
		ReplyResult string         `json:"ReplyResult" gorm:"column:reply_result;comment:回复状态"`       // 自动回复
		Reply       bool           `json:"Reply" gorm:"column:reply;comment:回复内容;type:TEXT(10000)"`   // 自动回复状态
		PassResult  string         `json:"PassResult" gorm:"column:pass_result;comment:跳过?"`          // Pass原因
		Pass        bool           `json:"Pass" gorm:"column:pass;comment:跳过原因"`                      // Pass状态
		CreatedAt   time.Time      `json:"CreatedAt" gorm:"column:created_at;comment:创建时间"`           // 创建时间
		UpdatedAt   time.Time      `json:"UpdatedAt" gorm:"column:updated_at;comment:更新时间"`           // 更新时间
		Deleted     gorm.DeletedAt `json:"Deleted" gorm:"column:deleted;comment:删除时间"`                // 删除时间
	}
)
