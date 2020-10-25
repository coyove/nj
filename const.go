package script

import "fmt"

const (
	regA   uint16 = 0x1fff // full 13 bits
	regNil uint16 = 0x7ff - 1
)

type opCode byte

const (
	_ opCode = iota
	OpSet
	OpStore
	OpLoad
	OpAdd
	OpConcat
	OpSub
	OpMul
	OpDiv
	OpInc
	OpMod
	OpNot
	OpEq
	OpNeq
	OpLess
	OpLessEq
	OpPow
	OpIf
	OpIfNot
	OpJmp
	OpLoadFunc
	OpPush
	OpPushV
	OpPopV
	OpPopVAll
	OpPopVClear
	OpCall
	OpRet
	OpYield
	OpLen
	OpEOB
)

type valueType byte

const (
	VNil       valueType = 0  // nil
	VNumber              = 3  // number
	VString              = 7  // string
	VStack               = 15 // stack
	VFunction            = 31 // function
	VInterface           = 63 // interface
	_NumNum              = VNumber * 2
	_StrStr              = VString * 2
)

var typeMappings = map[valueType]string{
	VNil:       "nil",
	VNumber:    "number",
	VString:    "string",
	VFunction:  "function",
	VInterface: "interface",
	VStack:     "stack",
}

func panicf(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}

func catchErr(err *error) {
	if r := recover(); r != nil {
		*err, _ = r.(error)
		if *err == nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
