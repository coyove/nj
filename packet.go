package potatolang

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

// +---------+----------+----------+
// | op (1b) | opA (3b) | opB (4b) |
// +---------+----------+----------+
// opA is only 24bit long, jmp offset is stored in opB

func makeop(op byte, a, b uint32) uint64 {
	return uint64(op)<<56 + uint64(a&0x00ffffff)<<32 + uint64(b)
}

func op(x uint64) (op byte, a, b uint32) {
	op = byte(x >> 56)
	a = uint32(x>>32) & 0x00ffffff
	b = uint32(x)
	return
}

func makeop2(a, b uint32, c uint16) uint64 {
	return uint64(a&0x00ffffff)<<40 + uint64(b&0x00ffffff)<<16 + uint64(c)
}

func op2(x uint64) (a, b uint32, c uint16) {
	a = uint32(x>>40) & 0x00ffffff
	b = uint32(x>>16) & 0x00ffffff
	c = uint16(x)
	return
}

func btob(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func slice64to8(p []uint64) []byte {
	r := reflect.SliceHeader{}
	r.Cap = cap(p) * 8
	r.Len = len(p) * 8
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

var filler = []byte{0, 0, 0, 0, 0, 0, 0}

func slice8to64(p []byte) []uint64 {
	if m := len(p) % 8; m != 0 {
		p = append(p, filler[:8-m]...)
	}
	r := reflect.SliceHeader{}
	r.Cap = cap(p) / 8
	r.Len = len(p) / 8
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

type packet struct {
	data   []uint64
	pos    []uint64
	source string
}

func newpacket() packet {
	return packet{data: make([]uint64, 0, 1), pos: make([]uint64, 0, 1)}
}

func (b *packet) Clear() {
	b.data = b.data[:0]
}

func (b *packet) Write(buf packet) {
	datalen := len(b.data)
	b.data = append(b.data, buf.data...)
	idx := len(b.pos)
	b.pos = append(b.pos, buf.pos...)
	for i := idx; i < len(b.pos); i++ {
		op, line, col := op2(b.pos[i])
		op += uint32(datalen)
		b.pos[i] = makeop2(op, line, col)
	}
	b.source = buf.source
}

func (b *packet) WriteRaw(buf []uint64) {
	b.data = append(b.data, buf...)
}

func (b *packet) Write64(v uint64) {
	b.data = append(b.data, v)
}

func (b *packet) WriteOP(op byte, opa, opb uint32) {
	b.data = append(b.data, makeop(op, opa, opb))
}

func (b *packet) WritePos(p parser.Position) {
	b.pos = append(b.pos, makeop2(uint32(len(b.data)), uint32(p.Line), uint16(p.Column)))
	if p.Source != "" {
		b.source = p.Source
	}
}

func (b *packet) WriteDouble(v float64) {
	d := *(*uint64)(unsafe.Pointer(&v))
	b.Write64(d)
}

func (b *packet) WriteString(v string) {
	b.Write64(uint64(len(v)))
	b.WriteRaw(slice8to64([]byte(v)))
}

func (b *packet) TruncateLast(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *packet) Len() int {
	return len(b.data)
}

func crRead(data []uint64, cursor *uint32, len int) []uint64 {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func crRead64(data []uint64, cursor *uint32) uint64 {
	*cursor++
	return data[*cursor-1]
}

func crReadDouble(data []uint64, cursor *uint32) float64 {
	d := crRead64(data, cursor)
	return *(*float64)(unsafe.Pointer(&d))
}

func crReadString(data []uint64, cursor *uint32) string {
	x := crRead64(data, cursor)
	buf := crRead(data, cursor, int((x+7)/8))
	return string(slice64to8(buf)[:x])
}

func cruRead64(data uintptr, cursor *uint32) uint64 {
	*cursor++
	return *(*uint64)(unsafe.Pointer(data + uintptr(*cursor-1)*8))
}

// little endian only
// equivalent to: op, opa, opb := op(cruRead64(...))
func cruop(data uintptr, cursor *uint32) (byte, uint32, uint32) {
	addr := uintptr(*cursor) * 8
	*cursor++
	return *(*byte)(unsafe.Pointer(data + addr + 7)),
		*(*uint32)(unsafe.Pointer(data + addr + 4)) & 0x00ffffff,
		*(*uint32)(unsafe.Pointer(data + addr))
}

var singleOp = map[byte]string{
	OP_ADD:     "add",
	OP_SUB:     "sub",
	OP_MUL:     "mul",
	OP_DIV:     "div",
	OP_MOD:     "mod",
	OP_EQ:      "eq",
	OP_NEQ:     "neq",
	OP_LESS:    "less",
	OP_LESS_EQ: "less-eq",
	OP_LEN:     "len",
	OP_DUP:     "dup",
	OP_LOAD:    "load",
	OP_STORE:   "store",
	OP_NOT:     "not",
	OP_BIT_NOT: "bit-not",
	OP_BIT_AND: "bit-and",
	OP_BIT_OR:  "bit-or",
	OP_BIT_XOR: "bit-xor",
	OP_BIT_LSH: "bit-lsh",
	OP_BIT_RSH: "bit-rsh",
	OP_TYPEOF:  "typeof",
	OP_MAKEMAP: "make-map",
}

func crHash(data []uint64) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(slice64to8(data))
	return e.Sum32()
}

func (c *Closure) crPrettify(tab int) string {
	sb := &bytes.Buffer{}
	spaces := strings.Repeat(" ", tab)
	metaprefix := spaces + "|    m    [-] "

	sb.WriteString(metaprefix + "args: " + strconv.Itoa(int(c.argsCount)) + "\n")
	sb.WriteString(metaprefix + "source: " + c.source + "\n")

	if len(c.preArgs) > 0 {
		sb.WriteString(metaprefix + "curried args:" + strconv.Itoa(len(c.preArgs)) + "\n")
	}
	if c.Isset(CLS_YIELDABLE) {
		sb.WriteString(metaprefix + "yieldable: true\n")
	}
	if c.Isset(CLS_HASRECEIVER) {
		sb.WriteString(metaprefix + "receiver: true\n")
	}
	if !c.Isset(CLS_NOENVESCAPE) {
		sb.WriteString(metaprefix + "envescaped: true\n")
	}

	sb.WriteString(metaprefix + fmt.Sprintf("consts: %d\n", len(c.consts)))

	hash := crHash(c.code)
	sb.WriteString(metaprefix + fmt.Sprintf("hash: %08x\n", hash))

	var cursor uint32

	readAddr := func(a uint32) string {
		if a == regA {
			return "$a"
		}
		return fmt.Sprintf("$%d$%d", a>>16, uint16(a))
	}
	readKAddr := func(a uint16) string {
		return fmt.Sprintf("k$%d(%+v)", a, c.consts[a])
	}

	oldpos := c.pos
	lastopcurosor := uint32(0xffffffff)
MAIN:
	for {
		bop, a, b := op(crRead64(c.code, &cursor))
		sb.WriteString(spaces)

		if len(c.pos) > 0 {
			op, line, col := op2(c.pos[0])
			// log.Println(cursor, op, unsafe.Pointer(&pos))
			for cursor > op {
				c.pos = c.pos[1:]
				if len(c.pos) == 0 {
					break
				}
				if op, line, col = op2(c.pos[0]); cursor <= op {
					break
				}
			}

			if op == cursor {
				x := fmt.Sprintf("%d:%d", line, col)
				if cursor == lastopcurosor+1 {
					sb.WriteString(fmt.Sprintf("| %-7s   [%d] ", x, cursor-1))
				} else {
					sb.WriteString(fmt.Sprintf("| %-7s ┴─[%d] ", x, cursor-1))
				}
				c.pos = c.pos[1:]
				lastopcurosor = op
			} else {
				if cursor == lastopcurosor+1 {
					sb.WriteString(fmt.Sprintf("|         ┌─[%d] ", cursor-1))
				} else {
					sb.WriteString(fmt.Sprintf("|         │ [%d] ", cursor-1))
				}
			}
		} else {
			sb.WriteString(fmt.Sprintf("|         . [%d] ", cursor-1))
		}

		switch bop {
		case OP_EOB:
			sb.WriteString("end\n")
			break MAIN
		case OP_SET:
			sb.WriteString(readAddr(a) + " = " + readAddr(b))
		case OP_SETK:
			sb.WriteString(readAddr(a) + " = " + readKAddr(uint16(b)))
		case OP_R0:
			sb.WriteString("r0 = " + readAddr(a))
		case OP_R0K:
			sb.WriteString("r0 = " + readKAddr(uint16(a)))
		case OP_R1:
			sb.WriteString("r1 = " + readAddr(a))
		case OP_R1K:
			sb.WriteString("r1 = " + readKAddr(uint16(a)))
		case OP_R2:
			sb.WriteString("r2 = " + readAddr(a))
		case OP_R2K:
			sb.WriteString("r2 = " + readKAddr(uint16(a)))
		case OP_R3:
			sb.WriteString("r3 = " + readAddr(a))
		case OP_R3K:
			sb.WriteString("r3 = " + readKAddr(uint16(a)))
		case OP_PUSH, OP_PUSHK:
			sb.WriteString("push ")
			switch bop {
			case OP_PUSH:
				sb.WriteString(readAddr(a))
			case OP_PUSHK:
				sb.WriteString(readKAddr(uint16(a)))
			}
		case OP_ASSERT:
			tt := crReadString(c.code, &cursor)
			sb.WriteString(tt)
		case OP_RET:
			sb.WriteString("ret " + readAddr(a))
		case OP_RETK:
			sb.WriteString("ret " + readKAddr(uint16(a)))
		case OP_YIELD:
			sb.WriteString("yield " + readAddr(a))
		case OP_YIELDK:
			sb.WriteString("yield " + readKAddr(uint16(a)))
		case OP_LAMBDA:
			sb.WriteString("$a = closure:\n\n")
			cls := crReadClosure(c.code, &cursor, nil, a, b)
			sb.WriteString(cls.crPrettify(tab + 4))
		case OP_CALL:
			sb.WriteString("call " + readAddr(a))
		case OP_JMP:
			pos := int32(b)
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IF, OP_IFNOT:
			addr := readAddr(a)
			pos := int32(b)
			pos2 := strconv.Itoa(int(int32(cursor) + pos))
			if bop == OP_IFNOT {
				sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			} else {
				sb.WriteString("if " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			}
		case OP_NOP:
			sb.WriteString("nop")
		case OP_INC:
			sb.WriteString("inc " + readAddr(a) + " " + readKAddr(uint16(b)))
		default:
			if bs, ok := singleOp[bop]; ok {
				sb.WriteString(bs)
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	c.pos = oldpos
	return sb.String()
}

func crReadClosure(code []uint64, cursor *uint32, env *Env, opa, opb uint32) *Closure {
	metadata := opb
	argsCount := byte(metadata >> 24)
	options := byte(metadata)
	constsLen := opa
	consts := make([]Value, constsLen+1)
	for i := uint32(1); i <= constsLen; i++ {
		switch crRead64(code, cursor) {
		case Tnumber:
			consts[i] = NewNumberValue(crReadDouble(code, cursor))
		case Tstring:
			consts[i] = NewStringValue(crReadString(code, cursor))
		default:
			panic("shouldn't happen")
		}
	}
	src := crReadString(code, cursor)
	pos := crRead(code, cursor, int(crRead64(code, cursor)))
	buf := crRead(code, cursor, int(crRead64(code, cursor)))
	cls := NewClosure(buf, consts, env, byte(argsCount))
	cls.pos = pos
	cls.options = options
	cls.source = src
	return cls
}
