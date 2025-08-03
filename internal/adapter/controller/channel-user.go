package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/DESA/internal/adapter/validator"
	"github.com/wang900115/DESA/internal/application/usecase"
	"github.com/wang900115/DESA/lib/common"
	iresponse "github.com/wang900115/DESA/lib/common/response"
)

type ChannelUserController struct {
	channel usecase.ChannelUsecase
	user    usecase.UserUsecase
	p2p     usecase.P2PUsecase
	resp    iresponse.IResponse
}

func NewChannelUserController(channel *usecase.ChannelUsecase, p2p *usecase.P2PUsecase, resp iresponse.IResponse) *ChannelUserController {
	return &ChannelUserController{channel: *channel, p2p: *p2p, resp: resp}
}

// 獲取該頻道的上線用戶資訊
func (cu *ChannelUserController) GetChannelUsers(c *gin.Context) {
	var req validator.GetChannelUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cu.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	users, err := cu.channel.GetChannelUsers(c, req.ChannelID)
	if err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cu.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"users": users,
	})
}

// 用戶註冊至該頻道
func (cu *ChannelUserController) RegisterToChannel(c *gin.Context) {
	userID := c.GetString("user_id")
	var req validator.RegisterToChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cu.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	if err := cu.p2p.Register(req.ChannelID, userID, req.PubKey); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	if err := cu.user.Register(c, userID, req.ChannelID); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	cu.resp.Success(c, common.REGISTER_SUCCESS)
}

// 用戶連線至該頻道
func (cu *ChannelUserController) JoinChannel(c *gin.Context) {
	peerID := c.GetString("user_id")
	var req validator.JoinChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cu.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	valid, err := cu.p2p.Verify(req.ChannelID, peerID, req.Message, req.Signature)
	if !valid || err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	if err := cu.p2p.Connect(c, peerID, req.ChannelID, req.URL); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	if err := cu.user.JoinChannel(c, peerID, req.ChannelID); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cu.resp.Success(c, common.JOIN_SUCCESS)
}

// 用戶從頻道斷線
func (cu *ChannelUserController) LeaveChannel(c *gin.Context) {
	peerID := c.GetString("user_id")
	var req validator.LeaveChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cu.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	if err := cu.p2p.DisConnect(peerID, req.ChannelID); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	if err := cu.user.LeaveChannel(c, peerID, req.ChannelID); err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cu.resp.Success(c, common.LEAVE_SUCCESS)
}
