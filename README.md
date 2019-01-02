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

BenchmarkConvertV4IntegerToAddress-8   	500000000	         3.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertV4AddressToInteger-8   	2000000000	         0.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseNetworkCIDR-8            	10000000	       171 ns/op	      72 B/op	       4 allocs/op
BenchmarkNetworkComesBefore-8          	50000000	        24.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetClassfulNetwork-8          	20000000	        90.9 ns/op	      68 B/op	       3 allocs/op
BenchmarkNextAddr-8                    	50000000	        26.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkShrinkNetwork-8               	 1000000	      1719 ns/op	    1256 B/op	      50 allocs/op
BenchmarkNextNetwork-8                 	  500000	      2660 ns/op	    2248 B/op	      61 allocs/op
BenchmarkGetAllAddresses-8             	   50000	     25808 ns/op	   29044 B/op	    1035 allocs/op
BenchmarkFindUnusedSubnets-8           	    3000	    436046 ns/op	  242548 B/op	    8748 allocs/op
BenchmarkFindInbetweenV4Subnets-8      	   10000	    147543 ns/op	  121664 B/op	    3267 allocs/op
```
