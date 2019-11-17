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
|`m = map.New(n)`| mutable map[string]Value |
|`ch = chan.Make(n)`| chan Value |
|`chan.Send(ch, v)`| ch <- v |
|`v = chan.Recv(ch)`| v := <-ch |
|`v, ch = chan.Select(ch1, ch2, ... "default")`| select {...} |
| `true == 1` and `false == 0` | bool |
| `go(foo, arg1, arg2 ...)` | go foo(arg1, arg2, ...) |
|`for i = range start, end {}      ` |`for i = start; i <= end; i++ {}`|
|`for i = range start, end, step {}` |`for i = start; i <= end; i += step {}`|
|`for true {}`  (true is required)| `for {}` |
|`a = {...}`| `a := []Value{...}` |
|`a = { k: ... }` | `a := struct{ k Value }{...}` |
|`a1, a2 = (func(){return b1, b2})()` | `a1, a2 = (func(){return b1, b2})()` |
|NOT SUPPORTED (n > 2) | `a1, ..., an = (func(){return b1, ..., bn})()` |
|`foo = func a = a + 1` | `foo := func(a) { return a + 1 }` |
|`addr := &var; addr[] = 1` | `addr := &var; *addr = 1` |
|`a && b ⎮⎮ c` | `if (a) { return b } else { return c }`|
|`a = a << {1}; a = a << {2,3}`| `a = append(a, 1, 2, 3)`|
|`a = append(a, {1, 2}...)`| `a = append(a, []Value{1, 2}...)`|
|`a[len(a)] = 1`|`a = append(a, 1)`|

### Scope
Unlike golang, you can only create a new variable scope in a closure, which means the following code will yield `2`:
```
if true {
    a := 1
    go(func() {
        time.Sleep(time.Second)
        fmt.Println(a)
    })
}
a = 2
```

### Operators
Basically the same, note that:
1. All bitwise operators are applied on int32 operands except `>>>` (unsigned rsh) which works on uint32.
7. `Slice` and `Struct` can be automatically and recursively compared using `==` and `!=`.

### Struct
1. `Struct` are like `map` in golang, but once you initized it in code you can't add any more keys into it nor iterate it. So its behaviors are more like a `struct`.
2. To record the keys of a `Struct`, you can use `Named Struct` whose keys are stored in an extra field `__fields`:
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
