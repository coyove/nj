package potatolang

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"sync"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

const (
	ERR_UNDECLARED_VARIABLE = "undeclared variable: %+v"
)

func compileSetOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	aVar, aValue := atoms[1], atoms[2]
	varIndex := uint32(0)
	buf := NewBytesWriter()
	var newYX uint32
	var ok bool

	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aValue.Type != parser.NTCompound {
			newYX = uint32(sp)
			sp++
		}
	} else {
		varIndex, ok = table.get(aVar.Value.(string))
		if !ok {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aVar)
			return
		}
		newYX = varIndex
	}

	switch aValue.Type {
	case parser.NTAtom:
		valueIndex, ok := table.get(aValue.Value.(string))
		if !ok {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aValue)
			return
		}
		buf.Write16(OP_SET)
		buf.Write32(newYX)
		buf.Write32(valueIndex)
	case parser.NTNumber, parser.NTString:
		buf.Write16(OP_SETK)
		buf.Write32(newYX)
		buf.Write16(table.addConst(aValue.Value))
	case parser.NTCompound:
		code, newYX, sp, err = compileCompoundIntoVariable(sp, aValue, table,
			atoms[0].Value.(string) == "set", varIndex)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		_, redecl := table.get(aVar.Value.(string))
		if redecl {
			err = fmt.Errorf("redeclare: %+v", aVar)
			return
		}
		table.put(aVar.Value.(string), uint16(newYX))
	}
	return buf.data, newYX, sp, nil
}

func compileRetOp(op, opk uint16) compileFunc {
	return func(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
		buf := NewBytesWriter()
		if len(atoms) == 1 {
			buf.Write16(op)
			buf.Write32(regA)
			return buf.data, yx, sp, nil
		}

		switch atom := atoms[1]; atom.Type {
		case parser.NTAtom, parser.NTNumber, parser.NTString, parser.NTAddr:
			err = fill(buf, atom, table, op, opk)
			if err != nil {
				return
			}
		case parser.NTCompound:
			code, yx, sp, err = extract(sp, atom, table)
			buf.Write(code)
			buf.Write16(op)
			buf.Write32(yx)
		}

		if op == OP_YIELD {
			table.y = true
		}
		return buf.data, yx, sp, nil
	}
}

func compileListOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	var buf *BytesWriter
	buf, sp, err = flaten(sp, atoms[1].Compound, table)
	if err != nil {
		return
	}

	for _, atom := range atoms[1].Compound {
		err = fill(buf, atom, table, OP_PUSH, OP_PUSHK)
		if err != nil {
			return
		}
	}

	buf.Write16(OP_LIST)
	return buf.data, regA, sp, nil
}

func compileMapOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	if len(atoms[1].Compound)%2 != 0 {
		err = fmt.Errorf("every key in map must have a value: %+v", atoms[1])
		return
	}

	var buf *BytesWriter
	buf, sp, err = flaten(sp, atoms[1].Compound, table)
	if err != nil {
		return
	}

	for _, atom := range atoms[1].Compound {
		err = fill(buf, atom, table, OP_PUSH, OP_PUSHK)
		if err != nil {
			return
		}
	}

	buf.Write16(OP_MAP)
	return buf.data, regA, sp, nil
}

func indexToOpR(index int) (a, b uint16) {
	switch index {
	case 0:
		return OP_R0, OP_R0K
	case 1:
		return OP_R1, OP_R1K
	case 2:
		return OP_R2, OP_R2K
	case 3:
		return OP_R3, OP_R3K
	}
	panic("shouldn't happen")
}

func flaten(sp uint16, atoms []*parser.Node, table *symtable) (buf *BytesWriter, newsp uint16, err error) {
	replacedAtoms := []*parser.Node{}
	buf = NewBytesWriter()

	for i, atom := range atoms {
		var yx uint32
		var code []uint16

		if atom.Type == parser.NTCompound {
			code, yx, sp, err = compileCompoundIntoVariable(sp, atom, table, true, 0)
			if err != nil {
				return
			}
			if table.im != nil {
				atoms[i] = &parser.Node{Type: parser.NTNumber, Value: *table.im}
				table.im = nil
			} else if table.ims != nil {
				atoms[i] = &parser.Node{Type: parser.NTString, Value: *table.ims}
				table.ims = nil
			} else {
				atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
				replacedAtoms = append(replacedAtoms, atoms[i])
				buf.Write(code)
			}
		}
	}

	if len(replacedAtoms) == 1 {
		cursor := uint32(buf.Len() - 2)
		replacedAtoms[0].Value = crRead32(buf.data, &cursor)
		buf.TruncateLast(5)
		sp--
	}

	return buf, sp, nil
}

func flatWrite(sp uint16, atoms []*parser.Node, table *symtable, bop uint16) (code []uint16, yx uint32, newsp uint16, err error) {

	var buf *BytesWriter
	buf, sp, err = flaten(sp, atoms[1:], table)
	if err != nil {
		return
	}

	immediateRet := func(n float64) {
		buf.Clear()
		buf.Write16(OP_SETK)
		buf.Write32(regA)
		buf.Write16(table.addConst(n))
		table.im = &n
	}

	v1 := func() float64 { return atoms[1].Value.(float64) }
	v2 := func() float64 { return atoms[2].Value.(float64) }

	switch bop {
	case OP_ADD:
		if atoms[1].Type == parser.NTString && atoms[2].Type == parser.NTString {
			buf.Clear()
			buf.Write16(OP_SETK)
			buf.Write32(regA)
			str := atoms[1].Value.(string) + atoms[2].Value.(string)
			buf.Write16(table.addConst(str))
			table.ims = &str
			return buf.data, regA, sp, nil
		}
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() + v2())
			return buf.data, regA, sp, nil
		}
	case OP_SUB:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() - v2())
			return buf.data, regA, sp, nil
		}
	case OP_MUL:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() * v2())
			return buf.data, regA, sp, nil
		}
	case OP_DIV:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() / v2())
			return buf.data, regA, sp, nil
		}
	}

	for i := 1; i < len(atoms); i++ {
		a, b := indexToOpR(i - 1)
		err = fill(buf, atoms[i], table, a, b)
		if err != nil {
			return
		}
	}
	buf.Write16(bop)

	if bop == OP_ASSERT {
		buf.WriteString(atoms[0].String())
	}

	return buf.data, regA, sp, nil
}

func compileFlatOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {

	head := atoms[0]
	switch head.Value.(string) {
	case "+":
		return flatWrite(sp, atoms, table, OP_ADD)
	case "-":
		return flatWrite(sp, atoms, table, OP_SUB)
	case "*":
		return flatWrite(sp, atoms, table, OP_MUL)
	case "/":
		return flatWrite(sp, atoms, table, OP_DIV)
	case "%":
		return flatWrite(sp, atoms, table, OP_MOD)
	case "<":
		return flatWrite(sp, atoms, table, OP_LESS)
	case "<=":
		return flatWrite(sp, atoms, table, OP_LESS_EQ)
	case "eq":
		return flatWrite(sp, atoms, table, OP_EQ)
	case "neq":
		return flatWrite(sp, atoms, table, OP_NEQ)
	case "assert":
		return flatWrite(sp, atoms, table, OP_ASSERT)
	case "not":
		return flatWrite(sp, atoms, table, OP_NOT)
	case "~":
		return flatWrite(sp, atoms, table, OP_BIT_NOT)
	case "&":
		return flatWrite(sp, atoms, table, OP_BIT_AND)
	case "|":
		return flatWrite(sp, atoms, table, OP_BIT_OR)
	case "^":
		return flatWrite(sp, atoms, table, OP_BIT_XOR)
	case "<<":
		return flatWrite(sp, atoms, table, OP_BIT_LSH)
	case ">>":
		return flatWrite(sp, atoms, table, OP_BIT_RSH)
	case "nil":
		return flatWrite(sp, atoms, table, OP_NIL)
	case "true":
		return flatWrite(sp, atoms, table, OP_TRUE)
	case "false":
		return flatWrite(sp, atoms, table, OP_FALSE)
	case "store":
		return flatWrite(sp, atoms, table, OP_STORE)
	case "load":
		return flatWrite(sp, atoms, table, OP_LOAD)
	}

	err = fmt.Errorf("invalid flat op %+v", head)
	return
}

func compileIncOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	subject, ok := table.get(atoms[1].Value.(string))
	buf := NewBytesWriter()
	if !ok {
		return nil, 0, 0, fmt.Errorf(ERR_UNDECLARED_VARIABLE, atoms[1])
	}
	buf.Write16(OP_INC)
	buf.Write32(subject)
	buf.Write16(table.addConst(atoms[2].Value))
	return buf.data, regA, sp, nil
}

func compileImmediateIntoVariable(sp uint16, node *parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	buf := NewBytesWriter()
	switch node.Type {
	case parser.NTNumber, parser.NTString:
		buf.Write16(OP_SETK)
		buf.Write32(uint32(sp))
		buf.Write16(table.addConst(node.Value))
		yx = uint32(sp)
		sp++
	default:
		panic("shouldn't happen")
	}
	return buf.data, yx, sp, nil
}

// [and a b] => if not a then false else use b end
// [or a b]  => if a then use a else use b end
func compileAndOrOp(bop uint16) compileFunc {
	return func(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
		a1, a2 := atoms[1], atoms[2]
		buf := NewBytesWriter()
		switch a1.Type {
		case parser.NTNumber, parser.NTString:
			code, yx, sp, err = compileImmediateIntoVariable(sp, a1, table)
			buf.Write(code)
			buf.Write16(bop)
			buf.Write32(yx)
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a1, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			buf.Write16(bop)
			buf.Write32(yx)
		}
		buf.Write32(0)
		c2 := buf.Len()

		switch a2.Type {
		case parser.NTNumber, parser.NTString:
			buf.Write16(OP_SETK)
			buf.Write32(regA)
			buf.Write16(table.addConst(a2.Value))
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a2, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			if yx != regA {
				buf.Write16(OP_SET)
				buf.Write32(regA)
				buf.Write32(yx)
			}
		}
		jmp := buf.Len() - c2

		buf.data[c2-2] = uint16(uint32(jmp) >> 16)
		buf.data[c2-1] = uint16(jmp)
		return buf.data, regA, sp, nil
	}
}

// [if condition [truechain ...] [falsechain ...]]
func compileIfOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := NewBytesWriter()

	switch condition.Type {
	case parser.NTNumber, parser.NTString:
		code, yx, sp, err = compileImmediateIntoVariable(sp, condition, table)
		buf.Write(code)
		buf.Write16(OP_IFNOT)
		buf.Write32(yx)
	case parser.NTAtom, parser.NTCompound:
		code, yx, sp, err = extract(sp, condition, table)
		if err != nil {
			return nil, 0, 0, err
		}

		buf.Write(code)
		buf.Write16(OP_IFNOT)
		buf.Write32(yx)

		var trueCode, falseCode []uint16
		trueCode, yx, sp, err = compileChainOp(sp, trueBranch, table)
		if err != nil {
			return
		}
		falseCode, yx, sp, err = compileChainOp(sp, falseBranch, table)
		if err != nil {
			return
		}
		if len(falseCode) > 0 {
			buf.Write32(uint32(len(trueCode)) + 3) // jmp (1) + offset (2)
			buf.Write(trueCode)
			buf.Write16(OP_JMP)
			buf.Write32(uint32(len(falseCode)))
			buf.Write(falseCode)
		} else {
			buf.Write32(uint32(len(trueCode)))
			buf.Write(trueCode)
		}
	}
	return buf.data, regA, sp, nil
}

// [call func-name [args ...]]
func compileCallOp(sp uint16, nodes []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	buf := NewBytesWriter()
	callee := nodes[1]
	name, _ := callee.Value.(string)
	switch name {
	case "error":
		table.e = true
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_ERROR)
	case "len":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_LEN)
	case "dup":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_DUP)
	case "typeof":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_TYPEOF)
	case "who":
		return []uint16{OP_WHO}, regA, sp, nil
	}

	atoms, replacedAtoms := nodes[2].Compound, []*parser.Node{}
	for i := 0; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.NTCompound {
			code, yx, sp, err = compileCompoundIntoVariable(sp, atom, table, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
			replacedAtoms = append(replacedAtoms, atoms[i])
			buf.Write(code)
		}
	}

	// note: [call [..] [..]] is different
	if len(replacedAtoms) == 1 && callee.Type != parser.NTCompound {
		cursor := uint32(buf.Len() - 2)
		replacedAtoms[0].Value = crRead32(buf.data, &cursor)
		buf.TruncateLast(5)
		sp--
	}

	var varIndex uint32
	var ok bool
	switch callee.Type {
	case parser.NTAtom:
		varIndex, ok = table.get(callee.Value.(string))
		if !ok {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, callee)
			return
		}
	case parser.NTCompound:
		code, yx, sp, err = compileCompoundIntoVariable(sp, callee, table, true, 0)
		if err != nil {
			return
		}
		varIndex = yx
		buf.Write(code)
	case parser.NTAddr:
		varIndex = callee.Value.(uint32)
	default:
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = fill(buf, atoms[i], table, OP_PUSH, OP_PUSHK)
		if err != nil {
			return
		}
	}

	buf.Write16(OP_CALL)
	buf.Write32(varIndex)
	return buf.data, regA, sp, nil
}

// [lambda [namelist] [chain ...]]
func compileLambdaOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	newLookup := newsymtable()
	newLookup.parent = table
	newLookup.lineInfo = table.lineInfo

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("invalid lambda parameters: %+v", atoms[0])
		return
	}

	for i, p := range params.Compound {
		newLookup.put(p.Value.(string), uint16(i))
	}

	ln := len(newLookup.sym)
	code, yx, _, err = compileChainOp(uint16(ln), atoms[2], newLookup)
	if err != nil {
		return
	}

	code = append(code, OP_EOB)
	buf := NewBytesWriter()
	buf.Write16(OP_LAMBDA)
	buf.Write16(uint16(ln))
	buf.Write16(btob(newLookup.y))
	buf.Write16(btob(newLookup.e))
	buf.Write16(uint16(len(newLookup.consts)))
	for _, k := range newLookup.consts {
		if k.ty == Tnumber {
			buf.Write16(Tnumber)
			buf.WriteDouble(k.value.(float64))
		} else if k.ty == Tstring {
			buf.Write16(Tstring)
			buf.WriteString(k.value.(string))
		} else {
			panic("shouldn't happen")
		}
	}
	buf.Write32(uint32(len(code)))
	buf.Write(code)
	return buf.data, regA, sp, nil
}

var staticWhileHack struct {
	sync.Mutex
	continueFlag []uint16
	breakFlag    []uint16
}

// [continue | break]
func compileContinueBreakOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()
	gen := func(p *[]uint16) {
		const ln = 5
		*p = make([]uint16, ln)
		buf := [ln * 2]byte{}
		if _, err := rand.Read(buf[:]); err != nil {
			panic(err)
		}
		copy(*p, (*(*[ln]uint16)(unsafe.Pointer(&buf)))[:])
	}

	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			gen(&staticWhileHack.continueFlag)
		}
		return staticWhileHack.continueFlag, regA, sp, nil
	}
	if staticWhileHack.breakFlag == nil {
		gen(&staticWhileHack.breakFlag)
	}
	return staticWhileHack.breakFlag, regA, sp, nil
}

// [while condition [chain ...]]
func compileWhileOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {

	condition := atoms[1]
	buf := NewBytesWriter()
	var varIndex uint32

	switch condition.Type {
	case parser.NTAddr:
		varIndex = condition.Value.(uint32)
	case parser.NTNumber, parser.NTString:
		buf.Write16(OP_SETK)
		varIndex = uint32(sp)
		sp++
		buf.Write32(varIndex)
		buf.Write16(table.addConst(condition.Value))
	case parser.NTCompound, parser.NTAtom:
		code, yx, sp, err = extract(sp, condition, table)
		if err != nil {
			return
		}
		buf.Write(code)
		varIndex = yx
	}

	code, yx, sp, err = compileChainOp(sp, atoms[2], table)
	if err != nil {
		return
	}

	buf.Write16(OP_IFNOT)
	buf.Write32(varIndex)
	buf.Write32(uint32(len(code)) + 3)
	buf.Write(code)
	buf.Write16(OP_JMP)
	buf.Write32(-uint32(buf.Len()) - 2)

	code = buf.data
	code2 := slice16to8(code)
	if staticWhileHack.continueFlag != nil {
		flag := slice16to8(staticWhileHack.continueFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 2
			code[idx] = OP_JMP
			code[idx+1] = uint16(uint32(-(idx + 3)) >> 16)
			code[idx+2] = uint16(-(idx + 3))
			code[idx+3] = OP_NOP
			code[idx+4] = OP_NOP
			i = idx*2 + 5
		}
	}

	if staticWhileHack.breakFlag != nil {
		flag := slice16to8(staticWhileHack.breakFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 2
			code[idx] = OP_JMP
			code[idx+1] = uint16(uint32(len(code)-idx-3) >> 16)
			code[idx+2] = uint16(len(code) - idx - 3)
			code[idx+3] = OP_NOP
			code[idx+4] = OP_NOP
			i = idx*2 + 5
		}
	}

	return code, regA, sp, nil
}
