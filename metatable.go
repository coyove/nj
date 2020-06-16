package potatolang

var (
	M__metatable = Str("__metatable")
	M__tostring  = Str("__tostring")
	M__newindex  = Str("__newindex")
	M__index     = Str("__index")
	M__eq        = Str("__eq")
	M__len       = Str("__len")
	M__lt        = Str("__lt")
	M__le        = Str("__le")
	M__add       = Str("__add")
	M__concat    = Str("__concat")
	M__sub       = Str("__sub")
	M__unm       = Str("__unm")
	M__mul       = Str("__mul")
	M__div       = Str("__div")
	M__mod       = Str("__mod")
	M__pow       = Str("__pow")
	M__call      = Str("__call")
	M__ipairs    = Str("__ipairs")

	nilMetatable *Table
	blnMetatable *Table
	strMetatable *Table
	numMetatable *Table
	funMetatable *Table
	upkMetatable *Table
)

func init() {
	strMetatable = (&Table{}).Puts("sub", NativeFun(func(env *Env) {
		i, j, s := int(env.In(1, NUM).Num()), -1, env.In(0, STR).Str()
		if len(env.stack) > 2 {
			j = int(env.Get(2).Expect(NUM).Num())
		}
		if i < 0 {
			i = len(s) + i + 1
		}
		if j < 0 {
			j = len(s) + j + 1
		}
		env.A = Str(s[i-1 : j])
	}))
}

func (v Value) GetMetatable() *Table {
	switch t := v.Type(); t {
	case TAB:
		return v.Tab().mt
	case ANY:
		i := v.Any()
		if f, ok := i.(interface{ GetMetatable() *Table }); ok {
			return f.GetMetatable()
		}
		return nil
	case NIL:
		return nilMetatable
	case BLN:
		return blnMetatable
	case STR:
		return strMetatable
	case NUM:
		return numMetatable
	case FUN:
		return funMetatable
	case UPK:
		return upkMetatable
	default:
		panic("corrupted value")
	}
}

func (v Value) GetMetamethod(name Value) Value {
	if mt := v.GetMetatable(); mt != nil {
		return mt.RawGet(name)
	}
	return Value{}
}

func (v Value) SetMetatable(mt *Table) {
	switch t := v.Type(); t {
	case TAB:
		v.Tab().mt = mt
	case ANY:
		i := v.Any()
		if f, ok := i.(interface{ SetMetatable(*Table) }); ok {
			f.SetMetatable(mt)
		}
	case NIL:
		nilMetatable = mt
	case BLN:
		blnMetatable = mt
	case STR:
		strMetatable = mt
	case NUM:
		numMetatable = mt
	case FUN:
		funMetatable = mt
	case UPK:
		upkMetatable = mt
	default:
		panic("corrupted value")
	}
}

func findmm(a, b Value, name Value) Value {
	m := a.GetMetamethod(name)
	if m.IsNil() {
		return b.GetMetamethod(name)
	}
	return m
}
