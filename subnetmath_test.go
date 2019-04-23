package subnetmath

import (
	"math/big"
	"net"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// BUG: why does reflect.DeepEqual not return expected results
func sliceOfSubnetsAreEqual(alpha, bravo []*net.IPNet) bool {
	if len(alpha) != len(bravo) {
		return false
	}
	for i := range alpha {
		if !NetworksAreIdentical(alpha[i], bravo[i]) {
			return false
		}
	}
	return true
}

func TestNetworksAreIdentical(t *testing.T) {
	if !NetworksAreIdentical(
		ParseNetworkCIDR("192.168.0.0/23"),
		ParseNetworkCIDR("192.168.0.0/23"),
	) {
		t.Error("\n",
			ParseNetworkCIDR("192.168.0.0/23"),
			"is identical to",
			ParseNetworkCIDR("192.168.0.0/23"),
		)
	}
}

func TestNetworkComesBefore(t *testing.T) {
	if !NetworkComesBefore(
		ParseNetworkCIDR("192.168.0.0/23"),
		ParseNetworkCIDR("192.168.8.0/23"),
	) {
		t.Error("\n",
			ParseNetworkCIDR("192.168.0.0/23"),
			"didn't come before",
			ParseNetworkCIDR("192.168.8.0/23"),
		)
	}
	if !NetworkComesBefore(
		ParseNetworkCIDR("192.168.0.0/23"),
		ParseNetworkCIDR("fe80::/64"),
	) {
		t.Error("\n",
			ParseNetworkCIDR("192.168.0.0/23"),
			"didn't come before",
			ParseNetworkCIDR("fe80::/64"),
		)
	}
}

func TestSubnetZeroAddr(t *testing.T) {
	input1, input2, _ := net.ParseCIDR("192.168.0.5/23")
	actualOutput := SubnetZeroAddr(input1, input2)
	expectedOutput := net.ParseIP("192.168.0.0")
	if !actualOutput.Equal(expectedOutput) {
		t.Error("\n",
			"<<<input>>>\n", input1, input2,
			"\n<<<actual_output>>>\n", actualOutput,
			"\n<<<expected_output>>>\n", expectedOutput,
		)
	}
}

func TestNextNetwork(t *testing.T) {
	input := ParseNetworkCIDR("192.168.0.0/23")
	actualOutput := NextNetwork(input)
	expectedOutput := ParseNetworkCIDR("192.168.2.0/23")
	if !NetworksAreIdentical(actualOutput, expectedOutput) {
		t.Error("\n",
			"<<<input>>>\n", input,
			"\n<<<actual_output>>>\n", actualOutput,
			"\n<<<expected_output>>>\n", expectedOutput,
		)
	}
}

func TestBroadcastAddr(t *testing.T) {
	input := ParseNetworkCIDR("192.168.0.0/23")
	actualOutput := BroadcastAddr(input)
	expectedOutput := net.ParseIP("192.168.1.255")
	if !actualOutput.Equal(expectedOutput) {
		t.Error("\n",
			"<<<input>>>\n", input,
			"\n<<<actual_output>>>\n", actualOutput,
			"\n<<<expected_output>>>\n", expectedOutput,
		)
	}
}

func TestIntToAddr(t *testing.T) {
	actualAddress := IntToAddr(big.NewInt(3232235778))
	expectedAddress := net.ParseIP("192.168.1.2")
	if !actualAddress.Equal(expectedAddress) {
		t.Error("\n",
			"<<<input>>>\n", "3232235778",
			"\n<<<actual_output>>>\n", actualAddress,
			"\n<<<expected_output>>>\n", expectedAddress,
		)
	}
}

func TestAddrToInt(t *testing.T) {
	actualInteger := AddrToInt(net.ParseIP("192.168.1.2"))
	expectedInteger := big.NewInt(3232235778)
	if actualInteger.Cmp(expectedInteger) != 0 {
		t.Error("\n",
			"<<<input>>>\n", "3232235778",
			"\n<<<actual_output>>>\n", actualInteger,
			"\n<<<expected_output>>>\n", expectedInteger,
		)
	}
}

func TestFindInbetweenSubnets(t *testing.T) {
	input := []net.IP{
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.2.2"),
	}
	output := FindInbetweenSubnets(input[0], input[1])
	expected := []*net.IPNet{
		ParseNetworkCIDR("192.168.1.2/31"),
		ParseNetworkCIDR("192.168.1.4/30"),
		ParseNetworkCIDR("192.168.1.8/29"),
		ParseNetworkCIDR("192.168.1.16/28"),
		ParseNetworkCIDR("192.168.1.32/27"),
		ParseNetworkCIDR("192.168.1.64/26"),
		ParseNetworkCIDR("192.168.1.128/25"),
		ParseNetworkCIDR("192.168.2.0/31"),
		ParseNetworkCIDR("192.168.2.2/32"),
	}
	if !sliceOfSubnetsAreEqual(output, expected) {
		t.Error("\n",
			"<<<input>>>\n", input,
			"\n<<<actual_output>>>\n", spew.Sdump(output),
			"\n<<<expected_output>>>\n", spew.Sdump(expected),
		)
	}
}

// 2001:400:: 2001:440:ffff:ffff:7fff:ffff:ffff:ffff

func TestFindInbetweenSubnetsV6(t *testing.T) {
	input := []net.IP{
		net.ParseIP("2001:400::"),
		net.ParseIP("2001:440:ffff:ffff:7fff:ffff:ffff:ffff"),
	}
	output := FindInbetweenSubnets(input[0], input[1])
	expected := []*net.IPNet{
		ParseNetworkCIDR("2001:400::/26"),
		ParseNetworkCIDR("2001:440::/33"),
		ParseNetworkCIDR("2001:440:8000::/34"),
		ParseNetworkCIDR("2001:440:c000::/35"),
		ParseNetworkCIDR("2001:440:e000::/36"),
		ParseNetworkCIDR("2001:440:f000::/37"),
		ParseNetworkCIDR("2001:440:f800::/38"),
		ParseNetworkCIDR("2001:440:fc00::/39"),
		ParseNetworkCIDR("2001:440:fe00::/40"),
		ParseNetworkCIDR("2001:440:ff00::/41"),
		ParseNetworkCIDR("2001:440:ff80::/42"),
		ParseNetworkCIDR("2001:440:ffc0::/43"),
		ParseNetworkCIDR("2001:440:ffe0::/44"),
		ParseNetworkCIDR("2001:440:fff0::/45"),
		ParseNetworkCIDR("2001:440:fff8::/46"),
		ParseNetworkCIDR("2001:440:fffc::/47"),
		ParseNetworkCIDR("2001:440:fffe::/48"),
		ParseNetworkCIDR("2001:440:ffff::/49"),
		ParseNetworkCIDR("2001:440:ffff:8000::/50"),
		ParseNetworkCIDR("2001:440:ffff:c000::/51"),
		ParseNetworkCIDR("2001:440:ffff:e000::/52"),
		ParseNetworkCIDR("2001:440:ffff:f000::/53"),
		ParseNetworkCIDR("2001:440:ffff:f800::/54"),
		ParseNetworkCIDR("2001:440:ffff:fc00::/55"),
		ParseNetworkCIDR("2001:440:ffff:fe00::/56"),
		ParseNetworkCIDR("2001:440:ffff:ff00::/57"),
		ParseNetworkCIDR("2001:440:ffff:ff80::/58"),
		ParseNetworkCIDR("2001:440:ffff:ffc0::/59"),
		ParseNetworkCIDR("2001:440:ffff:ffe0::/60"),
		ParseNetworkCIDR("2001:440:ffff:fff0::/61"),
		ParseNetworkCIDR("2001:440:ffff:fff8::/62"),
		ParseNetworkCIDR("2001:440:ffff:fffc::/63"),
		ParseNetworkCIDR("2001:440:ffff:fffe::/64"),
		ParseNetworkCIDR("2001:440:ffff:ffff::/65"),
	}
	if !sliceOfSubnetsAreEqual(output, expected) {
		t.Error("\n",
			"<<<input>>>\n", input,
			"\n<<<actual_output>>>\n", spew.Sdump(output),
			"\n<<<expected_output>>>\n", spew.Sdump(expected),
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
	if !sliceOfSubnetsAreEqual(output, expected) {
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
	if !sliceOfSubnetsAreEqual(output, expected) {
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
	if !sliceOfSubnetsAreEqual(output, expected) {
		t.Error(
			"\n<<<input>>>\n", "aggregate:", aggregate, "\n", subnets,
			"\n<<<actual_output>>>\n", output,
			"\n<<<expected_output>>>\n", expected,
		)
	}
}

func BenchmarkIntToAddr(b *testing.B) {
	val := big.NewInt(3232235778)
	for i := 0; i < b.N; i++ {
		IntToAddr(val)
	}
}

func BenchmarkAddrToInt(b *testing.B) {
	address := net.ParseIP("192.168.1.2")
	for i := 0; i < b.N; i++ {
		AddrToInt(address)
	}
}

func BenchmarkParseNetworkCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseNetworkCIDR("192.168.1.2/24")
	}
}

func BenchmarkNetworkComesBefore(b *testing.B) {
	network := ParseNetworkCIDR("192.168.0.0/28")
	otherNetwork := ParseNetworkCIDR("192.168.0.0/22")
	for i := 0; i < b.N; i++ {
		NetworkComesBefore(network, otherNetwork)
	}
}

func BenchmarkIPv4ClassfulNetwork(b *testing.B) {
	address := net.ParseIP("192.168.0.0")
	for i := 0; i < b.N; i++ {
		IPv4ClassfulNetwork(address)
	}
}

func BenchmarkNextAddr(b *testing.B) {
	address := net.ParseIP("192.168.0.0")
	for i := 0; i < b.N; i++ {
		address = NextAddr(address)
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

func BenchmarkFindInbetweenSubnets(b *testing.B) {
	alpha := net.ParseIP("192.168.0.0")
	bravo := net.ParseIP("192.168.3.255")
	for i := 0; i < b.N; i++ {
		FindInbetweenSubnets(alpha, bravo)
	}
}

func BenchmarkFindInbetweenSubnetsBuffered(b *testing.B) {
	alpha := net.ParseIP("192.168.0.0")
	bravo := net.ParseIP("192.168.3.255")
	buf := NewBuffer()
	for i := 0; i < b.N; i++ {
		buf.FindInbetweenSubnets(alpha, bravo)
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
