  var demos = {
      "Select a demo...": `-- Author: coyove
_, author = match(SOURCE_CODE, [[Author: (\\S+)]])
println("Author is:", author)

-- Print all global values, mainly functions
-- use doc(function) to view its documentation
local n, ...g = globals()

print(format("version {}, total global values: {}\\n", VERSION, n))

for i=1,#g do
    local name = str(g[i])
    if #name > 32 then
        name = name[1:16] .. '...' .. name[#name-16:#name]
    end

    title = i .. ": " .. name
    print((",----------------------------------------------------------------")[1:#title+3], '.')
    print("| ", title, " |")
    print(("'----------------------------------------------------------------")[1:#title+3], "'")
    print(type(g[i]) == "function" and doc(g[i]) or "\\tN/A")
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
    local a = "hello"
    ::inner::
    print("I'm in")
    print("a=", a, ", which is not inited")
end
`,
/* = = = = = = = = */
      "yield": `function yieldable(n)
    while n > 0 do
        return n, debug_state()
        n -= 1
    end
end

local c, state = yieldable(10)
while c do
    print(c)
    c, state = debug_resume(yieldable, state)
end
`,
/* = = = = = = = = */
      "time": `println("Unix timestamp:", time())
println("Go time.Time:", Go_time().Format("2006-01-02 15:04:05"))
println(strtime("Y-m-d H:i:s", Go_time()))
println(doc(Go_time))
`,
/* = = = = = = = = */
      "json": `local j = json(dict( a=1, b=2, array=dict( 1, 2, dict( inner="inner" ))))
assert(json_get(j, "a") == 1)
assert(json_get(j, "b") == 2)
local n, a, b, c = json_get(j, "array")
assert(n == 3 and a == 1 and b == 2 and json_get(c, "inner") == "inner")
assert(json_get(j, "array.2.inner")=="inner")
println(json_get(j, "array"))

--[[
json() uses https://github.com/tidwall/gjson
Learn its syntax at https://github.com/tidwall/gjson/blob/master/SYNTAX.md
]]

-- Create a true array from variables
local _, ...arr = array(1, 2, 3)
print(dict(arr))

--[[
dynamic object (dict) creation is a bit harder to write,
first we need to trick the parser with syntax "foo( [nil] = something )" so it knows we are feeding key-value pairs,
then we layout pairs sequentially, so it looks like: { [nil] = key1, value1, key2, value2, ... }
however the above sequence is misaligned: Nth key becomes Nth value and Nth value becomes (N+1)th key
so we have to prepend the sequence with a dummy value to correct the positions:
{ [nil] = dummy, key1, value1, key2, value2, ... }
]]
local n, ...kvpairs = array("key1", "value1", "key2", "value2")
println(dict( [nil] = array(kvpairs, "key3", "value3") ))

`,
/* = = = = = = = = */
      "call": `function veryComplexFunction(a, b, c, d, e, f, g, H, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, Z)
    [[This is a very complex function,
it will only print H and Z's values]]
    println(H, Z)
end

println(doc(veryComplexFunction))
veryComplexFunction(Z="world", ["H"]="hello")

-- This is a trick from 'json' demo
local _, ...args = array("Z", "世界")
veryComplexFunction([nil] = array(args, "H", "你好"))

function foo() return "hello", "world" end
function bar() return "world", "hello" end

local ...a = random() > 0.5 and foo() or bar()
println(a)

function foo(a)
    return 1 + a, 2 + a, 3 + a
end
println(foo(0), foo(3)) -- println(1, 2, 3, 4, 5, 6)
println( ( foo(0) ), foo(3) ) -- println(1, 4, 5, 6)
println( ( foo(0) ), ( foo(3) ) ) -- println(1, 4)
`,
/* = = = = = = = = */
      "http": `local code, headers, body = http(
    method="POST",
    url="http://httpbin.org/post",
    query="k1=v1",
    query="k2=v2",
    json=dict(
        name="Smith",
    ),
)
println("code:", code)
println("headers:", headers)

if body then
    local data = json_get(body, "data")
    println("name:", json_get(data, "name"))
    println("args:", json_get(body, "args"))
end`,

      "goquery": `local code, _, body = http("GET", "https://example.com")
if iserror(code) then
    return "request failed: " .. code.error()
end

local el = goquery(body, "div").nodes()
for i = 1,#el do
    println(el[i].text())
end

local code, _, body = http(url="https://bokete.jp/boke/recent")
if iserror(code) then
    return "request failed: " .. code.error()
end

local list = goquery(body, ".boke-list")
list = list.find(".boke")
local boke = list.nodes()

for i=1,#boke do
    println("#" .. i, trim(boke[i].find('.photo-content img').attr('src'), "//", "prefix"))
    println("  ", trim(boke[i].find('.boke-text').text()))
end
`,
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
            return arr[1:i-1], arr[i+1:#arr]
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
    return true, debug_state()
    while n do
       	n -= 1
        return n, debug_state()
    end
end

local _, ...state1 = countdown(4)
local _, ...state2 = countdown(5)
initstate1, initstate2 = state1, state2

function run(ss)
    local result
    local ...s = __g(ss)
    result, ...s = debug_resume(countdown, s)
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
        write(stdout(), tick, ' ')
    end
    print()
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
      "debug": `function debug_find(name, ...info) 
    --[[
    debug info are laid out as such:
    0, var_name1, var_value1, 1, var_name2, var_value2, 2, var_name3, ...
    index starts from 0 to align with the internal logic
    ]]
    for i=1,#info,3 do
        if info[i+1] == name then return info[i] end
    end
    return -1
end

function foo(a, b)
    local c = a + b -- c == 3

    println(debug_locals())

    local d = a * b -- d == 2

    local d_on_stack = debug_find("d", debug_locals()) -- find d's position on stack

    println(debug_locals())

    debug_set(d_on_stack, 100) -- alter d to 100

    return c / d
end

return foo(1,2)
`,
/* = = = = = = = = */
      "dict": `local m = dict(a=1, b=2)
for k, v in pair(m)() do
    assert(unicode("a") - 1 + v, unicode(k))
end`,
/* = = = = = = = = */
      "bing.com": `local _, _, body = http(url="https://cn.bing.com/HPImageArchive.aspx", queries=dict(format='js', n=10))
local n, ...items = json_get(body, "images")

for i =1,n do
    println("https://cn.bing.com/" .. json_get(items[i], "url"))
end`,
/* = = = = = = = = */
      "iterator": `function range(from, to, step, __cookie)
    step = step or 1
    if __cookie == nil then
        -- init
        return range, from - step, to, step, true
    end
    if step > 0 then
        if from + step <= to then
            return from + step, to, step, true
        end
    else
        if from + step >= to then
          	return from + step, to, step, true
        end
    end
end

-- https://www.lua.org/pil/7.2.html

for i in range(1, 10) do
    write(stdout(), i, " ")
end
println()

for i in range(from=10, to=1, step=-1) do
    write(stdout(), i, " ")
end
println()

function countdown(n)
    while n > 0 do
        return n, debug_state()
        n -= 1
    end
end

function exec(x, f, ...args)
    if x == nil then 
        local ...r = f(args) 
        return r[#r], r[1:#r-1]
    end
    local ...r = debug_resume(f, x)
    return r[#r], r[1:#r-1]
end

for _, i in exec(countdown,10) do
    write(stdout(), i, " ")
end
println()

function range2(x, ...)
    if #... == 0 then return end
    if x == nil then
        return 1, ...[1]
    end
    if #... == x then return end
    return x + 1, ...[x + 1]
end

local _, ...a = array(10, 9, 8, 7, 6, 5, 4, 3, 2, 1)
for i, e in range2(a) do
    write(stdout(), e, " ")
end
println()

local idx = 11
function selfiter()
    idx -= 1
    return idx or nil
end

for e in selfiter do
    write(stdout(), e, " ")
end
println()`,
/* = = = = = = = = */
      "eof": ""
  };
