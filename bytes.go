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

func btob(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

func slice16to8(p []uint16) []byte {
	r := reflect.SliceHeader{}
	r.Cap = cap(p) * 2
	r.Len = len(p) * 2
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

func slice8to16(p []byte) []uint16 {
	if len(p)%2 != 0 {
		p = append(p, 0)
	}
	r := reflect.SliceHeader{}
	r.Cap = cap(p) / 2
	r.Len = len(p) / 2
	r.Data = (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
	return *(*[]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(&r))))
}

// BytesWriter writes complex values into bytes slice
type BytesWriter struct {
	data []uint16
}

func NewBytesWriter() *BytesWriter {
	const defaultBufferSize = 128
	return &BytesWriter{
		make([]uint16, 0, defaultBufferSize),
	}
}

func (b *BytesWriter) Dup() *BytesWriter {
	b2 := *b
	return &b2
}

func (b *BytesWriter) Clear() {
	b.data = b.data[:0]
}

func (b *BytesWriter) Data() []uint16 {
	return b.data
}

func (b *BytesWriter) Write(buf []uint16) {
	b.data = append(b.data, buf...)
}

func (b *BytesWriter) Write16(v uint16) {
	b.data = append(b.data, v)
}

func (b *BytesWriter) Write32(v uint32) {
	b.data = append(b.data, uint16(v>>16), uint16(v))
}

func (b *BytesWriter) Write64(v uint64) {
	b.data = append(b.data, uint16(v>>48), uint16(v>>32), uint16(v>>16), uint16(v))
}

func (b *BytesWriter) WriteDouble(v float64) {
	d := *(*uint64)(unsafe.Pointer(&v))
	b.Write64(d)
}

func (b *BytesWriter) WriteString(v string) {
	b.Write32(uint32(len(v)))
	b.Write(slice8to16([]byte(v)))
}

func (b *BytesWriter) TruncateLast(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *BytesWriter) Len() int {
	return len(b.data)
}

func crRead(data []uint16, cursor *uint32, len int) []uint16 {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func crRead16(data []uint16, cursor *uint32) uint16 {
	*cursor++
	return data[*cursor-1]
}

func crRead32(data []uint16, cursor *uint32) uint32 {
	*cursor += 2
	return uint32(data[*cursor-2])<<16 + uint32(data[*cursor-1])
}

func crRead64(data []uint16, cursor *uint32) uint64 {
	*cursor += 4
	return uint64(data[*cursor-4])<<48 + uint64(data[*cursor-3])<<32 + uint64(data[*cursor-2])<<16 + uint64(data[*cursor-1])
}

func crReadDouble(data []uint16, cursor *uint32) float64 {
	d := crRead64(data, cursor)
	return *(*float64)(unsafe.Pointer(&d))
}

func crReadString(data []uint16, cursor *uint32) string {
	x := crRead32(data, cursor)
	buf := crRead(data, cursor, int((x+1)/2))
	return string(slice16to8(buf)[:x])
}

var singleOp = map[uint16]string{
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

func crHash(data []uint16) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(slice16to8(data))
	return e.Sum32()
}

func crPrettifyLambda(args, curry int, y, e bool, code []uint16, consts []Value, tab int) string {
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
	for i, k := range consts {
		sb.WriteString(spaces + fmt.Sprintf("<k$%d: %+v>\n", i, k))
	}
	sb.WriteString(crPrettify(code, consts, tab))
	return sb.String()
}

func crPrettify(data []uint16, consts []Value, tab int) string {
	sb := &bytes.Buffer{}
	pre := strings.Repeat(" ", tab)
	hash := crHash(data)
	sb.WriteString(pre)
	sb.WriteString(fmt.Sprintf("<%x>\n", hash))

	var cursor uint32

	readAddr := func() string {
		if a := crRead32(data, &cursor); a == regA {
			return "$a"
		} else {
			return fmt.Sprintf("$%d$%d", a>>16, uint16(a))
		}
	}
	readKAddr := func() string {
		a := crRead16(data, &cursor)
		return fmt.Sprintf("k$%d <%+v>", a, consts[a])
	}

	lastBop := uint16(OP_EOB)
MAIN:
	for {
		bop := crRead16(data, &cursor)

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
			sb.WriteString(readAddr() + " = " + readAddr())
		case OP_SETK:
			sb.WriteString(readAddr() + " = " + readKAddr())
		case OP_R0:
			sb.WriteString("r0 = " + readAddr())
		case OP_R0K:
			sb.WriteString("r0 = " + readKAddr())
		case OP_R1:
			sb.Truncate(lastIdx)
			sb.WriteString(", r1 = " + readAddr())
		case OP_R1K:
			sb.Truncate(lastIdx)
			sb.WriteString(", r1 = " + readKAddr())
		case OP_R2:
			sb.Truncate(lastIdx)
			sb.WriteString(", r2 = " + readAddr())
		case OP_R2K:
			sb.Truncate(lastIdx)
			sb.WriteString(", r2 = " + readKAddr())
		case OP_R3:
			sb.Truncate(lastIdx)
			sb.WriteString(", r3 = " + readAddr())
		case OP_R3K:
			sb.Truncate(lastIdx)
			sb.WriteString(", r3 = " + readKAddr())
		case OP_PUSH, OP_PUSHK:
			if lastBop == OP_PUSH || lastBop == OP_PUSHK {
				sb.Truncate(lastIdx)
				sb.WriteString(", ")
			} else {
				sb.WriteString("push ")
			}
			switch bop {
			case OP_PUSH:
				sb.WriteString(readAddr())
			case OP_PUSHK:
				sb.WriteString(readKAddr())
			}
		case OP_ASSERT:
			tt := crReadString(data, &cursor)
			sb.WriteString(tt)
		case OP_RET:
			sb.WriteString("ret " + readAddr())
		case OP_RETK:
			sb.WriteString("ret " + readKAddr())
		case OP_YIELD:
			sb.WriteString("yield " + readAddr())
		case OP_YIELDK:
			sb.WriteString("yield " + readKAddr())
		case OP_LAMBDA:
			sb.WriteString("$a = lambda (\n")
			argsCount := crRead16(data, &cursor)
			yieldable := crRead16(data, &cursor) == 1
			errorable := crRead16(data, &cursor) == 1
			constsLen := crRead16(data, &cursor)
			consts := make([]Value, constsLen)
			for i := uint16(0); i < constsLen; i++ {
				switch crRead16(data, &cursor) {
				case Tnumber:
					consts[i] = NewNumberValue(crReadDouble(data, &cursor))
				case Tstring:
					consts[i] = NewStringValue(crReadString(data, &cursor))
				default:
					panic("shouldn't happen")
				}
			}
			buf := crRead(data, &cursor, int(crRead32(data, &cursor)))

			sb.WriteString(crPrettifyLambda(int(argsCount), 0, yieldable, errorable, buf, consts, tab+4))
			sb.WriteString(pre + ")")
		case OP_CALL:
			sb.WriteString("call " + readAddr())
		case OP_JMP:
			pos := int32(crRead32(data, &cursor))
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IFNOT:
			addr := readAddr()
			pos := int32(crRead32(data, &cursor))
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IF:
			addr := readAddr()
			pos := int32(crRead32(data, &cursor))
			pos2 := uint32(int32(cursor) + pos)
			sb.WriteString("if " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_NOP:
			sb.WriteString("nop")
		case OP_INC:
			sb.WriteString("inc " + readAddr() + " " + readKAddr())
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
