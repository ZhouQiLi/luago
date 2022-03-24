package binary_chunk

import (
	"fmt"
	"luago/vm"
)

func PrintHeader(protoType *Prototype) {
	funcType := "main"
	if protoType.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if protoType.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n", funcType, protoType.Source, protoType.LineDefined, protoType.LastLineDefined, len(protoType.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ", protoType.NumParams, varargFlag, protoType.MaxStackSize, len(protoType.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n", len(protoType.LocVars), len(protoType.Constants), len(protoType.Protos))
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xff)
			} else {
				fmt.Printf("  %d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xff)
			} else {
				fmt.Printf("  %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf("  %d", bx)
		}
	case vm.IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", ax)
	}
}

func PrintCode(protoType *Prototype) {
	for pc, c := range protoType.Code {
		line := "-"
		if len(protoType.LineInfo) > 0 {
			line = fmt.Sprintf("%d", protoType.LineInfo[pc])
		}

		i := vm.Instruction(c)

		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%s", k)
	default:
		return "?"
	}
}

func upvalueName(protoType *Prototype, index int) string {
	if len(protoType.UpvalueNames) >= index {
		return protoType.UpvalueNames[index]
	}

	return "-"
}

func PrintDetail(protoType *Prototype) {
	fmt.Printf("constants (%d):\n", len(protoType.Constants))
	for i, k := range protoType.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(protoType.LocVars))
	for i, locVar := range protoType.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC, locVar.EndPC)
	}

	fmt.Printf("upvalues (%d):\n", len(protoType.Upvalues))
	for i, upvalue := range protoType.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalueName(protoType, i), upvalue.Instack, upvalue.Idx)
	}

}
