package vm

import (
	. "luago/api"
)

func tForLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a++
	if !vm.IsNil(a + 1) {
		vm.Copy(a+1, a)
		vm.AddPC(sBx)
	}
}
