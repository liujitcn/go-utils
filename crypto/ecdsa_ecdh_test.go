package crypto

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestECDSACrypto_EncryptAndVerify(t *testing.T) {
	crypto, err := NewECDSACrypto()
	if err != nil {
		t.Fatalf("创建 ECDSACrypto 实例失败: %v", err)
	}

	message := "test message"

	// 签名消息
	encrypted, err := crypto.Encrypt(message)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 验证签名
	isValid, err := crypto.Verify(message, encrypted)
	if err != nil {
		t.Fatalf("验证失败: %v", err)
	}

	if !isValid {
		t.Fatal("签名验证未通过")
	}
}

func TestECDHCrypto_EncryptAndVerify(t *testing.T) {
	crypto1, err := NewECDHCrypto()
	if err != nil {
		t.Fatalf("创建 ECDHCrypto 实例1失败: %v", err)
	}

	crypto2, err := NewECDHCrypto()
	if err != nil {
		t.Fatalf("创建 ECDHCrypto 实例2失败: %v", err)
	}

	// 获取 crypto1 的公钥字符串: ecdh$base64(pubKey)
	encrypted, err := crypto1.Encrypt("peer-public-key")
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	parts := strings.SplitN(encrypted, "$", 2)
	if len(parts) != 2 {
		t.Fatalf("公钥格式无效: %s", encrypted)
	}
	pubKey1 := parts[1]

	// crypto2 先根据 crypto1 公钥推导共享密钥，并将其作为 Verify 的 plainPassword 参数
	secret2, err := crypto2.DeriveSharedSecret(pubKey1)
	if err != nil {
		t.Fatalf("推导共享密钥失败: %v", err)
	}

	// 验证共享密钥
	isValid, err := crypto2.Verify(base64.StdEncoding.EncodeToString(secret2), encrypted)
	if err != nil {
		t.Fatalf("验证失败: %v", err)
	}

	if !isValid {
		t.Fatal("共享密钥验证未通过")
	}
}
