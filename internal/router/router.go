package router

import (
	"xlab-feishu-robot/internal/controller"
	"xlab-feishu-robot/internal/dispatcher"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
	r.GET("/feiShu/GetUserAccessToken", controller.GetCodeThenGetUserAccessToken)
}
