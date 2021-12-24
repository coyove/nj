package bas

import (
	"time"
	"unsafe"

	"github.com/coyove/nj/typ"
)

type SafeValue Value

func (v SafeValue) Str(defaultValue string) string {
	switch Value(v).Type() {
	case typ.String:
		return Value(v).Str()
	case typ.Array:
		buf, ok := Value(v).Array().Unwrap().([]byte)
		if ok {
			return *(*string)(unsafe.Pointer(&buf))
		}
	}
	return defaultValue
}

func (v SafeValue) Int(defaultValue int) int {
	return int(v.Int64(int64(defaultValue)))
}

func (v SafeValue) Int64(defaultValue int64) int64 {
	if Value(v).Type() == typ.Number {
		return Value(v).Int64()
	}
	return defaultValue
}

func (v SafeValue) Float64(defaultValue float64) float64 {
	if Value(v).Type() == typ.Number {
		return Value(v).Float64()
	}
	return defaultValue
}

func (v SafeValue) Array() *Array {
	if Value(v).Type() != typ.Array {
		return nil
	}
	return Value(v).Array()
}

func (v SafeValue) Object() *Object {
	if Value(v).Type() != typ.Object {
		return nil
	}
	return Value(v).Object()
}

func (v SafeValue) Bytes() []byte {
	switch Value(v).Type() {
	case typ.String:
		return []byte(Value(v).Str())
	case typ.Array:
		buf, _ := Value(v).Array().Unwrap().([]byte)
		return buf
	}
	return nil
}

func (v SafeValue) Duration(defaultValue time.Duration) time.Duration {
	if Value(v).Type() != typ.Number {
		return defaultValue
	}
	if Value(v).IsInt64() {
		return time.Duration(Value(v).Int64()) * time.Second
	}
	return time.Duration(Value(v).Float64() * float64(time.Second))
}

func (v SafeValue) Error() error {
	if Value(v).Type() != typ.Array || Value(v).Array().meta != errorArrayMeta {
		return nil
	}
	return Value(v).Array().Unwrap().(*ExecError)
}

func DeepEqual(a, b Value) bool {
	if a.Equal(b) {
		return true
	}
	if at, bt := a.Type(), b.Type(); at == bt {
		switch at {
		case typ.Array:
			flag := a.Array().Len() == b.Array().Len()
			if !flag {
				return false
			}
			a.Array().ForeachIndex(func(k int, v Value) bool {
				flag = DeepEqual(b.Array().Get(k), v)
				return flag
			})
			return flag
		case typ.Object:
			flag := a.Object().Len() == b.Object().Len()
			if !flag {
				return false
			}
			a.Object().Foreach(func(k Value, v *Value) int {
				flag = DeepEqual(b.Object().Get(k), *v)
				if flag {
					return typ.ForeachContinue
				}
				return typ.ForeachBreak
			})
			return flag
		}
	}
	return false
}

func lessStr(a, b Value) bool {
	if a.isSmallString() && b.isSmallString() {
		if a.v == b.v {
			return uintptr(a.p) < uintptr(b.p) // a is shorter than b
		}
		return a.v < b.v
	}
	return a.Str() < b.Str()
}

func Less(a, b Value) bool {
	at, bt := a.Type(), b.Type()
	if at != bt {
		return at < bt
	}
	switch at {
	case typ.Number:
		if a.IsInt64() && b.IsInt64() {
			return a.UnsafeInt64() < b.UnsafeInt64()
		}
		return a.Float64() < b.Float64()
	case typ.String:
		return lessStr(a, b)
	}
	return a.UnsafeAddr() < b.UnsafeAddr()
}

func IsPrototype(a Value, p *Object) bool {
	switch a.Type() {
	case typ.Nil:
		return p == nil
	case typ.Object:
		return a.Object().IsPrototype(p)
	case typ.Bool:
		return p == Proto.Bool
	case typ.Number:
		return p == Proto.Float || (a.IsInt64() && p == Proto.Int)
	case typ.String:
		return p == Proto.Str
	case typ.Array:
		return a.Array().meta.Proto.IsPrototype(p)
	}
	return false
}

func IsCallable(a Value) bool {
	return a.Type() == typ.Object && a.Object().IsCallable()
}

func IntersectObject(a, b *Object) {
}
