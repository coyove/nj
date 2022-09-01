package typ

import (
	"bytes"
	"math/bits"
)

type ValueType byte

const (
	Nil    ValueType = 0
	Bool   ValueType = 1
	Number ValueType = 3
	String ValueType = 7
	Object ValueType = 15
	Native ValueType = 31
)

func (t ValueType) String() string {
	if t > Native {
		return "?"
	}
	return [...]string{
		"nil", "bool", "?", "number",
		"?", "?", "?", "string",
		"?", "?", "?", "?",
		"?", "?", "?", "object",
		"?", "?", "?", "?",
		"?", "?", "?", "?",
		"?", "?", "?", "?",
		"?", "?", "?", "native"}[t]
}

const (
	RegA          uint16 = 0xffff
	RegPhantom    uint16 = 0xfffe
	RegLocalMask         = 0x7fff
	RegGlobalFlag        = 0x8000
	RegMaxAddress        = 0x7f00
)

type shape uint64

var Shape shape = 0xF

// F: start
// E:
// D: array repeat
// C: array end
// B: array start
// A: native or nil
// 9: native
// 8: object or nil
// 7: object
// 6: string or nil
// 5: string
// 4: number or nil
// 3: number
// 2: bool or nil
// 1: bool
// 0: any

func (s shape) check() {
	if s>>60 == 0xF {
		panic("shape too complex")
	}
}

func (s shape) Any() shape         { s.check(); s = s<<4 | 0; return s }
func (s shape) Bool() shape        { s.check(); s = s<<4 | 1; return s }
func (s shape) BoolOrNil() shape   { s.check(); s = s<<4 | 2; return s }
func (s shape) Num() shape         { s.check(); s = s<<4 | 3; return s }
func (s shape) NumOrNil() shape    { s.check(); s = s<<4 | 4; return s }
func (s shape) Str() shape         { s.check(); s = s<<4 | 5; return s }
func (s shape) StrOrNil() shape    { s.check(); s = s<<4 | 6; return s }
func (s shape) Object() shape      { s.check(); s = s<<4 | 7; return s }
func (s shape) ObjectOrNil() shape { s.check(); s = s<<4 | 8; return s }
func (s shape) Native() shape      { s.check(); s = s<<4 | 9; return s }
func (s shape) NativeOrNil() shape { s.check(); s = s<<4 | 10; return s }
func (s shape) Repeat() shape      { s.check(); s = s<<4 | 13; return s }

func (s shape) Array(s2 shape) shape {
	s.check()
	s = s<<4 | 11
	for i := bits.LeadingZeros64(uint64(s2)) + 4; i <= 64; i += 4 {
		s.check()
		s = s<<4 | (s2>>(64-i))&0xF
	}
	s.check()
	s = s<<4 | 12
	return s
}

func (s shape) String() string {
	buf := bytes.NewBufferString("(")
	trunc := func() {
		if b := buf.Bytes(); len(b) > 0 && b[len(b)-1] == ',' {
			buf.Truncate(buf.Len() - 1)
		}
	}
	for i := bits.LeadingZeros64(uint64(s)) + 4; i <= 64; i += 4 {
		switch (s >> (64 - i)) & 0xF {
		case 0xD: // array repeat
			trunc()
			buf.WriteString("...")
		case 0xC: // array end
			trunc()
			buf.WriteString(")")
		case 0xB: // array start
			buf.WriteString("(")
		case 0xA:
			buf.WriteString("native/nil,")
		case 0x9:
			buf.WriteString("native,")
		case 0x8:
			buf.WriteString("object/nil,")
		case 0x7:
			buf.WriteString("object,")
		case 0x6:
			buf.WriteString("string/nil,")
		case 0x5:
			buf.WriteString("string,")
		case 0x4:
			buf.WriteString("number/nil,")
		case 0x3:
			buf.WriteString("number,")
		case 0x2:
			buf.WriteString("bool/nil,")
		case 0x1:
			buf.WriteString("bool,")
		case 0x0:
			buf.WriteString("any,")
		}
	}
	trunc()
	buf.WriteString(")")
	return buf.String()
}
