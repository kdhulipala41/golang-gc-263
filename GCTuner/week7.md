# Key Ideas

- Transform the garbage.go code to have more interesting memory allocation pattern and longer timeline
  - currently is just a short-lived program parsing out a given file / string to an AST -- then proceeding to throw this AST away as garbage (nice to stress the GC as this contains many pointers to other different objects and nodes of the AST)
  - possibly make a REST API and just load test it with some concurrent requests?

# Next Steps

- Try implementing a dynamic GOGC tuner with the following resource from Uber: https://www.uber.com/blog/how-we-saved-70k-cores-across-30-mission-critical-services/
- Alternatively, try implementing a dynamic GOMEMLIMIT tuner with the following resource from Zomato: https://blog.zomato.com/go-beyond-building-performant-and-reliable-golang-applications

# Resources on the GC itself

- A useful article that explains the concurrent mark-and-sweep and some techniques for profiling the GC at runtime: https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html

# Some interesting benchmark results

- By using pprof, built-in testing benchmarks created by the Go developers I was able to get these benchmarks:
- The benchmark esbuild is actually adapted from the open-source tool / module bundler called ESBuild that was written in Go. By comparing the default GC settings, GOGC=100 (when the live heap doubles GC procs) vs. having the GC run at GOGC=200 (GC procs when live heap quadruples):

```
go tool pprof -normalize -diff_base=results/esbuild/defaultConf.debug/ESBuildReactAdminJS-2433549412-cpu.prof results/esbuild/highGoMemLimit.debug/ESBuildReactAdminJS-1689689172-cpu.prof
File: esbuild
Type: cpu
Time: Feb 20, 2025 at 3:43pm (PST)
Duration: 60.40s, Total samples = 129.82s (214.94%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for -2.80s, 2.16% of 129.82s total
Dropped 1 node (cum <= 0.65s)
Showing top 10 nodes out of 1063
      flat  flat%   sum%        cum   cum%
    -9.76s  7.52%  7.52%     -9.76s  7.52%  runtime.madvise
     8.20s  6.32%  1.20%      8.20s  6.32%  syscall.syscall
    -2.84s  2.19%  3.40%     -7.78s  5.99%  runtime.scanobject
     2.21s  1.70%  1.69%      2.21s  1.70%  runtime.usleep
    -2.11s  1.63%  3.32%     -2.61s  2.01%  runtime.greyobject
     1.27s  0.98%  2.34%      1.27s  0.98%  runtime.pthread_cond_wait
     1.25s  0.97%  1.38%      1.25s  0.97%  runtime.memclrNoHeapPointers
    -0.86s  0.66%  2.04%     -0.86s  0.66%  runtime.pthread_kill
    -0.82s  0.63%  2.67%     -1.43s  1.10%  runtime.findObject
     0.66s  0.51%  2.16%      0.66s  0.51%  runtime.pthread_cond_signal
```

- Although it may not seem extremely significant, even just toying with the GOGC default value can make a difference of ~9% in CPU usage during the program execution
- Will keep working to develop an actual dynamic adjustment of this value based off the current CPU usage and live memory value as next steps
