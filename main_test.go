package nj

import (
	"bytes"
	"encoding/json"
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

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	_parser "github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	log.SetFlags(log.Lshortfile | log.Ltime)
	bas.Globals.SetProp("G", bas.Int(1))
	BuildGlobalStack()
}

type testStruct struct {
	A int
}

func (ts testStruct) Foo() int { return ts.A }

func (ts *testStruct) SetFoo(a int) { ts.A = a }

type testStructEmbed struct {
	T testStruct
}

func runFile(t *testing.T, path string) {
	if !flag.Parsed() {
		flag.Parse()
	}

	b, err := LoadFile(path, &bas.Environment{
		Globals: bas.NewObject(0).
			SetProp("structAddrTest", bas.ValueOf(&testStruct{2})).
			SetProp("structAddrTest2", bas.ValueOf(testStruct{3})).
			SetProp("structAddrTestEmbed", bas.ValueOf(&testStructEmbed{testStruct{4}})).
			SetProp("nativeVarargTest", bas.ValueOf(func(a ...int) int {
				return len(a)
			})).
			SetProp("nativeVarargTest2", bas.ValueOf(func(b string, a ...int) string {
				return b + strconv.Itoa(len(a))
			})).
			SetProp("gomap", bas.ValueOf(func(m map[string]int, k string, v int) map[string]int {
				m[k] = v
				return m
			})).
			SetProp("intAlias", bas.ValueOf(func(d time.Duration) time.Time {
				return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(d)
			})).
			SetProp("boolConvert", bas.ValueOf(func(v bool) {
				if !v {
					panic("bad")
				}
			})).
			SetProp("findGlobal", bas.ValueOf(func(env *bas.Env) {
				v, err := env.Global.Get("G_FLAG")
				fmt.Println(err)
				if v.IsFalse() {
					panic("findGlobal failed")
				}
				env.Global.Set("G_FLAG", bas.Str("ok"))
				fmt.Println("find global")
			})).
			SetProp("G", bas.Str("test")),
	})
	if err != nil {
		t.Fatal(err)
	}

	if internal.IsDebug() {
		fmt.Println(b.GoString())
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
`, &bas.Environment{Globals: bas.NewObject(0).SetProp("init", bas.Int(1))})
		v, _ := cls.Run()
		if v := bas.Call(v.Object(), bas.Int64(10)); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := bas.Call(v.Object(), bas.Int64(100)); v.Int64() != 111 {
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
		if v := bas.Call(v.Object(), bas.NewArray(bas.Int64(1), bas.Int64(2), bas.Int64(3), bas.Int64(4)).ToValue()); v.Int64() != 11 {
			t.Fatal(v)
		}

		if v := bas.Call(v.Object(), bas.NewArray(bas.Int64(10), bas.Int64(20)).ToValue()); v.Int64() != 41 {
			t.Fatal(v)
		}

		if v := bas.Call(v.Object()); v.Int64() != 41 {
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

func BenchmarkCompiling(b *testing.B) {
	BuildGlobalStack()
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
	n := typ.RegMaxAddress/2 - bas.Globals.Len()

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

func TestFalsyValue(t *testing.T) {
	assert := func(b bool) {
		if !b {
			_, fn, ln, _ := runtime.Caller(1)
			t.Fatal(fn, ln)
		}
	}

	assert(bas.Float64(0).IsFalse())
	assert(bas.Float64(1 / math.Inf(-1)).IsFalse())
	assert(!bas.Float64(math.NaN()).IsFalse())
	assert(!bas.Bool(true).IsFalse())
	assert(bas.Bool(false).IsFalse())
	assert(bas.Str("").IsFalse())
	assert(bas.Bytes(nil).IsTrue())
	assert(bas.Bytes([]byte("")).IsTrue())
	assert(!bas.ValueOf([]byte("")).IsFalse())
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
`, &bas.Environment{Globals: bas.NewObject(0).SetProp("add", bas.ValueOf(add))})
	v, err := p2.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(p2.GoString())
	if v1 := v.Array().Values(); v1[0].Int64() != 101 || v1[1].Int64() != 102 || v1[2].Int64() != 103 {
		t.Fatal(v, v1, err, p2.GoString())
	}

}

func TestNumberLexer(t *testing.T) {
	assert := func(src string, v bas.Value) {
		_, fn, ln, _ := runtime.Caller(1)
		r := MustRun(LoadString(src, nil))
		if r != v {
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
	assert("lambda()end (1)", bas.Int(1))
	assert("lambda()end [1][0]", bas.Int(1))
	assert("lambda()1 end()", bas.Int(1))
	assert("lambda() -1 end()", bas.Int(-1))
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
	if !bas.Less(a, b) {
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
	bas.Call(foo.Object(), bas.Nil, bas.Int64(1), bas.Int64(2))

	foo = MustRun(LoadString(`function foo(a, b, m...)
	assert(a == 1 and #(m) == 0)
    end
    return foo`, nil))
	bas.Call(foo.Object(), bas.Int64(1))

	foo = MustRun(LoadString(`m = {a=1}
	function m.pow2()
		return this.a * this.a
	end
	a = new(m, {a=10})
    return a`, nil))
	v := bas.Call(foo.Object().Prop("pow2").Object())
	if v.Int64() != 100 {
		t.Fatal(v)
	}

	foo = MustRun(LoadString(`m.a = 11
    return m.pow2()`, &bas.Environment{
		Globals: bas.NewObject(0).
			SetProp("m", bas.NewObject(0).SetPrototype(bas.NewObject(0).
				SetProp("a", bas.Int64(0)).
				SetMethod("pow2", func(e *bas.Env) {
					i := e.Object(-1).Prop("a").Int64()
					e.A = bas.Int64(i * i)
				}, "")).ToValue()),
	}))
	if foo.Int64() != 121 {
		t.Fatal(foo)
	}

	foo = MustRun(LoadString(`function foo(m...)
	return sum(m.concat(m)...) + sum2(m.slice(0, 2)...)
    end
    return foo`, &bas.Environment{
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
			})),
	}))
	v = bas.Call(foo.Object(), bas.Int64(1), bas.Int64(2), bas.Int64(3))
	if v.Int64() != 15 {
		t.Fatal(v)
	}
}

func TestReflectedValue(t *testing.T) {
	v := bas.NewArray(bas.True, bas.False).ToValue()
	x := v.ReflectValue(reflect.TypeOf([2]bool{})).Interface().([2]bool)
	if x[0] != true || x[1] != false {
		t.Fatal(x)
	}
	v = bas.NewObject(2).SetProp("a", bas.Int64(1)).SetProp("b", bas.Int64(2)).ToValue()
	y := v.ReflectValue(reflect.TypeOf(map[string]byte{})).Interface().(map[string]byte)
	if y["a"] != 1 || y["b"] != 2 {
		t.Fatal(x)
	}

	p, _ := LoadString(`function foo(v, p)
	p[0] = 99
	return v, v + 1, nil
	end
	bar(foo)`, &bas.Environment{Globals: bas.NewObject(0).
		SetProp("bar", bas.ValueOf(func(cb func(a int, p []byte) (int, int, error)) {
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

var jsonTest = func() string {
	x := `{"data": {"feeds": [{"audios": [], "comment_num": 0, "content": "\n\t\rabc", "create_time": {"time": 1640074878, "time_desc": "20 mins ago"}, "dislike_num": 0, "id": "61c18e7e4a71d05c90f19951", "in_hq": false, "like_num": 0, "other_info": {"spotify_info": {}}, "pics": [], "tags": ["999"], "user_id": "60e3fa8af9bf8d3b66c470b1", "user_info": {"age": 31, "avatar": "99eb05ca-59d9-11e9-8672-00163e02deb4", "bio": "\u0e40\u0e18\u0e2d\u0e04\u0e37\u0e2d\u0e1a\u0e38\u0e04\u0e04\u0e25\u0e25\u0e36\u0e01\u0e25\u0e31\u0e1a", "country": "TH", "cover_photo": "60e3fa8af9bf8d3b66c470b1_coverPhoto_8_1640071635228.png", "create_time": 1625553546.0, "frame_fileid": "", "gender": "girl", "huanxin_id": "love131004954661603", "is_vip": false, "lit_id": 1119075350, "nickname": "\u7761", "party_level_info": {"received": {"avatar": "52b743c2-5705-11ec-80d6-00163e02a5e5", "diamonds": 0, "level": 18, "new_diamonds": 72255}, "sent": {"avatar": "99c83d2a-5d6e-11eb-8b66-00163e022423", "diamonds": 0, "level": 3, "new_diamonds": 50963}}, "party_top_three": -1, "removed": false, "role": 0, "tag_str": "_!EMPTY", "user_id": "60e3fa8af9bf8d3b66c470b1"}, "video": null, "video_length": 0, "visibility": 0}, {"audios": [], "comment_num": 0, "content": "\u7890", "create_time": {"time": 1640074781, "time_desc": "22 mins ago"}, "dislike_num": 0, "id": "61c18e1d4a71d05c992f578c", "in_hq": false, "like_num": 0, "other_info": {"spotify_info": {}}, "pics": [], "tags": ["999"], "user_id": "60e3fa8af9bf8d3b66c470b1", "user_info": {"age": 31, "avatar": "99eb05ca-59d9-11e9-8672-00163e02deb4", "bio": "\u0e40\u0e18\u0e2d\u0e04\u0e37\u0e2d\u0e1a\u0e38\u0e04\u0e04\u0e25\u0e25\u0e36\u0e01\u0e25\u0e31\u0e1a", "country": "TH", "cover_photo": "60e3fa8af9bf8d3b66c470b1_coverPhoto_8_1640071635228.png", "create_time": 1625553546.0, "frame_fileid": "", "gender": "girl", "huanxin_id": "love131004954661603", "is_vip": false, "lit_id": 1119075350, "nickname": "\u7761", "party_level_info": {"received": {"avatar": "52b743c2-5705-11ec-80d6-00163e02a5e5", "diamonds": 0, "level": 18, "new_diamonds": 72255}, "sent": {"avatar": "99c83d2a-5d6e-11eb-8b66-00163e022423", "diamonds": 0, "level": 3, "new_diamonds": 50963}}, "party_top_three": -1, "removed": false, "role": 0, "tag_str": "_!EMPTY", "user_id": "60e3fa8af9bf8d3b66c470b1"}, "video": null, "video_length": 0, "visibility": 0}], "has_next": false, "next_start": -1}, "result": 0, "success": true}`
	return "[" + strings.Join([]string{x, x, x, x}, ",") + "]"
}()

func TestParseJSON(t *testing.T) {
	v1, err := _parser.ParseJSON(jsonTest)
	internal.PanicErr(err)
	m := []interface{}{}
	json.Unmarshal([]byte(jsonTest), &m)
	v2 := parseJSON(m)
	if !bas.DeepEqual(v1, v2) {
		t.Fatal(v1, v2)
	}
}

func BenchmarkNativeJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_parser.ParseJSON(jsonTest)
	}
}

func BenchmarkGoJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseStrictJSON(jsonTest)
	}
}

// func BenchmarkGJSON(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		bas.ValueOf(gjson.Parse(jsonTest))
// 	}
// }

func TestRunTimeout(t *testing.T) {
	o := bas.NewObject(0)
	p, _ := LoadString("for i=0,1e8 do z.a = i end", &bas.Environment{
		Globals: bas.NewObject(0).SetProp("z", o.ToValue()),
	})

	p.Deadline = time.Now().Add(time.Second / 2)
	_, err := p.Run()
	if err.Error() != "timeout" {
		t.Fatal(err)
	}
	if v := o.Prop("a"); v.Safe().Int(0) == 0 {
		t.Fatal(v)
	}

	p.Deadline = time.Now().Add(time.Second / 2)
	_, err = p.Run()
	if err.Error() != "timeout" {
		t.Fatal(err)
	}
	if v := o.Prop("a"); v.Safe().Int(0) == 0 {
		t.Fatal(v)
	}
}
