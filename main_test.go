package potatolang

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"testing"
	// _ "net/http/pprof"
	"runtime"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(path, "import.txt") {
		log.Println(b.PrettyString())
	}

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
func TestSStringIndex(t *testing.T) {
	runFile(t, "tests/stringindex.txt")
}

func TestSStd(t *testing.T) {
	runFile(t, "tests/std.txt")
}

func TestSLoop(t *testing.T) {
	runFile(t, "tests/loop.txt")
}

func TestSPlaceholder(t *testing.T) {
	runFile(t, "tests/placeholder.txt")
}

func TestSImport(t *testing.T) {
	runFile(t, "tests/import.txt")
}

func TestArithmeticUnfold(t *testing.T) {
	cls, err := LoadString(`
		return 1 + 2 * 3 / 4;
`)
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
`)
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

func TestArithmeticNAN(t *testing.T) {
	cls, err := LoadString(`
		return (1 / 0 + 1) * 0;
`)
	if err != nil {
		t.Error(err)
	}

	if !math.IsNaN(cls.Exec(nil).AsNumber()) {
		t.Error("wrong answer")
	}
}

func TestImportLoop(t *testing.T) {
	os.MkdirAll("tmp/src", 0777)
	defer os.RemoveAll("tmp")

	ioutil.WriteFile("tmp/1.txt", []byte(`
		require "2.txt"; 
		require "src/3.txt";`), 0777)
	ioutil.WriteFile("tmp/2.txt", []byte(`require "src/3.txt";`), 0777)
	ioutil.WriteFile("tmp/src/3.txt", []byte(`var a = require "1.txt";`), 0777)
	ioutil.WriteFile("tmp/src/1.txt", []byte(`require "../1.txt";`), 0777)

	_, err := LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}

	ioutil.WriteFile("tmp/1.txt", []byte(`require "1.txt";`), 0777)
	_, err = LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}
}

func TestCopyCall(t *testing.T) {
	cls, err := LoadString("var a = copy(1);")
	if err != nil {
		t.Fatal(err)
	}
	code := cls.code
	o, a, _ := op(code[0])
	if o != OP_R0K || cls.consts[a].Num() != 1.0 {
		t.Fatal("error opcode 0")
	}

	o, a, _ = op(code[1])
	if o != OP_R1K || cls.consts[a].Num() != 1.0 {
		t.Fatal("error opcode 1")
	}

	o, a, _ = op(code[2])
	if o != OP_R2K || cls.consts[a].Num() != 0.0 {
		t.Fatal("error opcode 2")
	}

	cls, err = LoadString("var a = copy();")
	if err != nil {
		t.Fatal(err)
	}
	code = cls.code
	o, a, _ = op(code[0])
	if o != OP_R0K || cls.consts[a].Num() != 1.0 {
		t.Fatal("error opcode 0 1")
	}

	o, a, _ = op(code[2])
	if o != OP_R2K || cls.consts[a].Num() != 1.0 {
		t.Fatal("error opcode 2 1")
	}

	cls, err = LoadString("return copy();")
	if err != nil {
		t.Fatal(err)
	}
	code = cls.code
	o, a, _ = op(code[0])
	if o != OP_R0K || cls.consts[a].Num() != 1.0 {
		t.Fatal("error opcode 0 2")
	}

	o, a, _ = op(code[2])
	if o != OP_R2K || cls.consts[a].Num() != 2.0 {
		t.Fatal("error opcode 2 2")
	}

	cls, err = LoadString("copy();")
	if err != nil {
		t.Fatal(err)
	}
	code = cls.code
	o, a, _ = op(code[0])
	if o != OP_R0K || cls.consts[a].Num() != 0.0 {
		t.Fatal("error opcode 0 3")
	}

}

func BenchmarkCompiling(b *testing.B) {
	buf, _ := ioutil.ReadFile("tests/string.txt")
	src := "(func() {" + string(buf) + "})();"
	for i := 0; i < b.N; i++ {
		y := make([]byte, len(src)*i)
		for x := 0; x < i; x++ {
			copy(y[x*len(src):], src)
		}
		_, err := LoadString(string(y))
		if err != nil {
			b.Fatal(err)
		}
	}
}
