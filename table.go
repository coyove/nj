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
	"bytes"
	"unsafe"

	"github.com/coyove/script/typ"
)

type Table struct {
	parent    *Table
	hashCount uint32
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

func NewTable(size int) *Table {
	if size >= 8 {
		size *= 2
	}
	return &Table{hashItems: make([]hashItem, int64(size))}
}

func (m *Table) Size() int { return len(m.items) + len(m.hashItems) }

func (m *Table) Len() int { return int(m.count) + int(m.hashCount) }

func (m *Table) MapLen() int { return int(m.hashCount) }

func (m *Table) ArrayLen() int { return int(m.count) }

// Clear clears the table, where already allocated memory will be reused.
func (m *Table) Clear() { m.ClearArray(); m.ClearMap() }

func (m *Table) ClearArray() { m.items = m.items[:0]; m.count = 0 }

func (m *Table) ClearMap() { m.hashItems = m.hashItems[:0]; m.hashCount = 0 }

func (m *Table) Parent() *Table { return m.parent }

func (m *Table) SetParent(m2 *Table) { m.parent = m2 }

func (m *Table) SetFirstParent(m2 *Table) {
	if m.parent != nil {
		m2 = m2.Copy()
		m2.SetFirstParent(m.parent)
	}
	m.parent = m2
}

func (m *Table) GetString(k string) (v Value) {
	return m.getImpl(Str(k), true)
}

// Get retrieves the value for a given key.
func (m *Table) Get(k Value) (v Value) {
	return m.getImpl(k, true)
}

func (m *Table) getImpl(k Value, funcRecv bool) (v Value) {
	if k == Nil {
		return Nil
	}
	if k.IsInt() {
		if idx := k.Int(); idx >= 0 && idx < int64(len(m.items)) {
			v = m.items[idx]
			goto FINAL
		}
	}
	if idx := m.findHash(k); idx >= 0 {
		v = m.hashItems[idx].Val
	} else if m.parent != nil {
		v = m.parent.getImpl(k, funcRecv)
	}
FINAL:
	if funcRecv && v.Type() == typ.Func {
		f := *v.Func()
		f.Receiver = m.Value()
		v = f.Value()
	}
	return v
}

func (m *Table) findHash(k Value) int {
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

func (m *Table) Contains(k Value) bool {
	if k == Nil {
		return false
	}
	if k.IsInt() {
		if idx := k.Int(); idx >= 0 && idx < int64(len(m.items)) {
			return m.items[idx] != Nil
		}
	}
	return m.findHash(k) >= 0
}

func (m *Table) SetString(k string, v Value) (prev Value) {
	return m.Set(Str(k), v)
}

// Set inserts or updates a key/val pair into the Map. If val == Nil, then key will get deleted
func (m *Table) Set(k, v Value) (prev Value) {
	if k == Nil {
		panicf("table set with nil key")
	}

	if m.parent != nil && !m.Contains(k) {
		for p := m.parent; p != nil; p = p.parent {
			if p.Contains(k) {
				return p.Set(k, v)
			}
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

func (m *Table) RawSet(k, v Value) (prev Value) {
	old := m.parent
	m.parent = nil
	prev, m.parent = m.Set(k, v), old
	return prev
}

func (m *Table) RawGet(k Value) (v Value) {
	old := m.parent
	m.parent = nil
	v, m.parent = m.Get(k), old
	return v
}

func (m *Table) setHash(incoming hashItem) (prev Value, growed bool) {
	num := len(m.hashItems)
	idx := int(incoming.Key.HashCode() % uint64(num))

	for idxStart := idx; ; {

		e := &m.hashItems[idx]
		if e.Key == Nil {
			m.hashItems[idx] = incoming
			m.hashCount++
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
		if num < 8 {
			if idx == idxStart {
				m.resizeHash(num + 1)
				prev, _ = m.setHash(incoming)
				return prev, true
			}
		} else {
			if int(m.hashCount) >= num/2 || idx == idxStart {
				m.resizeHash(num*2 + 1)
				prev, _ = m.setHash(incoming)
				return prev, true
			}
		}
	}
}

func (m *Table) delHash(k Value) (prev Value) {
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
	m.hashCount--
	return prev
}

func (m *Table) Foreach(f func(k, v Value) bool) {
	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if !f(k, v) {
			return
		}
	}
}

func (m *Table) Next(k Value) (Value, Value) {
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

func (m *Table) ArrayPart() []Value {
	return m.items
}

func (m *Table) MapPart() map[Value]Value {
	g := make(map[Value]Value, len(m.hashItems))
	for _, i := range m.hashItems {
		if i.Key != Nil {
			g[i.Key] = i.Val
		}
	}
	return g
}

func (m *Table) String() string {
	return m.Value().String()
}

func (m *Table) rawPrint(p *bytes.Buffer, lv int, j, showParent bool) {
	if m.hashCount == 0 {
		p.WriteString(ifstr(j, "[", "{"))
		for _, a := range m.ArrayPart() {
			a.toString(p, lv+1, j)
			p.WriteString(",")
		}
		closeBuffer(p, ifstr(j, "]", "}"))
	} else {
		p.WriteString("{")
		for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
			k.toString(p, lv+1, j)
			p.WriteString(ifstr(j, ":", "="))
			v.toString(p, lv+1, j)
			p.WriteString(",")
		}
		closeBuffer(p, "}")
	}
	if m.parent != nil && showParent {
		p.WriteString("^")
		m.parent.rawPrint(p, lv+1, j, true)
	}
}

func (m *Table) Value() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(typ.Table), p: unsafe.Pointer(m)}
}

func (m *Table) Name() string {
	var n string
	old := m.parent
	m.parent = nil
	if f := m.GetString("__name"); f != Nil {
		n = f.String()
	} else {
		n = "table"
	}
	m.parent = old
	if m.parent != nil {
		n += "^" + m.parent.Name()
	}
	return n
}

func (m *Table) New() *Table {
	if f := m.GetString("__new"); f.Type() == typ.Func {
		res, err := f.Func().Call()
		panicErr(err)
		return res.MustTable("table.__new")
	}
	return m
}

func (m *Table) Copy() *Table {
	m2 := *m
	m2.hashItems = append([]hashItem{}, m.hashItems...)
	m2.items = append([]Value{}, m.items...)
	return &m2
}

func (m *Table) Merge(src *Table, kvs ...Value) *Table {
	if src == nil {
		m.resizeHash((m.Len() + len(kvs)) * 2)
	} else {
		m.resizeHash((m.Len() + src.Len() + len(kvs)) * 2)
		src.Foreach(func(k, v Value) bool { m.Set(k, v); return true })
	}
	for i := 0; i < len(kvs)/2*2; i += 2 {
		m.Set(kvs[i], kvs[i+1])
	}
	return m
}

func (m *Table) resizeHash(newSize int) {
	if newSize < len(m.hashItems) {
		panic("resizeHash: invalid size")
	}
	if newSize == len(m.hashItems) {
		return
	}
	tmp := Table{hashItems: make([]hashItem, newSize)}
	for _, e := range m.hashItems {
		if e.Key != Nil {
			e.Distance = 0
			tmp.setHash(e)
		}
	}
	m.hashItems = tmp.hashItems
}
