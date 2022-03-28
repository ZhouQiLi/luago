package main

import (
	"fmt"
	. "luago/api"
	"luago/state"
)

// func main() {
// 	if len(os.Args) > 1 {
// 		data, err := ioutil.ReadFile(os.Args[1])
// 		if err != nil {
// 			panic(err)
// 		}
// 		proto := binchunk.Undump(data)
// 		list(proto)
// 	}
// }

// func list(f *binchunk.Prototype) {
// 	binchunk.PrintHeader(f)
// 	binchunk.PrintCode(f)
// 	binchunk.PrintDetail(f)
// 	for _, p := range f.Protos {
// 		list(p)
// 	}
// }

// func main() {
// L := state.New()

// L.PushBoolean(true)
// printStack(L)

// L.PushInteger(10)
// printStack(L)

// L.PushNil()
// printStack(L)

// L.PushString("hello")
// printStack(L)

// L.PushValue(-4)
// printStack(L)

// L.Replace(3)
// printStack(L)

// L.SetTop(6)
// printStack(L)

// L.Remove(-3)
// printStack(L)

// L.SetTop(-5)
// printStack(L)

// }

func printStack(L LuaState) {
	top := L.GetTop()
	for index := 1; index <= top; index++ {
		t := L.Type(index)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", L.ToBoolean(index))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", L.ToNumber(index))
		case LUA_TSTRING:
			fmt.Printf("[%q]", L.ToString(index))
		default:
			fmt.Printf("[%s]", L.TypeName(t))
		}
	}
	fmt.Println()
}

func main() {
	L := state.New()
	L.PushInteger(1)
	L.PushString("2.0")
	L.PushString("3.0")
	L.PushNumber(4.0)
	printStack(L)

	L.Arith(LUA_OPADD)
	printStack(L)

	L.Arith(LUA_OPBNOT)
	printStack(L)

	L.Len(2)
	printStack(L)

	L.Concat(3)
	printStack(L)

	L.PushBoolean(L.Compare(1, 2, LUA_OPEQ))
	printStack(L)
}
