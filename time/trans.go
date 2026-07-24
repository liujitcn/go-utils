package time

import (
	"time"

	"github.com/liujitcn/go-utils/trans"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var DefaultTimeLocation *time.Location

var stringTimeLayouts = []string{
	Layout,
	DateLayout,
	ClockLayout,
}

var stringDateLayouts = []string{
	DateLayout,
	DateCompactLayout,
	DateSlashLayout,
	Layout,
	ClockLayout,
}

// RefreshDefaultTimeLocation 刷新默认时区配置。
func RefreshDefaultTimeLocation(name string) {
	DefaultTimeLocation, _ = time.LoadLocation(name)
}

// UnixMilliToStringPtr 毫秒时间戳 -> 字符串
func UnixMilliToStringPtr(tm *int64) *string {
	if tm == nil {
		return nil
	}
	return new(time.UnixMilli(*tm).Format(Layout))
}

// StringToUnixMilliInt64Ptr 字符串 -> 毫秒时间戳
func StringToUnixMilliInt64Ptr(tm string) int64 {
	theTime := StringTimeToTime(tm)
	if theTime.IsZero() {
		return 0
	}
	unixTime := theTime.UnixMilli()
	return unixTime
}

// StringTimeToTime 时间字符串 -> 时间
func StringTimeToTime(str string) *time.Time {
	return parseStringToTimeByLayouts(str, stringTimeLayouts)
}

// TimeToTimeString 时间 -> 时间字符串
func TimeToTimeString(tm time.Time) string {
	if tm.IsZero() {
		return ""
	}
	return tm.Format(Layout)
}

// StringDateToTime 字符串 -> 时间
func StringDateToTime(str *string) *time.Time {
	if str == nil {
		return nil
	}
	return parseStringToTimeByLayouts(*str, stringDateLayouts)
}

// TimeToDateString 时间 -> 日期字符串
func TimeToDateString(tm time.Time) string {
	if tm.IsZero() {
		return ""
	}
	return tm.Format(DateLayout)
}

// ensureDefaultTimeLocation 确保默认时区已经初始化。
func ensureDefaultTimeLocation() {
	if DefaultTimeLocation == nil {
		RefreshDefaultTimeLocation(DefaultTimeLocationName)
	}
}

// parseStringToTimeByLayouts 按给定格式列表顺序解析时间字符串。
func parseStringToTimeByLayouts(str string, layouts []string) *time.Time {
	if len(str) == 0 {
		return nil
	}

	ensureDefaultTimeLocation()

	var err error
	var theTime time.Time

	for _, layout := range layouts {
		// 按顺序尝试，优先匹配调用方定义的主格式。
		theTime, err = time.ParseInLocation(layout, str, DefaultTimeLocation)
		if err == nil {
			return &theTime
		}
	}

	return nil
}

// TimestamppbToTime timestamppb.Timestamp -> time.Time
func TimestamppbToTime(tm *timestamppb.Timestamp) *time.Time {
	if tm != nil {
		return new(trans.Ptr(tm.AsTime()))
	}
	return nil
}

// TimeToTimestamppb time.Time -> timestamppb.Timestamp
func TimeToTimestamppb(tm time.Time) *timestamppb.Timestamp {
	if !tm.IsZero() {
		return timestamppb.New(tm)
	}
	return nil
}
