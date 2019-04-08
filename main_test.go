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

	"github.com/coyove/common/rand"
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

func TestSR2(t *testing.T) {
	runFile(t, "tests/r2.txt")
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
		return 1 + 2 * 3 / 4
`)
	if err != nil {
		t.Error(err)
	}

	if len(cls.consts) != 2 || cls.consts[1].AsNumber() != 2.5 {
		t.Error("unfolding failed")
	}

	if cls.Exec(nil).AsNumber() != 2.5 {
		t.Error("exec failed")
	}
}

func TestRegisterOptimzation(t *testing.T) {
	cls, err := LoadString(`
		var a = 1, b = 2
		var c = 0
		if (0) {
			a = 2
			b = 3
			c = a + b
	}
		c = a + b
		return c
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
		return (1 / 0 + 1) * 0
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
		use "2.txt" 
		use "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/2.txt", []byte(`use "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/src/3.txt", []byte(`var a = use "1.txt"`), 0777)
	ioutil.WriteFile("tmp/src/1.txt", []byte(`use "../1.txt"`), 0777)

	_, err := LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}

	ioutil.WriteFile("tmp/1.txt", []byte(`use "1.txt"`), 0777)
	_, err = LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}
}

func TestCopyCall(t *testing.T) {
	cls, err := LoadString("var a = dup 1")
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
	if o != OP_R2K || cls.consts[a].Type() != Tnil {
		t.Fatal("error opcode 2")
	}

	cls, err = LoadString("(dup 1)")
	if err != nil {
		t.Fatal(err)
	}
	code = cls.code
	o, a, _ = op(code[0])
	if o != OP_R0K || cls.consts[a].Num() != 0.0 {
		t.Fatal("error opcode 0 3")
	}

}

func TestOverNested(t *testing.T) {
	_, err := LoadString(`
var a = 1
var foo = fun = fun = fun = fun = fun = fun = (fun () {  a = 2 })
foo()
`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadString(`
var a = 1
var foo = fun = fun = fun = fun = fun = fun = fun = (fun () {  a += 2 })
foo()
`)
	if err == nil || !strings.Contains(err.Error(), "too many levels") {
		t.FailNow()
	}
}

func TestPosVByte(t *testing.T) {
	p := posVByte{}
	p2 := [][3]uint32{}
	r := rand.New()

	for i := 0; i < 1e6; i++ {
		a, b, c := uint32(r.Uint64()), uint32(r.Uint64()), uint32(uint16(r.Uint64()))
		if r.Intn(2) == 1 {
			c = uint32(r.Intn(14))
		}
		p.appendABC(a, b, uint16(c))
		p2 = append(p2, [3]uint32{a, b, c})
	}

	i, j := 0, 0
	for i < len(p) {
		var a, b uint32
		var c uint16
		i, a, b, c = p.readABC(i)
		if [3]uint32{a, b, uint32(c)} != p2[j] {
			t.Fatal(a, b, c, p2[j])
		}
		j++
	}
}

func TestReusingTmps(t *testing.T) {
	cls, err := LoadString(`
	fun add (a, b, c) { return a + b + c }
	var d = 1
	var sum = add(1 + d, 2 + d, d + 3)
	assert sum == 9
	var sum = add(4 + d, 5 + d, 6 + d)
	assert sum == 18
`)
	if err != nil {
		t.Fatal(err)
	}
	ExecCursor(cls.lastenv, cls, 0)
	// all core libs + add + d + sum + 2 tmps
	if cls.lastenv.SSize() != len(CoreLibs)+1+1+1+2 {
		t.FailNow()
	}
}

func BenchmarkCompiling(b *testing.B) {
	buf, _ := ioutil.ReadFile("tests/string.txt")
	src := "(fun() {" + string(buf) + "})()"
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
