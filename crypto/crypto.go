package crypto

import (
	"errors"
	"strings"
	"sync"
)

// Crypto 密码加解密接口
type Crypto interface {
	// Encrypt 加密密码，返回加密后的字符串（包含算法标识和盐值）
	Encrypt(plainPassword string) (encrypted string, err error)

	// Verify 验证密码是否匹配
	Verify(plainPassword, encrypted string) (bool, error)
}

func NewCrypto(algorithm string) (Crypto, error) {
	algorithm = strings.ToLower(strings.TrimSpace(algorithm))
	switch algorithm {
	case "bcrypt":
		return NewBCryptCrypto(), nil
	case "pbkdf2":
		return NewPBKDF2Crypto(), nil
	case "argon2":
		return NewArgon2Crypto(), nil
	case "sha256":
		return NewSHA256Crypto(), nil
	case "sha512":
		return NewSHA512Crypto(), nil
	case "ecdsa":
		return NewECDSACrypto()
	case "ecdh":
		return NewECDHCrypto()
	default:
		return nil, errors.New("不支持的加密算法")
	}
}

const defaultAlgorithm = "bcrypt"

var (
	defaultCrypto     Crypto
	defaultCryptoErr  error
	defaultCryptoOnce sync.Once
)

// newDefaultCrypto 返回默认加密器的单例实例。
func newDefaultCrypto() {
	defaultCryptoOnce.Do(func() {
		defaultCrypto, defaultCryptoErr = NewCrypto(defaultAlgorithm)
	})
}

// Encrypt 加密数据
func Encrypt(plainPassword string) (string, error) {
	newDefaultCrypto()
	if defaultCryptoErr != nil {
		return "", defaultCryptoErr
	}
	return defaultCrypto.Encrypt(plainPassword)
}

// Verify 验证加密数据
func Verify(plainPassword, encrypted string) (bool, error) {
	newDefaultCrypto()
	if defaultCryptoErr != nil {
		return false, defaultCryptoErr
	}
	return defaultCrypto.Verify(plainPassword, encrypted)
}
