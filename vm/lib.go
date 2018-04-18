package vm

import (
	"hash/crc32"

	"github.com/coyove/bracket/base"
)

type LibFunc struct {
	name string
	args int
	f    func(env *base.Env) base.Value

	// If we have more than 8 arguments, call ff instead
	ff func(env *base.Env) base.Value
}

func (l *LibFunc) Name() string {
	return l.name
}

func (l *LibFunc) Args() int {
	return l.args
}

var LibLookup map[string]int
var Lib []LibFunc
var LibHash uint32

func init() {
	Lib = []LibFunc{
		lib_foreach,
		lib_go,

		lib_outprint,
		lib_outprintln,
		lib_outwrite,
		lib_errprint,
		lib_errprintln,
		lib_errwrite,

		lib_syncwaitgroupnew,
		lib_syncwaitgroupadd,
		lib_syncwaitgroupdone,
		lib_syncwaitgroupwait,
		lib_syncmutexnew,
		lib_syncmutexlock,
		lib_syncmutexunlock,

		lib_numlongbits,
		lib_numtostring,
		lib_listmakelen,
		lib_mathsqrt,
	}

	LibLookup = make(map[string]int)
	c := crc32.New(crc32.IEEETable)
	for i, l := range Lib {
		LibLookup[l.name] = i
		c.Write([]byte(l.name))
	}
	LibHash = c.Sum32()
}
