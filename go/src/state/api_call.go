package state

import (
	. "luago/api"
	"luago/binary_chunk"
	"luago/vm"
)

func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binary_chunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)

	if len(proto.Upvalues) > 0 {
		env := self.registry.get(LUA_RIDX_GLOBALS)
		c.upvalues[0] = &upvalue{&env}
	}

	return 0
}

func (self *luaState) Call(argsCount, resultCount int) {
	value := self.stack.get(-(argsCount + 1))
	c, ok := value.(*closure)

	if !ok {
		// 原本传入的参数中没有函数名(value)这一项, 对于有__call元表的需要传入函数名, 至于位置是(-(argsCount + 1)或(-(argsCount + 2)不重要因为内容是一样的。
		if metaField := getMetaField(value, "__call", self); metaField != nil {
			if c, ok = metaField.(*closure); ok {
				self.stack.push(value)
				self.Insert(-(argsCount + 2))
				argsCount++
			}
		}
	}

	if ok {
		if c.proto != nil {
			self.callLuaClosure(argsCount, resultCount, c)
		} else {
			self.callGoClosure(argsCount, resultCount, c)
		}
	}
}

func (self *luaState) callLuaClosure(argsCount, resultCount int, c *closure) {
	maxStackSize := int(c.proto.MaxStackSize)
	paramCount := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(maxStackSize+LUA_MINSTACK, self)
	newStack.closure = c

	funcAndArgs := self.stack.popN(argsCount + 1)
	newStack.pushN(funcAndArgs[1:], paramCount)
	newStack.top = maxStackSize
	if argsCount > paramCount && isVararg {
		newStack.varargs = funcAndArgs[paramCount+1:]
	}

	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()
	if resultCount != 0 {
		results := newStack.popN(newStack.top - maxStackSize)
		self.stack.check(len(results))
		self.stack.pushN(results, resultCount)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (self *luaState) callGoClosure(argsCount, resultCount int, c *closure) {
	newStack := newLuaStack(argsCount+LUA_MINSTACK, self)
	newStack.closure = c

	args := self.stack.popN(argsCount)
	newStack.pushN(args, argsCount)
	// 与lua函数的区别是这里调用的golang函数不需要从栈中获取, 但call指令会默认把调用的function和参数压入栈中, 所以这里要pop出无用的值
	self.stack.pop()

	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	if resultCount != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, resultCount)
	}
}
