package state

import (
	"luago/number"
	"math"
)

type luaTable struct {
	metatable *luaTable
	arr       []luaValue
	_map      map[luaValue]luaValue
	keys      map[luaValue]luaValue
	changed   bool
}

func newLuaTable(arrSize, mapSize int) *luaTable {
	t := &luaTable{}
	if arrSize > 0 {
		t.arr = make([]luaValue, 0, arrSize)
	}

	if mapSize > 0 {
		t._map = make(map[luaValue]luaValue, mapSize)
	}

	return t
}

func floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (self *luaTable) get(key luaValue) luaValue {
	key = floatToInteger(key)
	if index, ok := key.(int64); ok {
		if index >= 1 && index <= int64(len(self.arr)) {
			return self.arr[index-1]
		}
	}
	return self._map[key]
}

func (self *luaTable) shrinkArray() {
	for i := len(self.arr) - 1; i >= 0; i-- {
		if self.arr[i] == nil {
			self.arr = self.arr[0:i]
		} else {
			break
		}
	}
}

func (self *luaTable) expandArray() {
	for i := int64(len(self.arr)) + 1; true; i++ {
		if value, found := self._map[i]; found {
			delete(self._map, i)
			self.arr = append(self.arr, value)
		} else {
			break
		}
	}
}

func (self *luaTable) len() int {
	return len(self.arr)
}

func (self *luaTable) put(key, value luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN")
	}
	key = floatToInteger(key)

	if i, ok := key.(int64); ok && i >= 1 {
		arrLen := int64(len(self.arr))
		if i <= arrLen {
			self.arr[i-1] = value
			if i == arrLen && value == nil {
				self.shrinkArray()
			}
			return
		}

		if i == arrLen+1 {
			delete(self._map, key)
			if value != nil {
				self.arr = append(self.arr, value)
				self.expandArray()
			}
			return
		}
	}

	if value != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue, 8)
		}
		self._map[key] = value
	} else {
		delete(self._map, key)
	}
}

func (self *luaTable) hasMetaField(fieldName string) bool {
	return self.metatable != nil && self.metatable.get(fieldName) != nil
}

func (self *luaTable) nextKey(key luaValue) luaValue {
	if self.keys == nil || key == nil {
		self.initKeys()
		self.changed = false
	}

	return self.keys[key]
}

func (self *luaTable) initKeys() {
	self.keys = make(map[luaValue]luaValue)
	var key luaValue = nil
	for i, v := range self.arr {
		if v != nil {
			self.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}

	for k, v := range self._map {
		if v != nil {
			self.keys[key] = k
			key = k
		}
	}
}
