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
        set f = list 1, 2, 3 end
        set g = map "key1" = "value1", "key2" = "value2" end
        set h = map
            (function () return "key1" end)() = "value1"
        end
        set i, j = 1, 2   # converted to set i = 1; set j = 2
        set k, l = 1      # converted to set k = 1; set l = 1
        a, b = 1, 2       # illegal, sorry
    
        # to define a string block, use sss<ident> to start and end<ident> to end. <ident> should be a valid identifer name:

        set s = ssshello
            ... raw string literal ...
        endhello

2. `dup` is an important builtin function:

        set a = 1
        set b = dup(a)
        assert b == 1

        set c = list 1, 2 end
        set d = dup(c) 
        c[0] = 0    
        assert d[0] == 1   # d is another list now

        # to iterate over a map:
        set m = map
            "1" = 1,
            "2" = 2
        end
        dup(m, function(k, v)
            assert k == ("" & v)
        end)

        # return value
        set e = dup(d, function(i, n) return n + 1 end)
        assert e == list 2, 3 end and d == list 1, 2 end

        set f = dup(d, function(i, n) 
            if i == 1 then error(true) end 
            return n + 1 
        end)
        assert f == list 2 end

        # if you want to return multiple results:
        function foo()
            set a = 1
            set b = 2

            # dup() without arguments will return a copy of the current stack
            # if dup() is right after return, then no actual copy will be done
            return dup()
        end

        set r = foo()
        assert r[-2] == 1 and r[-1] == 2

        # the same trick can be used to accept varargs:
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
        set b = a + "2"         # runtime panic
        set c = list end + a    # [1]
        set d = list 0 end + c  # [0, [1]]
        set e = list 0 end & c  # [0, 1]
        set f = "" & a & "2"    # "12" 
                                # in potato there is no tostring

2. When a closure is loaded from a map or a list, it remembers it, and this trick can be used to simulate structs:

        set counter = map
            "counter" = 0,
            "add" = function()
                set this = who()
                this.counter = this.counter + 1
            end
        end

        set c = dup(counter)
        c.add()
        c.add()
        assert c.counter == 2