# subnetmath

```Go
func main() {
    aggregate := subnetmath.ParseNetworkCIDR("192.168.0.0/22")
    subnets := []*net.IPNet{
        subnetmath.ParseNetworkCIDR("192.168.1.0/24"),
        subnetmath.ParseNetworkCIDR("192.168.2.32/30"),
    }
    unused := subnetmath.UnusedSubnets(aggregate, subnets)
    // [ 192.168.0.0/24 192.168.2.0/27 192.168.2.36/30 192.168.2.40/29 192.168.2.48/28 192.168.2.64/26 192.168.2.128/25 192.168.3.0/24 ]
}
```

```Bash
Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz

BenchmarkIntToAddr-8              	30000000	        56.8 ns/op	      24 B/op	       2 allocs/op
BenchmarkAddrToInt-8              	20000000	        92.6 ns/op	      80 B/op	       2 allocs/op
BenchmarkParseNetworkCIDR-8       	10000000	       176 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8     	50000000	        25.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkIPv4ClassfulNetwork-8    	20000000	        94.4 ns/op	      68 B/op	       3 allocs/op
BenchmarkNextAddr-8               	10000000	       182 ns/op	     104 B/op	       4 allocs/op
BenchmarkShrinkNetwork-8          	 1000000	      1778 ns/op	    1256 B/op	      50 allocs/op
BenchmarkNextNetwork-8            	 3000000	       501 ns/op	     384 B/op	      12 allocs/op
BenchmarkGetAllAddresses-8        	   10000	    154467 ns/op	  123200 B/op	    3091 allocs/op
BenchmarkFindInbetweenSubnets-8   	  200000	     11712 ns/op	    6960 B/op	     194 allocs/op
BenchmarkFindUnusedSubnets-8      	    3000	    424571 ns/op	  186624 B/op	    7228 allocs/op
```
