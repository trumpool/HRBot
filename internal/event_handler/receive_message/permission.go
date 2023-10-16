package receiveMessage

import "xlab-feishu-robot/internal/config"

func hasPermission(messageevent *MessageEvent) bool {
	// 遍历白名单
	for _, id := range config.C.WhiteList {
		if id == messageevent.Sender.Sender_id.Open_id {
			return true
		}
	}
	return false
}
