package snowflake

import (
	"sync"
	"time"
)

const (
	epoch          = int64(1577808000000)              // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69*1024年
	timestampBits  = uint(51)                          // 时间戳占用位数
	sequenceBits   = uint(12)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	timestampShift = sequenceBits                      // 时间戳左移位数
)

type Snowflake struct {
	sync.Mutex
	timestamp int64
	sequence  int64
}

func NewSnowflake() (*Snowflake, error) {
	return &Snowflake{
		timestamp: 0,
		sequence:  0,
	}, nil
}

func (s *Snowflake) NextVal() int64 {
	s.Lock()
	now := time.Now().UnixMilli() // 毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return 0
	}
	s.timestamp = now
	r := (t)<<timestampShift | (s.sequence)
	s.Unlock()
	return r
}

// GetTimestamp 获取时间戳
func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

// GetGenTimestamp 获取创建ID时的时间戳
func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + epoch
	return
}

// GetGenTime 获取创建ID时的时间字符串(精度：秒)
func GetGenTime(sid int64) (t string) {
	// 需将GetGenTimestamp获取的时间戳转换成秒
	t = time.Unix(0, GetGenTimestamp(sid)*1e3).Format("2006-01-02 15:04:05.000")
	return
}
