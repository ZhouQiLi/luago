package vm

import (
	. "luago/api"
)

// 在指令执行时会预先设定好栈空间, 所以本函数不需要额外对栈空间做操作。
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	// nil值入栈, 并将指定栈值使用copy函数拷贝为nil值, 循环后将栈顶的nil值出栈。
	vm.PushNil()
	for i := a + 1; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++

	// 目标值在寄存器中, 所以要先入栈, 再通过Replace接口赋值给指定位置
	vm.PushBoolean(b != 0)
	vm.Replace(a)

	if c != 0 {
		vm.AddPC(1)
	}
}

// Bx指令只有18个比特, 所以最大的常量数量为262143个
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a++

	vm.GetConst(bx)
	vm.Replace(a)
}

// 先用ABx指令获取目标位置, 再配合EXTRAARG(iAx)指令获取指定常量的索引
// Ax操作数有26个比特, 可表示的最大常量数量为67108864个, 可以满足大部分情况了
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a++
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
