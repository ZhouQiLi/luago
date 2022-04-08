package state

import (
	"luago/binary_chunk"
)

type closure struct {
	proto *binary_chunk.Prototype
}

func newLuaClosure(proto *binary_chunk.Prototype) *closure {
	return &closure{proto}
}
