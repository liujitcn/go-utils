package id

import (
	"testing"
	"time"
)

func TestGetTimestamp(t *testing.T) {
	genTS := epoch + 123456789
	sequence := int64(7)
	sid := ((genTS - epoch) << timestampShift) | sequence

	if got := GetTimestamp(sid); got != genTS-epoch {
		t.Fatalf("GetTimestamp() = %d, want %d", got, genTS-epoch)
	}
	if got := GetGenTimestamp(sid); got != genTS {
		t.Fatalf("GetGenTimestamp() = %d, want %d", got, genTS)
	}
}

func TestGetGenTime(t *testing.T) {
	genTS := epoch + 123456789
	sid := (genTS - epoch) << timestampShift

	got := GetGenTime(sid)
	want := time.UnixMilli(genTS).Format("2006-01-02 15:04:05.000")
	if got != want {
		t.Fatalf("GetGenTime() = %s, want %s", got, want)
	}
}

func TestGenSnowflakeIDMonotonic(t *testing.T) {
	id1 := GenSnowflakeID()
	id2 := GenSnowflakeID()
	if id2 <= id1 {
		t.Fatalf("GenSnowflakeID() not monotonic: id1=%d id2=%d", id1, id2)
	}
}

