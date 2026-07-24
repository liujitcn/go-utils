package mapper

import (
	"time"

	"github.com/jinzhu/copier"
	_time "github.com/liujitcn/go-utils/time"
)

// TimeTypeConverter 时间与时间字符串转换器
type TimeTypeConverter struct {
}

// NewTimeTypeConverter 创建时间与时间字符串转换器
func NewTimeTypeConverter() *TimeTypeConverter {
	return &TimeTypeConverter{}
}

// ToEntity 将时间字符串转换为时间对象
func (m *TimeTypeConverter) ToEntity(dto *string) *time.Time {
	if dto == nil {
		return nil
	}

	// 解析失败时返回 nil，保持与其他转换器一致的空值语义。
	return _time.StringTimeToTime(*dto)
}

// ToDTO 将时间对象转换为时间字符串
func (m *TimeTypeConverter) ToDTO(entity *time.Time) *string {
	if entity == nil {
		return nil
	}

	return new(_time.TimeToTimeString(*entity))
}

// NewConverterPair 创建时间与时间字符串的双向转换器
func (m *TimeTypeConverter) NewConverterPair() []copier.TypeConverter {
	fromFn := func(src time.Time) string {
		value := m.ToDTO(&src)
		if value == nil {
			return ""
		}
		return *value
	}
	toFn := func(src string) time.Time {
		value := m.ToEntity(&src)
		if value == nil {
			return time.Time{}
		}
		return *value
	}

	return NewGenericTypeConverterPair(time.Time{}, "", fromFn, toFn)
}
