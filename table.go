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

package nj

import (
	"bytes"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Object struct {
	parent   *Object
	count    int64
	items    []hashItem
	receiver Value
	callable *FuncBody
}

// hashItem represents an entry in the Map.
type hashItem struct {
	Key, Val Value
	Distance int // How far item is from its best position.
}

func NewObject(size int) *Object {
	if size >= 8 {
		size *= 2
	}
	return &Object{items: make([]hashItem, int64(size))}
}

func (m *Object) Parent() *Object { return m.parent }

func (m *Object) SetParent(m2 *Object) *Object { m.parent = m2; return m }

func (m *Object) Size() int { return len(m.items) }

func (m *Object) Len() int { return int(m.count) }

// Clear clears the table, where already allocated memory will be reused.
func (m *Object) Clear() {
	m.items = m.items[:0]
	m.count = 0
}

func (m *Object) SetFirstParent(m2 *Object) *Object {
	if m.parent != nil {
		m2 = m2.Copy()
		m2.SetFirstParent(m.parent)
	}
	m.parent = m2
	return m
}

func (m *Object) Gets(k string) (v Value) {
	return m.getImpl(Str(k), true)
}

// Get retrieves the value for a given key.
func (m *Object) Get(k Value) (v Value) {
	return m.getImpl(k, true)
}

func (m *Object) getImpl(k Value, recv bool) (v Value) {
	if k == Nil {
		return Nil
	}
	if idx := m.findHash(k); idx >= 0 {
		v = m.items[idx].Val
	} else if m.parent != nil {
		v = m.parent.getImpl(k, recv)
	}
	if recv && v.IsObject() {
		f := *v.Object()
		f.receiver = m.Value()
		v = f.Value()
	}
	return v
}

func (m *Object) findHash(k Value) int {
	num := len(m.items)
	if num <= 0 {
		return -1
	}
	idx := int(k.HashCode() % uint64(num))
	idxStart := idx

	for {
		e := &m.items[idx]
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

func (m *Object) Contains(k Value) bool {
	if k == Nil {
		return false
	}
	return m.findHash(k) >= 0
}

func (m *Object) Sets(k string, v Value) (prev Value) {
	return m.Set(Str(k), v)
}

// Set upserts a key-value pair into the table
func (m *Object) Set(k, v Value) (prev Value) {
	if k == Nil {
		internal.Panic("table set with nil key")
	}

	if m.parent != nil && !m.Contains(k) {
		for p := m.parent; p != nil; p = p.parent {
			if p.Contains(k) {
				return p.Set(k, v)
			}
		}
	}

	if len(m.items) <= 0 {
		m.items = make([]hashItem, 8)
	}
	prev, _ = m.setHash(hashItem{Key: k, Val: v, Distance: 0})
	return
}

// Delete deletes a key-value pair from the table
func (m *Object) Delete(k Value) (prev Value) {
	if k == Nil {
		internal.Panic("table delete with nil key")
	}
	return m.delHash(k)
}

func (m *Object) RawSet(k, v Value) (prev Value) {
	old := m.parent
	m.parent = nil
	prev, m.parent = m.Set(k, v), old
	return prev
}

func (m *Object) RawGet(k Value) (v Value) {
	old := m.parent
	m.parent = nil
	v, m.parent = m.Get(k), old
	return v
}

func (m *Object) setHash(incoming hashItem) (prev Value, growed bool) {
	num := len(m.items)
	idx := int(incoming.Key.HashCode() % uint64(num))

	for idxStart := idx; ; {

		e := &m.items[idx]
		if e.Key == Nil {
			m.items[idx] = incoming
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
			incoming, m.items[idx] = m.items[idx], incoming
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
			if int(m.count) >= num/2 || idx == idxStart {
				m.resizeHash(num*2 + 1)
				prev, _ = m.setHash(incoming)
				return prev, true
			}
		}
	}
}

func (m *Object) delHash(k Value) (prev Value) {
	idx := m.findHash(k)
	if idx < 0 {
		return Nil
	}
	prev = m.items[idx].Val

	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + 1
		if next >= len(m.items) {
			next = 0
		}

		if next == idx { // Went all the way around.
			break
		}

		f := &m.items[next]
		if f.Key == Nil || f.Distance <= 0 {
			break
		}

		f.Distance--

		m.items[idx] = *f

		idx = next
	}

	m.items[idx] = hashItem{}
	m.count--
	return prev
}

func (m *Object) Foreach(f func(k, v Value) bool) {
	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if !f(k, v) {
			return
		}
	}
}

func (m *Object) Next(k Value) (Value, Value) {
	nextHashPair := func(start int) (Value, Value) {
		for i := start; i < len(m.items); i++ {
			if i := &m.items[i]; i.Key != Nil {
				return i.Key, i.Val
			}
		}
		return Nil, Nil
	}
	if k == Nil {
		return nextHashPair(0)
	}
	idx := m.findHash(k)
	if idx < 0 {
		return Nil, Nil
	}
	return nextHashPair(idx + 1)
}

func (m *Object) Map() map[Value]Value {
	g := make(map[Value]Value, len(m.items))
	for _, i := range m.items {
		if i.Key != Nil {
			g[i.Key] = i.Val
		}
	}
	return g
}

func (m *Object) String() string {
	return m.Value().String()
}

func (m *Object) rawPrint(p *bytes.Buffer, lv int, j, showParent bool) {
	if !j {
		p.WriteString(m.Name())
	}
	p.WriteString("{")
	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		k.toString(p, lv+1, j)
		p.WriteString(ifstr(j, ":", "="))
		v.toString(p, lv+1, j)
		p.WriteString(",")
	}
	closeBuffer(p, "}")
	if m.parent != nil && showParent {
		p.WriteString("^")
		m.parent.rawPrint(p, lv+1, j, true)
	}
}

func (m *Object) Value() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(typ.Object), p: unsafe.Pointer(m)}
}

func (m *Object) Name() string {
	if m.callable != nil {
		return m.callable.Name
	}
	return "object"
}

// Copy returns a new table with a copy of dataset
func (m *Object) Copy() *Object {
	m2 := *m
	m2.items = append([]hashItem{}, m.items...)
	return &m2
}

// Shadow returns a new table with preallocated memory large enough to hold the original dataset
func (m *Object) Shadow() *Object {
	return &Object{
		parent: m.parent,
		items:  make([]hashItem, len(m.items)),
	}
}

func (m *Object) Call(args ...Value) (v1 Value, err error) {
	if m.callable != nil {
		defer internal.CatchErrorFuncCall(&err, m.callable.Name)
	}
	return m.Apply(args...), nil
}

func (m *Object) Apply(args ...Value) Value {
	if m.callable == nil {
		return m.Value()
	}
	if m.receiver != Nil {
		return m.callable.Apply(m.receiver, args...)
	}
	return m.callable.Apply(m.Value(), args...)
}

func (m *Object) Merge(src *Object, kvs ...Value) *Object {
	if src == nil {
		m.resizeHash((m.Len() + len(kvs)) * 2)
	} else {
		m.resizeHash((m.Len() + src.Len() + len(kvs)) * 2)
		src.Foreach(func(k, v Value) bool { m.Set(k, renameFuncName(k, v)); return true })
		m.callable = src.callable
	}
	for i := 0; i < len(kvs)/2*2; i += 2 {
		m.Set(kvs[i], renameFuncName(kvs[i], kvs[i+1]))
	}
	return m
}

func (m *Object) resizeHash(newSize int) {
	if newSize < len(m.items) {
		panic("resizeHash: invalid size")
	}
	if newSize == len(m.items) {
		return
	}
	tmp := Object{items: make([]hashItem, newSize)}
	for _, e := range m.items {
		if e.Key != Nil {
			e.Distance = 0
			tmp.setHash(e)
		}
	}
	m.items = tmp.items
}

func renameFuncName(k, v Value) Value {
	if v.IsObject() {
		if cls := v.Object().callable; cls != nil && cls.Name == internal.UnnamedFunc {
			cls.Name = k.String()
		}
	}
	return v
}
