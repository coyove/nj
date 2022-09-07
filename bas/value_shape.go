package bas

import (
	"bytes"
	"fmt"
	"io"
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
	var err error
	v.Object().Foreach(func(k Value, v *Value) bool {
		if err = sa.key.assert(k, msg); err != nil {
			return false
		}
		if sa.value != nil {
			if err = sa.value.assert(*v, msg); err != nil {
				return false
			}
		}
		return true
	})
	return err
}

func (sa *shaperObject) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	if sa.key == nil && sa.value == nil {
	} else {
		buf.WriteString(sa.key.String())
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

type shaperPrototype struct {
	name string
}

func (sa *shaperPrototype) add(s shape) {
}

func (sa *shaperPrototype) assert(v Value, msg string) error {
	switch v.Type() {
	case typ.Native:
		if v.Native().meta.Name == sa.name {
			return nil
		}
		for p := v.Native().meta.Proto; p != nil; p = p.parent {
			if p.Name() == sa.name {
				return nil
			}
		}
		return fmt.Errorf("pattern %q expects native of prototype/name %v", msg, sa.name)
	case typ.Object:
		for p := v.Object(); p != nil; p = p.parent {
			if p.Name() == sa.name {
				return nil
			}
		}
		return fmt.Errorf("pattern %q expects object of prototype %v", msg, sa.name)
	default:
		return fmt.Errorf("pattern %q expects native or object, got %v", msg, detail(v))
	}
}

func (sa *shaperPrototype) String() string {
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
		ok = strings.IndexByte(sa.verbs, 'N') >= 0
	case typ.Bool:
		ok = strings.IndexByte(sa.verbs, 'b') >= 0
	case typ.Number:
		if strings.IndexByte(sa.verbs, 'i') >= 0 {
			ok = v.IsInt64()
		} else {
			ok = strings.IndexByte(sa.verbs, 'n') >= 0
		}
	case typ.String:
		ok = strings.IndexByte(sa.verbs, 's') >= 0
	case typ.Object:
		ok = strings.IndexByte(sa.verbs, 'o') >= 0
	}

	ok = ok || any

	if !ok && IsError(v) && strings.IndexByte(sa.verbs, 'E') >= 0 {
		ok = true
	}
	if !ok && IsBytes(v) && strings.IndexByte(sa.verbs, 'B') >= 0 {
		ok = true
	}
	if !ok && strings.IndexByte(sa.verbs, 'R') >= 0 {
		switch v.Interface().(type) {
		case string, []byte, io.Reader, *Object:
			ok = true
		}
	}
	if !ok && strings.IndexByte(sa.verbs, 'W') >= 0 {
		switch v.Interface().(type) {
		case io.Writer, *Object:
			ok = true
		}
	}
	if !ok && strings.IndexByte(sa.verbs, 'C') >= 0 {
		switch v.Interface().(type) {
		case io.Closer, *Object:
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("pattern %q expects %v, got %v", msg, sa, detail(v))
	}
	return nil
}

func (sa *shaperPrimitive) String() string {
	var buf []string
	for _, b := range sa.verbs {
		switch b {
		case 'i':
			buf = append(buf, "int")
		case 'n':
			buf = append(buf, "number")
		case 'b':
			buf = append(buf, "bool")
		case 's':
			buf = append(buf, "string")
		case 'o':
			buf = append(buf, "object")
		case 'N':
			buf = append(buf, "nil")
		case 'E':
			buf = append(buf, "@error")
		case 'B':
			buf = append(buf, "@bytes")
		case 'R':
			buf = append(buf, "Reader")
		case 'W':
			buf = append(buf, "Writer")
		case 'C':
			buf = append(buf, "Closer")
		default:
			buf = append(buf, "any")
		}
	}
	return "<" + strings.Join(buf, ",") + ">"
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

func NewShape(s string) func(v Value) error {
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
			sa2 := &shaperPrototype{name: token[1:]}
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

func TestShapeFast(v Value, shape string) (err error) {
	switch shape[0] {
	case '(', '[', '{', '<':
		err = NewShape(shape)(v)
	case '@':
		sp := shaperPrototype{name: shape[1:]}
		err = sp.assert(v, shape)
	default:
		sp := shaperPrimitive{verbs: shape}
		err = sp.assert(v, shape)
	}
	return
}

func (v Value) AssertShape(shape, msg string) Value {
	if err := TestShapeFast(v, shape); err != nil {
		if msg == "" {
			panic(err)
		}
		panic(fmt.Errorf("%s: %v", msg, err))
	}
	return v
}
