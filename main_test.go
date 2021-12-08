package nj

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coyove/nj/internal"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	log.SetFlags(log.Lshortfile | log.Ltime)
	Globals.SetProp("G", Int(1))
}

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path, &CompileOptions{
		Globals: NewObject(0).
			SetProp("nativeVarargTest", ValueOf(func(a ...int) int {
				return len(a)
			})).
			SetProp("nativeVarargTest2", ValueOf(func(b string, a ...int) string {
				return b + strconv.Itoa(len(a))
			})).
			SetProp("intAlias", ValueOf(func(d time.Duration) time.Time {
				return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(d)
			})).
			SetProp("boolConvert", ValueOf(func(v bool) {
				if !v {
					panic("bad")
				}
			})).
			SetProp("findGlobal", ValueOf(func(env *Env) {
				v, err := env.Global.Get("G_FLAG")
				fmt.Println(err)
				if v.IsFalse() {
					panic("findGlobal failed")
				}
				env.Global.Set("G_FLAG", Str("ok"))
				fmt.Println("find global")
			})).
			SetProp("G", Str("test")),
	})
	if err != nil {
		t.Fatal(err)
	}

	if internal.IsDebug() {
		fmt.Println(b.PrettyCode())
	}
	// log.Println(b.Symbols)

	_, err = b.Run()
	if err != nil {
		t.Fatal(err)
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
`, &CompileOptions{Globals: NewObject(0).SetProp("init", Int(1))})
		v, _ := cls.Run()
		if v := Call(v.Object(), Int64(10)); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := Call(v.Object(), Int64(100)); v.Int64() != 111 {
			t.Fatal(v)
		}
	}
	{
		cls, _ := LoadString(`
a = 1
function foo(x) 
if not x then return a end
for i=0,len(x) do
a=a+x[i]
end
return a
end
return foo
`, nil)
		v, _ := cls.Run()
		if v := Call(v.Object(), NewArray(Int64(1), Int64(2), Int64(3), Int64(4)).ToValue()); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := Call(v.Object(), NewArray(Int64(10), Int64(20)).ToValue()); v.Int64() != 41 {
			t.Fatal(v)
		}

		if v := Call(v.Object()); v.Int64() != 41 {
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

	_, err = cls.Run()
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

	if v, _ := cls.Run(); v.Float64() != 2.5 {
		t.Error("exec failed")
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

	if n, _ := cls.Run(); n.Int64() != 3 {
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

	if v, _ := cls.Run(); !math.IsNaN(v.Float64()) {
		t.Error("wrong answer")
	}
}

func BenchmarkCompiling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// _parser.Parse("(a+1)", "")
		LoadString("(a+1)", nil)
	}
}

func BenchmarkGoCompiling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseExpr("(a+1)")
	}
}

func BenchmarkRHMap10(b *testing.B)    { benchmarkRHMap(b, 10) }
func BenchmarkGoMap10(b *testing.B)    { benchmarkGoMap(b, 10) }
func BenchmarkRHMap20(b *testing.B)    { benchmarkRHMap(b, 20) }
func BenchmarkGoMap20(b *testing.B)    { benchmarkGoMap(b, 20) }
func BenchmarkRHMap50(b *testing.B)    { benchmarkRHMap(b, 50) }
func BenchmarkGoMap50(b *testing.B)    { benchmarkGoMap(b, 50) }
func BenchmarkRHMapUnc10(b *testing.B) { benchmarkRHMapUnconstrainted(b, 10) }
func BenchmarkGoMapUnc10(b *testing.B) { benchmarkGoMapUnconstrainted(b, 10) }

func benchmarkRHMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := NewObject(n)
	for i := 0; i < n; i++ {
		m.Set(Int64(int64(i)), Int64(int64(i)))
	}
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(n)
		if m.Get(Int64(int64(idx))) != Int64(int64(idx)) {
			b.Fatal(idx, m)
		}
	}
}

func benchmarkRHMapUnconstrainted(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := NewObject(1)
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			x := rand.Intn(n)
			m.Set(Int64(int64(x)), Int64(int64(i)))
		}
	}
}

func benchmarkGoMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := map[int]int{}
	for i := 0; i < n; i++ {
		m[i] = i
	}
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(n)
		if m[idx] == -1 {
			b.Fatal(idx, m)
		}
	}
}

func benchmarkGoMapUnconstrainted(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			idx := rand.Intn(n)
			m[idx] = i
		}
	}
}

func TestBigList(t *testing.T) {
	n := maxAddress/2 - Globals.Len()

	makeCode := func(n int) string {
		buf := bytes.Buffer{}
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d = %d\n", i, i))
		}
		buf.WriteString("return [")
		for i := 0; i < n; i++ {
			buf.WriteString(fmt.Sprintf("a%d,", i))
		}
		buf.Truncate(buf.Len() - 1)
		return buf.String() + "]"
	}

	f, _ := LoadString(makeCode(n), nil)
	v2, err := f.Run()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		if v2.Array().Get(i).Int() != i {
			t.Fatal(v2)
		}
	}

	{
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	_, err = LoadString(makeCode(maxAddress), nil)
	fmt.Println("load", time.Since(start))
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

	assert(Float64(0).IsFalse())
	assert(Float64(1 / math.Inf(-1)).IsFalse())
	assert(!Float64(math.NaN()).IsFalse())
	assert(!Bool(true).IsFalse())
	assert(Bool(false).IsFalse())
	assert(Str("").IsFalse())
	assert(Bytes(nil).IsTrue())
	assert(Bytes([]byte("")).IsTrue())
	assert(!ValueOf([]byte("")).IsFalse())
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
return [a + add(), a + add(), a + add()]
`, &CompileOptions{Globals: NewObject(0).SetProp("add", ValueOf(add))})
	v, err := p2.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(p2.PrettyCode())
	if v1 := v.Array().Values(); v1[0].Int64() != 101 || v1[1].Int64() != 102 || v1[2].Int64() != 103 {
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
	assert("1 + 2 ", Int64(3))
	assert("1+ 2 ", Int64(3))
	assert("-1+ 2 ", Int64(1))
	assert("1- 2 ", Int64(-1))
	assert("1 - 2 ", Int64(-1))
	assert("1.5 +2", Float64(3.5))
	assert("1.5+ 2 ", Float64(3.5))
	assert("12.5e-1+ 2 ", Float64(3.25))
	assert("1.5e+1+ 2", Float64(17))
	assert(".5+ 2 ", Float64(2.5))
	assert("-.5+ 2", Float64(1.5))
	assert("0x1+ 2", Int64(3))
	assert("0xE+1 ", Int64(15))
	assert(".5E+1 ", Int64(5))
	assert("0x1_2_e+1", Int64(0x12f))
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
		if Str(v).Str() != v {
			t.Fatal(Str(v).v, v)
		}
	}
}

func TestRHMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := Object{}
	m2 := map[int64]int64{}
	counter := int64(0)
	for i := 0; i < 1e6; i++ {
		x := rand.Int63()
		if x%2 == 0 {
			x = counter
			counter++
		}
		m.Set(Int64(x), Int64(x))
		m2[x] = x
	}
	for k := range m2 {
		delete(m2, k)
		m.Delete(Int64(k))
		if rand.Intn(10000) == 0 {
			break
		}
	}

	fmt.Println(m.count, len(m.items), len(m2))

	for k, v := range m2 {
		if m.Get(Int64(k)).Int64() != v {
			for _, e := range m.items {
				if e.Key.Int64() == k {
					t.Log(e)
				}
			}
			t.Fatal(m.Get(Int64(k)), k, v)
		}
	}

	if m.Len() != len(m2) {
		t.Fatal(m.Len(), len(m2))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if _, ok := m2[k.Int64()]; !ok {
			t.Fatal(k, v, len(m2))
		}
		delete(m2, k.Int64())
	}
	if len(m2) != 0 {
		t.Fatal(len(m2))
	}

	m.Clear()
	m.Set(Int64(0), Int64(0))
	m.Set(Int64(1), Int64(1))
	m.Set(Int64(2), Int64(2))

	for i := 4; i < 9; i++ {
		m.Set(Int64(int64(i*i)), Int64(0))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		fmt.Println(k, v)
	}
}

func TestACall(t *testing.T) {
	foo := MustRun(LoadString(`function foo(m...)
	print(m)
    assert(m[1] == 1 and m[2] == 2)
	a0 = 0 a1 = 1 a2 = 2 a3 = 3
    end
    return foo`, nil))
	Call(foo.Object(), Nil, Int64(1), Int64(2))

	foo = MustRun(LoadString(`function foo(a, b, m...)
	assert(a == 1 and len(m) == 0)
    end
    return foo`, nil))
	Call(foo.Object(), Int64(1))

	foo = MustRun(LoadString(`m = {a=1}
	function m.pow2()
		return this.a * this.a
	end
	a = new(m, {a=10})
    return a`, nil))
	v := Call(foo.Object().Prop("pow2").Object())
	if v.Int64() != 100 {
		t.Fatal(v)
	}

	foo = MustRun(LoadString(`m.a = 11
    return m.pow2()`, &CompileOptions{
		Globals: NewObject(0).
			SetProp("m", NewObject(0).SetPrototype(NewObject(0).
				SetProp("a", Int64(0)).
				SetMethod("pow2", func(e *Env) {
					i := e.Object(-1).Prop("a").Int64()
					e.A = Int64(i * i)
				}, "")).ToValue()),
	}))
	if foo.Int64() != 121 {
		t.Fatal(foo)
	}

	foo = MustRun(LoadString(`function foo(m...)
	return sum(m.concat(m)...) + sum2(m.slice(0, 2)...)
    end
    return foo`, &CompileOptions{
		Globals: NewObject(0).
			SetProp("sum", ValueOf(func(a ...int) int {
				s := 0
				for _, a := range a {
					s += a
				}
				return s
			})).
			SetProp("sum2", ValueOf(func(a, b int) int {
				return a + b
			})),
	}))
	v = Call(foo.Object(), Int64(1), Int64(2), Int64(3))
	if v.Int64() != 15 {
		t.Fatal(v)
	}
}

func TestReflectedValue(t *testing.T) {
	v := NewArray(True, False).ToValue()
	x := v.ReflectValue(reflect.TypeOf([2]bool{})).Interface().([2]bool)
	if x[0] != true || x[1] != false {
		t.Fatal(x)
	}
	v = NewObject(2).SetProp("a", Int64(1)).SetProp("b", Int64(2)).ToValue()
	y := v.ReflectValue(reflect.TypeOf(map[string]byte{})).Interface().(map[string]byte)
	if y["a"] != 1 || y["b"] != 2 {
		t.Fatal(x)
	}

	p, _ := LoadString(`function foo(v, p)
	p[0] = 99
	return v, v + 1, nil
	end
	bar(foo)`, &CompileOptions{Globals: NewObject(0).
		SetProp("bar", ValueOf(func(cb func(a int, p []byte) (int, int, error)) {
			buf := []byte{0}
			a, b, _ := cb(10, buf)
			if a != 10 || b != 11 || buf[0] != 99 {
				t.Fatal(a, b)
			}
		})),
	})
	_, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashcodeDist(t *testing.T) {
	z := map[uint64]int{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1e6; i++ {
		v := Int64(int64(i)).HashCode()
		z[v]++
	}
	fmt.Println(len(z))

	z = map[uint64]int{}
	for i := 0; i < 1e6; i++ {
		v := Int64(rand.Int63()).HashCode()
		z[v]++
	}
	fmt.Println(len(z))
}
