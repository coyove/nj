potatolang (pol) is a script language written in golang.

## Quick starter guide for gophers

### Basic Type
1. Nil
2. Number (float64)
3. String
4. Map (map + slice)
5. Pointer (unsafe.Pointer)
6. Closure (func)
7. No real `bool` type, we have `true == 1` and `false == 0`

### Variable
1. No need to declare them, just write `a = 1` directly
2. But you can only refer defined variables, e.g. `a = b` is illegal, should be `b = 1 a = b`

#### Variable scope
Since we don't have declarations, to create a variable specically inside a scope, we use `$`:
```
b = 1

fun foo(a) { b = a }
fun bar(a) { $b = a }

foo(2) // b == 2
bar(3) // b is still 2, $b is only available in bar()
```
Note that `$` only makes sense inside `fun`, here is another example:
```
$x = 1
if true {
    $x = 2
}
// $x == 2
```

### Operators
Basically the same, note that:
1. Bitwise not `^` is written as: `~`, just like C
2. All bitwise operators are applied on int32 operands, `>>>` (unsigned rsh) is the only exception that works on uint32
3. Logical not `!`, and `&&`, or `||` are written as: `not`, `and`, `or`

### Closure
1. The keyword is `fun`, not `func`

### Loop
Basically the same, with new syntax:
1. `for i = start, end { ... }` => ` for i := start; i < end; i++ { ... }`
2. `for i = start, step, end { ... }` => `for i := start; i <= end; i += step { ... }` or `for i := start; i >= end; i += step { ... }`
3. `for m, fun (k, v) { ... }` => `for k, v := range m { .. }`, inside the callback, `return 0` will terminate the iteration
