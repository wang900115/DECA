package usecase

import (
	"context"
	"encoding/base64"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/wang900115/DESA/lib/domain"
	"github.com/wang900115/DESA/lib/implement"
	"github.com/wang900115/DESA/pkg/utils/encrypto"
)

type P2PUsecase struct {
	channelP2P implement.ChannelP2PService
	userP2P    implement.UserP2PService
}

func NewP2PUsecase(channelP2P *implement.ChannelP2PService, userP2P *implement.UserP2PService) *P2PUsecase {
	return &P2PUsecase{channelP2P: *channelP2P, userP2P: *userP2P}
}

func (p *P2PUsecase) CreateChannelHost(channel domain.Channel) (string, string, string, string, error) {
	peerID, addr, pri, pub, err := p.channelP2P.CreateHost(channel)
	if err != nil {
		return "", "", "", "", err
	}
	priKey, err := encrypto.EncodePrivateKey(pri)
	if err != nil {
		return "", "", "", "", err
	}
	pubKey, err := encrypto.EncodePublicKey(pub)
	if err != nil {
		return "", "", "", "", err
	}
	return peerID, addr, encrypto.EncodeToString(priKey), encrypto.EncodeToString(pubKey), nil
}

func (p *P2PUsecase) ShutDownChannelHost(peer_id string) error {
	return p.channelP2P.ShutDownHost(peer_id)
}

func (p *P2PUsecase) CreateUserHost(user domain.User) (string, string, string, string, error) {
	peerID, addr, pri, pub, err := p.userP2P.CreateHost(user)
	if err != nil {
		return "", "", "", "", err
	}
	priKey, err := encrypto.EncodePrivateKey(pri)
	if err != nil {
		return "", "", "", "", err
	}
	pubKey, err := encrypto.EncodePublicKey(pub)
	if err != nil {
		return "", "", "", "", err
	}
	return peerID, addr, encrypto.EncodeToString(priKey), encrypto.EncodeToString(pubKey), nil
}

func (p *P2PUsecase) ShutDownUserHost(peer_id string) error {
	return p.userP2P.ShutDownHost(peer_id)
}

func (p *P2PUsecase) Register(channelPeerID string, userPeerID string, pub string) error {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return err
	}
	pubKey, err := encrypto.DecodePublicKey(pubKeyBytes)
	if err != nil {
		return err
	}
	user, err := peer.Decode(userPeerID)
	if err != nil {
		return err
	}
	return p.channelP2P.RegisterUser(channelPeerID, user, pubKey)
}

func (p *P2PUsecase) Verify(channelPeerID string, userPeerID string, msg []byte, sig []byte) (bool, error) {
	user, err := peer.Decode(userPeerID)
	if err != nil {
		return false, err
	}
	return p.channelP2P.Verify(channelPeerID, user, msg, sig)
}

func (p *P2PUsecase) Connect(c context.Context, userPeerID string, channelPeerID string, addr string) error {
	channel, err := peer.Decode(channelPeerID)
	if err != nil {
		return err
	}
	targetMaddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	return p.userP2P.ConnectToChannelPeer(c, userPeerID, channel, targetMaddr)
}

func (p *P2PUsecase) DisConnect(userPeerID string, channelPeerID string) error {
	channel, err := peer.Decode(channelPeerID)
	if err != nil {
		return err
	}
	return p.userP2P.DisconnectFromChannelPeer(userPeerID, channel)
}
