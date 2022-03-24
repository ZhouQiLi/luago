package main

import (
	"fmt"
	"luago/api"
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

func main() {
	L := state.New()

	L.PushBoolean(true)
	printStack(L)

	L.PushInteger(10)
	printStack(L)

	L.PushNil()
	printStack(L)

	L.PushString("hello")
	printStack(L)

	L.PushValue(-4)
	printStack(L)

	L.Replace(3)
	printStack(L)

	L.SetTop(6)
	printStack(L)

	L.Remove(-3)
	printStack(L)

	L.SetTop(-5)
	printStack(L)

}

func printStack(L api.LuaState) {
	top := L.GetTop()
	for index := 1; index <= top; index++ {
		t := L.Type(index)
		switch t {
		case api.LUA_TBOOLEAN:
			fmt.Printf("[%t]", L.ToBoolean(index))
		case api.LUA_TNUMBER:
			fmt.Printf("[%g]", L.ToNumber(index))
		case api.LUA_TSTRING:
			fmt.Printf("[%q]", L.ToString(index))
		default:
			fmt.Printf("[%s]", L.TypeName(t))
		}
	}
	fmt.Println()
}
