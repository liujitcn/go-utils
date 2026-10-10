package crypto

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

// ECDSACrypto 实现基于 ECDSA 的加密和验证
type ECDSACrypto struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewECDSACrypto 创建一个新的 ECDSACrypto 实例
func NewECDSACrypto() (*ECDSACrypto, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECDSACrypto{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// Encrypt 使用 ECDSA 对消息进行签名
func (e *ECDSACrypto) Encrypt(plainPassword string) (string, error) {
	if plainPassword == "" {
		return "", errors.New("密码不能为空")
	}

	hash := sha256.Sum256([]byte(plainPassword))
	signature, err := ecdsa.SignASN1(rand.Reader, e.privateKey, hash[:])
	if err != nil {
		return "", err
	}

	var signatureValues struct {
		R *big.Int
		S *big.Int
	}
	var remaining []byte
	remaining, err = asn1.Unmarshal(signature, &signatureValues)
	if err != nil {
		return "", err
	}
	if len(remaining) > 0 {
		return "", errors.New("签名格式无效")
	}
	return "ecdsa$" + signatureValues.R.String() + "$" + signatureValues.S.String(), nil
}

// Verify 验证消息的签名是否有效
func (e *ECDSACrypto) Verify(plainPassword, encrypted string) error {
	if plainPassword == "" || encrypted == "" {
		return errors.New("密码或加密字符串不能为空")
	}

	parts := strings.SplitN(encrypted, "$", 3)
	if len(parts) != 3 || parts[0] != "ecdsa" {
		return errors.New("加密字符串格式无效")
	}

	r := new(big.Int)
	s := new(big.Int)
	if _, ok := r.SetString(parts[1], 10); !ok {
		return errors.New("加密字符串格式无效")
	}
	if _, ok := s.SetString(parts[2], 10); !ok {
		return errors.New("加密字符串格式无效")
	}
	signature, err := asn1.Marshal(struct {
		R *big.Int
		S *big.Int
	}{r, s})
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(plainPassword))
	if !ecdsa.VerifyASN1(e.publicKey, hash[:], signature) {
		return errors.New("签名验证失败")
	}
	return nil
}

// ECDHCrypto 实现基于 ECDH 的密钥交换
type ECDHCrypto struct {
	privateKey *ecdh.PrivateKey
	publicKey  *ecdh.PublicKey
}

// NewECDHCrypto 创建一个新的 ECDHCrypto 实例
func NewECDHCrypto() (*ECDHCrypto, error) {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECDHCrypto{
		privateKey: privateKey,
		publicKey:  privateKey.PublicKey(),
	}, nil
}

// Encrypt 返回公钥作为加密结果
func (e *ECDHCrypto) Encrypt(plainPassword string) (string, error) {
	if plainPassword == "" {
		return "", errors.New("密码不能为空")
	}

	publicKeyBytes := e.publicKey.Bytes()
	return "ecdh$" + base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

// Verify 验证共享密钥是否一致
func (e *ECDHCrypto) Verify(plainPassword, encrypted string) error {
	if plainPassword == "" || encrypted == "" {
		return errors.New("密码或加密字符串不能为空")
	}

	parts := strings.SplitN(encrypted, "$", 2)
	if len(parts) != 2 || parts[0] != "ecdh" {
		return errors.New("加密字符串格式无效")
	}

	remotePublicKey, err := parseECDHPublicKey(parts[1])
	if err != nil {
		return err
	}

	var actualSecret []byte
	actualSecret, err = e.privateKey.ECDH(remotePublicKey)
	if err != nil {
		return errors.New("共享密钥计算失败")
	}
	actualSecret = bytes.TrimLeft(actualSecret, "\x00")
	var expectedSecret []byte
	expectedSecret, err = base64.StdEncoding.DecodeString(plainPassword)
	if err != nil {
		return errors.New("共享密钥格式无效")
	}
	if len(expectedSecret) != len(actualSecret) {
		return errors.New("共享密钥不匹配")
	}
	if subtle.ConstantTimeCompare(expectedSecret, actualSecret) != 1 {
		return errors.New("共享密钥不匹配")
	}
	return nil
}

// DeriveSharedSecret 根据对方公钥派生共享密钥。
func (e *ECDHCrypto) DeriveSharedSecret(publicKey string) ([]byte, error) {
	remotePublicKey, err := parseECDHPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	var sharedSecret []byte
	sharedSecret, err = e.privateKey.ECDH(remotePublicKey)
	if err != nil {
		return nil, err
	}
	return bytes.TrimLeft(sharedSecret, "\x00"), nil
}

// parseECDHPublicKey 解析 Base64 编码的压缩或未压缩 P-256 公钥。
func parseECDHPublicKey(encoded string) (*ecdh.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	if len(publicKeyBytes) == 33 {
		x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKeyBytes)
		if x == nil || y == nil {
			return nil, errors.New("无效的公钥")
		}
		uncompressedPublicKey := make([]byte, 65)
		uncompressedPublicKey[0] = 4
		x.FillBytes(uncompressedPublicKey[1:33])
		y.FillBytes(uncompressedPublicKey[33:])
		publicKeyBytes = uncompressedPublicKey
	}

	var publicKey *ecdh.PublicKey
	publicKey, err = ecdh.P256().NewPublicKey(publicKeyBytes)
	if err != nil {
		return nil, errors.New("无效的公钥")
	}
	return publicKey, nil
}
