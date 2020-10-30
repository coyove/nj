package script

import (
	"fmt"
	"strconv"
)

const regA uint16 = 0x1fff // full 13 bits

type opCode byte

const (
	_ opCode = iota
	OpSet
	OpStore
	OpLoad
	OpSlice
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
	OpCallMap
	OpRet
	OpYield
	OpLen
	OpJSON
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

func (t valueType) String() string {
	switch t {
	case VNil:
		return "nil"
	case VNumber:
		return "number"
	case VString:
		return "string"
	case VFunction:
		return "function"
	case VInterface:
		return "interface"
	case VStack:
		return "stack"
	default:
		return "corrupted$" + strconv.Itoa(int(t))
	}
}

func panicf(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}
