# Performance Tests

This folder contains some performance scripts based on [benchmarksgame](https://benchmarksgame-team.pages.debian.net/benchmarksgame/). All tests are done on an M1 Macbook Pro.

|       | nj | python | lua5.3 | perl5 |
| ----- | ----- | ----- | ----- | ----- |
|fib (n=35) | 0.87 | 1.77 | 0.71 | 3.74 |
|nbody (iteration=500000) | 1.89 | TODO | 1.83 | 0.59 |
|binarytree (depth=21)    | 2.12 | TODO | 1.28 | 2.20 |
|spectralnorms (size=550, cores=1) | 0.85 | TODO | 0.40 | 0.56 |

