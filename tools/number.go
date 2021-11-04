package commonUtil

import (
	"strconv"
)

func StrToInt(str string) int {
	if str == "" {
		return 0
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

func IntToStr(val int) string {
	return strconv.Itoa(val)
}

func Int64ToStr(val int64) string {
	return strconv.FormatInt(val, 10)
}

func UintToStr(val uint) string {
	return strconv.FormatInt(int64(val), 10)
}
