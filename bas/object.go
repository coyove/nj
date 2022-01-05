package bas

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Object struct {
	parent *Object
	fun    *Function
	count  int64
	items  []hashItem
	this   Value
}

// hashItem represents an entry in the object.
type hashItem struct {
	key, val Value
	dist     int32
	pDeleted bool
}

func (i hashItem) String() string {
	if i.pDeleted {
		return fmt.Sprintf("deleted(+%d)", i.dist)
	}
	return fmt.Sprintf("%v=%v(+%d)", i.key, i.val, i.dist)
}

func NewObject(preallocateSize int) *Object {
	preallocateSize *= 2
	obj := &Object{}
	if preallocateSize > 0 {
		obj.items = make([]hashItem, preallocateSize)
	}
	obj.this = obj.ToValue()
	obj.parent = &ObjectProto
	return obj
}

func NamedObject(name string, preallocateSize int) *Object {
	obj := NewObject(preallocateSize)
	obj.fun = &Function{Name: name, Dummy: true, Native: func(*Env) {}, obj: obj}
	return obj
}

func (m *Object) Prototype() *Object {
	if m == nil {
		return nil
	}
	return m.parent
}

func (m *Object) SetPrototype(m2 *Object) *Object {
	m.parent = m2
	return m
}

func (m *Object) IsPrototype(proto *Object) bool {
	for ; m != nil; m = m.parent {
		if m == proto {
			return true
		}
	}
	return false
}

func (m *Object) Size() int {
	if m == nil {
		return 0
	}
	return len(m.items)
}

func (m *Object) Len() int {
	if m == nil {
		return 0
	}
	return int(m.count)
}

// Clear clears all keys in the object, where already allocated memory will be reused.
func (m *Object) Clear() {
	for i := range m.items {
		m.items[i] = hashItem{}
	}
	m.count = 0
}

func (m *Object) Prop(k string) (v Value) {
	return m.Find(Str(k))
}

func (m *Object) SetProp(k string, v Value) *Object {
	m.Set(Str(k), v)
	return m
}

func (m *Object) SetMethod(k string, v func(*Env), d string) *Object {
	m.Set(Str(k), Func(k, v, d))
	return m
}

// Get retrieves the value for a given key locally.
func (m *Object) Get(k Value) (v Value) {
	if m == nil || k == Nil {
		return Nil
	}
	return m.find(k, false)
}

// Find finds the value for a given key (including prototypes)
func (m *Object) Find(k Value) (v Value) {
	if m == nil || k == Nil {
		return Nil
	}
	return m.find(k, true)
}

func (m *Object) find(k Value, findPrototype bool) (v Value) {
	if idx := m.findHash(k); idx >= 0 {
		v = m.items[idx].val
	} else if findPrototype && m.parent != nil {
		v = m.parent.find(k, findPrototype)
	}
	if m.parent != Proto.StaticObject && v.IsObject() && v.Object().IsCallable() {
		f := v.Object().Copy(false)
		f.this = m.ToValue()
		v = f.ToValue()
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
		if e.key == Nil {
			if !e.pDeleted {
				return -1
			}
		}

		if e.key.Equal(k) {
			return idx
		}

		idx = (idx + 1) % num
		if idx == idxStart {
			return -1
		}
	}
}

func (m *Object) Contains(k Value, includePrototypes bool) bool {
	if m == nil || k == Nil {
		return false
	}
	found := m.findHash(k) >= 0
	if !found && includePrototypes {
		found = m.parent.Contains(k, true)
	}
	return found
}

// Set upserts a key-value pair into the object
func (m *Object) Set(k, v Value) (prev Value) {
	if k == Nil {
		internal.Panic("object set with nil key")
	}
	if len(m.items) <= 0 {
		m.items = make([]hashItem, 8)
	}
	if int(m.count) >= len(m.items)/2+1 {
		resizeHash(m, len(m.items)*2)
	}
	return m.setHash(hashItem{key: k, val: v})
}

// Delete deletes a key-value pair from the object
func (m *Object) Delete(k Value) (prev Value) {
	if k == Nil {
		internal.Panic("object delete with nil key")
	}
	idx := m.findHash(k)
	if idx < 0 {
		return Nil
	}
	current := &m.items[idx]
	current.pDeleted = true
	current.key = Nil
	m.count--
	return current.val
}

func (m *Object) setHash(incoming hashItem) (prev Value) {
	num := len(m.items)
	idx := int(incoming.key.HashCode() % uint64(num))

	for idxStart := idx; ; {
		e := &m.items[idx]
		if e.pDeleted {
			// Shift the following keys forward
			this := idx
			for startIdx := this; ; {
				next := (this + 1) % num
				if m.items[next].dist > 0 {
					m.items[this] = m.items[next]
					m.items[this].dist--
					this = next
					if this != startIdx {
						continue
					}
				}
				break
			}
			m.items[this] = hashItem{}
			continue
		}

		if e.key == Nil {
			m.items[idx] = incoming
			m.count++
			return Nil
		}

		if e.key.Equal(incoming.key) {
			prev = e.val
			e.val, e.dist, e.pDeleted = incoming.val, incoming.dist, false
			return prev
		}

		// Swap if the incoming item is further from its best idx.
		if e.dist < incoming.dist {
			incoming, m.items[idx] = m.items[idx], incoming
		}

		incoming.dist++ // One step further away from best idx.
		idx = (idx + 1) % num

		if idx == idxStart {
			if internal.IsDebug() {
				fmt.Println(m.items)
			}
			panic("object space not enough")
		}
	}
}

func (m *Object) Foreach(f func(k Value, v *Value) bool) {
	if m == nil {
		return
	}
	for i := 0; i < len(m.items); i++ {
		ip := &m.items[i]
		if ip.key != Nil && !ip.pDeleted {
			if !f(ip.key, &ip.val) {
				return
			}
		}
	}
}

func (m *Object) nextHashPair(start int) (Value, Value) {
	for i := start; i < len(m.items); i++ {
		if p := &m.items[i]; p.key != Nil && !p.pDeleted {
			return p.key, p.val
		}
	}
	return Nil, Nil
}

func (m *Object) Next(k Value) (Value, Value) {
	if m == nil {
		return Nil, Nil
	}
	if k == Nil {
		return m.nextHashPair(0)
	}
	idx := m.findHash(k)
	if idx < 0 {
		return Nil, Nil
	}
	return m.nextHashPair(idx + 1)
}

func (m *Object) String() string {
	p := &bytes.Buffer{}
	m.rawPrint(p, typ.MarshalToString, false)
	return p.String()
}

func (m *Object) rawPrint(p io.Writer, j typ.MarshalType, showProto bool) {
	if m == nil {
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "null", "nil"))
		return
	}
	needComma := false
	if m.fun != nil {
		if j == typ.MarshalToJSON {
			internal.WriteString(p, `{"<f>":"`)
			internal.WriteString(p, m.fun.String())
			internal.WriteString(p, `"`)
			needComma = true
		} else {
			internal.WriteString(p, m.fun.String())
			if m.count == 0 {
				return
			}
			internal.WriteString(p, "{")
		}
	} else {
		if j == typ.MarshalToString {
			internal.WriteString(p, m.Name())
		}
		internal.WriteString(p, "{")
	}
	m.Foreach(func(k Value, v *Value) bool {
		internal.WriteString(p, internal.IfStr(needComma, ",", ""))
		k.Stringify(p, j)
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, ":", "="))
		v.Stringify(p, j)
		needComma = true
		return true
	})
	if m.parent != nil && showProto && m.parent != &ObjectProto {
		internal.WriteString(p, internal.IfStr(needComma, ",", ""))
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "\"<proto>\":", "<proto>="))
		m.parent.rawPrint(p, j, true)
	}
	internal.WriteString(p, "}")
}

func (m *Object) ToValue() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(typ.Object), p: unsafe.Pointer(m)}
}

func (m *Object) Name() string {
	if m != nil {
		if m.fun != nil {
			return m.fun.Name
		}
		if m.parent != nil {
			return m.parent.Name()
		}
	}
	return "object"
}

func (m *Object) Copy(copyData bool) *Object {
	if m == nil {
		return NewObject(0)
	}
	m2 := *m
	if copyData {
		m2.items = append([]hashItem{}, m.items...)
	}
	if m.fun != nil {
		c2 := *m.fun
		m2.fun = &c2
		m2.fun.obj = &m2
	}
	return &m2
}

func (m *Object) Merge(src *Object) *Object {
	if src != nil && src.Len() > 0 {
		resizeHash(m, (m.Len()+src.Len())*2)
		src.Foreach(func(k Value, v *Value) bool { m.Set(k, *v); return true })
	}
	return m
}

func (m *Object) Callable() Function {
	if m == nil || m.fun == nil || m.fun.Dummy {
		return Function{Dummy: true}
	}
	return *m.fun
}

func (m *Object) IsCallable() bool {
	if m == nil {
		return false
	}
	return m.fun != nil && !m.fun.Dummy
}

var resizeHash = func(m *Object, newSize int) {
	if newSize <= len(m.items) {
		return
	}
	tmp := Object{items: make([]hashItem, newSize)}
	for _, e := range m.items {
		if e.key != Nil {
			e.dist = 0
			tmp.setHash(e)
		}
	}
	m.items = tmp.items
}