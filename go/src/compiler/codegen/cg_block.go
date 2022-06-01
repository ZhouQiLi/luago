package codegen

import (
	. "luago/compiler/ast"
)

func cgBlock(fi *funcInfo, node *Block) {
	for _, stat := range node.Stats {
		cgStat(fi, stat)
	}

	if node.RetExps != nil {
		cgRetStat(fi, node.RetExps)
	}
}

func cgRetStat(fi *funcInfo, exps []Exp) {
	expCount := len(exps)
	if expCount == 0 {
		fi.emitReturn(0, 0)
		return
	}

	multRet := isVarargOrFuncCall(exps[expCount-1])
	for i, exp := range exps {
		r := fi.allocReg()
		if i == expCount-1 && multRet {
			cgExp(fi, exp, r, -1)
		} else {
			cgExp(fi, exp, r, 1)
		}
	}

	fi.freeRegs(expCount)

	a := fi.usedRegs
	if multRet {
		fi.emitReturn(a, -1)
	} else {
		fi.emitReturn(a, expCount)
	}
}

func isVarargOrFuncCall(exp Exp) bool {
	switch exp.(type) {
	case *VarargExp, *FuncCallExp:
		return true
	}
	return false
}
