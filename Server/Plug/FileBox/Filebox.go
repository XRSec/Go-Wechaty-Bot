package FileBox

import (
	"github.com/beevik/etree"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

func FileBox() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	if message.Type() != schemas.MessageTypeUnknown && message.Type() != schemas.MessageTypeAttachment {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	if message.Type() == schemas.MessageTypeRecalled {
		log.Errorf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
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
		log.Errorf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	log.Infof("FileBox, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
	fileType, fileName := FileType(message)
	switch fileType {
	case "pdf":
		log.Infof("[fileType:%v] [fileName:%v]", fileType, fileName)
	case "rar|zip|tar|gz":
		log.Infof("[fileType:%v] [fileName:%v]", fileType, fileName)
	default:
		log.Infof("[fileType:%v] [fileName:%v]", fileType, fileName)
	}
}

func FileType(message *user.Message) (string, string) {
	fileType := ""
	fileName := ""
	doc := etree.NewDocument()
	if err = doc.ReadFromString(message.MentionText()); err != nil {
		log.Errorf("FileType Error: [%v]", err)
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
