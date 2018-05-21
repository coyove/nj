package potatolang

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/coyove/potatolang/parser"
)

const (
	ERR_UNDECLARED_VARIABLE = "undeclared variable: %+v"
)

func compileSetOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	aVar := atoms[1]
	varIndex := int32(0)
	if len(atoms) < 3 {
		err = fmt.Errorf("can't set/move without value %+v", atoms[0])
		return
	}

	aValue := atoms[2]

	buf := NewBytesWriter()
	var newYX int32
	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aValue.Type != parser.NTCompound {
			newYX = int32(sp)
			sp++
		}
	} else {
		varIndex = table.GetRelPosition(aVar.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aVar)
			return
		}
		newYX = varIndex
	}

	switch aValue.Type {
	case parser.NTAtom:
		valueIndex := table.GetRelPosition(aValue.Value.(string))
		if valueIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aValue)
			return
		}

		buf.WriteByte(OP_SET)
		buf.WriteInt32(newYX)
		buf.WriteInt32(valueIndex)
	case parser.NTNumber:
		buf.WriteByte(OP_SET_NUM)
		buf.WriteInt32(newYX)
		buf.WriteDouble(aValue.Value.(float64))
	case parser.NTString:
		buf.WriteByte(OP_SET_STR)
		buf.WriteInt32(newYX)
		buf.WriteString(aValue.Value.(string))
	case parser.NTCompound:
		code, newYX, sp, err = compileCompoundIntoVariable(sp, aValue, table,
			atoms[0].Value.(string) == "set", varIndex)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		_, reset := table.M[aVar.Value.(string)]
		if reset {
			err = fmt.Errorf("redeclare: %+v", aVar)
			return
		}

		table.M[aVar.Value.(string)] = int16(newYX)
	}

	table.I = nil
	return buf.Bytes(), newYX, sp, nil
}

func compileRetOp(r, n, s byte) compileFunc {
	return func(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
		buf := NewBytesWriter()
		if len(atoms) == 1 {
			buf.WriteByte(r)
			buf.WriteInt32(REG_A)
			table.I = nil
			return buf.Bytes(), yx, sp, nil
		}

		atom := atoms[1]

		switch atom.Type {
		case parser.NTAtom, parser.NTNumber, parser.NTString, parser.NTAddr:
			err = fill1(buf, atom, table, r, n, s)
			if err != nil {
				return
			}
		case parser.NTCompound:
			code, yx, sp, err = extract(sp, atom, table)
			buf.Write(code)
			buf.WriteByte(r)
			buf.WriteInt32(yx)
		}

		if r == OP_YIELD {
			table.Y = true
		}

		return buf.Bytes(), yx, sp, nil
	}
}

func compileListOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	var buf *BytesWriter
	buf, sp, err = flaten(sp, atoms[1].Compound, table)
	if err != nil {
		return
	}

	for _, atom := range atoms[1].Compound {
		err = fill1(buf, atom, table, OP_PUSH, OP_PUSH_NUM, OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(OP_LIST)
	table.I = nil
	return buf.Bytes(), REG_A, sp, nil
}

func compileMapOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
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
		err = fill1(buf, atom, table, OP_PUSH, OP_PUSH_NUM, OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(OP_MAP)
	table.I = nil
	return buf.Bytes(), REG_A, sp, nil
}

func indexToOpR(index int) []byte {
	switch index {
	case 0:
		return []byte{OP_R0, OP_R0_NUM, OP_R0_STR}
	case 1:
		return []byte{OP_R1, OP_R1_NUM, OP_R1_STR}
	case 2:
		return []byte{OP_R2, OP_R2_NUM, OP_R2_STR}
	case 3:
		return []byte{OP_R3, OP_R3_NUM, OP_R3_STR}
	}
	panic("shouldn't happen")
}

func flaten(sp int16, atoms []*parser.Node, table *symtable) (buf *BytesWriter, newsp int16, err error) {

	replacedAtoms := []*parser.Node{}
	buf = NewBytesWriter()

	for i, atom := range atoms {

		var yx int32
		var code []byte

		if atom.Type == parser.NTCompound {
			code, yx, sp, err = compileCompoundIntoVariable(sp, atom, table, true, 0)
			if err != nil {
				return
			}
			if table.I != nil {
				atoms[i] = &parser.Node{Type: parser.NTNumber, Value: *table.I}
				table.I = nil
			} else if table.Is != nil {
				atoms[i] = &parser.Node{Type: parser.NTString, Value: *table.Is}
				table.Is = nil
			} else {
				atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
				replacedAtoms = append(replacedAtoms, atoms[i])
				buf.Write(code)
			}
		}
	}

	if len(replacedAtoms) == 1 {
		cursor := buf.Len() - 4
		replacedAtoms[0].Value = int32(binary.LittleEndian.Uint32(buf.Bytes()[cursor:]))
		buf.TruncateLastBytes(9)
		sp--
	}

	return buf, sp, nil
}

func flatWrite(sp int16, atoms []*parser.Node, table *symtable, bop byte) (code []byte, yx int32, newsp int16, err error) {

	var buf *BytesWriter
	buf, sp, err = flaten(sp, atoms[1:], table)
	if err != nil {
		return
	}

	immediateRet := func(n float64) {
		buf.Clear()
		buf.WriteByte(OP_SET_NUM)
		buf.WriteInt32(REG_A)
		table.I = &n
		buf.WriteDouble(n)
	}

	v1 := func() float64 { return atoms[1].Value.(float64) }
	v2 := func() float64 { return atoms[2].Value.(float64) }

	switch bop {
	case OP_ADD:
		if atoms[1].Type == parser.NTString && atoms[2].Type == parser.NTString {
			buf.Clear()
			buf.WriteByte(OP_SET_STR)
			buf.WriteInt32(REG_A)
			str := atoms[1].Value.(string) + atoms[2].Value.(string)
			buf.WriteString(str)
			table.Is = &str
			return buf.Bytes(), REG_A, sp, nil
		}
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() + v2())
			return buf.Bytes(), REG_A, sp, nil
		}
	case OP_SUB:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() - v2())
			return buf.Bytes(), REG_A, sp, nil
		}
	case OP_MUL:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() * v2())
			return buf.Bytes(), REG_A, sp, nil
		}
	case OP_DIV:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() / v2())
			return buf.Bytes(), REG_A, sp, nil
		}
	}

	for i := 1; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], table, indexToOpR(i-1)...)
		if err != nil {
			return
		}
	}
	buf.WriteByte(bop)

	if bop == OP_ASSERT {
		buf.WriteString(atoms[0].String())
	}

	return buf.Bytes(), REG_A, sp, nil
}

func compileFlatOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {

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

func compileIncOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	subject, step, buf := table.GetRelPosition(atoms[1].Value.(string)), atoms[2].Value.(float64), NewBytesWriter()
	if subject == -1 {
		return nil, 0, 0, fmt.Errorf(ERR_UNDECLARED_VARIABLE, atoms[1])
	}
	buf.WriteByte(OP_INC)
	buf.WriteInt32(subject)
	buf.WriteDouble(step)
	return buf.Bytes(), REG_A, sp, nil
}

func compileImmediateIntoVariable(sp int16, node *parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	buf := NewBytesWriter()
	switch node.Type {
	case parser.NTNumber:
		buf.WriteByte(OP_SET_NUM)
		buf.WriteInt32(int32(sp))
		yx = int32(sp)
		buf.WriteDouble(node.Value.(float64))
		sp++
	case parser.NTString:
		buf.WriteByte(OP_SET_STR)
		buf.WriteInt32(int32(sp))
		yx = int32(sp)
		buf.WriteString(node.Value.(string))
		sp++
	default:
		panic("shouldn't happen")
	}
	return buf.Bytes(), yx, sp, nil
}

// [and a b] => if not a then false else use b end
// [or a b]  => if a then use a else use b end
func compileAndOrOp(bop byte) compileFunc {
	return func(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
		a1, a2 := atoms[1], atoms[2]
		buf := NewBytesWriter()
		switch a1.Type {
		case parser.NTNumber, parser.NTString:
			code, yx, sp, err = compileImmediateIntoVariable(sp, a1, table)
			buf.Write(code)
			buf.WriteByte(bop)
			buf.WriteInt32(yx)
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a1, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			buf.WriteByte(bop)
			buf.WriteInt32(yx)
		}
		buf.WriteInt32(0)
		c2 := buf.Len()

		switch a2.Type {
		case parser.NTNumber:
			buf.WriteByte(OP_SET_NUM)
			buf.WriteInt32(REG_A)
			buf.WriteDouble(a2.Value.(float64))
		case parser.NTString:
			buf.WriteByte(OP_SET_STR)
			buf.WriteInt32(REG_A)
			buf.WriteString(a2.Value.(string))
		case parser.NTAtom, parser.NTCompound:
			code, yx, sp, err = extract(sp, a2, table)
			if err != nil {
				return nil, 0, 0, err
			}
			buf.Write(code)
			if yx != REG_A {
				buf.WriteByte(OP_SET)
				buf.WriteInt32(REG_A)
				buf.WriteInt32(yx)
			}
		}
		jmp := buf.Len() - c2

		code = buf.Bytes()
		binary.LittleEndian.PutUint32(code[c2-4:], uint32(jmp))
		return code, REG_A, sp, nil
	}
}

// [if condition [truechain ...] [falsechain ...]]
func compileIfOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("if statement should have at least a true branch: %+v", atoms[0])
		return
	}

	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := NewBytesWriter()

	switch condition.Type {
	case parser.NTNumber, parser.NTString:
		code, yx, sp, err = compileImmediateIntoVariable(sp, condition, table)
		buf.Write(code)
		buf.WriteByte(OP_IFNOT)
		buf.WriteInt32(yx)
	case parser.NTAtom, parser.NTCompound:
		code, yx, sp, err = extract(sp, condition, table)
		if err != nil {
			return nil, 0, 0, err
		}

		buf.Write(code)
		buf.WriteByte(OP_IFNOT)
		buf.WriteInt32(yx)

		var trueCode, falseCode []byte
		trueCode, yx, sp, err = compileChainOp(sp, trueBranch, table)
		if err != nil {
			return
		}
		falseCode, yx, sp, err = compileChainOp(sp, falseBranch, table)
		if err != nil {
			return
		}
		if len(falseCode) > 0 {
			buf.WriteInt32(int32(len(trueCode)) + 5) // jmp (1b) + offset (4b)
			buf.Write(trueCode)
			buf.WriteByte(OP_JMP)
			buf.WriteInt32(int32(len(falseCode)))
			buf.Write(falseCode)
		} else {
			buf.WriteInt32(int32(len(trueCode)))
			buf.Write(trueCode)
		}
	}
	return buf.Bytes(), REG_A, sp, nil
}

// [call func-name [args ...]]
func compileCallOp(sp int16, nodes []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	buf := NewBytesWriter()
	callee := nodes[1]

	name, _ := callee.Value.(string)
	switch name {
	case "error":
		table.E = true
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_ERROR)
	case "len":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_LEN)
	case "dup":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_DUP)
	case "typeof":
		return flatWrite(sp, append(nodes[1:2], nodes[2].Compound...), table, OP_TYPEOF)
	case "who":
		return []byte{OP_WHO}, REG_A, sp, nil
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
		cursor := buf.Len() - 4
		replacedAtoms[0].Value = int32(binary.LittleEndian.Uint32(buf.Bytes()[cursor:]))
		buf.TruncateLastBytes(9)
		sp--
	}

	var varIndex int32
	switch callee.Type {
	case parser.NTAtom:
		varIndex = table.GetRelPosition(callee.Value.(string))
		if varIndex == -1 {
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
		varIndex = callee.Value.(int32)
	case parser.NTString:
		buf.WriteByte(OP_SET_STR)
		buf.WriteInt32(int32(sp))
		buf.WriteString(callee.Value.(string))
		varIndex = int32(sp)
		sp++
	default:
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], table, OP_PUSH, OP_PUSH_NUM, OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(OP_CALL)
	buf.WriteInt32(varIndex)

	table.I = nil
	return buf.Bytes(), REG_A, sp, nil
}

// [lambda [namelist] [chain ...]]
func compileLambdaOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	newLookup := &symtable{M: make(map[string]int16)}
	newLookup.Parent = table
	newLookup.LineInfo = table.LineInfo

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("invalid lambda parameters: %+v", atoms[0])
		return
	}

	for i, p := range params.Compound {
		newLookup.M[p.Value.(string)] = int16(i)
	}

	ln := len(newLookup.M)
	code, yx, _, err = compileChainOp(int16(ln), atoms[2], newLookup)
	if err != nil {
		return
	}

	code = append(code, OP_EOB)
	buf := NewBytesWriter()
	buf.WriteByte(OP_LAMBDA)
	buf.WriteByte(byte(ln))
	buf.WriteByte(btob(newLookup.Y))
	buf.WriteByte(btob(newLookup.E))
	buf.WriteInt32(int32(len(code)))
	buf.Write(code)

	return buf.Bytes(), REG_A, sp, nil
}

var staticWhileHack struct {
	sync.Mutex
	continueFlag []byte
	breakFlag    []byte
}

// [continue | break]
func compileContinueBreakOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()
	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			staticWhileHack.continueFlag = make([]byte, 9)
			if _, err := rand.Read(staticWhileHack.continueFlag); err != nil {
				panic(err)
			}
		}
		return staticWhileHack.continueFlag, REG_A, sp, nil
	}

	if staticWhileHack.breakFlag == nil {
		staticWhileHack.breakFlag = make([]byte, 9)
		if _, err := rand.Read(staticWhileHack.breakFlag); err != nil {
			panic(err)
		}
	}
	return staticWhileHack.breakFlag, REG_A, sp, nil
}

// [while condition [chain ...]]
func compileWhileOp(sp int16, atoms []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("while statement should have condition and body: %+v", atoms[0])
		return
	}
	condition := atoms[1]
	buf := NewBytesWriter()
	var varIndex int32

	switch condition.Type {
	case parser.NTAddr:
		varIndex = condition.Value.(int32)
	case parser.NTNumber:
		buf.WriteByte(OP_SET_NUM)
		varIndex = int32(sp)
		sp++
		buf.WriteInt32(varIndex)
		buf.WriteDouble(condition.Value.(float64))
	case parser.NTString:
		buf.WriteByte(OP_SET_STR)
		varIndex = int32(sp)
		sp++
		buf.WriteInt32(varIndex)
		buf.WriteString(condition.Value.(string))
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

	buf.WriteByte(OP_IFNOT)
	buf.WriteInt32(varIndex)
	buf.WriteInt32(int32(len(code)) + 5)
	buf.Write(code)
	buf.WriteByte(OP_JMP)
	buf.WriteInt32(-int32(buf.Len()) - 4)

	code = buf.Bytes()
	i := 0
	for i < len(code) && staticWhileHack.continueFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.continueFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32(-(idx + 5)))
		copy(code[idx+5:], []byte{OP_NOP, OP_NOP, OP_NOP, OP_NOP})
		i = idx + 9
	}

	i = 0
	for i < len(code) && staticWhileHack.breakFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.breakFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32((len(code)-idx)-5))
		copy(code[idx+5:], []byte{OP_NOP, OP_NOP, OP_NOP, OP_NOP})
		i = idx + 9
	}

	return code, REG_A, sp, nil
}
