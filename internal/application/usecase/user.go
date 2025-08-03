package usecase

import (
	"context"

	"github.com/wang900115/DESA/lib/domain"
	"github.com/wang900115/DESA/lib/implement"
)

type UserUsecase struct {
	reader    implement.UserQueryService
	writer    implement.UserCommandService
	auth      implement.TokenAuthService
	secretKey string
}

func NewUserUsecase(reader *implement.UserQueryService, writer *implement.UserCommandService, auth *implement.TokenAuthService, secretKey string) *UserUsecase {
	return &UserUsecase{reader: *reader, writer: *writer, auth: *auth, secretKey: secretKey}
}

func (uu *UserUsecase) QueryUser(c context.Context) ([]domain.User, error) {
	return uu.reader.QueryUser(c)
}

func (uu *UserUsecase) DeleteUser(c context.Context, userID string) error {
	return uu.writer.DeleteUser(c, userID)
}

func (uu *UserUsecase) UpdateUser(c context.Context, toUpdate domain.User) (domain.User, error) {
	return uu.writer.UpdateUser(c, toUpdate)
}

func (uu *UserUsecase) CreateUser(c context.Context, toCreate domain.User) (domain.User, error) {
	return uu.writer.CreateUser(c, toCreate)
}

func (uu *UserUsecase) Register(c context.Context, userID string, channelID string) error {
	return uu.writer.RegisterChannel(c, userID, channelID)
}

func (uu *UserUsecase) Login(c context.Context, username, password string) (string, string, domain.User, error) {
	user, err := uu.reader.CheckPassword(c, username, password)
	if err != nil {
		return "", "", domain.User{}, err
	}

	accessToken, refreshToken, err := uu.auth.Generate(c, user.PeerID, user.Role, uu.secretKey)
	if err != nil {
		return "", "", domain.User{}, err
	}

	return accessToken, refreshToken, user, nil
}

func (uu *UserUsecase) Logout(c context.Context, userID string) error {
	return uu.auth.Delete(c, userID)
}

func (uu *UserUsecase) JoinChannel(c context.Context, userID string, channelID string) error {
	return uu.writer.JoinChannel(c, userID, channelID)
}

func (uu *UserUsecase) LeaveChannel(c context.Context, userID string, channelID string) error {
	return uu.writer.LeaveChannel(c, userID, channelID)
}
