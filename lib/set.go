package lib

import (
	"fmt"

	"github.com/coyove/nj"
)

type Set struct{ m map[uint64]nj.Value }

func init() {
	nj.AddGlobalValue("set", func(env *nj.Env) {
		s := Set{m: map[uint64]nj.Value{}}
		for _, e := range env.Get(0).MustTable("").ArrayPart() {
			s.Add(e)
		}
		env.A = nj.Val(s)
	},
		"$f() value",
		"$f(a: array) value",
		"\tcreate a unique set:",
		"\t\tunique_set.add(...v: value) int",
		"\t\tunique_set.delete(v: value) value",
		"\t\tunique_set.union(set2: value)",
		"\t\tunique_set.intersect(set2: value)",
		"\t\tunique_set.subtract(set2: value)",
		"\t\tunique_set.contains(v: value) bool",
		"\t\tunique_set.values() array",
		"\t\tunique_set.size() int",
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
