package ExportMessages

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"
)

var (
	err error
	db  *gorm.DB
)

/*
	ExportMessages()
	对消息内容进行存储
*/

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.
		OnLogin(onLogin).
		OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
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

	if err = db.Table(message.Date().Format("2006-01")).AutoMigrate(&MessageInfo{}); err != nil {
		log.Errorf("自动创建表失败, %s", err)
		return
	}
	if err = db.Table(message.Date().Format("2006-01")).Create(&m).Error; err != nil {
		log.Errorf("写入数据失败, %s", err)
		return
	}
}

func onLogin(context *wechaty.Context, user *user.ContactSelf) {
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       viper.GetString("MYSQL.HOST"), // DSN data source name
		DefaultStringSize:         256,                           // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                          // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                          // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                          // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                         // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
}
