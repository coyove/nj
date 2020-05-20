package potatolang

import (
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

func tablestringkey(k string) tablekey {
	return tablekey{str: k, g: Value{v: StringType}}
}

func (t *Table) Put(k, v Value) {
	switch k.Type() {
	case StringType:
		if t.m == nil {
			t.m = make(map[tablekey]Value)
		}
		if v.IsNil() {
			delete(t.m, tablestringkey(k.AsString()))
		} else {
			t.m[tablestringkey(k.AsString())] = v
		}
	case NumberType:
		idx := k.AsNumber()
		if float64(int(idx)) == idx {
			if idx >= 0 && int(idx) < len(t.a) {
				t.a[int(idx)] = v
				if v.IsNil() {
					t.Compact()
				}
				return
			}
			if int(idx) == len(t.a) {
				if !v.IsNil() {
					t.a = append(t.a, v)
				}
				delete(t.m, tablekey{g: v})
				return
			}
		}
		fallthrough
	default:
		if t.m == nil {
			t.m = make(map[tablekey]Value)
		}
		if v.IsNil() {
			delete(t.m, tablekey{g: k})
		} else {
			t.m[tablekey{g: k}] = v
		}
	}
}

func (t *Table) Gets(k string) Value {
	v := t.m[tablestringkey(k)]
	if v.IsNil() && t.mt != nil {
		return t.mt.Gets(k)
	}
	return v
}

func (t *Table) Get(k Value) (v Value) {
	switch k.Type() {
	case StringType:
		v = t.m[tablestringkey(k.AsString())]
	case NumberType:
		idx := k.AsNumber()
		if float64(int(idx)) == idx {
			if idx >= 0 && int(idx) < len(t.a) {
				v = t.a[int(idx)]
				break
			}
		}
		fallthrough
	default:
		v = t.m[tablekey{g: k}]
	}
	if v.IsNil() && t.mt != nil {
		return t.mt.Get(k)
	}
	return v
}

func (t *Table) Len() int {
	return len(t.a) + len(t.m)
}

func (t *Table) Compact() {
	// if len(t.a) < 32 { // small array, no need to compact
	// 	return
	// }

	holes := 0
	best := struct {
		ratio float64
		index int
	}{}

	for i, v := range t.a {
		if !v.IsNil() {
			continue
		}

		holes++
		ratio := float64(i-holes) / float64(i)

		if ratio > best.ratio {
			best.ratio = ratio
			best.index = i
		}
	}

	if holes == 0 {
		return
	}

	if best.ratio < 0.5 {
		for i, v := range t.a {
			if !v.IsNil() {
				t.m[tablekey{g: NewNumberValue(float64(i))}] = v
			}
		}
		t.a = nil
		return
	}

	for i := best.index + 1; i < len(t.a); i++ {
		t.m[tablekey{g: NewNumberValue(float64(i))}] = t.a[i]
	}
	t.a = t.a[:best.index+1]
}

type TableIterator struct {
	t     *Table
	miter *reflect.MapIter
	aiter int
}

func (t *Table) Iter() *TableIterator {
	i := &TableIterator{t: t, aiter: -1}
	if t.m != nil {
		i.miter = reflect.ValueOf(t.m).MapRange()
	}
	return i
}

func (iter *TableIterator) Next() bool {
	iter.aiter++
	if iter.aiter < len(iter.t.a) {
		return true
	}
	if iter.miter == nil {
		return false
	}
	return iter.miter.Next()
}

func (iter *TableIterator) Value() Value {
	if iter.aiter < len(iter.t.a) {
		return iter.t.a[iter.aiter]
	}
	if iter.miter == nil {
		return Value{}
	}
	return iter.miter.Value().Interface().(Value)
}

func (iter *TableIterator) Key() Value {
	if iter.aiter < len(iter.t.a) {
		return NewNumberValue(float64(iter.aiter))
	}
	if iter.miter == nil {
		return Value{}
	}
	tk := iter.miter.Key().Interface().(tablekey)
	if tk.g.v == StringType && tk.g.p == nil {
		return NewStringValue(tk.str)
	}
	return tk.g
}
