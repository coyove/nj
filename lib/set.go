package lib

import (
	"fmt"

	"github.com/coyove/nj/bas"
)

type Set struct{ m map[uint64]bas.Value }

func init() {
	bas.Globals.SetMethod("set", func(env *bas.Env) {
		s := Set{m: map[uint64]bas.Value{}}
		env.Array(0).ForeachIndex(func(k int, v bas.Value) bool {
			s.Add(v)
			return true
		})
		env.A = bas.ValueOf(s)
	}, "$f() -> go.Set\n"+
		"$f(a: array) -> go.Set\n"+
		"\tcreate a unique set:\n"+
		"\t\tgo.Set.add(...v: value) -> int\n"+
		"\t\tgo.Set.delete(v: value) -> value\n"+
		"\t\tgo.Set.union(set2: value)\n"+
		"\t\tgo.Set.intersect(set2: value)\n"+
		"\t\tgo.Set.subtract(set2: value)\n"+
		"\t\tgo.Set.contains(v: value) -> bool\n"+
		"\t\tgo.Set.values() -> array\n"+
		"\t\tgo.Set.size() -> int")
}

func (s Set) Add(v ...bas.Value) int {
	c := 0
	for _, v := range v {
		hash := v.HashCode()
		if _, ok := s.m[hash]; !ok {
			c++
		}
		s.m[hash] = v
	}
	return c
}

func (s Set) Delete(v bas.Value) bas.Value {
	hash := v.HashCode()
	v = s.m[hash]
	delete(s.m, hash)
	return v
}

func (s Set) Union(s2 Set) {
	for hash, v := range s2.m {
		s.m[hash] = v
	}
}

func (s Set) Intersect(s2 Set) {
	for hash := range s.m {
		if _, ok := s2.m[hash]; !ok {
			delete(s.m, hash)
		}
	}
}

func (s Set) Subtract(s2 Set) {
	for hash := range s.m {
		if _, ok := s2.m[hash]; ok {
			delete(s.m, hash)
		}
	}
}

func (s Set) Size() int {
	return len(s.m)
}

func (s Set) Contains(v bas.Value) bool {
	_, ok := s.m[v.HashCode()]
	return ok
}

func (s Set) Values() []bas.Value {
	v := make([]bas.Value, 0, s.Size())
	for _, sv := range s.m {
		v = append(v, sv)
	}
	return v
}

func (s Set) String() string {
	return fmt.Sprintf("set(%d)", s.Size())
}
