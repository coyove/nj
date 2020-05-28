package potatolang

import (
	"bytes"
	"fmt"
	"reflect"
)

type tablekey struct {
	g   Value
	str string
}

type Table struct {
	a  []Value
	m  map[tablekey]Value
	mt *Table
}

func maketk(k Value) tablekey {
	if k.Type() == STR {
		return tablekey{str: k.Str(), g: Value{v: STR}}
	}
	return tablekey{g: k}
}

func (t *Table) rawgetstr(name string) Value {
	if t == nil {
		return Value{}
	}
	return t.m[tablekey{str: name, g: Value{v: STR}}]
}

func (t *Table) Put(k, v Value, raw bool) {
	if k.Type() == NUM {
		idx := k.Num()
		if float64(int(idx)) == idx {
			if idx >= 1 && int(idx) <= len(t.a) {
				t.a[int(idx)-1] = v
				if v.IsNil() {
					t.Compact()
				}
				return
			}
			if int(idx) == len(t.a)+1 {
				if !raw && !t.mt.rawgetstr("__newindex").IsNil() {
					t.newindex(k, v)
					return
				}

				if !v.IsNil() {
					t.a = append(t.a, v)
				}
				delete(t.m, tablekey{g: k})
				return
			}
		}
	}

	if t.m == nil {
		t.m = make(map[tablekey]Value)
	}
	key := maketk(k)
	if !raw && !t.mt.rawgetstr("__newindex").IsNil() && t.m[key].IsNil() {
		t.newindex(k, v)
		return
	}
	if v.IsNil() {
		delete(t.m, key)
	} else {
		t.m[key] = v
	}

}

func (t *Table) Insert(k, v Value) {
	if k.Type() == NUM {
		idx := k.Num()
		if float64(int(idx)) == idx {
			if idx >= 1 && int(idx) <= len(t.a) {
				t.a = append(t.a[:int(idx)-1], append([]Value{v}, t.a[int(idx)-1:]...)...)
				if v.IsNil() {
					t.Compact()
				}
				return
			}
		}
	}
	t.Put(k, v, true)
}

func (t *Table) newindex(k, v Value) {
	switch ni := t.mt.rawgetstr("__newindex"); ni.Type() {
	case FUN:
		ni.Fun().Call(Tab(t), k, v)
	case TAB:
		ni.Tab().Put(k, v, false)
	default:
		panicf("invalid __newindex")
	}
}

func (t *Table) Puts(k string, v Value, raw bool) *Table {
	t.Put(Str(k), v, raw)
	return t
}

func (t *Table) Gets(k string, raw bool) Value {
	return t.Get(Str(k), raw)
}

func (t *Table) Get(k Value, raw bool) (v Value) {
	if k.Type() == NUM {
		idx := k.Num()
		if float64(int(idx)) == idx {
			if idx >= 1 && int(idx) <= len(t.a) {
				return t.a[int(idx)-1]
			}
		}
	}
	key := maketk(k)
	if !raw && !t.mt.rawgetstr("__index").IsNil() && t.m[key].IsNil() {
		switch ni := t.mt.rawgetstr("__index"); ni.Type() {
		case FUN:
			v, _ = ni.Fun().Call(Tab(t), k)
			return v
		case TAB:
			return ni.Tab().Get(k, false)
		default:
			panicf("invalid __index")
		}
	}
	return t.m[key]
}

func (t *Table) Len() int {
	return len(t.a)
}

func (t *Table) HashLen() int {
	return len(t.m)
}

func (t *Table) Compact() {
	for i := len(t.a) - 1; i >= 0; i-- {
		if t.a[i].IsNil() {
			t.a = t.a[:i]
		} else {
			break
		}
	}
}

type TableMapIterator struct {
	t     *Table
	miter *reflect.MapIter
}

func (t *Table) Iter() *TableMapIterator {
	i := &TableMapIterator{t: t}
	if t.m != nil {
		i.miter = reflect.ValueOf(t.m).MapRange()
	}
	return i
}

func (iter *TableMapIterator) Next() bool {
	if iter.miter == nil {
		return false
	}
	return iter.miter.Next()
}

func (iter *TableMapIterator) Value() Value {
	if iter.miter == nil {
		return Value{}
	}
	return iter.miter.Value().Interface().(Value)
}

func (iter *TableMapIterator) Key() Value {
	if iter.miter == nil {
		return Value{}
	}
	tk := iter.miter.Key().Interface().(tablekey)
	if tk.g.v == STR && tk.g.p == nil {
		return Str(tk.str)
	}
	return tk.g
}

func (t *Table) String() string {
	p := bytes.NewBufferString("{")
	for i := range t.a {
		p.WriteString(fmt.Sprintf("[%d]=%v,", i+1, t.a[i].toString(0, true)))
	}
	for k, v := range t.m {
		if k.g.v == STR {
			p.WriteString(fmt.Sprintf("[%q]=%v,", k.str, v))
		} else if k.g.Type() == NUM {
			p.WriteString(fmt.Sprintf("[%v]=%v,", k.g, v.toString(0, true)))
		} else {
			p.WriteString(fmt.Sprintf("[%v]=%v,", k.g, v))
		}
	}
	if p.Bytes()[p.Len()-1] == ',' {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString("}")
	return p.String()
}

func (t *Table) iterStringKeys(cb func(k string, v Value)) {
	for k, v := range t.m {
		if k.g.v == STR {
			cb(k.str, v)
		}
	}
}
