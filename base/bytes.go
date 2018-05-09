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
	data   []byte
	cursor int
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

func (b *BytesReader) WriteInt32AtIndex(start int, v int32) {
	binary.LittleEndian.PutUint32(b.data[start:], uint32(v))
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

func (b *BytesReader) SetCursor(c int) {
	b.cursor = c
}

func (b *BytesReader) GetCursor() int {
	return b.cursor
}

func (b *BytesReader) Read(len int) []byte {
	b.cursor += len
	return b.data[b.cursor-len : b.cursor]
}

func (b *BytesReader) ReadByte() byte {
	b.cursor++
	return b.data[b.cursor-1]
}

func (b *BytesReader) ReadInt32() int32 {
	b.cursor += 4
	return int32(binary.LittleEndian.Uint32(b.data[b.cursor-4 : b.cursor]))
}

func (b *BytesReader) ReadInt64() int64 {
	b.cursor += 8
	return int64(binary.LittleEndian.Uint64(b.data[b.cursor-8 : b.cursor]))
}

func (b *BytesReader) ReadDouble() float64 {
	d := b.ReadInt64()
	return *(*float64)(unsafe.Pointer(&d))
}

func (b *BytesReader) ReadString() string {
	ln := b.ReadInt32()
	return string(b.Read(int(ln)))
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
	OP_MAP:        "array",
	OP_LOAD:       "load",
	OP_STORE:      "store",
	OP_RLOAD:      "rload",
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
	OP_EXPAND:     "expand",
	OP_NIL:        "nil",
	OP_TRUE:       "true",
	OP_FALSE:      "false",
	OP_BYTES:      "bytes",
}

func (b *BytesReader) Hash() uint32 {
	e := crc32.New(crc32.IEEETable)
	e.Write(b.Bytes())
	return e.Sum32()
}

func (b *BytesReader) Prettify(tab int) string {
	sb := &bytes.Buffer{}
	pre := strings.Repeat(" ", tab)
	sb.WriteString(pre)
	sb.WriteString(fmt.Sprintf("<%x>\n", b.Hash()))

	xy := func(in int32) string {
		if in == REG_A {
			return "$a"
		}
		return fmt.Sprintf("$%d$%d", in>>16, int16(in))
	}

	readDouble := func() string {
		return strconv.FormatFloat(b.ReadDouble(), 'f', 9, 64)
	}

	for {
		bop := b.ReadByte()
		if bop == OP_EOB {
			break
		}

		sb.WriteString(pre + "[")
		sb.WriteString(strconv.Itoa(b.GetCursor() - 1))
		sb.WriteString("] ")
		switch bop {
		case OP_SET:
			sb.WriteString("set " + xy(b.ReadInt32()) + " " + xy(b.ReadInt32()))

		case OP_SET_NUM:
			sb.WriteString("seti " + xy(b.ReadInt32()) + " " + readDouble())

		case OP_SET_STR:
			sb.WriteString("sets " + xy(b.ReadInt32()) + " " + b.ReadString())
		case OP_R0:
			sb.WriteString("r0 " + xy(b.ReadInt32()))
		case OP_R0_NUM:
			sb.WriteString("r0n " + strconv.FormatFloat(b.ReadDouble(), 'f', 9, 64))
		case OP_R0_STR:
			sb.WriteString("r0s " + b.ReadString())
		case OP_R1:
			sb.WriteString("r1 " + xy(b.ReadInt32()))
		case OP_R1_NUM:
			sb.WriteString("r1n " + strconv.FormatFloat(b.ReadDouble(), 'f', 9, 64))
		case OP_R1_STR:
			sb.WriteString("r1s " + b.ReadString())
		case OP_R2:
			sb.WriteString("r2 " + xy(b.ReadInt32()))
		case OP_R2_NUM:
			sb.WriteString("r2n " + strconv.FormatFloat(b.ReadDouble(), 'f', 9, 64))
		case OP_R2_STR:
			sb.WriteString("r2s " + b.ReadString())
		case OP_R3:
			sb.WriteString("r3 " + xy(b.ReadInt32()))
		case OP_R3_NUM:
			sb.WriteString("r3n " + strconv.FormatFloat(b.ReadDouble(), 'f', 9, 64))
		case OP_R3_STR:
			sb.WriteString("r3s " + b.ReadString())
		case OP_PUSH:
			sb.WriteString("push " + xy(b.ReadInt32()))
		case OP_PUSH_NUM:
			sb.WriteString("pushi " + readDouble())
		case OP_PUSH_STR:
			sb.WriteString("pushs " + b.ReadString())
		case OP_ASSERT:
			sb.WriteString("assert " + b.ReadString())
		case OP_RET:
			sb.WriteString("ret " + xy(b.ReadInt32()))
		case OP_RET_NUM:
			sb.WriteString("reti " + readDouble())
		case OP_RET_STR:
			sb.WriteString("rets " + b.ReadString())
		case OP_LAMBDA:
			sb.WriteString("lambda ")
			sb.WriteString(strconv.Itoa(int(b.ReadInt32())))
			sb.WriteString(" (\n")
			sb.WriteString(NewBytesReader(b.Read(int(b.ReadInt32()))).Prettify(tab + 4))
			sb.WriteString(pre + ")")

		case OP_CALL:
			sb.WriteString("call " + xy(b.ReadInt32()))

		case OP_JMP:
			sb.WriteString("jmp " + strconv.Itoa(int(b.ReadInt32())))

		case OP_IF:
			sb.WriteString("if " + xy(b.ReadInt32()) + " " + strconv.Itoa(int(b.ReadInt32())))
		case OP_RSTORE:
			sb.WriteString("rstore " + xy(b.ReadInt32()) + " " + strconv.Itoa(int(b.ReadByte())) + " r0")
		case OP_NOP:
			sb.WriteString("nop")
		case OP_LIB_CALL:
			sb.WriteString("libcall " + strconv.Itoa(int(uint32(b.ReadInt32()))))
		case OP_LIB_CALL_EX:
			sb.WriteString("libcallex " + strconv.Itoa(int(uint32(b.ReadInt32()))))
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
