package potatolang

import (
	"bytes"
)

type unpacked struct {
	a []Value
}

func (t *unpacked) Put(idx int64, v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		t.a[idx] = v
	}
}

func (t *unpacked) Remove(idx int64) Value {
	idx--
	if idx >= int64(len(t.a)) || idx < 0 {
		return Value{}
	}
	v := t.a[idx]
	t.a = append(t.a[:idx], t.a[idx+1:]...)
	return v
}

func (t *unpacked) Get(idx int64) (v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		return t.a[idx]
	}
	return Value{}
}

func (t *unpacked) Len() int { return len(t.a) }

func (t *unpacked) String() string {
	return t.toString(0)
}

func (t *unpacked) toString(lv int) string {
	p := bytes.NewBufferString("{")
	for _, a := range t.a {
		p.WriteString(a.toString(lv + 1))
		p.WriteString(",")
	}
	if len(t.a) > 0 {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString("}")
	return p.String()
}
