package typ

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
