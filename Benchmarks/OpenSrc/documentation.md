This directory will contain benchmarks that we use to stress the golang garbage collector.
We test out open source benchmarks and expand to cover a wider array of stress tests as well as different parameter configurations.
The main knob we play with is GOGC, although other factors are at play, chief of which is the threshold at which the GC is triggered.
At present we augment this factor by allocating a certain amount of space before running our benchmarks with the GC enabled.
