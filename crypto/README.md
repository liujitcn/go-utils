# crypto

`crypto` 是独立 Go 子模块，提供密码哈希、消息认证、对称加密、非对称密钥和流式文件加解密能力。

## 安装

```bash
go get github.com/liujitcn/go-utils/crypto@latest
```

## 密码哈希与验证

`Crypto` 接口统一提供 `Encrypt` 和 `Verify`：

| 构造函数 | 算法 | 结果格式或默认参数 |
| --- | --- | --- |
| `NewBCryptCrypto` | bcrypt | 使用 `bcrypt.DefaultCost` |
| `NewArgon2Crypto` | Argon2id | 64 MiB、3 次迭代、2 路并行、16 字节盐、32 字节结果 |
| `NewPBKDF2Crypto` | PBKDF2-HMAC-SHA256 | 310000 次迭代、32 字节结果 |
| `NewPBKDF2WithSHA512` | PBKDF2-HMAC-SHA512 | 600000 次迭代、64 字节结果 |
| `NewSHA256Crypto` | 加盐 SHA-256 | 16 字节随机盐 |
| `NewSHA512Crypto` | 加盐 SHA-512 | 16 字节随机盐 |

也可以通过 `NewCrypto("bcrypt")`、`NewCrypto("argon2")`、`NewCrypto("pbkdf2")`、`NewCrypto("sha256")`、`NewCrypto("sha512")` 按名称创建。

```go
crypto, err := crypto.NewCrypto("argon2")
if err != nil {
	return err
}
encoded, err := crypto.Encrypt("correct horse battery staple")
if err != nil {
	return err
}
if err := crypto.Verify("correct horse battery staple", encoded); err != nil {
	return err
}
```

`HMACCrypto` 用固定密钥计算 HMAC-SHA256 十六进制摘要，适合校验消息完整性，不用于密码存储：

```go
signer := crypto.NewHMACCrypto("server-secret")
signature, err := signer.Encrypt("payload")
if err != nil {
	return err
}
err = signer.Verify("payload", signature)
```

## AES

- `GenerateAESKey(16|24|32)` 生成合法 AES 密钥。
- `AesGCMEncrypt`、`AesGCMDecrypt` 使用 AES-GCM；默认 nonce 长度为 `AESGCMNonceSize`（12 字节），返回密文后附认证标签。
- `AesGCMEncryptWithAAD`、`AesGCMDecryptWithAAD` 额外认证可见元数据，AAD 不会写入密文。
- `AesEncrypt`、`AesDecrypt` 使用 AES-CBC 和调用方提供的 IV，保留用于兼容已有数据。
- `DefaultAESKey` 仅为兼容旧调用保留，不应作为新数据的固定密钥。

```go
import (
	"crypto/rand"

	cryptoutil "github.com/liujitcn/go-utils/crypto"
)

func encryptAES(plaintext []byte) ([]byte, []byte, []byte, error) {
	key, err := cryptoutil.GenerateAESKey(32)
	if err != nil {
		return nil, nil, nil, err
	}
	nonce := make([]byte, cryptoutil.AESGCMNonceSize)
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, nil, err
	}
	ciphertext, err := cryptoutil.AesGCMEncrypt(plaintext, key, nonce)
	return ciphertext, key, nonce, err
}
```

每条消息都必须使用不可复用的随机 nonce。CBC 接口没有认证能力，新数据优先使用 GCM。

## SM4

`Sm4Encrypt`/`Sm4Decrypt` 使用 SM4-CBC 和 PKCS#5 填充；`Sm4GCMEncrypt`/`Sm4GCMDecrypt` 使用 SM4-GCM 并校验认证标签。密钥、IV 和 nonce 的长度必须符合底层算法要求。

## RSA

`RSACrypto` 使用 RSA-OAEP-SHA256：

```go
import cryptoutil "github.com/liujitcn/go-utils/crypto"

func roundTripRSA() (string, error) {
	rsaCrypto, err := cryptoutil.NewRSACrypto(2048)
	if err != nil {
		return "", err
	}
	encrypted, err := rsaCrypto.Encrypt("secret")
	if err != nil {
		return "", err
	}
	return rsaCrypto.Decrypt(encrypted)
}
```

支持从 PEM 私钥或公钥创建实例：

- `NewRSACryptoFromPrivateKeyPEM` 兼容 PKCS#1、PKCS#8 私钥。
- `NewRSACryptoFromPublicKeyPEM` 兼容 PKIX、PKCS#1 公钥。
- `ExportPrivateKey`、`ExportPrivateKeyPKCS8` 导出私钥。
- `ExportPublicKey`、`ExportPublicKeyPKIX` 导出公钥。
- `EncryptBytes`、`DecryptBytes` 用于二进制数据，密文均使用 Base64 编码。

## ECDSA 与 ECDH

- `NewECDSACrypto` 使用 P-256 对消息 SHA-256 后签名，`Encrypt` 返回 `ecdsa$r$s` 格式，`Verify` 校验签名。
- `NewECDHCrypto` 使用 P-256 密钥交换。`Encrypt` 返回 `ecdh$<Base64 公钥>`，`DeriveSharedSecret` 根据对方 Base64 公钥派生共享密钥。

ECDH 的典型流程是双方交换 `Encrypt` 返回的公钥，双方分别调用 `DeriveSharedSecret`，然后使用共享密钥建立后续加密通道。`ECDHCrypto.Verify` 的第一个参数要求传入共享密钥的 Base64 字符串。

## OpenSSL 文件格式

`NewOpenSSLFileCrypto` 提供与 OpenSSL `enc` 兼容的流式格式：`Salted__` 文件头、8 字节盐值、PBKDF2-HMAC-SHA256（100000 次迭代）派生 AES-256-CBC 密钥，并使用 PKCS#7 填充。

```go
import (
	"io"

	cryptoutil "github.com/liujitcn/go-utils/crypto"
)

func encryptFile(password string, source io.Reader, target io.Writer) error {
	fileCrypto := cryptoutil.NewOpenSSLFileCrypto()
	return fileCrypto.EncryptReader(password, source, target)
}
```

`EncryptReader` 和 `DecryptReader` 不负责关闭调用方传入的 `io.Reader` 或 `io.Writer`。

## 注意事项

- 密钥、密码和随机 nonce 由调用方安全保存，库不会代管密钥。
- 哈希结果包含随机盐，可直接持久化后交给同一个算法实例验证。
- AES-CBC、SM4-CBC 和 OpenSSL 兼容格式没有现代 AEAD 的完整性语义，仅用于兼容既有数据；新协议优先选择 GCM 或在上层增加认证。
