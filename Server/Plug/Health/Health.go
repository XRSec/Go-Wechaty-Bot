package Health

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	. "github.com/wechaty/go-wechaty/wechaty"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"time"
	"wechatBot/Plug/DingMessage"
	"wechatBot/Plug/Jinrishici"
)

var (
	err error
)

/*
	Health()
	健康监测
*/
func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	Do(plug.Wechaty)
	return plug
}

func Do(bot *Wechaty) {
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(nyc))
	if _, err := c.AddFunc("0 23 * * *", func() {
		var roomID _interface.IRoom
		//if roomID = bot.Room().Find(&schemas.RoomQueryFilter{Id: "roomID@chatroom"}); roomID == nil {
		if roomID = bot.Room().Find(viper.Get("BOT.GROUP")); roomID == nil {
			DingMessage.DingSend(viper.GetString("Bot.AdminID"), "RoomID Find Error")
			log.Infof("RoomID Find Error")
			return
		}

		if _, err := roomID.Say(Jinrishici.Do()); err != nil {
			DingMessage.DingSend(viper.GetString("Bot.AdminID"), "failed to send messages")
			log.Errorf("onHeartbeat Say Error: [%v]", err)
			return
		}
		log.Infof("Heartbeat Say Success")
	}); err != nil {
		DingMessage.DingSend(viper.GetString("Bot.AdminID"), "Heartbeat Cron Add Error: "+err.Error())
		log.Errorf("Heartbeat Cron Add Error: [%v]", err)
	}
	log.Infof("Health Cron Start")
	c.Start()
}
