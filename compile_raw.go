package potatolang

import (
	"github.com/coyove/potatolang/parser"
)

var _rawOP0 = map[string]byte{
	"assert": OpAssert,
	"store":  OpStore,
	"load":   OpLoad,
	"add":    OpAdd,
	"sub":    OpSub,
	"mul":    OpMul,
	"div":    OpDiv,
	"mod":    OpMod,
	"not":    OpNot,
	"eq":     OpEq,
	"neq":    OpNeq,
	"less":   OpLess,
	"lesseq": OpLessEq,
	"bnot":   OpBitNot,
	"band":   OpBitAnd,
	"bor":    OpBitOr,
	"bxor":   OpBitXor,
	"blsh":   OpBitLsh,
	"brsh":   OpBitRsh,
	"bursh":  OpBitURsh,
	"pop":    OpPop,
	"slice":  OpSlice,
	"len":    OpLen,
	"typeof": OpTypeof,
	"nop":    OpNOP,
	"eob":    OpEOB,
	"copy":   OpForeach,
}

func (table *symtable) compileRawOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	panic("TODO")
	//	defer func() {
	//		if err == nil {
	//			code.WritePos(atoms[0].Meta)
	//		}
	//	}()
	//	// $op(a, b, c)
	//	head := atoms[1]
	//	opname := head.Value.(string)[1:]
	//
	//	if o, ok := _rawOP0[opname]; ok {
	//		code.WriteOP(o, 0, 0)
	//		return
	//	}
	//
	//	for i := 0; i < 4; i++ {
	//		x := strconv.Itoa(i)
	//		if o, ok := _rawOP0[strings.TrimSuffix(opname, x)]; ok {
	//			code.WriteOP(o, uint16(i+1), 0)
	//			return
	//		}
	//	}
	//
	//	for _, arg := range atoms[2].C() {
	//		if arg.Type == parser.Ncompound {
	//			err = fmt.Errorf("%v: can't use complex value in raw opcode", arg)
	//			return
	//		}
	//	}
	//
	//	extract := func(idx int) (offset int, addr uint16, k uint16, isK bool, ok bool) {
	//		arg := atoms[2].Cx(idx)
	//		switch arg.Type {
	//		case parser.Nnumber:
	//			offset = int(arg.Value.(float64))
	//			k = table.addConst(arg.Value)
	//			isK = true
	//			ok = true
	//		case parser.Natom:
	//			str := arg.Value.(string)
	//			if str == "$a" {
	//				addr = regA
	//				ok = true
	//			} else {
	//				addr, ok = table.get(str)
	//			}
	//		case parser.Nstring:
	//			k = table.addConst(arg.Value)
	//			isK = true
	//			ok = true
	//		}
	//		return
	//	}
	//
	//	if atoms[2].Cn() == 0 {
	//		err = fmt.Errorf("%v: not enough arguments", atoms[0])
	//		return
	//	}
	//
	//	o, a, k, isk, ok := extract(0)
	//	if !ok {
	//		err = fmt.Errorf("%v: invalid raw argument", atoms[2].Cx(0))
	//	}
	//	// need 1 argument
	//	switch opname {
	//	case "jmp":
	//		code.WriteOP(OP_JMP, 0, uint16(o))
	//		return
	//	case "call":
	//		code.WriteOP(OP_CALL, a, 0)
	//		return
	//	case "call0":
	//		code.WriteOP(OP_CALL, a, 1)
	//		return
	//	case "call1":
	//		code.WriteOP(OP_CALL, a, 2)
	//		return
	//	case "call2":
	//		code.WriteOP(OP_CALL, a, 3)
	//		return
	//	case "call3":
	//		code.WriteOP(OP_CALL, a, 4)
	//		return
	//	case "makemap":
	//		code.WriteOP(OP_MAKEMAP, uint16(o), 0)
	//		return
	//	case "r0":
	//		if !isk {
	//			code.WriteOP(OP_R0, a, 0)
	//		} else {
	//			code.WriteOP(OP_R0K, uint16(k), 0)
	//		}
	//		return
	//	case "r1":
	//		if !isk {
	//			code.WriteOP(OP_R1, a, 0)
	//		} else {
	//			code.WriteOP(OP_R1K, uint16(k), 0)
	//		}
	//		return
	//	case "r2":
	//		if !isk {
	//			code.WriteOP(OP_R2, a, 0)
	//		} else {
	//			code.WriteOP(OP_R2K, uint16(k), 0)
	//		}
	//		return
	//	case "r3":
	//		if !isk {
	//			code.WriteOP(OP_R3, a, 0)
	//		} else {
	//			code.WriteOP(OP_R3K, uint16(k), 0)
	//		}
	//		return
	//	case "push":
	//		if !isk {
	//			code.WriteOP(OP_PUSH, a, 0)
	//		} else {
	//			code.WriteOP(OP_PUSHK, uint16(k), 0)
	//		}
	//		return
	//	case "ret":
	//		if !isk {
	//			code.WriteOP(OP_RET, a, 0)
	//		} else {
	//			code.WriteOP(OP_RETK, uint16(k), 0)
	//		}
	//		return
	//	case "yield":
	//		if !isk {
	//			code.WriteOP(OP_YIELD, a, 0)
	//		} else {
	//			code.WriteOP(OP_YIELDK, uint16(k), 0)
	//		}
	//		return
	//	}
	//
	//	if atoms[2].Cn() < 2 {
	//		err = fmt.Errorf("%v: not enough arguments", head)
	//		return
	//	}
	//
	//	o2, a2, k2, isk2, ok := extract(1)
	//	if !ok {
	//		err = fmt.Errorf("%v: invalid raw argument", atoms[2].Cx(1))
	//	}
	//	switch opname {
	//	case "if":
	//		code.WriteOP(OP_IF, a, uint16(o2))
	//		return
	//	case "ifnot":
	//		code.WriteOP(OP_IFNOT, a, uint16(o2))
	//		return
	//	case "set":
	//		if !isk2 {
	//			code.WriteOP(OP_SET, a, a2)
	//		} else {
	//			code.WriteOP(OP_SETK, a, uint16(k2))
	//		}
	//		return
	//	case "rx":
	//		code.WriteOP(OP_RX, uint16(o), uint16(o2))
	//		return
	//	}
	//
	//	err = fmt.Errorf("%v: unknow raw operations", head)
	//	return
}
