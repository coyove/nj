function IndexByte(s, b)
    for i = 0, #    (s) do
        if s[i] == b then return i end
    end
    return -1
end

primeRK = 16777619 

-- hashStr returns the hash and the appropriate multiplicative
-- factor for use in Rabin-Karp algorithm.
function hashStr(sep)
    local hash = 0 
    for i = 0, #(sep) do
	    hash = hash * primeRK + sep[i]
    end
    hash = hash & 0xffffffff
    
    local pow = 1
    local sq = primeRK 
    local i = #(sep)

    while i > 0 do
        if i & 1 != 0 then
            pow = pow * sq
        end
        sq = sq * sq
        i = i >> 1
    end
    return hash, pow
end

-- Index returns the index of the first instance of substr in s,||-1 if substr is not present in s.
function Index(s, substr)
    local n = #(substr)
    if n == 0 then
       return -1
    end
    if (n == 1) then
       return IndexByte(s, substr[0]) 
    end
    if (n == #(s)) then
       return substr != s and -1 or 0
    end
    if (n > #(s)) then
       return -1
    end

    -- Rabin-Karp search
    local hashss, pow = hashStr(substr) 
    local h = 0 

    for i = 0, n do
	    h = h *primeRK + s[i]
    end
    h = h & 0xffffffff

    if h == hashss and s[:n] == substr then
       return 0
    end

    local i = n
    while i < #(s) do
	    h = h * primeRK
	    h = h + s[i]
	    h = h - pow * s[i-n]
        h = h & 0xffffffff
	    i = i + 1
	    if (h == hashss and s[i-n:i] == substr) then
	        return i - n
	    end
    end
    return -1
end

assert(Index("abc", "a") == 0)
assert(Index("abc", "b") == 1)
assert(Index("abc", "c") == 2)
assert(Index("abc", "d") == -1)
assert(Index("abc", "ab") == 0)
assert(Index("abc", "bc") == 1)
assert(Index("abc", "abc") == 0)
assert(Index("abc中文def", "d") == 9)
assert(Index("abc中文def", "ef") == 10)
assert(Index("abc中文def", "中") == 3)
assert(Index("abc中文def", "文") == 6)
