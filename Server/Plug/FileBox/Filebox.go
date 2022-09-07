package FileBox

import (
	. "wechatBot/General"

	"github.com/beevik/etree"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	if message.Type() != schemas.MessageTypeUnknown && message.Type() != schemas.MessageTypeAttachment {
		log.Printf("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() == schemas.MessageTypeRecalled {
		log.Infof("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		return
	}

	/* TODO MessageType
	MessageTypeUnknown
	MessageTypeAttachment
	MessageTypeAudio
	MessageTypeContact
	MessageTypeChatHistory
	MessageTypeEmoticon
	MessageTypeImage
	MessageTypeText
	MessageTypeLocation
	MessageTypeMiniProgram
	MessageTypeGroupNote
	MessageTypeTransfer
	MessageTypeRedEnvelope
	MessageTypeRecalled
	MessageTypeURL
	MessageTypeVideo
	*/
	if message.Type() == schemas.MessageTypeUnknown && message.Talker().Name() == "微信团队" {
		log.Infof("Type Pass, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
		return
	}
	log.Infof("FileBox, Type: [%v]:[%v] CoptRight: [%s]", message.Type().String(), message.Talker().Name(), Copyright(make([]uintptr, 1)))
	fileType, fileName := FileType(message)
	switch fileType {
	case "pdf":
		log.Infof("[fileType:%v] [fileName:%v] CoptRight: [%s]", fileType, fileName, Copyright(make([]uintptr, 1)))
	case "rar|zip|tar|gz":
		log.Infof("[fileType:%v] [fileName:%v] CoptRight: [%s]", fileType, fileName, Copyright(make([]uintptr, 1)))
	default:
		log.Infof("[fileType:%v] [fileName:%v] CoptRight: [%s]", fileType, fileName, Copyright(make([]uintptr, 1)))
	}
}

func FileType(message *user.Message) (string, string) {
	fileType := ""
	fileName := ""
	doc := etree.NewDocument()
	if err = doc.ReadFromString(message.MentionText()); err != nil {
		log.Errorf("FileType Error: [%v] CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	}
	for _, t := range doc.FindElements("//fileext") {
		fileType = t.Text()
	}
	for _, t := range doc.FindElements("//title") {
		fileName = t.Text()
	}
	return fileType, fileName
}

func FileBoxPDF(message *user.Message) {

}
