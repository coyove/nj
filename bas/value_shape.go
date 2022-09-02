package bas

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type shape interface {
	assert(Value, string)
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

func (sa *shaperArray) assert(v Value, msg string) {
	arr, ok := v.Interface().([]Value)
	if !ok {
		panic(fmt.Sprintf("pattern %q expects array, got %v", msg, detail(v)))
	}
	if sa.fixed {
		if len(arr) != len(sa.shapes) {
			panic(fmt.Sprintf("pattern %q expects array with %d elements, got %d", msg, len(sa.shapes), len(arr)))
		}
		for i, s := range sa.shapes {
			s.assert(arr[i], msg)
		}
	} else {
		if len(sa.shapes) == 0 {
			return
		}
		if len(arr)%len(sa.shapes) != 0 {
			panic(fmt.Sprintf("pattern %q expects array with %dN elements, got %d", msg, len(sa.shapes), len(arr)))
		}
		for i := 0; i < len(arr); i += len(sa.shapes) {
			for j, s := range sa.shapes {
				s.assert(arr[i+j], msg)
			}
		}
	}
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
	panic(fmt.Errorf("expects object shape in form of '{key:value}', got too many shapes"))
}

func (sa *shaperObject) assert(v Value, msg string) {
	if !v.IsObject() {
		panic(fmt.Sprintf("pattern %q expects object, got %v", msg, detail(v)))
	}
	if sa.key == nil && sa.value == nil {
		return
	}
	v.Object().Foreach(func(k Value, v *Value) bool {
		if sa.key != nil {
			sa.key.assert(k, msg)
		}
		if sa.value != nil {
			sa.value.assert(*v, msg)
		}
		return true
	})
}

func (sa *shaperObject) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	if sa.key == nil && sa.value == nil {
	} else {
		if sa.key == nil {
			buf.WriteString("any")
		} else {
			buf.WriteString(sa.key.String())
		}
		buf.WriteByte(':')
		if sa.value == nil {
			buf.WriteString("any")
		} else {
			buf.WriteString(sa.value.String())
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

type shaperNative struct {
	name string
}

func (sa *shaperNative) add(s shape) {
}

func (sa *shaperNative) assert(v Value, msg string) {
	if v.Type() != typ.Native {
		panic(fmt.Sprintf("pattern %q expects native, got %v", msg, detail(v)))
	}
	if v.Native().meta.Name != sa.name {
		panic(fmt.Sprintf("pattern %q expects %s, got %s", msg, sa.name, v.Native().meta.Name))
	}
}

func (sa *shaperNative) String() string {
	return "@" + sa.name
}

type shaperPrimitive struct {
	verbs string
}

func (sa *shaperPrimitive) add(s shape) {
}

func (sa *shaperPrimitive) assert(v Value, msg string) {
	if sa.verbs == "" {
		return
	}
	ok := false
	switch v.Type() {
	case typ.Nil:
		ok = strings.IndexByte(sa.verbs, '_') >= 0 || strings.IndexByte(sa.verbs, 'v') >= 0
	case typ.Bool:
		ok = strings.IndexByte(sa.verbs, 'b') >= 0
	case typ.Number:
		if strings.IndexByte(sa.verbs, 'i') >= 0 {
			if !v.IsInt64() {
				panic(fmt.Sprintf("pattern %q expects integer, got %v", msg, v))
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
	if !ok {
		if IsError(v) && strings.IndexByte(sa.verbs, 'E') >= 0 {
			ok = true
		}
		if IsBytes(v) && strings.IndexByte(sa.verbs, 'B') >= 0 {
			ok = true
		}
		if !ok {
			panic(fmt.Sprintf("pattern %q expects %q, got %v", msg, sa.verbs, detail(v)))
		}
	}
}

func (sa *shaperPrimitive) String() string {
	return sa.verbs
}

func shapeNextToken(s string) (token, rest string) {
	s = strings.TrimSpace(s)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ',':
			return s[:i], s[i+1:]
		case '(', '[', '{', ':', ')', ']', '}', ' ':
			if i == 0 {
				return s[:1], s[1:]
			}
			return s[:i], s[i:]
		}
	}
	return s, ""
}

func Shape(s string) func(v Value) {
	var until byte
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return func(Value) {}
	}

	old := s
	switch s[0] {
	case '(':
		until, s = ')', s[1:]
	case '[':
		until, s = ']', s[1:]
	case '{':
		until, s = '}', s[1:]
	}

	x := shapeScan(&s, until)
	if x == nil {
		return func(Value) {}
	}

	return func(v Value) {
		x.assert(v, old)
	}
}

func shapeScan(s *string, until byte) shape {
	var sa shape
	switch until {
	case ')':
		sa = &shaperArray{fixed: true}
	case ']':
		sa = &shaperArray{}
	case '}':
		sa = &shaperObject{}
	}

	for len(*s) > 0 {
		var token string
		token, *s = shapeNextToken(*s)
		if token == "" {
			return nil
		}
		switch token[0] {
		case until:
			return sa
		case '(':
			sa.add(shapeScan(s, ')'))
		case '[':
			sa.add(shapeScan(s, ']'))
		case '{':
			sa.add(shapeScan(s, '}'))
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
