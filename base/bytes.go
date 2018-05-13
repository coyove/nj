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

type BytesReader struct {
	data []byte
}

const (
	defaultBufferSize = 128
)

func NewBytesReader(data []byte) *BytesReader {
	return &BytesReader{
		data: data,
	}
}

func NewBytesBuffer() *BytesReader {
	return NewBytesReader(make([]byte, 0, defaultBufferSize))
}

func (b *BytesReader) Dup() *BytesReader {
	b2 := *b
	return &b2
}

func (b *BytesReader) Clear() {
	b.data = b.data[:0]
}

func (b *BytesReader) Bytes() []byte {
	return b.data
}

func (b *BytesReader) Write(buf []byte) {
	b.data = append(b.data, buf...)
}

func (b *BytesReader) WriteByte(v byte) {
	b.data = append(b.data, v)
}

func (b *BytesReader) WriteInt32(v int32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	b.data = append(b.data, buf...)
}

func (b *BytesReader) WriteInt64(v int64) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	b.data = append(b.data, buf...)
}

func (b *BytesReader) WriteDouble(v float64) {
	d := *(*int64)(unsafe.Pointer(&v))
	b.WriteInt64(d)
}

func (b *BytesReader) WriteString(v string) {
	b.WriteInt32(int32(len(v)))
	b.Write([]byte(v))
}

func (b *BytesReader) Truncate(n int) {
	if len(b.data) > n {
		b.data = b.data[:len(b.data)-n]
	}
}

func (b *BytesReader) Len() int {
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
	OP_LIST:       "list",
	OP_LEN:        "len",
	OP_DUP:        "dup",
	OP_WHO:        "who",
	OP_VARARGS:    "varargs",
	OP_MAP:        "map",
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
	OP_NIL:        "nil",
	OP_TRUE:       "true",
	OP_FALSE:      "false",
}

func crHash(data []byte) uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(data)
	return e.Sum32()
}

func Prettify(data []byte) string {
	return crPrettify(data, 0)
}

func crPrettify(data []byte, tab int) string {
	sb := &bytes.Buffer{}
	pre := strings.Repeat(" ", tab)
	sb.WriteString(pre)
	sb.WriteString(fmt.Sprintf("<%x>\n", crHash(data)))

	var cursor uint32
	xy := func(in int32) string {
		if in == REG_A {
			return "$a"
		}
		return fmt.Sprintf("$%d$%d", in>>16, int16(in))
	}

	readDouble := func() string {
		return strconv.FormatFloat(crReadDouble(data, &cursor), 'f', 9, 64)
	}

	for {
		bop := crReadByte(data, &cursor)
		if bop == OP_EOB {
			break
		}

		sb.WriteString(pre + "[")
		sb.WriteString(strconv.Itoa(int(cursor) - 1))
		sb.WriteString("] ")
		switch bop {
		case OP_SET:
			sb.WriteString("set " + xy(crReadInt32(data, &cursor)) + " " + xy(crReadInt32(data, &cursor)))
		case OP_SET_NUM:
			sb.WriteString("seti " + xy(crReadInt32(data, &cursor)) + " " + readDouble())
		case OP_SET_STR:
			sb.WriteString("sets " + xy(crReadInt32(data, &cursor)) + " " + crReadString(data, &cursor))
		case OP_R0:
			sb.WriteString("r0 " + xy(crReadInt32(data, &cursor)))
		case OP_R0_NUM:
			sb.WriteString("r0n " + readDouble())
		case OP_R0_STR:
			sb.WriteString("r0s " + crReadString(data, &cursor))
		case OP_R1:
			sb.WriteString("r1 " + xy(crReadInt32(data, &cursor)))
		case OP_R1_NUM:
			sb.WriteString("r1n " + readDouble())
		case OP_R1_STR:
			sb.WriteString("r1s " + crReadString(data, &cursor))
		case OP_R2:
			sb.WriteString("r2 " + xy(crReadInt32(data, &cursor)))
		case OP_R2_NUM:
			sb.WriteString("r2n " + readDouble())
		case OP_R2_STR:
			sb.WriteString("r2s " + crReadString(data, &cursor))
		case OP_R3:
			sb.WriteString("r3 " + xy(crReadInt32(data, &cursor)))
		case OP_R3_NUM:
			sb.WriteString("r3n " + readDouble())
		case OP_R3_STR:
			sb.WriteString("r3s " + crReadString(data, &cursor))
		case OP_PUSH:
			sb.WriteString("push " + xy(crReadInt32(data, &cursor)))
		case OP_PUSH_NUM:
			sb.WriteString("pushi " + readDouble())
		case OP_PUSH_STR:
			sb.WriteString("pushs " + crReadString(data, &cursor))
		case OP_ASSERT:
			sb.WriteString("assert " + crReadString(data, &cursor))
		case OP_RET:
			sb.WriteString("ret " + xy(crReadInt32(data, &cursor)))
		case OP_RET_NUM:
			sb.WriteString("reti " + readDouble())
		case OP_RET_STR:
			sb.WriteString("rets " + crReadString(data, &cursor))
		case OP_YIELD:
			sb.WriteString("yield " + xy(crReadInt32(data, &cursor)))
		case OP_YIELD_NUM:
			sb.WriteString("yieldi " + readDouble())
		case OP_YIELD_STR:
			sb.WriteString("yields " + crReadString(data, &cursor))
		case OP_LAMBDA:
			sb.WriteString("lambda ")
			sb.WriteString(strconv.Itoa(int(crReadInt32(data, &cursor))))
			sb.WriteString(" (\n")
			sb.WriteString(crPrettify(crReadBytes(data, &cursor, int(crReadInt32(data, &cursor))), tab+4))
			sb.WriteString(pre + ")")

		case OP_CALL:
			sb.WriteString("call " + xy(crReadInt32(data, &cursor)))

		case OP_JMP:
			sb.WriteString("jmp " + strconv.Itoa(int(crReadInt32(data, &cursor))))

		case OP_IF:
			sb.WriteString("if " + xy(crReadInt32(data, &cursor)) + " " + strconv.Itoa(int(crReadInt32(data, &cursor))))
		case OP_NOP:
			sb.WriteString("nop")
		case OP_LIB_CALL:
			sb.WriteString("libcall " + strconv.Itoa(int(uint32(crReadInt32(data, &cursor)))))
		case OP_LIB_CALL_EX:
			sb.WriteString("libcallex " + strconv.Itoa(int(uint32(crReadInt32(data, &cursor)))))
		default:
			if bs, ok := singleOp[bop]; ok {
				sb.WriteString(bs)
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
