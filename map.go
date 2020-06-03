// Robin Hood Hashmap here is an adaption of https://github.com/tidwall/rhh

// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package potatolang

const (
	loadFactor  = 0.85                      // must be above 50%
	dibBitSize  = 16                        // 0xFFFF
	hashBitSize = 64 - dibBitSize           // 0xFFFFFFFFFFFF
	maxHash     = ^uint64(0) >> dibBitSize  // max 28,147,497,671,0655
	maxDIB      = ^uint64(0) >> hashBitSize // max 65,535
)

type entry struct {
	hdib uint64 // bitfield { hash:48 dib:16 }
	k, v Value
}

func (e *entry) dib() int { return int(e.hdib & maxDIB) }

func (e *entry) hash() int { return int(e.hdib >> dibBitSize) }

func (e *entry) setDIB(dib int) { e.hdib = e.hdib>>dibBitSize<<dibBitSize | uint64(dib)&maxDIB }

func (e *entry) setHash(hash int) { e.hdib = uint64(hash)<<dibBitSize | e.hdib&maxDIB }

func makeHDIB(hash, dib int) uint64 { return uint64(hash)<<dibBitSize | uint64(dib)&maxDIB }

// hash returns a 48-bit hash for 64-bit environments, or 32-bit hash for
// 32-bit environments.
func (m *Map) hash(key Value) int { return int(key.Hash() >> dibBitSize) }

type bitmap []byte

func (b *bitmap) set(i int) { (*b)[i/8] |= (1 << (i % 8)) }

func (b *bitmap) clear(i int) { (*b)[i/8] &^= (1 << (i % 8)) }

// Map is a hashmap. Like map[string]interface{}
type Map struct {
	cap      int
	length   int
	mask     int
	growAt   int
	shrinkAt int
	buckets  []entry
	jump     bitmap
}

// New returns a new Map. Like map[string]interface{}
func New(cap int) *Map {
	m := new(Map)
	m.cap = cap
	sz := 8
	for sz < m.cap {
		sz *= 2
	}
	m.buckets = make([]entry, sz)
	m.jump = make(bitmap, sz/8)
	m.mask = len(m.buckets) - 1
	m.growAt = int(float64(len(m.buckets)) * loadFactor)
	m.shrinkAt = int(float64(len(m.buckets)) * (1 - loadFactor))
	return m
}

func (m *Map) resize(newCap int) {
	nmap := New(newCap)
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			nmap.set(m.buckets[i].hash(), m.buckets[i].k, m.buckets[i].v)
		}
	}
	cap := m.cap
	*m = *nmap
	m.cap = cap
}

// Set assigns a value to a key.
// Returns the previous value, or false when no value was assigned.
func (m *Map) Put(key, value Value) {
	if value.IsNil() {
		m.delete(key)
		return
	}

	if len(m.buckets) == 0 {
		*m = *New(0)
	}
	if m.length >= m.growAt {
		m.resize(len(m.buckets) * 2)
	}
	m.set(m.hash(key), key, value)
}

func (m *Map) set(hash int, key, value Value) {
	e := entry{makeHDIB(hash, 1), key, value}
	i := e.hash() & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			m.buckets[i] = e
			m.jump.set(i)
			m.length++
			return
		}
		if e.hash() == m.buckets[i].hash() && e.k.Equal(m.buckets[i].k) {
			m.buckets[i].v = e.v
			return
		}
		if m.buckets[i].dib() < e.dib() {
			e, m.buckets[i] = m.buckets[i], e
		}
		i = (i + 1) & m.mask
		e.setDIB(e.dib() + 1)
	}
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map) Get(key Value) Value {
	if len(m.buckets) == 0 {
		return Value{}
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return Value{}
		}
		if m.buckets[i].hash() == hash && m.buckets[i].k.Equal(key) {
			return m.buckets[i].v
		}
		i = (i + 1) & m.mask
	}
}

// Len returns the number of values in map.
func (m *Map) Len() int {
	return m.length
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) delete(key Value) {
	if m == nil || len(m.buckets) == 0 {
		return
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return
		}
		if m.buckets[i].hash() == hash && m.buckets[i].k.Equal(key) {
			m.buckets[i].setDIB(0)
			for {
				pi := i
				i = (i + 1) & m.mask
				if m.buckets[i].dib() <= 1 {
					m.buckets[pi] = entry{}
					m.jump.clear(pi)
					break
				}
				m.buckets[pi] = m.buckets[i]
				m.buckets[pi].setDIB(m.buckets[pi].dib() - 1)
			}
			m.length--
			if len(m.buckets) > m.cap && m.length <= m.shrinkAt {
				m.resize(m.length)
			}
			return
		}
		i = (i + 1) & m.mask
	}
}

func findNextKey(m *Map, i int) (Value, Value) {
	for i < len(m.buckets) {
		e := m.buckets[i]
		if !e.k.IsNil() {
			return e.k, e.v
		}

		if i/8*8 == i {
			if m.jump[i/8] == 0 {
				i += 8
			} else {
				i++
			}
		} else {
			i++
		}
	}
	return Value{}, Value{}
}

func (m *Map) Next(key Value) (Value, Value) {
	if len(m.buckets) == 0 {
		return Value{}, Value{}
	}

	if key.IsNil() { // first non-nil key
		return findNextKey(m, 0)
	}

	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return Value{}, Value{}
		}

		if m.buckets[i].hash() == hash && m.buckets[i].k.Equal(key) {
			return findNextKey(m, i+1)
		}

		i = (i + 1) & m.mask
	}
	return Value{}, Value{}
}
