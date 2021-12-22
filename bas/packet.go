package bas

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Packet struct {
	Code []typ.Inst
	Pos  internal.VByte32
}

func (b *Packet) WriteInst(op byte, opa, opb uint16) {
	if opa == opb && op == typ.OpSet {
		return
	}
	b.Code = append(b.Code, typ.Inst{Opcode: op, A: opa, B: int32(opb)})
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) WriteJmpInst(op byte, d int) {
	b.Code = append(b.Code, typ.JmpInst(op, d))
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) WriteLineNum(line uint32) {
	if line == 0 {
		// Debug Code, used to detect a null meta struct
		internal.Panic("DEBUG: null line")
	}
	b.Pos.Append(uint32(len(b.Code)), line)
}

func (b *Packet) TruncLast() {
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

func pkPrettify(c *Function, p *Program, toplevel bool) string {
	sb := &bytes.Buffer{}
	sb.WriteString("+ START " + c.String() + "\n")

	readAddr := func(a uint16, rValue bool) string {
		if a == typ.RegA {
			return "$a"
		}

		suffix := ""
		if rValue {
			if a > typ.RegLocalMask || toplevel {
				suffix = ":" + simpleString((*p.stack)[a&typ.RegLocalMask])
			}
		}

		if a > typ.RegLocalMask {
			return fmt.Sprintf("g$%d", a&typ.RegLocalMask) + suffix
		}
		return fmt.Sprintf("$%d", a&typ.RegLocalMask) + suffix
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
			if b != typ.RegPhantom {
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
			if us, ok := typ.UnaryOpcode[bop]; ok {
				sb.WriteString(us + " " + readAddr(a, true))
			} else if bs, ok := typ.BinaryOpcode[bop]; ok {
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
