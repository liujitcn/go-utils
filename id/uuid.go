package id

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
)

// formatUUID 按是否带连字符格式化 UUID。
func formatUUID(u uuid.UUID, withHyphen bool) string {
	if withHyphen {
		return u.String()
	}

	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[:])
}

// NewGUIDv4 生成带连字符的 UUID v4。
func NewGUIDv4() string {
	return formatUUID(uuid.New(), true)
}

// NewGUIDv4NoHyphen 生成不带连字符的 UUID v4。
func NewGUIDv4NoHyphen() string {
	return formatUUID(uuid.New(), false)
}

// NewGUIDv7 生成带连字符的 UUID v7。
func NewGUIDv7() string {
	u, err := uuid.NewV7()
	if err != nil {
		// 系统时钟异常时回退到 v4，避免返回空值
		return NewGUIDv4()
	}

	return formatUUID(u, true)
}

// NewGUIDv7NoHyphen 生成不带连字符的 UUID v7。
func NewGUIDv7NoHyphen() string {
	u, err := uuid.NewV7()
	if err != nil {
		// 系统时钟异常时回退到 v4，避免返回空值
		return NewGUIDv4NoHyphen()
	}

	return formatUUID(u, false)
}

// NewShortUUID 生成短格式 UUID。
func NewShortUUID() string {
	return shortuuid.New()
}

// NewXID 生成 XID。
func NewXID() string {
	return xid.New().String()
}
