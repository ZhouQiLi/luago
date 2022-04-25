package state

import (
	. "luago/api"
)

type luaStack struct {
	slots []luaValue
	top   int

	prev         *luaStack
	closure      *closure
	varargs      []luaValue
	pc           int
	state        *luaState
	openUpvalues map[int]*upvalue
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

func (self *luaStack) push(value luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow")
	}
	self.slots[self.top] = value
	self.top++
}

func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("stack underflow")
	}

	self.top--
	value := self.slots[self.top]
	self.slots[self.top] = nil

	return value
}

func (self *luaStack) absIndex(index int) int {
	if index <= LUA_REGISTRY_INDEX {
		return index
	}
	if index > 0 {
		return index
	}

	return index + self.top + 1
}

func (self *luaStack) isValid(index int) bool {
	if index < LUA_REGISTRY_INDEX {
		upvalueIndex := LUA_REGISTRY_INDEX - index - 1
		c := self.closure
		return c != nil && upvalueIndex < len(c.upvalues)
	}
	if index == LUA_REGISTRY_INDEX {
		return true
	}
	absIndex := self.absIndex(index)
	return absIndex > 0 && absIndex <= self.top
}

func (self *luaStack) get(index int) luaValue {
	if index < LUA_REGISTRY_INDEX {
		upvalueIndex := LUA_REGISTRY_INDEX - index - 1
		c := self.closure
		if c == nil || upvalueIndex >= len(c.upvalues) {
			return nil
		}
		return *(c.upvalues[upvalueIndex].value)
	}

	if index == LUA_REGISTRY_INDEX {
		return self.state.registry
	}

	absIndex := self.absIndex(index)
	if absIndex > 0 && absIndex <= self.top {
		return self.slots[absIndex-1]
	}
	return nil
}

func (self *luaStack) set(index int, value luaValue) {
	if index < LUA_REGISTRY_INDEX {
		upvalueIndex := LUA_REGISTRY_INDEX - index - 1
		c := self.closure
		if c != nil && upvalueIndex < len(c.upvalues) {
			*(c.upvalues[upvalueIndex].value) = value
		}
		return
	}

	if index == LUA_REGISTRY_INDEX {
		self.state.registry = value.(*luaTable)
		return
	}

	absIndex := self.absIndex(index)
	if absIndex > 0 && absIndex <= self.top {
		self.slots[absIndex-1] = value
		return
	}
	panic("invalid index")
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		to--
		from++
	}
}

func (self *luaStack) popN(n int) []luaValue {
	values := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		values[i] = self.pop()
	}
	return values
}

func (self *luaStack) pushN(values []luaValue, n int) {
	valueLength := len(values)
	if n < 0 {
		n = valueLength
	}
	for i := 0; i < n; i++ {
		if i < valueLength {
			self.push(values[i])
		} else {
			self.push(nil)
		}
	}
}
