function curry(f, args...)
    assert.shape(f, "@function")
    local ac = f.argcount()
    if f.isvarg() then ac -= 1 end
    if #args >= ac then return f.call(args...) end

    local cf = function(args...)
        self.args.concat(args)
        if #args >= self.args_needed then
            return self.f.call(self.args...)
        end
        local new_cf = self.copy(true)
        new_cf.args_needed = self.args_needed - #args
        return new_cf
    end
    cf.f = f
    cf.args = args or []
    cf.args_needed = ac - #args
    return cf
end

curry(print, "begin currying")

function add(a, b, c)
    return a+b+c
end

assert(curry(add, 1, 2)(3), 6)
assert(curry(add, 1)(2)(3), 6)
assert(curry(add)(1)(2)(3), 6)

function sizeof(c...)
return #c
end

assert(curry(sizeof), 0)
assert(curry(sizeof, 1), 1)

function addv(a, b, c...)
    local res = a + b
    for _, ca in c do
        res += ca
    end
    return res
end

assert(curry(addv, 1, 2), 3)
assert(curry(addv, 1)(2), 3)
assert(curry(addv)(1)(2, 3), 6)

print(curry(addv))
