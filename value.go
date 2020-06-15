package potatolang

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/cespare/xxhash"
)

const (
	// NIL: nil
	NIL = 0
	// BLN: boolean
	BLN = 1
	// NUM: number
	NUM = 3
	// STR: string
	STR = 7
	// TAB: table
	TAB = 15
	// FUN: function
	FUN = 31
	// ANY: generic type
	ANY = 63
	// UPK: unpacked values
	UPK = 255
)

const (
	NilNil = NIL * 2
	NumNum = NUM * 2
	BlnBln = BLN * 2
	StrStr = STR * 2
	TabTab = TAB * 2
	FunFun = FUN * 2
	AnyAny = ANY * 2
)

// Value is the basic value used by the intepreter
// For numbers there are 4 NaNs which are not representable: [0xffffffff`fffffffc, 0xffffffff`ffffffff]
// An empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

const SizeOfValue = unsafe.Sizeof(Value{})

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() byte {
	if v.p == nil {
		if v.v <= 3 {
			// v.v==0: nil, v.v==1: true, v.v==3: false
			//      0: niltype   1: booltype   3 & 1 == 1 -> booltype
			return byte(v.v & 1)
		}
		return NUM
	}
	return byte(v.v)
}

// IsFalse tests whether value contains a falsy value
// Type assertion should be aligned with Type()
func (v Value) IsFalse() bool {
	// if v.p == nil {
	// 	return v.v == 0 || v.v == 3
	// }
	// return false
	x := uintptr(v.p) + uintptr(v.v)
	return x == 0 || x == 3
}

var (
	typeMappings = map[byte]string{
		NIL: "nil", BLN: "boolean", NUM: "number", STR: "string", FUN: "function", ANY: "any", TAB: "table", UPK: "unpacked",
	}
	emptyUPK = []Value{}
)

func init() {
	initCoreLibs()
}

func newUnpackedValue(stack []Value) Value {
	if stack == nil {
		stack = emptyUPK
	}
	return Value{v: UPK, p: unsafe.Pointer(&stack)}
}

// Num returns a number value
func Num(f float64) Value {
	x := *(*uint64)(unsafe.Pointer(&f))
	return Value{v: ^x}
}

// Bln returns a boolean value
func Bln(b bool) Value {
	if !b {
		return Value{v: 3}
	}
	return Value{v: 1}
}

// Tab returns a table value
func Tab(m *Table) Value {
	if m == nil {
		return Value{}
	}
	return Value{v: TAB, p: unsafe.Pointer(m)}
}

// Fun returns a closure value
func Fun(c *Closure) Value {
	return Value{v: FUN, p: unsafe.Pointer(c)}
}

// Str returns a string value
func Str(s string) Value {
	if len(s) <= 16 {
		b := [16]byte{}
		copy(b[:], s)
		return Value{v: uint64(len(s)+1)<<56 | STR, p: unsafe.Pointer(&b)}
	}
	return Value{v: STR, p: unsafe.Pointer(&s)}
}

func _StrBytes(s []byte) Value {
	return Value{v: STR, p: unsafe.Pointer(&s)}
}

func Any(i interface{}) Value {
	switch v := i.(type) {
	case bool:
		return Bln(v)
	case float64:
		return Num(v)
	case string:
		return Str(v)
	case *Table:
		return Tab(v)
	case *Closure:
		return Fun(v)
	}
	return Value{v: ANY, p: unsafe.Pointer(&i)}
}

// Str cast value to string
func (v Value) Str() string {
	if l := v.v >> 56; l > 0 {
		var ss string
		b := (*[2]uintptr)(unsafe.Pointer(&ss))
		(*b)[0] = uintptr(v.p)
		(*b)[1] = uintptr(l - 1)
		return ss
	}
	return *(*string)(v.p)
}

func (v Value) _StrBytes() []byte {
	var ss []byte
	b := (*[3]uintptr)(unsafe.Pointer(&ss))
	if l := v.v >> 56; l > 0 {
		(*b)[0] = uintptr(v.p)
		(*b)[1], (*b)[2] = uintptr(l-1), uintptr(l-1)
	} else {
		vpp := *(*[2]uintptr)(v.p)
		(*b)[0] = vpp[0]
		(*b)[1], (*b)[2] = vpp[1], vpp[1]
	}
	return ss
}

func (v Value) IsNil() bool { return v == Value{} }

// Num cast value to float64
func (v Value) Num() float64 { return math.Float64frombits(^v.v) }

// Bln cast value to bool
func (v Value) Bln() bool { return v.v == 1 }

// Tab cast value to map of values
func (v Value) Tab() *Table { return (*Table)(v.p) }

func (v Value) _Upk() []Value {
	if v.p == unsafe.Pointer(&emptyUPK) {
		return nil
	}
	return *(*[]Value)(v.p)
}

func (v Value) _TestUpkLen() (int, bool) {
	if v.Type() == UPK {
		return len(v._Upk()), true
	}
	return 1, false
}

func (v Value) _AppendTo(arr []Value) []Value {
	if v.Type() == UPK {
		return append(arr, v._Upk()...)
	}
	return append(arr, v)
}

// Fun cast value to closure
func (v Value) Fun() *Closure { return (*Closure)(v.p) }

// Any returns the interface{}
func (v Value) Any() interface{} {
	switch v.Type() {
	case BLN:
		return v.Bln()
	case NUM:
		return v.Num()
	case STR:
		return v.Str()
	case TAB:
		return v.Tab()
	case FUN:
		return v.Fun()
	case ANY:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) Expect(t byte) Value {
	if v.Type() != t {
		panicf("expect %s, got %s", typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

func (v Value) ExpectMsg(t byte, msg string) Value {
	if v.Type() != t {
		panicf("%s: expect %s, got %s", msg, typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

func (v Value) String() string { return v.toString(0) }

// Equal tests whether value is equal to another value
func (v Value) Equal(r Value) bool {
	switch v.Type() + r.Type() {
	case NumNum, BlnBln, NilNil:
		return v == r
	case StrStr:
		return r.Str() == v.Str()
	case TabTab:
		if v == r {
			return true
		}
	case FunFun:
		return v.Fun() == r.Fun()
	}
	if eq := findmm(v, r, M__eq); eq.Type() == FUN {
		e, _ := eq.Fun().Call(v, r)
		return !e.IsFalse()
	}
	return false
}

//go:nosplit
func (v Value) Hash() uint64 {
	if v.Type() == STR {
		return xxhash.Sum64(v._StrBytes())
	}
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = uintptr(unsafe.Pointer(&v))
	bh.Len = int(unsafe.Sizeof(v))
	bh.Cap = int(unsafe.Sizeof(v))
	return xxhash.Sum64(b)
}

func (v Value) toString(lv int) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}
	switch v.Type() {
	case BLN:
		return strconv.FormatBool(v.Bln())
	case NUM:
		return strconv.FormatFloat(v.Num(), 'f', -1, 64)
	case STR:
		return v.Str()
	case TAB:
		return v.Tab().String()
	case FUN:
		return v.Fun().String()
	case ANY:
		return fmt.Sprintf("<any:%v>", v.Any())
	case UPK:
		return fmt.Sprintf("<unpacked:%v>", v._Upk())
	}
	return "nil"
}
