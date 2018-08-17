// Copyright (c) 2015, Emir Pasic. All rights reserved.

// Use of this source code is governed by a BSD-style

// license that can be found in the LICENSE file.

package potatolang

import (
	"math/rand"
	"strconv"
	"testing"
)

const benchSize = 40

func BenchmarkMap(b *testing.B) {
	m := NewMap()
	for j := 0; j < benchSize; j++ {
		m.Put(NewNumberValue(float64(j)), NewNumberValue(float64(j)))
	}

	for i := 0; i < b.N; i++ {
		x := float64(rand.Intn(benchSize))
		y := NewNumberValue(x)
		if v, _ := m.Get(y); v.AsNumber() != x {
			b.Error("shouldn't happen")
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	m := NewMap()
	for j := 0; j < benchSize; j++ {
		m.Put(NewNumberValue(float64(j)), NewNumberValue(float64(j)))
	}

	for i := 0; i < b.N; i++ {
		x := NewNumberValue(benchSize)
		m.Put(x, x)
		m.Remove(x)
	}
}

func BenchmarkNativeSlice(b *testing.B) {
	m := make([]Value, benchSize)
	for j := 0; j < benchSize; j++ {
		m[j] = NewNumberValue(float64(j))
	}

	x := float64(rand.Intn(benchSize))
	y := NewNumberValue(x)
	for i := 0; i < b.N; i++ {
		if m[int(y.AsNumber())].AsNumber() != x {
			b.Error("shouldn't happen")
		}
	}
}

func BenchmarkNativeSliceAppend(b *testing.B) {
	m := make([]Value, benchSize)
	for j := 0; j < benchSize; j++ {
		m[j] = NewNumberValue(float64(j))
	}
	for i := 0; i < b.N; i++ {
		m = append(m, NewNumberValue(benchSize))
		m = m[:len(m)-1]
	}
}

func BenchmarkNativeMap(b *testing.B) {
	m := map[string]Value{}
	for j := 0; j < benchSize; j++ {
		m[strconv.Itoa(j)] = NewNumberValue(float64(j))
	}

	for i := 0; i < b.N; i++ {
		x := float64(rand.Intn(benchSize))
		if m[strconv.Itoa(int(x))].AsNumber() != x {
			b.Error("shouldn't happen")
		}
	}
}

func TestMap_Put(t *testing.T) {
	m := NewMap()
	for j := 0; j < 10; j++ {
		m.Put(NewStringValue(strconv.Itoa(j)), NewNumberValue(float64(j)))
	}
	m.Put(NewStringValue("0"), NewValue()) // overwrite

	if v, f := m.Get(NewStringValue("0")); v.Type() != Tnil || !f {
		t.Error(0)
	}
	for j := 1; j < 10; j++ {
		if v, _ := m.Get(NewStringValue(strconv.Itoa(j))); v.AsNumber() != float64(j) {
			t.Error(j, v)
		}
	}

	m.Put(NewStringValue("5"), NewNumberValue(5))
	if v, _ := m.getFromMap(NewStringValue("5")); v.AsNumber() != 5 {
		t.Error(m)
	}

}
