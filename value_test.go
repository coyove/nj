package potatolang

import (
	"io/ioutil"
	"runtime"
	"strconv"
	"testing"
)

func stringChannel(ch chan string, s Value, flag bool) {
	// prevent inlining
	switch s.Type() {
	case Tstring, Tnil:
	default:
		panic(1)
	}
	if flag {
		ch <- s.AsString()
	}
}

func TestNewStringValue(t *testing.T) {
	ch := make(chan string, 10)

	for i := 0; i < 10; i++ {
		stringChannel(ch, NewStringValue(strconv.Itoa(i)), true)
	}

	for i := 0; i < 1000000; i++ {
		stringChannel(ch, NewValue(), false)
	}
	close(ch)
	runtime.GC()

	i := 0
	for c := range ch {
		if c != strconv.Itoa(i) {
			t.Error(c)
		}
		i = i + 1
	}
}

func TestStringValueHash(t *testing.T) {
	buf, _ := ioutil.ReadFile("value_test.go")
	parts := make([]Value, 0)
	i := 0
	// read 8 bytes at a time, small enough to keep them in a single Value struct
	for i < len(buf) {
		if i+8 < len(buf) {
			parts = append(parts, NewStringValue(string(buf[i:i+8])))
			i += 8
			continue
		}

		parts = append(parts, NewStringValue(string(buf[i:])))
		break
	}

	str := ""
	for _, p := range parts {
		str += p.Str()
	}

	a := NewStringValue(str).Hash()
	b := NewStringValue(string(buf)).Hash()
	if a.a != b.a || a.b != b.b {
		t.Error("hash not matched")
	}

	// t.Error(NewStringValue("zzz").hashstr(), NewStringValue("zzy").hashstr())
}
