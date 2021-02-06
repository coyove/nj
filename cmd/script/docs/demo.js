  var demos = {
      "Select a demo...": `-- Author: coyove
{_, author} = match(SOURCE_CODE, [[Author: (\\S+)]])
println("Author is:", author)

-- Print all global values, mainly functions
-- use doc(function) to view its documentation
local _g = set(globals())
local g = debug_globals()

print(format("version {}, total global values: {}\\n", VERSION, #g/3))

for i=3,#g,3 do
    if _g.exists(g[i]) then
        local name = str(g[i])
        if #name > 32 then
            name = name[1:16] .. '...' .. name[#name-16:#name]
        end

        title = i/3 .. ": " .. name
        print((",----------------------------------------------------------------")[1:#title+3], '.')
        print("| ", title, " |")
        print(("'----------------------------------------------------------------")[1:#title+3], "'")
        print(type(g[i]) == "function" and doc(g[i]) or 'constant: ' .. g[i-1])
        print()
    end
end`,
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
      "json": `local j = json(dict( a=1, b=2, array={ 1, 2, dict( inner="inner" )}))
assert(json_get(j, "a") == 1)
assert(json_get(j, "b") == 2)
local {a, b, c} = json_get(j, "array")
assert(a == 1 and b == 2 and json_get(c, "inner") == "inner")
assert(json_get(j, "array.2.inner")=="inner")
println(json_get(j, "array"))

--[[
json() uses https://github.com/tidwall/gjson
Learn its syntax at https://github.com/tidwall/gjson/blob/master/SYNTAX.md
]]

println(dict({"key1", "value1", "key2", "value2"}))
`,
/* = = = = = = = = */
      "call": `function veryComplexFunction(a, b, c, d, e, f, g, H, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, Z)
    [[This is a very complex function,
it will only print H and Z's values]]
    println(H, Z)
end

println(doc(veryComplexFunction))
veryComplexFunction(Z="world", ["H"]="hello")
`,
/* = = = = = = = = */
      "http": `local {code, headers, body} = http(
    method="POST",
    url="http://httpbin.org/post",
    query={"k1=v1", "k2=v2",},
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

      "goquery": `local {code, _, body} = http("GET", "https://example.com")
if iserror(code) then
    return "request failed: " .. code.error()
end

local el = goquery(body, "div").nodes()
for i = 1,#el do
    println(el[i].text())
end

local {code, _, body} = http(url="https://bokete.jp/boke/recent")
if iserror(code) then
    return "request failed: " .. code.error()
end

local list = goquery(body, ".boke-list")
list = list.find(".boke")
local boke = list.nodes()

for i=1,#boke do
    println("#" .. i, trim(boke[i].find('.photo-content img').attr('src'), "//", "prefix"))
    println("  ", trim(boke[i].find('.boke-text').text()))
end`,
/* = = = = = = = = */
      "debug": `function debug_find(name, info) 
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
      "bing.com": `local {_, _, body} = http(url="https://cn.bing.com/HPImageArchive.aspx", queries=dict(format='js', n=10))
local items = json_get(body, "images")

for i =1,#items do
    println("https://cn.bing.com/" .. json_get(items[i], "url"))
end`,
/* = = = = = = = = */
      "eof": ""
  };
