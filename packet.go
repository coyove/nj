package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/coyove/script/parser"
)

func makeop(op opCode, a, b uint16) uint32 {
	// 6 + 13 + 13
	return uint32(op)<<26 + uint32(a&0x1fff)<<13 + uint32(b&0x1fff)
}

func makejmpop(op opCode, a uint16, dist int) uint32 {
	if dist < -(1<<12) || dist >= 1<<12 {
		panic("long jump")
	}
	// 6 + 13 + 13
	b := uint16(dist + 1<<12)
	return uint32(op)<<26 + uint32(a&0x1fff)<<13 + uint32(b&0x1fff)
}

func op(x uint32) (op opCode, a, b uint16) {
	op = opCode(x >> 26)
	a = uint16(x>>13) & 0x1fff
	b = uint16(x) & 0x1fff
	return
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

type packet struct {
	Funcs  []*Func
	Code   []uint32
	Pos    posVByte
	Source string
}

func (b *packet) write(buf packet) {
	datalen := len(b.Code)
	b.Code = append(b.Code, buf.Code...)
	i := 0
	for i < len(buf.Pos) {
		var idx, line uint32
		i, idx, line = buf.Pos.read(i)
		b.Pos.append(idx+uint32(datalen), line)
	}
	if b.Source == "" {
		b.Source = buf.Source
	}
}

func (b *packet) writeBytes(buf []uint32) {
	b.Code = append(b.Code, buf...)
}

func (b *packet) write32(v uint32) {
	b.Code = append(b.Code, v)
}

func (b *packet) writeOP(op opCode, opa, opb uint16) {
	b.Code = append(b.Code, makeop(op, opa, opb))
}

func (b *packet) writeJmpOP(op opCode, opa uint16, d int) {
	b.Code = append(b.Code, makejmpop(op, opa, d))
}

func (b *packet) writePos(p parser.Position) {
	if p.Line == 0 {
		// Debug Code, used to detect a null meta struct
		panicf("null line")
	}
	b.Pos.append(uint32(len(b.Code)), p.Line)
	if p.Source != "" {
		b.Source = p.Source
	}
}

func (b *packet) truncateLast() {
	if len(b.Code) > 0 {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (b *packet) Len() int {
	return len(b.Code)
}

func pkRead(data []uint32, len int, cursor *uint32) []uint32 {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func pkRead32(data []uint32, cursor *uint32) uint32 {
	*cursor++
	return data[*cursor-1]
}

var singleOp = map[opCode]parser.Symbol{
	OpConcat: parser.AConcat,
	OpAdd:    parser.AAdd,
	OpSub:    parser.ASub,
	OpMul:    parser.AMul,
	OpDiv:    parser.ADiv,
	OpMod:    parser.AMod,
	OpEq:     parser.AEq,
	OpNeq:    parser.ANeq,
	OpLess:   parser.ALess,
	OpLessEq: parser.ALessEq,
	OpLen:    parser.ALen,
	OpLoad:   parser.ALoad,
	OpStore:  parser.AStore,
	OpNot:    parser.ANot,
	OpPow:    parser.APow,
}

func pkPrettify(c *Func, tab int) string {
	sb := &bytes.Buffer{}
	spaces := strings.Repeat("        ", tab)
	spaces2 := ""
	if tab > 0 {
		spaces2 = strings.Repeat("        ", tab-1) + "+-------"
	}

	sb.WriteString(spaces2 + "+ START " + c.String() + " " + c.Source + "\n")

	var cursor uint32
	readAddr := func(a uint16) string {
		if a == regA {
			return "$a"
		}
		if a == regNil {
			return "nil"
		}
		if a>>10 == 7 {
			return fmt.Sprintf("k$%d(%v)", a&0x03ff, c.ConstTable[a&0x3ff].toString(0))
		}
		if a>>10 == 1 {
			return fmt.Sprintf("g$%d", a&0x03ff)
		}
		return fmt.Sprintf("$%d", a&0x03ff)
	}

	oldpos := c.Pos
	lastLine := uint32(0)
MAIN:
	for {
		bop, a, b := op(pkRead32(c.Code, &cursor))
		sb.WriteString(spaces)

		if len(c.Pos) > 0 {
			next, op, line := c.Pos.read(0)
			// log.Println(cursor, op, unsafe.Pointer(&Pos))
			for cursor > op {
				c.Pos = c.Pos[next:]
				if len(c.Pos) == 0 {
					break
				}
				if next, op, line = c.Pos.read(0); cursor <= op {
					break
				}
			}

			if op == cursor {
				x := "."
				if line != lastLine {
					x = strconv.Itoa(int(line))
					lastLine = line
				}
				sb.WriteString(fmt.Sprintf("|%-4s % 4d| ", x, cursor-1))
				c.Pos = c.Pos[next:]
			} else {
				sb.WriteString(fmt.Sprintf("|     % 4d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|$    % 4d| ", cursor-1))
		}

		switch bop {
		case OpEOB:
			sb.WriteString("end\n")
			break MAIN
		case OpSet:
			sb.WriteString(readAddr(a) + " = " + readAddr(b))
		case OpPopVClear:
			sb.WriteString("clear-v")
		case OpPopVAll:
			switch a {
			case 0:
				sb.WriteString("$a = popv-all")
			case 1:
				sb.WriteString("$a = popv-all-with-a")
			}
		case OpPopV:
			sb.WriteString("$a = pop-v")
		case OpPushV:
			sb.WriteString("pushv " + readAddr(a))
			if b != 0 {
				sb.WriteString(" cap=" + strconv.Itoa(int(b)))
			}
		case OpPush:
			sb.WriteString(fmt.Sprintf("push-%d %v", b, readAddr(a)))
			if a == regA {
				sb.WriteString(" $v")
			}
		case OpRet:
			sb.WriteString("ret " + readAddr(a))
		case OpYield:
			sb.WriteString("yield " + readAddr(a))
		case OpLoadFunc:
			sb.WriteString("$a = closure:\n")
			cls := c.Funcs[a]
			sb.WriteString(pkPrettify(cls, tab+1))
		case OpCall:
			if b == 1 {
				sb.WriteString("tail-call " + readAddr(a))
			} else {
				sb.WriteString("call " + readAddr(a))
			}
		case OpJmp:
			pos := int32(b) - 1<<12
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OpIf, OpIfNot:
			addr := readAddr(a)
			pos := int32(b) - 1<<12
			pos2 := strconv.Itoa(int(int32(cursor) + pos))
			if bop == OpIfNot {
				sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			} else {
				sb.WriteString("if " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			}
		case OpInc:
			sb.WriteString("inc " + readAddr(a) + " " + readAddr(uint16(b)))
		default:
			if bs, ok := singleOp[bop]; ok {
				sb.WriteString(bs.Text + " " + readAddr(a) + " " + readAddr(b))
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	c.Pos = oldpos

	sb.WriteString(spaces2 + "+ END " + c.String() + " " + c.Source)
	return sb.String()
}
