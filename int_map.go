package potatolang

import (
	"sort"
	"unsafe"
)

type IntMap []Value

func (m IntMap) Get(k Value) (v Value, ok bool) {
	offset := len(m) / 2
	j := offset

	if j > 0 && m[0] == k {
		return m[offset], true
	}

	if j > 1 && m[1] == k {
		return m[1+offset], true
	}
	//	return 0, false

	for i := 2; i < j; {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		k2 := &m[h]
		if *k2 == k {
			return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(k2)) + uintptr(offset)*SizeOfValue)), true
			// return m.values[h+offset], true
		}

		if k2.AsNumber() < k.AsNumber() {
			i = h + 1
		} else {
			j = h
		}
	}

	return Value{}, false
}

func (m *IntMap) Add(create bool, k, v Value) bool {
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

func (m *IntMap) BatchSet(kv []Value) {
	*m = append([]Value{}, kv...)
	sort.Sort(kvSwapper(*m))
}
