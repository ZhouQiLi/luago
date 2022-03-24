package state

type luaStack struct {
	slots []luaValue
	top   int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
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
	if index > 0 {
		return index
	}

	return index + self.top + 1
}

func (self *luaStack) isValid(index int) bool {
	absIndex := self.absIndex(index)
	return absIndex > 0 && absIndex <= self.top
}

func (self *luaStack) get(index int) luaValue {
	absIndex := self.absIndex(index)
	if absIndex > 0 && absIndex <= self.top {
		return self.slots[absIndex-1]
	}
	return nil
}

func (self *luaStack) set(index int, value luaValue) {
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
