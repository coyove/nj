# potatolang

potatolang is a C/js-like language written in golang. It only works on 64bit machine.

## Quick go through in 150 LOC

```javascript
/* 
 * Like Lua, array and map are all "map" in potatolang
 */

var a = {1, 2, 3};
// a[0] = 1, a[1] = 2, a[2] = 3

var a = {"key1": "value1", 2: "value2"};
// keys of the map can be any value

var a = {"key1": "value1"};
a[0] = 1;
a[1] = 2;
a[3] = "in the map";
// 'a' now contains an array of {1, 2} and a map of {"key1": "value1", 3: "in the map"}
// '3' is inside the map because it is out of the array's index range.
// in this case, the valid index range is: [0, 2], anything outside it will go into the map.

/* 
 * Use `&` to concatenate two values, `+` to add or append values
 */

var a = 1;
var b = a + "2";   // runtime panic
var c = {} + a;    // { 1 }
var d = { 0 } + c; // { 0, { 1 } }
var e = { 0 } & c; // { 0, 1 }
var g = 7 & 8;     // bitwise and: 0
var f = "" & 1;    // "1", this is the de facto tostring() in potatolang
var h = 0 & "1";   // 1, this is the de facto tonumber() in potatolang

/*
 * Builtin function 'copy'
 */

// 'copy' does shallow copy of a value
var a = { 1, 2, {1, 2} };
var b = copy(a);
a[0] = 0;
a[2][0] = 0;
assert b[0] == 1 && b[2][0] == 0;

// iterate over a map:
var m = {"1": 1, "2": 2 };
copy(m, func(k, v) { assert k == ("" & v); });

var c = copy(a, func(i, n) { return n + 1; });
assert c == { 2, 3 } && a == { 1, 2 };

// if you want to return multiple results, use this copy trick:
func foo(c) {
    var a = 1 + c;
    var b = 2 + c;
    // copy() without arguments will return a copy of the current stack
    return copy();
}
var r = foo(2);
// r is the stack of foo(2), the last two elements would be 'a' and 'b':
assert r[len(r)-2] == 3 && r[len(r)-1] == 4;

// the same trick can be used to accept varargs:
func sum() {
    var x = copy(); // this line must be the first line of the whole function
    var s = 0;
    copy(x, func(i, n) {s = s + n;});
    return s;
}
assert sum(1, 2, 3) == 6;
assert sum("a", "b", "c") == "abc";

// string is immutable, copy(str) will return an array of its bytes:
var a = "text";
a[0] = 96;  // won't work
var b = copy(a);
assert typeof(b, "map");            // ok
assert b == {0x74,0x65,0x73,0x74};  // ok

// don't remove items when iterating an array
var a = {1, 2, 3};
copy(a, func(i) { std.remove(a, i); });

// you can remove items when iterating a map, this is an expected behavior in golang
var a = {1: 1, 2: 2, 3: 3};
copy(a, func(i) { std.remove(a, i); });

/*
 * Builtin operator 'yield'
 */

func a() {
    yield 1; yield 2; yield 3;
    // if you don't explicitly return, func itself returns nil here
}

var b = copy(a);
assert a() == 1 && a() == 2 && a() == 3 && a() == nil;
assert b() == 1 && b() == 2 && b() == 3 && b() == nil;
// now a and b are back to the start state

/*
 * Use 'this' as a parameter to simulate member functions
 */

var counter = {
    "tick": 0,
    "add": func (step, this) { this.tick = this.tick + step; }
};

var c = copy(counter);
c.add(1);
c.add(1);
assert c.tick == 2;
// note that the order of 'this' and other arguments does not matter, e.g.: 
// func(this, a, b) and func(a, this, b) are the same,
// both of them will be compiled into: func(a, b, this)

/*
 * Statements 'if' and 'for' don't create a new namespace
 */

var a = 0;
if (...) var a = 1;
if (...) var a = 2;
for (...) var a = 3;
// all 'a' are the same 'a', now a == 3

/* 
 * Other things worth mentioning
 */

// null type is written as 'nil', not 'null'
// there is no boolean value, use '1' or '0' instead
// there is no 'switch' statement, write 'else if' instead
// there is no '++' or '--', write 'a = a + 1' instead, it has special optimization
// there is no conditional operator '?', use and-or trick instead: 'a && b || c'
// there is no 'while(cond)', write 'for(cond)' instead, but not 'for(;cond;)'
// there is no 'do while'
// bitwise operations are all based on signed 32bit integers
// there is no '>>>' 
// semi-colons can not be omitted
// 'var a = func foo() { ... };' is illegal
```