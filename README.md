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

BenchmarkIntToAddr-8                      	30000000	        56.1 ns/op	      24 B/op	       2 allocs/op
BenchmarkAddrToInt-8                      	20000000	        80.5 ns/op	      80 B/op	       2 allocs/op
BenchmarkParseNetworkCIDR-8               	10000000	       167 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8             	50000000	        30.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkIPv4ClassfulNetwork-8            	20000000	        89.1 ns/op	      68 B/op	       3 allocs/op
BenchmarkNextAddr-8                       	10000000	       165 ns/op	     104 B/op	       4 allocs/op
BenchmarkShrinkNetwork-8                  	 2000000	       747 ns/op	     104 B/op	      26 allocs/op
BenchmarkNextNetwork-8                    	 5000000	       351 ns/op	     256 B/op	       9 allocs/op
BenchmarkFindInbetweenSubnets-8           	  200000	      7118 ns/op	    4336 B/op	     133 allocs/op
BenchmarkFindInbetweenSubnetsBuffered-8   	  300000	      4255 ns/op	     168 B/op	       9 allocs/op
BenchmarkFindUnusedSubnets-8              	   50000	     31529 ns/op	    6720 B/op	     480 allocs/op
```
