# Performance Tests

This folder contains some performance tests based on [benchmarksgame](https://benchmarksgame-team.pages.debian.net/benchmarksgame/).

Here are some single-threaded results by the 'real' output of `time ./prog`:

|       | potatolang | tengo | lua5.3 | perl5 |
| ----- | ---------- | ----- | ------ | ----- |
|fib (n=35) | 1.5 | 2.9 | 1.6 | 3.4 |
|nbody (iteration=500000) | 8.1 | 9.1 | 1.8 | 4.3 |
|binarytree (depth=23)    | 10.0 | 23.4 | 13.4 | 19.1 |
|spectralnorms (size=2000, cores=8) | 5.0 | TODO | TODO | 2.9 |
|spectralnorms (size=2000, cores=1) | 24.9 | 51.1 | 14.8 | 11.5 |

Basically speaking, potatolang is very slow, 10x ~ 50x slower than a native Go program, 2x slower (sometimes faster) than normal interpreted languages. Without JIT it is also impossible to compete with other implementations like luajit or v8.

However it is relatively faster than other script languages written in Go.

All tests are done on a 2018 15-inch Macbook Pro.
