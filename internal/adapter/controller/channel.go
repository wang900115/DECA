package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/DESA/internal/adapter/validator"
	"github.com/wang900115/DESA/internal/application/usecase"
	"github.com/wang900115/DESA/lib/common"
	iresponse "github.com/wang900115/DESA/lib/common/response"
	"github.com/wang900115/DESA/lib/domain"
)

type ChannelController struct {
	p2p     usecase.P2PUsecase
	channel usecase.ChannelUsecase
	resp    iresponse.IResponse
}

func NewChannelController(channel *usecase.ChannelUsecase, p2p *usecase.P2PUsecase, resp iresponse.IResponse) *ChannelController {
	return &ChannelController{channel: *channel, p2p: *p2p, resp: resp}
}

// 獲取所有頻道資訊
func (cc *ChannelController) GetAllChannels(c *gin.Context) {
	channels, err := cc.channel.GetAllChannels(c)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cc.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"channels": channels,
	})
}

// 建立頻道節點
func (cc *ChannelController) Create(c *gin.Context) {
	var req validator.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	channel := domain.Channel{
		ChannelName: req.ChannelName,
		ChannelType: req.ChannelType,
	}

	peerID, multiAddr, priv, pub, err := cc.p2p.CreateChannelHost(channel)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	channel.PeerID = peerID
	channel.MultiAddr = multiAddr
	channel.PrivateKey = &priv

	created, err := cc.channel.CreateChannel(c, channel)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	cc.resp.SuccessWithData(c, common.CREATE_SUCCESS, map[string]interface{}{
		"channel": created,
		"key":     pub,
	})
}

// 刪除頻道節點
func (cc *ChannelController) Delete(c *gin.Context) {
	var req validator.DeleteChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	if err := cc.p2p.ShutDownChannelHost(req.PeerID); err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	if err := cc.channel.DeleteChannel(c, req.PeerID); err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cc.resp.Success(c, common.DELETE_SUCCESS)
}
