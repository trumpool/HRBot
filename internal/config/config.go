package config

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xlab-feishu-robot/internal/pkg"
)

type Config struct {
	Feishu FeishuConfig
	Server struct {
		Port int

		// add your configuration fields here
		ExampleField1 string
	}

	DepartmentID     string
	DepartmentIdType string

	WhiteList []string
}

var C Config

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		logrus.Error("Failed to unmarshal config")
	}

	logrus.Info("Configuration file loaded")
}

func SetupFeishuApiClient() {
	pkg.Client = lark.NewClient(C.Feishu.AppId, C.Feishu.AppSecret, lark.WithEnableTokenCache(true))
}

type FeishuConfig struct {
	AppId             string
	AppSecret         string
	VerificationToken string
	EncryptKey        string
}
