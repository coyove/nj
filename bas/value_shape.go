package bas

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type shape interface {
	assert(Value, string) error
	add(shape)
	String() string
}

type shaperArray struct {
	fixed  bool
	shapes []shape
}

func (sa *shaperArray) add(s shape) {
	sa.shapes = append(sa.shapes, s)
}

func (sa *shaperArray) assert(v Value, msg string) error {
	arr, ok := v.Interface().([]Value)
	if !ok {
		return fmt.Errorf("pattern %q expects array, got %v", msg, detail(v))
	}
	if sa.fixed {
		if len(arr) != len(sa.shapes) {
			return fmt.Errorf("pattern %q expects array with %d elements, got %d", msg, len(sa.shapes), len(arr))
		}
		for i, s := range sa.shapes {
			if err := s.assert(arr[i], msg); err != nil {
				return err
			}
		}
	} else {
		if len(sa.shapes) == 0 {
			return nil
		}
		if len(arr)%len(sa.shapes) != 0 {
			return fmt.Errorf("pattern %q expects array with %d*N elements, got %d", msg, len(sa.shapes), len(arr))
		}
		for i := 0; i < len(arr); i += len(sa.shapes) {
			for j, s := range sa.shapes {
				if err := s.assert(arr[i+j], msg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sa *shaperArray) String() string {
	buf := &bytes.Buffer{}
	if sa.fixed {
		buf.WriteByte('(')
	} else {
		buf.WriteByte('[')
	}
	for _, s := range sa.shapes {
		buf.WriteString(s.String())
		buf.WriteByte(' ')
	}
	if sa.fixed {
		internal.CloseBuffer(buf, ")")
	} else {
		internal.CloseBuffer(buf, "]")
	}
	return buf.String()
}

type shaperObject struct {
	key, value shape
}

func (sa *shaperObject) add(s shape) {
	if sa.key == nil {
		sa.key = s
		return
	}
	if sa.value == nil {
		sa.value = s
		return
	}
}

func (sa *shaperObject) assert(v Value, msg string) error {
	if !v.IsObject() {
		return fmt.Errorf("pattern %q expects object, got %v", msg, detail(v))
	}
	if sa.key == nil && sa.value == nil {
		return nil
	}
	if sa.key != nil && sa.value == nil {
		sp, ok := sa.key.(*shaperNative)
		if !ok {
			return fmt.Errorf("invalid pattern %q, expects object form: {prototype}", msg)
		}
		for p := v.Object(); p != nil; p = p.parent {
			if p.Name() == sp.name {
				return nil
			}
		}
		return fmt.Errorf("pattern %q expects object of prototype %v", msg, sp.name)
	}

	var err error
	v.Object().Foreach(func(k Value, v *Value) bool {
		if err = sa.key.assert(k, msg); err != nil {
			return false
		}
		if err = sa.value.assert(*v, msg); err != nil {
			return false
		}
		return true
	})
	return err
}

func (sa *shaperObject) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	if sa.key == nil && sa.value == nil {
	} else if sa.key != nil && sa.value == nil {
		buf.WriteString("@")
		buf.WriteString(sa.key.String())
	} else {
		buf.WriteString(sa.key.String())
		buf.WriteByte(':')
		buf.WriteString(sa.value.String())
	}
	buf.WriteByte('}')
	return buf.String()
}

type shaperNative struct {
	name string
}

func (sa *shaperNative) add(s shape) {
}

func (sa *shaperNative) assert(v Value, msg string) error {
	if v.Type() != typ.Native {
		return fmt.Errorf("pattern %q expects native, got %v", msg, detail(v))
	}
	if v.Native().meta.Name != sa.name {
		return fmt.Errorf("pattern %q expects %s, got %s", msg, sa.name, v.Native().meta.Name)
	}
	return nil
}

func (sa *shaperNative) String() string {
	return "@" + sa.name
}

type shaperPrimitive struct {
	verbs string
}

func (sa *shaperPrimitive) add(s shape) {
}

func (sa *shaperPrimitive) assert(v Value, msg string) error {
	if sa.verbs == "" {
		return nil
	}
	ok := false
	any := strings.IndexByte(sa.verbs, '_') >= 0 || strings.IndexByte(sa.verbs, 'v') >= 0

	switch v.Type() {
	case typ.Nil:
		ok = any
	case typ.Bool:
		ok = strings.IndexByte(sa.verbs, 'b') >= 0
	case typ.Number:
		if strings.IndexByte(sa.verbs, 'i') >= 0 {
			if !v.IsInt64() {
				return fmt.Errorf("pattern %q expects integer, got %v", msg, v)
			}
			ok = true
		} else {
			ok = strings.IndexByte(sa.verbs, 'n') >= 0
		}
	case typ.String:
		ok = strings.IndexByte(sa.verbs, 's') >= 0
	case typ.Object:
		ok = strings.IndexByte(sa.verbs, 'o') >= 0
	}

	ok = ok || any
	if !ok {
		if IsError(v) && strings.IndexByte(sa.verbs, 'E') >= 0 {
			ok = true
		}
		if IsBytes(v) && strings.IndexByte(sa.verbs, 'B') >= 0 {
			ok = true
		}
		if !ok {
			return fmt.Errorf("pattern %q expects %q, got %v", msg, sa.verbs, detail(v))
		}
	}
	return nil
}

func (sa *shaperPrimitive) String() string {
	return sa.verbs
}

type shaperOr struct {
	shapes []shape
}

func (sa *shaperOr) add(s shape) {
	sa.shapes = append(sa.shapes, s)
}

func (sa *shaperOr) assert(v Value, msg string) error {
	for _, s := range sa.shapes {
		if s.assert(v, msg) == nil {
			return nil
		}
	}
	return fmt.Errorf("pattern %q expects %v, got %v", msg, sa, detail(v))
}

func (sa *shaperOr) String() string {
	x := make([]string, len(sa.shapes))
	for i := range sa.shapes {
		x[i] = sa.shapes[i].String()
	}
	return "<" + strings.Join(x, ",") + ">"
}

func shapeNextToken(s string) (token, rest string) {
	s = strings.TrimSpace(s)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ',':
			return s[:i], s[i+1:]
		case '<', '(', '[', '{', ':', ')', ']', '}', '>', ' ':
			if i == 0 {
				return s[:1], s[1:]
			}
			return s[:i], s[i:]
		}
	}
	return s, ""
}

var shapeCache sync.Map

func Shape(s string) func(v Value) error {
	var until byte
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return func(Value) error { return nil }
	}

	old := s
	switch s[0] {
	case '(':
		until, s = ')', s[1:]
	case '[':
		until, s = ']', s[1:]
	case '{':
		until, s = '}', s[1:]
	case '<':
		until, s = '>', s[1:]
	}

	if f, ok := shapeCache.Load(old); ok {
		return f.(func(Value) error)
	}

	x := shapeScan(old, &s, until)
	if x == nil {
		return func(Value) error { return nil }
	}

	f := func(v Value) error {
		return x.assert(v, old)
	}
	shapeCache.Store(old, f)
	return f
}

func shapeScan(p string, s *string, until byte) shape {
	var sa shape
	switch until {
	case ')':
		sa = &shaperArray{fixed: true}
	case ']':
		sa = &shaperArray{}
	case '}':
		sa = &shaperObject{}
	case '>':
		sa = &shaperOr{}
	}

	for len(*s) > 0 {
		var token string
		token, *s = shapeNextToken(*s)
		if token == "" {
			panic("invalid shape form: " + strconv.Quote(p))
		}
		switch token[0] {
		case until:
			return sa
		case '(':
			sa.add(shapeScan(p, s, ')'))
		case '[':
			sa.add(shapeScan(p, s, ']'))
		case '{':
			sa.add(shapeScan(p, s, '}'))
		case '<':
			sa.add(shapeScan(p, s, '>'))
		case ':':
		case '@':
			sa2 := &shaperNative{name: token[1:]}
			if sa == nil {
				return sa2
			}
			sa.add(sa2)
		default:
			sa2 := &shaperPrimitive{verbs: token}
			if sa == nil {
				return sa2
			}
			sa.add(sa2)
		}
	}

	return sa
}

func (v Value) AssertShape(shape, msg string) Value {
	if err := Shape(shape)(v); err != nil {
		if msg == "" {
			panic(err)
		}
		panic(fmt.Errorf("%s: %v", msg, err))
	}
	return v
}
