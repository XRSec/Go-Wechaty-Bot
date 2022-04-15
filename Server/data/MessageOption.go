package data

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
	"strings"
)

func AddFriend(messages MessageInfo, message *user.Message) {
	if messages.UserID == viper.GetString("bot.adminid") {
		addUserName := strings.TrimRight(strings.Replace(strings.Replace(messages.Content, " add ", "", 1), "@", "", 1), " ")
		//addUserName := strings.Replace(strings.Replace(strings.Replace(messages.Content, "add", "", 1), "@", "", 1), " ", "", -1)
		if member, err := message.Room().Member(addUserName); err != nil && member != nil {
			ErrorFormat("搜索用户名ID失败: ["+addUserName+"]"+member.String(), err)
		} else {
			if err = message.GetWechaty().Friendship().Add(member, "你好!"); err != nil {
				ErrorFormat("添加好友失败: ["+addUserName+"]", err)
				log.Println(member)
			} else {
				SuccessFormat("添加好友成功!")
			}
		}
	}
}

func DeleteUser(messages MessageInfo, message *user.Message) {
	if messages.UserID == viper.GetString("bot.adminid") {
		deleteUserName := strings.TrimRight(strings.Replace(strings.Replace(messages.Content, " delete ", "", 1), "@", "", 1), " ")
		//deleteUserName := strings.Replace(strings.Replace(strings.Replace(messages.Content, "delete", "", 1), "@", "", 1), " ", "", -1)
		if member, err := message.Room().Member(deleteUserName); err != nil && member != nil {
			ErrorFormat("搜索用户名ID失败: ["+deleteUserName+"]"+member.String(), err)
		} else {
			//if err = message.GetWechaty().Friendship().Delete(member); err != nil {
			//	ErrorFormat("删除好友失败: ["+deleteUserName+"]", err)
			//	log.Println(member)
			//} else {
			//	SuccessFormat("删除好友成功!")
			//}
		}
	}
}

func InviteUser(messages MessageInfo, message *user.Message) {
	if message.Text() == viper.GetString("InviteGroupkeywords") {
		if _, err = message.Say("群暗号正确，已发起邀请，（可能需要群管理员同意，请您耐心等待~）"); err != nil {
			ErrorFormat("发送邀请入群消息失败: ["+message.Text()+"]", err)
			return
		} else {
			SuccessFormat("发送邀请入群消息成功!")
		}
		inviteUserName := strings.TrimRight(strings.Replace(strings.Replace(messages.Content, " invite ", "", 1), "@", "", 1), " ")
		if group, err := message.GetPuppet().RoomInvitationPayload(viper.GetString("InviteGroupID")); err != nil {
			ErrorFormat("邀请用户进群失败: "+messages.UserName, err)
			return
		} else {
			SuccessFormat("")
			fmt.Println(group.Id)
			//	这样写吗？
		}
		//inviteUserName := strings.Replace(strings.Replace(strings.Replace(messages.Content, "invite", "", 1), "@", "", 1), " ", "", -1)
		if member, err := message.Room().Member(inviteUserName); err != nil && member != nil {
			ErrorFormat("搜索用户名ID失败: ["+inviteUserName+"]"+member.String(), err)
		} else {
			//if err = message.GetWechaty().Room().Invite(message.Room(), member); err != nil {
			//	ErrorFormat("邀请好友失败: ["+inviteUserName+"]", err)
			//	log.Println(member)
			//} else {
			//	SuccessFormat("邀请好友成功!")
			//}
		}
	}
}
