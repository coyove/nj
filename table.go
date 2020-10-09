package potatolang

import (
	"bytes"
	"fmt"
)

type Table struct {
	unpacked bool
	a        []Value
	m        Map
}

func (t *Table) Put(k, v Value) {
	if k.Type() == NUM {
		_, idx, _ := k.Num()
		if idx >= 1 && int(idx) <= len(t.a) {
			t.a[int(idx)-1] = v
			if v.IsNil() {
				t.Compact()
			}
			return
		}
		if int(idx) == len(t.a)+1 {
			if !v.IsNil() {
				t.a = append(t.a, v)
			}
			t.m.Put(k, Value{})
			return
		}
	}

	t.m.Put(k, v)
}

func (t *Table) Insert(k, v Value) {
	if k.Type() == NUM {
		_, idx, _ := k.Num()
		if idx >= 1 && int(idx) <= len(t.a) {
			t.a = append(t.a[:int(idx)-1], append([]Value{v}, t.a[int(idx)-1:]...)...)
			if v.IsNil() {
				t.Compact()
			}
			return
		}
	}
	t.Put(k, v)
}

func (t *Table) Remove(idx int) Value {
	v := t.a[idx-1]
	t.a = append(t.a[:int(idx)-1], t.a[idx-1+1:]...)
	return v
}

func (t *Table) Get(k Value) (v Value) { return t._get(k, false) }

func (t *Table) _get(k Value, lookinarray bool) (v Value) {
	if t == nil {
		return
	}
	if k.Type() == NUM {
		_, idx, _ := k.Num()
		if idx >= 1 && int(idx) <= len(t.a) {
			return t.a[int(idx)-1]
		}
	}
	if lookinarray {
		return Value{}
	}
	return t.m.Get(k)
}

func (t *Table) Len() int { return len(t.a) }

func (t *Table) HashLen() int { return t.m.Len() }

func (t *Table) Compact() {
	for i := len(t.a) - 1; i >= 0; i-- {
		if t.a[i].IsNil() {
			t.a = t.a[:i]
		} else {
			break
		}
	}
}

func (t *Table) String() string {
	p := bytes.NewBufferString("{")
	for i := range t.a {
		p.WriteString(fmt.Sprintf("%v,", t.a[i]))
	}
	if len(t.a) > 0 {
		p.Truncate(p.Len() - 1)
		if t.m.Len() > 0 {
			p.WriteString(";")
		}
	}
	for k, v := t.m.Next(Value{}); !k.IsNil(); k, v = t.m.Next(k) {
		p.WriteString(fmt.Sprintf("[%v]=%v,", k, v))
	}
	if p.Bytes()[p.Len()-1] == ',' {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString("}")
	return p.String()
}

func (t *Table) iterStringKeys(cb func(k string, v Value)) {
	for k, v := t.m.Next(Value{}); !k.IsNil(); k, v = t.m.Next(k) {
		if k.Type() == STR {
			cb(k.Str(), v)
		}
	}
}

func (t *Table) Next(k Value) (Value, Value) {
	if k.IsNil() {
		for i := 1.0; ; i++ {
			v := t.Get(Num(i))
			if v.IsNil() {
				break
			}
			return Num(float64(i)), v
		}
		return t.m.Next(Value{})
	}

	if k.Type() == NUM {
		_, idx, _ := k.Num()
		if v := t._get(Int(idx+1), true); !v.IsNil() {
			return Int(idx + 1), v
		}

		if v := t._get(Int(idx), true); !v.IsNil() {
			return t.m.Next(Value{})
		}
	}
	return t.m.Next(k)
}

func buildtable(args ...interface{}) *Table {
	t := &Table{m: *NewMap(len(args) / 2)}
	for i := 0; i < len(args); i += 2 {
		k := args[i].(string)
		switch x := args[i+1].(type) {
		case func(env *Env):
			f := Native(x)
			f.Fun().Name = k
			t.Put(Str(k), f)
		case *Table:
			t.Put(Str(k), Tab(x))
		case Value:
			t.Put(Str(k), x)
		default:
			panicf("build: %T", x)
		}
	}
	return t
}
