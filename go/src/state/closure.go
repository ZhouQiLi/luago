package state

import (
	. "luago/api"
	"luago/binary_chunk"
)

type closure struct {
	proto  *binary_chunk.Prototype
	goFunc GoFunction
}

func newLuaClosure(proto *binary_chunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f GoFunction) *closure {
	return &closure{goFunc: f}
}
