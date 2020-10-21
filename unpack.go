package potatolang

import (
	"bytes"
	"fmt"
)

type unpacked struct {
	a []Value
}

func (t *unpacked) Put(idx int64, v Value) {
	idx--
	if idx < int64(len(t.a)) {
		t.a[idx] = v
	}
}

func (t *unpacked) Remove(idx int64) Value {
	idx--
	if idx >= int64(len(t.a)) {
		return Value{}
	}
	v := t.a[idx]
	t.a = append(t.a[:idx], t.a[idx+1:]...)
	return v
}

func (t *unpacked) Get(idx int64) (v Value) {
	idx--
	if idx < int64(len(t.a)) {
		return t.a[idx]
	}
	return Value{}
}

func (t *unpacked) Len() int { return len(t.a) }

func (t *unpacked) String() string {
	p := bytes.NewBufferString("{")
	for i := range t.a {
		p.WriteString(fmt.Sprintf("%v,", t.a[i]))
	}
	if len(t.a) > 0 {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString("}")
	return p.String()
}
