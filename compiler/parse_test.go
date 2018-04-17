package compiler

import (
	"bytes"
	"log"
	"testing"

	"github.com/coyove/eugine/base"
)

func TestParsingExtFlat(t *testing.T) {
	tr, _ := newTokenReader(`
		[var a [var b 1]]
		   [set a [- b [b/and a 1]]]
		   [var c [+ b]]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	expected := []byte{
		10, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63,
		9, 2, 0, 0, 0, 1, 0, 0, 0,
		250, 52, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63,
		251, 33, 1, 0, 0, 0, 255, 255, 255, 255,
		9, 2, 0, 0, 0, 255, 255, 255, 255,
		192, 1, 0, 0, 0, 17, 0, 0, 0, 132, 1, 0, 1, 0, 132, 0, 0, 0, 0, 32, 140, 255, 255, 255, 255, 255, 9, 4, 0, 0, 0, 255, 255, 255, 255, 255,
	}

	log.Println(base.NewBytesReader(buf).Prettify(0))

	if !bytes.Equal(expected, buf) {
		t.Error("failed")
	}
}

func TestParsingStoreLoadSugar(t *testing.T) {
	tr, _ := newTokenReader(`
		[var a [list 1 2 3]]
		[assert [a:1] 2]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(buf, base.NewBytesReader(buf).Prettify(0))
}

func TestParsingCall(t *testing.T) {
	tr, _ := newTokenReader(`
		[var a [lambda [a] {var b [+ a 1]}]]
		[a 1]
		[[lambda [] [ret [+ 1 1]]]]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(buf, base.NewBytesReader(buf).Prettify(0))
}

func TestParsingInc(t *testing.T) {
	tr, _ := newTokenReader(`
		[var a 1] [inc a 1]
		[set a [list 1 2]]
		[inc [a:1] 2]
		[inc [a:1] [a:0]]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(buf, base.NewBytesReader(buf).Prettify(0))
}

func TestParsingIf(t *testing.T) {
	tr, _ := newTokenReader(` [if 1 [var a "a"] else [var b "b"]]
		 [set a 0]
		 [if [eq a [- 1 1]] [assert 1] else [assert 0]]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(buf, base.NewBytesReader(buf).Prettify(0))
}

func TestParsingWhile(t *testing.T) {
	tr, _ := newTokenReader(`[var i 0]
		[while [< [- i 0] 50] 
			[inc i 1]
			[if [eq i 10] [break]]
			[if [eq i 12] [continue]]
		]`, true)

	buf, err := tr.parse()
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(buf, base.NewBytesReader(buf).Prettify(0))
}
