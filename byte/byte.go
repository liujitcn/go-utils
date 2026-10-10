package byte

import (
	"encoding/binary"
)

// IntToBytes 将int转换为[]byte
func IntToBytes(n int) []byte {
	return binary.BigEndian.AppendUint64(nil, uint64(int64(n)))
}

// BytesToInt 将[]byte转换为int
func BytesToInt(bys []byte) int {
	if len(bys) < 8 {
		return 0
	}
	return int(int64(binary.BigEndian.Uint64(bys)))
}

// ByteToLower 将 ASCII 大写字节转换为小写。
func ByteToLower(b byte) byte {
	if b <= '\u007F' {
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		return b
	}
	return b
}

// ByteToUpper 将 ASCII 小写字节转换为大写。
func ByteToUpper(b byte) byte {
	if b <= '\u007F' {
		if 'a' <= b && b <= 'z' {
			b -= 'a' - 'A'
		}
		return b
	}
	return b
}
