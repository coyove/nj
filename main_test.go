package script

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
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
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

	i, i2 := b.Call()
	t.Log(i, i2)
}

func TestSMain(t *testing.T) {
	runFile(t, "tests/test.txt")
}

func TestSR2(t *testing.T) {
	runFile(t, "tests/r2.txt")
}

func TestReturnFunction(t *testing.T) {
	cls, _ := LoadString(`
a = 1
return function(n) 
a+=n
return a
end
`)
	v, _ := cls.Call()
	if v, _ := v.Function().Call(Int(10)); v.Int() != 11 {
		t.Fatal(v)
	}

	if v, _ := v.Function().Call(Int(100)); v.Int() != 111 {
		t.Fatal(v)
	}
}

func TestTailCallPanic(t *testing.T) {
	cls, _ := LoadString(`
x = 0
function foo()
x+=1
if x == 1e6 then assert(false) end
foo()
end
foo()
`)
	if s := cls.PrettyString(); !strings.Contains(s, "tail-call") {
		t.Fatal(s)
	}

	_, _, err := cls.PCall()
	if err == nil {
		t.FailNow()
	}
	if len(err.Error()) > 1e6 { // error too long, which means tail call is not effective
		t.Fatal(len(err.Error()))
	}
}

func TestArithmeticUnfold(t *testing.T) {
	cls, err := LoadString(`
		return 1 + 2 * 3 / 4
`)
	if err != nil {
		t.Error(err)
	}

	if len(cls.ConstTable) != 1 || cls.ConstTable[0].Float() != 2.5 {
		t.Error("unfolding failed")
	}

	if v, _ := cls.Call(); v.Float() != 2.5 {
		t.Error("exec failed")
	}
}

func TestPCallStackSize(t *testing.T) {
	cls, _ := LoadString(`
_, err = pcall(function() 
local a, b, c = 1, 2, 3
assert(a, b, c)
return a
end)
print(err)
assert(match(err.Error(), "overflow" ))
`)
	WithMaxStackSize(cls, 7+int64(len(g)))
	cls.Call()

	cls, _ = LoadString(`
a = ""
for i = 1,1e3 do
a = a .. i
end
return a
`)
	WithMaxStackSize(cls, int64(len(g))+10) // 10: a small value
	res, _, err := cls.PCall()
	if !strings.Contains(err.Error(), "string overflow") {
		t.Fatal(res)
	}
}

func TestRegisterOptimzation(t *testing.T) {
	cls, err := LoadString(`
		a = 1
		b = 2
		c = 0
		if (0) then
			a = 2
			b = 3
			c = a + b
	end
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

	if n, _ := cls.Call(); n.Int() != 3 {
		t.Error("exec failed:", n, cls)
	}
}

func TestArithmeticNAN(t *testing.T) {
	cls, err := LoadString(`
a = 0 
		return (1 / a + 1) * a
`)
	if err != nil {
		t.Error(err)
	}

	if v, _ := cls.Call(); !math.IsNaN(v.Float()) {
		t.Error("wrong answer")
	}
}

func TestImportLoop(t *testing.T) {
	os.MkdirAll("tmp/src", 0777)
	defer os.RemoveAll("tmp")

	ioutil.WriteFile("tmp/1.txt", []byte(`
		require "2.txt" 
		require "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/2.txt", []byte(`require "src/3.txt"`), 0777)
	ioutil.WriteFile("tmp/src/3.txt", []byte(`require "1.txt"`), 0777)
	ioutil.WriteFile("tmp/src/1.txt", []byte(`require  "../1.txt"`), 0777)

	_, err := LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "including each other") {
		t.Error("something wrong")
	}

	ioutil.WriteFile("tmp/1.txt", []byte(`require "1.txt"`), 0777)
	_, err = LoadFile("tmp/1.txt")
	if !strings.Contains(err.Error(), "including each other") {
		t.Error("something wrong")
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

func TestFalsyValue(t *testing.T) {
	assert := func(b bool) {
		if !b {
			_, fn, ln, _ := runtime.Caller(1)
			t.Fatal(fn, ln)
		}
	}

	assert(Float(0).IsFalse())
	assert(!Float(1 / math.Inf(-1)).IsFalse())
	assert(!Float(math.NaN()).IsFalse())

	s := Bool(true)
	assert(!s.IsFalse())
	s = Bool(false)
	assert(s.IsFalse())
}
