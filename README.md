# potatolang

potatolang is a script language written in golang. 

## Quick go through in 150 LOC

```javascript
/* 
 * Like Lua, array and map are all "map" in potatolang
 */

var a = {1, 2, 3}
// a[0] = 1, a[1] = 2, a[2] = 3

var a = {"key1": "value1", 2: "value2"}
// keys of the map can be any value

var a = {"key1": "value1"}
a[0] = 1  a[1] = 2  a[3] = "in the map"
// 'a' now contains an array of {1, 2} and a map of {"key1": "value1", 3: "in the map"}
// '3' is inside the map because it is out of the array's index range.
// in this case, the valid index range is: [0, 2], anything outside it will go into the map.

/* 
 * Use `&` to concatenate two values, `+` to add or append values
 */

var a = 1
var b = a + "2"   // runtime panic
var c = {} + a    // c == { 1 }
var d = { 0 } + c // d == { 0, { 1 } }
var e = { 0 } & c // e == { 0, 1 }
var g = 7 & 8     // int32 bitwise and: 0
var f = "" & 1    // "1", this is the de facto tostring() in potatolang
var h = 0 & "1"   // 1, this is the de facto tonumber() in potatolang

// note '+' and '&' on map will modify the original value
var a = { 1, 2 }
var b = a + 3
assert a == b and a == { 1, 2, 3 } and b == { 1, 2, 3 }

/*
 * Builtin function 'copy'
 */

// 'copy' does shallow copy of a value
var a = { 1, 2, {1, 2} }
var b = copy(a)
a[0] = 0
a[2][0] = 0
assert b[0] == 1 and b[2][0] == 0

// it is the only way to iterate over a map in pol:
var m = {"1": 1, "2": 2 }
copy(m, fun(k, v) { assert k == ("" & v) })

var c = copy(a, fun i, n = n + 1)
assert c == { 2, 3 } and a == { 1, 2 }

// varargs:
fun sum() {
    // copy() without arguments will return a copy of the current execution stack
    var x = copy() // normally this line MUST be the first line of the whole function
    var s = 0
    copy(x, fun(i, n) {s = s + n})
    return s
}
assert sum(1, 2, 3) == 6
assert sum("a", "b", "c") == "abc"

// string is immutable, copy(str) will return an array of its bytes:
var a = "text"
a[0] = 96  // won't work
var b = copy(a)
assert typeof(b, "map")            // ok
assert b == {0x74,0x65,0x73,0x74}  // ok

// don't remove items when iterating an array
var a = {1, 2, 3}
copy(a, fun i = std.remove(a, i))

// but you can remove items when iterating a map (only) since it's an expected behavior in golang
var a = {1: 1, 2: 2, 3: 3}
copy(a, fun i = std.remove(a, i))

/*
 * Builtin operator 'yield'
 */

fun a() {
    yield 1 yield 2 yield 3
    // if you don't explicitly return, function itself returns nil here
    // return nil
}

var b = copy(a)
assert a() == 1 and a() == 2 and a() == 3 and a() == nil
assert b() == 1 and b() == 2 and b() == 3 and b() == nil
// now a and b are back to the start state

/*
 * Use 'this' as a parameter to simulate member functions
 */

var counter = {
    "tick": 0,
    "add": fun (step, this) { this.tick = this.tick + step }
}

var c = copy(counter)
c.add(1)
c.add(1)
assert c.tick == 2
// note that the order of 'this' and other arguments does not matter, e.g.: 
// fun(this, a, b) and fun(a, this, b) are the same,
// both of them will be compiled into: fun(a, b, this)

/*
 * Only 'fun' creates a new namespace
 */

fun foo() {
    var a = 0
    if ... var a = 1 else var a = 2
    for ... var a = 3
    while ... var a = 4

    {
        var a = 5
        {
            var a = 6
        }
    }
    return a
}
// all 'a' are the same 'a', foo() == 6

/* 
 * Other things worth mentioning
 */

// there is no boolean value, use '1' or '0' instead, actually, strings like "true" or "false" are also recommended
// there is no 'switch' statement, write 'else if' instead
// there is no '+=, -=, *= ...', write 'a = a + b' instead, however you can write 'a++' or 'a--', they have special optimizations
// there is no conditional operator 'a ? b : c', use and-or trick instead: 'a and b or c'
// bitwise operations are identical to javascript 
```