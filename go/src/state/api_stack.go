package state

func (self *luaState) GetTop() int {
	return self.stack.top
}

// 将相对索引转换成绝对索引。
func (self *luaState) AbsIndex(index int) int {
	return self.stack.absIndex(index)
}

// 使栈的大小有n个槽, 不会返回false
func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true
}

func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

// 从fromIndex的位置上的值拷贝到toIndex
func (self *luaState) Copy(fromIndex, toIndex int) {
	value := self.stack.get(fromIndex)
	self.stack.set(toIndex, value)
}

// 将index位置的值push到栈顶
func (self *luaState) PushValue(index int) {
	value := self.stack.get(index)
	self.stack.push(value)
}

// 将栈顶的值出栈并赋值到index位置
func (self *luaState) Replace(index int) {
	value := self.stack.pop()
	self.stack.set(index, value)
}

// 将栈顶值弹出并插入到指定index位置中
func (self *luaState) Insert(index int) {
	self.Rotate(index, 1)
}

// 删除指定index的值, 并使index之上的值下移
func (self *luaState) Remove(index int) {
	self.Rotate(index, -1)
	self.Pop(1)
}

// 旋转, 使用三次反转来处理
func (self *luaState) Rotate(index, n int) {
	topIndex := self.stack.top - 1
	startIndex := self.stack.absIndex(index) - 1
	var midIndex int
	if n >= 0 {
		midIndex = topIndex - n
	} else {
		midIndex = startIndex - n - 1
	}
	self.stack.reverse(startIndex, midIndex)
	self.stack.reverse(midIndex+1, topIndex)
	self.stack.reverse(startIndex, topIndex)
}

// 若当前栈顶索引大于index, 则将栈空间缩小至index, 否则使用nil值填充栈至index大小
func (self *luaState) SetTop(index int) {
	newTop := self.AbsIndex(index)
	if newTop < 0 {
		panic("stack underflow!")
	}
	top := self.stack.top
	if top == newTop {
		return
	} else if top > newTop {
		self.Pop(top - newTop)
	} else {
		for i := top; i < newTop; i++ {
			self.stack.push(nil)
		}
	}
}
