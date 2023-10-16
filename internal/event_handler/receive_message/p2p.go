package receiveMessage

import (
	"github.com/sirupsen/logrus"
	"strings"
	"xlab-feishu-robot/internal/controller"
)

func p2p(messageevent *MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		p2pTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive p2p message, but this type is not supported")
	}
}

func p2pTextMessage(messageevent *MessageEvent) {
	// get the pure text message
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive p2p TEXT message")
	content := messageevent.Message.Content
	switch {
	case strings.Contains(content, "批量加人"):
		if !hasPermission(messageevent) {
			logrus.Warn("Receive p2p TEXT message, but the sender does not have permission")
			return
		}
		controller.AddPeople(content)
	default:
		logrus.Errorf("Receive p2p TEXT message, but this type is not supported")
	}
}
