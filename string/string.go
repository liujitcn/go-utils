package string

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func ConvertJsonStringToInt64Array(s string) []int64 {
	var res []int64
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res
	}
	return res
}

func ConvertJsonStringToStringArray(s string) []string {
	var res []string
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res
	}
	return res
}

func ConvertInt64ArrayToString(arr []int64) string {
	if len(arr) == 0 {
		return "[]"
	}
	marshal, err := json.Marshal(arr)
	if err != nil {
		return "[]"
	}
	return string(marshal)
}

func ConvertStringArrayToString(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	marshal, err := json.Marshal(arr)
	if err != nil {
		return "[]"
	}
	return string(marshal)
}

func ConvertStringToInt64Array(s string) []int64 {
	// 处理空字符串的情况，返回空数组
	if s == "" {
		return []int64{}
	}

	parts := strings.Split(s, ",")
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		num, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, num)
	}

	return result
}

func GetRandomString(length int) string {
	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		digits[i] = byte(n.Int64() + '0') // 生成 0~9 的 ASCII 字符
	}
	randomString := string(digits)
	return randomString
}

func DesensitizePhone(phone string) string {
	// 校验是否为 11 位数字
	if len(phone) != 11 {
		return phone
	}
	for _, c := range phone {
		if !unicode.IsDigit(c) {
			return phone
		}
	}

	// 脱敏处理
	return phone[:3] + "****" + phone[7:]
}

func ConvertAnyToJsonString(s interface{}) string {
	t := reflect.TypeOf(s)
	if t == nil || s == nil {
		return "{}"
	}

	marshal, err := json.Marshal(s)
	if err != nil {
		if t.Kind() == reflect.Slice {
			return "[]"
		}

		return "{}"
	}
	res := string(marshal)
	if res == "null" {
		if t.Kind() == reflect.Slice {
			return "[]"
		}

		return "{}"
	}

	return res
}

// ConvertYuanStringToFen 将金额元（字符串）转换为分，按“分”四舍五入
func ConvertYuanStringToFen(yuan string) int64 {
	trimmed := strings.TrimSpace(yuan)
	if trimmed == "" {
		return 0
	}

	rat := new(big.Rat)
	if _, ok := rat.SetString(trimmed); !ok {
		return 0
	}

	rat.Mul(rat, big.NewRat(100, 1))

	num := new(big.Int).Set(rat.Num())
	den := new(big.Int).Set(rat.Denom())

	q := new(big.Int).Quo(num, den)
	rem := new(big.Int).Rem(num, den)

	absDoubleRem := new(big.Int).Mul(new(big.Int).Abs(rem), big.NewInt(2))
	absDen := new(big.Int).Abs(den)
	if absDoubleRem.Cmp(absDen) >= 0 {
		if num.Sign() >= 0 {
			q.Add(q, big.NewInt(1))
		} else {
			q.Sub(q, big.NewInt(1))
		}
	}

	if !q.IsInt64() {
		return 0
	}
	return q.Int64()
}

// ConvertFenToYuanString 将金额分转换为金额元字符串（保留 2 位小数）
func ConvertFenToYuanString(fen int64) string {
	sign := ""
	if fen < 0 {
		sign = "-"
		fen = -fen
	}

	intPart := fen / 100
	fracPart := fen % 100
	return fmt.Sprintf("%s%d.%02d", sign, intPart, fracPart)
}
