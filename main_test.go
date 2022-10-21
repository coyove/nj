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
	"sync"
	"testing"
	"time"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	_parser "github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	log.SetFlags(log.Lshortfile | log.Ltime)
	bas.AddTopValue("G", bas.Int(1))
}

type testStruct struct {
	A int
}

func (ts testStruct) Foo() int { return ts.A }

func (ts *testStruct) SetFoo(a int) { ts.A = a }

type testStructEmbed struct {
	T testStruct
}

type testStructPtr struct {
	T    *testStruct
	Next *testStructPtr
}

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path, &LoadOptions{
		Globals: bas.NewObject(0).
			SetProp("syncMap", bas.ValueOf(&sync.Map{})).
			SetProp("structAddrTest", bas.ValueOf(&testStruct{2})).
			SetProp("structAddrTest2", bas.ValueOf(testStruct{3})).
			SetProp("structAddrTestEmbed", bas.ValueOf(&testStructEmbed{testStruct{4}})).
			SetProp("nativeVarargTest", bas.ValueOf(func(a ...int) int {
				return len(a)
			})).
			SetProp("nativeVarargTest2", bas.ValueOf(func(b string, a ...int) string {
				return b + strconv.Itoa(len(a))
			})).
			SetProp("nativeVarargTest3", bas.ValueOf(func(s testStructPtr) int {
				return s.T.A
			})).
			SetProp("nativeVarargTest4", bas.ValueOf(func(s *testStructPtr, v int) {
				s.T.A = v
			})).
			SetProp("gomap", bas.ValueOf(func(m map[string]int, k string, v int) map[string]int {
				m[k] = v
				return m
			})).
			SetProp("intAlias", bas.ValueOf(func(d time.Duration) time.Time {
				return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(d)
			})).
			SetProp("goVarg", bas.ValueOf(func(a int, f func(a int, b ...int) int) int {
				return f(a, a+1, a+2)
			})).
			SetProp("boolConvert", bas.ValueOf(func(v bool) {
				if !v {
					panic("bad")
				}
			})).
			SetProp("findGlobal", bas.ValueOf(func(env *bas.Env) {
				v, err := env.MustProgram().Get("G_FLAG")
				fmt.Println(err)
				if v.IsFalse() {
					panic("findGlobal failed")
				}
				env.MustProgram().Set("G_FLAG", bas.Str("ok"))
				fmt.Println("find global")
			})).
			SetProp("G", bas.Str("test")).ToMap(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if internal.IsDebug() {
		fmt.Println(b.GoString())
	}

	_, err = b.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileTest(t *testing.T) { runFile(t, "tests/test.nj.lua") }

func TestFileStruct(t *testing.T) { runFile(t, "tests/struct.nj.lua") }

func TestFileString(t *testing.T) { runFile(t, "tests/string.nj.lua") }

func TestFileGoto(t *testing.T) { runFile(t, "tests/goto.nj.lua") }

func TestFileR2(t *testing.T) { runFile(t, "tests/r2.nj.lua") }

func TestFileStringIndex(t *testing.T) { runFile(t, "tests/indexstr.nj.lua") }

func TestFileCurry(t *testing.T) { runFile(t, "tests/curry.nj.lua") }

func TestFileEvaluator(t *testing.T) { runFile(t, "tests/eval.nj.lua") }

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
`, &LoadOptions{Globals: bas.NewObject(0).SetProp("init", bas.Int(1)).ToMap()})
		v, _ := cls.Run()
		if v := v.Object().Call(nil, bas.Int64(10)); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := v.Object().Call(nil, bas.Int64(100)); v.Int64() != 111 {
			t.Fatal(v)
		}
	}
	{
		cls, _ := LoadString(`
a = 1
function foo(x) 
if not x then return a end
for i=0,#(x) do
a=a+x[i]
end
return a
end
return foo
`, nil)
		v, _ := cls.Run()
		if v := v.Object().Call(nil, bas.Array(bas.Int64(1), bas.Int64(2), bas.Int64(3), bas.Int64(4))); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := v.Object().Call(nil, bas.Array(bas.Int64(10), bas.Int64(20))); v.Int64() != 41 {
			t.Fatal(v)
		}
	}
}

func TestTailCallPanic(t *testing.T) {
	cls, err := LoadString(`
x = 0
function foo()
x=x+1
if x == 1e5 then assert(false) end
foo()
end
foo()
`, nil)
	fmt.Println(err, cls.GoString())
	if s := cls.GoString(); !strings.Contains(s, "tailcall") {
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

func init() {
	if os.Getenv("njb") == "1" {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		fmt.Println("cpuprofile")
		for i := 0; i < 1e5; i++ {
			// _parser.Parse("(a+1)", "")
			LoadString("(a+1)", nil)
		}
		pprof.StopCPUProfile()
	}
}

func BenchmarkParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_parser.Parse("(a+1)", "")
	}
}

func BenchmarkCompiling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadString("(a+1)", nil)
	}
}

func BenchmarkGoCompiling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseExpr("(a+1)")
	}
}

func TestBigList(t *testing.T) {
	n := typ.RegMaxAddress/2 - bas.TopSymbols().Len()

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
	// fmt.Println(f.GoString())
	v2, err := f.Run()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		if v2.Native().Get(i).Int() != i {
			t.Fatal(v2)
		}
	}

	start := time.Now()
	_, err = LoadString(makeCode(typ.RegMaxAddress), nil)
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
`, &LoadOptions{Globals: bas.NewObject(0).SetProp("add", bas.ValueOf(add)).ToMap()})
	v, err := p2.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(p2.GoString())
	if v1 := v.Native().Values(); v1[0].Int64() != 101 || v1[1].Int64() != 102 || v1[2].Int64() != 103 {
		t.Fatal(v, v1, err, p2.GoString())
	}

	add = MustRun(LoadString("function foo(a) panic(a) end return function(b) foo(b) + 1 end", nil))
	_, err = add.Object().TryCall(nil, bas.Int(10))
	if err.(*bas.ExecError).GetCause() != bas.Int(10) {
		t.Fatal(err)
	}
	fmt.Println(err)
}

func TestNumberLexer(t *testing.T) {
	assert := func(src string, v bas.Value) {
		_, fn, ln, _ := runtime.Caller(1)
		r := MustRun(LoadString(src, nil))
		if r != v {
			n, _ := _parser.Parse(src, "")
			n.Dump(os.Stdout)
			t.Fatal(fn, ln, r, v)
		}
	}
	assert("1 + 2 ", bas.Int64(3))
	assert("1+ 2 ", bas.Int64(3))
	assert("-1+ 2 ", bas.Int64(1))
	assert("1- 2 ", bas.Int64(-1))
	assert("(1+1)- 2 ", bas.Zero)
	assert("1 - 2 ", bas.Int64(-1))
	assert("1 - -2 ", bas.Int64(3))
	assert("1- -2 ", bas.Int64(3))
	assert("1-2 ", bas.Int64(-1))
	assert("1.5 +2", bas.Float64(3.5))
	assert("1.5+ 2 ", bas.Float64(3.5))
	assert("12.5e-1+ 2 ", bas.Float64(3.25))
	assert("1.5e+1+ 2", bas.Float64(17))
	assert(".5+ 2 ", bas.Float64(2.5))
	assert("-.5+ 2", bas.Float64(1.5))
	assert("0x1+ 2", bas.Int64(3))
	assert("0xE+1 ", bas.Int64(15))
	assert(".5E+1 ", bas.Int64(5))
	assert("0x1_2_e+1", bas.Int64(0x12f))
	assert("([[1]])[0]", bas.Int(49))
	assert("'1'[0]", bas.Int(49))
	assert("([ [1] ])[0][0]", bas.Int(1))
	assert("([ [1]])[0][0]", bas.Int(1))
	assert("[0,1,2][1]", bas.Int(1))
	assert("[[%d]]['format'](1)", bas.Str("1"))
	assert("function()end (1)", bas.Int(1))
	assert("function()end [1][0]", bas.Int(1))
	assert("function() return 1 end()", bas.Int(1))
	assert("function() return-1 end()", bas.Int(-1))
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
		if bas.Str(v).Str() != v {
			t.Fatal(bas.Str(v).UnsafeInt64(), v)
		}
	}
}

func TestStrLess(t *testing.T) {
	a := bas.Str("a")
	b := bas.Str("a\x00")
	t.Log(a, b)
	if !a.Less(b) {
		t.FailNow()
	}
}

func TestACall(t *testing.T) {
	foo := MustRun(LoadString(`function foo(m...)
	print(m)
    assert(m[1] == 1 and m[2] == 2)
	a0 = 0 a1 = 1 a2 = 2 a3 = 3
    end
    return foo`, nil))
	foo.Object().Call(nil, bas.Nil, bas.Int64(1), bas.Int64(2))

	foo = MustRun(LoadString(`function foo(a, b, m...)
	assert(a == 1 and #(m) == 0)
    end
    return foo`, nil))
	foo.Object().Call(nil, bas.Int64(1), bas.Nil)

	foo = MustRun(LoadString(`m = {a=1}
	function m.pow2()
		return this.a * this.a
	end
	a = new(m, {a=10})
    return a`, nil))
	v := foo.Object().Get(bas.Str("pow2")).Object().Call(nil)
	if v.Int64() != 100 {
		t.Fatal(v)
	}

	foo = MustRun(LoadString(`m.a = 11
    return m.pow2()`, &LoadOptions{
		Globals: bas.NewObject(0).
			SetProp("m", bas.NewObject(0).SetPrototype(bas.NewObject(0).
				SetProp("a", bas.Int64(0)).
				AddMethod("pow2", func(e *bas.Env) {
					i := e.Object(-1).Get(bas.Str("a")).Int64()
					e.A = bas.Int64(i * i)
				})).ToValue()).
			ToMap(),
	}))
	if foo.Int64() != 121 {
		t.Fatal(foo)
	}

	foo = MustRun(LoadString(`function foo(m...)
	return sum(m.concat(m)...) + sum2(m[:2]...)
    end
    return foo`, &LoadOptions{
		Globals: bas.NewObject(0).
			SetProp("sum", bas.ValueOf(func(a ...int) int {
				s := 0
				for _, a := range a {
					s += a
				}
				return s
			})).
			SetProp("sum2", bas.ValueOf(func(a, b int) int {
				return a + b
			})).
			ToMap(),
	}))
	v = foo.Object().Call(nil, bas.Int64(1), bas.Int64(2), bas.Int64(3))
	if v.Int64() != 15 {
		t.Fatal(v)
	}
}

func TestReflectedValue(t *testing.T) {
	v := bas.Array(bas.True, bas.False)
	x := v.ToType(reflect.TypeOf([2]bool{})).Interface().([2]bool)
	if x[0] != true || x[1] != false {
		t.Fatal(x)
	}
	v = bas.NewObject(2).SetProp("a", bas.Int64(1)).SetProp("b", bas.Int64(2)).ToValue()
	y := v.ToType(reflect.TypeOf(map[string]byte{})).Interface().(map[string]byte)
	if y["a"] != 1 || y["b"] != 2 {
		t.Fatal(x)
	}

	p, _ := LoadString(`function foo(v, p)
	p[0] = 99
	return v, v + 1, nil
	end
	bar(foo)`, &LoadOptions{Globals: bas.NewObject(0).
		SetProp("bar", bas.ValueOf(func(cb func(a int, p []byte) (int, int, error)) {
			buf := []byte{0}
			a, b, _ := cb(10, buf)
			if a != 10 || b != 11 || buf[0] != 99 {
				t.Fatal(a, b)
			}
		})).ToMap(),
	})
	_, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// func TestRunTimeout(t *testing.T) {
// 	o := bas.NewObject(0)
// 	p, _ := LoadString("for i=0,1e8 do z.a = i end", &bas.Environment{
// 		Globals: bas.NewObject(0).SetProp("z", o.ToValue()),
// 	})
//
// 	p.Deadline = time.Now().Add(time.Second / 2)
// 	_, err := p.Run()
// 	if err.Error() != "timeout" {
// 		t.Fatal(err)
// 	}
// 	if v := o.Prop("a"); v.Maybe().Int(0) == 0 {
// 		t.Fatal(v)
// 	}
//
// 	p.Deadline = time.Now().Add(time.Second / 2)
// 	_, err = p.Run()
// 	if err.Error() != "timeout" {
// 		t.Fatal(err)
// 	}
// 	if v := o.Prop("a"); v.Maybe().Int64(0) == 0 {
// 		t.Fatal(v)
// 	}
// }
