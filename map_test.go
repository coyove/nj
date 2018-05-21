// Copyright (c) 2015, Emir Pasic. All rights reserved.

// Use of this source code is governed by a BSD-style

// license that can be found in the LICENSE file.

package potatolang

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMap()
		for j := 0; j < 4; j++ {
			m.Put(strconv.Itoa(j), NewNumberValue(float64(j)))
		}
		if v, _ := m.Get(strconv.Itoa(rand.Intn(4))); v.AsNumberUnsafe() == -1 {
			b.Error("won't happen")
		}
	}
}

func BenchmarkNativeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]Value{}
		for j := 0; j < 4; j++ {
			m[strconv.Itoa(j)] = NewNumberValue(float64(j))
		}
		if m[strconv.Itoa(rand.Intn(4))].AsNumberUnsafe() == -1 {
			b.Error("won't happen")
		}
	}
}

func TestMap_Put(t *testing.T) {
	m := NewMap()
	for j := 0; j < mapentrySize; j++ {
		m.Put(strconv.Itoa(j), NewNumberValue(float64(j)))
	}
	m.Put("0", NewValue()) // overwrite

	if m.t == nil {
		t.Error("t is nil")
	}

	if v, f := m.Get("0"); v.ty != Tnil || !f {
		t.Error(0)
	}
	for j := 0; j < mapentrySize; j++ {
		if v, _ := m.Get(strconv.Itoa(j)); v.AsNumberUnsafe() != float64(j) {
			t.Error(j)
		}
	}

	m.Put("5", NewValue())
	if m.t != nil {
		t.Error("t is not nil")
	}

	keys := "0:1:2:3:5:"
	for k := range m.m {
		if !strings.Contains(keys, k+":") {
			t.Error(k)
		}
	}
}
