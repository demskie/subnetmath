package subnetmath

import (
	"math/big"
	"net"
	"sync"
)

type Buffer struct {
	mtx           *sync.Mutex
	bigIntAlpha   *big.Int
	bigIntBravo   *big.Int
	bigIntCharlie *big.Int
	bigIntDelta   *big.Int
}

func NewBuffer() *Buffer {
	return &Buffer{
		mtx:           &sync.Mutex{},
		bigIntAlpha:   new(big.Int),
		bigIntBravo:   new(big.Int),
		bigIntCharlie: new(big.Int),
		bigIntDelta:   new(big.Int),
	}
}

// NetworkComesBefore returns a bool with regards to numerical network order.
// Note that IPv4 networks come before IPv6 networks.
func (b *Buffer) NetworkComesBefore(first, second *net.IPNet) bool {
	if first != nil && second != nil {
		if first.IP.Equal(second.IP) {
			firstMask, _ := first.Mask.Size()
			secondMask, _ := second.Mask.Size()
			if firstMask < secondMask {
				return true
			}
			return false
		}
		return b.AddressComesBefore(first.IP, second.IP)
	}
	return false
}

// AddressComesBefore returns a bool with regards to numerical address order.
// Note that IPv4 addresses come before IPv6 addresses.
func (b *Buffer) AddressComesBefore(firstIP, secondIP net.IP) bool {
	if firstIP.To4() == nil && secondIP.To4() != nil {
		return true
	} else if firstIP.To4() != nil && secondIP.To4() == nil {
		return false
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.addrToIntAlpha(firstIP).Cmp(b.addrToIntBravo(secondIP)) < 0 {
		return true
	}
	return false
}

// NetworkContainsSubnet validates that the network is a valid supernet
func (b *Buffer) NetworkContainsSubnet(network *net.IPNet, subnet *net.IPNet) bool {
	if network != nil && subnet != nil {
		b.mtx.Lock()
		defer b.mtx.Unlock()
		supernetInt := b.addrToIntAlpha(network.IP)
		subnetInt := b.addrToIntBravo(subnet.IP)
		if supernetInt.Cmp(subnetInt) <= 0 {
			supernetInt.Add(supernetInt, b.addressCountCharlie(network))
			subnetInt.Add(subnetInt, b.addressCountCharlie(subnet))
			if supernetInt.Cmp(subnetInt) >= 0 {
				return true
			}
		}
	}
	return false
}

func (b *Buffer) addrToIntAlpha(address net.IP) *big.Int {
	v4addr := address.To4()
	if v4addr != nil {
		b.bigIntAlpha.SetBytes(v4addr)
	} else {
		b.bigIntAlpha.SetBytes(address.To16())
	}
	return b.bigIntAlpha
}

func (b *Buffer) addrToIntBravo(address net.IP) *big.Int {
	v4addr := address.To4()
	if v4addr != nil {
		b.bigIntBravo.SetBytes(v4addr)
	} else {
		b.bigIntBravo.SetBytes(address.To16())
	}
	return b.bigIntBravo
}

func (b *Buffer) addressCountCharlie(network *net.IPNet) *big.Int {
	if network != nil {
		ones, bits := network.Mask.Size()
		return b.bigIntCharlie.Exp(bigTwo, b.bigIntDelta.SetInt64(int64(bits-ones)), nil)
	}
	return nil
}
