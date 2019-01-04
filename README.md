# subnetmath

```Go
func main() {
    aggregate := subnetmath.ParseNetworkCIDR("192.168.0.0/22")
    subnets := []*net.IPNet{
        subnetmath.ParseNetworkCIDR("192.168.1.0/24"),
        subnetmath.ParseNetworkCIDR("192.168.2.32/30"),
    }
    unused := subnetmath.FindUnusedSubnets(aggregate, subnets)
    // [ 192.168.0.0/24 192.168.2.0/27 192.168.2.36/30 192.168.2.40/29 192.168.2.48/28 192.168.2.64/26 192.168.2.128/25 192.168.3.0/24 ]
}
```

```Bash
Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz

BenchmarkIntToAddr-8              	20000000	        57.8 ns/op	      24 B/op	       2 allocs/op
BenchmarkAddrToInt-8              	20000000	        86.8 ns/op	      80 B/op	       2 allocs/op
BenchmarkParseNetworkCIDR-8       	10000000	       180 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8     	50000000	        25.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkIPv4ClassfulNetwork-8    	20000000	        92.9 ns/op	      68 B/op	       3 allocs/op
BenchmarkNextAddr-8               	10000000	       175 ns/op	     104 B/op	       4 allocs/op
BenchmarkShrinkNetwork-8          	 2000000	       743 ns/op	     104 B/op	      26 allocs/op
BenchmarkNextNetwork-8            	 3000000	       494 ns/op	     384 B/op	      12 allocs/op
BenchmarkFindInbetweenSubnets-8   	  200000	     11676 ns/op	    6960 B/op	     194 allocs/op
BenchmarkFindUnusedSubnets-8      	   50000	     34091 ns/op	   10304 B/op	     564 allocs/op
```
