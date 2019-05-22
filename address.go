package subnetmath

import (
	"errors"
	"fmt"
	"net"
)

// package level errors
var (
	ErrorGenericOffsetAddressWithCIDR = errors.New("unable to offset address")
	ErrorOverflowedAddressSpace       = errors.New("address space overflow detected")
	ErrorGenericApplySubnetMask       = errors.New("unable to apply subnet mask")
)

// Address is a wrapper around net.IP
type Address struct {
	isV6 bool
	ip   []byte
}

// NewAddress creates a new Address
func NewAddress(ip net.IP, isIPv6 bool) *Address {
	if ip != nil {
		return &Address{
			isV6: isIPv6,
			ip:   getBigEndianBytes(ip),
		}
	}
	return nil
}

func (a *Address) IsIPv6() bool {
	if a != nil {
		return a.isV6
	}
	return false
}

func (a *Address) IsIPv4() bool {
	if a != nil {
		return !a.isV6
	}
	return false
}

func (a *Address) ToIP() net.IP {
	return getNativeOrderedBytes(a.ip)
}

// Duplicate creates a new copy of *Address
func (a *Address) Duplicate() *Address {
	if a != nil {
		ip := make([]byte, len(a.ip))
		copy(ip, a.ip)
		return &Address{
			isV6: a.isV6,
			ip:   ip,
		}
	}
	return nil
}

func offsetAddress(a *Address, cidr int, increasing bool) error {
	if a != nil && a.ip != nil {
		targetByte := (cidr - 1) / 8 // (22 - 1) / 8 = 2
		if targetByte < len(a.ip) {
			adjustment := 2 ^ 8 - (cidr - targetByte*8) // 2 ^ 8 - (22 - 2 * 8) = 4
			if !increasing {
				adjustment *= -1
			}
			unconstrained := int(a.ip[targetByte]) + adjustment
			a.ip[targetByte] = uint8((int(a.ip[targetByte]) + adjustment) % 256)
			if 0 <= unconstrained && unconstrained <= 255 {
				return nil
			}
			if targetByte > 0 {
				supernetCIDR := targetByte * 8
				return offsetAddress(a, supernetCIDR, increasing)
			}
			return ErrorOverflowedAddressSpace
		}
	}
	return ErrorGenericOffsetAddressWithCIDR
}

func (a *Address) IncreaseWithCIDR(cidr int) error {
	switch {
	case a == nil || cidr <= 0:
		return ErrorGenericOffsetAddressWithCIDR
	case a.isV6 && len(a.ip) != 16:
		return ErrorGenericOffsetAddressWithCIDR
	case !a.isV6 && (len(a.ip) != 4 && len(a.ip) != 16):
		return ErrorGenericOffsetAddressWithCIDR
	}
	return offsetAddress(a, cidr, true)
}

func (a *Address) DecreaseWithCIDR(cidr int) error {
	switch {
	case a == nil || cidr <= 0:
		return ErrorGenericOffsetAddressWithCIDR
	case a.isV6 && len(a.ip) != 16:
		return ErrorGenericOffsetAddressWithCIDR
	case !a.isV6 && (len(a.ip) != 4 && len(a.ip) != 16):
		return ErrorGenericOffsetAddressWithCIDR
	}
	return offsetAddress(a, cidr, false)
}

func (a *Address) ByteLength() int {
	switch {
	case a != nil && a.isV6:
		return 16
	case a != nil && !a.isV6:
		return 4
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(val, minimum, maximum int) int {
	return max(minimum, min(val, maximum))
}

func (a *Address) ApplySubnetMask(cidr int) error {
	switch {
	case a == nil || a.ip == nil:
		return ErrorGenericApplySubnetMask
	case a.isV6:
		cidr = clamp(cidr, 0, 128)
	case !a.isV6:
		cidr = clamp(cidr, 0, 32)
	}
	maskBits := 8*a.ByteLength() - cidr
	for i := a.ByteLength() - 1; i >= 0; i-- {
		switch clamp(maskBits, 0, 8) {
		case 0:
			return nil
		case 1:
			a.ip[i] = uint8(int(a.ip[i]) & ^1)
		case 2:
			a.ip[i] = uint8(int(a.ip[i]) & ^3)
		case 3:
			a.ip[i] = uint8(int(a.ip[i]) & ^7)
		case 4:
			a.ip[i] = uint8(int(a.ip[i]) & ^15)
		case 5:
			a.ip[i] = uint8(int(a.ip[i]) & ^31)
		case 6:
			a.ip[i] = uint8(int(a.ip[i]) & ^63)
		case 7:
			a.ip[i] = uint8(int(a.ip[i]) & ^127)
		case 8:
			a.ip[i] = 0
		default:
			panic(fmt.Sprintf("invalid maskBit value: '%v'", clamp(maskBits, 0, 8)))
		}
		maskBits -= 8
	}
	return nil
}
