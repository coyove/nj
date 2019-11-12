package potatolang

import (
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

type base struct {
	ptype byte
	ptag  uint32
}

type baseSlice struct {
	base
	l []Value
}

type baseStruct struct {
	base
	l IntMap
}

type baseString struct {
	base
	s []byte
}

type Closure struct {
	base
	Code       []uint32
	Pos        posVByte
	source     string
	ConstTable []Value
	Env        *Env
	ArgsCount  byte
	options    byte
	lastp      uint32
	lastenv    *Env
	native     func(env *Env) Value
}

type basePointer struct {
	base
	ptr unsafe.Pointer
}

// NewSlice creates a new map
func NewSlice() *baseSlice {
	return &baseSlice{l: make([]Value, 0)}
}

// NewSliceSize creates a new map with pre-allocated slice
func NewSliceSize(n int) *baseSlice {
	return &baseSlice{l: make([]Value, n)}
}

// Dup duplicates the map
func (m *baseSlice) Dup() *baseSlice {
	m2 := &baseSlice{}
	m2.l = make([]Value, len(m.l))
	for i, x := range m.l {
		m2.l[i] = x.Dup()
	}
	return m2
}

// Equal compares two maps
func (m *baseSlice) Equal(m2 *baseSlice) bool {
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
func (m *baseSlice) Put(idx int, value Value) *baseSlice {
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

// Get gets the corresponding value with the key
func (m *baseSlice) Get(idx int) Value {
	return m.l[idx]
}

// Remove removes the key from map and return the corresponding value
func (m *baseSlice) Remove(idx int) Value {
	if idx < len(m.l) {
		v := m.l[idx]
		m.l = append(m.l[:idx], m.l[idx+1:]...)
		return v
	}
	return Value{}
}

// Size returns the size of map
func (m *baseSlice) Size() int {
	return len(m.l)
}

func NewStruct() *baseStruct {
	return &baseStruct{}
}

// Dup duplicates the map
func (m *baseStruct) Dup() *baseStruct {
	m2 := &baseStruct{l: make(IntMap, len(m.l))}
	offset := len(m.l) / 2
	copy(m2.l, m.l[:offset])

	for i := 0; i < offset; i++ {
		if m2.l[i] == NewNumberValue(parser.FieldsField) {
			m2.l[i+offset] = m.l[i+offset]
		} else {
			m2.l[i+offset] = m.l[i+offset].Dup()
		}
	}

	return m2
}

// Equal compares two maps
func (m *baseStruct) Equal(m2 *baseStruct) bool {
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
func (m *baseStruct) Put(key string, value Value) *baseStruct {
	m.l.Add(true, NewNumberValue(parser.HashString(key)), value)
	return m
}

// Get gets the corresponding value with the key
func (m *baseStruct) Get(key Value) (Value, bool) {
	return m.l.Get(key)
}

// Size returns the size of map
func (m *baseStruct) Size() int {
	return len(m.l)
}
