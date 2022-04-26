package state

import (
	. "luago/api"
)

func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t, name, false)
}

func (self *luaState) GetMetatable(index int) bool {
	value := self.stack.get(index)
	if metatable := getMetatable(value, self); metatable != nil {
		self.stack.push(metatable)
		return true
	} else {
		return false
	}
}
