package bas

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

func randString() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func randInt(len, idx int) int {
	buf := make([]byte, 6)
	rand.Read(buf)
	v := rand.Int()
	if Int(v).HashCode()%uint32(len) == uint32(idx) {
		return v
	}
	return randInt(len, idx)
}

func TestMapForeachDelete(t *testing.T) {
	rand.Seed(time.Now().Unix())
	check := func(o *Map, idx int, k int, dist int32) {
		i := o.items[idx]
		if i.key.Int() != k || i.dist != dist {
			t.Fatal(o.items, string(debug.Stack()))
		}
	}

	o := newMap(1)
	o.noresize = true
	a := randInt(2, 1)
	b := randInt(2, 1)
	c := randInt(2, 1)
	o.Set(Int(a), Int(a)) // [null, a+0]
	o.Set(Int(b), Int(b)) // [b+1, a+0]
	o.Delete(Int(a))      // [b+1, deleted+0]
	if o.items[0].dist != 1 || !o.items[1].pDeleted {
		t.Fatal(o.items)
	}
	o.Set(Int(a), Int(a)) // [a+1, b+0]
	check(o, 0, a, 1)
	check(o, 1, b, 0)

	o.Delete(Int(a))      // [deleted+1, b+0]
	o.Set(Int(c), Int(c)) // [c+1, b+0]
	check(o, 0, c, 1)
	check(o, 1, b, 0)

	o = newMap(2)
	o.noresize = true
	a = randInt(4, 1)
	b = randInt(4, 1)
	c = randInt(4, 1)
	d := randInt(4, 1)
	o.Set(Int(a), Int(a))
	o.Set(Int(b), Int(b))
	o.Set(Int(c), Int(c))
	o.Set(Int(d), Int(d)) // [d+3, a+0, b+1, c+2]
	check(o, 0, d, 3)
	check(o, 1, a, 0)
	check(o, 2, b, 1)
	check(o, 3, c, 2)

	o.Delete(Int(b))      // [d+3, a+0, deleted+1, c+2]
	o.Set(Int(b), Int(b)) // [b+3, a+0, c+1, d+2]
	check(o, 0, b, 3)
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 2)

	o.Delete(Int(a))
	o.Delete(Int(c)) // [b+3, deleted+0, deleted+1, d+2]
	loopCount := 0
	o.Foreach(func(k Value, v *Value) bool { loopCount++; return true })
	if loopCount != 2 {
		t.Fatal(loopCount, o.items)
	}
	for k, _ := o.FindNext(Nil); k != Nil; k, _ = o.FindNext(k) {
		loopCount--
	}
	if loopCount != 0 {
		t.Fatal(o.items)
	}
	a = randInt(4, 2)
	o.Set(Int(a), Int(a)) // [a+2, deleted+0, d+1, b+2]
	check(o, 0, a, 2)
	check(o, 2, d, 1)
	check(o, 3, b, 2)

	if o.Get(Int(a)).Int() != a {
		t.Fatal(o.items)
	}

	o = newMap(4)
	o.noresize = true
	a = randInt(8, 1)
	b = randInt(8, 1)
	c = randInt(8, 1)
	d = randInt(8, 2)
	e := randInt(8, 4)
	o.Set(Int(a), Int(a))
	o.Set(Int(b), Int(b))
	o.Set(Int(c), Int(c))
	o.Set(Int(d), Int(d))
	o.Set(Int(e), Int(e)) // [nil, a+0, b+1, c+2, d+2, e+1, nil, nil]
	check(o, 1, a, 0)
	check(o, 2, b, 1)
	check(o, 3, c, 2)
	check(o, 4, d, 2)
	check(o, 5, e, 1)

	o.Delete(Int(b))
	o.Delete(Int(d)) // [nil, a+0, deleted+1, c+2, deleted+2, e+1, nil, nil]
	if o.Get(Int(b)) != Nil {
		t.Fatal(o.items)
	}
	if o.Get(Int(c)) != Int(c) {
		t.Fatal(o.items)
	}
	o.Set(Int(d), Int(d)) // [nil, a+0, c+1, d+1, e+0, nil, nil, nil]
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 1)
	check(o, 4, e, 0)

	o.Delete(Int(e))
	o.Set(Int(e), Int(e))
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 1)
	check(o, 4, e, 0)
}

func BenchmarkRHMap10(b *testing.B)      { benchmarkRHMap(b, 10) }
func BenchmarkGoMap10(b *testing.B)      { benchmarkGoMap(b, 10) }
func BenchmarkRHMap20(b *testing.B)      { benchmarkRHMap(b, 20) }
func BenchmarkGoMap20(b *testing.B)      { benchmarkGoMap(b, 20) }
func BenchmarkRHMap50(b *testing.B)      { benchmarkRHMap(b, 50) }
func BenchmarkGoMap50(b *testing.B)      { benchmarkGoMap(b, 50) }
func BenchmarkRHMap5000(b *testing.B)    { benchmarkRHMap(b, 5000) }
func BenchmarkGoMap5000(b *testing.B)    { benchmarkGoMap(b, 5000) }
func BenchmarkRHMapUnc10(b *testing.B)   { benchmarkRHMapUnconstrainted(b, 10) }
func BenchmarkGoMapUnc10(b *testing.B)   { benchmarkGoMapUnconstrainted(b, 10) }
func BenchmarkRHMapUnc1000(b *testing.B) { benchmarkRHMapUnconstrainted(b, 1000) }
func BenchmarkGoMapUnc1000(b *testing.B) { benchmarkGoMapUnconstrainted(b, 1000) }

func benchmarkRHMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := newMap(n)
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
	m := newMap(1)
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

func TestRHMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := newMap(0)
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

	fmt.Println(m.Len(), m.Size(), len(m2))

	for k, v := range m2 {
		if m.Get(Int64(k)).Int64() != v {
			m.Foreach(func(mk Value, mv *Value) bool {
				if mk.Int64() == k {
					t.Log(mk, *mv)
				}
				return true
			})
			t.Fatal(m.Get(Int64(k)), k, v)
		}
	}

	if m.Len() != len(m2) {
		t.Fatal(m.Len(), len(m2))
	}

	for k, v := m.FindNext(Nil); k != Nil; k, v = m.FindNext(k) {
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

	for k, v := m.FindNext(Nil); k != Nil; k, v = m.FindNext(k) {
		fmt.Println(k, v)
	}
}

func TestMapDistance(t *testing.T) {
	test := func(sz int) {
		o := newMap(sz)
		for i := 0; i < sz; i++ {
			o.Set(Int(randInt(sz, i)), Int(i))
		}
		for _, i := range o.items {
			if i.key != Nil {
				if i.dist != 0 {
					t.Fatal(o.items)
				}
			}
		}
	}
	for i := 1; i < 16; i++ {
		test(i)
	}
	test = func(sz int) {
		o := newMap(sz / 2)
		o.noresize = true
		for i := 0; i < sz; i++ {
			o.Set(Int(randInt(sz, i)), Int(i))
		}
		for _, i := range o.items {
			if i.dist != 0 {
				t.Fatal(o.items)
			}
		}
	}
	for i := 2; i <= 16; i += 2 {
		test(i)
	}
}

func TestHashcodeDist(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for _, a := range []string{"a", "b", "c", "z", randString(), randString(), randString(), randString()} {
		fmt.Println(Str(a).HashCode() % 32)
	}

	z := map[uint32]int{}
	m := newMap(0)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1e6; i++ {
		v := Int64(int64(i)).HashCode() % 32
		z[v]++
		m.Set(Int(i), Int(i))
	}
	fmt.Println(z, m.density(), m.Size())

	z = map[uint32]int{}
	for i := 0; i < 1e6; i++ {
		v := Int64(rand.Int63()).HashCode() % 32
		z[v]++
	}
	fmt.Println((z))

	z = map[uint32]int{}
	for i := 0; i < 1e6; i++ {
		v := Str(randString()).HashCode() % 32
		z[v]++
	}
	fmt.Println((z))

	z = map[uint32]int{}
	m.Clear()
	for i := 0; i < 1e6; i++ {
		x := fmt.Sprintf("%016x", rand.Uint64())
		v := Str(x).HashCode() % 32
		z[v]++
		m.Set(Str(x), Str(x))
	}
	fmt.Println(z, m.density(), m.Size())

	m = newMap(0)
	for i := 0; i < 20; i++ {
		m.Set(Int(i), Int(i))
	}
	fmt.Println(m.DebugString())
}

func BenchmarkStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Str("\x00")
	}
}

func BenchmarkStrHashCode(b *testing.B) {
	x := Str(randString() + randString())
	for i := 0; i < b.N; i++ {
		x.HashCode()
	}
}

func BenchmarkContains(b *testing.B) {
	m := newMap(0)
	k2 := []Value{}
	for i := 0; i < 1e3; i++ {
		k := randString()
		m.Set(Str(k), Int(0))
		k2 = append(k2, Str(randString()))
	}
	for i := 0; i < b.N; i++ {
		if m.Contains(k2[rand.Intn(len(k2))]) {
			b.Fatal()
		}
	}
}

func BenchmarkContainsNative(b *testing.B) {
	k2 := []string{}
	m := map[string]bool{}
	for i := 0; i < 1e3; i++ {
		k := randString()
		m[k] = true
		k2 = append(k2, randString())
	}
	for i := 0; i < b.N; i++ {
		if m[k2[rand.Intn(len(k2))]] {
			b.Fatal()
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
	assert(Str("\x00").IsTrue())
	assert(Str("\x00\x00").IsTrue())
	assert(Str("\x00\x00\x00").IsTrue())
	assert(Str("\x00\x00\x00\x00").IsTrue())
	assert(Str("\x00\x00\x00\x00\x00").IsTrue())
	assert(Str(strings.Repeat("\x00", 6)).IsTrue())
	assert(Str(strings.Repeat("\x00", 7)).IsTrue())
	assert(Str(strings.Repeat("\x00", 8)).IsTrue())
	assert(Byte('a') == Str("a"))
	assert(Rune('a') == Str("a"))
	assert(Rune('\u263a') == Str("\u263a"))
	assert(Rune('\U0001f60a') == Str("\U0001f60a"))
	assert(Bytes(nil).IsTrue())
	assert(Bytes([]byte("")).IsTrue())
	assert(!ValueOf([]byte("")).IsFalse())
	assert(Str("\x00\x00\x00\x00\x00\x00\x00").Less(Str("\x00\x00\x00\x00\x00\x00\x00\x00")))
	assert(newArray().ToValue().IsArray())
}

func TestFormat(t *testing.T) {
	sprintf := func(format string, args ...interface{}) string {
		p := &bytes.Buffer{}
		internal.Fprintf(p, format, args...)
		return p.String()
	}

	type payload struct {
		f   string
		v   interface{}
		res string
	}

	payloads := []payload{
		{"%d", uint(12345), "12345"},
		{"%d", int(-12345), "-12345"},
		{"%d", ^uint8(0), "255"},
		{"%d", ^uint16(0), "65535"},
		{"%d", ^uint32(0), "4294967295"},
		{"%d", ^uint64(0), "18446744073709551615"},
		{"%d", int8(-1 << 7), "-128"},
		{"%d", int16(-1 << 15), "-32768"},
		{"%d", int32(-1 << 31), "-2147483648"},
		{"%d", int64(-1 << 63), "-9223372036854775808"},
		{"%.d", 0, ""},
		{"%.0d", 0, ""},
		{"%6.0d", 0, "      "},
		{"%06.0d", 0, "      "},
		{"% d", 12345, " 12345"},
		{"%+d", 12345, "+12345"},
		{"%+d", -12345, "-12345"},
		{"%b", 7, "111"},
		{"%b", -6, "-110"},
		{"%#b", 7, "0b111"},
		{"%#b", -6, "-0b110"},
		{"%b", ^uint32(0), "11111111111111111111111111111111"},
		{"%b", ^uint64(0), "1111111111111111111111111111111111111111111111111111111111111111"},
		{"%o", 01234, "1234"},
		{"%o", -01234, "-1234"},
		{"%#o", 01234, "01234"},
		{"%#o", -01234, "-01234"},
		{"%O", 01234, "0o1234"},
		{"%O", -01234, "-0o1234"},
		{"%o", ^uint32(0), "37777777777"},
		{"%o", ^uint64(0), "1777777777777777777777"},
		{"%#X", 0, "0X0"},
		{"%x", 0x12abcdef, "12abcdef"},
		{"%X", 0x12abcdef, "12ABCDEF"},
		{"%x", ^uint32(0), "ffffffff"},
		{"%X", ^uint64(0), "FFFFFFFFFFFFFFFF"},
		{"%.20b", 7, "00000000000000000111"},
		{"%10d", 12345, "     12345"},
		{"%10d", -12345, "    -12345"},
		{"%+10d", 12345, "    +12345"},
		{"%010d", 12345, "0000012345"},
		{"%010d", -12345, "-000012345"},
		{"%20.8d", 1234, "            00001234"},
		{"%20.8d", -1234, "           -00001234"},
		{"%020.8d", 1234, "            00001234"},
		{"%020.8d", -1234, "           -00001234"},
		{"%-20.8d", 1234, "00001234            "},
		{"%-20.8d", -1234, "-00001234           "},
		{"%-#20.8x", 0x1234abc, "0x01234abc          "},
		{"%-#20.8X", 0x1234abc, "0X01234ABC          "},
		{"%-#20.8o", 01234, "00001234            "},
		{"%+.3e", 0.0, "+0.000e+00"},
		{"%+.3e", 1.0, "+1.000e+00"},
		{"%+.3x", 0.0, "+0x0.000p+00"},
		{"%+.3x", 1.0, "+0x1.000p+00"},
		{"%+.3f", -1.0, "-1.000"},
		{"%+.3F", -1.0, "-1.000"},
		{"%+.3F", float32(-1.0), "-1.000"},
		{"%+07.2f", 1.0, "+001.00"},
		{"%+07.2f", -1.0, "-001.00"},
		{"%-07.2f", 1.0, "1.00   "},
		{"%-07.2f", -1.0, "-1.00  "},
		{"%+-07.2f", 1.0, "+1.00  "},
		{"%+-07.2f", -1.0, "-1.00  "},
		{"%-+07.2f", 1.0, "+1.00  "},
		{"%-+07.2f", -1.0, "-1.00  "},
		{"%+10.2f", +1.0, "     +1.00"},
		{"%+10.2f", -1.0, "     -1.00"},
		{"% .3E", -1.0, "-1.000E+00"},
		{"% .3e", 1.0, " 1.000e+00"},
		{"% .3X", -1.0, "-0X1.000P+00"},
		{"% .3x", 1.0, " 0x1.000p+00"},
		{"%+.3g", 0.0, "+0"},
		{"%+.3g", 1.0, "+1"},
		{"%+.3g", -1.0, "-1"},
		{"% .3g", -1.0, "-1"},
		{"% .3g", 1.0, " 1"},
		{"%b", float32(1.0), "8388608p-23"},
		{"%b", 1.0, "4503599627370496p-52"},
		// Test sharp flag used with floats.
		{"%#g", 1e-323, "1.00000e-323"},
		{"%#g", -1.0, "-1.00000"},
		{"%#g", 1.1, "1.10000"},
		{"%#g", 123456.0, "123456."},
		{"%#g", 1234567.0, "1.234567e+06"},
		{"%#g", 1230000.0, "1.23000e+06"},
		{"%#g", 1000000.0, "1.00000e+06"},
		{"%#.0f", 1.0, "1."},
		{"%#.0e", 1.0, "1.e+00"},
		{"%#.0x", 1.0, "0x1.p+00"},
		{"%#.0g", 1.0, "1."},
		{"%#.0g", 1100000.0, "1.e+06"},
		{"%#.4f", 1.0, "1.0000"},
		{"%#.4e", 1.0, "1.0000e+00"},
		{"%#.4x", 1.0, "0x1.0000p+00"},
		{"%#.4g", 1.0, "1.000"},
		{"%#.4g", 100000.0, "1.000e+05"},
		{"%#.4g", 1.234, "1.234"},
		{"%#.4g", 0.1234, "0.1234"},
		{"%#.4g", 1.23, "1.230"},
		{"%#.4g", 0.123, "0.1230"},
		{"%#.4g", 1.2, "1.200"},
		{"%#.4g", 0.12, "0.1200"},
		{"%#.4g", 10.2, "10.20"},
		{"%#.4g", 0.0, "0.000"},
		{"%#.4g", 0.012, "0.01200"},
		{"%#.0f", 123.0, "123."},
		{"%#.0e", 123.0, "1.e+02"},
		{"%#.0x", 123.0, "0x1.p+07"},
		{"%#.0g", 123.0, "1.e+02"},
		{"%#.4f", 123.0, "123.0000"},
		{"%#.4e", 123.0, "1.2300e+02"},
		{"%#.4x", 123.0, "0x1.ec00p+06"},
		{"%#.4g", 123.0, "123.0"},
		{"%#.4g", 123000.0, "1.230e+05"},
		{"%#9.4g", 1.0, "    1.000"},
		// The sharp flag has no effect for binary float format.
		{"%#b", 1.0, "4503599627370496p-52"},
		// Precision has no effect for binary float format.
		{"%.4b", float32(1.0), "8388608p-23"},
		{"%.4b", -1.0, "-4503599627370496p-52"},
		// Test correct f.intbuf boundary checks.
		// float infinites and NaNs
		{"%f", math.Inf(1), "+Inf"},
		{"%.1f", math.Inf(-1), "-Inf"},
		{"% f", math.NaN(), " NaN"},
		{"%20f", math.Inf(1), "                +Inf"},
		{"% 20F", math.Inf(1), "                 Inf"},
		{"% 20e", math.Inf(-1), "                -Inf"},
		{"% 20x", math.Inf(-1), "                -Inf"},
		{"%+20E", math.Inf(-1), "                -Inf"},
		{"%+20X", math.Inf(-1), "                -Inf"},
		{"% +20g", math.Inf(-1), "                -Inf"},
		{"%+-20G", math.Inf(1), "+Inf                "},
		{"%20e", math.NaN(), "                 NaN"},
		{"%20x", math.NaN(), "                 NaN"},
		{"% +20E", math.NaN(), "                +NaN"},
		{"% +20X", math.NaN(), "                +NaN"},
		{"% -20g", math.NaN(), " NaN                "},
		{"%+-20G", math.NaN(), "+NaN                "},
		// Zero padding does not apply to infinities and NaN.
		{"%+020e", math.Inf(1), "                +Inf"},
		{"%+020x", math.Inf(1), "                +Inf"},
		{"%-020f", math.Inf(-1), "-Inf                "},
		{"%-020E", math.NaN(), "NaN                 "},
		{"%-020X", math.NaN(), "NaN                 "},
	}

	for _, p := range payloads {
		if v := sprintf(p.f, p.v); v != p.res {
			t.Fatal(p.f, p.v, "->", v, p.res)
		}
	}
}

func TestShape(t *testing.T) {
	assertError := func(f bool, err error) {
		if (err == nil) == f {
			t.Fatal(err, string(debug.Stack()))
		}
	}

	NewShape("")
	// Shape(",")
	// Shape("[")
	// Shape("]")
	// Shape("[,]")

	assertError(false, NewShape("i")(Int(1)))
	assertError(false, NewShape("n")(Int(1)))
	assertError(false, NewShape("n")(Float64(1)))
	assertError(false, NewShape("(i)")(Array(Int(1))))
	assertError(false, NewShape("(i i)")(Array(Int(1), Int(2))))
	assertError(false, NewShape("[i i]")(Array(Int(1), Int(2))))
	assertError(false, NewShape("[i i]")(Array(Int(1), Int(2), Int(3), Int(4))))
	assertError(false, NewShape("[i,i]")(Array()))
	assertError(false, NewShape("(i,is)")(Array(Int(1), Int(2))))
	assertError(false, NewShape("(i,is)")(Array(Int(1), Str("2"))))
	assertError(false, NewShape("([] is)")(Array(Array(Int(1)), Str("2"))))
	assertError(false, NewShape("([@*int] is)")(Array(Array(NewNative(new(int)).ToValue()), Str("2"))))
	assertError(false, NewShape("E")(Error(nil, fmt.Errorf("test"))))
	assertError(false, NewShape("@error")(Error(nil, fmt.Errorf("test"))))
	assertError(true, NewShape("[i]")(Array(Int(1), Float64(0.5))))
	assertError(false, NewShape("{} ")(NewObject(0).ToValue()))

	o := NewObject(10)
	o.Set(Int(1), Array())
	o.Set(Int(2), Array())
	assertError(false, NewShape("({i:[]} is)")(Array(o.ToValue(), Str("2"))))

	o.Clear()
	o.Set(Int(1), True)
	o.Set(Int(2), Str(""))
	assertError(false, NewShape("{}")(o.ToValue()))
	assertError(false, NewShape("{i}")(o.ToValue()))
	assertError(false, NewShape("{i v}")(o.ToValue()))

	o.Clear()
	assertError(false, NewShape("{i}")(o.ToValue()))
	assertError(false, NewShape("<i,{}>")(o.ToValue()))
	assertError(false, NewShape("<i,@object>")(o.ToValue()))

	p := NewNamedObject("test", 0)
	assertError(true, NewShape("@test")(o.ToValue()))
	o.SetPrototype(p)
	assertError(false, NewShape("@test")(o.ToValue()))
	assertError(false, NewShape("@object")(o.ToValue()))

	assertError(false, NewShape("R")(ValueOf(&bytes.Buffer{})))
	assertError(true, NewShape("C")(ValueOf(&bytes.Buffer{})))

	assertError(false, NewShape("<[],{}>")(o.ToValue()))
}

func BenchmarkShape(b *testing.B) {
	v := Array(Int(1), Int(2))
	s := NewShape("[i,i]")
	for i := 0; i < b.N; i++ {
		s(v)
	}
}

func BenchmarkShapeSimple(b *testing.B) {
	v := Int(1)
	for i := 0; i < b.N; i++ {
		v.AssertShape("si", "")
	}
}

func BenchmarkShapeType(b *testing.B) {
	v := Int(1)
	for i := 0; i < b.N; i++ {
		switch v.Type() {
		case typ.Number, typ.String:
			v.AssertNumber("")
		}
	}
}
