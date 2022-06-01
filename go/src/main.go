package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "luago/api"
	. "luago/compiler/lexer"
	"luago/compiler/parser"
	"luago/state"

	// . "luago/vm"
	"os"
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

// func main() {
// L := state.New()
// L.PushInteger(1)
// L.PushString("2.0")
// L.PushString("3.0")
// L.PushNumber(4.0)
// printStack(L)

// L.Arith(LUA_OPADD)
// printStack(L)

// L.Arith(LUA_OPBNOT)
// printStack(L)

// L.Len(2)
// printStack(L)

// L.Concat(3)
// printStack(L)

// L.PushBoolean(L.Compare(1, 2, LUA_OPEQ))
// printStack(L)
// }

// func main() {
// 	if len(os.Args) > 1 {
// 		data, err := ioutil.ReadFile(os.Args[1])
// 		if err != nil {
// 			panic(err)
// 		}
// 		proto := binchunk.Undump(data)
// 		luaMain(proto)
// 	}
// }

// func luaMain(proto *binchunk.Prototype) {
// 	maxStackSize := int(proto.MaxStackSize)
// 	L := state.New(maxStackSize+8, proto)
// 	L.SetTop(maxStackSize)
// 	for {
// 		pc := L.PC()
// 		inst := Instruction(L.Fetch())
// 		if inst.Opcode() != OP_RETURN {
// 			inst.Execute(L)
// 			fmt.Printf("[%02d] %s", pc+1, inst.OpName())
// 			printStack(L)
// 		} else {
// 			break
// 		}
// 	}
// }

func print(L LuaState) int {
	argsCount := L.GetTop()
	for i := 1; i <= argsCount; i++ {
		if L.IsBoolean(i) {
			fmt.Printf("%t", L.ToBoolean(i))
		} else if L.IsString(i) {
			fmt.Print(L.ToString(i))
		} else {
			fmt.Print(L.TypeName(L.Type(i)))
		}

		if i < argsCount {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func getMetatable(L LuaState) int {
	if !L.GetMetatable(1) {
		L.PushNil()
	}
	return 1
}

func setMetatable(L LuaState) int {
	L.SetMetatable(1)
	return 1
}

func next(L LuaState) int {
	L.SetTop(2)
	if L.Next(1) {
		return 2
	} else {
		L.PushNil()
		return 1
	}
}

func pairs(L LuaState) int {
	L.PushGoFunction(next)
	L.PushValue(1)
	L.PushNil()
	return 3
}

func _ipairsAux(L LuaState) int {
	i := L.ToInteger(2) + 1
	L.PushInteger(i)
	if L.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

func ipairs(L LuaState) int {
	L.PushGoFunction(_ipairsAux)
	L.PushValue(1)
	L.PushInteger(0)
	return 3
}

func _error(L LuaState) int {
	return L.Error()
}

func pCall(L LuaState) int {
	argsCount := L.GetTop() - 1
	status := L.PCall(argsCount, -1, 0)
	L.PushBoolean(status == LUA_OK)
	L.Insert(1)
	return L.GetTop()
}

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		L := state.New()
		L.Register("print", print)
		L.Register("getmetatable", getMetatable)
		L.Register("setmetatable", setMetatable)
		L.Register("next", next)
		L.Register("pairs", pairs)
		L.Register("ipairs", ipairs)
		L.Register("error", _error)
		L.Register("pcall", pCall)
		L.Load(data, os.Args[1], "bt")
		L.Call(0, 0)
	}
}

func testLexer(chunk, chunkName string) {
	lexer := NewLexer(chunk, chunkName)
	for {
		line, kind, token := lexer.NextToken()
		fmt.Printf("[%2d] [%-10s] %s\n", line, kindToCategory(kind), token)
		if kind == TOKEN_EOF {
			break
		}
	}
}

func kindToCategory(kind int) string {
	switch {
	case kind < TOKEN_SEP_SEMI:
		return "other"
	case kind <= TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= TOKEN_OP_NOT:
		return "operator"
	case kind <= TOKEN_KW_WHILE:
		return "keyword"
	case kind == TOKEN_IDENTIFIER:
		return "identifier"
	case kind == TOKEN_NUMBER:
		return "number"
	case kind == TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}

func testParser(chunk, chunkName string) {
	ast := parser.Parse(chunk, chunkName)
	b, err := json.Marshal(ast)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

// func main() {
// 	if len(os.Args) > 1 {
// 		data, err := ioutil.ReadFile(os.Args[1])
// 		if err != nil {
// 			panic(err)
// 		}
// 		testParser(string(data), os.Args[1])
// 	}
// }
