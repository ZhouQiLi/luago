package api

type LuaType = int

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
}