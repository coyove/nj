package nj

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

func inst(op byte, a, b uint16) typ.Inst {
	return typ.Inst{Opcode: op, A: a, B: int32(b)}
}

func jmpInst(op byte, dist int) typ.Inst {
	if dist < -(1<<30) || dist >= 1<<30 {
		panic("long jump")
	}
	return typ.Inst{Opcode: op, B: int32(dist)}
}

type Packet struct {
	Code []typ.Inst
	Pos  internal.VByte32
}

func (b *Packet) writeInst(op byte, opa, opb uint16) {
	if opa == opb && op == typ.OpSet {
		return
	}
	b.Code = append(b.Code, inst(op, opa, opb))
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) writeJmpInst(op byte, d int) {
	b.Code = append(b.Code, jmpInst(op, d))
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) writeLineNum(line uint32) {
	if line == 0 {
		// Debug Code, used to detect a null meta struct
		internal.Panic("DEBUG: null line")
	}
	b.Pos.Append(uint32(len(b.Code)), line)
}

func (b *Packet) truncLast() {
	if len(b.Code) > 0 {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (b *Packet) Len() int {
	return len(b.Code)
}

func (b *Packet) LastInst() typ.Inst {
	return b.Code[len(b.Code)-1]
}

var (
	biOp = map[byte]string{
		typ.OpAdd:     typ.AAdd,
		typ.OpSub:     typ.ASub,
		typ.OpMul:     typ.AMul,
		typ.OpDiv:     typ.ADiv,
		typ.OpIDiv:    typ.AIDiv,
		typ.OpMod:     typ.AMod,
		typ.OpEq:      typ.AEq,
		typ.OpNeq:     typ.ANeq,
		typ.OpLess:    typ.ALess,
		typ.OpLessEq:  typ.ALessEq,
		typ.OpLoad:    typ.ALoad,
		typ.OpStore:   typ.AStore,
		typ.OpBitAnd:  typ.ABitAnd,
		typ.OpBitOr:   typ.ABitOr,
		typ.OpBitXor:  typ.ABitXor,
		typ.OpBitLsh:  typ.ABitLsh,
		typ.OpBitRsh:  typ.ABitRsh,
		typ.OpBitURsh: typ.ABitURsh,
		typ.OpNext:    typ.ANext,
		typ.OpIsProto: typ.AIs,
	}
	uOp = map[byte]string{
		typ.OpBitNot:     typ.ABitNot,
		typ.OpNot:        typ.ANot,
		typ.OpRet:        typ.AReturn,
		typ.OpLen:        typ.ALen,
		typ.OpPush:       "push",
		typ.OpPushUnpack: "pushvararg",
	}
)

func pkPrettify(c *function, p *Program, toplevel bool) string {
	sb := &bytes.Buffer{}
	sb.WriteString("+ START " + c.String() + "\n")

	readAddr := func(a uint16, rValue bool) string {
		if a == regA {
			return "$a"
		}

		suffix := ""
		if rValue {
			if a > regLocalMask || toplevel {
				suffix = ":" + simpleString((*p.stack)[a&regLocalMask])
			}
		}

		if a > regLocalMask {
			return fmt.Sprintf("g$%d", a&regLocalMask) + suffix
		}
		return fmt.Sprintf("$%d", a&regLocalMask) + suffix
	}

	oldpos := c.CodeSeg.Pos
	lastLine := uint32(0)

	for i, inst := range c.CodeSeg.Code {
		cursor := uint32(i) + 1
		bop, a, b := inst.Opcode, inst.A, uint16(inst.B)

		if c.CodeSeg.Pos.Len() > 0 {
			op, line := c.CodeSeg.Pos.Pop()
			// log.Println(cursor, splitInst, unsafe.Pointer(&Pos))
			for uint32(cursor) > op && c.CodeSeg.Pos.Len() > 0 {
				if op, line = c.CodeSeg.Pos.Pop(); uint32(cursor) <= op {
					break
				}
			}

			if op == uint32(cursor) {
				x := "."
				if line != lastLine {
					x = strconv.Itoa(int(line))
					lastLine = line
				}
				sb.WriteString(fmt.Sprintf("|%-4s % 4d| ", x, cursor-1))
			} else {
				sb.WriteString(fmt.Sprintf("|     % 4d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|$    % 4d| ", cursor-1))
		}

		switch bop {
		case typ.OpSet:
			sb.WriteString(readAddr(a, false) + " = " + readAddr(b, true))
		case typ.OpCreateArray:
			sb.WriteString("array")
		case typ.OpCreateObject:
			sb.WriteString("createobject")
		case typ.OpLoadFunc:
			cls := p.functions[a]
			sb.WriteString("loadfunc " + cls.fun.Name + "\n")
			sb.WriteString(pkPrettify(cls.fun, p, false))
		case typ.OpTailCall, typ.OpCall:
			if b != regPhantom {
				sb.WriteString("push " + readAddr(b, true) + " -> ")
			}
			if bop == typ.OpTailCall {
				sb.WriteString("tail")
			}
			sb.WriteString("call " + readAddr(a, true))
		case typ.OpIfNot, typ.OpJmp:
			pos := inst.B
			pos2 := uint32(int32(cursor) + pos)
			if bop == typ.OpIfNot {
				sb.WriteString("if not $a ")
			}
			sb.WriteString(fmt.Sprintf("jmp %d to %d", pos, pos2))
		case typ.OpInc:
			sb.WriteString("inc " + readAddr(a, false) + " " + readAddr(b, true))
		default:
			if us, ok := uOp[bop]; ok {
				sb.WriteString(us + " " + readAddr(a, true))
			} else if bs, ok := biOp[bop]; ok {
				sb.WriteString(bs + " " + readAddr(a, true) + " " + readAddr(b, true))
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	c.CodeSeg.Pos = oldpos

	sb.WriteString("+ END " + c.String())
	return sb.String()
}
