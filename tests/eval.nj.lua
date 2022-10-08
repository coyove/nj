local env = createprototype("env", function(p)
    this.store = {}
    this.funcs = {}
    this.top = if(p, p.top, this)
    this.parent = p
end)

env._break = []
env._continue = []
env._return = {}

function env.set(k, v)
    this.store[k] = v
end

function env.setp(k, v)
    if this.store.hasownproperty(k) then
        this.store[k] = v
        return true
    end
    if this.parent then
        if this.parent.setp(k, v) then
            return false
        end
    end
    this.store[k] = v
    return true
end

function env.get(k)
    if k == "nil" then return nil, true end
    if k == "true" then return true, true end
    if k == "false" then return false, true end
    local v = this.store[k:self]
    if v != self then
        return v, true
    end
    if this.parent then
        return this.parent.get(k)
    end
    if this.top == this then
        local v = this.funcs[k:self]
        if v != self then
            return v, true
        end
    end
    return nil, false
end

function runUnary(node, e)
    if node.Op == eval.op.ret then
        panic(new(env._return, {result=run(node.A, e)}))
    end
end

function runBinary(node, e)
    if node.Op == eval.op.add then
        return run(node.A, e) + run(node.B, e)
    elseif node.Op == eval.op.sub then
        return run(node.A, e) - run(node.B, e)
    elseif node.Op == eval.op.less then
        return run(node.A, e) < run(node.B, e)
    elseif node.Op == eval.op.inc then
        local res = run(node.A, e) + run(node.B, e)
        e.setp(node.A.Name, res)
        return res
    end
end

runners = {}

function run(node, e)
    local run = runners[node.typename()]
    if run then return run(node, e) end
    local tmp = buffer()
    node.Dump(tmp)
    panic("unknown node: %v".format(tmp.value()))
end

runners["*parser.Prog"] = function(node, e)
    local e2 = if(node.DoBlock, env(e), e)
    local last
    for _, stat in node.Stats do
        last = run(stat, e2) 
        if last == env._break or last == env._continue then return last end
    end
    return last
end
runners["*parser.If"] = function(node, e)
    if run(node.Cond, e) then
        res = run(node.True, env(e))
    elseif node.False then
        res = run(node.False, env(e))
    end
    return res
end
runners["*parser.Loop"] = function(node, e)
    while true do
        local res = run(node.Body, e)
        if res == env._break then break end
        if res == env._continue then 
            run(node.Continue, e)
        end
    end
end
runners["*parser.And"] = function(node, e)
    return run(node.A, e) and run(node.B, e)
end
runners["*parser.Or"] = function(node, e)
    return run(node.A, e) or run(node.B, e)
end
runners["*parser.Declare"] = function(node, e)
    e.set(node.Name.Name, run(node.Value, e))
end
runners["*parser.Assign"] = function(node, e)
    e.setp(node.Name.Name, run(node.Value, e))
end
runners["*parser.Release"] = function(node, e)
    -- omit
end
runners["*parser.LoadConst"] = function(node, e)
    for f in node.natptrat("Funcs") do e.top.funcs[f] = nil end
end
runners["parser.Primitive"] = function(node, e)
    return node.Value()
end
runners["*parser.Binary"] = function(node, e)
    return runBinary(node, e)
end
runners["*parser.Unary"] = function(node, e)
    return runUnary(node, e)
end
runners["*parser.BreakContinue"] = function(node, e)
    return if(node.Break, env._break, env._continue)
end
runners["*parser.Symbol"] = function(node, e)
    local v, ok = e.get(node.Name)
    assert(ok, true, "unknown symbol: %v".format(node.Name))
    return v
end
runners["parser.ExprList"] = function(node, e)
    lst = []
    for _, n in node do
        lst[#lst] = run(n, e)
    end
    return lst
end
runners["parser.ExprAssignList"] = function(node, e)
    lst = {}
    for _, n in node do
        lst[run(n[0], e)] = run(n[1], e)
    end
    return lst
end
runners["*parser.Call"] = function(node, e)
    local callee = run(node.Callee, e)
    if callee is callable then
        return callee(run(node.Args, e)...)
    end
    local e2 = env(e.top)
    for i, name in callee.Args do
        e2.set(name.Name, if(i < #node.Args, run(node.Args[i], e), nil))
    end
    local res = run.try(callee.Body, e2)
    if res is error and res.error() is env._return then
        return res.error().result
    end
    return res
end
runners["*parser.Function"] = function(node, e)
    e.top.funcs[node.Name] = node
    e.store[node.Name] = node
    return node
end

function newenv(a)
    local e = env(nil)
    e.store.merge(globals())
    e.store.merge(a)
    return e
end

function loadstring(expr, e)
    res = run.try(eval.parse(expr), e)
    if res is error and res.error() is env._return then
        return res.error().result
    end
    return res
end

local N=1e3
local res = loadstring("a=0 for i=1,%v+1 do a+=i end a".format(N), newenv(nil))
assert(res, N*(N+1)/2)

start = time()
code = [[
    function fib(n)
        if n < 2 then return n end
        return fib(n - 1) + fib(n-2)
    end
    return (fib(N))
]]
res = loadstring(code, newenv({N=20}))
assert(res, eval(code, {N=20}))
print(time() - start, ' res=', res)
