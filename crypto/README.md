# 密码加解密

算法列表

| 算法名                                           | 特点                            | 用途                     |
|-----------------------------------------------|-------------------------------|------------------------|
| Bcrypt                                        | 基于 Blowfish 算法，内置盐值，支持调整计算成本。 | 密码存储，防止暴力破解。           |
| Argon2                                        | 内存硬性哈希算法，支持多种参数调整（内存、时间、并行度）。 | 码存储，适合高安全性需求。          |
| PBKDF2                                        | 基于 HMAC 的密钥派生函数，支持多种哈希算法。     | 密码存储，兼容性好。             |
| SHA-256/SHA-512                               | 快速单向哈希算法，无内置盐值。               | 数据完整性校验，需手动添加盐值用于密码存储。 |
| AES (Advanced Encryption Standard)            | 对称加密算法，支持 128/192/256 位密钥。    | 数据加密传输或存储。             |
| RSA                                           | 非对称加密算法，基于大整数分解难题。            | 数据加密、数字签名。             |
| ECDSA/ECDH (椭圆曲线算法)                           | 基于椭圆曲线的非对称加密，密钥更短但安全性高。       | 数据完整性和认证。              |
| HMAC (Hash-based Message Authentication Code) | 基于哈希算法的消息认证码。                 | 数据完整性和认证。              |

## RSA

- `NewRSACrypto(keySize int)`：生成 RSA 密钥对。
- `NewRSACryptoFromPrivateKeyPEM(privateKeyPEM string)`：从 PEM 私钥创建 RSA 加解密实例，兼容 PKCS#1 与 PKCS#8。
- `NewRSACryptoFromPublicKeyPEM(publicKeyPEM string)`：从 PEM 公钥创建 RSA 加密实例，兼容 PKIX 与 PKCS#1。
- `Encrypt(data string)` / `Decrypt(encryptedData string)`：使用 RSA-OAEP-SHA256 加解密字符串，密文使用 base64 编码。
- `EncryptBytes(data []byte)` / `DecryptBytes(encryptedData string)`：使用 RSA-OAEP-SHA256 加解密二进制数据，密文使用 base64 编码。
- `ExportPrivateKey()`：导出 PKCS#1 私钥 PEM。
- `ExportPrivateKeyPKCS8()`：导出 PKCS#8 私钥 PEM。
- `ExportPublicKey()`：导出 PKIX 公钥 PEM，PEM 块类型为 `RSA PUBLIC KEY`。
- `ExportPublicKeyPKIX()`：导出 PKIX 公钥 PEM，PEM 块类型为 `PUBLIC KEY`，适合浏览器 WebCrypto `spki` 导入。

## AES

- `GenerateAESKey(length int)`：生成 16、24 或 32 字节 AES 密钥。
- `AesEncrypt(plainText, key, iv []byte)` / `AesDecrypt(cryptedText, key, iv []byte)`：使用 AES-CBC 加解密，保留用于兼容已有调用。
- `AesGCMEncrypt(plainText, key, nonce []byte)` / `AesGCMDecrypt(cipherText, key, nonce []byte)`：使用 AES-GCM 加解密，推荐用于需要防篡改校验的传输加密场景。
