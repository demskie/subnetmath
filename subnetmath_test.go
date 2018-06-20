package subnetmath

import (
	"net"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
		BlindlyParseCIDR("192.168.9.0/24"),
	}

	if reflect.DeepEqual(output, expectedOutput) == false {
		t.Error("\n",
			"<<<input>>>\n", "aggregate:", aggregate, "\n", spew.Sdump(subnets),
			"<<<actual_output>>>\n", spew.Sdump(output),
			"<<<expected_output>>>\n", spew.Sdump(expectedOutput),
		)
	}
}
