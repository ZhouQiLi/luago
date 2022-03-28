package number

import (
	"strconv"
)

// TODO 自己实现解析和转换
func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

// TODO 自己实现解析和转换
func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
