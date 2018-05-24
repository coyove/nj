package potatolang

import (
	"testing"
)

func TestBytesWriter(t *testing.T) {
	const v = 10

	w := NewBytesWriter()
	w.Write(v)
	w.Write32(v)
	w.WriteDouble(v)
	w.WriteInt64(v)
	w.WriteString("10")

	buf := w.Bytes()
	cursor := uint32(0)
	if crRead(buf, &cursor) != v {
		t.Error(cursor)
	}
	if crRead32(buf, &cursor) != v {
		t.Error(cursor)
	}
	if crReadDouble(buf, &cursor) != v {
		t.Error(cursor)
	}
	if crReadInt64(buf, &cursor) != v {
		t.Error(cursor)
	}
	if crReadString(buf, &cursor) != "10" {
		t.Error(cursor)
	}

}
