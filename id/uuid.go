package id

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewGUIDv4(withHyphen bool) string {
	u := uuid.New()

	if withHyphen {
		return u.String()
	}

	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[:])
}

func NewGUIDv7(withHyphen bool) string {
	u, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if system clock is unreliable
		return NewGUIDv4(withHyphen)
	}

	if withHyphen {
		return u.String()
	}

	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[:])
}

func NewShortUUID() string {
	return shortuuid.New()
}

func NewKSUID() string {
	return ksuid.New().String()
}

func NewXID() string {
	return xid.New().String()
}

func NewMongoObjectID() string {
	objID := bson.NewObjectID()
	return objID.String()
}
