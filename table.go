package potatolang

import (
	"bytes"
	"fmt"
)

type Table struct {
	a  []Value
	m  Map
	mt *Table
}

func (t *Table) _put(k, v Value, raw bool) {
	var ni Value
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
				if !raw {
					if ni = t.mt.RawGet(__newindex); !ni.IsNil() {
						goto newindex
					}
				}

				if !v.IsNil() {
					t.a = append(t.a, v)
				}
				t.m.Put(k, Value{})
				return
			}
		}
	}

	if !raw && t.m.Get(k).IsNil() {
		if ni = t.mt.RawGet(__newindex); !ni.IsNil() {
			goto newindex
		}
	}

	t.m.Put(k, v)
	return

newindex:
	switch ni.Type() {
	case FUN:
		ni.Fun().Call(Tab(t), k, v)
	case TAB:
		if ni.Tab() == t {
			panicf("invalid __newindex, recursive delegation")
		}
		ni.Tab().Put(k, v)
	default:
		panicf("invalid __newindex, expect table or function")
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
	t.RawPut(k, v)
}

func (t *Table) Remove(idx int) Value {
	v := t.a[idx-1]
	t.a = append(t.a[:int(idx)-1], t.a[idx-1+1:]...)
	return v
}

func (t *Table) Put(k, v Value) *Table { t._put(k, v, false); return t }

func (t *Table) RawPut(k, v Value) *Table { t._put(k, v, true); return t }

func (t *Table) Puts(k string, v Value) *Table { return t.Put(Str(k), v) }

func (t *Table) RawPuts(k string, v Value) *Table { return t.RawPut(Str(k), v) }

func (t *Table) Gets(k string) Value { return t.Get(Str(k)) }

func (t *Table) Get(k Value) (v Value) { return t._get(k, false) }

func (t *Table) RawGet(k Value) (v Value) { return t._get(k, true) }

func (t *Table) _get(k Value, raw bool) (v Value) {
	if t == nil {
		return
	}
	if k.Type() == NUM {
		idx := k.Num()
		if float64(int(idx)) == idx {
			if idx >= 1 && int(idx) <= len(t.a) {
				return t.a[int(idx)-1]
			}
		}
	}
	if !raw && t.m.Get(k).IsNil() {
		switch ni := t.mt.RawGet(__index); ni.Type() {
		case FUN:
			v, _ = ni.Fun().Call(Tab(t), k)
			return v
		case TAB:
			if ni.Tab() == t {
				panicf("invalid __index, recursive delegation")
			}
			return ni.Tab().Get(k)
		case NIL:
		default:
			panicf("invalid __index, expect table or function")
		}
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
