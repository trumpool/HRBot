package receiveMessage

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/store"
)

// dispatch message, according to Chat type
func Receive(event map[string]any) {
	messageevent := store.MessageEvent{}
	map2struct(event, &messageevent)
	switch messageevent.Message.Chat_type {
	case "p2p":
		p2p(&messageevent)
	case "group":
		group(&messageevent)
	default:
		logrus.WithFields(logrus.Fields{"chat type": messageevent.Message.Chat_type}).Warn("Receive message, but this chat type is not supported")
	}
}

func map2struct(m map[string]interface{}, stru interface{}) {
	bytes, _ := json.Marshal(m)
	_ = json.Unmarshal(bytes, stru)
}
