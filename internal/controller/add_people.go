package controller

import (
	"context"
	"encoding/json"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"strings"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/pkg"
	"xlab-feishu-robot/internal/store"
)

func AddPeople(messageEvent *store.MessageEvent) {
	// 检查权限
	if !HasPermission(messageEvent) {
		logrus.Warn("No permission")
		return
	}
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
	dataRecord, err := inviteUserToGroupChat(peopleID, groupsID)
	if err != nil {
		logrus.Error(err)
		return
	}

	checkInviteResult(dataRecord, messageEvent)
}

func inviteUserToGroupChat(peopleID []string, groupsID []string) ([]*larkim.CreateChatMembersRespData, error) {
	dataRecord := make([]*larkim.CreateChatMembersRespData, 0)
	for _, groupID := range groupsID {
		// 创建请求对象
		req := larkim.NewCreateChatMembersReqBuilder().
			ChatId(groupID).
			MemberIdType("open_id").
			// 将参数中可用的 ID 全部拉入群聊，返回拉群成功的响应，并展示剩余不可用的 ID 及原因
			SucceedType(1).
			Body(larkim.NewCreateChatMembersReqBodyBuilder().
				IdList(peopleID).
				Build()).
			Build()
		// 发起请求
		resp, err := pkg.Client.Im.ChatMembers.Create(context.Background(), req)
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

// checkInviteResult 检查邀请是否成功，如果失败则向用户发送失败信息
func checkInviteResult(dataRecord []*larkim.CreateChatMembersRespData, messageEvent *store.MessageEvent) {
	invalidIDList := make([]string, 0)
	notExistedIDList := make([]string, 0)
	for _, data := range dataRecord {
		for _, v := range data.InvalidIdList {
			invalidIDList = append(invalidIDList, v)
		}
		for _, v := range data.NotExistedIdList {
			notExistedIDList = append(notExistedIDList, v)
		}
	}

	failedMessage := "以下用户未被邀请成功：\n"
	failedMessage += "无效的ID：\n"
	for _, v := range invalidIDList {
		failedMessage += fmt.Sprintf("%s\n", v)
	}
	failedMessage += "不存在的ID：\n"
	for _, v := range notExistedIDList {
		failedMessage += fmt.Sprintf("%s\n", v)
	}

	failedMessage += "请联系机器人管理员，将您的输入和错误信息一起反馈，谢谢！"

	msgContent := map[string]interface{}{
		"text": failedMessage,
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

func getPeopleID(wantedPeople []string) ([]string, error) {
	allPeople, err := getAllPeopleInDepartment()
	if err != nil {
		return nil, err
	}
	var result []string

	wantedPeopleMap := make(map[string]bool)
	for _, v := range wantedPeople {
		wantedPeopleMap[v] = true
	}

	for _, v := range allPeople {
		if wantedPeopleMap[*v.Name] {
			result = append(result, *v.OpenId)
		}
	}
	return result, nil
}

func getGroupsID(wantedGroup []string) ([]string, error) {
	allGroups, err := getBotGroupList()
	if err != nil {
		return nil, err
	}
	var result []string

	wantedGroupMap := make(map[string]bool)
	for _, v := range wantedGroup {
		wantedGroupMap[v] = true
	}

	for _, v := range allGroups {
		if wantedGroupMap[*v.Name] {
			result = append(result, *v.ChatId)
		}
	}
	return result, nil
}

func getAllPeopleInDepartment() ([]*larkcontact.User, error) {
	// 创建请求对象
	req := larkcontact.NewFindByDepartmentUserReqBuilder().
		UserIdType("open_id").
		DepartmentIdType(config.C.DepartmentIdType).
		DepartmentId(config.C.DepartmentID).
		Build()
	// 发起请求
	tenantAccessToken, err := GetTenantAccessToken()
	if err != nil {
		return nil, err
	}
	resp, err := pkg.Client.Contact.User.FindByDepartment(context.Background(), req, larkcore.WithTenantAccessToken(tenantAccessToken))
	// 处理错误
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
	}

	result := resp.Data.Items
	for *resp.Data.HasMore {
		req = larkcontact.NewFindByDepartmentUserReqBuilder().
			UserIdType("open_id").
			DepartmentIdType(config.C.DepartmentIdType).
			DepartmentId(config.C.DepartmentID).
			PageToken(*resp.Data.PageToken).
			Build()
		resp, err = pkg.Client.Contact.User.FindByDepartment(context.Background(), req, larkcore.WithTenantAccessToken(tenantAccessToken))
		if err != nil {
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
		}

		result = append(result, resp.Data.Items...)
	}

	return result, nil
}

// getBotGroupList 获取机器人所在的所有群
func getBotGroupList() ([]*larkim.ListChat, error) {
	tenantAccessToken, err := GetTenantAccessToken()
	if err != nil {
		return nil, err
	}

	req := larkim.NewListChatReqBuilder().
		Build()
	resp, err := pkg.Client.Im.Chat.List(context.Background(), req, larkcore.WithTenantAccessToken(tenantAccessToken))
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
	}

	result := resp.Data.Items
	for *resp.Data.HasMore {
		req = larkim.NewListChatReqBuilder().
			PageToken(*resp.Data.PageToken).
			Build()
		resp, err = pkg.Client.Im.Chat.List(context.Background(), req, larkcore.WithTenantAccessToken(tenantAccessToken))
		if err != nil {
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
		}

		result = append(result, resp.Data.Items...)
	}

	return result, nil
}

func parsePeopleAndGroup(content string) (people []string, group []string) {
	// content格式：批量加人. 张三, 李四, 王五. 推送群, 答疑群, 交流群
	// 1. 以.分割
	tmp := strings.Split(content, ".")
	peopleStr, groupStr := tmp[1], tmp[2]
	// 2. 以,分割
	people = strings.Split(peopleStr, ",")
	group = strings.Split(groupStr, ",")

	// 3. 去除空格
	for i := 0; i < len(people); i++ {
		people[i] = strings.TrimSpace(people[i])
	}
	for i := 0; i < len(group); i++ {
		group[i] = strings.TrimSpace(group[i])
	}

	return
}
