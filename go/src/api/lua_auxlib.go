package api

type FuncReg map[string]GoFunction

type AuxLib interface {
	// 报错相关
	Error2(fmt string, a ...interface{}) int
	ArgError(arg int, extraMag string) int
	// 参数检查
	CheckStack2(sz int, msg string)
	ArgCheck(cond bool, arg int, extraMsg string)
	CheckAny(arg int)
	CheckType(arg int, t LuaType)
	CheckInteger(arg int) int64
	CheckNumber(arg int) float64
	CheckString(arg int) string
	OptInteger(arg int, d int64) int64
	OptNumber(arg int, d float64) float64
	OptString(arg int, d string) string
	// 加载相关
	DoFile(filename string) bool
	DoString(str string) bool
	LoadFile(filename string) int
	LoadFileX(filename, mode string) int
	LoadString(s string) int
	// 其他
	TypeName2(index int) string
	ToString2(index int) string
	Len2(index int) int64
	GetSubTable(index int, fname string) bool
	GetMetafield(obj int, e string) LuaType
	CallMeta(obj int, e string) bool
	OpenLibs()
	RequireF(modname string, openf GoFunction, glb bool)
	NewLib(l FuncReg)
	NewLibTable(l FuncReg)
	SetFuncs(l FuncReg, nup int)
}
