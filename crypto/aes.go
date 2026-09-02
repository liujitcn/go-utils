package crypto

import (
	"bytes"
	"fmt"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// DefaultAESKey 默认AES密钥(16字节)
var DefaultAESKey = []byte("f51d66a73d8a0927")

// AESGCMNonceSize 是标准 AES-GCM 的 nonce 长度。
const AESGCMNonceSize = 12

// GenerateAESKey 生成AES密钥
func GenerateAESKey(length int) ([]byte, error) {
	if length != 16 && length != 24 && length != 32 {
		return nil, fmt.Errorf("invalid key length: %d, must be 16, 24, or 32 bytes", length)
	}
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// PKCS5Padding 填充明文
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

// PKCS5UnPadding 去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt AES加密
func AesEncrypt(plainText, key, iv []byte) ([]byte, error) {
	if plainText == nil {
		return nil, fmt.Errorf("plain text is nil")
	}
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()

	if iv == nil {
		// 初始向量的长度必须等于块block的长度16字节
		iv = key[:blockSize]
	}

	plainText = PKCS5Padding(plainText, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryptedText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cryptedText, plainText)
	return cryptedText, nil
}

// AesDecrypt AES解密
func AesDecrypt(cryptedText, key, iv []byte) ([]byte, error) {
	if cryptedText == nil {
		return nil, fmt.Errorf("crypted text is nil")
	}
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()

	if iv == nil {
		// 初始向量的长度必须等于块block的长度16字节
		iv = key[:blockSize]
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	plainText := make([]byte, len(cryptedText))
	blockMode.CryptBlocks(plainText, cryptedText)
	plainText = PKCS5UnPadding(plainText)
	return plainText, nil
}

// AesGCMEncrypt 使用 AES-GCM 加密明文。
func AesGCMEncrypt(plainText, key, nonce []byte) ([]byte, error) {
	return AesGCMEncryptWithAAD(plainText, key, nonce, nil)
}

// AesGCMEncryptWithAAD 使用 AES-GCM 加密明文，并认证附加数据。
func AesGCMEncryptWithAAD(plainText, key, nonce, additionalData []byte) ([]byte, error) {
	if plainText == nil {
		return nil, fmt.Errorf("plain text is nil")
	}
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}
	if nonce == nil {
		return nil, fmt.Errorf("nonce is nil")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var aesGCM cipher.AEAD
	aesGCM, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aesGCM.NonceSize() {
		return nil, fmt.Errorf("invalid nonce length: %d, must be %d bytes", len(nonce), aesGCM.NonceSize())
	}
	cipherText := aesGCM.Seal(nil, nonce, plainText, additionalData)
	return cipherText, nil
}

// AesGCMDecrypt 使用 AES-GCM 解密密文。
func AesGCMDecrypt(cipherText, key, nonce []byte) ([]byte, error) {
	return AesGCMDecryptWithAAD(cipherText, key, nonce, nil)
}

// AesGCMDecryptWithAAD 使用 AES-GCM 解密密文，并校验附加数据。
func AesGCMDecryptWithAAD(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	if cipherText == nil {
		return nil, fmt.Errorf("cipher text is nil")
	}
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}
	if nonce == nil {
		return nil, fmt.Errorf("nonce is nil")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var aesGCM cipher.AEAD
	aesGCM, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aesGCM.NonceSize() {
		return nil, fmt.Errorf("invalid nonce length: %d, must be %d bytes", len(nonce), aesGCM.NonceSize())
	}
	plainText, err := aesGCM.Open(nil, nonce, cipherText, additionalData)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
