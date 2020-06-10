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
	data   []uint32
	pos    posVByte
	source string
}

func newpacket() packet {
	return packet{data: make([]uint32, 0, 1), pos: make(posVByte, 0, 1)}
}

func (b *packet) Clear() {
	b.data = b.data[:0]
}

func (p *packet) Write(buf packet) {
	datalen := len(p.data)
	p.data = append(p.data, buf.data...)
	i := 0
	for i < len(buf.pos) {
		var a, b uint32
		var c uint16
		i, a, b, c = buf.pos.readABC(i)
		p.pos.append(a+uint32(datalen), b, c)
	}
	if p.source == "" {
		p.source = buf.source
	}
}

func (b *packet) WriteRaw(buf []uint32) {
	b.data = append(b.data, buf...)
}

func (b *packet) Write32(v uint32) {
	b.data = append(b.data, v)
}

func (b *packet) WriteOP(op _Opcode, opa, opb uint16) {
	b.data = append(b.data, makeop(op, opa, opb))
}

func (b *packet) WriteJmpOP(op _Opcode, opa uint16, d int) {
	b.data = append(b.data, makejmpop(op, opa, d))
}

func (b *packet) WritePos(p parser.Position) {
	if p.Line == 0 {
		// Debug Code, used to detect a null meta struct
		panicf("null line")
	}
	b.pos.append(uint32(len(b.data)), p.Line, uint16(p.Column))
	if p.Source != "" {
		b.source = p.Source
	}
}

func (b *packet) WriteString(v string) {
	b.Write32(uint32(len(v)))
	b.WriteRaw(u32FromBytes([]byte(v)))
}

func (b *packet) TruncateLast(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *packet) WriteConsts(consts []interface{}) {
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
			buf.WriteByte(3)
			if k {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		default:
			panic("shouldn't happen")
		}
	}
	b.Write32(uint32(buf.Len()))
	b.WriteRaw(u32FromBytes(buf.Bytes()))
}

func crReadConsts(code []uint32, cursor *uint32, n int) []Value {
	res := make([]Value, 0, n)
	buf := crReadBytesLen(code, int(crRead32(code, cursor)), cursor)
	for init := buf; len(buf) > 0; {
		switch buf[0] {
		case 0: // float
			x := math.Float64frombits(binary.BigEndian.Uint64(buf[1:]))
			res = append(res, Num(x))
			buf = buf[9:]
		case 1: // int
			v, n := binary.Varint(buf[1:])
			res = append(res, Num(float64(v)))
			buf = buf[1+n:]
		case 2: // string
			l, n := binary.Uvarint(buf[1:])
			res = append(res, _StrBytes(buf[1+n:1+n+int(l)]))
			buf = buf[1+n+int(l):]
		case 3: // bool
			res = append(res, Bln(buf[1] == 1))
			buf = buf[2:]
		default:
			panicf("invalid const table entry: %v in %x", buf[0], init)
		}
	}
	return res
}

func (b *packet) Len() int {
	return len(b.data)
}

func crRead(data []uint32, len int, cursor *uint32) []uint32 {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func crRead32(data []uint32, cursor *uint32) uint32 {
	*cursor++
	return data[*cursor-1]
}

func crReadBytesLen(data []uint32, length int, cursor *uint32) []byte {
	buf := crRead(data, int((length+3)/4), cursor)
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

func (c *Closure) crPrettify(tab int) string {
	sb := &bytes.Buffer{}
	spaces := strings.Repeat("        ", tab)
	spaces2 := ""
	if tab > 0 {
		spaces2 = strings.Repeat("        ", tab-1) + "+-------"
	}

	sb.WriteString(spaces2 + "+ START " + string(c.source) + "\n")

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
		return fmt.Sprintf("$%d$%d", a>>10, a&0x03ff)
	}

	oldpos := c.Pos
MAIN:
	for {
		bop, a, b := op(crRead32(c.Code, &cursor))
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
			sb.WriteString("$a = popv")
		case OpPushV:
			sb.WriteString("pushv " + readAddr(a))
			if b != 0 {
				sb.WriteString(" cap=" + strconv.Itoa(int(b)))
			}
		case OpPush:
			sb.WriteString("push " + readAddr(a))
		case OpPush2:
			sb.WriteString("push2 " + readAddr(a) + " " + readAddr(b))
		case OpRet:
			sb.WriteString("ret " + readAddr(a))
		case OpYield:
			sb.WriteString("yield " + readAddr(a))
		case OpLambda:
			sb.WriteString("$a = closure:\n")
			cls := crReadClosure(c.Code, &cursor, nil, a, b)
			sb.WriteString(cls.crPrettify(tab + 1))
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
		case OpMakeTable:
			switch a {
			case 1:
				sb.WriteString("make-hash")
			case 2:
				sb.WriteString("make-array")
			case 3:
				sb.WriteString("make-hash-a")
			}
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

	sb.WriteString(spaces2 + "+ END " + string(c.source))
	return sb.String()
}

func crReadClosure(code []uint32, cursor *uint32, env *Env, opa, opb uint16) *Closure {
	opaopb := uint32(opa)<<13 | uint32(opb)
	numParam := byte(opaopb >> 18)
	options := byte(opaopb >> 10)

	consts := crReadConsts(code, cursor, int(uint16(opaopb&0x3ff)))
	src := crReadBytesLen(code, int(crRead32(code, cursor)), cursor)
	pos := crReadBytesLen(code, int(crRead32(code, cursor)), cursor)
	clsCode := crRead(code, int(crRead32(code, cursor)), cursor)

	return &Closure{
		Code:       clsCode,
		ConstTable: consts,
		Env:        env,
		NumParam:   numParam,
		Pos:        pos,
		options:    options,
		source:     src,
	}
}
