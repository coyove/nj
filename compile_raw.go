package potatolang

import (
	"fmt"

	"github.com/coyove/potatolang/parser"
)

var _rawOP0 = map[string]byte{
	"assert": OP_ASSERT,
	"store":  OP_STORE,
	"load":   OP_LOAD,
	"add":    OP_ADD,
	"sub":    OP_SUB,
	"mul":    OP_MUL,
	"div":    OP_DIV,
	"mod":    OP_MOD,
	"not":    OP_NOT,
	"eq":     OP_EQ,
	"neq":    OP_NEQ,
	"less":   OP_LESS,
	"lesseq": OP_LESS_EQ,
	"bnot":   OP_BIT_NOT,
	"band":   OP_BIT_AND,
	"bor":    OP_BIT_OR,
	"bxor":   OP_BIT_XOR,
	"blsh":   OP_BIT_LSH,
	"brsh":   OP_BIT_RSH,
	"pop":    OP_POP,
	"slice":  OP_SLICE,
	"len":    OP_LEN,
	"typeof": OP_TYPEOF,
	"nop":    OP_NOP,
	"eob":    OP_EOB,
	"r0r2":   OP_R0R2,
	"r1r2":   OP_R1R2,
	"copy":   OP_COPY,
}

func (table *symtable) compileRawOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	defer func() {
		if err == nil {
			code.WritePos(atoms[0].Meta)
		}
	}()
	// $op(a, b, c)
	head := atoms[1]
	opname := head.Value.(string)[1:]

	if o, ok := _rawOP0[opname]; ok {
		code.WriteOP(o, 0, 0)
		return
	}

	for _, arg := range atoms[2].C() {
		if arg.Type == parser.Ncompound {
			err = fmt.Errorf("%v: can't use complex value in raw opcode", arg)
			return
		}
	}

	extract := func(idx int) (offset int32, addr uint32, k uint16, isK bool, ok bool) {
		arg := atoms[2].Cx(idx)
		switch arg.Type {
		case parser.Nnumber:
			offset = int32(arg.Value.(float64))
			k = table.addConst(arg.Value)
			isK = true
			ok = true
		case parser.Natom:
			str := arg.Value.(string)
			if str == "$a" {
				addr = regA
				ok = true
			} else {
				addr, ok = table.get(str)
			}
		case parser.Nstring:
			k = table.addConst(arg.Value)
			isK = true
			ok = true
		}
		return
	}

	if atoms[2].Cn() == 0 {
		err = fmt.Errorf("%v: not enough arguments", atoms[0])
		return
	}

	o, a, k, isk, ok := extract(0)
	if !ok {
		err = fmt.Errorf("%v: invalid raw argument", atoms[2].Cx(0))
	}
	// need 1 argument
	switch opname {
	case "jmp":
		code.WriteOP(OP_JMP, 0, uint32(o))
		return
	case "call":
		code.WriteOP(OP_CALL, a, 0)
		return
	case "makemap":
		code.WriteOP(OP_MAKEMAP, uint32(o), 0)
		return
	case "r0":
		if !isk {
			code.WriteOP(OP_R0, a, 0)
		} else {
			code.WriteOP(OP_R0K, uint32(k), 0)
		}
		return
	case "r1":
		if !isk {
			code.WriteOP(OP_R1, a, 0)
		} else {
			code.WriteOP(OP_R1K, uint32(k), 0)
		}
		return
	case "r2":
		if !isk {
			code.WriteOP(OP_R2, a, 0)
		} else {
			code.WriteOP(OP_R2K, uint32(k), 0)
		}
		return
	case "r3":
		if !isk {
			code.WriteOP(OP_R3, a, 0)
		} else {
			code.WriteOP(OP_R3K, uint32(k), 0)
		}
		return
	case "push":
		if !isk {
			code.WriteOP(OP_PUSH, a, 0)
		} else {
			code.WriteOP(OP_PUSHK, uint32(k), 0)
		}
		return
	case "ret":
		if !isk {
			code.WriteOP(OP_RET, a, 0)
		} else {
			code.WriteOP(OP_RETK, uint32(k), 0)
		}
		return
	case "yield":
		if !isk {
			code.WriteOP(OP_YIELD, a, 0)
		} else {
			code.WriteOP(OP_YIELDK, uint32(k), 0)
		}
		return
	}

	if atoms[2].Cn() < 2 {
		err = fmt.Errorf("%v: not enough arguments", head)
		return
	}

	o2, a2, k2, isk2, ok := extract(1)
	if !ok {
		err = fmt.Errorf("%v: invalid raw argument", atoms[2].Cx(1))
	}
	switch opname {
	case "if":
		code.WriteOP(OP_IF, a, uint32(o2))
		return
	case "ifnot":
		code.WriteOP(OP_IFNOT, a, uint32(o2))
		return
	case "set":
		if !isk2 {
			code.WriteOP(OP_SET, a, a2)
		} else {
			code.WriteOP(OP_SETK, a, uint32(k2))
		}
		return
	}

	err = fmt.Errorf("%v: unknow raw operations", head)
	return
}
