package crypto

import (
	"crypto/cipher"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

// Sm4Encrypt 使用 SM4-CBC 加密明文，返回带 PKCS#5 填充的密文。
func Sm4Encrypt(plainText, key, iv []byte) ([]byte, error) {
	if plainText == nil {
		return nil, fmt.Errorf("plain text is nil")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv length: %d, must be %d bytes", len(iv), block.BlockSize())
	}
	plainText = PKCS5Padding(plainText, block.BlockSize())
	cipherText := make([]byte, len(plainText))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

// Sm4Decrypt 使用 SM4-CBC 解密并校验 PKCS#5 填充。
func Sm4Decrypt(cipherText, key, iv []byte) ([]byte, error) {
	if cipherText == nil {
		return nil, fmt.Errorf("cipher text is nil")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv length: %d, must be %d bytes", len(iv), block.BlockSize())
	}
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("invalid cipher text length")
	}
	plainText := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainText, cipherText)
	return unpadSM4(plainText, block.BlockSize())
}

// Sm4GCMEncrypt 使用 SM4-GCM 加密明文，返回 ciphertext||tag。
func Sm4GCMEncrypt(plainText, key, nonce []byte) ([]byte, error) {
	if plainText == nil {
		return nil, fmt.Errorf("plain text is nil")
	}
	if len(nonce) == 0 {
		return nil, fmt.Errorf("nonce is nil")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("invalid nonce length: %d, must be %d bytes", len(nonce), aead.NonceSize())
	}
	return aead.Seal(nil, nonce, plainText, nil), nil
}

// Sm4GCMDecrypt 使用 SM4-GCM 解密并校验密文完整性。
func Sm4GCMDecrypt(cipherText, key, nonce []byte) ([]byte, error) {
	if cipherText == nil {
		return nil, fmt.Errorf("cipher text is nil")
	}
	if len(nonce) == 0 {
		return nil, fmt.Errorf("nonce is nil")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("invalid nonce length: %d, must be %d bytes", len(nonce), aead.NonceSize())
	}
	return aead.Open(nil, nonce, cipherText, nil)
}

// unpadSM4 校验并去除 SM4-CBC 使用的 PKCS#5 填充。
func unpadSM4(value []byte, blockSize int) ([]byte, error) {
	if len(value) == 0 || len(value)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded text length")
	}
	padding := int(value[len(value)-1])
	if padding == 0 || padding > blockSize || padding > len(value) {
		return nil, fmt.Errorf("invalid padding")
	}
	for _, item := range value[len(value)-padding:] {
		if int(item) != padding {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return value[:len(value)-padding], nil
}
