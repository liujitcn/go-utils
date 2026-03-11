package id

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
)

func formatUUID(u uuid.UUID, withHyphen bool) string {
	if withHyphen {
		return u.String()
	}

	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[:])
}

func NewGUIDv4() string {
	return formatUUID(uuid.New(), true)
}

func NewGUIDv4NoHyphen() string {
	return formatUUID(uuid.New(), false)
}

func NewGUIDv7() string {
	u, err := uuid.NewV7()
	if err != nil {
		// 系统时钟异常时回退到 v4，避免返回空值
		return NewGUIDv4()
	}

	return formatUUID(u, true)
}

func NewGUIDv7NoHyphen() string {
	u, err := uuid.NewV7()
	if err != nil {
		// 系统时钟异常时回退到 v4，避免返回空值
		return NewGUIDv4NoHyphen()
	}

	return formatUUID(u, false)
}

func NewShortUUID() string {
	return shortuuid.New()
}

func NewXID() string {
	return xid.New().String()
}
