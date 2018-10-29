package subnetmath

import (
	"net"
	"reflect"
	"testing"
)

func TestGetAllAddresses(t *testing.T) {
	input := ParseNetworkCIDR("2607:fb38:10:1::/64")
	output := GetAllAddresses(input)
	if len(output) != 0 {
		t.Error("\n",
			"<<<input>>>\n", "2607:fb38:10:1::/64",
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", "length: 0",
		)
	}
}

func TestFindUnusedSubnets(t *testing.T) {
	aggregate := ParseNetworkCIDR("10.71.8.0/21")
	subnets := []*net.IPNet{}
	output := FindUnusedSubnets(aggregate, subnets...)
	expected := []*net.IPNet{
		ParseNetworkCIDR("10.71.8.0/21"),
	}
	if reflect.DeepEqual(output, expected) == false {
		t.Error(
			"\n<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expected,
		)
	}
	aggregate = ParseNetworkCIDR("192.168.0.0/22")
	subnets = []*net.IPNet{
		ParseNetworkCIDR("192.168.1.0/24"),
		ParseNetworkCIDR("192.168.2.32/30"),
	}
	output = FindUnusedSubnets(aggregate, subnets...)
	expected = []*net.IPNet{
		ParseNetworkCIDR("192.168.0.0/24"),
		ParseNetworkCIDR("192.168.2.0/27"),
		ParseNetworkCIDR("192.168.2.36/30"),
		ParseNetworkCIDR("192.168.2.40/29"),
		ParseNetworkCIDR("192.168.2.48/28"),
		ParseNetworkCIDR("192.168.2.64/26"),
		ParseNetworkCIDR("192.168.2.128/25"),
		ParseNetworkCIDR("192.168.3.0/24"),
	}
	if reflect.DeepEqual(output, expected) == false {
		t.Error(
			"\n<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expected,
		)
	}
	aggregate = ParseNetworkCIDR("172.16.0.0/16")
	subnets = []*net.IPNet{
		ParseNetworkCIDR("172.16.11.0/24"),
		ParseNetworkCIDR("172.16.20.0/24"),
		ParseNetworkCIDR("172.16.24.0/24"),
		ParseNetworkCIDR("172.16.25.0/24"),
		ParseNetworkCIDR("172.16.26.0/26"),
		ParseNetworkCIDR("172.16.26.64/26"),
		ParseNetworkCIDR("172.16.26.192/26"),
		ParseNetworkCIDR("172.16.27.0/26"),
		ParseNetworkCIDR("172.16.28.0/24"),
		ParseNetworkCIDR("172.16.29.0/24"),
		ParseNetworkCIDR("172.16.30.0/23"),
		ParseNetworkCIDR("172.16.33.0/24"),
		ParseNetworkCIDR("172.16.40.0/22"),
		ParseNetworkCIDR("172.16.44.0/22"),
		ParseNetworkCIDR("172.16.48.0/22"),
		ParseNetworkCIDR("172.16.52.0/22"),
		ParseNetworkCIDR("172.16.56.0/22"),
		ParseNetworkCIDR("172.16.64.0/22"),
		ParseNetworkCIDR("172.16.255.1/32"),
		ParseNetworkCIDR("172.16.255.2/32"),
		ParseNetworkCIDR("172.16.255.3/32"),
		ParseNetworkCIDR("172.16.255.4/32"),
		ParseNetworkCIDR("172.16.255.5/32"),
		ParseNetworkCIDR("172.16.255.6/32"),
		ParseNetworkCIDR("172.16.255.16/30"),
		ParseNetworkCIDR("172.16.255.20/30"),
		ParseNetworkCIDR("172.16.255.24/30"),
		ParseNetworkCIDR("172.16.255.28/30"),
		ParseNetworkCIDR("172.16.255.32/30"),
		ParseNetworkCIDR("172.16.255.36/30"),
		ParseNetworkCIDR("172.16.255.40/30"),
		ParseNetworkCIDR("172.16.255.44/30"),
		ParseNetworkCIDR("172.16.255.48/30"),
		ParseNetworkCIDR("172.16.255.52/30"),
		ParseNetworkCIDR("172.16.255.56/30"),
	}
	expected = []*net.IPNet{
		ParseNetworkCIDR("172.16.0.0/21"),
		ParseNetworkCIDR("172.16.8.0/23"),
		ParseNetworkCIDR("172.16.10.0/24"),
		ParseNetworkCIDR("172.16.12.0/22"),
		ParseNetworkCIDR("172.16.16.0/22"),
		ParseNetworkCIDR("172.16.21.0/24"),
		ParseNetworkCIDR("172.16.22.0/23"),
		ParseNetworkCIDR("172.16.26.128/26"),
		ParseNetworkCIDR("172.16.27.64/26"),
		ParseNetworkCIDR("172.16.27.128/25"),
		ParseNetworkCIDR("172.16.32.0/24"),
		ParseNetworkCIDR("172.16.34.0/23"),
		ParseNetworkCIDR("172.16.36.0/22"),
		ParseNetworkCIDR("172.16.60.0/22"),
		ParseNetworkCIDR("172.16.68.0/22"),
		ParseNetworkCIDR("172.16.72.0/21"),
		ParseNetworkCIDR("172.16.80.0/20"),
		ParseNetworkCIDR("172.16.96.0/19"),
		ParseNetworkCIDR("172.16.128.0/18"),
		ParseNetworkCIDR("172.16.192.0/19"),
		ParseNetworkCIDR("172.16.224.0/20"),
		ParseNetworkCIDR("172.16.240.0/21"),
		ParseNetworkCIDR("172.16.248.0/22"),
		ParseNetworkCIDR("172.16.252.0/23"),
		ParseNetworkCIDR("172.16.254.0/24"),
		ParseNetworkCIDR("172.16.255.0/32"),
		ParseNetworkCIDR("172.16.255.7/32"),
		ParseNetworkCIDR("172.16.255.8/29"),
		ParseNetworkCIDR("172.16.255.60/30"),
		ParseNetworkCIDR("172.16.255.64/26"),
		ParseNetworkCIDR("172.16.255.128/25"),
	}
	output = FindUnusedSubnets(aggregate, subnets...)
	if reflect.DeepEqual(output, expected) == false {
		t.Error(
			"\n<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expected,
		)
	}
}

func BenchmarkNetworkComesBefore(b *testing.B) {
	network := ParseNetworkCIDR("192.168.0.0/28")
	otherNetwork := ParseNetworkCIDR("192.168.0.0/22")
	for i := 0; i < b.N; i++ {
		NetworkComesBefore(network, otherNetwork)
	}
}

func BenchmarkGetClassfulNetwork(b *testing.B) {
	addr := net.ParseIP("192.168.0.0")
	for i := 0; i < b.N; i++ {
		GetClassfulNetwork(addr)
	}
}

func BenchmarkAddToAddr(b *testing.B) {
	addr := net.ParseIP("192.168.0.0")
	for i := 0; i < b.N; i++ {
		addr = AddToAddr(addr, 1)
	}
}

func BenchmarkNextAddr(b *testing.B) {
	addr := net.ParseIP("192.168.0.0")
	for i := 0; i < b.N; i++ {
		addr = NextAddr(addr)
	}
}

func BenchmarkShrinkNetwork(b *testing.B) {
	startNetwork := ParseNetworkCIDR("8.0.0.0/8")
	for i := 0; i < b.N; i++ {
		network := DuplicateNetwork(startNetwork)
		for j := 0; j < 24; j++ {
			network = ShrinkNetwork(network)
		}
	}
}

func BenchmarkNextNetwork(b *testing.B) {
	network := ParseNetworkCIDR("192.168.0.0/28")
	for i := 0; i < b.N; i++ {
		network = NextNetwork(network)
	}
}

func BenchmarkGetAllAddresses(b *testing.B) {
	network := ParseNetworkCIDR("192.168.0.0/22")
	for i := 0; i < b.N; i++ {
		GetAllAddresses(network)
	}
}

func BenchmarkFindUnusedSubnets(b *testing.B) {
	aggregate := ParseNetworkCIDR("192.168.0.0/22")
	subnets := []*net.IPNet{
		ParseNetworkCIDR("192.168.1.0/24"),
		ParseNetworkCIDR("192.168.2.32/30"),
	}
	for i := 0; i < b.N; i++ {
		FindUnusedSubnets(aggregate, subnets...)
	}
}
