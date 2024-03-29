function makecls(f, c)
    a = function(x...)
        o = debug.self()
        return o._f(o._c, x...)
    end.copy()
    a.setproto({_f=f, _c=c}['setproto'](func))
    return a
end

a = makecls(function(x, step) x.c=x.c+step return x.c end, {c=0})
b = makecls(function(x, step) x.c=x.c+step return x.c end, {c=0})
assert(a(1), 1)
assert(b(2), 2)
assert(a(1), 2)
assert(b(2), 4)
assert(a is callable)

timer = {}
function timer.reset()
    this.c = 0
end
function timer.add(v)
    this.c = this.c + v + this.X
end
T1, T2 = new(timer, {c=0,X=10}), new(timer, {c=0,X=10})
T1.add(2)
T2.add(3)
assert(T1.c == 12 and T2.c == 13)

println(T1.c, T2.c)
T1.X = 100
println(T1.c, T2.c)
T1.add(20)
println(T1.c, T2.c)
T2.add(30)
assert(T1.c == 132 and T2.c == 53)

print(timer.add)

assert(structAddrTestEmbed.T.SetFoo, nil)
structAddrTestEmbed.natptrat('T').SetFoo(10)
assert(structAddrTestEmbed.T.A, 10)

function dict() return new(new(self)) end

function dict.next2()
    local k, v = this.next(this.proto()._iter)
    this.proto()._iter = k
    return k, v
end

m = dict()
m[0] = 0
m[true] = false
m.zz = 'zz'
while true do
    k, v = m.next2()
    print(k ,v)
    if k == nil then
        break
    end
    assert(m[k] == v)
    println(k ,v)
end

CarPrototype = {}
function CarPrototype.getBrand() return this.brand end
function CarPrototype.setBrand(brand) this.brand = brand; return this end
function CarPrototype.toString() return ('%s %s').format(this.color, this.brand) end

Car = new(CarPrototype, {brand=''})

carA = new(Car.copy(), {color='red'}).setBrand('ferrari')
carB = new(Car.copy(), {color='orange'}).setBrand('mclaren')

assert(carA.toString(), "red ferrari")
assert(carB.toString(), "orange mclaren")

cnt = 0
k = 99
for k, v in m do
    cnt = cnt + 1
    assert(m[k] == v)
end
assert(cnt == 3)
assert(k == 99) -- k shouldn't be altered

function worker()
    time.sleep(1)
    print("worker after 1s")
    return 'finished'
end
assert(worker.map([0], os.numcpus)[0] == 'finished')

-- closure
function foo(a)
    return function(b)
        return function(c)
            self.a = self.a + 1
            if c == 6 then panic(self.a + self.b + c) end
            return self.a + self.b + c
        end
    end
end
bar = foo(2)(3)
assert(bar(4), 10)
assert(bar(5), 12)

bar = foo(2)(3)
assert(bar(4), 10)
assert(bar(5), 12)
assert(bar.try(6).error(), 14)
local x, _, _ = bar.try(6)
assert(x.error(), 15)

function addcls(a)
    local a = 3
    return function(b)
        local b = 2
        return function(c)
            return self.a + self.b + c
        end
    end
end
assert(addcls(1)(2)(3), 8)

print.try("native try call")
assert(function(x) assert(x, 1) return x end.try(1) is (1))
assert(function(x) assert(x, 1) return x end.try(2) is error)

-- syntax test
m = {}
function m.a() return 'a' end
assert(m.a() == 'a')

cls = {}
function cls.__str() return this.v end
cls2 = {}
function cls2.__str() return this.v + '2' end
cls2 = new(cls, cls2)
obj = new(cls, {v='obj'})
obj2 = new(cls2, {v='obj'})
assert((obj).__str(), 'obj')
assert((obj2).__str(), 'obj2')
assert(obj2 is cls)
print(obj2)

t = {a=1, b=2}
t2 = t.copy(true)
t.b =3
assert(t2.b, 2)

t = ([1,2,3,4]).filter(function(v) return v >= 2 end)
assert(t[0] == 2 and t[2] == 4)


f1 = open("tests/a.txt", "w+")
f1.write("hello zzz")
f1.close()

f1 = open("tests/a.txt", "r")
f2 = open("tests/b.txt", "w+")
f2.write(f1.read(5))
f1.close()
f2.write(" world")
f2.close()

assert(open("tests/b.txt").read().unsafestr(), "hello world")
open.close()


f1 = open("tests/a.txt", "r")
f2 = open("tests/b.txt", "a+")
f2.pipe(f1)
f1.close()
f2.close()

assert(open("tests/b.txt").read().unsafestr(), "hello worldhello zzz")
open.close()

os.remove("tests/a.txt")
os.remove("tests/b.txt")


m = gomap({a=1,b=2}, 'c', 3)
print(type(m))
assert(m is nativemap)
assert(m.c , 3)
assert(#m, 3)
m.merge({d=4, e=5})
assert(m.contains('e'))
print(m.keys())


ch = channel(1)
ch.send(1)
assert(ch.recv()[0], 1)

ch2 = channel(1)
flag = false
function()
    ch.send(1)
    time.sleep(1)
    ch2.send(2)
    flag = true
end.go()

local rch, rv = channel.recvmulti([ch, ch2])
assert(rch, ch)
assert(rv, 1)
local rch, rv = channel.recvmulti([ch, ch2])
assert(rch, ch2)
assert(rv, 2)
assert(flag)

loadfile("tests/array.nj.lua")

assert(goVarg(10, function(a, b...)
    for _, b in b do
        a *= b
    end
    println(a, b)
    return a 
end), 10 * 11 * 12)

