# subnetmath

```Go
func main() {
    
}
```

Benchmark results
```
demskie$ go test github.com/demskie/subnetmath -bench=.
goos: darwin
goarch: amd64
pkg: github.com/demskie/subnetmath
BenchmarkNetworkComesBefore-8   	50000000	        26.9 ns/op
BenchmarkGetClassfulNetwork-8   	20000000	        87.6 ns/op
BenchmarkAddToAddr-8            	300000000	         5.11 ns/op
BenchmarkNextAddr-8             	50000000	        29.3 ns/op
BenchmarkShrinkNetwork-8        	 1000000	      1691 ns/op
BenchmarkNextNetwork-8          	 5000000	       286 ns/op
BenchmarkGetAllAddresses-8      	   30000	     44225 ns/op
BenchmarkFindUnusedSubnets-8    	    5000	    330295 ns/op
```
