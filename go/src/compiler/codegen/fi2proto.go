package codegen

import (
	. "luago/binary_chunk"
)

func toProto(fi *funcInfo) *Prototype {
	proto := &Prototype{
		NumParams:    byte(fi.numParams),
		MaxStackSize: byte(fi.maxRegs),
		Code:         fi.insts,
		Constants:    getConstants(fi),
		Upvalues:     getUpvalues(fi),
		Protos:       toProtos(fi.subFuncs),
		LineInfo:     []uint32{},
		LocVars:      []LocVar{},
		UpvalueNames: []string{},
	}
	if fi.isVararg {
		proto.IsVararg = 1
	}
	return proto
}

func toProtos(fis []*funcInfo) []*Prototype {
	protos := make([]*Prototype, len(fis))
	for i, fi := range fis {
		protos[i] = toProto(fi)
	}
	return protos
}

func getConstants(fi *funcInfo) []interface{} {
	constants := make([]interface{}, len(fi.constants))
	for k, index := range fi.constants {
		constants[index] = k
	}
	return constants
}

func getUpvalues(fi *funcInfo) []Upvalue {
	upvalues := make([]Upvalue, len(fi.upvalues))
	for _, uv := range fi.upvalues {
		if uv.locVarSlot >= 0 {
			upvalues[uv.index] = Upvalue{1, byte(uv.locVarSlot)}
		} else {
			upvalues[uv.index] = Upvalue{0, byte(uv.upvalIndex)}
		}
	}
	return upvalues
}
