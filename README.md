# potatolang

potatolang is a C/js-like language written in golang.

## Quick go through

1. Multi-declaration:

        var a, b;         // compile into: var a = nil; var b = nil;
        var a, b = 1, 2;  // compile into: var a = 1; var b = 2;
        var a, b, c = 1;  // compile into: var a = 1; var b = 1; var c = 1;
        a, b = 1, 2;      // illegal, the above are just special syntax sugar, you can't use it here

2. Array and map:

        // like Lua, array and map are all "map" in potatolang

        var a = {1, 2, 3};
        // a[0] = 1, a[1] = 2, a[2] = 3

        var a = {"key1": "value1", 2: "value2"};
        // keys of the map can be any value

        var a = {"key1": "value1"};
        a[0] = 1;
        // a now contains an array of {1} and a map of {key1: value1}

2. Builtin function `copy`:

        var a = { 1, 2 };
        var b = copy(a);
        a[0] = 0;
        assert b[0] == 1;  // d is another map now

        // iterate over a map:
        var m = {"1": 1, "2": 2 };
        copy(m, func(k, v) {
            assert k == ("" & v);
        });

        // provide a second argument to copy:
        var c = copy(a, func(i, n) { return n + 1; });
        assert c == { 2, 3 } and a == { 1, 2 };

2. Advanced `copy`:

        // if you want to return multiple results, use this trick:
        func foo(c) {
            var a = 1 + c;
            var b = 2 + c;

            // copy() without arguments will return a copy of the current stack
            return copy();
        }

        //         0   1   2
        //  foo  +---+---+---+
        // stack | c | a | b |
        //       +---+---+---+

        var r = foo(2);
        assert r[len(r)-2] == 3 and r[len(r)-1] == 4;

        // the same trick can be used to accept varargs:
        func sum() {
            var x = copy();
            var s = 0;
            copy(x, func(i, n) {s = s + n;});
            return s;
        }

        assert sum(1, 2, 3) == 6;
        assert sum("a", "b", "c") == "cba";  // goto "Misc" section to see why

        // string is immutable, copy(str) will return an array of its bytes:
        var a = "text";
        a[0] = 96;  // won't work

        var b = copy(a);
        assert typeof(b, "map");            // ok
        assert b == {0x74,0x65,0x73,0x74};  // ok

2. Yield:

        func a() {
            yield 1;
            yield 2;
            yield 3;
        }
        var b = copy(a);
        assert a() == 1 && a() == 2 && a() == 3 && a() == nil;
        assert b() == 1 && b() == 2 && b() == 3 && b() == nil;

        // a and b are dead, and back to the start state, so now: a() == 1 && b() == 1

2. Use `&` to concatenate two values, `+` to add or append values:

        var a = 1;
        var b = a + "2";         // runtime panic
        var c = {} + a;          // { 1 }
        var d = { 0 } + c;       // { 0, { 1 } }
        var e = { 0 } & c;       // { 0, 1 }
        var f = "" & a & "2";    // "12" 
                                 // this is the de facto tostring() in potatolang
        var g = 7 & 8;           // bit and: 0
        var h = 0 & "1";         // h = 1
                                 // this is the de facto tonumber() in potatolang

2. Use `this` as an argument name to simulate member functions:

        var counter = {
            "tick": 0,
            "add": func (step, this) {
                this.tick = this.tick + step;
            }
        };

        var c = copy(counter);
        c.add(1);
        c.add(1);
        assert c.tick == 2;

        // note that the order of 'this' and other arguments does not matter, e.g.: 
        // func(this, a, b) and func(a, this, b) are the same,
        // both of them will be compiled into: func(a, b, this)

2. Other than `func`, `if` and `for` don't create a new namespace:

        if (1) var a = 0; 
        if (1) var a = 1; 
        for (1) var a = 1;
        // all 'a' are the same 'a'
        // a == 1

2. Misc:

        There is no 'switch' statement, 'if' is enough;
        There is no '++' or '--' expressions;
        When using 'copy' to iterate over an array, it loops in a reversed order;
        And-or trick: 'a && b || c' means 'a ? b : c';
        Bitwise operations are all based on signed 32bit integers;
