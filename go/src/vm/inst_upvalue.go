package vm

import (
	. "luago/api"
)

func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	b++

	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}


func getUpvalue(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	b++

	vm.Copy(LuaUpvalueIndex(b), a)
}

func setUpvalue(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.Copy(a, LuaUpvalueIndex(b))
}
