package usecase

import (
	"context"

	"github.com/wang900115/DESA/lib/domain"
	"github.com/wang900115/DESA/lib/implement"
)

type ChannelUsecase struct {
	reader implement.ChannelQueryService
	writer implement.ChannelCommandService
}

func NewChannelUsecase(reader *implement.ChannelQueryService, writer *implement.ChannelCommandService) *ChannelUsecase {
	return &ChannelUsecase{reader: *reader, writer: *writer}
}

func (cu *ChannelUsecase) GetAllChannels(c context.Context) ([]domain.Channel, error) {
	return cu.reader.QueryChannel(c)
}

func (cu *ChannelUsecase) GetUserChannels(c context.Context, userID string) ([]domain.Channel, error) {
	return cu.reader.QueryCertainChannel(c, userID)
}

func (cu *ChannelUsecase) GetChannelUsers(c context.Context, channelID string) ([]domain.User, error) {
	return cu.reader.QueryUser(c, channelID)
}

func (cu *ChannelUsecase) CreateChannel(c context.Context, toCreate domain.Channel) (domain.Channel, error) {
	return cu.writer.CreateChannel(c, toCreate)
}

func (cu *ChannelUsecase) DeleteChannel(c context.Context, peerID string) error {
	return cu.writer.DeleteChannel(c, peerID)
}

func (cu *ChannelUsecase) UpdateChannel(c context.Context, toUpdate domain.Channel) (domain.Channel, error) {
	return cu.writer.UpdateChannel(c, toUpdate)
}
