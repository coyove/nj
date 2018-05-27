package potatolang

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func makeop(op byte, a, b uint32) uint64 {
	return uint64(op)<<56 + uint64(a&0x00ffffff)<<32 + uint64(b)
}

func op(x uint64) (op byte, a, b uint32) {
	op = byte(x >> 56)
	a = uint32(x>>32) & 0x00ffffff
	b = uint32(x)
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

// BytesWriter writes complex values into bytes slice
type BytesWriter struct {
	data []uint64
}

func NewBytesWriter() *BytesWriter {
	const defaultBufferSize = 128
	return &BytesWriter{
		make([]uint64, 0, defaultBufferSize),
	}
}

func (b *BytesWriter) Dup() *BytesWriter {
	b2 := *b
	return &b2
}

func (b *BytesWriter) Clear() {
	b.data = b.data[:0]
}

func (b *BytesWriter) Data() []uint64 {
	return b.data
}

func (b *BytesWriter) Write(buf []uint64) {
	b.data = append(b.data, buf...)
}

func (b *BytesWriter) Write64(v uint64) {
	b.data = append(b.data, v)
}

func (b *BytesWriter) WriteDouble(v float64) {
	d := *(*uint64)(unsafe.Pointer(&v))
	b.Write64(d)
}

func (b *BytesWriter) WriteString(v string) {
	b.Write64(uint64(len(v)))
	b.Write(slice8to64([]byte(v)))
}

func (b *BytesWriter) TruncateLast(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *BytesWriter) Len() int {
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
	OP_LESS_EQ: "lesseq",
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
	OP_ERROR:   "error",
	OP_TYPEOF:  "typeof",
	OP_NIL:     "nil",
	OP_TRUE:    "true",
	OP_FALSE:   "false",
	OP_WHO:     "who",
	OP_MAP:     "map",
	OP_LIST:    "list",
}

func crHash(data []uint64) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(slice64to8(data))
	return e.Sum32()
}

func crPrettifyLambda(args, curry int, y, e, esc bool, code []uint64, consts []Value, tab int) string {
	sb := &bytes.Buffer{}
	spaces := strings.Repeat(" ", tab)
	sb.WriteString(spaces + "<args: " + strconv.Itoa(args) + ">\n")
	if curry > 0 {
		sb.WriteString(spaces + "<curry: " + strconv.Itoa(curry) + ">\n")
	}
	if y {
		sb.WriteString(spaces + "<yieldable>\n")
	}
	if e {
		sb.WriteString(spaces + "<errorable>\n")
	}
	if esc {
		sb.WriteString(spaces + "<env escaped>\n")
	}
	for i, k := range consts {
		sb.WriteString(spaces + fmt.Sprintf("<k$%d: %+v>\n", i, k))
	}
	sb.WriteString(crPrettify(code, consts, tab))
	return sb.String()
}

func crPrettify(data []uint64, consts []Value, tab int) string {

	sb := &bytes.Buffer{}
	pre := strings.Repeat(" ", tab)
	hash := crHash(data)
	sb.WriteString(pre)
	sb.WriteString(fmt.Sprintf("<%08x>\n", hash))

	var cursor uint32

	readAddr := func(a uint32) string {
		if a == regA {
			return "$a"
		}
		return fmt.Sprintf("$%d$%d", a>>16, uint16(a))
	}
	readKAddr := func(a uint16) string {
		return fmt.Sprintf("k$%d <%+v>", a, consts[a])
	}

	lastBop := byte(OP_EOB)
MAIN:
	for {
		bop, a, b := op(crRead64(data, &cursor))

		lastIdx := sb.Len() - 1
		sb.WriteString(pre + "[")
		sb.WriteString(strconv.Itoa(int(cursor) - 1))
		sb.WriteString("] ")
		switch bop {
		case OP_LINE:
			sb.WriteString(fmt.Sprintf("---- <%x> %s", hash, crReadString(data, &cursor)))
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
			// sb.Truncate(lastIdx)
			sb.WriteString("r1 = " + readAddr(a))
		case OP_R1K:
			// sb.Truncate(lastIdx)
			sb.WriteString("r1 = " + readKAddr(uint16(a)))
		case OP_R2:
			// sb.Truncate(lastIdx)
			sb.WriteString("r2 = " + readAddr(a))
		case OP_R2K:
			// sb.Truncate(lastIdx)
			sb.WriteString("r2 = " + readKAddr(uint16(a)))
		case OP_R3:
			// sb.Truncate(lastIdx)
			sb.WriteString("r3 = " + readAddr(a))
		case OP_R3K:
			// sb.Truncate(lastIdx)
			sb.WriteString("r3 = " + readKAddr(uint16(a)))
		case OP_PUSH, OP_PUSHK:
			if lastBop == OP_PUSH || lastBop == OP_PUSHK {
				sb.Truncate(lastIdx)
				sb.WriteString(", ")
			} else {
				sb.WriteString("push ")
			}
			switch bop {
			case OP_PUSH:
				sb.WriteString(readAddr(a))
			case OP_PUSHK:
				sb.WriteString(readKAddr(uint16(a)))
			}
		case OP_ASSERT:
			tt := crReadString(data, &cursor)
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
			sb.WriteString("$a = lambda (\n")
			argsCount := byte(b >> 24)
			yieldable := byte(b>>16) == 1
			errorable := byte(b>>8) == 1
			noenvescape := byte(b) == 1
			constsLen := a
			consts := make([]Value, constsLen)
			for i := uint32(0); i < constsLen; i++ {
				switch crRead64(data, &cursor) {
				case Tnumber:
					consts[i] = NewNumberValue(crReadDouble(data, &cursor))
				case Tstring:
					consts[i] = NewStringValue(crReadString(data, &cursor))
				default:
					panic("shouldn't happen")
				}
			}
			buf := crRead(data, &cursor, int(crRead64(data, &cursor)))
			sb.WriteString(crPrettifyLambda(int(argsCount), 0, yieldable, errorable, !noenvescape, buf, consts, tab+4))
			sb.WriteString(pre + ")")
		case OP_CALL:
			sb.WriteString("call " + readAddr(a))
		case OP_JMP:
			pos := int32(b)
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IFNOT:
			addr := readAddr(a)
			pos := int32(b)
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IF:
			addr := readAddr(a)
			pos := int32(b)
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("if " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
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
		lastBop = bop
	}

	return sb.String()
}
