package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSACrypto 实现 RSA 加密和解密
type RSACrypto struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSACrypto 创建一个新的 RSACrypto 实例
func NewRSACrypto(keySize int) (*RSACrypto, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}
	return &RSACrypto{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewRSACryptoFromPrivateKeyPEM 根据 PEM 私钥创建 RSACrypto 实例。
func NewRSACryptoFromPrivateKeyPEM(privateKeyPEM string) (*RSACrypto, error) {
	privateKey, err := ParseRSAPrivateKeyPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	return &RSACrypto{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewRSACryptoFromPublicKeyPEM 根据 PEM 公钥创建 RSACrypto 实例。
func NewRSACryptoFromPublicKeyPEM(publicKeyPEM string) (*RSACrypto, error) {
	publicKey, err := ParseRSAPublicKeyPEM(publicKeyPEM)
	if err != nil {
		return nil, err
	}
	return &RSACrypto{
		publicKey: publicKey,
	}, nil
}

// Encrypt 使用公钥加密数据
func (r *RSACrypto) Encrypt(data string) (string, error) {
	return r.EncryptBytes([]byte(data))
}

// EncryptBytes 使用公钥加密二进制数据，并返回 base64 密文。
func (r *RSACrypto) EncryptBytes(data []byte) (string, error) {
	if r == nil || r.publicKey == nil {
		return "", errors.New("rsa public key is nil")
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, data, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// Decrypt 使用私钥解密数据
func (r *RSACrypto) Decrypt(encryptedData string) (string, error) {
	decryptedBytes, err := r.DecryptBytes(encryptedData)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

// DecryptBytes 使用私钥解密 base64 密文，并返回原始二进制数据。
func (r *RSACrypto) DecryptBytes(encryptedData string) ([]byte, error) {
	if r == nil || r.privateKey == nil {
		return nil, errors.New("rsa private key is nil")
	}
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, decodedData, nil)
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

// ExportPrivateKey 导出私钥为 PEM 格式
func (r *RSACrypto) ExportPrivateKey() (string, error) {
	if r == nil || r.privateKey == nil {
		return "", errors.New("rsa private key is nil")
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(r.privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM), nil
}

// ExportPrivateKeyPKCS8 导出 PKCS#8 私钥为 PEM 格式。
func (r *RSACrypto) ExportPrivateKeyPKCS8() (string, error) {
	if r == nil || r.privateKey == nil {
		return "", errors.New("rsa private key is nil")
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(r.privateKey)
	if err != nil {
		return "", err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM), nil
}

// ExportPublicKey 导出公钥为 PEM 格式
func (r *RSACrypto) ExportPublicKey() (string, error) {
	if r == nil || r.publicKey == nil {
		return "", errors.New("rsa public key is nil")
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(r.publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}

// ExportPublicKeyPKIX 导出 PKIX 公钥为 PEM 格式。
func (r *RSACrypto) ExportPublicKeyPKIX() (string, error) {
	if r == nil || r.publicKey == nil {
		return "", errors.New("rsa public key is nil")
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(r.publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}

// ParseRSAPrivateKeyPEM 从 PEM 字符串解析 RSA 私钥，兼容 PKCS#1 与 PKCS#8。
func ParseRSAPrivateKeyPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("invalid rsa private key pem")
	}
	privateKey, pkcs1Err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if pkcs1Err == nil {
		return privateKey, nil
	}
	key, pkcs8Err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if pkcs8Err != nil {
		return nil, pkcs8Err
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("pem private key is not rsa")
	}
	return privateKey, nil
}

// ParseRSAPublicKeyPEM 从 PEM 字符串解析 RSA 公钥，兼容 PKIX 与 PKCS#1。
func ParseRSAPublicKeyPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("invalid rsa public key pem")
	}
	key, pkixErr := x509.ParsePKIXPublicKey(block.Bytes)
	if pkixErr == nil {
		publicKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("pem public key is not rsa")
		}
		return publicKey, nil
	}
	publicKey, pkcs1Err := x509.ParsePKCS1PublicKey(block.Bytes)
	if pkcs1Err != nil {
		return nil, pkcs1Err
	}
	return publicKey, nil
}
