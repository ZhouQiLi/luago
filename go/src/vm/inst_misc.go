package vm

import (
	. "luago/api"
)

func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.Copy(b, a)
}

func jump(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		vm.CloseUpvalues(a)
	}
}
