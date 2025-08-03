package validator

type CreateMessageRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
	// UserID    uint `json:"userID" binding:"required"`

	MsgType string `json:"msgType" binding:"required"`
	Status  string `json:"status" binding:"required"`

	ReplyToID *uint `json:"replyToID" binding:"omitempty"`

	Content   string `json:"content"  binding:"required"`
	AttachURL string `json:"attachURL" binding:"omitempty"`

	PublicKey string `json:"publicKey" binding:"required"`
}

type UpdateMessageRequest struct {
	ID        uint   `json:"id" binding:"required"`
	ChannelID string `json:"channelID" binding:"required"`
	// UserID    uint `json:"userID" binding:"required"`

	MsgType string `json:"msgType" binding:"omitempty"`
	Status  string `json:"status" binding:"omitempty"`

	ReplyToID *uint `json:"replyToID" binding:"omitempty"`

	Content   string `json:"content" binding:"omitempty"`
	AttachURL string `json:"attachURL" binding:"omitempty"`
}

type DeleteMessageRequest struct {
	MeesageID uint
}

type GetChannelMessagesRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
}

type GetChannelUserMessagesRequest struct {
	ChannelID string `json:"channelID" binding:"required"`
	UserID    string `json:"userID" binding:"required"`
}
