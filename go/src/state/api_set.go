package state

import (
	. "luago/api"
)

func (self *luaState) SetGlobal(name string) {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v)
}

func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}
