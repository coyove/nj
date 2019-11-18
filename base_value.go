package potatolang

import (
	"sort"
	"unsafe"

	"github.com/coyove/potatolang/hash50"
)

type Base struct {
	ptype byte
	ptag  uint32
}

type Slice struct {
	Base
	l []Value
}

type Struct struct {
	Base
	l treeMap
}

type String struct {
	Base
	s []byte
}

type Pointer struct {
	Base
	ptr unsafe.Pointer
}

// NewSlice creates a new map
func NewSlice() *Slice {
	return &Slice{l: make([]Value, 0)}
}

// NewSliceSize creates a new map with pre-allocated slice
func NewSliceSize(n int) *Slice {
	return &Slice{l: make([]Value, n)}
}

// Dup duplicates the map
func (m *Slice) Dup() *Slice {
	m2 := &Slice{}
	m2.l = make([]Value, len(m.l))
	for i, x := range m.l {
		m2.l[i] = x.Dup()
	}
	return m2
}

// Equal compares two maps
func (m *Slice) Equal(m2 *Slice) bool {
	if len(m2.l) != len(m.l) {
		return false
	}
	for i, x := range m.l {
		if !x.Equal(m2.l[i]) {
			return false
		}
	}
	return true
}

// Put puts a new entry into the map
func (m *Slice) Put(idx int, value Value) *Slice {
	ln := len(m.l)
	if idx < ln {
		m.l[idx] = value
	} else if idx == ln {
		m.l = append(m.l, value)
	} else {
		panic("index out of range")
	}
	return m
}

// hashGet gets the corresponding value with the key
func (m *Slice) Get(idx int) Value {
	return m.l[idx]
}

// Remove removes the key from map and return the corresponding value
func (m *Slice) Remove(idx int) Value {
	if idx < len(m.l) {
		v := m.l[idx]
		m.l = append(m.l[:idx], m.l[idx+1:]...)
		return v
	}
	return Value{}
}

// Size returns the size of map
func (m *Slice) Size() int {
	return len(m.l)
}

func NewStruct() *Struct {
	return &Struct{}
}

// Dup duplicates the map
func (m *Struct) Dup() *Struct {
	m2 := &Struct{l: make(treeMap, len(m.l))}
	offset := len(m.l) / 2
	copy(m2.l, m.l[:offset])

	for i := 0; i < offset; i++ {
		m2.l[i+offset] = m.l[i+offset].Dup()
	}

	return m2
}

// Equal compares two maps
func (m *Struct) Equal(m2 *Struct) bool {
	if len(m2.l) != len(m.l) {
		return false
	}
	for i, x := range m.l {
		if !x.Equal(m2.l[i]) {
			return false
		}
	}
	return true
}

// Put puts a new entry into the map
func (m *Struct) Put(key string, value Value) *Struct {
	m.l.Add(true, NewNumberValue(float64(hash50.HashString(key))), value)
	return m
}

// hashGet gets the corresponding value with the key
func (m *Struct) Get(key string) (Value, bool) {
	return m.hashGet(NewNumberValue(float64(hash50.HashString(key))))
}

func (m *Struct) MustGet(key string) Value {
	v, ok := m.Get(key)
	if !ok {
		panic(key + " not found")
	}
	return v
}

func (m *Struct) hashGet(key Value) (Value, bool) {
	return m.l.Get(key)
}

// Size returns the size of map
func (m *Struct) Size() int {
	return len(m.l) / 2
}

type treeMap []Value

func (m treeMap) Get(k Value) (v Value, ok bool) {
	offset := len(m) / 2
	j := offset

	if j > 0 && m[0] == k {
		return m[offset], true
	}

	if j > 1 && m[1] == k {
		return m[1+offset], true
	}
	//	return 0, false

	kn := k.AsNumber()
	for i := 2; i < j; {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		k2 := &m[h]
		if *k2 == k {
			return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(k2)) + uintptr(offset)*SizeOfValue)), true
			// return m.values[h+offset], true
		}

		if k2.AsNumber() < kn {
			i = h + 1
		} else {
			j = h
		}
	}

	if int(kn) < len(m)/2 {
		return m[int(kn)+len(m)/2], false
	}
	return Value{}, false
}

func (m *treeMap) Add(create bool, k, v Value) bool {
	offset := len(*m) / 2
	i, j := 0, offset

	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if (*m)[h] == k {
			(*m)[h+offset] = v
			return true
		}

		if (*m)[h].AsNumber() < k.AsNumber() {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	if !create {
		return false
	}

	*m = append(*m, Value{}, Value{})
	copy((*m)[i+offset+2:], (*m)[i+offset:])
	(*m)[i+offset+1] = v
	copy((*m)[i+1:], (*m)[i:i+offset])
	(*m)[i] = k

	return false
}

type kvSwapper []Value

func (s kvSwapper) Len() int {
	return len(s) / 2
}

func (s kvSwapper) Less(i, j int) bool {
	return (s)[i].AsNumber() < (s)[j].AsNumber()
}

func (s kvSwapper) Swap(i, j int) {
	(s)[i], (s)[j] = (s)[j], (s)[i]
	i, j = i+len(s)/2, j+len(s)/2
	(s)[i], (s)[j] = (s)[j], (s)[i]
}

func (m *treeMap) BatchSet(kv []Value) {
	*m = append([]Value{}, kv...)
	sort.Sort(kvSwapper(*m))
}
