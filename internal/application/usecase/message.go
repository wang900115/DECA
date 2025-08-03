package usecase

import (
	"context"

	"github.com/wang900115/DESA/lib/domain"
	"github.com/wang900115/DESA/lib/implement"
)

type MessageUsecase struct {
	reader implement.MessageQueryService
	writer implement.MessageCommandService
}

func NewMessageUsecase(reader *implement.MessageQueryService, writer *implement.MessageCommandService) *MessageUsecase {
	return &MessageUsecase{reader: *reader, writer: *writer}
}

func (mu *MessageUsecase) GetChannelMessages(c context.Context, channelID string) ([]domain.Message, error) {
	return mu.reader.QueryMessage(c, channelID)
}

func (mu *MessageUsecase) GetChannelUserMessages(c context.Context, channelID string, userID string) ([]domain.Message, error) {
	return mu.reader.QueryCertainMessage(c, channelID, userID)
}

func (mu *MessageUsecase) CreateMessage(c context.Context, toCreate domain.Message) (domain.Message, error) {
	return mu.writer.CreateMessage(c, toCreate)
}

func (mu *MessageUsecase) UpdateMessage(c context.Context, toUpdate domain.Message) (domain.Message, error) {
	return mu.writer.UpdateMessage(c, toUpdate)
}

func (mu *MessageUsecase) DeleteMessage(c context.Context, messageID uint) error {
	return mu.writer.DeleteMessage(c, messageID)
}
