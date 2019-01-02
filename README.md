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

BenchmarkConvertV4IntegerToAddress-8   	500000000	         3.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertV4AddressToInteger-8   	200000000	         9.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseNetworkCIDR-8            	10000000	       173 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8          	50000000	        28.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetClassfulNetwork-8          	20000000	        91.0 ns/op	      68 B/op	       3 allocs/op
BenchmarkNextAddr-8                    	50000000	        26.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkShrinkNetwork-8               	 1000000	      1744 ns/op	    1256 B/op	      50 allocs/op
BenchmarkNextNetwork-8                 	 3000000	       417 ns/op	     372 B/op	      11 allocs/op
BenchmarkGetAllAddresses-8             	   50000	     25530 ns/op	   29044 B/op	    1035 allocs/op
BenchmarkFindUnusedSubnets-8           	    5000	    365557 ns/op	  162620 B/op	    6440 allocs/op
BenchmarkFindInbetweenV4Subnets-8      	  200000	     10246 ns/op	    6880 B/op	     193 allocs/op
```
