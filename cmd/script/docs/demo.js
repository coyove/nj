  var demos = {
      "Select a demo...": "print([[ hello world ]])",
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
      "n-body": `local PI = 3.141592653589793
local SOLAR_MASS = 4 * PI * PI
local DAYS_PER_YEAR = 365.24

sun = struct("x", 0.0, "y", 0.0, "z", 0.0, "vx", 0.0, "vy", 0.0, "vz", 0.0, "mass", SOLAR_MASS)

jupiter = struct("x" , 4.84143144246472090e+00,
"y" , -1.16032004402742839e+00,
"z" , -1.03622044471123109e-01,
"vx" , 1.66007664274403694e-03 * DAYS_PER_YEAR,
"vy" , 7.69901118419740425e-03 * DAYS_PER_YEAR,
"vz" , -6.90460016972063023e-05 * DAYS_PER_YEAR,
"mass" , 9.54791938424326609e-04 * SOLAR_MASS)

saturn = struct("x" , 8.34336671824457987e+00,
"y" , 4.12479856412430479e+00,
"z" , -4.03523417114321381e-01,
"vx" , -2.76742510726862411e-03 * DAYS_PER_YEAR,
"vy" , 4.99852801234917238e-03 * DAYS_PER_YEAR,
"vz" , 2.30417297573763929e-05 * DAYS_PER_YEAR,
"mass" , 2.85885980666130812e-04 * SOLAR_MASS)

uranus = struct("x" , 1.28943695621391310e+01,
"y" , -1.51111514016986312e+01,
"z" , -2.23307578892655734e-01,
"vx" , 2.96460137564761618e-03 * DAYS_PER_YEAR,
"vy" , 2.37847173959480950e-03 * DAYS_PER_YEAR,
"vz" , -2.96589568540237556e-05 * DAYS_PER_YEAR,
"mass" , 4.36624404335156298e-05 * SOLAR_MASS)

neptune = struct("x" , 1.53796971148509165e+01,
"y" , -2.59193146099879641e+01,
"z" , 1.79258772950371181e-01,
"vx" , 2.68067772490389322e-03 * DAYS_PER_YEAR,
"vy" , 1.62824170038242295e-03 * DAYS_PER_YEAR,
"vz" , -9.51592254519715870e-05 * DAYS_PER_YEAR,
"mass" , 5.15138902046611451e-05 * SOLAR_MASS)

function array(...) return ... end
local ...bodies = array(sun,jupiter,saturn,uranus,neptune)

local function advance(nbody, dt)
  for i=1,nbody do
    local bi = bodies[i]
    local bix, biy, biz, bimass = bi.x, bi.y, bi.z, bi.mass
    local bivx, bivy, bivz = bi.vx, bi.vy, bi.vz
    for j=i+1,nbody do
      local bj = bodies[j]
      local dx, dy, dz = bix-bj.x, biy-bj.y, biz-bj.z
      local dist2 = dx*dx + dy*dy + dz*dz
      local mag = sqrt(dist2)
      mag = dt / (mag * dist2)
      local bm = bj.mass*mag
      bivx = bivx - (dx * bm)
      bivy = bivy - (dy * bm)
      bivz = bivz - (dz * bm)
      bm = bimass*mag
      bj.vx = bj.vx + (dx * bm)
      bj.vy = bj.vy + (dy * bm)
      bj.vz = bj.vz + (dz * bm)
    end
    bi.vx = bivx
    bi.vy = bivy
    bi.vz = bivz
    bi.x = bix + dt * bivx
    bi.y = biy + dt * bivy
    bi.z = biz + dt * bivz
  end
end

local function energy(nbody)
  local e = 0
  for i=1,nbody do
    local bi = bodies[i]
    local vx, vy, vz, bim = bi.vx, bi.vy, bi.vz, bi.mass
    e = e + (0.5 * bim * (vx*vx + vy*vy + vz*vz))
    for j=i+1,nbody do
      local bj = bodies[j]
      local dx, dy, dz = bi.x-bj.x, bi.y-bj.y, bi.z-bj.z
      local distance = sqrt(dx*dx + dy*dy + dz*dz)
      e = e - ((bim * bj.mass) / distance)
    end
  end
  return e
end

local function offsetMomentum(nbody)
  local px, py, pz = 0, 0, 0
  for i=1,nbody do
    local bi = bodies[i]
    local bim = bi.mass
    px = px + (bi.vx * bim)
    py = py + (bi.vy * bim)
    pz = pz + (bi.vz * bim)
  end
  bodies[1].vx = -px / SOLAR_MASS
  bodies[1].vy = -py / SOLAR_MASS
  bodies[1].vz = -pz / SOLAR_MASS
end

local N = 5000
local nbody = #bodies

offsetMomentum(nbody)
println(energy(nbody))
for i=1,N do advance(nbody, 0.01) end
println(energy(nbody))`,
      "goto": `goto inner
if false then
    ::inner::
    print("I'm in")
end
`,
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
      "time": `println("Unix timestamp:", time())
println("Go time.Time:", Go_time().Format("2006-01-02 15:04:05"))`,

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

      "call": `function veryComplexFunction(a, b, c, d, e, f, g, H, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, Z)
    println(H, Z)
end
veryComplexFunction(Z="world", ["H"]="hello")
`,

      "eof": ""
  };
