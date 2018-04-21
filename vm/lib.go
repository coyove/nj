package vm

import (
	"hash/crc32"
	"log"

	"github.com/coyove/bracket/base"
)

type LibFunc struct {
	name string
	args int

	// Arugments are passed by registers
	f func(env *base.Env) base.Value

	// Arguments are passed by stack
	ff func(env *base.Env) base.Value
}

func (l *LibFunc) Name() string {
	return l.name
}

func (l *LibFunc) Args() int {
	return l.args
}

func (l *LibFunc) IsFF() bool {
	return l.ff != nil
}

var LibLookup map[string]int
var Lib []LibFunc
var LibHash uint32

func init() {
	Lib = []LibFunc{
		lib_foreach,
		lib_typeof,
		lib_go,
		lib_dup,

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
		lib_stringtonum,
		lib_listmakelen,
		lib_mathsqrt,

		lib_osargs,
		lib_startprocess,
		lib_startprocessbg,
		lib_createfile,
		lib_writefile,
		lib_closefile,
	}

	LibLookup = make(map[string]int)
	c := crc32.New(crc32.IEEETable)
	for i, l := range Lib {
		if l.ff != nil && l.f != nil {
			log.Panicf("%s: can't implement f and ff both at same time", l.name)
		}

		LibLookup[l.name] = i
		c.Write([]byte(l.name))
	}
	LibHash = c.Sum32()
}
