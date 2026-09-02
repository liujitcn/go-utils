package crypto

import (
	"bytes"
	"testing"
)

func TestAesGCMWithAAD(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	nonce := bytes.Repeat([]byte{2}, AESGCMNonceSize)
	additionalData := []byte("config")
	plaintext := []byte("redis-password")

	ciphertext, err := AesGCMEncryptWithAAD(plaintext, key, nonce, additionalData)
	if err != nil {
		t.Fatal(err)
	}

	var decrypted []byte
	decrypted, err = AesGCMDecryptWithAAD(ciphertext, key, nonce, additionalData)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text = %q, want %q", decrypted, plaintext)
	}
}

func TestAesGCMWithAADRejectsDifferentAdditionalData(t *testing.T) {
	key := bytes.Repeat([]byte{3}, 32)
	nonce := bytes.Repeat([]byte{4}, AESGCMNonceSize)

	ciphertext, err := AesGCMEncryptWithAAD([]byte("value"), key, nonce, []byte("config"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err = AesGCMDecryptWithAAD(ciphertext, key, nonce, []byte("other")); err == nil {
		t.Fatal("decrypt with different additional data unexpectedly succeeded")
	}
}

func TestAesGCMWithoutAADRemainsCompatible(t *testing.T) {
	key := bytes.Repeat([]byte{5}, 32)
	nonce := bytes.Repeat([]byte{6}, AESGCMNonceSize)
	plaintext := []byte("value")

	ciphertext, err := AesGCMEncrypt(plaintext, key, nonce)
	if err != nil {
		t.Fatal(err)
	}

	var decrypted []byte
	decrypted, err = AesGCMDecrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text = %q, want %q", decrypted, plaintext)
	}
}
