package parser

import (
	. "luago/compiler/ast"
	. "luago/compiler/lexer"
	"luago/number"
)

func parseExpList(lexer *Lexer) []Exp {
	exps := make([]Exp, 0, 4)
	exps = append(exps, parseExp(lexer))
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		exps = append(exps, parseExp(lexer))
	}

	return exps
}

// exp   ::= exp12
// exp12 ::= exp11 {or exp11}
// exp11 ::= exp10 {and exp10}
// exp10 ::= exp9 {('<' | '>' | '<=' | '>=' | '~=' | '==') exp9}
// exp9  ::= exp8 {'|' exp8}
// exp8  ::= exp7 {'~' exp7}
// exp7  ::= exp6 {'&' exp6}
// exp6  ::= exp5 {('<<' | '>>') exp5}
// exp5  ::= exp4 {'..' exp4}
// exp4  ::= exp3 {('+' | '-' | '*' | '/' | '//' | '%') exp3}
// exp2  ::= {('not' | '#' | '-' | '~')} exp1
// exp1  ::= exp0 {'^' exp2}
// exp0  ::= nil | false | true | Numeral | LiteralString
//         | '...' | functiondef | prefixexp | tableconstructor
func parseExp(lexer *Lexer) Exp {
	return parseExp12(lexer)
}

func parseExp12(lexer *Lexer) Exp {
	exp := parseExp11(lexer)
	for lexer.LookAhead() == TOKEN_OP_OR {
		line, op, _ := lexer.NextToken()
		exp = optimizeLogicalOr(&BinopExp{line, op, exp, parseExp11(lexer)})
	}

	return exp
}

func parseExp11(lexer *Lexer) Exp {
	exp := parseExp10(lexer)
	for lexer.LookAhead() == TOKEN_OP_AND {
		line, op, _ := lexer.NextToken()
		exp = optimizeLogicalAnd(&BinopExp{line, op, exp, parseExp10(lexer)})
	}

	return exp
}

func parseExp10(lexer *Lexer) Exp {
	exp := parseExp9(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_LT, TOKEN_OP_GT, TOKEN_OP_LE, TOKEN_OP_GE, TOKEN_OP_EQ, TOKEN_OP_NE:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{line, op, exp, parseExp9(lexer)}
		default:
			return exp
		}
	}
}

func parseExp9(lexer *Lexer) Exp {
	exp := parseExp8(lexer)
	for lexer.LookAhead() == TOKEN_OP_BOR {
		line, op, _ := lexer.NextToken()
		exp = optimizeBitwiseBinaryOp(&BinopExp{line, op, exp, parseExp8(lexer)})
	}
	return exp
}

func parseExp8(lexer *Lexer) Exp {
	exp := parseExp7(lexer)
	for lexer.LookAhead() == TOKEN_OP_BXOR {
		line, op, _ := lexer.NextToken()
		exp = optimizeBitwiseBinaryOp(&BinopExp{line, op, exp, parseExp7(lexer)})
	}
	return exp
}

func parseExp7(lexer *Lexer) Exp {
	exp := parseExp6(lexer)
	for lexer.LookAhead() == TOKEN_OP_BAND {
		line, op, _ := lexer.NextToken()
		exp = optimizeBitwiseBinaryOp(&BinopExp{line, op, exp, parseExp6(lexer)})
	}
	return exp
}

func parseExp6(lexer *Lexer) Exp {
	exp := parseExp5(lexer)

	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_SHL, TOKEN_OP_SHR:
			line, op, _ := lexer.NextToken()
			exp = optimizeBitwiseBinaryOp(&BinopExp{line, op, exp, parseExp5(lexer)})
		default:
			return exp
		}
	}
}

func parseExp5(lexer *Lexer) Exp {
	exp := parseExp4(lexer)
	if lexer.LookAhead() != TOKEN_OP_CONCAT {
		return exp
	}

	line := 0
	exps := []Exp{exp}

	for lexer.LookAhead() == TOKEN_OP_CONCAT {
		line, _, _ = lexer.NextToken()
		exps = append(exps, parseExp4(lexer))
	}

	return &ConcatExp{line, exps}
}

func parseExp4(lexer *Lexer) Exp {
	exp := parseExp3(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_ADD, TOKEN_OP_SUB:
			line, op, _ := lexer.NextToken()
			exp = optimizeArithBinaryOp(&BinopExp{line, op, exp, parseExp3(lexer)})
		default:
			return exp
		}
	}
}

func parseExp3(lexer *Lexer) Exp {
	exp := parseExp2(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_MUL, TOKEN_OP_DIV, TOKEN_OP_IDIV, TOKEN_OP_MOD:
			line, op, _ := lexer.NextToken()
			exp = optimizeArithBinaryOp(&BinopExp{line, op, exp, parseExp2(lexer)})
		default:
			return exp
		}
	}
}

func parseExp2(lexer *Lexer) Exp {
	switch lexer.LookAhead() {
	case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
		line, op, _ := lexer.NextToken()
		exp := &UnopExp{line, op, parseExp2(lexer)}
		return optimizeUnaryOp(exp)
	}
	return parseExp1(lexer)
}

func parseExp1(lexer *Lexer) Exp {
	exp := parseExp0(lexer)
	if lexer.LookAhead() == TOKEN_OP_POW {
		line, op, _ := lexer.NextToken()
		exp = optimizePow(&BinopExp{line, op, exp, parseExp2(lexer)})
	}

	return exp
}

func parseExp0(lexer *Lexer) Exp {
	switch lexer.LookAhead() {
	case TOKEN_VARARG:
		line, _, _ := lexer.NextToken()
		return &VarargExp{line}

	case TOKEN_KW_NIL:
		line, _, _ := lexer.NextToken()
		return &NilExp{line}
	case TOKEN_KW_TRUE:
		line, _, _ := lexer.NextToken()
		return &TrueExp{line}
	case TOKEN_KW_FALSE:
		line, _, _ := lexer.NextToken()
		return &FalseExp{line}
	case TOKEN_STRING:
		line, _, token := lexer.NextToken()
		return &StringExp{line, token}
	case TOKEN_NUMBER:
		return parseNumberExp(lexer)
	case TOKEN_SEP_LCURLY:
		return parseTableConstructorExp(lexer)
	case TOKEN_KW_FUNCTION:
		lexer.NextToken()
		return parseFuncDefExp(lexer)
	default:
		return parsePrefixExp(lexer)
	}
}

func parseNumberExp(lexer *Lexer) Exp {
	line, _, token := lexer.NextToken()
	if i, ok := number.ParseInteger(token); ok {
		return &IntegerExp{line, i}
	} else if f, ok := number.ParseFloat(token); ok {
		return &FloatExp{line, f}
	} else {
		panic("not a number: " + token)
	}
}

func parseFuncDefExp(lexer *Lexer) *FuncDefExp {
	line := lexer.Line()
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN)
	parList, isVararg := _parseParList(lexer)
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)
	block := parseBlock(lexer)
	lastLine, _ := lexer.NextTokenOfKind(TOKEN_KW_END)
	return &FuncDefExp{line, lastLine, parList, isVararg, block}
}

func _parseParList(lexer *Lexer) (names []string, isVararg bool) {
	switch lexer.LookAhead() {
	case TOKEN_SEP_RPAREN:
		return nil, false
	case TOKEN_VARARG:
		lexer.NextToken()
		return nil, true
	}
	_, name := lexer.NextIdentifier()
	names = append(names, name)
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		if lexer.LookAhead() == TOKEN_IDENTIFIER {
			_, name := lexer.NextIdentifier()
			names = append(names, name)
		} else {
			lexer.NextTokenOfKind(TOKEN_VARARG)
			isVararg = true
			break
		}
	}

	return
}

func parseTableConstructorExp(lexer *Lexer) *TableConstructorExp {
	line := lexer.Line()
	lexer.NextTokenOfKind(TOKEN_SEP_LCURLY)
	keyExps, valExps := _parseFieldList(lexer)
	lexer.NextTokenOfKind(TOKEN_SEP_RCURLY)
	lastLine := lexer.Line()
	return &TableConstructorExp{line, lastLine, keyExps, valExps}
}

func _parseFieldList(lexer *Lexer) (keyExps, valExps []Exp) {
	if lexer.LookAhead() != TOKEN_SEP_RCURLY {
		k, v := _parseField(lexer)
		keyExps = append(keyExps, k)
		valExps = append(valExps, v)
		for _isFieldSep(lexer.LookAhead()) {
			lexer.NextToken()
			if lexer.LookAhead() != TOKEN_SEP_RCURLY {
				k, v := _parseField(lexer)
				keyExps = append(keyExps, k)
				valExps = append(valExps, v)
			} else {
				break
			}
		}
	}

	return
}

func _isFieldSep(tokenKind int) bool {
	return tokenKind == TOKEN_SEP_COMMA || tokenKind == TOKEN_SEP_SEMI
}

func _parseField(lexer *Lexer) (k, v Exp) {
	if lexer.LookAhead() == TOKEN_SEP_LBRACK {
		lexer.NextToken()
		k = parseExp(lexer)
		lexer.NextTokenOfKind(TOKEN_SEP_RBRACK)
		lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)
		v = parseExp(lexer)
		return
	}

	exp := parseExp(lexer)
	if nameExp, ok := exp.(*NameExp); ok {
		if lexer.LookAhead() == TOKEN_OP_ASSIGN {
			lexer.NextToken()
			k = &StringExp{nameExp.Line, nameExp.Name}
			v = parseExp(lexer)
			return
		}
	}

	return nil, exp
}
