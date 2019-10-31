potatolang (pol) is a script language written in golang.

## Quick starter guide for gophers

### Basic Type
1. Nil (nil)
2. Number (float64)
3. String (string + []byte)
4. Map (map + slice)
5. Pointer (unsafe.Pointer)
6. Closure (func)
7. Phantom (special void value)
8. No real `bool` type, we have `true == 1` and `false == 0`

### Variable
1. No need to declare them, just write `a = 1` directly.
2. But you can only refer defined variables, e.g. `a = b` is illegal, should be `b = 1 a = b`.
3. To initiate an array, you write `a = {1, 2, 3}`, to initiate a map, you write `a = {"k": v}`.
4. Since we don't have declarations, to create a variable specically inside a scope, we usually prepend it with a `$`, e.g.:
```
func foo(b) {
    $a = 1
    (func() {
        $a = b
        io.println("inner: ", $a)
    })()
    io.println("outer: ", $a)
}
foo(2)
// outputs:
//      inner: 2
//      outer: 1
```
Note there are two exceptions as shown below where the variable `a` is never touched:
```
a = 1
func foo(a) {
    a = 2 // a is local, because it's the parameter of foo
} 
foo(2)

func bar() {
    func a() {}
    a = 2 // closures are always local, so here we are overriding it with '2'
}
bar()
```

### Phantom
The equivalent of `undefined` in JS, written as `#nil`.

### Operators
Basically the same, note that:
1. Bitwise not `^` is written as: `~`, just like C.
2. All bitwise operators are applied on int32 operands, `>>>` (unsigned rsh) is the only exception that works on uint32.
3. Lua trick: `a && b || c` => `if (a) { return b } else { return c }`
4. To delete a key from a map, assign the Phantom value to it: `m["key"] = #nil`.
5. To pop the last value of a slice, use `#` (as you may notice, pop a nil will give you the Phantom value), e.g.:
```
a = {1, 2, 3}
b = #a
// a == {1, 2} && b == 3
```
6. `Map` can be automatically and recursively compared using `==` and `!=`.

### Loop
Basically the same, with new syntax:
1. `for i = start, end { ... }` => ` for i := start; i < end; i++ { ... }`.
2. `for i = start, step, end { ... }` => `for i := start; i <= end; i += step { ... }` or `for i := start; i >= end; i += step { ... }`.
3. `for m, func (k, v) { ... }` => `for k, v := range m { .. }`, inside the callback, `return false` will terminate the iteration.

### String
Strings are mutable by syntax, but behind the stage we convert it to `[]byte` anyway, e.g.:
```
a = "hello"
a[0] = 'H'
// a == "Hello"
a[4] = "o world"
// a == "Hello world"
a[5] = ""
// a == "Helloworld"
```
