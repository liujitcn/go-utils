package crypto

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestOpenSSLFileCryptoRoundTrip 验证 OpenSSL 兼容文件加解密的流式往返。
func TestOpenSSLFileCryptoRoundTrip(t *testing.T) {
	plainText := bytes.Repeat([]byte("kratos-admin-backup\x00\n"), 4097)
	cipher := NewOpenSSLFileCrypto()
	var encrypted bytes.Buffer
	if err := cipher.EncryptReader("backup-password", bytes.NewReader(plainText), &encrypted); err != nil {
		t.Fatalf("EncryptReader() error = %v", err)
	}
	if !bytes.HasPrefix(encrypted.Bytes(), []byte(opensslFileMagic)) {
		t.Fatalf("encrypted data does not have OpenSSL header")
	}
	var decrypted bytes.Buffer
	if err := cipher.DecryptReader("backup-password", bytes.NewReader(encrypted.Bytes()), &decrypted); err != nil {
		t.Fatalf("DecryptReader() error = %v", err)
	}
	if !bytes.Equal(decrypted.Bytes(), plainText) {
		t.Fatalf("decrypted data does not match plaintext")
	}
}

// TestOpenSSLFileCryptoRejectsInvalidCiphertext 验证非法文件头、密文长度和填充会被拒绝。
func TestOpenSSLFileCryptoRejectsInvalidCiphertext(t *testing.T) {
	cipher := NewOpenSSLFileCrypto()
	tests := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: nil},
		{name: "invalid-header", data: []byte("invalid!")},
		{name: "short-salt", data: []byte(opensslFileMagic)},
		{name: "short-ciphertext", data: append([]byte(opensslFileMagic), bytes.Repeat([]byte{1}, opensslFileSaltLength+1)...)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := cipher.DecryptReader("backup-password", bytes.NewReader(test.data), &output); err == nil {
				t.Fatal("DecryptReader() unexpectedly succeeded")
			}
		})
	}
}

// TestOpenSSLFileCryptoInterop 验证 Go 与 OpenSSL enc 的双向文件兼容性。
func TestOpenSSLFileCryptoInterop(t *testing.T) {
	opensslPath, err := exec.LookPath("openssl")
	if err != nil {
		t.Skipf("openssl is unavailable: %v", err)
	}
	directory := t.TempDir()
	plainPath := filepath.Join(directory, "plain.bin")
	goCipherPath := filepath.Join(directory, "go.bin")
	goPlainPath := filepath.Join(directory, "go-plain.bin")
	opensslCipherPath := filepath.Join(directory, "openssl.bin")
	opensslPlainPath := filepath.Join(directory, "openssl-plain.bin")
	plainText := bytes.Repeat([]byte("openssl-compatible-pbkdf2\n"), 257)
	if err = os.WriteFile(plainPath, plainText, 0o600); err != nil {
		t.Fatalf("write plaintext: %v", err)
	}

	input, err := os.Open(plainPath)
	if err != nil {
		t.Fatalf("open plaintext: %v", err)
	}
	output, err := os.OpenFile(goCipherPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		_ = input.Close()
		t.Fatalf("create Go ciphertext: %v", err)
	}
	if err = NewOpenSSLFileCrypto().EncryptReader("interop-password", input, output); err != nil {
		_ = input.Close()
		_ = output.Close()
		t.Fatalf("Go encryption: %v", err)
	}
	if err = input.Close(); err != nil {
		t.Fatalf("close plaintext: %v", err)
	}
	if err = output.Close(); err != nil {
		t.Fatalf("close Go ciphertext: %v", err)
	}
	runOpenSSL(t, opensslPath, "-d", goCipherPath, opensslPlainPath)
	assertFileContent(t, opensslPlainPath, plainText)

	runOpenSSL(t, opensslPath, "", plainPath, opensslCipherPath)
	input, err = os.Open(opensslCipherPath)
	if err != nil {
		t.Fatalf("open OpenSSL ciphertext: %v", err)
	}
	output, err = os.OpenFile(goPlainPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		_ = input.Close()
		t.Fatalf("create Go plaintext: %v", err)
	}
	if err = NewOpenSSLFileCrypto().DecryptReader("interop-password", input, output); err != nil {
		_ = input.Close()
		_ = output.Close()
		t.Fatalf("Go decryption: %v", err)
	}
	if err = input.Close(); err != nil {
		t.Fatalf("close OpenSSL ciphertext: %v", err)
	}
	if err = output.Close(); err != nil {
		t.Fatalf("close Go plaintext: %v", err)
	}
	assertFileContent(t, goPlainPath, plainText)
}

// runOpenSSL 使用固定协议参数执行 OpenSSL enc 命令。
func runOpenSSL(t *testing.T, command, mode, input, output string) {
	t.Helper()
	args := []string{"enc"}
	if mode != "" {
		args = append(args, mode)
	}
	args = append(args, "-aes-256-cbc", "-pbkdf2", "-iter", "100000", "-md", "sha256", "-saltlen", "8", "-in", input, "-out", output, "-pass", "pass:interop-password")
	commandRun := exec.Command(command, args...)
	if outputValue, err := commandRun.CombinedOutput(); err != nil {
		t.Fatalf("openssl %v: %v: %s", mode, err, outputValue)
	}
}

// assertFileContent 验证文件内容与期望值一致。
func assertFileContent(t *testing.T, path string, expected []byte) {
	t.Helper()
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("file %s does not match expected content", path)
	}
}

// TestPBKDF2CryptoRoundTrip 保持现有 PBKDF2 密码哈希模式可用。
func TestPBKDF2CryptoRoundTrip(t *testing.T) {
	cipher := NewPBKDF2Crypto()
	encrypted, err := cipher.Encrypt("password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if err = cipher.Verify("password", encrypted); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if err = cipher.Verify("wrong-password", encrypted); err == nil {
		t.Fatal("Verify() unexpectedly accepted a wrong password")
	}
}
