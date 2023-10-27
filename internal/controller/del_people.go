package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
	"xlab-feishu-robot/internal/store"
)

var userAccessToken string

type CodeResponse struct {
	RedirectRrl string
	Code        string
	State       string
}

// todo: this part has not been tested yet!
func Login(messageEvent *store.MessageEvent) {
	redirectUrl := "https://81762kq506.goho.co/feiShu/GetUserAccessToken"
	appID := "cli_a591c1b02279900e"
	//loginLink := "https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=https://81762kq506.goho.co/feiShu/Event&app_id=cli_a591c1b02279900e"
	loginLink := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s", redirectUrl, appID)
	SendMessage(messageEvent.Sender.Sender_id.Open_id, "请点击以下链接进行登录：\n"+loginLink)

	return
}

func GetCodeThenGetUserAccessToken(c *gin.Context) {
	code := c.Query("code")
	fmt.Println(code)
	// setp1 get app_access_token
	appID := "cli_a591c1b02279900e"
	appSecret := "3Qrm0ChE6c71gpaOpxEqJgD26GiArMIq"
	client := lark.NewClient(appID, appSecret)
	req := larkauthen.NewCreateAccessTokenReqBuilder().
		Body(larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType("authorization_code").
			Code(code).
			Build()).
		Build()
	// 发起请求
	resp, err := client.Authen.AccessToken.Create(context.Background(), req)
	if err != nil {
		logrus.Error("Cannot Get User Access Token ", req)
		return
	}
	userAccessToken = *resp.Data.AccessToken
	fmt.Println(userAccessToken)
	return
}
func DelPeople(messageEvent *store.MessageEvent) {
	// 检查权限
	if !HasPermission(messageEvent) {
		logrus.Warn("No permission")
		return
	}
	//预处理人名 组名
	people, group := parsePeopleAndGroup(messageEvent.Message.Content)
	logrus.Infof("people:%v, group:%v", people, group)
	// 获取所有人的ID
	peopleID, err := getPeopleID(people)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("peopleID:%v", peopleID)

	// 获得所有群的ID
	groupsID, err := getGroupsID(group)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("groupsID:%v", groupsID)

	checkAllIDFound(peopleID, groupsID, people, group, messageEvent)

	checkWhetherBotInGroup(groupsID, messageEvent)

	// 删人
	dataRecord, err := deleteUserInGroupChat(peopleID, groupsID, messageEvent.Sender.Sender_id.Open_id)
	if err != nil {
		logrus.Error(err)
		return
	}

	checkDeleteResult(dataRecord, messageEvent)
}
func deleteUserInGroupChat(peopleMap map[string]string, groupsMap map[string]string, receiverID string) ([]*larkim.DeleteChatMembersRespData, error) {
	dataRecord := make([]*larkim.DeleteChatMembersRespData, 0)
	IDList := make([]string, 0)
	for _, v := range peopleMap {
		IDList = append(IDList, v)
	}

	for _, groupID := range groupsMap {
		// 创建请求对象
		req := larkim.NewDeleteChatMembersReqBuilder().
			ChatId(groupID).
			MemberIdType("open_id").
			Body(larkim.NewDeleteChatMembersReqBodyBuilder().
				IdList(IDList).
				Build()).
			Build()
		// 发起请求
		//userAccessToken := "empty"
		//userAccessToken = getUserAccessToken()
		resp, err := pkg.Client.Im.ChatMembers.Delete(context.Background(), req, larkcore.WithUserAccessToken(userAccessToken))
		// 处理错误
		if err != nil || !resp.Success() {
			SendMessage(receiverID, fmt.Sprintf("邀请失败，错误信息：%s, response: %v", err.Error(), resp))
		}

		dataRecord = append(dataRecord, resp.Data)

		fmt.Println(larkcore.Prettify(resp))
	}

	return dataRecord, nil
}

// checkInviteResult 检查删除是否成功，如果失败则向用户发送失败信息
func checkDeleteResult(dataRecord []*larkim.DeleteChatMembersRespData, messageEvent *store.MessageEvent) {
	invalidIDList := make([]string, 0)
	for _, data := range dataRecord {
		invalidIDList = append(invalidIDList, data.InvalidIdList...)
		//only returns invalid id
	}

	message := "以下用户未被删除成功：\n"
	message += "无效的ID：\n"
	for _, v := range invalidIDList {
		message += fmt.Sprintf("%s\n", v)
	}

	message += "请联系机器人管理员，将您的输入和错误信息一起反馈，谢谢！"

	if len(invalidIDList) == 0 {
		message = "所有用户均已成功从群聊中删除！"
	}

	SendMessage(messageEvent.Sender.Sender_id.Open_id, message)

}

func getUserAccessToken() (userAccessToken string) {
	return
}
