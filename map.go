package potatolang

import "unsafe"

// Map represents the map structure in potatolang
// Like lua it has a linear slice and a hash table.
// We use a 128bit hash to identify the key,
// and since 128bit is large enough, we will consider it impossible to collide
// in a foreseeable period of time (though this is not a good practice).
type Map struct {
	l     []Value
	m     map[interface{}]Value
	ptr   unsafe.Pointer
	ptype byte
	ptag  uint32
}

// NewMap creates a new map
func NewMap() *Map {
	return &Map{l: make([]Value, 0)}
}

// NewMapSize creates a new map with pre-allocated slice
func NewMapSize(n int) *Map {
	return &Map{l: make([]Value, n)}
}

// Dup duplicates the map
func (m *Map) Dup(duper func(Value, Value) Value) *Map {
	m2 := &Map{}
	m2.l = make([]Value, len(m.l))
	if duper == nil {
		copy(m2.l, m.l)
	} else {
		for i, x := range m.l {
			m2.l[i] = duper(NewNumberValue(float64(i)), x)
		}
	}

	if m.m != nil {
		m2.m = make(map[interface{}]Value)
		for k, v := range m.m {
			if duper != nil {
				v = duper(NewValueFromInterface(k), v)
			}
			m2.m[k] = v
		}
	}
	return m2
}

// Equal compares two maps
func (m *Map) Equal(m2 *Map) bool {
	if len(m2.l) != len(m.l) {
		return false
	}
	for i, x := range m.l {
		if !x.Equal(m2.l[i]) {
			return false
		}
	}
	for k, v := range m.m {
		if v2, ok := m2.m[k]; !ok || !v2.Equal(v) {
			return false
		}
	}
	return true
}

// Put puts a new entry into the map
func (m *Map) Put(key Value, value Value) *Map {
	if key.Type() == Tnumber {
		idx, ln := int(key.AsNumber()), len(m.l)
		if idx < ln {
			m.l[idx] = value
			return m
		} else if idx == ln {
			m.l = append(m.l, value)
			return m
		}
	}
	if m.m == nil {
		m.m = make(map[interface{}]Value)
	}
	m.m[key.I()] = value
	return m
}

func (m *Map) putIntoMap(key Value, value Value) *Map {
	if m.m == nil {
		m.m = make(map[interface{}]Value)
	}
	m.m[key.I()] = value
	return m
}

// Puts puts a new entry with a string key into the map
func (m *Map) Puts(key string, value Value) *Map {
	return m.Put(NewStringValue(key), value)
}

// Get gets the corresponding value with the key
func (m *Map) Get(key Value) (value Value, found bool) {
	if key.Type() == Tnumber {
		if idx, ln := int(key.AsNumber()), len(m.l); idx < ln {
			return m.l[idx], true
		}
	}
	if m.m == nil {
		return Value{}, false
	}
	v, ok := m.m[key.I()]
	return v, ok
}

func (m *Map) getFromMap(key Value) (value Value, found bool) {
	v, ok := m.m[key.I()]
	return v, ok
}

// Remove removes the key from map and return the corresponding value
func (m *Map) Remove(key Value) Value {
	if key.Type() == Tnumber {
		if idx, ln := int(key.AsNumber()), len(m.l); idx < ln {
			v := m.l[idx]
			m.l = append(m.l[:idx], m.l[idx+1:]...)
			return v
		}
	}
	if m.m == nil {
		return Value{}
	}
	hash := key.I()
	v := m.m[hash]
	delete(m.m, hash)
	return v
}

// Size returns the size of map
func (m *Map) Size() int {
	return len(m.l) + len(m.m)
}
