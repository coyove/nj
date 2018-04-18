package compiler

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/coyove/bracket/base"
)

var staticWhileHack struct {
	sync.Mutex
	continueFlag []byte
	breakFlag    []byte
}

func compileContinueBreakOp(
	stackPtr int16,
	atoms []*token,
	varLookup *base.CMap,
) (code []byte, yx int32, newStackPtr int16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()

	if atoms[0].v.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			staticWhileHack.continueFlag = make([]byte, 9)
			if _, err := rand.Read(staticWhileHack.continueFlag); err != nil {
				panic(err)
			}
		}
		return staticWhileHack.continueFlag, base.REG_A, stackPtr, nil
	}

	if staticWhileHack.breakFlag == nil {
		staticWhileHack.breakFlag = make([]byte, 9)
		if _, err := rand.Read(staticWhileHack.breakFlag); err != nil {
			panic(err)
		}
	}
	return staticWhileHack.breakFlag, base.REG_A, stackPtr, nil
}

func compileWhileOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("while statement should have condition and body: %+v", atoms[0])
		return
	}

	condition := atoms[1]
	buf := base.NewBytesBuffer()
	var varIndex int32

	switch condition.ty {
	case TK_addr:
		varIndex = condition.v.(int32)
	case TK_number:
		buf.WriteByte(base.OP_SET_NUM)
		varIndex = int32(stackPtr)
		stackPtr++
		buf.WriteInt32(varIndex)
		buf.WriteDouble(condition.v.(float64))
	case TK_string:
		buf.WriteByte(base.OP_SET_STR)
		varIndex = int32(stackPtr)
		stackPtr++
		buf.WriteInt32(varIndex)
		buf.WriteString(condition.v.(string))
	case TK_compound, TK_atomic:
		code, yx, stackPtr, err = extract(stackPtr, condition, varLookup)
		if err != nil {
			return
		}
		buf.Write(code)
		varIndex = yx
	}

	code, yx, stackPtr, err = compile(stackPtr, atoms[2:], varLookup)
	if err != nil {
		return
	}

	code = truncLastByte(code)

	buf.WriteByte(base.OP_IF)
	buf.WriteInt32(varIndex)
	buf.WriteInt32(int32(len(code)) + 5)
	buf.Write(code)
	buf.WriteByte(base.OP_JMP)
	buf.WriteInt32(-int32(buf.Len()) - 4)

	code = buf.Bytes()
	i := 0
	for i < len(code) && staticWhileHack.continueFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.continueFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = base.OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32(-(idx + 5)))
		copy(code[idx+5:], []byte{base.OP_NOP, base.OP_NOP, base.OP_NOP, base.OP_NOP})
		i = idx + 9
		// 				buf.WriteInt(i + 1, bop == Op.JMP_BREAK ?
		// 						(f.code.size() - i) - 5 :
		// 						-(i + 5));
	}

	i = 0
	for i < len(code) && staticWhileHack.breakFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.breakFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = base.OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32((len(code)-idx)-5))
		copy(code[idx+5:], []byte{base.OP_NOP, base.OP_NOP, base.OP_NOP, base.OP_NOP})
		i = idx + 9
	}

	return code, base.REG_A, stackPtr, nil
}
