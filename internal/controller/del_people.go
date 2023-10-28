package controller

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
	"xlab-feishu-robot/internal/store"
)

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
	if err != nil || dataRecord == nil {
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
		userAccessToken, err := GetUserAccessToken(receiverID)
		if err != nil || userAccessToken == "" {
			SendMessage(receiverID, fmt.Sprintf("您尚未登陆， 发送 开始使用 以登录"))
			return nil, err
		}
		resp, err := pkg.Client.Im.ChatMembers.Delete(context.Background(), req, larkcore.WithUserAccessToken(userAccessToken))
		// 处理错误
		if err != nil {
			SendMessage(receiverID, fmt.Sprintf("机器人发送消息时错误，错误信息：%s, response: %v", err.Error(), resp))
		} else if !resp.Success() {
			SendMessage(receiverID, fmt.Sprintf("删除成员失败，错误信息：%d, response: %v", resp.Code, resp.Msg))
		} else {
			SendMessage(receiverID, fmt.Sprintf("删除结束"))
		}

		dataRecord = append(dataRecord, resp.Data)

	}

	return dataRecord, nil
}

// checkInviteResult 检查删除是否成功，如果失败则向用户发送失败信息
func checkDeleteResult(dataRecord []*larkim.DeleteChatMembersRespData, messageEvent *store.MessageEvent) {
	invalidIDList := make([]string, 0)
	for _, data := range dataRecord {
		if (data == nil) || (data.InvalidIdList == nil) {
			continue
		}
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
		//message = "所有用户均已成功从群聊中删除！"
		//这里不能用invalidIDList来判断，如果有错误在前面就会返回，这里只处理非法ID的情况
		return
	}

	SendMessage(messageEvent.Sender.Sender_id.Open_id, message)

}
