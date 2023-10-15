package router

import (
	"xlab-feishu-robot/internal/dispatcher"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
}
