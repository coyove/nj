assert( "abc" == "abc")
-- assert( "abc" == "\97\98\99")
assert( "abc" == "\x61\x62\x63")
local abc = "abc" 
assert( "abc吱吱吱" == abc + "吱吱" + "吱")
assert( "abc吱吱吱" == "abc\u5431\u5431\u5431")


local a, b = re([["(\d+)"]]).findall("a\"12\" \"2\"3")
assert(a[0] == [["12"]] and a[1] == "12" and b[0] == [["2"]] and b[1] == "2")

local err = error("err")
print("reflectLoad: ", err.error())
assert(err.error() == "err")
assert("eerrrr" == re("(e|r)").replace(err.error(), "$1$1"))

assert("x" == ("嘻x嘻").trim("嘻"))
assert("嘻x" == ("嘻x嘻").trimright("嘻"))
assert("x嘻" == ("嘻x嘻").trimleft("嘻"))

-- json generator
do
	function str.json_get(p) return json.parse(this)[p] end
	local j = json.stringify({
	      a = 'A',
	      "b-2": true,
	})
    assert(j.json_get('a'), 'A')
    assert(j.json_get('b-2'), true)
end

function str.reverse()
	x = ""
    tmp = this
    while #tmp > 0 do
        local r, sz = tmp.decodeutf8()
        x = chr(r) + x
        tmp = tmp[sz:]
    end
	return x
end
assert(("abc")==(("cba").reverse()))
assert(("あbc")==(("cbあ").reverse()))

function foo(x)
	local a= x.split(",").copy(1, 2,["x","y","z"]) 
	assert(a.istyped())
	return a
end

function array.equals(rhs)
	if #rhs != #this then return false end
	for i, v in rhs do
		if v != this[i] then return false end
	end
	return true
end

assert(foo("1,2,3").equals(["1","x","3"]))

assert(json.parse("null") is nil)
assert(json.parse("true") is true)
assert(json.parse("false") is bool)
assert(json.parse.try("nil") is error)
