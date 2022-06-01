package state

import (
	. "luago/api"
	binchunk "luago/binary_chunk"
)

type upvalue struct {
	value *luaValue
}

type closure struct {
	proto    *binchunk.Prototype
	goFunc   GoFunction
	upvalues []*upvalue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if upvalueCount := len(proto.Upvalues); upvalueCount > 0 {
		c.upvalues = make([]*upvalue, upvalueCount)
	}
	return c
}

func newGoClosure(f GoFunction, upvalueCount int) *closure {
	c := &closure{goFunc: f}
	if upvalueCount > 0 {
		c.upvalues = make([]*upvalue, upvalueCount)
	}

	return c
}
