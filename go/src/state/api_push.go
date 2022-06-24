package state

import (
	"fmt"
	. "luago/api"
)

func (self *luaState) PushNil() {
	self.stack.push(nil)
}

func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

func (self *luaState) PushFString(fmtStr string, a ...interface{}) {
	str := fmt.Sprintf(fmtStr, a...)
	self.stack.push(str)
}

func (self *luaState) PushGoFunction(f GoFunction) {
	self.stack.push(newGoClosure(f, 0))
}

func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

func (self *luaState) PushGoClosure(f GoFunction, upvalueCount int) {
	closure := newGoClosure(f, upvalueCount)
	for i := upvalueCount; i > 0; i-- {
		value := self.stack.pop()
		closure.upvalues[i-1] = &upvalue{&value}
	}
	self.stack.push(closure)
}
