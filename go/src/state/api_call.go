package state

import (
	"fmt"
	"luago/binary_chunk"
	"luago/vm"
)

func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binary_chunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)
	return 0
}

func (self *luaState) Call(argsCount, resultCount int) {
	value := self.stack.get(-(argsCount + 1))
	if c, ok := value.(*closure); ok {
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
		self.callLuaClosure(argsCount, resultCount, c)
	} else {
		panic("not function")
	}
}

func (self *luaState) callLuaClosure(argsCount, resultCount int, c *closure) {
	maxStackSize := int(c.proto.MaxStackSize)
	paramCount := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(maxStackSize + 20)
	newStack.closure = c

	funcAndArgs := self.stack.popN(argsCount + 1)
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
