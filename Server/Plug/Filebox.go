package Plug

import (
	"github.com/beevik/etree"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func FileBox(message *user.Message) {
	if message.Type() != schemas.MessageTypeUnknown && message.Type() != schemas.MessageTypeAttachment {
		log.Printf("Type Pass, Type: [%v]:[%v]", message.Type().String(), message.Talker().Name())
		return
	}
	fileType, fileName := FileType(message)
	switch fileType {
	case "pdf":
		log.Printf(fileType)
	case "rar|zip|tar|gz":
		log.Printf(fileType)
	default:
		log.Printf(fileType)
	}
}

func FileType(message *user.Message) string, string {
	fileType := ""
	fileName := ""
	doc := etree.NewDocument()
	if err := doc.ReadFromString(message.MentionText()); err != nil {
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
