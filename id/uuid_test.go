package id

import (
	"regexp"
	"testing"
)

var (
	guidWithHyphenRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	guidNoHyphenRe   = regexp.MustCompile(`^[0-9a-f]{32}$`)
)

func TestNewGUIDv4Format(t *testing.T) {
	guid := NewGUIDv4()
	if !guidWithHyphenRe.MatchString(guid) {
		t.Fatalf("NewGUIDv4() format invalid: %s", guid)
	}

	guidNoHyphen := NewGUIDv4NoHyphen()
	if !guidNoHyphenRe.MatchString(guidNoHyphen) {
		t.Fatalf("NewGUIDv4NoHyphen() format invalid: %s", guidNoHyphen)
	}
}

func TestNewGUIDv7Format(t *testing.T) {
	guid := NewGUIDv7()
	if !guidWithHyphenRe.MatchString(guid) {
		t.Fatalf("NewGUIDv7() format invalid: %s", guid)
	}
	if guid[14] != '7' {
		// v7 的版本位应为 7（含连字符格式位于索引 14）
		t.Fatalf("NewGUIDv7() version invalid: %s", guid)
	}

	guidNoHyphen := NewGUIDv7NoHyphen()
	if !guidNoHyphenRe.MatchString(guidNoHyphen) {
		t.Fatalf("NewGUIDv7NoHyphen() format invalid: %s", guidNoHyphen)
	}
	if guidNoHyphen[12] != '7' {
		// v7 的版本位应为 7（无连字符格式位于索引 12）
		t.Fatalf("NewGUIDv7NoHyphen() version invalid: %s", guidNoHyphen)
	}
}
