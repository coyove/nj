package potatolang

type tablekey struct {
	g   Value
	str string
}

type Table struct {
	a  []Value
	m  map[tablekey]Value
	mt *Table
}

func (t *Table) Put(k, v Value) {
	switch k.Type() {
	case StringType:
		if t.m == nil {
			t.m = make(map[tablekey]Value)
		}
		if v.IsNil() {
			delete(t.m, tablekey{str: k.AsString()})
		} else {
			t.m[tablekey{str: k.AsString()}] = v
		}
	case NumberType:
		idx := k.AsNumber()
		if float64(int(idx)) == idx {
			if idx >= 0 && int(idx) < len(t.a) {
				t.a[int(idx)] = v
				if v.IsNil() {
					t.Compact()
				}
				return
			}
			if int(idx) == len(t.a) && !v.IsNil() {
				t.a = append(t.a, v)
				return
			}
		}
		fallthrough
	default:
		if t.m == nil {
			t.m = make(map[tablekey]Value)
		}
		if v.IsNil() {
			delete(t.m, tablekey{g: k})
		} else {
			t.m[tablekey{g: k}] = v
		}
	}
}

func (t *Table) Gets(k string) Value {
	return t.m[tablekey{str: k}]
}

func (t *Table) Get(k Value) Value {
	switch k.Type() {
	case StringType:
		return t.m[tablekey{str: k.AsString()}]
	case NumberType:
		idx := k.AsNumber()
		if float64(int(idx)) == idx {
			if idx >= 0 && int(idx) < len(t.a) {
				return t.a[int(idx)]
			}
		}
		fallthrough
	default:
		return t.m[tablekey{g: k}]
	}
}

func (t *Table) Compact() {
	if len(t.a) < 32 { // small array, no need to compact
		return
	}

	holes := 0
	best := struct {
		ratio float64
		index int
	}{}

	for i, v := range t.a {
		if !v.IsNil() {
			continue
		}

		holes++
		ratio := float64(i-holes) / float64(i)

		if ratio > best.ratio {
			best.ratio = ratio
			best.index = i
		}
	}

	if holes == 0 {
		return
	}

	if best.ratio < 0.5 {
		for i, v := range t.a {
			t.m[tablekey{g: NewNumberValue(float64(i))}] = v
		}
		t.a = nil
		return
	}

	for i := best.index + 1; i < len(t.a); i++ {
		t.m[tablekey{g: NewNumberValue(float64(i))}] = t.a[i]
	}
	t.a = t.a[:best.index+1]
}
