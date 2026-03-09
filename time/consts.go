package time

import "time"

const (
	DateLayout  = "2006-01-02"
	ClockLayout = "15:04:05"
	Layout      = DateLayout + " " + ClockLayout

	DefaultTimeLocationName = "Asia/Shanghai"
)

var ReferenceTimeValue = time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.FixedZone("MST", -7*60*60))
