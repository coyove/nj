# potatolang

potatolang is inspired by lua, particularly [gopher-lua](https://github.com/yuin/gopher-lua). However I am not really a fan of lua's syntax, so I decided to make my own one.

## Quick go through

1. `set a = 1` declares a variable named `a` with its value being number `1`. Value is a must so `set a` is not legal.
2. Declare before use. (I will mix the use of "declare" and "define")
2. There is no multi-assignment in potato. However you can write `set a, b = 1, 2`, but keep in mind this is just a syntax sugar (`set a, b = 1` will expand to `set a = 1 set b = 1`).
2. To define comlex structures: 

        set l = list 1, 2, 3 end
        set m = map "key1" = "value1", "key2" = "value2" end

    note that the key doesn't have to be an immediate value, but it must be a string:

        set m = map
            (function () return "key1" end) = "value1"
        end
    
2. To define a string block, use `sss<ident>` to start and `end<ident>` to end (`<ident>` should be a valid identifer name):

        set s = ssshello
            raw string literal
        endhello

2. `dup` is an important builtin function:

        set a = 1
        set b = dup(a)
        assert b == 1
        set c = list 1, 2 end
        set d = dup(c) 
        c[0] = 0    
        assert d[0] == 1
        set e = dup(d, function(i, n) return n + 1 end)
        assert e == list 2, 3 end
        set f = dup(d, function(i, n) if i == 1 then error(true) end return n + 1 end)
        assert f == list 2 end

2. If you want to return multiple results to the caller:

        function foo()
            set a = 1
            set b = 2

            # dup() without arguments will return a copy of the current stack
            # if dup() is right after return, then no actual copy will be done
            return dup()
        end

        set r = foo()
        assert r[-2] == 1 and r[-1] == 2

    The same trick can be used to accept varargs:

        function sum()
            set x = dup() 
            set s = 0
            dup(x, function(i, n) s = s + n end)
            return s
        end

        assert sum(1, -1) == 0
        assert sum(1, 2, 3) == 6 

2. Yield:

        function a()
            yield 1
            yield 2
            yield 3
        end

        set b = dup(a)

        assert a() == 1 and a() == 2 and a() == 3 and a() == nil
        assert b() == 1 and b() == 2 and b() == 3 and b() == nil

2. Strings and numbers can't be added together. however `&` can (concat):

        set a = 1
        set b = a + "2"         # panic
        set c = list end + a    # [1]
        set d = list 0 end + c  # [0, [1]]
        set e = list 0 end & c  # [0, 1]
        set f = "" & a & "2"    # "12"

2. To iterate over a map:

        set m = map
            "1" = 1,
            "2" = 2
        end
        dup(m, function(k, v)
            assert k == ("" & v)
        end)