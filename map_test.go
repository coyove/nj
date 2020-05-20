package potatolang

import (
	"testing"
)

func BenchmarkTable(b *testing.B) {
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		delete(m, i)
	}
}
