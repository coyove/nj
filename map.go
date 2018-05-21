package potatolang

type mapentry struct {
	k string
	v Value
}

// for elements less than mapentrySize, they will be stored in linear slice
// otherwise they will be moved to map[string]Value.
// don't make this number too big because of O(n) performance, 4 is good, 8 may also be ok
const mapentrySize = 4

// Map represents the map structure in potatolang
type Map struct {
	t []mapentry
	m map[string]Value
}

// NewMap returns a new map
// please use this function instead of new(Map)
func NewMap() *Map {
	return &Map{t: make([]mapentry, 0, 1)}
}

// Dup duplicates the map
func (m *Map) Dup(duper func(string, Value) Value) *Map {
	m2 := &Map{}
	if m.t != nil {
		m2.t = make([]mapentry, len(m.t))
		if duper == nil {
			copy(m2.t, m.t)
		} else {
			for i, x := range m.t {
				m2.t[i] = mapentry{x.k, duper(x.k, x.v)}
			}
		}
		return m2
	}
	m2.m = make(map[string]Value)
	for k, v := range m.m {
		if duper == nil {
			m2.m[k] = v
		} else {
			m2.m[k] = duper(k, v)
		}
	}
	return m2
}

// Equal compares two maps
func (m *Map) Equal(m2 *Map) bool {
	if m.t != nil && m2.t != nil {
		for _, x := range m.t {
			if v, ok := m2.Get(x.k); !ok || !v.Equal(x.v) {
				return false
			}
		}
		return true
	}
	for k, v := range m.m {
		if !m2.m[k].Equal(v) {
			return false
		}
	}
	return true
}

// Put puts a new entry into the map
func (m *Map) Put(key string, value Value) *Map {
	if m.t != nil {
		for i, x := range m.t {
			if x.k == key {
				m.t[i].v = value
				return m
			}
		}

		if len(m.t) < mapentrySize {
			m.t = append(m.t, mapentry{key, value})
			return m
		}

		m.m = make(map[string]Value)
		for _, x := range m.t {
			m.m[x.k] = x.v
		}
		m.t = nil
	}
	m.m[key] = value
	return m
}

// Get gets the corresponding value with the key
func (m *Map) Get(key string) (value Value, found bool) {
	if m.t != nil {
		for _, x := range m.t {
			if x.k == key {
				return x.v, true
			}
		}
		return NewValue(), false
	}
	v, ok := m.m[key]
	return v, ok
}

// Remove removes the key from map and return the corresponding value
func (m *Map) Remove(key string) Value {
	if m.t != nil {
		for i, x := range m.t {
			if x.k == key {
				m.t = append(m.t[:i], m.t[i+1:]...)
				return x.v
			}
		}
		return NewValue()
	}
	v := m.m[key]
	delete(m.m, key)
	return v
}

// Size returns the size of map
func (m *Map) Size() int {
	return len(m.t) + len(m.m)
}

// Values returns the underlay data
// it will either return a map, nil and nil, representing that all data are inside the map
// or it will return nil and two slices: keys and values separately
func (m *Map) Values() (map[string]Value, []string, []Value) {
	if m.m != nil {
		return m.m, nil, nil
	}

	keys := make([]string, 0, mapentrySize)
	values := make([]Value, 0, mapentrySize)

	for _, x := range m.t {
		keys = append(keys, x.k)
		values = append(values, x.v)
	}

	return nil, keys, values
}
