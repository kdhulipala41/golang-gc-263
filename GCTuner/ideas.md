# Key Ideas

- Transform the garbage.go code to have more interesting memory allocation pattern and longer timeline
  - currently is just a short-lived program parsing out a given file / string to an AST -- then proceeding to throw this AST away as garbage (nice to stress the GC as this contains many pointers to other different objects and nodes of the AST)
  - possibly make a REST API and just load test it with some concurrent requests?

# Next Steps

- Try implementing a dynamic GOGC tuner with the following resource from Uber: https://www.uber.com/blog/how-we-saved-70k-cores-across-30-mission-critical-services/
- Alternatively, try implementing a dynamic GOMEMLIMIT tuner with the following resource from Zomato: https://blog.zomato.com/go-beyond-building-performant-and-reliable-golang-applications

# Resources on the GC itself

- A useful article that explains the concurrent mark-and-sweep and some techniques for profiling the GC at runtime: https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html
