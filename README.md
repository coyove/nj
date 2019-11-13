potatolang (pol) is a golang-dialect script language written in golang itself. Currently it only runs on 64bit platforms.

For benchmarks, refer to [here](https://github.com/coyove/potatolang/blob/master/tests/bench/perf.md).

## Quick starter guide for gophers

### Basic Type
1. Nil (nil)
2. Number (float64)
3. String (immutable []byte)
4. Slice ([]Value)
5. Pointer (unsafe.Pointer)
6. Closure (func)
7. Struct (immutable map[string]Value)
8. No real `bool` type, we have `true == 1` and `false == 0`

### Variable
1. No need to declare them, just write `a = 1` directly.
2. You can only refer defined variables, e.g. `a = b` is illegal, should be `b = <something> a = b`.
2. NO way to return multiple values.
3. To initiate a slice, you write `a = {1, 2, 3}`, to initiate a struct, you write `a = {k: 1}`. A struct's fields are immutable (more on that later):
```
a = { k : 1 }
a.k++
assert(a.k == 2)// ok
a.k2 = 2        // panic
```
4. Since we don't have declarations, to create a variable specifically inside a scope, we use `:=`:
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
    a = 2 // closure is always local (a := func() {})
}
bar()
```

### Operators
Basically the same, note that:
1. All bitwise operators are applied on int32 operands except `>>>` (unsigned rsh) which works on uint32.
2. Lua trick: `a && b || c` => `if (a) { return b } else { return c }`
3. `<<` can also be used to `append` some values, e.g. to delete a value inside a slice: 
```
a = {1, 2, 3} 
a = a[:1] << a[2:] // or if you prefer the builtin function append: a = append(a[:1], a[2:]...)
a == {1, 3}
```
6. However, to append a single value, this way is more preferred:
```
a = {1, 2, 3}
a[len(a)] = 4 // a == {1, 2, 3, 4}
a[10] = 10    // index out of range
```
7. `Slice` and `Struct` can be automatically and recursively compared using `==` and `!=`.

### Loop
Basically the same, with some new syntax:
1. `for i = v { ... }                => for i = 0; i < len(v); i++ { ... }`.
2. `for i = start, end { ... }       => for i = start; i < end; i++ { ... }`.
3. `for i = start, end, step { ... } => for i = start; i <= end; i += step { ... }`.
4. `for true { ... }`, unlike golang, don't forget `true`.

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
