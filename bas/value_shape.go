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
	assert(Value) error
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

func (sa *shaperArray) assert(v Value) error {
	if !v.IsArray() {
		return fmt.Errorf("%v expects array, got %v", sa, v.simple())
	}
	arr := v.Native()
	if sa.fixed {
		if arr.Len() != len(sa.shapes) {
			return fmt.Errorf("%v expects array with %d elements, got %d", sa, len(sa.shapes), arr.Len())
		}
		for i, s := range sa.shapes {
			if err := s.assert(arr.Get(i)); err != nil {
				return err
			}
		}
	} else {
		if len(sa.shapes) == 0 {
			return nil
		}
		if arr.Len()%len(sa.shapes) != 0 {
			return fmt.Errorf("%v expects array with a multiple of %d elements, got %d", sa, len(sa.shapes), arr.Len())
		}
		for i := 0; i < arr.Len(); i += len(sa.shapes) {
			for j, s := range sa.shapes {
				if err := s.assert(arr.Get(i + j)); err != nil {
					return fmt.Errorf("array shape %v, value at index %d: %v", sa, i+j, err)
				}
			}
		}
	}
	return nil
}

func (sa *shaperArray) String() string {
	buf := bytes.NewBufferString(internal.IfStr(sa.fixed, "(", "["))
	for _, s := range sa.shapes {
		buf.WriteString(s.String())
		buf.WriteByte(',')
	}
	internal.CloseBuffer(buf, internal.IfStr(sa.fixed, ")", "]"))
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

func (sa *shaperObject) assert(v Value) error {
	if !v.IsObject() {
		return fmt.Errorf("%v expects object, got %v", sa, v.simple())
	}
	if sa.key == nil && sa.value == nil {
		return nil
	}
	var err error
	v.Object().Foreach(func(k Value, v *Value) bool {
		if err = sa.key.assert(k); err != nil {
			err = fmt.Errorf("object shape %v key error: %v", sa, err)
			return false
		}
		if sa.value != nil {
			if err = sa.value.assert(*v); err != nil {
				err = fmt.Errorf("object shape %v value error: %v", sa, err)
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

func (sa *shaperPrototype) assert(v Value) error {
	return assertShapePrototype(v, sa.name)
}

func assertShapePrototype(v Value, name string) error {
	switch v.Type() {
	case typ.Native:
		if v.Native().meta.Name == name {
			return nil
		}
		for p := v.Native().meta.Proto; p != nil; p = p.parent {
			if p.Name() == name {
				return nil
			}
		}
		return fmt.Errorf("expects native of prototype/name %v, got %v", name, v.Native().meta.Name)
	case typ.Object:
		for p := v.Object(); p != nil; p = p.parent {
			if p.Name() == name {
				return nil
			}
		}
		return fmt.Errorf("expects object of prototype %v, got %v", name, v.Object().Name())
	default:
		return fmt.Errorf("expects native or object, got %v", v.simple())
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

func (sa *shaperPrimitive) assert(v Value) error {
	return assertShapePrimitive(v, sa.verbs)
}

func assertShapePrimitive(v Value, verbs string) error {
	if verbs == "" || verbs == "_" || verbs == "v" {
		return nil
	}

	bm := [128]bool{}
	for i := 0; i < len(verbs); i++ {
		bm[verbs[i]] = true
	}

	ok := false
	switch v.Type() {
	case typ.Nil:
		ok = bm['N']
	case typ.Bool:
		ok = bm['b']
	case typ.Number:
		if bm['i'] {
			ok = v.IsInt64()
		} else {
			ok = bm['n']
		}
	case typ.String:
		ok = bm['s'] || bm['R'] || bm['G']
	case typ.Object:
		ok = bm['o'] || bm['R'] || bm['W'] || bm['C']
	case typ.Native:
		if bm['E'] {
			ok = IsError(v)
		} else if bm['B'] {
			ok = IsBytes(v)
		} else if bm['C'] {
			_, ok = v.Native().Unwrap().(io.Closer)
		} else if bm['W'] {
			_, ok = v.Native().Unwrap().(io.Writer)
		} else if bm['R'] {
			x := v.Native().Unwrap()
			_, ok = x.(io.Reader)
			if !ok {
				_, ok = x.([]byte)
			}
		}
	}
	if !ok {
		return fmt.Errorf("%v can't match %v", &shaperPrimitive{verbs}, v.simple())
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
		case 'G':
			buf = append(buf, "goto")
		default:
			buf = append(buf, "any")
		}
	}
	if len(buf) == 1 {
		return buf[0]
	}
	return "<" + strings.Join(buf, ",") + ">"
}

type shaperOr struct {
	shapes []shape
}

func (sa *shaperOr) add(s shape) {
	sa.shapes = append(sa.shapes, s)
}

func (sa *shaperOr) assert(v Value) error {
	for _, s := range sa.shapes {
		if s.assert(v) == nil {
			return nil
		}
	}
	return fmt.Errorf("%v can't match %v", sa, v.simple())
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
			if i == 0 {
				return shapeNextToken(s[1:])
			}
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
	s = strings.TrimSpace(s)
	if f, ok := shapeCache.Load(s); ok {
		return f.(func(Value) error)
	}

	x := buildShape(s)
	if x == nil {
		return func(Value) error { return nil }
	}

	f := func(v Value) error {
		if err := x.assert(v); err != nil {
			return fmt.Errorf("%q: %v", s, err)
		}
		return nil
	}
	shapeCache.Store(s, f)
	return f
}

func buildShape(s string) shape {
	var until byte
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
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
	return shapeScan(old, &s, until)
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
		err = assertShapePrototype(v, shape[1:])
	default:
		err = assertShapePrimitive(v, shape)
	}
	return
}

func (v Value) AssertShape(shape, msg string) Value {
	if err := TestShapeFast(v, shape); err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
	return v
}

func (v Value) AssertNumber(msg string) Value {
	if v.Type() != typ.Number {
		internal.Panic("%s: expects number, got %v", msg, v.simple())
	}
	return v
}

func (v Value) AssertString(msg string) string {
	if v.Type() != typ.String {
		internal.Panic("%s: expects string, got %v", msg, v.simple())
	}
	return v.Str()
}

func (v Value) AssertObject(msg string) *Object {
	if v.Type() != typ.Object {
		internal.Panic("%s: expects object, got %v", msg, v.simple())
	}
	return v.Object()
}
