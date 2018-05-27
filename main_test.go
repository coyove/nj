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

	t.Log(b.String())

	i := b.Exec(nil)
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

func TestArithmeticUnfold(t *testing.T) {
	cls, err := LoadString(`
		return 1 + 2 * 3 / 4
`, false)
	if err != nil {
		t.Error(err)
	}

	// 1st const: 2 * 3 = 6
	// 2nd const: 6 / 4 = 1.5
	// 3rd const: 1 + 1.5 = 2.5

	if len(cls.consts) != 3 || cls.consts[2].AsNumber() != 2.5 {
		t.Error("unfolding failed")
	}

	if cls.Exec(nil).AsNumber() != 2.5 {
		t.Error("exec failed")
	}
}
