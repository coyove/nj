//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//  http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package script

import (
	"sync/atomic"
	"unsafe"
)

const (
	growRate = 1.25
	errRace  = "concurrent map read and map write"
)

type RHMap struct {
	Parent    *RHMap
	rwlock    int32
	count     uint32
	hashItems []hashItem
	items     []Value
}

// hashItem represents an entry in the Map.
type hashItem struct {
	Key      Value
	Val      Value
	Distance int // How far item is from its best position.
}

func NewMap(size int) *RHMap {
	return &RHMap{hashItems: make([]hashItem, int64(size)*int64(growRate*16)/16+1)}
}

func (m *RHMap) Len() int {
	return int(m.count)
}

func (m *RHMap) MapLen() int {
	return int(m.count) - len(m.items)
}

// Clear clears Map, where already allocated memory will be reused.
func (m *RHMap) Clear() {
	if atomic.LoadInt32(&m.rwlock) > 0 {
		panicf(errRace)
	}
	m.hashItems = m.hashItems[:0]
	m.count = 0
	m.items = m.items[:0]
}

func (m *RHMap) GetString(k string) (v Value) {
	return m.Get(Str(k))
}

// Get retrieves the val for a given key.
func (m *RHMap) Get(k Value) (v Value) {
	if k == Nil {
		return Nil
	}
	atomic.AddInt32(&m.rwlock, 1)
	defer atomic.AddInt32(&m.rwlock, -1)

	if k.IsInt() {
		if idx := k.Int(); idx >= 0 && idx < int64(len(m.items)) {
			return m.items[idx]
		}
	}
	if idx := m.findHash(k); idx >= 0 {
		return m.hashItems[idx].Val
	}
	if m.Parent != nil {
		return m.Parent.Get(k)
	}
	return Nil
}

func (m *RHMap) findHash(k Value) int {
	num := len(m.hashItems)
	if num <= 0 {
		return -1
	}
	idx := int(k.HashCode() % uint64(num))
	idxStart := idx

	for {
		e := &m.hashItems[idx]
		if e.Key == Nil {
			return -1
		}

		if e.Key.Equal(k) {
			return idx
		}

		idx++
		if idx >= num {
			idx = 0
		}

		if idx == idxStart { // Went all the way around.
			return -1
		}
	}
}

func (m *RHMap) Contains(k Value) bool {
	atomic.AddInt32(&m.rwlock, 1)
	defer atomic.AddInt32(&m.rwlock, -1)
	if k == Nil {
		return false
	}
	if k.IsInt() {
		if idx := k.Int(); idx >= 0 && idx < int64(len(m.items)) {
			return true
		}
	}
	return m.findHash(k) >= 0
}

func (m *RHMap) ParentContains(k Value) *RHMap {
	if k == Nil {
		return nil
	}
	if m.Parent != nil {
		p := m.Parent.ParentContains(k)
		if p != nil {
			return p
		}
	}
	if m.Contains(k) {
		return m
	}
	return nil
}

func (m *RHMap) SetString(k string, v Value) (prev Value) {
	return m.Set(Str(k), v)
}

// Set inserts or updates a key/val pair into the Map. If val == Nil, then key will get deleted
func (m *RHMap) Set(k, v Value) (prev Value) {
	if k == Nil {
		panicf("table set with nil key")
	}
	if atomic.LoadInt32(&m.rwlock) > 0 {
		panicf(errRace)
	}
	return m.unsafeSet(k, v)
}

func (m *RHMap) unsafeSet(k, v Value) (prev Value) {
	if m.Parent != nil && v.Type() != FUNC {
		if x := m.ParentContains(k); x != nil && x != m {
			return x.Set(k, v)
		}
	}

	if k.IsInt() {
		idx := k.Int()
		if idx >= 0 && idx < int64(len(m.items)) {
			prev, m.items[idx] = m.items[idx], v
			if v == Nil && prev != Nil {
				m.count--
			} else if v != Nil && prev == Nil {
				m.count++
			}
			return prev
		}
		if idx == int64(len(m.items)) {
			m.delHash(k)
			if v != Nil {
				m.items = append(m.items, v)
				m.count++
				return Nil
			}
			return Nil
		}
	}

	if v == Nil {
		return m.delHash(k)
	}

	if len(m.hashItems) <= 0 {
		m.hashItems = make([]hashItem, 8)
	}

	prev, _ = m.setHash(hashItem{Key: k, Val: v, Distance: 0})
	return
}

func (m *RHMap) setHash(incoming hashItem) (prev Value, growed bool) {
	num := len(m.hashItems)
	idx := int(incoming.Key.HashCode() % uint64(num))

	for idxStart := idx; ; {

		e := &m.hashItems[idx]
		if e.Key == Nil {
			m.hashItems[idx] = incoming
			m.count++
			return Nil, false
		}

		if e.Key.Equal(incoming.Key) {
			prev = e.Val
			e.Val, e.Distance = incoming.Val, incoming.Distance
			return prev, false
		}

		// Swap if the incoming item is further from its best idx.
		if e.Distance < incoming.Distance {
			incoming, m.hashItems[idx] = m.hashItems[idx], incoming
		}

		incoming.Distance++ // One step further away from best idx.
		idx = (idx + 1) % num

		// Grow if distances become big or we went all the way around.
		if float64(num)/float64(m.count) < growRate || idx == idxStart {
			m.grow(num * int(growRate*16) / 16)
			prev, _ = m.setHash(incoming)
			return prev, true
		}
	}
}

func (m *RHMap) delHash(k Value) (prev Value) {
	idx := m.findHash(k)
	if idx < 0 {
		return Nil
	}
	prev = m.hashItems[idx].Val

	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + 1
		if next >= len(m.hashItems) {
			next = 0
		}

		if next == idx { // Went all the way around.
			break
		}

		f := &m.hashItems[next]
		if f.Key == Nil || f.Distance <= 0 {
			break
		}

		f.Distance--

		m.hashItems[idx] = *f

		idx = next
	}

	m.hashItems[idx] = hashItem{}
	m.count--
	return prev
}

func (m *RHMap) Foreach(f func(k, v Value) bool) {
	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if !f(k, v) {
			return
		}
	}
}

func (m *RHMap) Next(k Value) (Value, Value) {
	atomic.AddInt32(&m.rwlock, 1)
	defer atomic.AddInt32(&m.rwlock, -1)
	nextHashPair := func(start int) (Value, Value) {
		for i := start; i < len(m.hashItems); i++ {
			if i := &m.hashItems[i]; i.Key != Nil {
				return i.Key, i.Val
			}
		}
		return Nil, Nil
	}
	if k == Nil {
		if len(m.items) == 0 {
			return nextHashPair(0)
		}
		return Int(0), m.items[0]
	}
	if k.IsInt() {
		n := k.Int()
		if n >= 0 && n < int64(len(m.items))-1 {
			for n++; n < int64(len(m.items)); n++ {
				if m.items[n] != Nil {
					return Int(n), m.items[n]
				}
			}
		}
		return nextHashPair(m.findHash(k) + 1)
	}
	idx := m.findHash(k)
	if idx < 0 {
		return Nil, Nil
	}
	return nextHashPair(idx + 1)
}

func (m *RHMap) Array() []Value {
	return m.items
}

func (m *RHMap) String() string {
	return m.Value().String()
}

func (m *RHMap) Value() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(MAP), p: unsafe.Pointer(m)}
}

func (m *RHMap) grow(newSize int) {
	tmp := RHMap{hashItems: make([]hashItem, newSize)}
	for _, e := range m.hashItems {
		if e.Key != Nil {
			e.Distance = 0
			tmp.setHash(e)
		}
	}
	m.hashItems = tmp.hashItems
}
