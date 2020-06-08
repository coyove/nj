package potatolang

var (
	__metatable = Str("__metatable")
	__tostring  = Str("__tostring")
	__newindex  = Str("__newindex")
	__index     = Str("__index")
	__eq        = Str("__eq")
	__len       = Str("__len")
	__lt        = Str("__lt")
	__le        = Str("__le")
	__add       = Str("__add")
	__concat    = Str("__concat")
	__sub       = Str("__sub")
	__mul       = Str("__mul")
	__div       = Str("__div")
	__mod       = Str("__mod")
	__pow       = Str("__pow")
	__call      = Str("__call")

	nilMetatable *Table
	blnMetatable *Table
	strMetatable *Table
	numMetatable *Table
	funMetatable *Table
	upkMetatable *Table
)

func init() {
	strMetatable = (&Table{}).RawPuts("sub", NativeFun(2, func(env *Env) {
		i, j, s := int(env.In(1, NUM).Num()), -1, env.In(0, STR).Str()
		if len(env.V) > 0 {
			j = int(env.V[0].Expect(NUM).Num())
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
