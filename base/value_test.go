package base

import (
	"testing"
)

func TestValue_AttachDetach(t *testing.T) {
	v := NewNumberValue(1993)
	a0 := NewBoolValue(true)
	b0 := NewNumberValue(12)
	c0 := NewBoolValue(false)
	v.Attach(a0)
	v.Attach(b0)
	v.Attach(c0)

	c := v.Detach()
	b := v.Detach()
	a := v.Detach()

	if !a.Equal(a0) || !b.Equal(b0) || !c.Equal(c0) {
		t.Error(a, a0, b, b0, c, c0)
	}

	v.Attach(a0)
	v.Attach(b0)
	b = v.Detach()
	v.Attach(c0)

	c = v.Detach()
	a = v.Detach()

	if !a.Equal(a0) || !b.Equal(b0) || !c.Equal(c0) {
		t.Error(a, a0, b, b0, c, c0)
	}

	v.Attach(a0)
	v.Attach(b0)
	v.Attach(NewNumberValue(float64(0x100))) // overflow 255
	if v.Attachments() != 2 {
		t.Error(v)
	}

	b = v.Detach()
	a = v.Detach()
	if !a.Equal(a0) || !b.Equal(b0) || !c.Equal(c0) {
		t.Error(a, a0, b, b0, c, c0)
	}

	v.Attach(NewNumberValue(float64(0x100)))
	i := v.Detach()
	if i.AsNumberUnsafe() != 0x100 {
		t.Error(i)
	}
}
