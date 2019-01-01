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

BenchmarkConvertIntegerIPv4-8   	20000000	        60.6 ns/op	      48 B/op	       2 allocs/op
BenchmarkParseNetworkCIDR-8     	10000000	       178.1 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8   	50000000	        31.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetClassfulNetwork-8   	20000000	        93.8 ns/op	      68 B/op	       3 allocs/op
BenchmarkAddToAddr-8            	300000000	         4.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkNextAddr-8             	50000000	        26.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkShrinkNetwork-8        	 1000000	      1745 ns/op	    1256 B/op	      50 allocs/op
BenchmarkNextNetwork-8          	  500000	      2670 ns/op	    2248 B/op	      61 allocs/op
BenchmarkGetAllAddresses-8      	   50000	     25531 ns/op	   29044 B/op	    1035 allocs/op
BenchmarkFindUnusedSubnets-8    	    3000	    467050 ns/op	  242548 B/op	    8748 allocs/op
```
