package potatolang

import (
	"runtime"
	"strconv"
	"testing"
)

func stringChannel(ch chan string, s Value, flag bool) {
	// prevent inlining
	switch s.ty {
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

	for i := 0; i < 10000; i++ {
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
