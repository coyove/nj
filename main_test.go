package potatolang

import (
	"flag"
	"log"
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

	log.Println(b.String())

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

func TestSStd(t *testing.T) {
	runFile(t, "tests/std.txt")
}

func TestSLoop(t *testing.T) {
	runFile(t, "tests/loop.txt")
}

func TestArithmeticUnfold(t *testing.T) {
	cls, err := LoadString(`
		return 1 + 2 * 3 / 4;
`, false)
	if err != nil {
		t.Error(err)
	}

	// 0: nil
	// 1st const: 2 * 3 = 6
	// 2nd const: 6 / 4 = 1.5
	// 3rd const: 1 + 1.5 = 2.5

	if len(cls.consts) != 4 || cls.consts[3].AsNumber() != 2.5 {
		t.Error("unfolding failed")
	}

	if cls.Exec(nil).AsNumber() != 2.5 {
		t.Error("exec failed")
	}
}

func TestRegisterOptimzation(t *testing.T) {
	cls, err := LoadString(`
		var a = 1, b = 2;
		var c = 0;
		if (0) {
			a = 2;
			b = 3;
			c = a + b;
	}
		c = a + b;
		return c;
`, false)
	if err != nil {
		t.Error(err)
	}

	// At the end of the if block, the op code will be like:
	// R0 = a, R1 = b -> Add
	// But after the if block, there is another c = a + b, we can't re-use the registers R0 and R1
	// because they will not contain the value we want as the if block was not executed at all.

	if n := cls.Exec(nil).AsNumber(); n != 3 {
		t.Error("exec failed:", n, cls)
	}
}
