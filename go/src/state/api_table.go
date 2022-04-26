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
func (self *luaState) getTable(t, key luaValue, raw bool) LuaType {
	if tb, ok := t.(*luaTable); ok {
		value := tb.get(key)
		if raw || value != nil || !tb.hasMetaField("__index") {
			self.stack.push(value)
			return typeOf(value)
		}
	}
	if !raw {
		if metaField := getMetaField(t, "__index", self); metaField != nil {
			switch x := metaField.(type) {
			case *luaTable:
				return self.getTable(x, key, false)
			case *closure:
				self.stack.push(metaField)
				self.stack.push(t)
				self.stack.push(key)
				self.Call(2, 1)
				result := self.stack.get(-1)
				return typeOf(result)
			}
		}
	}

	panic("index error!")
}

func (self *luaState) setTable(t, key, value luaValue, raw bool) {
	if tb, ok := t.(*luaTable); ok {
		if raw || tb.get(key) != nil || !tb.hasMetaField("__newindex") {
			tb.put(key, value)
			return
		}
	}

	if !raw {
		if metaField := getMetaField(t, "__newindex", self); metaField != nil {
			switch x := metaField.(type) {
			case *luaTable:
				self.setTable(x, key, value, false)
				return
			case *closure:
				self.stack.push(metaField)
				self.stack.push(t)
				self.stack.push(key)
				self.stack.push(value)
				self.Call(3, 0)
				return
			}
		}
	}

	panic("index error!")
}

func (self *luaState) GetTable(index int) LuaType {
	t := self.stack.get(index)
	key := self.stack.pop()
	return self.getTable(t, key, false)
}

func (self *luaState) SetTable(index int) {
	t := self.stack.get(index)
	value := self.stack.pop()
	key := self.stack.pop()
	self.setTable(t, key, value, false)
}

func (self *luaState) GetI(index int, i int64) LuaType {
	t := self.stack.get(index)
	return self.getTable(t, i, false)
}

func (self *luaState) SetI(index int, i int64) {
	t := self.stack.get(index)
	value := self.stack.pop()
	self.setTable(t, i, value, false)
}

func (self *luaState) GetField(index int, key string) LuaType {
	t := self.stack.get(index)
	return self.getTable(t, key, false)
}

func (self *luaState) SetField(index int, key string) {
	t := self.stack.get(index)
	value := self.stack.pop()
	self.setTable(t, key, value, false)
}

func (self *luaState) RawGet(index int) LuaType {
	t := self.stack.get(index)
	key := self.stack.pop()
	return self.getTable(t, key, true)
}

func (self *luaState) RawSet(index int) {
	t := self.stack.get(index)
	value := self.stack.pop()
	key := self.stack.pop()
	self.setTable(t, key, value, true)
}

func (self *luaState) RawGetI(index int, i int64) LuaType {
	t := self.stack.get(index)
	return self.getTable(t, i, true)
}

func (self *luaState) RawSetI(index int, i int64) {
	t := self.stack.get(index)
	value := self.stack.pop()
	self.setTable(t, i, value, true)
}

func (self *luaState) RawEqual(index1, index2 int) bool {
	if !self.stack.isValid(index1) || !self.stack.isValid(index2) {
		return false
	}
	value1 := self.stack.get(index1)
	value2 := self.stack.get(index2)
	return eq(value1, value2, nil)
}
