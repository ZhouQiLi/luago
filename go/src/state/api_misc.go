package state

import "luago/number"

func (self *luaState) Len(index int) {
	value := self.stack.get(index)
	if s, ok := value.(string); ok {
		self.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(value, value, "__len", self); ok {
		self.stack.push(result)
	} else if t, ok := value.(*luaTable); ok {
		self.stack.push(int64(t.len()))
	} else {
		panic("length error")
	}
}

func (self *luaState) RawLen(index int) uint {
	value := self.stack.get(index)
	switch x := value.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}

			b := self.stack.pop()
			a := self.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", self); ok {
				self.stack.push(result)
				continue
			}

			panic("concatenation error")
		}
	}
}

func (self *luaState) Next(index int) bool {
	value := self.stack.get(index)
	if tab, ok := value.(*luaTable); ok {
		key := self.stack.pop()
		if nextKey := tab.nextKey(key); nextKey != nil {
			self.stack.push(nextKey)
			self.stack.push(tab.get(nextKey))
			return true
		}
		return false
	}
	panic("table expected!")
}

func (self *luaState) Error() int {
	err := self.stack.pop()
	panic(err)
}

func (self *luaState) StringToNumber(s string) bool {
	if n, ok := number.ParseInteger(s); ok {
		self.PushInteger(n)
		return true
	}
	if n, ok := number.ParseFloat(s); ok {
		self.PushNumber(n)
		return true
	}
	return false
}
