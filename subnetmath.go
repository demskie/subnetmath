package subnetmath

import (
	"bytes"
	"net"
)

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

func NetworkComesBefore(first, second *net.IPNet) bool {
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

func GetClassfulNetwork(oldIP net.IP) *net.IPNet {
	var newIP net.IP
	var newMask net.IPMask
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

func SubnetZeroAddr(subnet *net.IPNet) net.IP {
	return subnet.IP.Mask(subnet.Mask)
}

func DuplicateAddr(addr net.IP) net.IP {
	newIP := make(net.IP, len(addr))
	copy(newIP, addr)
	return newIP
}

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

func GetAllSubnetHosts(subnet *net.IPNet, maximum int, includeBroadcast bool) []net.IP {
	hosts := make([]net.IP, 0)
	currentIP := SubnetZeroAddr(subnet)
	firstIP := DuplicateAddr(currentIP)
	for {
		if len(hosts) >= maximum {
			break
		}
		if subnet.Contains(currentIP) &&
			bytes.Equal(currentIP, firstIP) == false || len(hosts) == 0 {
			hosts = append(hosts, currentIP)
			currentIP = NextAddr(currentIP)
		} else {
			break
		}
	}
	ones, _ := subnet.Mask.Size()
	if ones <= 30 && subnet.IP.To4() != nil {
		hosts = hosts[:len(hosts)-1]
	}
	return hosts
}

func BlindlyParseCIDR(cidr string) *net.IPNet {
	addr, network, _ := net.ParseCIDR(cidr)
	if network.IP.Equal(addr) {
		return network
	}
	return nil
}

func FindNetworkIntersection(network *net.IPNet, otherNetworks ...*net.IPNet) *net.IPNet {
	for _, otherNetwork := range otherNetworks {
		if network.Contains(otherNetwork.IP) || otherNetwork.Contains(network.IP) {
			return otherNetwork
		}
	}
	return nil
}

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

func FindMaskWithoutIntersection(network *net.IPNet, otherNetworks ...*net.IPNet) *net.IPNet {
	currentNetwork := DuplicateNetwork(network)
	for {
		currentNetwork = ShrinkNetwork(currentNetwork)
		if currentNetwork == nil {
			break
		}
		if FindNetworkIntersection(currentNetwork, otherNetworks...) == nil {
			if SubnetZeroAddr(currentNetwork).Equal(currentNetwork.IP) {
				return currentNetwork
			}
		}
	}
	return nil
}

func NextNetwork(network *net.IPNet) *net.IPNet {
	newNetwork := DuplicateNetwork(network)
	ones, bits := newNetwork.Mask.Size()
	hosts := 2 << uint((bits-1)-ones)
	for i := 0; i < hosts; i++ {
		for octet := len(newNetwork.IP) - 1; octet >= 0; octet-- {
			newNetwork.IP[octet]++
			if uint8(newNetwork.IP[octet]) > 0 {
				break
			}
		}
	}
	return newNetwork
}

func UnusedSubnets(aggregate *net.IPNet, subnets ...*net.IPNet) (unused []*net.IPNet) {
	if FindNetworkIntersection(aggregate, subnets...) == nil {
		return
	}
	newSubnet := DuplicateNetwork(aggregate)
	for aggregate.Contains(newSubnet.IP) {
		canidateSubnet := FindMaskWithoutIntersection(newSubnet, subnets...)
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
