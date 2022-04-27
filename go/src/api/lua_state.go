package api

type LuaType = int
type ArithOp = int
type CompareOp = int // 比较操作符
type GoFunction func(LuaState) int

func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRY_INDEX - i
}

type LuaState interface {
	// 堆栈操作
	GetTop() int
	AbsIndex(index int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIndex, toIndex int)
	PushValue(index int)
	Replace(index int)
	Insert(index int)
	Remove(index int)
	Rotate(index, n int)
	SetTop(index int)

	// 出栈操作 stack -> go
	TypeName(tp LuaType) string
	Type(index int) LuaType
	IsNone(index int) bool
	IsNil(index int) bool
	IsNoneOrNil(index int) bool
	IsBoolean(index int) bool
	IsInteger(index int) bool
	IsNumber(index int) bool
	IsString(index int) bool
	ToBoolean(index int) bool
	ToInteger(index int) int64
	ToIntegerX(index int) (int64, bool)
	ToNumber(index int) float64
	ToNumberX(index int) (float64, bool)
	ToString(index int) string
	ToStringX(index int) (string, bool)

	// 入栈操作 go -> stack
	PushNil()
	PushBoolean(b bool)
	PushInteger(i int64)
	PushNumber(n float64)
	PushString(s string)

	// 操作符
	Arith(op ArithOp)
	Compare(index1, index2 int, op CompareOp) bool
	Len(index int)
	Concat(n int)

	// 表操作
	NewTable()
	CreateTable(arrSize, mapSize int)
	GetTable(index int) LuaType
	GetField(index int, key string) LuaType
	GetI(index int, i int64) LuaType
	SetTable(index int)
	SetField(index int, k string)
	SetI(index int, i int64)

	// 函数调用
	Load(chunk []byte, chunkName, mode string) int
	Call(argsCount, results int)

	// golang函数调用
	PushGoFunction(f GoFunction)
	IsGoFunction(index int) bool
	ToGoFunction(index int) GoFunction

	// 全局环境
	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)

	PushGoClosure(f GoFunction, n int)

	// 元表
	GetMetatable(index int) bool
	SetMetatable(index int)
	RawLen(index int) uint
	RawEqual(index1, index2 int) bool
	RawGet(index int) LuaType
	RawSet(index int)
	RawGetI(index int, i int64) LuaType
	RawSetI(index int, i int64)

	// 迭代器
	Next(index int) bool
}
