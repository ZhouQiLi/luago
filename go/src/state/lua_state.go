package state

import "luago/binary_chunk"

type luaState struct {
	stack *luaStack
	proto *binary_chunk.Prototype
	pc    int
}

func New(stackSize int, proto *binary_chunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
