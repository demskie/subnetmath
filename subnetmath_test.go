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
