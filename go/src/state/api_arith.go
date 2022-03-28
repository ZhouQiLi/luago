package state

import (
	. "luago/api"
	"luago/number"
	"math"
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	{iadd, fadd},
	{isub, fsub},
	{imul, fmul},
	{imod, fmod},
	{nil, pow},
	{nil, div},
	{iidiv, fidiv},
	{band, nil},
	{bor, nil},
	{bxor, nil},
	{shl, nil},
	{shr, nil},
	{iunm, funm},
	{bnot, nil},
}

func arith(a, b luaValue, op operator) luaValue {
	if op.integerFunc != nil {
		if ia, ok := convertToInteger(a); ok {
			if ib, ok := convertToInteger(b); ok {
				return op.integerFunc(ia, ib)
			}
		}
	}

	if fa, ok := convertToFloat(a); ok {
		if fb, ok := convertToFloat(b); ok {
			return op.floatFunc(fa, fb)
		}
	}

	return nil
}

func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = self.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		a = b
	}
	operator := operators[op]
	if result := arith(a, b, operator); result != nil {
		self.stack.push(result)
	} else {
		panic("arithmetic error")
	}
}
