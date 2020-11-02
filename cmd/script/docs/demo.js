  var demos = {
      "Select a demo...": `-- Print all global values, mainly functions
-- use doc(function) to view its documentation
local n, ...g = globals()

print(format("version {}, total global values: {}\\n", VERSION, n))

for i=1,#g do
    local name = str(g[i])
    if #name > 32 then
        name = name[1,16] .. '...' .. name[#name-16,#name]
    end

    title = i .. ": " .. name
    print((",----------------------------------------------------------------")[1,#title+3], '.')
    print("| ", title, " |")
    print(("'----------------------------------------------------------------")[1,#title+3], "'")
    print(type(g[i]) == "function" and doc(g[i]) or "\tN/A")
    print()
end
`,
/* = = = = = = = = */
      "fib": `function fib(n)
    if n == 0 then
        return 0
    elseif n == 1 then
        return 1
    end
    return fib(n-1) + fib(n-2)
end
return fib(10)
`,
/* = = = = = = = = */
      "goto": `goto inner
if false then
    ::inner::
    print("I'm in")
end
`,
/* = = = = = = = = */
      "yield": `function yieldable(n)
    while n > 0 do
        yield n
        n -= 1
    end
end

local c, ...state = yieldable(10)
while c do
    print(c)
    c, ...state = resume(yieldable, state)
end
`,
/* = = = = = = = = */
      "time": `println("Unix timestamp:", time())
println("Go time.Time:", Go_time().Format("2006-01-02 15:04:05"))
println(doc(Go_time))
`,
/* = = = = = = = = */
      "json": `local j = { a=1, b=2, array={ 1, 2, { inner="inner" } } }
--[[
There is no table type, code above will actually generate the
correspondent JSON STRING: '{"a":1,"b":2,"array":[1,2,{"inner":"inner"}]}'
]]

assert(json(j, "a") == 1)
assert(json(j, "b") == 2)
local n, a, b, c = json(j, "array")
assert(n == 3 and a == 1 and b == 2 and json(c, "inner") == "inner")
assert(json(j, "array.2.inner")=="inner")
println(json(j, "array"))

--[[
json() uses https://github.com/tidwall/gjson
Learn its syntax at https://github.com/tidwall/gjson/blob/master/SYNTAX.md
]]

-- Create JSON string from variables
local _, ...arr = array(1, 2, 3)
print({arr})

--[[
JSON dynamic object creation is a bit harder to write,
first we need to trick the parser with syntax "{ [nil] = something }" so it knows this is an object,
then we layout the key-value pairs sequentially, so it looks like:
{ [nil] = key1, value1, key2, value2, ... }
however the above sequence is misaligned: Nth key becomes Nth value and Nth value becomes (N+1)th key
so we have to prepend the sequence with a dummy value to correct the positions:
{ [nil] = dummy, key1, value1, key2, value2, ... }
]]
local n, ...kvpairs = array("key1", "value1", "key2", "value2")
println({ [nil] = array(kvpairs, "key3", "value3") })

`,
/* = = = = = = = = */
      "call": `function veryComplexFunction(a, b, c, d, e, f, g, H, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, Z)
    println(H, Z)
end
veryComplexFunction(Z="world", ["H"]="hello")

-- This is a trick from 'json' demo
local _, ...args = array("Z", "世界")
veryComplexFunction([nil] = array(args, "H", "你好"))

function foo() return "hello", "world" end
function bar() return "world", "hello" end

local ...a = random() > 0.5 and foo() or bar()
println(a)
`,
/* = = = = = = = = */
      "http": `local code, headers, body = http(
    method="POST",
    url="http://httpbin.org/post",
    query="k1=v1",
    query="k2=v2",
    json={
        name="Smith",
    },
)
println("code:", code)
println("headers:", headers)

if body then
    local data = json(body, "data")
    println("name:", json(data, "name"))
    println("args:", json(body, "args"))
end`,

      "goquery": `local code, _, body = http("GET", "https://example.com")
if iserror(code) then
    return "request failed: " .. code.error()
end

local el = goquery(body, "div").nodes()
for i = 1,#el do
    println(el[i].text())
end`,
      /* = = = = = = = = */
      "array-tricks": `--[[
There is no support for table, array, no anything similar,
but there are some tricks to simulate part of them (with performance penalty)
]]

local n, ...arr = array(1, 2, "3", 4)
println(n, "elements:")
for i = 1, #arr do
    print('* ', arr[i])
end

n, ...arr = array(arr, 5)
println("now there are", n, "elements:", arr)

function remove(el, ...arr)
    for i=1,#arr do
        if arr[i] == el then
            return arr[1,i-1], arr[i+1,#arr]
        end
    end
    return arr
end

println("remove 3 from array:", remove("3", arr))
println("remove nothing from array:", remove(random(), arr))
`,
/* = = = = = = = = */
      "zip-strings": `function alloc(n)
    if n == 0 then return end
    local _, ...r = array(nil)
    for i = 2,n do
        _, ...r = array(r, nil) 
    end
    return r
end

function zip(l ,r)
    assert(#l, #r, "input sources sizes not matched")
    local ...z = alloc(#l)

    for i = 1, #z do
        z[i] = l[i] .. r[i]
    end

    println("zip:", l, "+++", r, "===", z)
    return z
end

local a, b, c = (zip(
    l='a', r='1',
    l='b', r='2',
    l='c', r=1+2,
))
assert(a == "a1" and b == "b2" and c == "c3")

local _, ...kv = array()
for i = 0, 9 do
    _, ...kv = array(kv, 'l', char(unicode('a') + i), 'r', i)
end

zip([nil] = array(kv))`,
/* = = = = = = = = */
      "countdown": `function countdown(n)
    yield true
    while n do
       	n -= 1
        yield n
    end
end

local _, ...state1 = countdown(4)
local _, ...state2 = countdown(5)
initstate1, initstate2 = state1, state2

function run(ss)
    local result
    local ...s = __g(ss)
    result, ...s = resume(countdown, s)
    if s and #s then 
        __g(ss, s) 
        return result, true
    end
    __g(ss, __g("init" .. ss))
    return result, false
end

function runCountdown(s)
    while true do
        local tick, ok = run(s)
        if not ok then break end
        println(tick)
    end
end

println("Run 1st countdown")
runCountdown("state1")
println("Run 2nd countdown")
runCountdown("state2")
println("Re-run 1st countdown")
runCountdown("state1")
println("Re-run 1st countdown again")
runCountdown("state1")`,
/* = = = = = = = = */
      "eof": ""
  };
