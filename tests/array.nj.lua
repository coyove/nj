a = [1,2,3,4]
b = a[:2]
b += "a"
assert(a[2], 'a')
print(a)


a = bytes(16)
function(i) return i end.map(a, os.numcpus)

function foo() return a end

b = foo()[:10]
b += 100
assert(a[10], 100)

b = [foo()][(0)][:10]
b += 100
assert(a[10], 100)


syncMap.Store("a", 1)
assert(syncMap.Load("a")[0], 1)


