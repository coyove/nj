package potatolang

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	// _ "net/http/pprof"
	"runtime"

	"math/rand"
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
		//t.Log(b.PrettyString())
	}

	i, i2 := b.Exec(nil)
	t.Log(i, i2)
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

func TestSReturn2(t *testing.T) {
	runFile(t, "tests/return2.txt")
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

	if len(cls.ConstTable) != 1 || cls.ConstTable[0].Num() != 2.5 {
		t.Error("unfolding failed")
	}

	if v, _ := cls.Exec(nil); v.Num() != 2.5 {
		t.Error("exec failed")
	}
}

func TestRegisterOptimzation(t *testing.T) {
	cls, err := LoadString(`
		a = 1
		b = 2
		c = 0
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

	// At the end of the if block, the op Code will be like:
	// R0 = a, R1 = b -> Add
	// But after the if block, there is another c = a + b, we can't re-use the registers R0 and R1
	// because they will not contain the value we want as the if block was not executed at all.

	if n, _ := cls.Exec(nil); n.Num() != 3 {
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

	if v, _ := cls.Exec(nil); !math.IsNaN(v.Num()) {
		t.Error("wrong answer")
	}
}

func TestImportLoop(t *testing.T) {
	os.MkdirAll("tmp/src", 0777)
	defer os.RemoveAll("tmp")

	ioutil.WriteFile("tmp/1.txt", []byte(`
		import "2.txt" 
		import "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/2.txt", []byte(`import "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/src/3.txt", []byte(`a = import "1.txt"`), 0777)
	ioutil.WriteFile("tmp/src/1.txt", []byte(`import  "../1.txt"`), 0777)

	_, err := LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}

	ioutil.WriteFile("tmp/1.txt", []byte(`import "1.txt"`), 0777)
	_, err = LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "importing each other") {
		t.Error("something wrong")
	}
}

func TestInc(t *testing.T) {
	cls, err := LoadString("a = 1; a = -1 + a; a = 2 - a")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(cls.PrettyString())

	code := cls.Code
	o, _, s := op(code[1])
	if o != OpInc || cls.ConstTable[s&0x3ff].Num() != -1.0 {
		t.Fatal("error opcode 0")
	}

	o, _, s = op(code[2])
	if o != OpSub {
		t.Fatal("error opcode 1")
	}
}

func TestOverNested(t *testing.T) {
	_, err := LoadString(`
a = 1
foo = func = func = func = func = func = (func () {  a = 2 })
foo()
`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadString(`
a = 1
foo = func = func = func = func = func = func = (func () {  a += 2 })
foo()
`)
	if err == nil || !strings.Contains(err.Error(), "too many levels") {
		t.FailNow()
	}
}

func TestPosVByte(t *testing.T) {
	p := posVByte{}
	p2 := [][3]uint32{}
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < 1e6; i++ {
		a, b, c := uint32(r.Uint64()), uint32(r.Uint64()), uint32(uint16(r.Uint64()))
		if r.Intn(2) == 1 {
			c = uint32(r.Intn(14))
		}
		p.append(a, b, uint16(c))
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

func BenchmarkCompiling(b *testing.B) {
	buf, _ := ioutil.ReadFile("tests/string.txt")
	src := "(func() {" + string(buf) + "})()"
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

func BenchmarkReturnClosure(b *testing.B) {

	b.StopTimer()
	rand.Seed(time.Now().Unix())

	buf := bytes.Buffer{}
	for i := 0; i < 1000; i++ {
		if rand.Intn(2) == 0 {
			x := make([]byte, rand.Intn(16)+16)
			rand.Read(x)
			buf.WriteString(fmt.Sprintf("v%d = \"%s\"\n", i, hex.EncodeToString(x)))
		} else {
			buf.WriteString(fmt.Sprintf("v%d = %d\n", i, rand.Int()))
		}
	}

	src := "(func() {" + buf.String() + "})()"
	cls, err := LoadString(src)
	if err != nil || cls == nil {
		b.Fatal(cls, err)
	}

	b.StartTimer()

	l := cls.lastenv
	for i := 0; i < b.N; i++ {
		cls.lastenv = l
		cls.Exec(nil)
	}
}
