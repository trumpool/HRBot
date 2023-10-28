package controller

import (
	"xlab-feishu-robot/internal/store"
)

func HelpP2P(messageEvent *store.MessageEvent) {
	helpMessage := `输入"help"查看帮助，发送“开始使用”以开始使用机器人。
	(如果您距离上一次使用已经过去了一段时间，请发送"开始使用"以重新认证！)
	功能1: 批量加人
		用法："批量加人. 张三, 李四, 王五. 答疑群, 水群, 正式群"
		注意事项：
			1. 机器人使用"."和","分割输入文本，标点符号为英文标点。
			2. 要想拉人进某个群，必须先把机器人拉进该群。
	功能2: 批量删人
		用法："批量加人. 张三, 李四, 王五. 答疑群, 水群, 正式群"
		注意事项：
			1. 机器人使用"."和","分割输入文本，标点符号为英文标点。
			2. 机器人需要在群组中且操作者需要有删除成员的权限。
	`
	SendMessage(messageEvent.Sender.Sender_id.Open_id, helpMessage)
}
