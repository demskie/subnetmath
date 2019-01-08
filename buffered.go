package subnetmath

import (
	"bytes"
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
	bigIntEcho    *big.Int
	ipSubZero     [16]byte
}

func NewBuffer() *Buffer {
	return &Buffer{
		mtx:           &sync.Mutex{},
		bigIntAlpha:   new(big.Int),
		bigIntBravo:   new(big.Int),
		bigIntCharlie: new(big.Int),
		bigIntDelta:   new(big.Int),
		bigIntEcho:    new(big.Int),
		ipSubZero:     [16]byte{},
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
			supernetInt.Add(supernetInt, b.addressCountCharlieDelta(network))
			subnetInt.Add(subnetInt, b.addressCountCharlieDelta(subnet))
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

func (b *Buffer) addressCountCharlieDelta(network *net.IPNet) *big.Int {
	if network != nil {
		ones, bits := network.Mask.Size()
		return b.bigIntCharlie.Exp(bigTwo, b.bigIntDelta.SetInt64(int64(bits-ones)), nil)
	}
	return nil
}

func (b *Buffer) nextNetworkEcho(network *net.IPNet) *net.IPNet {
	if network != nil {
		nextNetwork := DuplicateNetwork(network)
		v4addr := network.IP.To4()
		if v4addr != nil {
			b.bigIntEcho.SetBytes(v4addr)
		} else {
			b.bigIntEcho.SetBytes(network.IP.To16())
		}
		b.bigIntEcho.Add(b.bigIntEcho, b.addressCountCharlieDelta(network))
		nextNetwork.IP = IntToAddr(b.bigIntEcho)
		return nextNetwork
	}
	return nil
}

// FindInbetweenSubnets returns a slice of subnets given a range of IP addresses.
// Note that the delimiter 'stop' is inclusive. In other words, it will be included in the result.
func (b *Buffer) FindInbetweenSubnets(start, stop net.IP) []*net.IPNet {
	if sameAddrType(start, stop) && b.AddressComesBefore(start, stop) {
		var subnets []*net.IPNet
		maskBits := maskBitLength(start)
		current := DuplicateAddr(start)
		stopInt := b.addrToIntAlpha(stop)
		for {
			currentSubnet := &net.IPNet{
				IP:   current,
				Mask: make(net.IPMask, maskBits/8),
			}
			for ones := 1; ones <= maskBits; ones++ {
				currentSubnet.Mask = recreateMask(currentSubnet.Mask, ones, maskBits)
				increment := b.addressCountCharlieDelta(currentSubnet)
				addressInt := b.addrToIntBravo(currentSubnet.IP)
				addressInt.Add(addressInt, increment)
				addressInt.Sub(addressInt, bigOne)
				if addressInt.Cmp(stopInt) > 0 {
					continue
				}
				if b.SubnetZeroAddr(currentSubnet.IP, currentSubnet).Equal(currentSubnet.IP) {
					break
				}
			}
			subnets = append(subnets, currentSubnet)
			current = b.nextNetworkEcho(currentSubnet).IP
			if b.AddressComesBefore(current, start) {
				break
			}
			if b.AddressComesBefore(stop, current) && !current.Equal(stop) {
				break
			}
		}
		return subnets
	}
	return nil
}

func allFF(b []byte) bool {
	for _, c := range b {
		if c != 0xff {
			return false
		}
	}
	return true
}

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

func applyMaskDirectly(ip net.IP, mask net.IPMask) net.IP {
	if len(mask) == net.IPv6len && len(ip) == net.IPv4len && allFF(mask[:12]) {
		mask = mask[12:]
	}
	if len(mask) == net.IPv4len && len(ip) == net.IPv6len && bytes.Equal(ip[:12], v4InV6Prefix) {
		ip = ip[12:]
	}
	if len(ip) == len(mask) {
		for i := 0; i < len(ip); i++ {
			ip[i] = ip[i] & mask[i]
		}
		return ip
	}
	return nil
}

// SubnetZeroAddr returns the subnet zero address
func (b *Buffer) SubnetZeroAddr(address net.IP, network *net.IPNet) net.IP {
	if network != nil {
		b.mtx.Lock()
		defer b.mtx.Unlock()
		for i := 0; i < len(address); i++ {
			b.ipSubZero[i] = address[i]
		}
		return applyMaskDirectly(b.ipSubZero[:len(address)], network.Mask)
	}
	return nil
}

func recreateMask(mask net.IPMask, ones, bits int) net.IPMask {
	if bits != 8*net.IPv4len && bits != 8*net.IPv6len {
		return nil
	}
	if ones < 0 || ones > bits {
		return nil
	}
	l := bits / 8
	n := uint(ones)
	for i := 0; i < l; i++ {
		if n >= 8 {
			mask[i] = 0xff
			n -= 8
			continue
		}
		mask[i] = ^byte(0xff >> n)
		n = 0
	}
	return mask
}
