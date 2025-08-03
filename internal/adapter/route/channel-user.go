package route

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/DESA/internal/adapter/controller"
	middlewareJWT "github.com/wang900115/DESA/lib/common/middleware/jwt"
	middlewareRole "github.com/wang900115/DESA/lib/common/middleware/role"

	"github.com/wang900115/DESA/lib/common/router"
)

type ChannelUserRouter struct {
	channelUser controller.ChannelUserController
	jwt         middlewareJWT.JWT
	role        middlewareRole.Permission
}

func NewChannelUserRouter(channelUser *controller.ChannelUserController, jwt *middlewareJWT.JWT, role *middlewareRole.Permission) router.IRoute {
	return &ChannelUserRouter{channelUser: *channelUser, jwt: *jwt, role: *role}
}

func (cu *ChannelUserRouter) Setup(router *gin.RouterGroup) {
	channelUserGroup := router.Group("v1/channel-user/")
	{
		channelUserGroup.POST("/join", cu.jwt.Middleware, cu.channelUser.JoinChannel)   // 用戶進去該頻道 (連通 websocket 連線該節點)
		channelUserGroup.POST("/query", cu.channelUser.GetChannelUsers)                 // 獲取該頻道的用戶
		channelUserGroup.POST("/leave", cu.jwt.Middleware, cu.channelUser.LeaveChannel) // 用戶離開該頻道 (斷開 websocket 斷開該節點)
	}
}
