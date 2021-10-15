package lib

import (
	"fmt"

	"github.com/coyove/script"
)

type Set struct {
	m map[uint64]struct{}
	v []script.Value
	p *script.Program
}

func init() {
	script.AddGlobalValue("set", func(env *script.Env) {
		s := &Set{m: map[uint64]struct{}{}, p: env.Global}
		for _, e := range env.Get(0).MustMap("").ArrayPart() {
			s.Add(e)
		}
		env.A = script.Val(s)
	},
		"set() => unique_set",
		"set({ e1, e2, ..., en }) => unique_set",
		"\tcreate a unique set:",
		"\t\tunique_set.Add(value) => bool",
		"\t\tunique_set.Exists(value) => bool",
		"\t\tunique_set.Values() => { e1, e2, ... en }",
		"\t\tunique_set.Size() => int",
	)
}

func (s *Set) Add(v script.Value) bool {
	hash := v.HashCode()
	_, exist := s.m[hash]
	if exist {
		return false
	}
	s.m[hash] = struct{}{}
	s.v = append(s.v, v)
	return true
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Exists(v script.Value) bool {
	_, ok := s.m[v.HashCode()]
	return ok
}

func (s *Set) Values() []script.Value {
	return s.v
}

func (s *Set) String() string {
	return fmt.Sprintf("set(%d)", s.Size())
}
