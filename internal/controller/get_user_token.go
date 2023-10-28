package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
)

// 全局变量userAccessToken
// map[user_open_id]userAccessToken
var userAccessTokenMap map[string]string = make(map[string]string)

func GetCodeThenGetUserAccessToken(c *gin.Context) {
	code := c.Query("code")
	user_open_id := c.Query("state")
	fmt.Println(code)
	//拿userAccessToken
	req := larkauthen.NewCreateAccessTokenReqBuilder().
		Body(larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType("authorization_code").
			Code(code).
			Build()).
		Build()
	resp, err := pkg.Client.Authen.AccessToken.Create(context.Background(), req)
	SetUserAccessToken(user_open_id, *resp.Data.AccessToken)
	if err != nil {
		logrus.Error("Cannot Get User Access Token ", req)
		return
	} else {
		SendMessage(user_open_id, "登录授权成功, User access token: "+*resp.Data.AccessToken)
	}
	return
}

func SetUserAccessToken(openID string, userAccessToken string) {
	userAccessTokenMap[openID] = userAccessToken

}

func GetUserAccessToken(openID string) (token string, err error) {
	if token, ok := userAccessTokenMap[openID]; ok {
		return token, nil
	} else {
		logrus.Error("Cannot Get User Access Token ", err)
		return "", errors.New("the User Access Token is not exist")
	}
}
