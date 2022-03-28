package state

import (
	"luago/api"
	"strconv"
)

func (self *luaState) TypeName(luaType api.LuaType) string {
	switch luaType {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (self *luaState) Type(index int) api.LuaType {
	if self.stack.isValid(index) {
		value := self.stack.get(index)
		return typeOf(value)
	}

	return api.LUA_TNONE
}

func (self *luaState) IsNone(index int) bool {
	return self.Type(index) == api.LUA_TNONE
}

func (self *luaState) IsNil(index int) bool {
	return self.Type(index) == api.LUA_TNIL
}

func (self *luaState) IsNoneOrNil(index int) bool {
	return self.Type(index) <= api.LUA_TNIL
}

func (self *luaState) IsBoolean(index int) bool {
	return self.Type(index) == api.LUA_TBOOLEAN
}

func (self *luaState) IsString(index int) bool {
	t := self.Type(index)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

func (self *luaState) IsNumber(index int) bool {
	return false
}

func (self *luaState) IsInteger(index int) bool {
	value := self.stack.get(index)
	_, ok := value.(int64)
	return ok
}

func covertToBoolean(value luaValue) bool {
	switch x := value.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (self *luaState) ToBoolean(index int) bool {
	value := self.stack.get(index)
	return covertToBoolean(value)
}

func (self *luaState) ToNumberX(index int) (float64, bool) {
	value := self.stack.get(index)
	return convertToFloat(value)
}

func (self *luaState) ToNumber(index int) float64 {
	value, _ := self.ToNumberX(index)
	return value
}

func (self *luaState) ToIntegerX(index int) (int64, bool) {
	value := self.stack.get(index)
	return convertToInteger(value)
}

func (self *luaState) ToInteger(index int) int64 {
	value, _ := self.ToIntegerX(index)
	return value
}

func (self *luaState) ToStringX(index int) (string, bool) {
	value := self.stack.get(index)
	switch x := value.(type) {
	case string:
		return x, true
	case int64:
		return strconv.FormatInt(x, 10), true
	case float64:
		return strconv.FormatFloat(x, 'e', 30, 32), true
	default:
		return "", false
	}
}

func (self *luaState) ToString(index int) string {
	value, _ := self.ToStringX(index)
	return value
}
