package vm

import (
	. "luago/api"
)

func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a++

	vm.LoadProto(bx)
	vm.Replace(a)
}

func fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}

	// 原本的数据并没有出栈而是在栈顶, 所以使用旋转将数据放置到正确的位置上。
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func pushFunctionAndArgs(a, b int, vm LuaVM) int {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func popResults(a, c int, vm LuaVM) {
	if c == 1 {

	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++

	argsCount := pushFunctionAndArgs(a, b, vm)
	vm.Call(argsCount, c-1)
	popResults(a, c, vm)
}

func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	if b == 1 {

	} else if b > 1 {
		for index := a; index <= a+b-2; index++ {
			vm.PushValue(index)
		}
	} else {
		fixStack(a, vm)
	}
}

func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	if b != 1 {
		vm.LoadVararg(b - 1)
		popResults(a, b, vm)
	}
}

func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	c := 0
	argsCount := pushFunctionAndArgs(a, b, vm)
	vm.Call(argsCount, c-1)
	popResults(a, c, vm)
}

// 正常操作先获取function再添加table作为参数之一要两条move指令, self只用一条。
func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	b++

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
