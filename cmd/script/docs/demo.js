  var demos = {
      "Select a demo...": "print([[ hello world ]])",
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
println("Go time.Time:", Go_time().Format("2006-01-02 15:04:05"))`,
/* = = = = = = = = */
      "json": `local j = { a=1, b=2, array={ 1, 2, { inner="inner" } } }
--[[
There is no table type, code above will actually generate the
correspondent JSON STRING: '{"a":1,"b":2,....}'
]]

assert(json(j, "a") == 1)
assert(json(j, "b") == 2)
local n, a, b, c = json(j, "array")
assert(n == 3 and a == 1 and b == 2 and json(c, "inner")== "inner")
assert(json(j, "array.2.inner")=="inner")
println(json(j, "array"))

--[[
json() uses https://github.com/tidwall/gjson
Learn its syntax at https://github.com/tidwall/gjson/blob/master/SYNTAX.md
]]

`,
/* = = = = = = = = */
      "call": `function veryComplexFunction(a, b, c, d, e, f, g, H, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, Z)
    println(H, Z)
end
veryComplexFunction(Z="world", ["H"]="hello")
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
println("code=", code)
println("headers=", headers)

if body then
    local data = json(body, "data")
    println(json(data, "name"))
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
      "eof": ""
  };
