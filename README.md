`script` is a simple script engine written in golang with Lua-like syntax.

(If you are looking for a Lua 5.2 compatible engine, refer to tag `v0.2`)

## Differ from Lua

- Functions should be declared in the topmost scope:
	- `do function foo() ... end end` is invalid.
	- `function foo() function bar() ... end end` is invalid.
- Use `lambda` to declare anonymous functions:
	- `local foo = lambda(x) ... end`.
	- Last `return` can be omitted: `lambda (x) x=x+1; x end <=> lambda (x) x=x+1; return x end`.
- Syntax of calling functions strictly requires no spaces between callee and '(':
	- `print(1)` is the only right way of calling a function.
	- `print (1)` literally means two things: 1) get value of `print` and discard it, 2) evaluate `(1)`.
- To write variadic functions:
	- `function foo(a, b...) end`.
	- `args = {1, 2, 3}; foo(args...)`.
- Simple keyword arguments syntax sugar:
	- `foo(a, b=2, c=3)` will be converted to `foo(a, { b=2, c=3 })`.
	- There is no REAL keyword argument in script, e.g.: `function foo(a) return type(a) end`:
		- `foo(1)` returns "number".
		- `foo(a=1) <=> foo({a=1})` returns "table".
- Returning multiple arguments will be translated into returning a `table`, e.g.:
	- `function f() return 1, 2 end; local a, b = f()`.
	- `function f() return {1, 2} end; local tmp = f(); local a, b = tmp[0], tmp[1]`.
	- `local a, b, c = d <=> local a, b, c = d[0], d[1], d[2]`.
- Everything starts at ZERO. For-loops start inclusively and end exclusively, e.g.:
	- `a={1, 2}; assert(a[0] == 1)`.
	- `for i=0,n do ... end`.
	- `for i=n-1,-1,-1 do ... end`.
- Functions loaded from `table` will have a self-like argument at first, to not-to-have it, use `:` operator, e.g.:
	- `a={foo=lambda(this) print(this.v) end, v=1} a.foo()` will print `1`.
	- `a={foo=lambda(this) print(this) end} a:foo()` will print nil.
	- That's to say, you should call most lib functions using `:`.
- You can define up to 32000 variables (varies depending on the number of temporal variables generated by interpreter) in a function.
- Numbers are `int64 + float64` internally, interpreter may promote it to `float64` when needed and downgrade it to `int64` when possible.
- You can `return` anywhere inside functions, `continue` inside for-loops, `goto` any label within the same function.

## Run

```golang
program, err := script.LoadString("return 1")
v, err := program.Run() // v == 1
```

### Global Values

```golang
script.AddGlobalValue("G", func() int { return 1 })

program, _ := script.LoadString("return G() + 1")
v, err := program.Run() // v == 2

program, _ = script.LoadString("return G() + 2")
v, err = program.Run() // v == 3

program, _ = script.LoadString("return G + 2", &CompileOptions{
	GlobalKeyValues: {
		"G": 10, // override the global 'G'
	},
})
v, err = program.Run() // v == 12
```

## Benchmarks

Refer to [here](https://github.com/coyove/potatolang/blob/master/tests/bench/perf.md).

