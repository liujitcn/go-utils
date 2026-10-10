package xdb

import (
	"encoding/binary"
	"fmt"
)

const (
	Structure20      = 2
	Structure30      = 3
	HeaderInfoLength = 256
	VectorIndexRows  = 256
	VectorIndexCols  = 256
	VectorIndexSize  = 8
)

// --- Index policy define

type IndexPolicy int

const (
	VectorIndexPolicy IndexPolicy = 1
	BTreeIndexPolicy  IndexPolicy = 2
)

// String 返回头信息的字符串描述。
func (i IndexPolicy) String() string {
	switch i {
	case VectorIndexPolicy:
		return "VectorIndex"
	case BTreeIndexPolicy:
		return "BtreeIndex"
	default:
		return "unknown"
	}
}

// --- Header define

type Header struct {
	// data []byte
	Version       uint16
	IndexPolicy   IndexPolicy
	CreatedAt     uint32
	StartIndexPtr uint32
	EndIndexPtr   uint32

	// since IPv6 supporting
	IPVersion       int
	RuntimePtrBytes int
}

// NewHeader 从字节缓冲区解析数据库头信息。
func NewHeader(input []byte) (*Header, error) {
	if len(input) < 20 {
		return nil, fmt.Errorf("invalid input buffer")
	}

	return &Header{
		Version:         binary.LittleEndian.Uint16(input[0:2]),
		IndexPolicy:     IndexPolicy(binary.LittleEndian.Uint16(input[2:4])),
		CreatedAt:       binary.LittleEndian.Uint32(input[4:8]),
		StartIndexPtr:   binary.LittleEndian.Uint32(input[8:12]),
		EndIndexPtr:     binary.LittleEndian.Uint32(input[12:16]),
		IPVersion:       int(binary.LittleEndian.Uint16(input[16:])),
		RuntimePtrBytes: int(binary.LittleEndian.Uint16(input[18:20])),
	}, nil
}

// String 返回头信息的字符串描述。
func (h *Header) String() string {
	return fmt.Sprintf(
		"{version:%d, index_policy:%d, created_at:%d, start_index_ptr:%d, end_index_ptr:%d, ip_version:%d, runtime_ptr_bytes:%d}",
		h.Version, h.IndexPolicy, h.CreatedAt, h.StartIndexPtr, h.EndIndexPtr, h.IPVersion, h.RuntimePtrBytes,
	)
}
