package controller

import (
	"context"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/log"
	"xlab-feishu-robot/internal/pkg"
)

func Test_getAllPeopleInDepartment(t *testing.T) {
	setupForTest()
	people, err := getAllPeopleInDepartment()
	assert.NoError(t, err)
	for _, v := range people {
		logrus.Info(*v.Name)
	}
}

func Test_getPeopleID(t *testing.T) {
	setupForTest()
	IDs, err := getPeopleID([]string{"牛马", "鼠鼠"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(IDs))
	for _, v := range IDs {
		logrus.Info(v)
	}
}

func Test_getBotGroupList(t *testing.T) {
	setupForTest()
	groups, err := getBotGroupList()
	assert.NoError(t, err)
	for _, v := range groups {
		logrus.Info(*v.Name)
	}
}

func Test_getGroupsID(t *testing.T) {
	setupForTest()
	IDs, err := getGroupsID([]string{"测试知识树提醒"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(IDs))
	for _, v := range IDs {
		logrus.Info(v)
	}
}

func Test_inviteUserToGroupChat(t *testing.T) {
	setupForTest()
	_, err := inviteUserToGroupChat([]string{"123"}, []string{"123"})
	assert.NoError(t, err)
}

func Test_SearchEmployee(t *testing.T) {
	setupForTest()

	tenantAccessToken, err := GetTenantAccessToken()
	assert.NoError(t, err)
	logrus.Info(tenantAccessToken)

	resp, err := pkg.Client.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod: http.MethodPost,
			ApiPath:    "https://open.feishu.cn/open-apis/corehr/v2/employees/search",
			Body: map[string]interface{}{
				"key_word": "牛马",
			},
			QueryParams: larkcore.QueryParams{
				"page_size":    []string{"10"},
				"user_id_type": []string{"open_id"},
			},
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
		},
	)

	if err != nil {
		logrus.Errorf("failed to search employee, err: %v", err)
		return
	}

	logrus.Info(resp)
}

func setupForTest() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../config/")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	if err := viper.Unmarshal(&config.C); err != nil {
		logrus.Error("Failed to unmarshal config")
	}

	logrus.Info("Configuration file loaded")

	// log
	log.SetupLogrus()
	logrus.Info("Robot starts up")

	// feishu api client
	config.SetupFeishuApiClient()
}
