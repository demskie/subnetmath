# subnetmath

```Go
func main() {
    aggregate := subnetmath.BlindlyParseCIDR("192.168.0.0/22")
    subnets := []*net.IPNet{
        subnetmath.BlindlyParseCIDR("192.168.1.0/24"),
        subnetmath.BlindlyParseCIDR("192.168.2.32/30"),
    }
    unused := subnetmath.UnusedSubnets(aggregate, subnets)
    // [ 192.168.0.0/24 192.168.2.0/27 192.168.2.36/30 192.168.2.40/29 192.168.2.48/28 192.168.2.64/26 192.168.2.128/25 192.168.3.0/24 ]
}
```

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
