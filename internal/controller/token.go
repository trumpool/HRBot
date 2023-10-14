package controller

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauth "github.com/larksuite/oapi-sdk-go/v3/service/auth/v3"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/pkg"
)

func GetTenantAccessToken() (string, error) {
	// 创建请求对象
	req := larkauth.NewInternalTenantAccessTokenReqBuilder().
		Body(larkauth.NewInternalTenantAccessTokenReqBodyBuilder().
			AppId(config.C.Feishu.AppId).
			AppSecret(config.C.Feishu.AppId).
			Build()).
		Build()
	// 发起请求
	resp, err := pkg.Client.Auth.TenantAccessToken.Internal(context.Background(), req)

	// 处理错误
	if err != nil {
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		return "", fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
	}

	logrus.Info("GetTenantAccessToken success")
	return larkcore.Prettify(resp), nil
}
