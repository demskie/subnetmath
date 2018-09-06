package subnetmath

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
	"net"
)

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
	if first.IP.To4() == nil && second.IP.To4() != nil {
		return true
	} else if first.IP.To4() != nil && second.IP.To4() == nil {
		return false
	}
	firstBytes, secondBytes := []byte(first.IP), []byte(second.IP)
	difference := bytes.Compare(firstBytes, secondBytes)
	if difference > 0 {
		return true
	} else if difference < 0 {
		return false
	} else {
		firstMask, _ := first.Mask.Size()
		secondMask, _ := second.Mask.Size()
		if firstMask < secondMask {
			return true
		}
	}
	return false
}

// GetClassfulNetwork either return the classful network given an IPv4 address or
// return nil if given a multicast address or IPv6 address
func GetClassfulNetwork(oldIP net.IP) *net.IPNet {
	if oldIP.To4() == nil {
		return nil
	}
	var (
		newIP   net.IP
		newMask net.IPMask
	)
	switch {
	case uint8(oldIP[0]) < 128:
		newIP = net.IPv4(uint8(oldIP[0]), 0, 0, 0)
		newMask = net.IPv4Mask(255, 0, 0, 0)
	case uint8(oldIP[0]) < 192:
		newIP = net.IPv4(uint8(oldIP[0]), uint8(oldIP[1]), 0, 0)
		newMask = net.IPv4Mask(255, 255, 0, 0)
	case uint8(oldIP[0]) < 224:
		newIP = net.IPv4(uint8(oldIP[0]), uint8(oldIP[1]), uint8(oldIP[2]), 0)
		newMask = net.IPv4Mask(255, 255, 255, 0)
	default:
		return nil
	}
	return &net.IPNet{IP: newIP, Mask: newMask}
}

// DuplicateNetwork returns a new copy of *net.IPNet
func DuplicateNetwork(network *net.IPNet) *net.IPNet {
	newIP := make(net.IP, len(network.IP))
	newMask := make(net.IPMask, len(network.Mask))
	copy(newIP, network.IP)
	copy(newMask, network.Mask)
	return &net.IPNet{
		IP:   newIP,
		Mask: newMask,
	}
}

// SubnetZeroAddr returns the subnet zero address
func SubnetZeroAddr(subnet *net.IPNet) net.IP {
	return subnet.IP.Mask(subnet.Mask)
}

// DuplicateAddr creates a new copy of net.IP
func DuplicateAddr(addr net.IP) net.IP {
	newIP := make(net.IP, len(addr))
	copy(newIP, addr)
	return newIP
}

// AddToAddr returns the same net.IP with its address incremented by val
func AddToAddr(addr net.IP, val int) net.IP {
	for i := 0; i < val; i++ {
		for octet := len(addr) - 1; octet >= 0; octet-- {
			if val > 0 {
				addr[octet]++
			} else {
				addr[octet]--
			}
			if uint8(addr[octet]) > 0 {
				break
			}
		}
	}
	return addr
}

// NextAddr returns a new net.IP that is the next address
func NextAddr(addr net.IP) net.IP {
	newIP := DuplicateAddr(addr)
	for octet := len(newIP) - 1; octet >= 0; octet-- {
		newIP[octet]++
		if uint8(newIP[octet]) > 0 {
			break
		}
	}
	return newIP
}

func maxInt() int {
	if binary.Size(int(0)) == 4 {
		return math.MaxInt32
	}
	return math.MaxInt64
}

// AddressCountInt will return -1 if the number of addresses would overflow an Int
func AddressCountInt(network *net.IPNet) int {
	val := AddressCountBigInt(network)
	limit := big.NewInt(int64(maxInt()))
	if val.Cmp(limit) > 0 {
		return -1
	}
	return int(val.Int64())
}

// AddressCountBigInt returns the number of addresses
func AddressCountBigInt(network *net.IPNet) *big.Int {
	ones, bits := network.Mask.Size()
	return new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits-ones)), nil)
}

// GetAllAddresses will return a slice of net.IPs for the subnet specified.
// However nil is returned instead if the slice length is greater than the max int value
func GetAllAddresses(subnet *net.IPNet) []net.IP {
	numAddresses := AddressCountInt(subnet)
	if numAddresses < 0 {
		return nil
	}
	results := make([]net.IP, numAddresses)
	currentIP := SubnetZeroAddr(subnet)
	for i := range results {
		results[i] = currentIP
		currentIP = NextAddr(currentIP)
	}
	return results
}

// ParseNetworkCIDR is a convienence function that will return either the *net.IPNet
// or nil if the supplied cidr is invalid
func ParseNetworkCIDR(cidr string) *net.IPNet {
	addr, network, err := net.ParseCIDR(cidr)
	if err == nil && network != nil && network.IP.Equal(addr) {
		return network
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

// ShrinkNetwork increases the mask size by one
func ShrinkNetwork(network *net.IPNet) *net.IPNet {
	ones, bits := network.Mask.Size()
	if ones < bits {
		return &net.IPNet{
			IP:   network.IP,
			Mask: net.CIDRMask(ones+1, bits),
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
		if findNetworkIntersection(currentNetwork, otherNetworks...) == nil {
			if SubnetZeroAddr(currentNetwork).Equal(currentNetwork.IP) {
				return currentNetwork
			}
		}
	}
	return nil
}

func isZero(val *big.Int) bool {
	return val.Cmp(big.NewInt(0)) == 0
}

func addOne(val *big.Int) *big.Int {
	return val.Add(val, big.NewInt(1))
}

// NextNetwork returns the next network of the same size
func NextNetwork(network *net.IPNet) *net.IPNet {
	newNetwork := DuplicateNetwork(network)
	maskSize := AddressCountBigInt(newNetwork)
	for i := big.NewInt(0); isZero(i) || i.Cmp(maskSize) < 0; addOne(i) {
		for octet := len(newNetwork.IP) - 1; octet >= 0; octet-- {
			newNetwork.IP[octet]++
			if uint8(newNetwork.IP[octet]) > 0 {
				break
			}
		}
	}
	return newNetwork
}

// FindUnusedSubnets returns a slice of unused subnets given the aggregate and sibling subnets
func FindUnusedSubnets(aggregate *net.IPNet, subnets ...*net.IPNet) (unused []*net.IPNet) {
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
