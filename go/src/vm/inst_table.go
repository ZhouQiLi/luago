package vm

import (
	. "luago/api"
)

const FIELDS_PER_FLUSH = 50

func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	b++
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

/*
setList命令是批量写入作为数组的table中。入栈的元素会紧跟在table后, A操作数表示table在栈中的索引, B操作数表示有n个要写入table的元素。C操作数为写入table的起始下标。
因为ABC模式中C操作数只有9位, 所以只用512表示table的下标肯定是不足的, 所以对于量大的写入操作会执行批次处理, 分多次将元素写入table。C操作数会进行偏移计算。当C值为0时,
偏移的值会写入下一条Ax指令中, 否则使用C值-1来进行下一步具体偏移值的计算。得到C值后起始的索引计算公式为: C * FIELDS_PER_FLUSH。

举个例子, 要往table中写入1000个元素
会有的指令为：
SETLIST 0 500 1
SETLIST 0 500 11

第一次批量写入的起始下标为 (1 - 1) * 50, 得到结果为0, 所以每个元素写入的下标为 0+1, 0+2 ... 0+n
第二次批量写入的起始下标为 (11 - 1) * 50, 得到的结果为500, 所以每个元素写入的下标为 500+1, 500+2 ... 500+n
*/
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a++
	if c > 0 {
		c--
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	// 当操作数b为0时, 说明要初始化的元素都在栈顶
	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	index := int64(c * FIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		index++
		vm.PushValue(a + j)
		vm.SetI(a, index)
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			index++
			vm.PushValue(j)
			vm.SetI(a, index)
		}

		vm.SetTop(vm.RegisterCount())
	}
}
