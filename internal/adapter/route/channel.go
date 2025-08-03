package route

import (
	"github.com/wang900115/DESA/internal/adapter/controller"
	middlewareRole "github.com/wang900115/DESA/lib/common/middleware/role"
	"github.com/wang900115/DESA/lib/common/router"

	"github.com/gin-gonic/gin"
)

type ChannelRouter struct {
	channel controller.ChannelController
	role    middlewareRole.Permission
}

func NewChannelRouter(channel *controller.ChannelController) router.IRoute {
	return &ChannelRouter{channel: *channel}
}

func (cr *ChannelRouter) Setup(router *gin.RouterGroup) {
	channelGroup := router.Group("v1/channel/")
	{
		channelGroup.POST("/create", cr.role.RequireRoles("admin"), cr.channel.Create) // 新建頻道
		channelGroup.POST("/query", cr.channel.GetAllChannels)                         // 獲取頻道
		channelGroup.POST("/delete", cr.channel.Delete)                                // 刪除頻道
	}
}
