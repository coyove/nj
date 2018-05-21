package potatolang

import (
	"flag"
	"testing"
	// _ "net/http/pprof"
	"runtime"
)

var lineinfo = flag.Bool("li", false, "toggle lineinfo")

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path, *lineinfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(Prettify(b))

	i := Exec(NewTopEnv(), b)
	t.Log(i.I())
}

func TestSMain(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	runFile(t, "tests/test.txt")
}

func TestSBuiltin(t *testing.T) {
	runFile(t, "tests/builtin.txt")
}

func TestSString(t *testing.T) {
	runFile(t, "tests/string.txt")
}
