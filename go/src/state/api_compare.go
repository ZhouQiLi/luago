package state

import (
	. "luago/api"
)

func eq(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	default:
		return a == b
	}
}

func lt(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		case float64:
			return x < y
		default:
			return false
		}
	}
	panic("comparison error")
}

func le(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		case float64:
			return x <= y
		default:
			return false
		}
	}
	panic("comparison error")
}

func (self *luaState) Compare(index1, index2, op CompareOp) bool {
	a := self.stack.get(index1)
	b := self.stack.get(index2)
	switch op {
	case LUA_OPEQ:
		return eq(a, b)
	case LUA_OPLT:
		return lt(a, b)
	case LUA_OPLE:
		return le(a, b)
	default:
		panic("invalid compare op")
	}
}
