package subnetmath

import (
	"bytes"
	"math/big"
	"net"
)

// commonly used bigint values
var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

// ParseNetworkCIDR is a convienence function that will return either the *net.IPNet
// or nil if the supplied cidr is invalid
func ParseNetworkCIDR(cidr string) *net.IPNet {
	addr, network, err := net.ParseCIDR(cidr)
	if err == nil && network != nil && network.IP.Equal(addr) {
		return network
	}
	return nil
}

// NetworksAreIdentical returns a bool with regards to the two networks being equal
func NetworksAreIdentical(first, second *net.IPNet) bool {
	if first != nil && second != nil && first.IP.Equal(second.IP) {
		firstSize, _ := first.Mask.Size()
		secondSize, _ := second.Mask.Size()
		if firstSize == secondSize {
			return true
		}
	}
	return false
}

// NetworkComesBefore returns a bool with regards to numerical network order.
// Note that IPv4 networks come before IPv6 networks.
func NetworkComesBefore(first, second *net.IPNet) bool {
	if first.IP.Equal(second.IP) {
		firstMask, _ := first.Mask.Size()
		secondMask, _ := second.Mask.Size()
		if firstMask < secondMask {
			return true
		}
		return false
	}
	return AddressComesBefore(first.IP, second.IP)
}

// AddressComesBefore returns a bool with regards to numerical address order.
// Note that IPv4 addresses come before IPv6 addresses.
func AddressComesBefore(firstIP, secondIP net.IP) bool {
	if firstIP.To4() == nil && secondIP.To4() != nil {
		return true
	} else if firstIP.To4() != nil && secondIP.To4() == nil {
		return false
	}
	difference := bytes.Compare([]byte(firstIP), []byte(secondIP))
	if difference > 0 {
		return false
	}
	return true
}

// DuplicateNetwork returns a new copy of *net.IPNet
func DuplicateNetwork(network *net.IPNet) *net.IPNet {
	if network != nil {
		newIP := make(net.IP, len(network.IP))
		newMask := make(net.IPMask, len(network.Mask))
		copy(newIP, network.IP)
		copy(newMask, network.Mask)
		return &net.IPNet{
			IP:   newIP,
			Mask: newMask,
		}
	}
	return nil
}

// DuplicateAddr creates a new copy of net.IP
func DuplicateAddr(addr net.IP) net.IP {
	newIP := make(net.IP, len(addr))
	copy(newIP, addr)
	return newIP
}

// SubnetZeroAddr returns the subnet zero address
func SubnetZeroAddr(address net.IP, network *net.IPNet) net.IP {
	if network != nil {
		return address.Mask(network.Mask)
	}
	return nil
}

// NextNetwork returns the next network of the same size
func NextNetwork(network *net.IPNet) *net.IPNet {
	if network != nil {
		nextNetwork := DuplicateNetwork(network)
		networkInt := AddrToInt(network.IP)
		networkInt.Add(networkInt, addressCount(network))
		nextNetwork.IP = IntToAddr(networkInt)
		return nextNetwork
	}
	return nil
}

// BroadcastAddr returns the broadcast address
func BroadcastAddr(network *net.IPNet) net.IP {
	if network != nil {
		networkInt := AddrToInt(network.IP)
		networkInt.Add(networkInt, addressCount(network))
		networkInt.Sub(networkInt, bigOne)
		return IntToAddr(networkInt)
	}
	return nil
}

// NextAddr returns a new net.IP that is the next address
func NextAddr(addr net.IP) net.IP {
	addrInt := AddrToInt(addr)
	addrInt.Add(addrInt, bigOne)
	return IntToAddr(addrInt)
}

func addressCount(network *net.IPNet) *big.Int {
	if network != nil {
		ones, bits := network.Mask.Size()
		return new(big.Int).Exp(bigTwo, big.NewInt(int64(bits-ones)), nil)
	}
	return nil
}

// GetAllNetworkAddresses will return a limited slice of net.IPs for the subnet specified
func GetAllNetworkAddresses(network *net.IPNet, limit int) []net.IP {
	if network != nil {
		originalAddrInt := AddrToInt(network.IP)
		incrementInt := addressCount(network)
		var results []net.IP
		for i := new(big.Int).Set(bigZero); i.Cmp(incrementInt) < 0; i.Add(i, bigOne) {
			val := new(big.Int).Set(originalAddrInt)
			results = append(results, IntToAddr(val.Add(val, i)))
			if len(results) >= limit {
				break
			}
		}
		return results
	}
	return nil
}

// ShrinkNetwork increases the mask size by one
func ShrinkNetwork(network *net.IPNet) *net.IPNet {
	if network != nil {
		ones, bits := network.Mask.Size()
		if ones < bits {
			return &net.IPNet{
				IP:   network.IP,
				Mask: net.CIDRMask(ones+1, bits),
			}
		}
	}
	return nil
}

func sameAddrType(first, second net.IP) bool {
	alpha, bravo := first.To4(), second.To4()
	if alpha == nil && bravo == nil || alpha != nil && bravo != nil {
		return true
	}
	return false
}

func maskBitLength(address net.IP) int {
	if address.To4() != nil {
		return 32
	}
	return 128
}

// FindInbetweenSubnets returns a slice of subnets given a range of IP addresses.
// Note that the delimiter 'stop' is inclusive. In other words, it will be included in the result.
func FindInbetweenSubnets(start, stop net.IP) []*net.IPNet {
	if sameAddrType(start, stop) && AddressComesBefore(start, stop) {
		var subnets []*net.IPNet
		maskBits := maskBitLength(start)
		current := DuplicateAddr(start)
		stopInt := AddrToInt(stop)
		for {
			currentSubnet := &net.IPNet{IP: current}
			for ones := 1; ones <= maskBits; ones++ {
				currentSubnet.Mask = net.CIDRMask(ones, maskBits)
				increment := addressCount(currentSubnet)
				addressInt := AddrToInt(currentSubnet.IP)
				addressInt.Add(addressInt, increment)
				addressInt.Sub(addressInt, bigOne)
				if addressInt.Cmp(stopInt) > 0 {
					continue
				}
				if SubnetZeroAddr(currentSubnet.IP, currentSubnet).Equal(currentSubnet.IP) {
					break
				}
			}
			subnets = append(subnets, currentSubnet)
			current = NextNetwork(currentSubnet).IP
			if !current.Equal(stop) && AddressComesBefore(stop, current) ||
				AddressComesBefore(current, start) {
				break
			}
		}
		return subnets
	}
	return nil
}

func findNetworkIntersection(network *net.IPNet, otherNetworks ...*net.IPNet) *net.IPNet {
	for _, otherNetwork := range otherNetworks {
		if network.Contains(otherNetwork.IP) || otherNetwork.Contains(network.IP) {
			return otherNetwork
		}
	}
	return nil
}

func findMaskWithoutIntersection(network *net.IPNet, otherNetworks ...*net.IPNet) *net.IPNet {
	currentNetwork := DuplicateNetwork(network)
	for {
		currentNetwork = ShrinkNetwork(currentNetwork)
		if currentNetwork == nil {
			break
		}
		if findNetworkIntersection(currentNetwork, otherNetworks...) == nil &&
			SubnetZeroAddr(currentNetwork.IP, currentNetwork).Equal(currentNetwork.IP) {
			return currentNetwork
		}
	}
	return nil
}

// FindUnusedSubnets returns a slice of unused subnets given the aggregate and sibling subnets
func FindUnusedSubnets(aggregate *net.IPNet, subnets ...*net.IPNet) (unused []*net.IPNet) {
	// BUG: need to refactor using the new logic
	if len(subnets) == 0 {
		return []*net.IPNet{aggregate}
	}
	if findNetworkIntersection(aggregate, subnets...) == nil {
		return nil
	}
	newSubnet := DuplicateNetwork(aggregate)
	var canidateSubnet *net.IPNet
	for aggregate.Contains(newSubnet.IP) {
		canidateSubnet = findMaskWithoutIntersection(newSubnet, subnets...)
		if canidateSubnet != nil {
			unused = append(unused, canidateSubnet)
			newSubnet = NextNetwork(canidateSubnet)
			newSubnet.Mask = aggregate.Mask
		} else {
			newSubnet.IP = NextAddr(newSubnet.IP)
		}
	}
	return unused
}

// IntToAddr will return the net.IP of the big.Int represented address
func IntToAddr(intAddress *big.Int) net.IP {
	intBytes := intAddress.Bytes()
	if len(intBytes) == 4 {
		return net.IPv4(intBytes[0], intBytes[1], intBytes[2], intBytes[3])
	}
	return intBytes
}

// AddrToInt will return the *bit.Int of a given IPv6 address
func AddrToInt(address net.IP) *big.Int {
	v4addr := address.To4()
	if v4addr != nil {
		return big.NewInt(0).SetBytes(v4addr)
	}
	return big.NewInt(0).SetBytes(address.To16())
}

// IPv4ClassfulNetwork eithers return the classful network given an IPv4 address or
// returns nil if given a multicast address or IPv6 address
func IPv4ClassfulNetwork(address net.IP) *net.IPNet {
	if address.To4() != nil {
		var newIP net.IP
		var newMask net.IPMask
		switch {
		case uint8(address[0]) < 128:
			newIP = net.IPv4(uint8(address[0]), 0, 0, 0)
			newMask = net.IPv4Mask(255, 0, 0, 0)
		case uint8(address[0]) < 192:
			newIP = net.IPv4(uint8(address[0]), uint8(address[1]), 0, 0)
			newMask = net.IPv4Mask(255, 255, 0, 0)
		case uint8(address[0]) < 224:
			newIP = net.IPv4(uint8(address[0]), uint8(address[1]), uint8(address[2]), 0)
			newMask = net.IPv4Mask(255, 255, 255, 0)
		default:
			return nil
		}
		return &net.IPNet{IP: newIP, Mask: newMask}
	}
	return nil
}
