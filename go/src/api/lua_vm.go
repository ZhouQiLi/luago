package api

type LuaVM interface {
	LuaState
	PC() int            // 返回当前PC值
	AddPC(n int)        // 修改PC(指令跳转)
	Fetch() uint32      // 获取当前指令;将PC指向下一条指令
	GetConst(index int) // 将指定常量推入栈顶
	GetRK(rk int)       // 将指定常量或栈值推入栈顶

	RegisterCount() int
	LoadVararg(n int)
	LoadProto(index int)
}
