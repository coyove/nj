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
	RegA          = 0x8000
	RegPhantom    = 0xffff
	RegLocalMask  = 0x7fff
	RegNil        = 0x8001
	RegMaxAddress = 0x7ff0
)
