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

func compileSetOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	aDest, aSrc := atoms[1], atoms[2]
	varIndex := uint32(0)
	buf := newopwriter()
	var newYX uint32
	var ok bool

	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aSrc.Type != parser.NTCompound {
			newYX = uint32(sp)
			sp++
		}
	} else {
		varIndex, ok = table.get(aDest.Value.(string))
		if !ok {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aDest)
			return
		}
		newYX = varIndex
		table.clearRegRecord(varIndex)
	}

	switch aSrc.Type {
	case parser.NTAtom:
		if aSrc.Value.(string) == "nil" {
			buf.WriteOP(OP_SETK, newYX, 0)
		} else {
			valueIndex, ok := table.get(aSrc.Value.(string))
			if !ok {
				err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aSrc)
				return
			}
			buf.WriteOP(OP_SET, newYX, valueIndex)
		}
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, newYX, uint32(table.addConst(aSrc.Value)))
	case parser.NTCompound:
		code, newYX, sp, err = compileCompoundIntoVariable(sp, aSrc, table,
			atoms[0].Value.(string) == "set", varIndex)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		_, redecl := table.get(aDest.Value.(string))
		if redecl {
			err = fmt.Errorf("redeclare: %+v", aDest)
			return
		}
		table.put(aDest.Value.(string), uint16(newYX))
	}
	return buf.data, newYX, sp, nil
}

func compileRetOp(op, opk byte) compileFunc {
	return func(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
		buf := newopwriter()
		if len(atoms) == 1 {
			buf.WriteOP(op, regA, 0)
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
			buf.WriteOP(op, yx, 0)
		}

		if op == OP_YIELD {
			table.y = true
		}
		return buf.data, yx, sp, nil
	}
}

func compileMapOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	if len(atoms[1].Compound)%2 != 0 {
		err = fmt.Errorf("every key in map must have a value: %+v", atoms[1])
		return
	}
	var buf *opwriter
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
	buf.WriteOP(OP_MAKEMAP, 0, 0)
	return buf.data, regA, sp, nil
}

func flaten(sp uint16, atoms []*parser.Node, table *symtable) (buf *opwriter, newsp uint16, err error) {
	replacedAtoms := []*parser.Node{}
	buf = newopwriter()

	for i, atom := range atoms {
		var yx uint32
		var code []uint64

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

	if len(replacedAtoms) > 0 {
		_, _, replacedAtoms[len(replacedAtoms)-1].Value = op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		sp--
	}

	return buf, sp, nil
}

func flatWrite(sp uint16, atoms []*parser.Node, table *symtable, bop byte) (code []uint64, yx uint32, newsp uint16, err error) {
	var buf *opwriter
	buf, sp, err = flaten(sp, atoms[1:], table)
	if err != nil {
		return
	}

	immediateRet := func(n float64) {
		buf.Clear()
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(n)))
		table.im = &n
	}

	v1 := func() float64 { return atoms[1].Value.(float64) }
	v2 := func() float64 { return atoms[2].Value.(float64) }

	switch bop {
	case OP_ADD:
		if atoms[1].Type == parser.NTString && atoms[2].Type == parser.NTString {
			str := atoms[1].Value.(string) + atoms[2].Value.(string)
			buf.Clear()
			buf.WriteOP(OP_SETK, regA, uint32(table.addConst(str)))
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

	var op = [4]uint16{OP_R0<<8 | OP_R0K, OP_R1<<8 | OP_R1K, OP_R2<<8 | OP_R2K, OP_R3<<8 | OP_R3K}
	switch bop {
	case OP_LEN, OP_STORE, OP_LOAD:
		op[0], op[3] = op[3], op[0]
		op[1], op[2] = op[2], op[1]
	}

	count := buf.Len()
	for i := 1; i < len(atoms); i++ {
		if err = fill(buf, atoms[i], table, byte(op[i-1]>>8), byte(op[i-1])); err != nil {
			return
		}
	}
	if buf.Len() == count {
		// TODO:
		// No argument opcode was written into buf, which means either this op accepts 0 argument,
		// or all arguments have been inside the registers already.
		// So we should evaluate the case that A still holds the result and skip the current op, e.g.:

		// set a = (i + j) / (i + j + 1) =>
		// R0 = i, R1 = j, A = i + j
	}
	buf.WriteOP(bop, 0, 0)

	if bop == OP_ASSERT {
		buf.WriteString(atoms[0].String())
	}

	return buf.data, regA, sp, nil
}

func compileFlatOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	head := atoms[0].Value.(string)
	op := flatOpMapping[head]
	if op == 0 {
		err = fmt.Errorf("invalid flat op %+v", atoms[0])
		return
	}
	return flatWrite(sp, atoms, table, op)
}

func compileIncOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	subject, ok := table.get(atoms[1].Value.(string))
	buf := newopwriter()
	if !ok {
		return nil, 0, 0, fmt.Errorf(ERR_UNDECLARED_VARIABLE, atoms[1])
	}
	table.clearRegRecord(subject)
	buf.WriteOP(OP_INC, subject, uint32(table.addConst(atoms[2].Value)))
	return buf.data, regA, sp, nil
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func compileAndOrOp(bop byte) compileFunc {
	return func(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
		a1, a2 := atoms[1], atoms[2]
		buf := newopwriter()
		switch a1.Type {
		case parser.NTNumber, parser.NTString:
			buf.WriteOP(OP_SETK, regA, uint32(table.addConst(a1.Value)))
			buf.WriteOP(bop, regA, 0)
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a1, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			buf.WriteOP(OP_SET, regA, yx)
			buf.WriteOP(bop, regA, 0)
		}
		c2 := buf.Len()

		switch a2.Type {
		case parser.NTNumber, parser.NTString:
			buf.WriteOP(OP_SETK, regA, uint32(table.addConst(a2.Value)))
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a2, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			if yx != regA {
				buf.WriteOP(OP_SET, regA, yx)
			}
		}
		jmp := buf.Len() - c2

		_, yx, _ = op(buf.data[c2-1])
		buf.data[c2-1] = makeop(bop, yx, uint32(int32(jmp)))
		return buf.data, regA, sp, nil
	}
}

// [if condition [truechain ...] [falsechain ...]]
func compileIfOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := newopwriter()

	switch condition.Type {
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(condition.Value)))
		yx = regA
	case parser.NTAtom, parser.NTCompound:
		code, yx, sp, err = extract(sp, condition, table)
		if err != nil {
			return nil, 0, 0, err
		}
		buf.Write(code)
	}
	condyx := yx
	var trueCode, falseCode []uint64
	trueCode, yx, sp, err = compileChainOp(sp, trueBranch, table)
	if err != nil {
		return
	}
	falseCode, yx, sp, err = compileChainOp(sp, falseBranch, table)
	if err != nil {
		return
	}
	if len(falseCode) > 0 {
		buf.WriteOP(OP_IFNOT, condyx, uint32(len(trueCode))+1)
		buf.Write(trueCode)
		buf.WriteOP(OP_JMP, 0, uint32(len(falseCode)))
		buf.Write(falseCode)
	} else {
		buf.WriteOP(OP_IFNOT, condyx, uint32(len(trueCode)))
		buf.Write(trueCode)
	}
	table.clearAllRegRecords()
	return buf.data, regA, sp, nil
}

// [call func-name [args ...]]
func compileCallOp(sp uint16, nodes []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	buf := newopwriter()
	callee := nodes[1]
	name, _ := callee.Value.(string)
	switch name {
	case "error":
		table.e = true
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_ERROR)
	case "len":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_LEN)
	case "dup":
		x := append(nodes[1:2], nodes[2].Compound...)
		if y, ok := x[3].Value.(float64); ok && y == 2 {
			// return stack, env is escaped
			table.envescape = true
		}
		return flatWrite(sp, x, table, OP_DUP)
	case "typeof":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_TYPEOF)
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
		_, _, replacedAtoms[0].Value = op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
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
		if len(replacedAtoms) == 0 {
			_, _, varIndex = op(code[len(code)-1])
			code = code[:len(code)-1]
			sp--
		}
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

	buf.WriteOP(OP_CALL, varIndex, 0)
	return buf.data, regA, sp, nil
}

// [lambda [namelist] [chain ...]]
func compileLambdaOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	table.envescape = true
	newtable := newsymtable()
	newtable.parent = table
	newtable.lineInfo = table.lineInfo

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("%+v: invalid arguments list", atoms[0])
		return
	}

	var this bool
	i := 0
	for _, p := range params.Compound {
		argname := p.Value.(string)
		if argname == "this" {
			this = true
			continue
		}
		newtable.put(argname, uint16(i))
		i++
	}

	ln := len(newtable.sym)
	if this {
		newtable.put("this", uint16(ln))
	}

	if ln > 255 {
		return nil, 0, 0, fmt.Errorf("do you really need more than 255 arguments?")
	}
	code, yx, _, err = compileChainOp(uint16(ln), atoms[2], newtable)
	if err != nil {
		return
	}

	code = append(code, makeop(OP_EOB, 0, 0))
	buf := newopwriter()
	buf.WriteOP(OP_LAMBDA, uint32(len(newtable.consts)),
		uint32(byte(ln))<<24+
			uint32(btob(newtable.y))<<20+
			uint32(btob(newtable.e))<<16+
			uint32(btob(!newtable.envescape))<<12+
			uint32(btob(this))<<8)
	for _, k := range newtable.consts {
		if k.ty == Tnumber {
			buf.Write64(Tnumber)
			buf.WriteDouble(k.value.(float64))
		} else if k.ty == Tstring {
			buf.Write64(Tstring)
			buf.WriteString(k.value.(string))
		} else {
			panic("shouldn't happen")
		}
	}
	buf.Write64(uint64(len(code)))
	buf.Write(code)
	return buf.data, regA, sp, nil
}

var staticWhileHack struct {
	sync.Mutex
	continueFlag []uint64
	breakFlag    []uint64
}

// [continue | break]
func compileContinueBreakOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()
	gen := func(p *[]uint64) {
		const ln = 2
		*p = make([]uint64, ln)
		buf := [ln * 8]byte{}
		if _, err := rand.Read(buf[:]); err != nil {
			panic(err)
		}
		copy(*p, (*(*[ln]uint64)(unsafe.Pointer(&buf)))[:])
	}

	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			gen(&staticWhileHack.continueFlag)
		}
		buf := newopwriter()
		code, yx, sp, err = compileChainOp(sp, table.continueNode[len(table.continueNode)-1], table)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.Write(staticWhileHack.continueFlag)
		return buf.data, regA, sp, nil
	}
	if staticWhileHack.breakFlag == nil {
		gen(&staticWhileHack.breakFlag)
	}
	return staticWhileHack.breakFlag, regA, sp, nil
}

// [for condition incr [chain ...]]
func compileWhileOp(sp uint16, atoms []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	table.clearAllRegRecords()
	condition := atoms[1]
	buf := newopwriter()
	var varIndex uint32

	switch condition.Type {
	case parser.NTAddr:
		varIndex = condition.Value.(uint32)
	case parser.NTNumber, parser.NTString:
		varIndex = uint32(sp)
		sp++
		buf.WriteOP(OP_SETK, varIndex, uint32(table.addConst(condition.Value)))
	case parser.NTCompound, parser.NTAtom:
		code, yx, sp, err = extract(sp, condition, table)
		if err != nil {
			return
		}
		buf.Write(code)
		varIndex = yx
	}

	table.continueNode = append(table.continueNode, atoms[2])
	code, yx, sp, err = compileChainOp(sp, atoms[3], table)
	if err != nil {
		return
	}
	var icode []uint64
	icode, yx, sp, err = compileChainOp(sp, atoms[2], table)
	if err != nil {
		return
	}
	table.continueNode = table.continueNode[:len(table.continueNode)-1]

	code = append(code, icode...)
	buf.WriteOP(OP_IFNOT, varIndex, uint32(len(code))+1)
	buf.Write(code)
	buf.WriteOP(OP_JMP, 0, -uint32(buf.Len())-1)

	code = buf.data
	code2 := slice64to8(code)
	if staticWhileHack.continueFlag != nil {
		flag := slice64to8(staticWhileHack.continueFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 8
			code[idx] = makeop(OP_JMP, 0, uint32(int32(-idx-1)))
			code[idx+1] = makeop(OP_NOP, 0, 0)
			i = idx*8 + 2
		}
	}

	if staticWhileHack.breakFlag != nil {
		flag := slice64to8(staticWhileHack.breakFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 8
			code[idx] = makeop(OP_JMP, 0, uint32(int32(len(code)-idx-1)))
			code[idx+1] = makeop(OP_NOP, 0, 0)
			i = idx*8 + 2
		}
	}
	table.clearAllRegRecords()
	return buf.data, regA, sp, nil
}
