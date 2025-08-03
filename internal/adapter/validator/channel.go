package validator

type CreateChannelRequest struct {
	ChannelName string `json:"channelName" binding:"required, min:1 max:30"`
	ChannelType string `json:"channelType" binding:"required"`
}

type DeleteChannelRequest struct {
	PeerID string `json:"peerID" binding:"required"`
}

type GetAllChannelsRequest struct {
}

type GetUserChannelsRequest struct {
}

type GetChannelUsersRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
}

type RegisterToChannelRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
	PubKey    string `json:"pubKey" binding:"required"`
}

type JoinChannelRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Message   []byte `json:"message" binding:"required"`
	Signature []byte `json:"signature" binding:"required"`
}

type LeaveChannelRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
}
