package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"xlab-feishu-robot/internal/config"
	"xlab-feishu-robot/internal/log"
)

func TestGetPeopleID(t *testing.T) {
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

	IDs, err := getPeopleID([]string{"牛马", "鼠鼠"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(IDs))
	for _, v := range IDs {
		logrus.Info(v)
	}
}
