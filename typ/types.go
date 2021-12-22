package typ

import "fmt"

type ValueType byte

const (
	Nil    ValueType = 0
	Bool   ValueType = 1
	Number ValueType = 3
	String ValueType = 7
	Object ValueType = 15
	Array  ValueType = 17
	Native ValueType = 19
)

func (t ValueType) String() string {
	if t > Native {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "object", "?", "array", "?", "native"}[t]
}

const (
	RegA          uint16 = 0xffff
	RegPhantom    uint16 = 0xfffe
	RegLocalMask         = 0x7fff
	RegGlobalFlag        = 0x8000
	RegMaxAddress        = 0x7f00
)

type Symbol struct {
	Address uint16
}

func (s *Symbol) String() string {
	return fmt.Sprintf("symbol:%d", s.Address)
}
