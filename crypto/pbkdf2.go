package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"
)

const (
	// opensslFileMagic 是 OpenSSL enc 文件的固定文件头。
	opensslFileMagic = "Salted__"
	// opensslFileSaltLength 是 OpenSSL 兼容文件使用的盐值长度。
	opensslFileSaltLength = 8
	// opensslFileIterations 是 OpenSSL 兼容文件使用的 PBKDF2 迭代次数。
	opensslFileIterations = 100000
	// opensslFileDerivedKeyLength 是 AES-256 密钥和 CBC 初始向量的总长度。
	opensslFileDerivedKeyLength = 32 + aes.BlockSize
)

// PBKDF2Crypto 实现 PBKDF2-HMAC 密码哈希算法
type PBKDF2Crypto struct {
	// 可配置参数，默认使用推荐值
	Iterations int
	KeyLength  int
	Hash       func() hash.Hash
	HashName   string
}

// NewPBKDF2Crypto 创建带默认参数的 PBKDF2 加密器 (SHA256)
func NewPBKDF2Crypto() *PBKDF2Crypto {
	return &PBKDF2Crypto{
		Iterations: 310000, // NIST 推荐最小值
		KeyLength:  32,     // 256-bit
		Hash:       sha256.New,
		HashName:   "sha256",
	}
}

// NewPBKDF2WithSHA512 创建使用 SHA512 的 PBKDF2 加密器
func NewPBKDF2WithSHA512() *PBKDF2Crypto {
	return &PBKDF2Crypto{
		Iterations: 600000, // SHA512 需要更多迭代
		KeyLength:  64,     // 512-bit
		Hash:       sha512.New,
		HashName:   "sha512",
	}
}

// OpenSSLFileCrypto 实现与 OpenSSL enc 兼容的 AES-256-CBC 文件加解密。
type OpenSSLFileCrypto struct {
	iterations int
	saltLength int
	hash       func() hash.Hash
}

// NewOpenSSLFileCrypto 创建固定参数的 OpenSSL 兼容文件加解密器。
func NewOpenSSLFileCrypto() *OpenSSLFileCrypto {
	return &OpenSSLFileCrypto{
		iterations: opensslFileIterations,
		saltLength: opensslFileSaltLength,
		hash:       sha256.New,
	}
}

// Encrypt 实现密码加密
func (p *PBKDF2Crypto) Encrypt(password string) (string, error) {
	// 生成随机盐值 (16 bytes 推荐最小值)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 生成密钥
	key := pbkdf2Key([]byte(password), salt, p.Iterations, p.KeyLength, p.Hash)

	// 格式: pbkdf2:<hash>:<iterations>:<base64-salt>:<base64-key>
	return fmt.Sprintf(
		"pbkdf2:%s:%d:%s:%s",
		p.HashName,
		p.Iterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

// Verify 验证密码
func (p *PBKDF2Crypto) Verify(password, encrypted string) error {
	// 解析哈希字符串
	parts := strings.Split(encrypted, ":")
	if len(parts) != 5 || parts[0] != "pbkdf2" {
		return errors.New("无效的 PBKDF2 哈希格式")
	}

	// 解析参数
	hashName := parts[1]
	iterations, err := strconv.Atoi(parts[2])
	if err != nil {
		return errors.New("无效的迭代次数")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return err
	}

	expectedKey, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	// 根据哈希名称选择哈希函数
	hashFunc, ok := getHashFunction(hashName)
	if !ok {
		return fmt.Errorf("不支持的哈希算法: %s", hashName)
	}

	// 生成新密钥
	keyLength := len(expectedKey)
	newKey := pbkdf2Key([]byte(password), salt, iterations, keyLength, hashFunc)

	// 安全比较
	if !hmac.Equal(newKey, expectedKey) {
		return errors.New("密码不匹配")
	}
	return nil
}

// EncryptReader 使用 OpenSSL enc 兼容格式加密输入流。
func (c *OpenSSLFileCrypto) EncryptReader(password string, source io.Reader, target io.Writer) error {
	if c == nil {
		return errors.New("OpenSSL 文件加密器为空")
	}
	if source == nil {
		return errors.New("待加密输入流为空")
	}
	if target == nil {
		return errors.New("加密输出流为空")
	}
	salt := make([]byte, c.saltLength)
	var err error
	if _, err = rand.Read(salt); err != nil {
		return fmt.Errorf("生成加密盐值失败: %w", err)
	}
	if err = writeAll(target, []byte(opensslFileMagic)); err != nil {
		return fmt.Errorf("写入加密文件头失败: %w", err)
	}
	if err = writeAll(target, salt); err != nil {
		return fmt.Errorf("写入加密盐值失败: %w", err)
	}
	derivedKey := pbkdf2Key([]byte(password), salt, c.iterations, opensslFileDerivedKeyLength, c.hash)
	block, err := aes.NewCipher(derivedKey[:32])
	if err != nil {
		return fmt.Errorf("创建 AES 加密器失败: %w", err)
	}
	if err = encryptCBC(source, target, block, derivedKey[32:]); err != nil {
		return fmt.Errorf("加密文件内容失败: %w", err)
	}
	return nil
}

// DecryptReader 使用 OpenSSL enc 兼容格式解密输入流。
func (c *OpenSSLFileCrypto) DecryptReader(password string, source io.Reader, target io.Writer) error {
	if c == nil {
		return errors.New("OpenSSL 文件解密器为空")
	}
	if source == nil {
		return errors.New("待解密输入流为空")
	}
	if target == nil {
		return errors.New("解密输出流为空")
	}
	header := make([]byte, len(opensslFileMagic))
	var err error
	if _, err = io.ReadFull(source, header); err != nil {
		return fmt.Errorf("读取加密文件头失败: %w", err)
	}
	if string(header) != opensslFileMagic {
		return errors.New("加密文件头不是 OpenSSL enc 格式")
	}
	salt := make([]byte, c.saltLength)
	if _, err = io.ReadFull(source, salt); err != nil {
		return fmt.Errorf("读取加密盐值失败: %w", err)
	}
	derivedKey := pbkdf2Key([]byte(password), salt, c.iterations, opensslFileDerivedKeyLength, c.hash)
	block, err := aes.NewCipher(derivedKey[:32])
	if err != nil {
		return fmt.Errorf("创建 AES 解密器失败: %w", err)
	}
	if err = decryptCBC(source, target, block, derivedKey[32:]); err != nil {
		return fmt.Errorf("解密文件内容失败: %w", err)
	}
	return nil
}

// pbkdf2Key 实现 PBKDF2 核心算法
func pbkdf2Key(password, salt []byte, iterations, keyLength int, hashFunc func() hash.Hash) []byte {
	prf := hmac.New(hashFunc, password)
	hashLength := prf.Size()
	blockCount := (keyLength + hashLength - 1) / hashLength

	output := make([]byte, 0, blockCount*hashLength)
	for i := 1; i <= blockCount; i++ {
		// U1 = PRF(password, salt || INT(i))
		prf.Reset()
		prf.Write(salt)
		prf.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
		u := prf.Sum(nil)

		// F = U1 ⊕ U2 ⊕ ... ⊕ U_iterations
		f := make([]byte, len(u))
		copy(f, u)

		for j := 1; j < iterations; j++ {
			prf.Reset()
			prf.Write(u)
			u = prf.Sum(nil)
			for k := 0; k < len(f); k++ {
				f[k] ^= u[k]
			}
		}

		output = append(output, f...)
	}

	return output[:keyLength]
}

// encryptCBC 对输入流执行带 PKCS#7 填充的 AES-CBC 加密。
func encryptCBC(source io.Reader, target io.Writer, block cipher.Block, iv []byte) error {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return fmt.Errorf("AES-CBC 初始向量长度无效: %d", len(iv))
	}
	pending := make([]byte, 0, blockSize*2)
	buffer := make([]byte, 32*1024)
	encrypter := cipher.NewCBCEncrypter(block, iv)
	var err error
	for {
		var readCount int
		readCount, err = source.Read(buffer)
		if readCount > 0 {
			pending = append(pending, buffer[:readCount]...)
			processLength := len(pending) - blockSize
			processLength -= processLength % blockSize
			if processLength > 0 {
				if err = encryptAndWrite(pending[:processLength], target, encrypter); err != nil {
					return err
				}
				pending = pending[processLength:]
			}
		}
		if errors.Is(err, io.EOF) {
			padded := pkcs7Pad(pending, blockSize)
			if err = encryptAndWrite(padded, target, encrypter); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("读取待加密内容失败: %w", err)
		}
	}
}

// decryptCBC 对输入流执行 AES-CBC 解密并校验 PKCS#7 填充。
func decryptCBC(source io.Reader, target io.Writer, block cipher.Block, iv []byte) error {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return fmt.Errorf("AES-CBC 初始向量长度无效: %d", len(iv))
	}
	buffer := make([]byte, blockSize)
	decrypter := cipher.NewCBCDecrypter(block, iv)
	var previous []byte
	var err error
	for {
		_, err = io.ReadFull(source, buffer)
		if errors.Is(err, io.EOF) {
			if previous == nil {
				return errors.New("加密文件没有密文内容")
			}
			unpadded, unpadErr := pkcs7Unpad(previous, blockSize)
			if unpadErr != nil {
				return unpadErr
			}
			return writeAll(target, unpadded)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return errors.New("加密文件密文长度不是 AES 分组长度的整数倍")
		}
		if err != nil {
			return fmt.Errorf("读取待解密内容失败: %w", err)
		}
		plainText := make([]byte, blockSize)
		decrypter.CryptBlocks(plainText, buffer)
		if previous != nil {
			if err = writeAll(target, previous); err != nil {
				return fmt.Errorf("写入解密内容失败: %w", err)
			}
		}
		previous = plainText
	}
}

// encryptAndWrite 加密完整分组并写入目标流。
func encryptAndWrite(plainText []byte, target io.Writer, encrypter cipher.BlockMode) error {
	if len(plainText) == 0 || len(plainText)%encrypter.BlockSize() != 0 {
		return errors.New("待加密内容不是完整 AES 分组")
	}
	cipherText := make([]byte, len(plainText))
	encrypter.CryptBlocks(cipherText, plainText)
	return writeAll(target, cipherText)
}

// pkcs7Pad 为数据添加 PKCS#7 填充。
func pkcs7Pad(value []byte, blockSize int) []byte {
	padding := blockSize - len(value)%blockSize
	result := make([]byte, len(value)+padding)
	copy(result, value)
	for index := len(value); index < len(result); index++ {
		result[index] = byte(padding)
	}
	return result
}

// pkcs7Unpad 校验并移除 PKCS#7 填充。
func pkcs7Unpad(value []byte, blockSize int) ([]byte, error) {
	if len(value) == 0 || len(value)%blockSize != 0 {
		return nil, errors.New("解密内容长度无效")
	}
	padding := int(value[len(value)-1])
	if padding == 0 || padding > blockSize || padding > len(value) {
		return nil, errors.New("解密内容 PKCS#7 填充无效")
	}
	for _, item := range value[len(value)-padding:] {
		if int(item) != padding {
			return nil, errors.New("解密内容 PKCS#7 填充无效")
		}
	}
	return value[:len(value)-padding], nil
}

// writeAll 确保完整写入目标流。
func writeAll(target io.Writer, value []byte) error {
	written, err := target.Write(value)
	if err != nil {
		return err
	}
	if written != len(value) {
		return io.ErrShortWrite
	}
	return nil
}

// getHashFunction 根据名称获取哈希函数
func getHashFunction(name string) (func() hash.Hash, bool) {
	switch name {
	case "sha256":
		return sha256.New, true
	case "sha512":
		return sha512.New, true
	default:
		return nil, false
	}
}
