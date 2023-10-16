package receiveMessage

import (
	"github.com/sirupsen/logrus"
	"strings"
	"xlab-feishu-robot/internal/store"
)

func group(messageevent *store.MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		groupTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive group message, but this type is not supported")
	}
}

func groupTextMessage(messageevent *store.MessageEvent) {
	// get the pure text message, without @xxx
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	messageevent.Message.Content = messageevent.Message.Content[strings.Index(messageevent.Message.Content, " ")+1:]
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive group TEXT message")
}
