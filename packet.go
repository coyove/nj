package potatolang

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

func makeop(op _Opcode, a, b uint16) uint32 {
	// 6 + 13 + 13
	return uint32(op)<<26 + uint32(a&0x1fff)<<13 + uint32(b&0x1fff)
}

func makejmpop(op _Opcode, a uint16, dist int) uint32 {
	if dist < -(1<<12) || dist >= 1<<12 {
		panic("long jump")
	}
	// 6 + 13 + 13
	b := uint16(dist + 1<<12)
	return uint32(op)<<26 + uint32(a&0x1fff)<<13 + uint32(b&0x1fff)
}

func op(x uint32) (op _Opcode, a, b uint16) {
	op = _Opcode(x >> 26)
	a = uint16(x>>13) & 0x1fff
	b = uint16(x) & 0x1fff
	return
}

func u32Bytes(p []uint32) []byte {
	r := reflect.SliceHeader{}
	r.Cap = cap(p) * 4
	r.Len = len(p) * 4
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

func u32FromBytes(p []byte) []uint32 {
	if m := len(p) % 4; m != 0 {
		p = append(p, 0, 0, 0, 0)
		p = p[:len(p)+1-m]
	}
	r := reflect.SliceHeader{}
	r.Cap = cap(p) / 4
	r.Len = len(p) / 4
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

type posVByte []byte

func (p *posVByte) append(idx uint32, line uint32, col uint16) {
	trunc := func(v uint32) (w byte, buf [4]byte) {
		if v < 256 {
			w = 0
		} else if v < 256*256 {
			w = 1
		} else if v < 256*256*256 {
			w = 2
		} else {
			w = 3
		}
		binary.LittleEndian.PutUint32(buf[:], v)
		return
	}
	iw, ib := trunc(idx)
	lw, lb := trunc(line)
	x := iw<<6 + lw<<4
	cw := 0
	col--
	if col < 14 {
		x |= byte(col)
	} else if col < 256 {
		cw = 1
		x |= 0xe
	} else {
		cw = 2
		x |= 0xf
	}

	*p = append(*p, x)
	*p = append(*p, ib[:iw+1]...)
	*p = append(*p, lb[:lw+1]...)
	if cw == 1 {
		*p = append(*p, byte(col))
	} else if cw == 2 {
		*p = append(*p, byte(col), byte(col>>8))
	}
}

func (p posVByte) readABC(i int) (next int, a, b uint32, c uint16) {
	x := p[i]

	i++
	buf := [4]byte{}
	copy(buf[:], p[i:i+1+int(x>>6)])
	a = binary.LittleEndian.Uint32(buf[:])

	i = i + 1 + int(x>>6)
	buf = [4]byte{}
	copy(buf[:], p[i:i+1+int(x<<2>>6)])
	b = binary.LittleEndian.Uint32(buf[:])

	i = i + 1 + int(x<<2>>6)
	x &= 0xf
	if x < 14 {
		c = uint16(x + 1)
		next = i
	} else {
		buf = [4]byte{}
		copy(buf[:], p[i:i+int(x-13)])
		c = binary.LittleEndian.Uint16(buf[:2]) + 1
		next = i + int(x-13)
	}
	return
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
		var a, b_ uint32
		var c uint16
		i, a, b_, c = buf.Pos.readABC(i)
		b.Pos.append(a+uint32(datalen), b_, c)
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

func (b *packet) writeOP(op _Opcode, opa, opb uint16) {
	b.Code = append(b.Code, makeop(op, opa, opb))
}

func (b *packet) writeJmpOP(op _Opcode, opa uint16, d int) {
	b.Code = append(b.Code, makejmpop(op, opa, d))
}

func (b *packet) writePos(p parser.Position) {
	if p.Line == 0 {
		// Debug Code, used to detect a null meta struct
		panicf("null line")
	}
	b.Pos.append(uint32(len(b.Code)), p.Line, uint16(p.Column))
	if p.Source != "" {
		b.Source = p.Source
	}
}

func (b *packet) truncateLast() {
	if len(b.Code) > 0 {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (b *packet) writeConstTable(consts []interface{}) {
	buf := bytes.Buffer{}
	for _, k := range consts {
		switch k := k.(type) {
		case float64:
			if float64(int64(k)) == k {
				buf.WriteByte(1)
				buf.WriteString("0000000000")
				buf.Truncate(buf.Len() - 10 + binary.PutVarint(buf.Bytes()[buf.Len()-10:], int64(k)))
			} else {
				buf.WriteByte(0)
				binary.Write(&buf, binary.BigEndian, math.Float64bits(k))
			}
		case string:
			buf.WriteByte(2)
			buf.WriteString("0000000000")
			buf.Truncate(buf.Len() - 10 + binary.PutUvarint(buf.Bytes()[buf.Len()-10:], uint64(len(k))))
			buf.WriteString(k)
		case bool:
			if k {
				buf.WriteByte(3)
			} else {
				buf.WriteByte(4)
			}
		default:
			panic("shouldn't happen")
		}
	}
	b.write32(uint32(buf.Len()))
	b.writeBytes(u32FromBytes(buf.Bytes()))
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

func pkReadBytes(data []uint32, length int, cursor *uint32) []byte {
	buf := pkRead(data, int((length+3)/4), cursor)
	return u32Bytes(buf)[:length]
}

var singleOp = map[_Opcode]parser.Symbol{
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
MAIN:
	for {
		bop, a, b := op(pkRead32(c.Code, &cursor))
		sb.WriteString(spaces)

		if len(c.Pos) > 0 {
			next, op, line, col := c.Pos.readABC(0)
			// log.Println(cursor, op, unsafe.Pointer(&Pos))
			for cursor > op {
				c.Pos = c.Pos[next:]
				if len(c.Pos) == 0 {
					break
				}
				if next, op, line, col = c.Pos.readABC(0); cursor <= op {
					break
				}
			}

			if op == cursor {
				x := fmt.Sprintf("%d:%d", line, col)
				sb.WriteString(fmt.Sprintf("|%-7s %d| ", x, cursor-1))
				c.Pos = c.Pos[next:]
			} else {
				sb.WriteString(fmt.Sprintf("|        %d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|      . %d| ", cursor-1))
		}

		switch bop {
		case OpEOB:
			sb.WriteString("end\n")
			break MAIN
		case OpSet:
			sb.WriteString(readAddr(a) + " = " + readAddr(b))
		case OpPopV:
			switch a {
			case 0:
				sb.WriteString("$a = popv-last-and-clear-rest")
			case 1:
				sb.WriteString("$a = popv")
			case 2:
				sb.WriteString("$a = popv-all")
			case 3:
				sb.WriteString("popv-clear")
			case 4:
				sb.WriteString("$a = popv-all-with-a")
			}
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
		case OpLambda:
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
