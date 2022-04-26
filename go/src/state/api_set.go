package state

import (
	. "luago/api"
)

func (self *luaState) SetGlobal(name string) {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v, false)
}

func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

func (self *luaState) SetMetatable(index int) {
	value := self.stack.get(index)
	metatable := self.stack.pop()

	if metatable == nil {
		setMetatable(value, nil, self)
	} else if mt, ok := metatable.(*luaTable); ok {
		setMetatable(value, mt, self)
	} else {
		panic("table expected!")
	}
}
