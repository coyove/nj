package potatolang

type hash128 struct {
	a, b uint64
}

// Map represents the map structure in potatolang
// Like lua it has a linear slice and a hash table
type Map struct {
	l []Value
	m map[hash128][2]Value
}

// NewMap returns a new map
func NewMap() *Map {
	return &Map{l: make([]Value, 0), m: make(map[hash128][2]Value)}
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

	m2.m = make(map[hash128][2]Value)
	for k, v := range m.m {
		if duper != nil {
			v[1] = duper(v[0], v[1])
		}
		m2.m[k] = v
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
		if v2, ok := m2.m[k]; !ok || !v2[1].Equal(v[1]) {
			return false
		}
	}
	return true
}

// Put puts a new entry into the map
func (m *Map) Put(key Value, value Value) *Map {
	if key.ty == Tnumber {
		idx, ln := int(key.AsNumber()), len(m.l)
		if idx < ln {
			m.l[idx] = value
			return m
		} else if idx == ln {
			m.l = append(m.l, value)
			return m
		}
	}
	m.m[key.Hash()] = [2]Value{key, value}
	return m
}

func (m *Map) putIntoMap(key Value, value Value) *Map {
	m.m[key.Hash()] = [2]Value{key, value}
	return m
}

// Puts puts a new entry with a string key into the map
func (m *Map) Puts(key string, value Value) *Map {
	return m.Put(NewStringValue(key), value)
}

// Get gets the corresponding value with the key
func (m *Map) Get(key Value) (value Value, found bool) {
	if key.ty == Tnumber {
		if idx, ln := int(key.AsNumber()), len(m.l); idx < ln {
			return m.l[idx], true
		}
	}
	v, ok := m.m[key.Hash()]
	return v[1], ok
}

func (m *Map) getFromMap(key Value) (value Value, found bool) {
	v, ok := m.m[key.Hash()]
	return v[1], ok
}

// Remove removes the key from map and return the corresponding value
func (m *Map) Remove(key Value) Value {
	if key.ty == Tnumber {
		if idx, ln := int(key.AsNumber()), len(m.l); idx < ln {
			v := m.l[idx]
			m.l = append(m.l[:idx], m.l[idx+1:]...)
			return v
		}
	}
	hash := key.Hash()
	v := m.m[hash]
	delete(m.m, hash)
	return v[1]
}

// Size returns the size of map
func (m *Map) Size() int {
	return len(m.l) + len(m.m)
}
