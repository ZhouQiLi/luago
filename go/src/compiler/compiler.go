package compiler

import (
	binchunk "luago/binary_chunk"
	"luago/compiler/codegen"
	"luago/compiler/parser"
)

func Compile(chunk, chunkName string) *binchunk.Prototype {
	ast := parser.Parse(chunk, chunkName)
	return codegen.GenProto(ast)
}
