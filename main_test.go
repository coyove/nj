package script

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
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

	b, err := LoadFile(path, &CompileOptions{
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
			"findGlobal": func(env *Env) {
				v, _ := env.Global.Get("G_FLAG")
				if v.IsFalse() {
					panic("findGlobal failed")
				}
				env.Global.Set("G_FLAG", String("ok"))
				env.Global.Println("find global")
			},
			"mapFunc": NativeWithParamMap("mapFunc", func(env *Env) {
				m := env.A.Map()
				if m.Get(String("a")) != Nil {
					env.A = String("a")
				}
				if m.Get(String("b")) != Nil {
					env.A = m.Get(String("b"))
				}
				if m.Get(String("c")) != Nil {
					env.A = String(m.Get(String("c")).String())
				}
				if m.Get(String("d")) != Nil {
					env.A = String(env.A.String() + m.Get(String("d")).String())
				}
			}, "DocString...", "a", "b", "c", "d"),
			"G": "test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// log.Println(b.PrettyCode())

	_, err = b.Call()
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("crab") != "" {
		fmt.Println(b.PrettyCode())
	}
}

func TestFileTest(t *testing.T) { runFile(t, "tests/test.txt") }

func TestFileStruct(t *testing.T) { runFile(t, "tests/struct.txt") }

func TestFileString(t *testing.T) { runFile(t, "tests/string.txt") }

func TestFileGoto(t *testing.T) { runFile(t, "tests/goto.txt") }

func TestFileR2(t *testing.T) { runFile(t, "tests/r2.txt") }

func TestFileStringIndex(t *testing.T) { runFile(t, "tests/indexstr.txt") }

func TestReturnFunction(t *testing.T) {
	{
		cls, _ := LoadString(`
print(init)
a = init
function foo(n) 
a=a+n
return a
end
return foo
`, &CompileOptions{GlobalKeyValues: map[string]interface{}{"init": 1}})
		v, _ := cls.Call()
		if v, _ := v.Function().Call(Int(10)); v.Int() != 11 {
			t.Fatal(v)
		}

		if v, _ := v.Function().Call(Int(100)); v.Int() != 111 {
			t.Fatal(v)
		}
	}
	{
		cls, _ := LoadString(`
a = 1
function foo(x) 
for i=0,len(x) do
a=a+x[i]
end
return a
end
return foo
`, nil)
		v, _ := cls.Call()
		if v, _ := v.Function().Call(Array(Int(1), Int(2), Int(3), Int(4))); v.Int() != 11 {
			t.Fatal(v)
		}

		if v, _ := v.Function().Call(Array(Int(10), Int(20))); v.Int() != 41 {
			t.Fatal(v)
		}

		if v, _ := v.Function().Call(); v.Int() != 41 {
			t.Fatal(v)
		}
	}
}

func TestTailCallPanic(t *testing.T) {
	cls, err := LoadString(`
x = 0
function foo()
x=x+1
if x == 1e6 then assert(false) end
foo()
end
foo()
`, nil)
	fmt.Println(err)
	if s := cls.PrettyCode(); !strings.Contains(s, "tailcall") {
		t.Fatal(s)
	}

	_, err = cls.Call()
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
`, nil)
	if err != nil {
		t.Error(err)
	}

	if v, _ := cls.Call(); v.Float() != 2.5 {
		t.Error("exec failed")
	}
}

func TestPCallStackSize(t *testing.T) {
	cls, _ := LoadString(`
a = ""
for i = 1,1e3 do
a = a + i
end
return a
`, nil)
	cls.SetDeadsize(int64(len(g)) + 10)
	res, err := cls.Call()
	if !strings.Contains(err.Error(), "deadsize") {
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
`, nil)
	if err != nil {
		t.Error(err)
	}

	// At the end of the if block, the splitInst Code will be like:
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
`, nil)
	if err != nil {
		t.Error(err)
	}

	if v, _ := cls.Call(); !math.IsNaN(v.Float()) {
		t.Error("wrong answer")
	}
}

func BenchmarkCompiling(b *testing.B) {
	buf, _ := ioutil.ReadFile("tests/string.txt")
	y := string(bytes.Repeat(buf, 100))
	for i := 0; i < b.N; i++ {
		p, err := LoadString(string(y), nil)
		if err != nil {
			b.Fatal(err)
		}
		p.Stdout = ioutil.Discard
		// p.Run()
	}
}

func TestBigList(t *testing.T) {
	n := 2000 - len(g)

	makeCode := func(n int) string {
		buf := bytes.Buffer{}
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d = %d\n", i, i))
		}
		buf.WriteString("return {")
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d,", i))
		}
		buf.Truncate(buf.Len() - 1)
		return buf.String() + "}"
	}

	f, _ := LoadString(makeCode(n), nil)
	v2, err := f.Call()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		if v2.Map().Get(Int(int64(i))).Int() != int64(i) {
			t.Fatal(v2)
		}
	}

	_, err = LoadString(makeCode(2000), nil)
	if !strings.Contains(err.Error(), "too many") {
		t.Fatal(err)
	}

	{
		buf := bytes.NewBufferString("function foo(")
		for i := 0; i < 256; i++ {
			buf.WriteString(fmt.Sprintf("a%d,", i))
		}
		buf.WriteString("x) end")
		_, err = LoadString(buf.String(), nil)
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

func TestPlainReturn(t *testing.T) {
	if _, err := LoadString("return", nil); err != nil {
		t.FailNow()
	}
	if _, err := LoadString("return ", nil); err != nil {
		t.FailNow()
	}
	if _, err := LoadString("return \n ", nil); err != nil {
		t.FailNow()
	}
}

func TestFunctionClosure(t *testing.T) {
	p, _ := LoadString(` local a = 0
function add ()
a=a+1
return a
end
return add`, nil)
	add, _ := p.Run()

	p2, _ := LoadString(`
local a = 100
return {a + add(), a + add(), a + add()}
`, &CompileOptions{GlobalKeyValues: map[string]interface{}{"add": add}})
	v, err := p2.Run()
	if v1 := v.Map().Array(); v1[0].Int() != 101 || v1[1].Int() != 102 || v1[2].Int() != 103 {
		t.Fatal(v, v1, err, p2.PrettyCode())
	}

}

func TestNumberLexer(t *testing.T) {
	assert := func(src string, v Value) {
		_, fn, ln, _ := runtime.Caller(1)
		r := MustRun(LoadString("return "+src, nil))
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
	assert("0x1_2_e+1", Int(0x12f))
}

func TestSmallString(t *testing.T) {
	rand.Seed(time.Now().Unix())
	randString := func() string {
		buf := make([]byte, rand.Intn(10))
		for i := range buf {
			buf[i] = byte(rand.Intn(256))
		}
		return string(buf)
	}
	for i := 0; i < 1e6; i++ {
		v := randString()
		if String(v).rawStr() != v {
			t.Fatal(String(v).v, v)
		}
	}
}

func TestRHMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := Map{}
	m2 := map[int64]int64{}
	counter := int64(0)
	for i := 0; i < 1e6; i++ {
		x := rand.Int63()
		if x%2 == 0 {
			x = counter
			counter++
		}
		_, memSpace := m.Set(Int(int64(x)), Int(int64(x)))
		_ = memSpace
		m2[x] = x
	}
	for k := range m2 {
		delete(m2, k)
		m.Set(Int(k), Nil)
		if rand.Intn(10000) == 0 {
			break
		}
	}

	fmt.Println(m.count, len(m.hashItems), len(m2))

	for k, v := range m2 {
		if m.Get(Int(k)).Int() != v {
			for _, e := range m.hashItems {
				if e.Key.Int() == k {
					t.Log(e)
				}
			}
			t.Fatal(m.Get(Int(k)), k, v)
		}
	}

	if m.Len() != len(m2) {
		t.Fatal(m.Len(), len(m2))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if _, ok := m2[k.Int()]; !ok {
			t.Fatal(k, v, len(m2))
		}
		delete(m2, k.Int())
	}
	if len(m2) != 0 {
		t.Fatal(len(m2))
	}

	m.Clear()
	m.Set(Int(0), Int(0))
	m.Set(Int(1), Int(1))
	m.Set(Int(2), Int(2))

	for i := 4; i < 9; i++ {
		m.Set(Int(int64(i*i)), Int(0))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		fmt.Println(k, v)
	}
}
