package lib

import (
	"fmt"

	"github.com/coyove/nj"
)

type Set struct{ m map[uint64]nj.Value }

func init() {
	nj.AddGlobalValue("set", func(env *nj.Env) {
		s := Set{m: map[uint64]nj.Value{}}
		env.Array(0).Foreach(func(k int, v nj.Value) bool {
			s.Add(v)
			return true
		})
		env.A = nj.ValueOf(s)
	},
		"$f() -> go.Set",
		"$f(a: array) -> go.Set",
		"\tcreate a unique set:",
		"\t\tgo.Set.add(...v: value) -> int",
		"\t\tgo.Set.delete(v: value) -> value",
		"\t\tgo.Set.union(set2: value)",
		"\t\tgo.Set.intersect(set2: value)",
		"\t\tgo.Set.subtract(set2: value)",
		"\t\tgo.Set.contains(v: value) -> bool",
		"\t\tgo.Set.values() -> array",
		"\t\tgo.Set.size() -> int",
	)
}

func (s Set) Add(v ...nj.Value) int {
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

func (s Set) Delete(v nj.Value) nj.Value {
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

func (s Set) Contains(v nj.Value) bool {
	_, ok := s.m[v.HashCode()]
	return ok
}

func (s Set) Values() []nj.Value {
	v := make([]nj.Value, 0, s.Size())
	for _, sv := range s.m {
		v = append(v, sv)
	}
	return v
}

func (s Set) String() string {
	return fmt.Sprintf("set(%d)", s.Size())
}
