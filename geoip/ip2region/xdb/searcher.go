package xdb

import (
	"encoding/binary"
	"fmt"
)

type Searcher struct {
	version *Version

	ioCount int

	// use it only when this feature enabled.
	// Preload the vector index will reduce the number of IO operations
	// thus speedup the search process
	vectorIndex []byte

	// content buffer.
	// running with the whole xdb file cached
	contentBuff []byte
}

// NewWithBuffer 基于完整内容缓冲区创建搜索器。
func NewWithBuffer(version *Version, cBuff []byte) (*Searcher, error) {
	return NewSearcher(version, nil, cBuff)
}

// NewSearcher 根据索引与内容缓冲区创建搜索器。
func NewSearcher(version *Version, vIndex []byte, cBuff []byte) (*Searcher, error) {
	if version == nil || (version.Bytes != 4 && version.Bytes != 16) || version.SegmentIndexSize < version.Bytes*2+6 || version.IPCompare == nil {
		return nil, fmt.Errorf("invalid IP version configuration")
	}
	if cBuff == nil && vIndex == nil {
		return nil, fmt.Errorf("searcher requires an index or content buffer")
	}
	if vIndex != nil && len(vIndex) < VectorIndexRows*VectorIndexCols*VectorIndexSize {
		return nil, fmt.Errorf("vector index buffer is too small")
	}
	if cBuff != nil && len(cBuff) < HeaderInfoLength+VectorIndexRows*VectorIndexCols*VectorIndexSize {
		return nil, fmt.Errorf("content buffer is too small")
	}

	// content buff first
	if cBuff != nil {
		return &Searcher{
			version:     version,
			vectorIndex: nil,
			contentBuff: cBuff,
		}, nil
	}

	return &Searcher{
		version:     version,
		vectorIndex: vIndex,
	}, nil
}

// Close 关闭搜索器并释放资源。
func (s *Searcher) Close() {
}

// IPVersion return the ip version
func (s *Searcher) IPVersion() *Version {
	return s.version
}

// GetIOCount return the global io count for the last search
func (s *Searcher) GetIOCount() int {
	return s.ioCount
}

// SearchByStr find the region for the specified ip string
func (s *Searcher) SearchByStr(str string) (string, error) {
	ip, err := ParseIP(str)
	if err != nil {
		return "", err
	}

	return s.Search(ip)
}

// Search 根据 IP 字节切片查询归属地信息。
func (s *Searcher) Search(ip []byte) (string, error) {
	if s == nil || s.version == nil {
		return "", fmt.Errorf("searcher is not initialized")
	}

	// ip version check
	if len(ip) != s.version.Bytes {
		return "", fmt.Errorf("invalid ip address(%s expected)", s.version.Name)
	}

	// reset the global ioCount
	s.ioCount = 0

	// locate the segment index block based on the vector index
	var il0, il1 = int(ip[0]), int(ip[1])
	var idx = il0*VectorIndexCols*VectorIndexSize + il1*VectorIndexSize
	var sPtr, ePtr uint32
	if s.vectorIndex != nil {
		if idx+VectorIndexSize > len(s.vectorIndex) {
			return "", fmt.Errorf("vector index offset %d is out of range", idx)
		}
		sPtr = binary.LittleEndian.Uint32(s.vectorIndex[idx:])
		ePtr = binary.LittleEndian.Uint32(s.vectorIndex[idx+4:])
	} else if s.contentBuff != nil {
		if HeaderInfoLength+idx+VectorIndexSize > len(s.contentBuff) {
			return "", fmt.Errorf("vector index offset %d is out of range", idx)
		}
		sPtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx:])
		ePtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx+4:])
	} else {
		// read the vector index block
		var buff = make([]byte, VectorIndexSize)
		err := s.read(int64(HeaderInfoLength+idx), buff)
		if err != nil {
			return "", fmt.Errorf("read vector index block at %d: %w", HeaderInfoLength+idx, err)
		}

		sPtr = binary.LittleEndian.Uint32(buff)
		ePtr = binary.LittleEndian.Uint32(buff[4:])
	}

	// fmt.Printf("sPtr=%d, ePtr=%d\n", sPtr, ePtr)

	// binary search the segment index to get the region
	var bytes, dBytes = len(ip), len(ip) << 1
	var segIndexSize = uint32(s.version.SegmentIndexSize)
	if segIndexSize < uint32(dBytes+6) {
		return "", fmt.Errorf("invalid segment index size: %d", segIndexSize)
	}
	if ePtr < sPtr || (ePtr-sPtr)%segIndexSize != 0 {
		return "", fmt.Errorf("invalid segment index range: %d-%d", sPtr, ePtr)
	}
	var dataLen, dataPtr = 0, uint32(0)
	var buff = make([]byte, segIndexSize)
	var l, h = 0, int((ePtr - sPtr) / segIndexSize)
	for l <= h {
		m := (l + h) >> 1
		p := sPtr + uint32(m)*segIndexSize
		err := s.read(int64(p), buff)
		if err != nil {
			return "", fmt.Errorf("read segment index at %d: %w", p, err)
		}

		// decode the data step by step to reduce the unnecessary operations
		if s.version.IPCompare(ip, buff[0:bytes]) < 0 {
			h = m - 1
		} else if s.version.IPCompare(ip, buff[bytes:dBytes]) > 0 {
			l = m + 1
		} else {
			dataLen = int(binary.LittleEndian.Uint16(buff[dBytes:]))
			dataPtr = binary.LittleEndian.Uint32(buff[dBytes+2:])
			break
		}
	}

	// fmt.Printf("dataLen: %d, dataPtr: %d\n", dataLen, dataPtr)
	if dataLen == 0 {
		return "", nil
	}

	// load and return the region data
	var regionBuff = make([]byte, dataLen)
	err := s.read(int64(dataPtr), regionBuff)
	if err != nil {
		return "", fmt.Errorf("read region at %d: %w", dataPtr, err)
	}

	return string(regionBuff), nil
}

// read 从内容缓冲区读取指定范围的数据。
func (s *Searcher) read(offset int64, buff []byte) error {
	if s.contentBuff != nil {
		if offset < 0 || offset > int64(len(s.contentBuff)) || int64(len(buff)) > int64(len(s.contentBuff))-offset {
			return fmt.Errorf("read range %d-%d is out of bounds", offset, offset+int64(len(buff)))
		}
		copy(buff, s.contentBuff[int(offset):int(offset)+len(buff)])
		return nil
	}
	return fmt.Errorf("searcher has no content buffer")
}
