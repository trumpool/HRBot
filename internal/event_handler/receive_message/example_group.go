package receiveMessage

import (
	_ "xlab-feishu-robot/internal/pkg"

	_ "github.com/sirupsen/logrus"
)

func init() {
	groupMessageRegister(groupHelpMenu, "help")
}

func groupHelpMenu(messageevent *MessageEvent) {
}
