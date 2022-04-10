package data

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
)

func AddFriend(weixinID string, friendship *user.Friendship) {
	log.Println("========================AddFriend👇========================")
	var test = schemas.RoomQueryFilter{
		Id: weixinID,
	}
	log.Println(friendship.GetWechaty().Room().Find(test).ID())
}
