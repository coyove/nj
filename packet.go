package potatolang

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

func makeop(op byte, a, b uint16) uint32 {
	// 6 + 13 + 13
	return uint32(op)<<26 + uint32(a&0x1fff)<<13 + uint32(b&0x1fff)
}

func op(x uint32) (op byte, a, b uint16) {
	op = byte(x >> 26)
	a = uint16(x>>13) & 0x1fff
	b = uint16(x) & 0x1fff
	return
}

func btob(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func slice64to8(p []uint32) []byte {
	r := reflect.SliceHeader{}
	r.Cap = cap(p) * 4
	r.Len = len(p) * 4
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

var filler = []byte{0, 0, 0, 0, 0, 0, 0}

func slice8to64(p []byte) []uint32 {
	if m := len(p) % 4; m != 0 {
		p = append(p, filler[:4-m]...)
	}
	r := reflect.SliceHeader{}
	r.Cap = cap(p) / 4
	r.Len = len(p) / 4
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

type packet struct {
	data   []uint32
	pos    []uint32
	source string
}

func newpacket() packet {
	return packet{data: make([]uint32, 0, 1), pos: make([]uint32, 0, 1)}
}

func (b *packet) Clear() {
	b.data = b.data[:0]
}

func (b *packet) Write(buf packet) {
	datalen := len(b.data)
	b.data = append(b.data, buf.data...)
	idx := len(b.pos)
	b.pos = append(b.pos, buf.pos...)
	for i := idx; i < len(b.pos); i += 2 {
		b.pos[i] += uint32(datalen)
	}
	b.source = buf.source
}

func (b *packet) WriteRaw(buf []uint32) { b.data = append(b.data, buf...) }

func (b *packet) Write64(v uint64) { b.data = append(b.data, uint32(v>>32), uint32(v)) }

func (b *packet) Write32(v uint32) { b.data = append(b.data, v) }

func (b *packet) WriteOP(op byte, opa, opb uint16) { b.data = append(b.data, makeop(op, opa, opb)) }

func (b *packet) WritePos(p parser.Meta) {
	b.pos = append(b.pos, uint32(len(b.data)))
	b.pos = append(b.pos, (p.Line<<12)|uint32(p.Column&0x0fff))
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

func (b *packet) WriteConsts(consts []kinfo) {
	// const table struct:
	// all values are placed sequentially
	// for numbers other than MaxUint64, they will be written directly
	// for MaxUint64, it will be written twice
	// for strings, a MaxUint64 will be written first, then the string
	for _, k := range consts {
		if k.ty == Tnumber {
			n := k.value.(float64)
			if math.Float64bits(n) == math.MaxUint64 {
				b.Write64(math.MaxUint64)
				b.Write64(math.MaxUint64)
			} else {
				b.WriteDouble(n)
			}
		} else {
			b.Write64(math.MaxUint64)
			b.WriteString(k.value.(string))
		}
	}
}

func (b *packet) Len() int {
	return len(b.data)
}

func crRead(data []uint32, cursor *uint32, len int) []uint32 {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func crRead32(data []uint32, cursor *uint32) uint32 {
	*cursor++
	return data[*cursor-1]
}

func crRead64(data []uint32, cursor *uint32) uint64 {
	*cursor += 2
	return uint64(data[*cursor-2])<<32 + uint64(data[*cursor-1])
}

func crReadDouble(data []uint32, cursor *uint32) float64 {
	d := crRead64(data, cursor)
	return *(*float64)(unsafe.Pointer(&d))
}

func crReadString(data []uint32, cursor *uint32) string {
	x := crRead32(data, cursor)
	return crReadStringLen(data, int(x), cursor)
}

func crReadStringLen(data []uint32, length int, cursor *uint32) string {
	buf := crRead(data, cursor, int((length+3)/4))
	return string(slice64to8(buf)[:length])
}

var singleOp = map[byte]string{
	OP_ASSERT:   "assert",
	OP_ADD:      "add",
	OP_SUB:      "sub",
	OP_MUL:      "mul",
	OP_DIV:      "div",
	OP_MOD:      "mod",
	OP_EQ:       "eq",
	OP_NEQ:      "neq",
	OP_LESS:     "less",
	OP_LESS_EQ:  "less-eq",
	OP_LEN:      "len",
	OP_COPY:     "copy",
	OP_LOAD:     "load",
	OP_STORE:    "store",
	OP_NOT:      "not",
	OP_BIT_NOT:  "bit-not",
	OP_BIT_AND:  "bit-and",
	OP_BIT_OR:   "bit-or",
	OP_BIT_XOR:  "bit-xor",
	OP_BIT_LSH:  "bit-lsh",
	OP_BIT_RSH:  "bit-rsh",
	OP_BIT_URSH: "bit-ursh",
	OP_TYPEOF:   "typeof",
	OP_SLICE:    "slice",
	OP_POP:      "pop",
}

func crHash(data []uint32) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(slice64to8(data))
	return e.Sum32()
}

func (c *Closure) crPrettify(tab int) string {
	sb := &bytes.Buffer{}
	spaces := strings.Repeat("        ", tab)
	spaces2 := ""
	if tab > 0 {
		spaces2 = strings.Repeat("        ", tab-1) + "+-------"
	}

	sb.WriteString(spaces2 + "+ " + c.source + "\n")

	var cursor uint32
	readAddr := func(a uint16) string {
		if a == regA {
			return "$a"
		}
		return fmt.Sprintf("$%d$%d", a>>10, a&0x03ff)
	}
	readKAddr := func(a uint16) string {
		return fmt.Sprintf("k$%d(%+v)", a, c.consts[a])
	}

	oldpos := c.pos
MAIN:
	for {
		bop, a, b := op(crRead32(c.code, &cursor))
		sb.WriteString(spaces)

		if len(c.pos) > 0 {
			op, line, col := c.pos[0], uint32(c.pos[1]>>12), uint32(c.pos[1]&0x0fff)
			// log.Println(cursor, op, unsafe.Pointer(&pos))
			for cursor > op {
				c.pos = c.pos[2:]
				if len(c.pos) == 0 {
					break
				}
				if op, line, col = c.pos[0], uint32(c.pos[1]>>12), uint32(c.pos[1]&0x0fff); cursor <= op {
					break
				}
			}

			if op == cursor {
				x := fmt.Sprintf("%d:%d", line, col)
				sb.WriteString(fmt.Sprintf("|%-7s %d| ", x, cursor-1))
				c.pos = c.pos[2:]
			} else {
				sb.WriteString(fmt.Sprintf("|        %d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|      . %d| ", cursor-1))
		}

		switch bop {
		case OP_EOB:
			sb.WriteString("end\n")
			break MAIN
		case OP_SET:
			sb.WriteString(readAddr(a) + " = " + readAddr(b))
		case OP_SETK:
			sb.WriteString(readAddr(a) + " = " + readKAddr(uint16(b)))
		case OP_R0, OP_R1, OP_R2, OP_R3:
			sb.WriteString("r" + strconv.Itoa(int(bop-OP_R0)/2) + " = " + readAddr(a))
		case OP_R0K, OP_R1K, OP_R2K, OP_R3K:
			sb.WriteString("r" + strconv.Itoa(int(bop-OP_R0K)/2) + " = " + readKAddr(uint16(a)))
		case OP_PUSH:
			sb.WriteString("push " + readAddr(a))
		case OP_PUSHK:
			sb.WriteString("push " + readKAddr(uint16(a)))
		case OP_RET:
			sb.WriteString("ret " + readAddr(a))
		case OP_RETK:
			sb.WriteString("ret " + readKAddr(uint16(a)))
		case OP_YIELD:
			sb.WriteString("yield " + readAddr(a))
		case OP_YIELDK:
			sb.WriteString("yield " + readKAddr(uint16(a)))
		case OP_LAMBDA:
			sb.WriteString("$a = closure:\n")
			cls := crReadClosure(c.code, &cursor, nil, a, b)
			sb.WriteString(cls.crPrettify(tab + 1))
		case OP_CALL:
			sb.WriteString("call " + readAddr(a))
			if b > 0 {
				sb.WriteString(" -> r" + strconv.Itoa(int(b)-1))
			}
		case OP_JMP:
			pos := int32(b) - 1<<12
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IF, OP_IFNOT:
			addr := readAddr(a)
			pos := int32(b) - 1<<12
			pos2 := strconv.Itoa(int(int32(cursor) + pos))
			if bop == OP_IFNOT {
				sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			} else {
				sb.WriteString("if " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + pos2)
			}
		case OP_RX:
			sb.WriteString("r" + strconv.Itoa(int(a)) + " = r" + strconv.Itoa(int(b)))
		case OP_NOP:
			sb.WriteString("nop")
		case OP_INC:
			sb.WriteString("inc " + readAddr(a) + " " + readKAddr(uint16(b)))
		case OP_MAKEMAP:
			if a == 1 {
				sb.WriteString("make-array")
			} else {
				sb.WriteString("make-map")
			}
		default:
			if bs, ok := singleOp[bop]; ok {
				sb.WriteString(bs)
				if a > 0 {
					sb.WriteString(" -> r" + strconv.Itoa(int(a)-1))
				}
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	c.pos = oldpos

	sb.WriteString(spaces2 + "+ " + c.source)
	return sb.String()
}

func crReadClosure(code []uint32, cursor *uint32, env *Env, opa, opb uint16) *Closure {
	argsCount := byte(opa)
	options := byte(opb)
	constsLen := uint16(crRead32(code, cursor))
	consts := make([]Value, constsLen+1)
	for i := uint16(1); i <= constsLen; i++ {
		x := crRead64(code, cursor)
		if x != math.MaxUint64 {
			consts[i] = NewNumberValue(math.Float64frombits(x))
			continue
		}
		x = crRead64(code, cursor)
		if x == math.MaxUint64 {
			consts[i] = NewNumberValue(math.Float64frombits(x))
			continue
		}
		consts[i] = NewStringValue(crReadStringLen(code, int(x), cursor))
	}

	xlen := crRead64(code, cursor)
	poslen, codelen, srclen := uint32(xlen>>38), uint32(xlen<<26>>38), uint16(xlen<<52>>52)
	src := crReadStringLen(code, int(srclen), cursor)
	pos := crRead(code, cursor, int(poslen))
	buf := crRead(code, cursor, int(codelen))
	cls := NewClosure(buf, consts, env, byte(argsCount))
	cls.pos = pos
	cls.options = options
	cls.source = src
	return cls
}
