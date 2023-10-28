package controller

import (
	"fmt"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/store"
)

/*
Login is used to send a login link to the user.
*/
func Login(messageEvent *store.MessageEvent) {
	redirectUrl := config.C.RedirectUrl
	appID := config.C.Feishu.AppId
	loginLink := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s", redirectUrl, appID, messageEvent.Sender.Sender_id.Open_id)
	SendMessage(messageEvent.Sender.Sender_id.Open_id, "请点击以下链接进行登录：\n"+loginLink)
	return
}
