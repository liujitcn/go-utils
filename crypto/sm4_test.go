package crypto

import (
	"bytes"
	"testing"
)

// TestSm4CBC 验证 SM4-CBC 加解密和填充校验。
func TestSm4CBC(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("abcdef0123456789")
	plainText := []byte("sm4 cbc test")

	cipherText, err := Sm4Encrypt(plainText, key, iv)
	if err != nil {
		t.Fatalf("Sm4Encrypt() error = %v", err)
	}
	decrypted, err := Sm4Decrypt(cipherText, key, iv)
	if err != nil {
		t.Fatalf("Sm4Decrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plainText) {
		t.Fatalf("Sm4Decrypt() = %q, want %q", decrypted, plainText)
	}
}

// TestSm4GCM 验证 SM4-GCM 加解密和完整性校验。
func TestSm4GCM(t *testing.T) {
	key := []byte("0123456789abcdef")
	nonce := []byte("0123456789ab")
	plainText := []byte("sm4 gcm test")

	cipherText, err := Sm4GCMEncrypt(plainText, key, nonce)
	if err != nil {
		t.Fatalf("Sm4GCMEncrypt() error = %v", err)
	}
	decrypted, err := Sm4GCMDecrypt(cipherText, key, nonce)
	if err != nil {
		t.Fatalf("Sm4GCMDecrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plainText) {
		t.Fatalf("Sm4GCMDecrypt() = %q, want %q", decrypted, plainText)
	}

	cipherText[len(cipherText)-1] ^= 1
	if _, err = Sm4GCMDecrypt(cipherText, key, nonce); err == nil {
		t.Fatal("Sm4GCMDecrypt() expected tampered ciphertext error")
	}
}
