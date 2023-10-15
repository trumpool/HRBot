package controller

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/sirupsen/logrus"
	"strings"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/pkg"
)

func AddPeople(content string) {
	people, group := parsePeopleAndGroup(content)
	logrus.Infof("people:%v, group:%v", people, group)
	// 获取所有人的ID

	// 获得所有群的ID

	// 将所有人加入所有群

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
		result = append(result, resp.Data.Items...)
	}

	// 服务端错误处理
	if !resp.Success() {
		return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
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
