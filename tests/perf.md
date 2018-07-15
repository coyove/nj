# Performance Tests

This folder contains some performance tests based on [benchmarksgame](https://benchmarksgame-team.pages.debian.net/benchmarksgame/).

Basically speaking, potatolang is very slow, 10x ~ 100x slower than a native Go program. Without JIT it is also impossible to compete with other scritpt languages like lua or javascript.

However sometimes it can outperform poorly-written python/perl code. By which i mean, potatolang is not that slow as you might think, but to achieve the best performance, you have to write ugly code.

