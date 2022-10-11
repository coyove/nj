assert(G == "test")

assert(-1, -1)
assert(([1,0,2])[1]-1, -1)
assert(({a=0}).a-1, -1)
assert((0)-1, -1)
assert((2)+-1, 1)
assert(({a=-1}).a, -1)
a = 1
assert(a-2, -1)

do
	local j = { a= 1, b= 2, array=[1, 2, {inner="inner"}]}
	assert(j.b == 2)
    j = json.parse("{\"a\":[[1]]}")
    assert(j.a[0][0], 1)
end


print("s2")

assert(true and true or false)
assert(false and false or true)

function deepadd(v)
    if v == 1e6 then
        return v
    end
    return deepadd(v + 1)
end
assert(deepadd(0) == 1e6)

do
    local sum = 0
    for i = 1,10 do
        for j=1,10,2 do
           sum = sum * j
           if j == 3 then break end
        end
        if i == 5 then break end
        sum = sum + i
    end
    assert(sum == 174)

    sum = 0
    for _, v in [1,2,3,4,5,6] do
        if v == 5 then continue end
        sum = sum + v
    end
    assert(sum == 16)

    sum = 0
    lst = [1,2,3,4,5,6,7,8,9]
    local k = 'not touched'
    for k, v in lst do
        if v%3 == 1 then continue end
        sum = sum + v
    end
    assert(sum, 2 + 3 + 5 + 6 + 8 + 9)
    assert(k, 'not touched')
end

do
	local arr = [1,2,3]
	for i =0,#(arr) do
	    assert( i+1 == arr[i])
	end
end

function f(a, b) assert(true); return [ b, a ] end
    function g(a, b ,c)
    b, c = f(b, c)
    return [c, b, a]
    end

do
    a,b,c = g(1,2,3)
    assert(a == 2 and b == 3 and c == 1)
end

function syntax_return_void() return end
function syntax_return_value(en) return en end
function syntax_return_void2(en) return
en end
assert(syntax_return_void2(1), nil)

ex = [ 1, 2 ]
ex[0], ex[1], ex3 = nil, ex[0], ex[1]
assert(ex3 == 2 and ex[0] == nil and ex[1] == 1)

assert(0x7fff_fffffffffffd < 0x7fffffffffffffff)
assert(0x7fff_ffff_ffff_fffd + 2 == 0x7fff_ffff_ffff_ffff)
assert(0x7fff_fffffffffffe + 1 == 0x7fffffffffffffff)

do
    scope = 1
    assert(scope == 1)
end
assert(scope == 1)

function callstr(a) 
    return a + "a"
end

assert(callstr("a") == "aa")

a = 0
assert(a == 0)

local a , b = 1, 2
assert(a == 1)
assert(b == 2)

if false then
    assert(false)
elseif a == 1 then
    local a = 3
    a = a + 2 - b
    assert(a == 3)
elseif true then
    assert(false)
end

assert(a == 1)

function add(a,b) return a + b end

function fib2(a, b, n)
    x = []
    while n do
        c = add(a, b)
        a = b
        b = c
        x.append(c)
        n=n-1
    end
    return x
end

do
    fib_seq = [1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711, 28657, 46368, 75025, 121393, 196418, 317811, 514229, 832040, 1346269, 2178309, 3524578, 5702887, 9227465]

    local s = fib2(0, 1, 33)
    for i = 0,33 do
        print(s[i])
        assert(fib_seq[i] == s[i])
    end
end

function deepadd2(a)
    if (a <= 0) then return 0 end
    return add(a, deepadd2(a - 1))
end

e = 100000
assert( deepadd2(e) == (1 + e) * (e / 2 ))
 
a = 2
assert( 1 + 2 *3/4 == 2.5)
assert( 1 /2+ 2.5 * (a + 1) * 5 == 38)
assert((a + 1) % 2 == 1)
assert(math.mod(a + 1.25, 2) == 1.25)
assert(math.remainder(a + 1.25, 2) == -0.75)
assert(math.mod(-a - 1.25, 2) == -1.25)

do
    local mark = 0
    for i=0x7000000000000707,-1,-0x1000000000000101 do
	mark=mark+0x1000000000000001
	println(i, mark)
    end
    assert(mark == 0x8000000000000008)
end

assert(0==nativeVarargTest())
assert(1==nativeVarargTest(1))
assert(2==nativeVarargTest(1, 2))
assert("10"==nativeVarargTest2("1"))
assert("11"==nativeVarargTest2("1", 2))

assert(intAlias(1e9).Unix() == 1)
print(time.now().Format("2006-01-02"))

println("boolean test")
boolConvert(true)

function returnnil()
local _ = 1 + 1
end
assert(returnnil() == nil)

G_FLAG = 1
findGlobal()
assert(G_FLAG == 'ok')

flag = "G_FLA" + "G"
print(flag, {}, {a=2}, {"a": 2})

function double(args...)
    for i=0,#(args) do
        args[i] = args[i] * 2
    end
    return args
end

x = double()
assert(#(x) == 0)
x = double(1)
assert(x[0] == 2)
x = double(1, 2, 3)
assert(x[0] + x[1] + x[2] == 12)
print(x)

function double2(k, args...)
    for i=0,#(args) do
        args[i] = args[i] * k
    end
    return args
end

x = double2(3, 1)
assert(x[0] == 3)
x = double2(4, 1, 2, 3)
assert(x[0] + x[1] + x[2] == 24)

function test_return(a)
end_ = 10+a
return end_
end
assert(test_return((1)),11)
assert(test_return(((2))),12)


function bar() return 2 end
function foo() return bar end
assert(foo()(), 2)

-- test TLParen parsing
foo()
({a=1}).a  -- should run flawlessly
({a=(1+bar())}).a

assert((function(x) return (x+1) end)(1), 2)
assert((function(x) end)(1), nil)
assert((function(x) return x.a end)({a=100}), 100)

a = 0 
a += 1
assert(a, 1)
a -= 2
assert(a, -1)
a*=a
assert(a, 1)
a=[a]
a[0]/=2
assert(a[0], 0.5)

a[0]*=4
a[0]<<=3
assert(a[0], 16)
a[0] %= 7
assert(a[0], 2)

a=-a[0]
assert(a, -2)
a=a -1
assert(a, -2)
a=a- 1
assert(a, -3)
a=a - 1
assert(a, -4)
a=-[a,a - 1][1]
assert(a, 5)
a=-function() return -2 end()
assert(a, 2)

a=-(a-1)-a
assert(a, -3)
a=1-a-a
assert(a, 7)
a=1-a -a
assert(a, -6)
a=1-a - a
assert(a, 13)
a=1-
a
assert(a, -12)
a=a-
-1
assert(a, -11)
a=[a][0]- -1
assert(a, -10)
a=a- --comment
-1
assert(a, -9)
a=a - --comment
-1
assert(a, -8)
a=a -
-(a)
assert(a, -16)
a=a+
-a
assert(a, 0)

assert(if(true, 10, panic("bad")), 10)

if (false)then end
if [0][0]then end

assert(if(true, "a", "")[0], "a".decodeutf8()[0])

m = gomap({a=1,b=2}, "a", 10)
assert(m.a , 10)

cnt = 0
for k,v in m do
    print(k, v)
    cnt += 1
    assert(v, if(k=='a', 10, if(k=='b', 2, nil)))
end
assert(cnt, 2)

m.c = 3
m.d = 4
for k,v in m do
if k == 'b' then m.delete(k) end
end
assert(#m, 3)
assert(not m.contains('b'))

function func1() end
function func2() end
function func1.a() end
function func2.a() end

a = 0
assert(a < 1)
assert(a < 1<<15-1)
assert(a < 1<<15)
assert(a < 1<<15+1)

a = 1e7
assert(1 < a)
assert(1<<15-1 < a)
assert(1<<15 < a)
assert(1<<15+1 < a)

assert(a-1, 9999999)
assert(1-a, -9999999)
assert(a+1, 10000001)
assert(1+a, 10000001)
assert(a+1<<15, 1e7+1<<15)
assert(1<<15-a, 1<<15-1e7)

tailcallobj = {}
function tailcallobj.foo(a) this.val += a.val; if this.val > 100 then return this.val end; return a.foo(this) end

tailcallobj1 = new(tailcallobj, {val=1})
tailcallobj2 = new(tailcallobj, {val=1})

-- 1+1+3+8+21+55...
-- 1+2+5+13+34+89...
assert(tailcallobj1.foo(tailcallobj2), 144)

function vargdefault(a... b0,b1)
    if b1 then return (a + b0) * b1 end
    if b0 then return a + b0 end
    return a * 2
end

assert(vargdefault(1), 2)
assert(vargdefault(1, 2), 3)
assert(vargdefault(1, 2, 3), 9)
assert(vargdefault(1, 2, 3, 4), 9)

function vargdefault2(... b0,b1,b2)
if not b0 and not b1 and not b2 then return 777 end
return b0 or b1 or b2
end
assert(vargdefault2(), 777)
assert(vargdefault2(10), 10)
assert(vargdefault2(false, 11), 11)
