package controller

import (
	"context"
	"encoding/json"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
	"xlab-feishu-robot/internal/store"
)

func HelpP2P(messageEvent *store.MessageEvent) {
	helpMessage := `输入"help"查看帮助

	功能1: 批量加人
		用法："批量加人. 张三, 李四, 王五. 答疑群, 水群, 正式群"
	`
	msgContent := map[string]interface{}{
		"text": helpMessage,
	}
	// help JSON message
	msgContentJSON, err := json.Marshal(msgContent)
	if err != nil {
		logrus.Error(err)
		return
	}
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(messageEvent.Sender.Sender_id.Open_id).
			MsgType("text").
			Content(string(msgContentJSON)).
			Build()).
		Build()

	resp, err := pkg.Client.Im.Message.Create(context.Background(), req)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		logrus.Error(resp.Code, resp.Msg)
		return
	}
}
