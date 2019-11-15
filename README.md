potatolang (pol) is a golang-dialect script language written in golang itself. Currently it only runs on 64bit platforms.

For benchmarks, refer to [here](https://github.com/coyove/potatolang/blob/master/tests/bench/perf.md).

## Quick starter guide for gophers

|Basic concept |Golang equivalent|
|--------|------|
|Type `Nil`     | nil |
|Type `Number`  | float64 |
|Type `String`  | immutable []byte |
|Type `Slice`   | []Value |
|Type `Pointer` | unsafe.Pointer |
|Type `Closure` | func |
|Type `Struct`  | immutable map[string]Value |
|NOT SUPPORTED | mutable map[Value]Value |
| `true == 1` and `false == 0` | bool |
|No need to declare first, just define `a = 1` | `a := 1` |
|Refer defined variables: `b = 1; a = b` | `b := 1; a := b`|
|`for i = v {}               ` |`for i = 0; i < len(v); i++ {}`|
|`for i = start, end {}      ` |`for i = start; i < end; i++ {}`|
|`for i = start, end, step {}` |`for i = start; i <= end; i += step {}`|
|Basically identical | Other for-loop syntax|
|`for true {}`  (true is required)| `for {}` |
|`a = {...}`| `a := []Value{...}` |
|`a = { k: ... }` | `a := struct{ k Value }{...}` |
|`a, b = b, a` | `a, b = b, a` |
|NOT SUPPORTED (n > 2) | `a1, ..., an = b1, ..., bn` |
|`a, b = (func(){return 1, 2})()` | `a, b = (func(){return 1, 2})()` |
|NOT SUPPORTED (n > 2) | `a1, ..., an = (func(){return b1, ..., bn})()` |
|`foo = func a = a + 1` | `foo := func(a) { return a + 1 }` |
|`addr := &var; addr[] = 1` | `addr := &var; *addr = 1` |
|`a && b ⎮⎮ c` | `if (a) { return b } else { return c }`|
|`a = a << {1}; a = a << {2,3}`| `a = append(a, 1, 2, 3)`|
|`a = append(a, {1, 2}...)`| `a = append(a, []Value{1, 2}...)`|
|`a[len(a)] = 1`|`a = append(a, 1)`|

### Variable Scope
Since we don't have declarations, to create a variable specifically inside a scope, we use `:=`:
```
func foo(b) {
    a := 1
    (func() {
        a := b
        io.println("inner: ", a)
    })()
    io.println("outer: ", a)
}
foo(2)
// outputs:
//      inner: 2
//      outer: 1
```
Note there are two exceptions as shown below where the topmost variable `a` is never touched:
```
a = 1
func foo(a) {
    a = 2 // a is local, because it's the parameter of foo
} 
foo(2)

func bar() {
    func a() {}
    a = 2 // closure is always local: a := func() {})
}
bar()
```

### Operators
Basically the same, note that:
1. All bitwise operators are applied on int32 operands except `>>>` (unsigned rsh) which works on uint32.
7. `Slice` and `Struct` can be automatically and recursively compared using `==` and `!=`.

### Struct
1. `Struct` are like `map` in golang, but once you initized it in code you can't add any more keys into it nor iterate it. So its behaviors are more like a `struct`.
2. To record the keys of a `Struct`, you can use `Named Struct`:
```
a = {
    k: nil,
    k2: 0,
}
a.__fields // nil
a = struct {
    k: nil,
    k2: 0,
}
a.__fields // ["k", "k2"]
```
