package strmatch

import (
	"log"
	"testing"
)

func TestMatch(t *testing.T) {
	text := "a12"
	p := "(%d)+"

	log.Println(Find(text, p, 0, false))

	log.Println(Match("hellow 123 world", "%d+ %w+", 0))
}
