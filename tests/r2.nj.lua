a = 2
b = (a + 1) / (a - 1)
assert( b == 3)
assert( "sep0")

assert( (a + 1) / (a - 1) == 3)
assert( "sep1")

assert( 3 == (a + 1) / (a - 1))
assert( "sep2")

assert( (a + 1) == (a + 1) / (a - 1))
assert( "sep3")

function sub1(a) return a - 1 end
assert( (a + 1) == (a + 1) / sub1(a))
assert( "sep4")

assert( a + (a + (a + (a + 1))) == 9)
assert( "sep5")

assert( 1 + (1 + (1 + (1 + a))) == 6)
assert( "sep6")

c = [ 1,2,3,4 ]
assert( 1 + (1 + (1 + (1 + a))) + #(c) , 10)
assert( "sep7")

 a = 10
assert( 1 + (1 + a) == 12)

function foo1(a, b, c)
return a, b, c
end

local a, b, c = foo1(1, 2, 3)
assert(a == 1 and b + c == 5)

function foo2(a, b, c, x)
    local tmp = [a, b, c]
	if x then
		for i = 0, #(x) do tmp.append(x[i]) end
	end
    return tmp
end

local _, _, _, x,y = foo2(1,2,3,["4",5])
assert(x == "4" and y == 5)

function remove(el, arr)
    for i=0, #(arr) do
	    if arr[i] == el then
            for i=i,#(arr)-1 do
                arr[i] = arr[i+1]
            end
            arr[#(arr) - 1] = nil
	        return arr
	    end
    end
    return arr
end

do
	local r = remove(1, [1, 2, 3])
	assert(r[0] == 2 and r[1] == 3)
	local r = remove(2, [1, 2, 3])
	assert(r[0] == 1 and r[1] == 3)
	local r = remove(3, [1, 2, 3])
	println(r)
	assert(r[0] == 1 and r[1] == 2)
end

function bad(n) panic("bad" + n) end
do
	local err = bad.try("boy")
	assert(err.error() == "badboy")
end

function intersect(a, b)
	a.foreach(function(k, v)
		if not self.b.contains(k) then return false end
	end)
	return a
end

a = intersect({a=1, b=1, c=3}, {a=2, c:4})
assert(a.a, 1)
assert(a.c, 3)

assert(structAddrTest.Foo(), structAddrTest.A)
assert(structAddrTest2.Foo(), structAddrTest2.A)


function foo(a, b, c) return if(a,b,c) end

assert(foo(true, 10, 20), 10)

print('---R2---')

assert(1 is not object)
assert([1,2] is array)
assert([1,2] is native)
assert([1,2] is not nativemap)
assert([1,2] is not object)
assert(native is not object)


function foo3(k) debug.self()[k] = k end
for i=0,10 do foo3(i) end
for i=0,10 do assert(foo3[i], i) end

oneline = {}

function oneline.write(b)
    b = b.unsafestr()
    x = #b
    b = b.replace('\n', '').replace('\t', '')
    b = re([[\s+]]).replace(b, '')
    this.buf += b
    return x
end

ol = new(oneline, {buf=''})
os.shell("ls -al", {stdout=ol})
print(ol.buf)

ol = buffer()
os.shell("whoami", {stdout=ol})
print(ol.value().trim())

local bigint = createprototype("bigint", function(s)
    this.buf = bytes(0)
    for i=#s-1,-1,-1 do
        this.buf.append(s[i] - 48)
    end
end, {})

function bigint.tostring()
    local tmp = ''
    for _, b in this.buf do
        tmp = chr(48 + b) + tmp
    end
    return tmp
end

function bigint.equals(rhs)
    if #this.buf != #rhs.buf then return false end
    for i=0,#this.buf do
        if this.buf[i] != rhs.buf[i] then return false end
    end
    return true
end

function bigint.add(rhs)
    for i=#this.buf,#rhs.buf do this.buf.append(0) end
    local carry = 0
    for i=0,#this.buf do
        if i < #rhs.buf then
            this.buf[i] += rhs.buf[i] + carry
        else
            if carry == 0 then break end
            this.buf[i] += carry
        end

        if this.buf[i] < 10 then 
            carry = 0
            continue
        end
        carry = 1
        this.buf[i] -= 10
    end
    if carry != 0 then
        this.buf.append(1)
    end
    return this
end

assert(bigint("1234").tostring(), "1234")
assert(bigint("9").add(bigint("2")).equals(bigint("11")))
assert(bigint("9").add(bigint("99")).equals(bigint("108")))
assert(bigint("9999999999999999999999999999").add(bigint("2")).
equals(bigint("10000000000000000000000000001")))

local f = eval("a=1 return function(n) a +=n return a end")
assert(f(1), 2)
assert(f(2), 4)

function clsrec(a)
    a += 1
    function N()
        if self.a == 1 then return 1 end 
        self.a -= 1
        return self.a * self.N()
    end
    return N
end
assert(clsrec(5)(), 120)

function clsarray(a, n)
    local res = []
    for i=0,n do
        res.append(function()
            do
                return self.a * (self.i + 1)
            end
        end)
    end
    return res
end
assert(#clsarray.caplist(), 0)

for i, f in clsarray(10, 10) do
    assert(f(), (i + 1) * 10)
end

function createlinklist(a...)
    if #a == 0 then return end
    local head = [a[0], nil]
    local tmp = head
    for i=1,#a do
        tmp[1] = [a[i], nil]
        tmp = tmp[1]
    end
    return head
end

function reverselinklist(lst)
    if not lst then return lst end
    local dummy = ['lead', nil]
    while lst do
        local old = dummy[1]
        local x = [lst[0], old]
        dummy[1] = x
        lst = lst[1]
    end
    return dummy[1]
end

ll = createlinklist(1, 2, 3, 4)
print(ll, reverselinklist(ll))

function andortest(a, b, c)
    if c then return a[0] or b[1] end
    return a[0] and b[1]
end

assert(andortest([0], [0, 1], true), 1)
assert(andortest([0], [0, 1], false), 0)
assert(andortest([2], [0, "1"], false), "1")

function firstK(a, b, k)
    a.sort()
    b.sort()

    while k do
        if #a == 0 and #b == 0 then break end
        if #a == 0 then
            b = b[1:]
        elseif #b == 0 then
            a = a[1:]
        elseif a[0] < b[0] then
            a = a[1:]
        else
            b = b[1:]
        end
        k -= 1
    end

    if k >= #a and k >= #b then return "not found" end
    if #a and #b then return math.min(a[0], b[0]) end
    if #a then return a[0] end
    if #b then return b[0] end
    panic("shouldn't happen")
end

assert(firstK([1,2,3], [4,5,6], 1), 2)
assert(firstK([1,2,3], [0.4,1.5,2.6], 2), 1.5)
assert(firstK([1,2,3], [0.7], 2), 2)

local tmp = []
for i=0,1e3 do tmp += math.random() end
ki = int(math.random() * #tmp)
k = tmp.clone().sort()[ki]
a, b = tmp[:#tmp / 2], tmp[#tmp/ 2:]
assert(firstK(a, b, ki), k)

function fibwrapper() end
function fibwrapper.next(ab)
    if ab[0] == nil then
        this.v = 1
        return 1, this.v
    end
    this.v = ab[0] + ab[1]
    ab[0], ab[1] = ab[1], this.v
    return ab
end

local nw = createnativewrapper(fibwrapper)
res = []
for k, v in nw({}) do 
    res.append(v)
    if #res > 35 then break end
end
assert(res[35], 24157817)

function runtimejump(v)
    if v == 1 then return jump("label1") end
    if v == 2 then return jump("label2") end
    return jump("labeldefault")
end

runtimejump(3)
::label1:: assert(false)
::label2:: assert(false)
::labeldefault:: assert(true) print("runtimejump 3")

function test_jump()
    local ch1 = channel()
    local ch2 = channel()

    function() time.sleep(0.2) self.ch1.send("ch1") end.go()

    channel.recvmulti2(local out[2], {ch1: goto ch1, ch2: goto ch2})
    ::ch2::
    assert(false)
    if::ch1::then
    assert(out[0], "ch1")
    print(out)
    end
end
test_jump()
