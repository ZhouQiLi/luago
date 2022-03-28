package number

import "math"

// lua中的整除是向负无穷方向取整, golang是将结果截断向0取整, 所以对于整除操作符要有特殊处理
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

func ShiftLeft(a, n int64) int64 {
	if n > 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}

// golang本身的右移运算符是有符号右移(补1), 但在lua中需要的是无符号右移(补0)
// 所以需要先转成无符号整数再做运算。
func ShiftRight(a, n int64) int64 {
	if n > 0 {
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

// 仅当浮点数部分为0且整数部分在lua整数能够表示的范围内才成功
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}
