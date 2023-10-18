package controller

import (
	"context"
	"encoding/json"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
)

// SendMessage sends message to receiver
// receiverID: open_id of receiver
func SendMessage(receiverID string, message string) {
	msgContent := map[string]interface{}{
		"text": message,
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
			ReceiveId(receiverID).
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
