package vm

import (
	. "luago/api"
)

func binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

func unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

func add(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPADD)
}

func sub(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPSUB)
}

func mul(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPMUL)
}

func mod(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPMOD)
}

func pow(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPPOW)
}

func div(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPDIV)
}

func idiv(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPIDIV)
}

func band(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPBAND)
}

func bor(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPBOR)
}

func bxor(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPBXOR)
}

func shl(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPSHL)
}

func shr(i Instruction, vm LuaVM) {
	binaryArith(i, vm, LUA_OPSHR)
}

func unm(i Instruction, vm LuaVM) {
	unaryArith(i, vm, LUA_OPUNM)
}

func bnot(i Instruction, vm LuaVM) {
	unaryArith(i, vm, LUA_OPBNOT)
}

func valueLength(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	b++

	vm.Len(b)
	vm.Replace(a)
}

func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	b++
	c++
	n := c - b + 1
	// 要将拼接的值推入栈顶后才能拼接, 因为不确定要拼接的变量数量, 所以要预先分配好足够的栈空间。
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

func compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, c := i.ABC()
	vm.GetRK(b)
	vm.GetRK(c)

	// a操作数记录的是预期结果, 如果Compare的返回值与预期不一致, 则跳过下一条指令。
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}

	// 将刚才入栈的两个值出栈
	vm.Pop(2)
}

func eq(i Instruction, vm LuaVM) {
	compare(i, vm, LUA_OPEQ)
}

func lt(i Instruction, vm LuaVM) {
	compare(i, vm, LUA_OPLT)
}

func le(i Instruction, vm LuaVM) {
	compare(i, vm, LUA_OPLE)
}

func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a++
	b++

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

func testset(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	b++
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a+1)
	} else {
		vm.AddPC(1)
	}
}

func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a++

	// 一般下一条指令是jmp, 所以当ac一致时要跳过jmp指令执行
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a++
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	vm.AddPC(sBx)
}

func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a++
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)

	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) >= 0
	if (isPositiveStep && vm.Compare(a, a+1, LUA_OPLE)) || (!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE)) {
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}
