package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"unsafe"
)

// BytesWriter writes complex values into bytes slice
type BytesWriter struct {
	data []byte
}

const (
	defaultBufferSize = 128
)

func NewBytesWriter() *BytesWriter {
	return &BytesWriter{make([]byte, 0, defaultBufferSize)}
}

func (b *BytesWriter) Dup() *BytesWriter {
	b2 := *b
	return &b2
}

func (b *BytesWriter) Clear() {
	b.data = b.data[:0]
}

func (b *BytesWriter) Bytes() []byte {
	return b.data
}

func (b *BytesWriter) Write(buf []byte) {
	b.data = append(b.data, buf...)
}

func (b *BytesWriter) WriteByte(v byte) {
	b.data = append(b.data, v)
}

func (b *BytesWriter) WriteInt32(v int32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	b.data = append(b.data, buf...)
}

func (b *BytesWriter) WriteInt64(v int64) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	b.data = append(b.data, buf...)
}

func (b *BytesWriter) WriteDouble(v float64) {
	d := *(*int64)(unsafe.Pointer(&v))
	b.WriteInt64(d)
}

func (b *BytesWriter) WriteString(v string) {
	b.WriteInt32(int32(len(v)))
	b.Write([]byte(v))
}

func (b *BytesWriter) TruncateLastBytes(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *BytesWriter) Len() int {
	return len(b.data)
}

func crReadBytes(data []byte, cursor *uint32, len int) []byte {
	*cursor += uint32(len)
	return data[*cursor-uint32(len) : *cursor]
}

func crReadByte(data []byte, cursor *uint32) byte {
	*cursor++
	return data[*cursor-1]
}

func crReadInt32(data []byte, cursor *uint32) int32 {
	*cursor += 4
	return int32(binary.LittleEndian.Uint32(data[*cursor-4 : *cursor]))
}

func crReadInt64(data []byte, cursor *uint32) int64 {
	*cursor += 8
	return int64(binary.LittleEndian.Uint64(data[*cursor-8 : *cursor]))
}

func crReadDouble(data []byte, cursor *uint32) float64 {
	d := crReadInt64(data, cursor)
	return *(*float64)(unsafe.Pointer(&d))
}

func crReadString(data []byte, cursor *uint32) string {
	ln := crReadInt32(data, cursor)
	return string(crReadBytes(data, cursor, int(ln)))
}

var singleOp = map[byte]string{
	OP_ADD:        "add",
	OP_SUB:        "sub",
	OP_MUL:        "mul",
	OP_DIV:        "div",
	OP_MOD:        "mod",
	OP_EQ:         "eq",
	OP_NEQ:        "neq",
	OP_LESS:       "less",
	OP_LESS_EQ:    "less-eq",
	OP_MORE:       "more",
	OP_MORE_EQ:    "more-eq",
	OP_LEN:        "len",
	OP_DUP:        "dup",
	OP_WHO:        "who",
	OP_LOAD:       "load",
	OP_STORE:      "store",
	OP_SAFE_LOAD:  "sload",
	OP_SAFE_STORE: "sstore",
	OP_NOT:        "not",
	OP_AND:        "and",
	OP_OR:         "or",
	OP_XOR:        "xor",
	OP_BIT_NOT:    "bit-not",
	OP_BIT_AND:    "bit-and",
	OP_BIT_OR:     "bit-or",
	OP_BIT_XOR:    "bit-xor",
	OP_BIT_LSH:    "bit-lsh",
	OP_BIT_RSH:    "bit-rsh",
	OP_ERROR:      "error",
}

func crHash(data []byte) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(data)
	return e.Sum32()
}

// Prettify prettifies the code to somehow human readable
func Prettify(code []byte) string {
	return crPrettify(code, 0)
}

func crPrettify(data []byte, tab int) string {
	sb := &bytes.Buffer{}
	pre := strings.Repeat(" ", tab)
	sb.WriteString(pre)
	sb.WriteString(fmt.Sprintf("<%x>\n", crHash(data)))

	var cursor uint32

	readDouble := func() string {
		n := crReadDouble(data, &cursor)
		if float64(int64(n)) == n {
			return strconv.Itoa(int(n))
		}
		return strconv.FormatFloat(n, 'f', 9, 64)
	}
	readString := func() string {
		return strconv.Quote(crReadString(data, &cursor))
	}
	readAddr := func() string {
		if a := crReadInt32(data, &cursor); a == REG_A {
			return "$a"
		} else {
			return fmt.Sprintf("$%d$%d", a>>16, int16(a))
		}
	}

	lastBop := byte(OP_EOB)
MAIN:
	for {
		bop := crReadByte(data, &cursor)

		lastIdx := sb.Len() - 1
		sb.WriteString(pre + "[")
		sb.WriteString(strconv.Itoa(int(cursor) - 1))
		sb.WriteString("] ")
		switch bop {
		case OP_EOB:
			sb.WriteString("end\n")
			break MAIN
		case OP_SET:
			sb.WriteString(readAddr() + " = " + readAddr())
		case OP_SET_NUM:
			sb.WriteString(readAddr() + " = " + readDouble())
		case OP_SET_STR:
			sb.WriteString(readAddr() + " = " + readString())
		case OP_R0:
			sb.WriteString("r0 = " + readAddr())
		case OP_R0_NUM:
			sb.WriteString("r0 = " + readDouble())
		case OP_R0_STR:
			sb.WriteString("r0 = " + readString())
		case OP_R1:
			sb.Truncate(lastIdx)
			sb.WriteString(", r1 = " + readAddr())
		case OP_R1_NUM:
			sb.Truncate(lastIdx)
			sb.WriteString(", r1 = " + readDouble())
		case OP_R1_STR:
			sb.Truncate(lastIdx)
			sb.WriteString(", r1 = " + readString())
		case OP_R2:
			sb.Truncate(lastIdx)
			sb.WriteString(", r2 = " + readAddr())
		case OP_R2_NUM:
			sb.Truncate(lastIdx)
			sb.WriteString(", r2 = " + readDouble())
		case OP_R2_STR:
			sb.Truncate(lastIdx)
			sb.WriteString(", r2 = " + readString())
		case OP_R3:
			sb.Truncate(lastIdx)
			sb.WriteString(", r3 = " + readAddr())
		case OP_R3_NUM:
			sb.Truncate(lastIdx)
			sb.WriteString(", r3 = " + readDouble())
		case OP_R3_STR:
			sb.Truncate(lastIdx)
			sb.WriteString(", r3 = " + readString())
		case OP_PUSH, OP_PUSH_NUM, OP_PUSH_STR:
			if lastBop == OP_PUSH || lastBop == OP_PUSH_NUM || lastBop == OP_PUSH_STR {
				sb.Truncate(lastIdx)
				sb.WriteString(", ")
			} else {
				sb.WriteString("push ")
			}
			switch bop {
			case OP_PUSH:
				sb.WriteString(readAddr())
			case OP_PUSH_NUM:
				sb.WriteString(readDouble())
			case OP_PUSH_STR:
				sb.WriteString(readString())
			}
		case OP_ASSERT:
			sb.Truncate(lastIdx)
			sb.WriteString(" -> " + crReadString(data, &cursor))
		case OP_RET:
			sb.WriteString("ret " + readAddr())
		case OP_RET_NUM:
			sb.WriteString("ret " + readDouble())
		case OP_RET_STR:
			sb.WriteString("ret " + readString())
		case OP_YIELD:
			sb.WriteString("yield " + readAddr())
		case OP_YIELD_NUM:
			sb.WriteString("yield " + readDouble())
		case OP_YIELD_STR:
			sb.WriteString("yield " + readString())
		case OP_LAMBDA:
			sb.WriteString("$a = lambda (\n")
			sb.WriteString(strings.Repeat(" ", tab+4) + "<" + strconv.Itoa(int(crReadInt32(data, &cursor))) + " args>\n")
			if crReadByte(data, &cursor) == 1 {
				sb.WriteString(strings.Repeat(" ", tab+4) + "<yieldable>\n")
			}
			sb.WriteString(crPrettify(crReadBytes(data, &cursor, int(crReadInt32(data, &cursor))), tab+4))
			sb.WriteString(pre + ")")
		case OP_CALL:
			sb.WriteString("call " + readAddr())
		case OP_JMP:
			pos := crReadInt32(data, &cursor)
			pos2 := cursor + uint32(pos)
			sb.WriteString("jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_IF:
			addr := readAddr()
			pos := crReadInt32(data, &cursor)
			pos2 := cursor + uint32(pos)
			sb.WriteString("if not " + addr + " jmp " + strconv.Itoa(int(pos)) + " to " + strconv.Itoa(int(pos2)))
		case OP_NOP:
			sb.WriteString("nop")
		case OP_NIL:
			sb.WriteString("$a = nil")
		case OP_TRUE:
			sb.WriteString("$a = true")
		case OP_FALSE:
			sb.WriteString("$a = false")
		case OP_WHO:
			sb.WriteString("$a = who")
		case OP_STACK:
			sb.WriteString("$a = stack")
		case OP_MAP:
			sb.WriteString("$a = map")
		case OP_LIST:
			sb.WriteString("$a = list")
		default:
			if bs, ok := singleOp[bop]; ok {
				sb.Truncate(lastIdx)
				sb.WriteString(" -> [" + bs + "]")
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
		lastBop = bop
	}

	return sb.String()
}
