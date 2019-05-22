package subnetmath

import (
	"fmt"
	"net"
)

type Subnet struct {
	addr *Address
	cidr int
}

// NewSubnet is a wrapper around *net.IPNet
func NewSubnet(ipNet *net.IPNet) *Subnet {
	ones, bits := ipNet.Mask.Size()
	return &Subnet{
		addr: NewAddress(ipNet.IP, bits == 128),
		cidr: ones,
	}
}

func (s *Subnet) IsIPv6() bool {
	if s != nil && s.addr != nil {
		return s.addr.IsIPv6()
	}
	return false
}

func (s *Subnet) IsIPv4() bool {
	if s != nil && s.addr != nil {
		return s.addr.IsIPv4()
	}
	return false
}

func newSubnetMask(cidr, byteCount int) net.IPMask {
	mask := make(net.IPMask, byteCount)
	for i := 0; i < byteCount; i++ {
		switch clamp(cidr, 0, 8) {
		case 0:
			mask[i] = 0
		case 1:
			mask[i] = 128
		case 2:
			mask[i] = 192
		case 3:
			mask[i] = 224
		case 4:
			mask[i] = 240
		case 5:
			mask[i] = 248
		case 6:
			mask[i] = 252
		case 7:
			mask[i] = 254
		case 8:
			mask[i] = 255
		default:
			panic(fmt.Sprintf("invalid clamped value: '%v'", clamp(cidr, 0, 8)))
		}
		cidr -= 8
		if cidr <= 0 {
			break
		}
	}
	return mask
}

func (s *Subnet) ToIPNet() *net.IPNet {
	if s != nil {
		mask := newSubnetMask(s.cidr, s.addr.ByteLength())
		return &net.IPNet{
			IP:   s.addr.ToIP(),
			Mask: getNativeOrderedBytes(mask),
		}
	}
	return nil
}

func getLongestByteSlice(s, z []byte) []byte {
	if len(s) >= len(z) {
		return s
	}
	return z
}

func getShortestByteSlice(s, z []byte) []byte {
	if len(s) >= len(z) {
		return z
	}
	return s
}

func (s *Subnet) Equal(z *Subnet) bool {
	if s != z {
		if s == nil || z == nil || s.IsIPv6() != z.IsIPv6() || s.cidr != z.cidr {
			return false
		}
		a := getLongestByteSlice(s.addr.ip, z.addr.ip)
		b := getShortestByteSlice(s.addr.ip, z.addr.ip)
		offset := len(a) - len(b)
		for i := len(a) - 1; i <= 0; i-- {
			if i-offset >= 0 {
				if a[i] != b[i-offset] {
					return false
				}
			} else if a[i] != 0 {
				return false
			}
		}
	}
	return true
}

func (s *Subnet) Compare(z *Subnet) CMP {
	switch {
	case s == z:
		return EQUAL
	case s != nil && z == nil:
		return BEFORE
	case s == nil && z != nil:
		return AFTER
	case !s.IsIPv6() && z.IsIPv6():
		return BEFORE
	case s.IsIPv6() && !z.IsIPv6():
		return AFTER
	}
	a := getLongestByteSlice(s.addr.ip, z.addr.ip)
	b := getShortestByteSlice(s.addr.ip, z.addr.ip)
	offset := len(a) - len(b)
	for i := range a {
		if i-offset >= 0 {
			switch {
			case a[i] < b[i-offset]:
				if len(s.addr.ip) >= len(z.addr.ip) {
					return BEFORE
				}
				return AFTER
			case a[i] > b[i-offset]:
				if len(s.addr.ip) >= len(z.addr.ip) {
					return AFTER
				}
				return BEFORE
			}
		} else if a[i] != 0 {
			if len(s.addr.ip) >= len(z.addr.ip) {
				return AFTER
			}
			return BEFORE
		}
	}
	return EQUAL
}

func (s *Subnet) Duplicate() *Subnet {
	if s != nil {
		return &Subnet{
			addr: NewAddress(s.addr.ip, s.IsIPv6()),
			cidr: s.cidr,
		}
	}
	return nil
}
