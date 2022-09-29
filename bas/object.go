package bas

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Object struct {
	parent *Object
	fun    *funcbody
	count  int
	items  []hashItem
	this   Value
}

// hashItem represents an entry in the object.
type hashItem struct {
	key, val Value
	dist     int32
	hash16   uint16
	pDeleted bool
}

func NewObject(size int) *Object {
	size *= 2
	obj := &Object{}
	if size > 0 {
		obj.items = make([]hashItem, size)
	}
	obj.this = obj.ToValue()
	obj.parent = &Proto.Object
	obj.fun = objEmptyFunc
	return obj
}

func NewNamedObject(name string, size int) *Object {
	return NewObject(size).setName(name)
}

func (m *Object) setName(name string) *Object {
	m.fun = &funcbody{name: name, native: func(*Env) {}}
	return m
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

func (m *Object) HasPrototype(proto *Object) bool {
	for ; m != nil; m = m.parent {
		if m == proto {
			return true
		}
	}
	return false
}

// Size returns the total slots in the object, one slot is 40-bytes on 64bit platform.
func (m *Object) Size() int {
	if m == nil {
		return 0
	}
	return len(m.items)
}

// Len returns the count of local properties in the object.
func (m *Object) Len() int {
	if m == nil {
		return 0
	}
	return int(m.count)
}

// Clear clears all local properties in the object, where already allocated memory will be reused.
func (m *Object) Clear() {
	for i := range m.items {
		m.items[i] = hashItem{}
	}
	m.count = 0
}

// SetProp sets property by string 'name', which is short for Set(Str(name), v).
func (m *Object) SetProp(name string, v Value) *Object {
	m.Set(Str(name), v)
	return m
}

// AddMethod binds function 'fun' to property 'name' in the object, making 'fun' a method of the object.
// This differs from 'Set(name, Func(name, fun))' because the latter one,
// as not being a method, can't use 'this' argument when called.
func (m *Object) AddMethod(name string, fun func(*Env)) *Object {
	f := Func(m.Name()+"."+name, fun)
	f.Object().fun.method = true
	m.Set(Str(name), f)
	return m
}

// Get retrieves the property by 'name'.
func (m *Object) Get(name Value) (v Value) {
	if m == nil || name == Nil {
		return Nil
	}
	v, _ = m.find(name, true)
	return v
}

// GetDefault retrieves the property by 'name', returns 'defaultValue' if not found.
func (m *Object) GetDefault(name, defaultValue Value) (v Value) {
	if m == nil || name == Nil {
		return defaultValue
	}
	if v, ok := m.find(name, true); ok {
		return v
	}
	return defaultValue
}

func (m *Object) find(k Value, setReceiver bool) (v Value, ok bool) {
	if idx := m.findValue(k); idx >= 0 {
		v, ok = m.items[idx].val, true
	} else if m.parent != nil {
		v, ok = m.parent.find(k, false)
	}
	if setReceiver && v.IsObject() {
		if obj := v.Object(); obj.fun.method {
			f := obj.Copy(false)
			f.this = m.ToValue()
			v = f.ToValue()
		}
	}
	return
}

func (m *Object) findValue(k Value) int {
	num := len(m.items)
	if num <= 0 {
		return -1
	}
	idx := int(k.HashCode() % uint32(num))
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

// Contains returns true if object contains 'name', inherited properties will also be checked.
func (m *Object) Contains(name Value) bool {
	if m == nil || name == Nil {
		return false
	}
	found := m.findValue(name) >= 0
	if !found {
		found = m.parent.Contains(name)
	}
	return found
}

// HasOwnProperty returns true if 'name' is a local property in the object.
func (m *Object) HasOwnProperty(name Value) bool {
	if m == nil || name == Nil {
		return false
	}
	return m.findValue(name) >= 0
}

// Set sets a local property in the object. Inherited property with the same name will be shadowed.
func (m *Object) Set(name, v Value) (prev Value) {
	if name == Nil {
		internal.Panic("property name can't be nil")
	}
	if len(m.items) <= 0 {
		m.items = make([]hashItem, 8)
	}
	if int(m.count) >= len(m.items)*3/4 {
		resizeHash(m, len(m.items)*2)
	}
	return m.setHash(hashItem{key: name, val: v})
}

// Delete deletes a local property from the object. Inherited properties are omitted and never deleted.
func (m *Object) Delete(name Value) (prev Value) {
	if name == Nil {
		return Nil
	}
	idx := m.findValue(name)
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
	idx := int(incoming.key.HashCode() % uint32(num))

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

// Foreach iterates all local properties in the object, for each of them, 'f(name, &value)' will be
// called. Values are passed by pointers and it is legal to manipulate them directly in 'f'.
func (m *Object) Foreach(f func(Value, *Value) bool) {
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

func (m *Object) internalNext(kv Value) Value {
	if kv == Nil {
		kv = Array(Nil, Nil)
	}
	nk, nv := m.FindNext(kv.Native().Get(0))
	if nk == Nil {
		return Nil
	}
	kv.Native().Set(0, nk)
	kv.Native().Set(1, nv)
	return kv
}

// FindNext finds the next property after 'name', the output will be stable between object changes, e.g. Delete()
func (m *Object) FindNext(name Value) (Value, Value) {
	if m == nil {
		return Nil, Nil
	}
	if name == Nil {
		return m.nextHashPair(0)
	}
	idx := m.findValue(name)
	if idx < 0 {
		return Nil, Nil
	}
	return m.nextHashPair(idx + 1)
}

func (m *Object) String() string {
	p := &bytes.Buffer{}
	m.rawPrint(p, typ.MarshalToString)
	return p.String()
}

func (m *Object) rawPrint(p io.Writer, j typ.MarshalType) {
	if m == nil {
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "null", "nil"))
		return
	}
	needComma := false
	if j == typ.MarshalToJSON {
		if m.fun == objEmptyFunc {
			internal.WriteString(p, `{`)
		} else if m.fun != nil {
			internal.WriteString(p, `{"<f>":"`)
			internal.WriteString(p, m.funcSig())
			internal.WriteString(p, `"`)
			needComma = true
		}
	} else {
		if m.fun != objEmptyFunc && m.fun != nil {
			internal.WriteString(p, m.funcSig())
		}
		internal.WriteString(p, "{")
	}
	m.Foreach(func(k Value, v *Value) bool {
		internal.WriteString(p, internal.IfStr(needComma, ",", ""))
		k.Stringify(p, j.NoRec())
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, ":", "="))
		v.Stringify(p, j.NoRec())
		needComma = true
		return true
	})
	internal.WriteString(p, "}")
}

func (m *Object) ToValue() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(typ.Object), p: unsafe.Pointer(m)}
}

func (m *Object) Name() string {
	if m == &Proto.Object {
		return objEmptyFunc.name
	}
	if m != nil && m.fun != nil {
		if m.fun.name == objEmptyFunc.name {
			return m.parent.Name()
		}
		return m.fun.name
	}
	return objEmptyFunc.name
}

func (m *Object) Copy(copyData bool) *Object {
	if m == nil {
		return NewObject(0)
	}
	m2 := *m
	if copyData {
		m2.items = append([]hashItem{}, m.items...)
	}
	if m2.fun == nil {
		// Some empty Objects don't have proper structures,
		// normally they are declared directly instead of using NewObject.
		m2.fun = objEmptyFunc
		m2.parent = &Proto.Object
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

func (m *Object) density() float64 {
	num := len(m.items)
	if num <= 0 || m.count <= 0 {
		return math.NaN()
	}

	var maxRun int
	for i := 0; i < num; {
		if m.items[i].key == Nil {
			i++
			continue
		}
		run := 1
		for i++; i < num; i++ {
			if m.items[i].key != Nil {
				run++
			} else {
				break
			}
		}
		if run > maxRun {
			maxRun = run
		}
	}
	return float64(maxRun) / (float64(num) / float64(m.count))
}

func (m *Object) DebugString() string {
	p := bytes.Buffer{}
	for idx, i := range m.items {
		p.WriteString(strconv.Itoa(idx) + ":")
		if i.pDeleted {
			p.WriteString("\t" + strings.Repeat(".", int(i.dist)) + "deleted\n")
		} else if i.key == Nil {
			p.WriteString("\t-\n")
		} else {
			at := i.key.HashCode() % uint32(len(m.items))
			if i.dist > 0 {
				p.WriteString(fmt.Sprintf("^%d", at))
			}
			p.WriteString("\t" + strings.Repeat(".", int(i.dist)) + fmt.Sprintf("%v\n", i.key))
		}
	}
	return p.String()
}
