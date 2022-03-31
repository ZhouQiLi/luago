package state

import (
	. "luago/api"
)

func (self *luaState) CreateTable(arrSize, mapSize int) {
	t := newLuaTable(arrSize, mapSize)
	self.stack.push(t)
}

func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}
func (self *luaState) getTable(t, key luaValue) LuaType {
	if tb, ok := t.(*luaTable); ok {
		value := tb.get(key)
		self.stack.push(value)
		return typeOf(value)
	} else {
		// TODO
		panic("is not table!")
	}
}

func (self *luaState) setTable(t, key, value luaValue) {
	if tb, ok := t.(*luaTable); ok {
		tb.put(key, value)
		return
	} else {
		// TODO
		panic("is not table!")
	}
}

func (self *luaState) GetTable(index int) LuaType {
	t := self.stack.get(index)
	key := self.stack.pop()
	return self.getTable(t, key)
}

func (self *luaState) GetField(index int, key string) LuaType {
	t := self.stack.get(index)
	return self.getTable(t, key)
}

func (self *luaState) GetI(index int, i int64) LuaType {
	t := self.stack.get(index)
	return self.getTable(t, i)
}

func (self *luaState) SetTable(index int) {
	t := self.stack.get(index)
	value := self.stack.pop()
	key := self.stack.pop()
	self.setTable(t, key, value)
}

func (self *luaState) SetField(index int, key string) {
	t := self.stack.get(index)
	value := self.stack.pop()
	self.setTable(t, key, value)
}

func (self *luaState) SetI(index int, i int64) {
	t := self.stack.get(index)
	value := self.stack.pop()
	self.setTable(t, i, value)
}
