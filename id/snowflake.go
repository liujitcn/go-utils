package id

import (
	"sync"
	"time"
)

var (
	sf     *snowflake
	sfOnce sync.Once
)

const (
	epoch          = int64(1577808000000)              // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69*1024年
	timestampBits  = uint(51)                          // 时间戳占用位数
	sequenceBits   = uint(12)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	timestampShift = sequenceBits                      // 时间戳左移位数
)

func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + epoch
	return
}

func GenSnowflakeID() int64 {
	return newSnowflake().nextVal()
}

func GetGenTime(sid int64) (t string) {
	// GetGenTimestamp 返回的是毫秒时间戳，这里直接按毫秒转换
	t = time.UnixMilli(GetGenTimestamp(sid)).Format("2006-01-02 15:04:05.000")
	return
}

type snowflake struct {
	sync.Mutex
	timestamp int64
	sequence  int64
}

func newSnowflake() *snowflake {
	sfOnce.Do(func() {
		sf = &snowflake{
			timestamp: 0,
			sequence:  0,
		}
	})
	return sf
}

func (s *snowflake) nextVal() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixMilli() // 毫秒
	if now < s.timestamp {
		// 时钟回拨时等待到上一次时间戳，避免出现重复或倒序 ID
		now = waitNextMillis(s.timestamp)
	}
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			now = waitNextMillis(s.timestamp)
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		return 0
	}
	s.timestamp = now
	return (t)<<timestampShift | (s.sequence)
}

func waitNextMillis(lastTimestamp int64) int64 {
	now := time.Now().UnixMilli()
	for now <= lastTimestamp {
		// 短暂休眠避免空转占满 CPU
		time.Sleep(time.Microsecond * 200)
		now = time.Now().UnixMilli()
	}
	return now
}
