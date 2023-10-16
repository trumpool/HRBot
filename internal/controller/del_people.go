package controller

import (
	"context"
	"encoding/json"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/internal/pkg"
	"xlab-feishu-robot/internal/store"
)

//todo: this part has not been tested yet!

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

	// 将所有人加入所有群
	dataRecord, err := deleteUserInGroupChat(peopleID, groupsID)
	if err != nil {
		logrus.Error(err)
		return
	}

	checkDeleteResult(dataRecord, messageEvent)
}
func deleteUserInGroupChat(peopleID []string, groupsID []string) ([]*larkim.DeleteChatMembersRespData, error) {
	dataRecord := make([]*larkim.DeleteChatMembersRespData, 0)
	for _, groupID := range groupsID {
		// 创建请求对象
		req := larkim.NewDeleteChatMembersReqBuilder().
			ChatId(groupID).
			MemberIdType("open_id").
			Body(larkim.NewDeleteChatMembersReqBodyBuilder().
				IdList(peopleID).
				Build()).
			Build()

		resp, err := pkg.Client.Im.ChatMembers.Delete(context.Background(), req)
		// 处理错误
		if err != nil {
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
		}

		dataRecord = append(dataRecord, resp.Data)
	}

	return dataRecord, nil
}

// checkInviteResult 检查删除是否成功，如果失败则向用户发送失败信息
func checkDeleteResult(dataRecord []*larkim.DeleteChatMembersRespData, messageEvent *store.MessageEvent) {
	invalidIDList := make([]string, 0)
	//notExistedIDList := make([]string, 0)
	for _, data := range dataRecord {
		invalidIDList = append(invalidIDList, data.InvalidIdList...)
		//only returns invalid id
		//notExistedIDList = append(notExistedIDList, data.InvalidIdList...)
	}

	message := "以下用户未被邀请成功：\n"
	message += "无效的ID：\n"
	for _, v := range invalidIDList {
		message += fmt.Sprintf("%s\n", v)
	}

	message += "请联系机器人管理员，将您的输入和错误信息一起反馈，谢谢！"

	if len(invalidIDList) == 0 {
		message = "所有用户均已成功加入群聊！"
	}

	msgContent := map[string]interface{}{
		"text": message,
	}
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
