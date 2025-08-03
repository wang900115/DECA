package encrypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/wang900115/DESA/lib/common"
	"github.com/wang900115/DESA/lib/constant"
	"golang.org/x/crypto/argon2"
)

func HashPasswordArgon2id(password string) (string, error) {
	salt := make([]byte, constant.SALTLENGTH)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.Key([]byte(password), salt, constant.ITERATIONS, constant.MEMORY, constant.PARALLELISM, constant.KEYLENGTH)
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, constant.MEMORY, constant.ITERATIONS, constant.PARALLELISM, b64salt, b64hash)

	return encoded, nil
}

func VerifyPasswordArgon2id(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, common.HashPassword
	}
	var version, m, t, p int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	result := argon2.IDKey([]byte(password), salt, uint32(t), uint32(m), uint8(p), uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, result) == 1 {
		return true, nil
	}
	return false, nil
}
