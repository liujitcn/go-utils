# go-utils

`go-utils` 是一个通用 Go 工具库，当前包含以下能力：

- 密码加密与校验（`crypto`）
- 字符串与 JSON/切片转换（`str`）
- `[]int64` 集合等价比较（`slice`）
- 简单雪花 ID 生成（`snowflake`）
- 时间格式化、区间、差值、protobuf 时间转换（`timeutil`）
- 指针/值互转与基础集合转换（`trans`）
- 一些聚合工具函数（根包 `utils`）

## 安装

```bash
go get github.com/liujitcn/go-utils
```

## 包说明

### 根包 `utils`

- `CalcGrowthRate(prev, current int64) int64`：增长率（万分比，10000 表示 100%）
- `GetCreatedAt(timeType TimeType) (time.Time, time.Time)`：按类型返回统计起止时间
- `FormatDate(timeType TimeType, key int) string`：格式化周/月/日展示文本
- `ConvertYuanToFen(yuan string) int64`：元转分

`TimeType`：

- `TimeTypeWeek`
- `TimeTypeMonth`

示例：

```go
package main

import (
    "fmt"

    utils "github.com/liujitcn/go-utils"
)

func main() {
    fmt.Println(utils.CalcGrowthRate(100, 120)) // 2000 => 20.00%
    fmt.Println(utils.ConvertYuanToFen("12.34"))
}
```

### `crypto`

- `HashPassword(password string) (string, error)`：bcrypt 加密
- `CheckPasswordHash(password, hash string) bool`：bcrypt 校验

```go
hash, err := crypto.HashPassword("123456")
ok := crypto.CheckPasswordHash("123456", hash)
_ = err
_ = ok
```

### `str`

主要函数：

- JSON 字符串与数组互转：
- `ConvertJsonStringToInt64Array`
- `ConvertJsonStringToStringArray`
- `ConvertInt64ArrayToString`
- `ConvertStringArrayToString`
- 普通字符串转 `[]int64`：`ConvertStringToInt64Array`（逗号分隔）
- 生成随机数字串：`GetRandomString(length int)`
- 手机号脱敏：`DesensitizePhone(phone string)`
- 任意值转 JSON 字符串：`ConvertAnyToJsonString(any)`

```go
ids := str.ConvertStringToInt64Array("1,2, 3")
jsonStr := str.ConvertInt64ArrayToString(ids)
phone := str.DesensitizePhone("13812345678")
_, _ = jsonStr, phone
```

### `slice`

- `Equal(ids1, ids2 []int64) bool`：比较两组 `int64` 是否等价（内部会排序）
- `EqualByString(ids1Str, ids2Str string) bool`：比较两组 JSON 字符串形式的 `[]int64`

```go
same := slice.Equal([]int64{3, 1, 2}, []int64{1, 2, 3})
_ = same
```

### `snowflake`

- `NewSnowflake() (*Snowflake, error)`
- `(*Snowflake).NextVal() int64`
- `GetTimestamp(sid int64) int64`
- `GetGenTimestamp(sid int64) int64`
- `GetGenTime(sid int64) string`

```go
node, _ := snowflake.NewSnowflake()
id := node.NextVal()
createdAt := snowflake.GetGenTime(id)
_, _ = id, createdAt
```

### `timeutil`

常量：

- `DateLayout = "2006-01-02"`
- `ClockLayout = "15:04:05"`
- `TimeLayout = "2006-01-02 15:04:05"`

功能分类：

- 时间格式化与 duration：
- `ReferenceTime`
- `FormatTimer`
- `FormatTimerf`
- `DurationHMS`
- protobuf duration 转换：
- `Float64ToDurationpb`
- `SecondToDurationpb`
- `DurationpbSecond`
- 字符串/时间/毫秒时间戳转换：
- `UnixMilliToStringPtr`
- `StringToUnixMilliInt64Ptr`
- `StringTimeToTime`
- `StringDateToTime`
- `TimeToTimeString`
- `TimeToDateString`
- protobuf timestamp 转换：
- `TimestamppbToTime`
- `TimeToTimestamppb`
- 日期差值：
- `DayDifferenceHours`
- `StringDifferenceDays`
- `DayTimeDifferenceHours`
- `TimeDifferenceDays`
- `DaySecondsDifferenceHours`
- `SecondsDifferenceDays`
- 常用时间区间（昨天/今天/上月/本月/去年/今年）：
- `GetYesterdayRangeTime`
- `GetTodayRangeTime`
- `GetLastMonthRangeTime`
- `GetCurrentMonthRangeTime`
- `GetLastYearRangeTime`
- `GetCurrentYearRangeTime`
- 及对应 `RangeDateString` / `RangeTimeString` 版本

```go
start, end := timeutil.GetTodayRangeTimeString()
fmt.Println(start, end)
```

### `trans`

基础值与指针互转：

- `String/StringValue`
- `Int/IntValue`
- `Int8/Int8Value`
- `Int16/Int16Value`
- `Int32/Int32Value`
- `Int64/Int64Value`
- `Uint/UintValue`
- `Uint8/Uint8Value`
- `Uint16/Uint16Value`
- `Uint32/Uint32Value`
- `Uint64/Uint64Value`
- `Float32/Float32Value`
- `Float64/Float64Value`
- `Bool/BoolValue`
- `Time/TimeValue`

切片互转：

- `IntSlice/IntValueSlice`、`Int8Slice/Int8ValueSlice`、`Int16Slice/Int16ValueSlice`
- `Int32Slice/Int32ValueSlice`、`Int64Slice/Int64ValueSlice`
- `UintSlice/UintValueSlice`、`Uint8Slice/Uint8ValueSlice`、`Uint16Slice/Uint16ValueSlice`
- `Uint32Slice/Uint32ValueSlice`、`Uint64Slice/Uint64ValueSlice`
- `Float32Slice/Float32ValueSlice`、`Float64Slice/Float64ValueSlice`
- `StringSlice/StringSliceValue`
- `BoolSlice/BoolSliceValue`

其他：

- `Enum/EnumValue`
- `MapKeys/MapValues`
- `Ptr`、`SliceOfPtrs`

```go
name := trans.String("alice")
nameVal := trans.StringValue(name)
_ = nameVal
```

## 测试

```bash
go test ./...
```

如果你的环境出现 `exec format error`，通常是本地 Go 工具链与测试二进制架构不一致，先检查 `go env GOOS GOARCH` 与当前机器架构。
