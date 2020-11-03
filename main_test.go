package script

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	log.SetFlags(log.Lshortfile | log.Ltime)
	AddGlobalValue("G", 1)
}

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path, CompileOptions{
		GlobalKeyValues: map[string]interface{}{
			"nativeVarargTest": func(a ...int) int {
				return len(a)
			},
			"nativeVarargTest2": func(b string, a ...int) string {
				return b + strconv.Itoa(len(a))
			},
			"intAlias": func(d time.Duration) time.Time {
				return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(d)
			},
			"boolConvert": func(v bool) {
				if !v {
					panic("bad")
				}
			},
			"mapFunc": NativeWithParamMap("mapFunc", func(env *Env, in Arguments) {
				if !in["a"].IsNil() {
					env.A = _str("a")
				}
				if !in["b"].IsNil() {
					env.A = in["b"]
				}
				if !in["c"].IsNil() {
					env.A = _str(in["c"].String())
				}
				if !in["d"].IsNil() {
					env.A = _str(env.A.String() + in["d"].String())
				}
			}, "doc...", "a", "b", "c", "d"),
			"G": "test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// log.Println(b.PrettyCode())

	i, i2, err := b.Call()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(i, i2, err, "str alloc:", b.Survey.TotalStringAlloc)
}

func TestFileTest(t *testing.T) { runFile(t, "tests/test.txt") }

func TestFileString(t *testing.T) { runFile(t, "tests/string.txt") }

func TestFileGoto(t *testing.T) { runFile(t, "tests/goto.txt") }

func TestFileR2(t *testing.T) { runFile(t, "tests/r2.txt") }

func TestFileStringIndex(t *testing.T) { runFile(t, "tests/indexstr.txt") }

func TestFileVararg(t *testing.T) { runFile(t, "tests/vararg.txt") }

func TestReturnFunction(t *testing.T) {
	{
		cls, _ := LoadString(`
print(init)
a = init
function foo(n) 
a+=n
return a
end
return foo
`, CompileOptions{GlobalKeyValues: map[string]interface{}{"init": 1}})
		v, _, _ := cls.Call()
		if v, _, _ := v.Function().Call(Int(10)); v.Int() != 11 {
			t.Fatal(v)
		}

		if v, _, _ := v.Function().Call(Int(100)); v.Int() != 111 {
			t.Fatal(v)
		}
	}
	{
		cls, _ := LoadString(`
a = 1
function foo(...x) 
for i=1,#x do
a+=x[i]
end
return a
end
return foo
`)
		v, _, _ := cls.Call()
		if v, _, _ := v.Function().Call(Int(1), Int(2), Int(3), Int(4)); v.Int() != 11 {
			t.Fatal(v)
		}

		if v, _, _ := v.Function().Call(Int(10), Int(20)); v.Int() != 41 {
			t.Fatal(v)
		}

		if v, _, _ := v.Function().Call(); v.Int() != 41 {
			t.Fatal(v)
		}
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
	if s := cls.PrettyCode(); !strings.Contains(s, "tail-call") {
		t.Fatal(s)
	}

	_, _, err := cls.Call()
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

	if v, _, _ := cls.Call(); v.Float() != 2.5 {
		t.Error("exec failed")
	}
}

func TestPCallStackSize(t *testing.T) {
	cls, _ := LoadString(`
function test() 
local a, b, c = 1, 2, 3
assert(a, b, c)
return a
end
_, err = pcall(test)
print(err)
assert(match(err.Error(), "overflow" ))
`)
	cls.MaxStackSize = 7 + int64(len(g))
	cls.Call()

	cls, _ = LoadString(`
a = ""
for i = 1,1e3 do
a = a .. i
end
return a
`)
	cls.MaxStringSize = (int64(len(g)) + 10) * 16 // 10: a small value
	res, _, err := cls.Call()
	if !strings.Contains(err.Error(), "string overflow") {
		t.Fatal(res, err)
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

	// At the end of the if block, the splitInst Code will be like:
	// R0 = a, R1 = b -> Add
	// But after the if block, there is another c = a + b, we can't re-use the registers R0 and R1
	// because they will not contain the value we want as the if block was not executed at all.

	if n, _, _ := cls.Call(); n.Int() != 3 {
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

	if v, _, _ := cls.Call(); !math.IsNaN(v.Float()) {
		t.Error("wrong answer")
	}
}

func BenchmarkCompiling(b *testing.B) {
	buf, _ := ioutil.ReadFile("tests/string.txt")
	y := string(bytes.Repeat(buf, 100))
	for i := 0; i < b.N; i++ {
		p, err := LoadString(string(y))
		if err != nil {
			b.Fatal(err)
		}
		p.Stdout = ioutil.Discard
		// p.Run()
	}
}

func TestTooManyVariables(t *testing.T) {
	n := 2000 - len(g)

	makeCode := func(n int) string {
		buf := bytes.Buffer{}
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d = %d\n", i, i))
		}
		buf.WriteString("return ")
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d,", i))
		}
		buf.Truncate(buf.Len() - 1)
		return buf.String()
	}

	f, _ := LoadString(makeCode(n))
	_, v2, err := f.Call()
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i < n; i++ {
		if v2[i-1].Int() != int64(i) {
			t.Fatal(v2)
		}
	}

	_, err = LoadString(makeCode(2000))
	if !strings.Contains(err.Error(), "too many") {
		t.Fatal(err)
	}

	{
		buf := bytes.NewBufferString("function foo(")
		for i := 0; i < 256; i++ {
			buf.WriteString(fmt.Sprintf("a%d,", i))
		}
		buf.WriteString("x) end")
		_, err = LoadString(buf.String())
		if !strings.Contains(err.Error(), "too many") {
			t.Fatal(err)
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
	assert(Float(1 / math.Inf(-1)).IsFalse())
	assert(!Float(math.NaN()).IsFalse())

	s := Bool(true)
	assert(!s.IsFalse())
	s = Bool(false)
	assert(s.IsFalse())
}

func TestFunctionClosure(t *testing.T) {
	p, _ := LoadString(` local a = 0
function add ()
a+=1
return a
end
return add`)
	add, _, _ := p.Run()

	p2, _ := LoadString(`
local a = 100
return a + add(), a + add(), a + add()
`, CompileOptions{GlobalKey: "add", GlobalValue: add})
	v, v1, err := p2.Run()
	if v.Int() != 101 || v1[0].Int() != 102 || v1[1].Int() != 103 {
		t.Fatal(v, v1, err, p2.PrettyCode())
	}

}

func TestNumberLexer(t *testing.T) {
	assert := func(src string, v Value) {
		_, fn, ln, _ := runtime.Caller(1)
		r, _ := MustRun(LoadString("return " + src))
		if r != v {
			t.Fatal(fn, ln, r, v)
		}
	}
	assert("1 + 2 ", Int(3))
	assert("1+ 2 ", Int(3))
	assert("-1+ 2 ", Int(1))
	assert("1- 2 ", Int(-1))
	assert("1 - 2 ", Int(-1))
	assert("1.5 +2", Float(3.5))
	assert("1.5+ 2 ", Float(3.5))
	assert("12.5e-1+ 2 ", Float(3.25))
	assert("1.5e+1+ 2", Float(17))
	assert(".5+ 2 ", Float(2.5))
	assert("-.5+ 2", Float(1.5))
	assert("0x1+ 2", Int(3))
	assert("0xE+1 ", Int(15))
	assert(".5E+1 ", Int(5))
}
