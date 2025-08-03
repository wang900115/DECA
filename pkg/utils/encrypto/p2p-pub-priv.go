package encrypto

import (
	"encoding/base64"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/wang900115/DESA/lib/constant"
)

func GenerateKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, constant.KEY_BIT)
	return priv, pub, err
}

func EncodePrivateKey(priv crypto.PrivKey) ([]byte, error) {
	return crypto.MarshalPrivateKey(priv)
}

func DecodePrivateKey(data []byte) (crypto.PrivKey, error) {
	return crypto.UnmarshalPrivateKey(data)
}

func EncodePublicKey(pub crypto.PubKey) ([]byte, error) {
	return crypto.MarshalPublicKey(pub)
}

func DecodePublicKey(data []byte) (crypto.PubKey, error) {
	return crypto.UnmarshalPublicKey(data)
}

func EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
