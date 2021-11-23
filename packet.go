package nj

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

func inst(op byte, a, b uint16) _inst {
	return _inst{op: op, a: a, b: int32(b)}
}

func jmpInst(op byte, dist int) _inst {
	if dist < -(1<<30) || dist >= 1<<30 {
		panic("long jump")
	}
	return _inst{op: op, b: int32(dist)}
}

func splitInst(i _inst) (op byte, a, b uint16) {
	return i.op, i.a, uint16(i.b)
}

type posVByte []byte

func (p *posVByte) append(idx uint32, line uint32) {
	v := func(v uint64) {
		*p = append(*p, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		n := binary.PutUvarint((*p)[len(*p)-10:], v)
		*p = (*p)[:len(*p)-10+n]
	}
	v(uint64(idx))
	v(uint64(line))
}

func (p posVByte) read(i int) (next int, idx, line uint32) {
	rd := p[i:]
	a, n := binary.Uvarint(rd)
	b, n2 := binary.Uvarint(rd[n:])
	if n == 0 || n2 == 0 {
		next = len(p) + 1
		return
	}
	return i + n + n2, uint32(a), uint32(b)
}

type _inst struct {
	op byte
	a  uint16
	b  int32
}

type packet struct {
	Code []_inst
	Pos  posVByte
}

func (b *packet) writeInst(op byte, opa, opb uint16) {
	if opa == opb && op == typ.OpSet {
		return
	}
	b.Code = append(b.Code, inst(op, opa, opb))
}

func (b *packet) writeJmpInst(op byte, d int) {
	b.Code = append(b.Code, jmpInst(op, d))
}

func (b *packet) writePos(p parser.Position) {
	if p.Line == 0 {
		// Debug Code, used to detect a null meta struct
		internal.Panic("DEBUG: null line")
	}
	b.Pos.append(uint32(len(b.Code)), p.Line)
}

func (b *packet) truncateLast() {
	if len(b.Code) > 0 {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (b *packet) Len() int {
	return len(b.Code)
}

var (
	biOp = map[byte]string{
		typ.OpAdd:     parser.AAdd,
		typ.OpSub:     parser.ASub,
		typ.OpMul:     parser.AMul,
		typ.OpDiv:     parser.ADiv,
		typ.OpIDiv:    parser.AIDiv,
		typ.OpMod:     parser.AMod,
		typ.OpEq:      parser.AEq,
		typ.OpNeq:     parser.ANeq,
		typ.OpLess:    parser.ALess,
		typ.OpLessEq:  parser.ALessEq,
		typ.OpLoad:    parser.ALoad,
		typ.OpStore:   parser.AStore,
		typ.OpBitAnd:  parser.ABitAnd,
		typ.OpBitOr:   parser.ABitOr,
		typ.OpBitXor:  parser.ABitXor,
		typ.OpBitLsh:  parser.ABitLsh,
		typ.OpBitRsh:  parser.ABitRsh,
		typ.OpBitURsh: parser.ABitURsh,
	}
	uOp = map[byte]string{
		typ.OpBitNot:     parser.ABitNot,
		typ.OpNot:        parser.ANot,
		typ.OpRet:        parser.AReturn,
		typ.OpPush:       "push",
		typ.OpPushUnpack: "pushvararg",
	}
)

func pkPrettify(c *FuncBody, p *Program, toplevel bool) string {
	sb := &bytes.Buffer{}
	sb.WriteString("+ START " + c.String() + "\n")

	readAddr := func(a uint16, rValue bool) string {
		if a == regA {
			return "$a"
		}

		suffix := ""
		if rValue {
			if a > regLocalMask || toplevel {
				switch v := (*p.Stack)[a&regLocalMask]; v.Type() {
				case typ.Number, typ.Bool, typ.String:
					text := v.JSONString()
					if len(text) > 120 {
						text = text[:60] + "..." + text[len(text)-60:]
					}
					suffix = "(" + text + ")"
				case typ.Object:
					suffix = "{" + v.Object().Name() + "}"
				case typ.Native:
					suffix = "<" + v.String() + ">"
				}
			}
		}

		if a > regLocalMask {
			return fmt.Sprintf("g$%d", a&regLocalMask) + suffix
		}
		return fmt.Sprintf("$%d", a&regLocalMask) + suffix
	}

	oldpos := c.Code.Pos
	lastLine := uint32(0)

	for i, inst := range c.Code.Code {
		cursor := uint32(i) + 1
		bop, a, b := splitInst(inst)

		if len(c.Code.Pos) > 0 {
			next, op, line := c.Code.Pos.read(0)
			// log.Println(cursor, splitInst, unsafe.Pointer(&Pos))
			for uint32(cursor) > op {
				c.Code.Pos = c.Code.Pos[next:]
				if len(c.Code.Pos) == 0 {
					break
				}
				if next, op, line = c.Code.Pos.read(0); uint32(cursor) <= op {
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
				c.Code.Pos = c.Code.Pos[next:]
			} else {
				sb.WriteString(fmt.Sprintf("|     % 4d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|$    % 4d| ", cursor-1))
		}

		switch bop {
		case typ.OpSet:
			sb.WriteString(readAddr(a, false) + " = " + readAddr(b, true))
		case typ.OpArray:
			sb.WriteString("array")
		case typ.OpCreateObject:
			sb.WriteString("createobject")
		case typ.OpLoadFunc:
			cls := p.Functions[a]
			sb.WriteString("loadfunc " + cls.callable.Name + "\n")
			sb.WriteString(pkPrettify(cls.callable, p, false))
		case typ.OpTailCall, typ.OpCall:
			if b != regPhantom {
				sb.WriteString("push " + readAddr(b, true) + " -> ")
			}
			if bop == typ.OpTailCall {
				sb.WriteString("tail")
			}
			sb.WriteString("call " + readAddr(a, true))
		case typ.OpIfNot, typ.OpJmp:
			pos := inst.b
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

	c.Code.Pos = oldpos

	sb.WriteString("+ END " + c.String())
	return sb.String()
}
