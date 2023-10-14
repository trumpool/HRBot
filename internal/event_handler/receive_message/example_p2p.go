package receiveMessage

import (
	_ "xlab-feishu-robot/internal/pkg"

	_ "github.com/sirupsen/logrus"
)

func init() {
	p2pMessageRegister(p2pHelpMenu, "help")
}

func p2pHelpMenu(messageevent *MessageEvent) {
}
