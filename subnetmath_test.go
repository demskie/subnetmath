package subnetmath

import (
	"net"
	"reflect"
	"testing"
)

func TestUnusedSubnets(t *testing.T) {
	aggregate := BlindlyParseCIDR("192.168.0.0/22")
	subnets := []*net.IPNet{
		BlindlyParseCIDR("192.168.1.0/24"),
		BlindlyParseCIDR("192.168.2.32/30"),
	}
	output := UnusedSubnets(aggregate, subnets...)
	expectedOutput := []*net.IPNet{
		BlindlyParseCIDR("192.168.0.0/24"),
		BlindlyParseCIDR("192.168.2.0/27"),
		BlindlyParseCIDR("192.168.2.36/30"),
		BlindlyParseCIDR("192.168.2.40/29"),
		BlindlyParseCIDR("192.168.2.48/28"),
		BlindlyParseCIDR("192.168.2.64/26"),
		BlindlyParseCIDR("192.168.2.128/25"),
		BlindlyParseCIDR("192.168.3.0/24"),
	}
	if reflect.DeepEqual(output, expectedOutput) == false {
		t.Error("\n",
			"<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expectedOutput,
		)
	}
	aggregate = BlindlyParseCIDR("172.16.0.0/16")
	subnets = []*net.IPNet{
		BlindlyParseCIDR("172.16.11.0/24"),
		BlindlyParseCIDR("172.16.20.0/24"),
		BlindlyParseCIDR("172.16.24.0/24"),
		BlindlyParseCIDR("172.16.25.0/24"),
		BlindlyParseCIDR("172.16.26.0/26"),
		BlindlyParseCIDR("172.16.26.64/26"),
		BlindlyParseCIDR("172.16.26.192/26"),
		BlindlyParseCIDR("172.16.27.0/26"),
		BlindlyParseCIDR("172.16.28.0/24"),
		BlindlyParseCIDR("172.16.29.0/24"),
		BlindlyParseCIDR("172.16.30.0/23"),
		BlindlyParseCIDR("172.16.33.0/24"),
		BlindlyParseCIDR("172.16.40.0/22"),
		BlindlyParseCIDR("172.16.44.0/22"),
		BlindlyParseCIDR("172.16.48.0/22"),
		BlindlyParseCIDR("172.16.52.0/22"),
		BlindlyParseCIDR("172.16.56.0/22"),
		BlindlyParseCIDR("172.16.64.0/22"),
		BlindlyParseCIDR("172.16.255.1/32"),
		BlindlyParseCIDR("172.16.255.2/32"),
		BlindlyParseCIDR("172.16.255.3/32"),
		BlindlyParseCIDR("172.16.255.4/32"),
		BlindlyParseCIDR("172.16.255.5/32"),
		BlindlyParseCIDR("172.16.255.6/32"),
		BlindlyParseCIDR("172.16.255.16/30"),
		BlindlyParseCIDR("172.16.255.20/30"),
		BlindlyParseCIDR("172.16.255.24/30"),
		BlindlyParseCIDR("172.16.255.28/30"),
		BlindlyParseCIDR("172.16.255.32/30"),
		BlindlyParseCIDR("172.16.255.36/30"),
		BlindlyParseCIDR("172.16.255.40/30"),
		BlindlyParseCIDR("172.16.255.44/30"),
		BlindlyParseCIDR("172.16.255.48/30"),
		BlindlyParseCIDR("172.16.255.52/30"),
		BlindlyParseCIDR("172.16.255.56/30"),
	}
	expectedOutput = []*net.IPNet{
		BlindlyParseCIDR("172.16.0.0/21"),
		BlindlyParseCIDR("172.16.8.0/23"),
		BlindlyParseCIDR("172.16.10.0/24"),
		BlindlyParseCIDR("172.16.12.0/22"),
		BlindlyParseCIDR("172.16.16.0/22"),
		BlindlyParseCIDR("172.16.21.0/24"),
		BlindlyParseCIDR("172.16.22.0/23"),
		BlindlyParseCIDR("172.16.26.128/26"),
		BlindlyParseCIDR("172.16.27.64/26"),
		BlindlyParseCIDR("172.16.27.128/25"),
		BlindlyParseCIDR("172.16.32.0/24"),
		BlindlyParseCIDR("172.16.34.0/23"),
		BlindlyParseCIDR("172.16.36.0/22"),
		BlindlyParseCIDR("172.16.60.0/22"),
		BlindlyParseCIDR("172.16.68.0/22"),
		BlindlyParseCIDR("172.16.72.0/21"),
		BlindlyParseCIDR("172.16.80.0/20"),
		BlindlyParseCIDR("172.16.96.0/19"),
		BlindlyParseCIDR("172.16.128.0/18"),
		BlindlyParseCIDR("172.16.192.0/19"),
		BlindlyParseCIDR("172.16.224.0/20"),
		BlindlyParseCIDR("172.16.240.0/21"),
		BlindlyParseCIDR("172.16.248.0/22"),
		BlindlyParseCIDR("172.16.252.0/23"),
		BlindlyParseCIDR("172.16.254.0/24"),
		BlindlyParseCIDR("172.16.255.0/32"),
		BlindlyParseCIDR("172.16.255.7/32"),
		BlindlyParseCIDR("172.16.255.8/29"),
		BlindlyParseCIDR("172.16.255.60/30"),
		BlindlyParseCIDR("172.16.255.64/26"),
		BlindlyParseCIDR("172.16.255.128/25"),
	}
	output = UnusedSubnets(aggregate, subnets...)
	if reflect.DeepEqual(output, expectedOutput) == false {
		t.Error("\n",
			"<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expectedOutput,
		)
	}
}

func BenchmarkNetworkComesBefore(b *testing.B) {
	network := BlindlyParseCIDR("192.168.0.0/28")
	otherNetwork := BlindlyParseCIDR("192.168.0.0/22")
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
	startNetwork := BlindlyParseCIDR("8.0.0.0/8")
	for i := 0; i < b.N; i++ {
		network := DuplicateNetwork(startNetwork)
		for j := 0; j < 24; j++ {
			network = ShrinkNetwork(network)
		}
	}
}

func BenchmarkNextNetwork(b *testing.B) {
	network := BlindlyParseCIDR("192.168.0.0/28")
	for i := 0; i < b.N; i++ {
		network = NextNetwork(network)
	}
}

func BenchmarkGetAllAddresses(b *testing.B) {
	network := BlindlyParseCIDR("192.168.0.0/22")
	for i := 0; i < b.N; i++ {
		GetAllAddresses(network)
	}
}

func BenchmarkFindUnusedSubnets(b *testing.B) {
	aggregate := BlindlyParseCIDR("192.168.0.0/22")
	subnets := []*net.IPNet{
		BlindlyParseCIDR("192.168.1.0/24"),
		BlindlyParseCIDR("192.168.2.32/30"),
	}
	for i := 0; i < b.N; i++ {
		UnusedSubnets(aggregate, subnets...)
	}
}
