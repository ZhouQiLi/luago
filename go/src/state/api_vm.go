package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++
	return i
}

// 将常量数据写入栈顶
func (self *luaState) GetConst(index int) {
	c := self.stack.closure.proto.Constants[index]
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

func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}

	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

func (self *luaState) LoadProto(index int) {
	stack := self.stack
	proto := stack.closure.proto.Protos[index]
	closure := newLuaClosure(proto)
	stack.push(closure)

	for i, upvalueInfo := range proto.Upvalues {
		upvalueIndex := int(upvalueInfo.Idx)
		if upvalueInfo.Instack == 1 {
			if stack.openUpvalues == nil {
				stack.openUpvalues = map[int]*upvalue{}
			}

			if openUpvalue, found := stack.openUpvalues[upvalueIndex]; found {
				closure.upvalues[i] = openUpvalue
			} else {
				closure.upvalues[i] = &upvalue{&stack.slots[upvalueIndex]}
				stack.openUpvalues[upvalueIndex] = closure.upvalues[i]
			}
		} else {
			closure.upvalues[i] = stack.closure.upvalues[upvalueIndex]
		}
	}
}

func (self *luaState) CloseUpvalues(a int) {
	for i, openUpvalue := range self.stack.openUpvalues {
		if i >= a-1 {
			value := *openUpvalue.value
			openUpvalue.value = &value
			delete(self.stack.openUpvalues, i)
		}
	}
}
