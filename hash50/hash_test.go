package hash50

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func randString() []byte {
	n := rand.Intn(4) + 6
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = revLookup[rand.Intn(len(revLookup))]
	}
	return buf
}

func TestHash50(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1e5; i++ {
		buf := randString()
		h := Hash(buf)
		if !bytes.Equal(FindStringHash(h), buf) {
			t.Fatal(h, string(buf))
		}
	}
}
