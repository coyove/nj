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
	local  Map
	this   Value
}

func NewObject(size int) *Object {
	obj := newObjectInplace(Map{})
	obj.local.Init(size)
	return obj
}

func newObjectInplace(m Map) *Object {
	obj := &Object{}
	obj.local = m
	obj.this = obj.ToValue()
	obj.parent = &Proto.Object
	obj.fun = objDefaultFun
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

// Cap returns the capacity of the object.
func (m *Object) Cap() int {
	if m == nil {
		return 0
	}
	return m.local.Cap()
}

// Len returns the count of local properties in the object.
func (m *Object) Len() int {
	if m == nil {
		return 0
	}
	return m.local.Len()
}

// Clear clears all local properties in the object.
func (m *Object) Clear() {
	m.local.Clear()
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

// Find retrieves the property by 'name', returns false as the second argument if not found.
func (m *Object) Find(name Value) (v Value, exists bool) {
	if m == nil {
		return Nil, false
	}
	return m.find(name, true)
}

// Get retrieves the property by 'name'.
func (m *Object) Get(name Value) (v Value) {
	if m == nil {
		return Nil
	}
	v, _ = m.find(name, true)
	return v
}

// GetDefault retrieves the property by 'name', returns 'defaultValue' if not found.
func (m *Object) GetDefault(name, defaultValue Value) (v Value) {
	if m == nil {
		return defaultValue
	}
	if v, ok := m.find(name, true); ok {
		return v
	}
	return defaultValue
}

func (m *Object) find(k Value, setReceiver bool) (v Value, ok bool) {
	v, ok = m.local.Get(k)
	if !ok && m.parent != nil {
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

// Contains returns true if object contains property 'name', inherited properties will also be checked.
func (m *Object) Contains(name Value) bool {
	if m == nil {
		return false
	}
	found := m.local.Contains(name)
	if !found {
		found = m.parent.Contains(name)
	}
	return found
}

// HasOwnProperty returns true if 'name' is a local property in the object.
func (m *Object) HasOwnProperty(name Value) bool {
	if m == nil {
		return false
	}
	return m.local.Contains(name)
}

// Set sets a local property in the object. Inherited property with the same name will be shadowed.
func (m *Object) Set(name, v Value) (prev Value) {
	return m.local.Set(name, v)
}

// Delete deletes a local property from the object. Inherited properties are omitted and never deleted.
func (m *Object) Delete(name Value) (prev Value) {
	return m.local.Delete(name)
}

// Foreach iterates all local properties in the object, for each of them, 'f(name, &value)' will be
// called. Values are passed by pointers and it is legal to manipulate them directly in 'f'.
func (m *Object) Foreach(f func(Value, *Value) bool) {
	if m == nil {
		return
	}
	m.local.Foreach(f)
}

func (m *Object) internalNext(kv Value) Value {
	if kv == Nil {
		kv = Array(Nil, Nil)
	}
	nk, nv := m.local.FindNext(kv.Native().Get(0))
	if nk == Nil {
		return Nil
	}
	kv.Native().Set(0, nk)
	kv.Native().Set(1, nv)
	return kv
}

func (m *Object) String() string {
	return m.local.String()
}

func (m *Object) rawPrint(p io.Writer, j typ.MarshalType) {
	if m == nil {
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "null", "nil"))
		return
	}
	if j != typ.MarshalToJSON {
		if m.fun != objDefaultFun && m.fun != nil {
			internal.WriteString(p, m.funcSig())
		}
	}
	m.local.rawPrint(p, j)
}

func (m *Object) ToValue() Value {
	if m == nil {
		return Nil
	}
	return Value{v: uint64(typ.Object), p: unsafe.Pointer(m)}
}

func (m *Object) Name() string {
	if m == &Proto.Object {
		return objDefaultFun.name
	}
	if m != nil && m.fun != nil {
		if m.fun.name == objDefaultFun.name {
			return m.parent.Name()
		}
		return m.fun.name
	}
	return objDefaultFun.name
}

func (m *Object) Copy(copyData bool) *Object {
	if m == nil {
		return NewObject(0)
	}
	m2 := *m
	if copyData {
		m2.local = m.local.Copy()
	}
	if m2.fun == nil {
		// Some empty objects don't have proper structures,
		// normally they are declared directly instead of using NewObject.
		m2.fun = objDefaultFun
		m2.parent = &Proto.Object
	}
	return &m2
}

func (m *Object) Merge(src *Object) *Object {
	if src != nil && src.Len() > 0 {
		m.local.Merge(&src.local)
	}
	return m
}

func (m *Object) ToMap() Map {
	if m == nil {
		return Map{}
	}
	return m.local
}

type Map struct {
	noresize bool
	count    uint32
	items    []hashItem
}

// hashItem represents a slot in the map.
type hashItem struct {
	key, val Value
	dist     int32
	hash16   uint16
	pDeleted bool
}

func newMap(size int) *Map {
	obj := &Map{}
	obj.Init(size)
	return obj
}

// Init pre-allocates enough memory for 'count' key and clears all old data.
func (m *Map) Init(count int) *Map {
	if count > 0 {
		m.count = 0
		m.items = make([]hashItem, count*2)
	}
	return m
}

// Cap returns the capacity of the map in terms of key-value pairs, one pair is (ValueSize * 2 + 8) bytes.
func (m Map) Cap() int {
	return len(m.items)
}

// Len returns the count of keys in the map.
func (m Map) Len() int {
	return int(m.count)
}

// Clear clears all keys in the map, where already allocated memory will be reused.
func (m *Map) Clear() {
	for i := range m.items {
		m.items[i] = hashItem{}
	}
	m.count = 0
}

// Get retrieves the value by 'k', returns false as the second argument if not found.
func (m Map) Get(k Value) (v Value, exists bool) {
	if idx := m.findValue(k); idx >= 0 {
		return m.items[idx].val, true
	}
	return Nil, false
}

func (m *Map) findValue(k Value) int {
	num := len(m.items)
	if num <= 0 || k == Nil {
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

// Contains returns true if the map contains 'k'.
func (m Map) Contains(k Value) bool {
	return m.findValue(k) >= 0
}

// Set upserts a key-value pair in the map. Nil key is not allowed.
func (m *Map) Set(k, v Value) (prev Value) {
	if k == Nil {
		internal.Panic("key can't be nil")
	}
	if len(m.items) <= 0 {
		m.items = make([]hashItem, 8)
	}
	if int(m.count) >= len(m.items)*3/4 {
		m.resizeHash(len(m.items) * 2)
	}
	return m.setHash(hashItem{key: k, val: v})
}

// Delete deletes a key from the map, returns deleted value if existed
func (m *Map) Delete(k Value) (prev Value) {
	idx := m.findValue(k)
	if idx < 0 {
		return Nil
	}
	current := &m.items[idx]
	current.pDeleted = true
	current.key = Nil
	m.count--
	return current.val
}

func (m *Map) setHash(incoming hashItem) (prev Value) {
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

// Foreach iterates all keys in the map, for each of them, 'f(key, &value)' will be
// called. Values are passed by pointers and it is legal to manipulate them directly in 'f'.
// Deletions are allowed during Foreach(), but the iteration may be incomplete therefore.
func (m Map) Foreach(f func(Value, *Value) bool) {
	for i := 0; i < len(m.items); i++ {
		ip := &m.items[i]
		if ip.key != Nil && !ip.pDeleted {
			if !f(ip.key, &ip.val) {
				return
			}
		}
	}
}

func (m *Map) nextHashPair(start int) (Value, Value) {
	for i := start; i < len(m.items); i++ {
		if p := &m.items[i]; p.key != Nil && !p.pDeleted {
			return p.key, p.val
		}
	}
	return Nil, Nil
}

// FindNext finds the next key after 'k', returns nil if not found.
// The output is stable between map changes (e.g. Delete).
func (m Map) FindNext(k Value) (Value, Value) {
	if k == Nil {
		return m.nextHashPair(0)
	}
	idx := m.findValue(k)
	if idx < 0 {
		return Nil, Nil
	}
	return m.nextHashPair(idx + 1)
}

func (m Map) String() string {
	p := &bytes.Buffer{}
	m.rawPrint(p, typ.MarshalToString)
	return p.String()
}

func (m Map) rawPrint(p io.Writer, j typ.MarshalType) {
	needComma := false
	internal.WriteString(p, "{")
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

func (m Map) Copy() Map {
	m.items = append([]hashItem{}, m.items...)
	return m
}

func (m *Map) Merge(src *Map) *Map {
	if src.Len() > 0 {
		m.resizeHash((m.Len() + src.Len()) * 2)
		src.Foreach(func(k Value, v *Value) bool { m.Set(k, *v); return true })
	}
	return m
}

func (m *Map) resizeHash(newSize int) {
	if m.noresize {
		return
	}
	if newSize <= len(m.items) {
		return
	}
	tmp := Map{items: make([]hashItem, newSize)}
	for _, e := range m.items {
		if e.key != Nil {
			e.dist = 0
			tmp.setHash(e)
		}
	}
	m.items = tmp.items
}

func (m Map) density() float64 {
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

func (m Map) DebugString() string {
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

func (m Map) ToObject() *Object {
	return newObjectInplace(m)
}
