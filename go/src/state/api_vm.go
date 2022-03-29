package state

func (self *luaState) PC() int {
	return self.pc
}

func (self *luaState) AddPC(n int) {
	self.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++
	return i
}

func (self *luaState) GetConst(index int) {
	c := self.proto.Constants[index]
	self.stack.push(c)
}

func (self *luaState) GetRK(rk int) {
	// 只有255个寄存器, 当rk值大于255则表示是常量
	if rk > 0xff {
		self.GetConst(rk & 0xff)
	} else {
		// 栈索引是从1开始的, 所以对于寄存器的值要+1
		self.PushValue(rk + 1)
	}
}
