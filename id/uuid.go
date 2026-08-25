package id

import (
	"encoding/base64"
	"encoding/hex"
	"uuid"

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
	return formatUUID(uuid.NewV4(), true)
}

// NewGUIDv4NoHyphen 生成不带连字符的 UUID v4。
func NewGUIDv4NoHyphen() string {
	return formatUUID(uuid.NewV4(), false)
}

// NewGUIDv7 生成带连字符的 UUID v7。
func NewGUIDv7() string {
	return formatUUID(uuid.NewV7(), true)
}

// NewGUIDv7NoHyphen 生成不带连字符的 UUID v7。
func NewGUIDv7NoHyphen() string {
	return formatUUID(uuid.NewV7(), false)
}

// NewShortUUID 生成短格式 UUID。
func NewShortUUID() string {
	u := uuid.NewV4()
	return base64.RawURLEncoding.EncodeToString(u[:])
}

// NewXID 生成 XID。
func NewXID() string {
	return xid.New().String()
}
