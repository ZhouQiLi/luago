package state

import (
	. "luago/api"
)

type luaState struct {
	registry *luaTable
	stack    *luaStack
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	L := &luaState{registry: registry}
	L.pushLuaStack(newLuaStack(LUA_MINSTACK, L))

	return L
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}
