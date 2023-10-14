package controller

import (
	"context"
	"encoding/json"
	"fmt"
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
			AppSecret(config.C.Feishu.AppSecret).
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

	var result map[string]interface{}
	err = json.Unmarshal(resp.ApiResp.RawBody, &result)
	if err != nil {
		logrus.Error("GetTenantAccessToken failed")
		return "", err
	}
	logrus.Info("GetTenantAccessToken success")
	return result["tenant_access_token"].(string), nil
}
