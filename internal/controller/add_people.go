package controller

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
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

func GetAllPeopleInDepartment() ([]*larkcontact.User, error) {
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

	// 服务端错误处理
	if !resp.Success() {
		return nil, fmt.Errorf("resp failed, code:%d, msg:%s", resp.Code, resp.Msg)
	}

	return resp.Data.Items, nil
}
