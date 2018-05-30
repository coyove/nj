# potatolang

potatolang is inspired by lua, particularly [gopher-lua](https://github.com/yuin/gopher-lua), but with some heavy modifications to the syntax and other designs.

## Quick go through

1. Declare before use.

        a = 1             # compile error
        set a = 1         # declare first
        set b = "hello world"
        set c = nil
        set d = true
        set e = 0x123
        set f = { 1, 2, 3 }
        set g = { "key1" = "value1", "key2" = "value2" }
        set h = {
            (function () return "key1" end)() = "value1"
        }
        set i, j = 1, 2   # converted to set i = 1; set j = 2
        set k, l = 1      # converted to set k = 1; set l = 1
        a, b = 1, 2       # illegal, sorry
    
        # to define a string block, use sss<ident> to start and end<ident> to end. <ident> should be a valid identifer name:
        set s = ssshello
            ... raw string literal ...
        endhello

2. `dup` is an important builtin function, it can create a duplicate of a map or a string. During the process of duplication, you can map the value to a new one:

        set a = 1
        set b = dup(a)
        assert b == 1      # assert

        set c = { 1, 2 }
        set d = dup(c) 
        c[0] = 0    
        assert d[0] == 1   # d is another list now

        # to iterate over a map:
        set m = {
            "1" = 1,
            "2" = 2
        }
        dup(m, function(k, v)
            assert k == ("" & v)
        end)

        # map value in dup
        set e = dup(d, function(i, n) return n + 1 end)
        assert e == { 2, 3 } and d == { 1, 2 }

        # use error(...) to interrupt the dup
        set f = dup(d, function(i, n) 
            if i == 1 then error(true) end 
            return n + 1 
        end)
        assert f == {2}

2. Advance `dup`:

        # if you want to return multiple results:
        function foo(c)
            set a = 1 + c
            set b = 2 + c

            # dup() without arguments will return a copy of the current stack
            return dup()
        end

        #         0   1   2
        #  foo  +---+---+---+
        # stack | c | a | b |
        #       +---+---+---+

        set r = foo(2)
        assert r[len(r)-2] == 3 and r[len(r)-1] == 4

        # the same trick can be used to accept varargs:
        function sum()
            set x = dup() 
            set s = 0
            dup(x, function(i, n) s = s + n end)
            return s
        end

        assert sum(1, -1) == 0    # when calling dup() in sum(), since there is no
                                  # other variables yet on the stack, x = [1, -1]
        assert sum(1, 2, 3) == 6  # x = [1, 2, 3]

        # string is immutable, dup(str) will return an array of its bytes
        set a = "text"
        a[0] = 96                  # panic

        set b = dup(a)
        assert typeof(b, "map")    # ok
        b[0] = 96                  # ok

2. Yield:

        function a()
            yield 1
            yield 2
            yield 3
        end

        set b = dup(a)

        assert a() == 1 and a() == 2 and a() == 3 and a() == nil
        assert b() == 1 and b() == 2 and b() == 3 and b() == nil

        # now the life of a is terminated, everything starts from the begining

        assert a() == 1 and a() == 2 and a() == 3 and a() == nil

2. Strings and numbers can't be summed up. however `&` (concat) can:

        set a = 1
        set b = a + "2"         # runtime panic
        set c = {} + a          # { 1 }
        set d = { 0 } + c       # { 0, { 1 } }
        set e = { 0 } & c       # { 0, 1 }
        set f = "" & a & "2"    # "12" 
                                # this is the de facto tostring() in potatolang
        set g = 7 & 8           # bit and: 0

2. Closure knows whether it is a member of the map or not. This acknowledgement can be used to simulate member functions:

        set counter = map
            "tick" = 0,
            "add" = function()
                set this = who() # who's your dad
                this.tick = this.tick + 1
            end
        end

        set c = dup(counter)
        c.add()
        c.add()
        assert c.tick == 2

        set c2 = dup(c)
        c2.add()
        assert c2.tick == 3

        set c3 = dup(counter)
        c3.add()
        assert c3.tick == 1

2. potatolang can't return multiple values, nor throw an exception, its error handling will be like:

        function foo(n)
            if n < 2 then error("n is less than 2") end
            return n + 1
        end

        foo(1)
        assert error() == "n is less than 2"
        foo(2)
        assert error() == nil

        # note that calling error(...) will not interrupt the execution of its following code
        # you should manually return after it
        set ans = foo(1)
        assert error() == "n is less than 2"   # as expected
        assert ans == 2                        # still as expected (return n + 1)

        # also keep in mind: do check the error if the callee tends to return an error
        # otherwise the error gets propagated implicitly

        function bar(n)
            error("error" + n)
        end

        function foo(n)
            bar(n)
            # not calling error() to check the error
            # now the error (if any) belongs to foo(n)
        end

        bar(1)
        assert error() == "error1"   # as expected
        foo(1)
        assert error() == "error1"   # oops