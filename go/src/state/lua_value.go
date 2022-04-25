package state

import (
	"fmt"
	. "luago/api"
	"luago/number"
)

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *closure:
		return LUA_TFUNCTION
	default:
		panic("Todo!")
	}
}

func convertToFloat(value luaValue) (float64, bool) {
	switch x := value.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func stringToInteger(str string) (int64, bool) {
	if i, ok := number.ParseInteger(str); ok {
		return i, ok
	}
	if f, ok := number.ParseFloat(str); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}

func convertToInteger(value luaValue) (int64, bool) {
	switch x := value.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return stringToInteger(x)
	default:
		return 0, false
	}
}

func convertToBoolean(value luaValue) bool {
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}

func getMetatable(value luaValue, L *luaState) *luaTable {
	if t, ok := value.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(value))
	if mt := L.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}

	return nil
}

func setMetatable(value luaValue, mt *luaTable, L *luaState) {
	if t, ok := value.(*luaTable); ok {
		t.metatable = mt
		return
	}
	key := fmt.Sprintf("_MT%d", typeOf(value))
	L.registry.put(key, mt)
}

func getMetaField(value luaValue, fieldName string, L *luaState) luaValue {
	if metatable := getMetatable(value, L); metatable != nil {
		return metatable.get(fieldName)
	}
	return nil
}

func callMetamethod(a, b luaValue, metamethodName string, L *luaState) (luaValue, bool) {
	var metamethod luaValue
	if metamethod = getMetaField(a, metamethodName, L); metamethod == nil {
		if metamethod = getMetaField(b, metamethodName, L); metamethod == nil {
			return nil, false
		}
	}

	L.stack.check(4)
	L.stack.push(metamethod)
	L.stack.push(a)
	L.stack.push(b)
	L.Call(2, 1)
	return L.stack.pop(), true
}
