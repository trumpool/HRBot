package receiveMessage

import (
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/store"
)

func hasPermission(messageevent *store.MessageEvent) bool {
	// 遍历白名单
	for _, id := range config.C.WhiteList {
		if id == messageevent.Sender.Sender_id.Open_id {
			return true
		}
	}
	return false
}
