package controller

import (
	"context"
	"fmt"
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
	foundPeopleMap, err := getPeopleID(people)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("foundPeopleMap:%v", foundPeopleMap)

	// 获得所有群的ID
	foundGroupMap, err := getGroupsID(group)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("foundGroupMap:%v", foundGroupMap)

	// 检查是否有未找到的ID
	checkAllIDFound(foundPeopleMap, foundGroupMap, people, group, messageEvent)

	// 将所有人加入所有群
	dataRecord, err := inviteUserToGroupChat(foundPeopleMap, foundGroupMap)
	if err != nil {
		logrus.Error(err)
		return
	}

	checkInviteResult(dataRecord, messageEvent)
}

func inviteUserToGroupChat(peopleMap map[string]string, groupsMap map[string]string) ([]*larkim.CreateChatMembersRespData, error) {
	dataRecord := make([]*larkim.CreateChatMembersRespData, 0)
	IDList := make([]string, 0)
	for _, v := range peopleMap {
		IDList = append(IDList, v)
	}

	for _, groupID := range groupsMap {
		// 创建请求对象
		req := larkim.NewCreateChatMembersReqBuilder().
			ChatId(groupID).
			MemberIdType("open_id").
			// 将参数中可用的 ID 全部拉入群聊，返回拉群成功的响应，并展示剩余不可用的 ID 及原因
			SucceedType(1).
			Body(larkim.NewCreateChatMembersReqBodyBuilder().
				IdList(IDList).
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
		invalidIDList = append(invalidIDList, data.InvalidIdList...)
		notExistedIDList = append(notExistedIDList, data.NotExistedIdList...)
	}

	message := "以下用户未被邀请成功：\n"
	message += "无效的ID：\n"
	for _, v := range invalidIDList {
		message += fmt.Sprintf("%s\n", v)
	}
	message += "不存在的ID：\n"
	for _, v := range notExistedIDList {
		message += fmt.Sprintf("%s\n", v)
	}

	message += "请联系机器人管理员，将您的输入和错误信息一起反馈，谢谢！"

	if len(invalidIDList) == 0 && len(notExistedIDList) == 0 {
		message = "所有用户均已成功加入群聊！"
	}

	SendMessage(messageEvent.Sender.Sender_id.Open_id, message)
}

func checkAllIDFound(foundPeopleMap map[string]string, foundGroupMap map[string]string, peopleNameList []string, groupNameList []string, messageEvent *store.MessageEvent) {
	message := ""
	for _, v := range peopleNameList {
		if _, ok := foundPeopleMap[v]; !ok {
			message += fmt.Sprintf("未找到用户：%s\n", v)
		}
	}

	for _, v := range groupNameList {
		if _, ok := foundGroupMap[v]; !ok {
			message += fmt.Sprintf("未找到群：%s\n", v)
		}
	}

	if message != "" {
		SendMessage(messageEvent.Sender.Sender_id.Open_id, message)
	}
}

func getPeopleID(wantedPeople []string) (map[string]string, error) {
	allPeople, err := getAllPeopleInDepartment()
	if err != nil {
		return nil, err
	}

	wantedPeopleMap := make(map[string]bool)
	for _, v := range wantedPeople {
		wantedPeopleMap[v] = true
	}

	foundPeopleMap := make(map[string]string)
	for _, v := range allPeople {
		if wantedPeopleMap[*v.Name] {
			foundPeopleMap[*v.Name] = *v.OpenId
		}
	}
	return foundPeopleMap, nil
}

func getGroupsID(wantedGroup []string) (map[string]string, error) {
	allGroups, err := getBotGroupList()
	if err != nil {
		return nil, err
	}

	wantedGroupMap := make(map[string]bool)
	for _, v := range wantedGroup {
		wantedGroupMap[v] = true
	}

	foundGroupMap := make(map[string]string)
	for _, v := range allGroups {
		if wantedGroupMap[*v.Name] {
			foundGroupMap[*v.Name] = *v.ChatId
		}
	}
	return foundGroupMap, nil
}

func getAllPeopleInDepartment() ([]*larkcontact.User, error) {
	// 创建请求对象
	req := larkcontact.NewFindByDepartmentUserReqBuilder().
		UserIdType("open_id").
		DepartmentIdType(config.C.DepartmentIdType).
		DepartmentId(config.C.DepartmentID).
		Build()
	// 发起请求
	resp, err := pkg.Client.Contact.User.FindByDepartment(context.Background(), req)
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
		resp, err = pkg.Client.Contact.User.FindByDepartment(context.Background(), req)
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
	req := larkim.NewListChatReqBuilder().
		Build()
	resp, err := pkg.Client.Im.Chat.List(context.Background(), req)
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
		resp, err = pkg.Client.Im.Chat.List(context.Background(), req)
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
