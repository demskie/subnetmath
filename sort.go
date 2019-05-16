package subnetmath

import "net"

// SortNetworks radix sorts the provided IPv4 and IPv6 network slice
func SortNetworks(networks []*net.IPNet) {
	// create array of buckets to reuse between iterations
	var (
		runningPrefixSum = [256]int{}
		offsetPrefixSum  = [256]int{}
		counts           = runningPrefixSum
	)

	// determine maxByteLength
	maxByteLength := 0
	for _, network := range networks {
		if len(network.IP) > maxByteLength {
			maxByteLength = len(network.IP)
		}
	}

	// iterate through and swap each byte in place until runningPrefixSum is
	for byteIndex := -1; byteIndex <= maxByteLength; byteIndex++ {
		// count each occurance of byte value
		byteValue := uint8(0)
		for _, network := range networks {
			switch byteIndex {
			case -1:
				// sort based on length of byte slice
				byteValue = uint8(len(network.IP))
			case maxByteLength:
				// sort based on cidr size
				cidr, _ := network.Mask.Size()
				byteValue = uint8(cidr)
			default:
				// sort based on uint value of byte
				byteValue = network.IP[byteIndex]
			}
			counts[byteValue]++
		}

		// building both prefixSums
		total := 0
		for i := 0; i < 256; i++ {
			oldCount := counts[i]
			runningPrefixSum[i] = total
			total += oldCount
			if i > 0 && i < 255 {
				offsetPrefixSum[i-1] = runningPrefixSum[i]
			}
		}
		offsetPrefixSum[255] = offsetPrefixSum[254]

		// in place swap and sort by value
		idx := 0
		for idx < len(networks) {
			switch byteIndex {
			case -1:
				// sort based on length of byte slice
				byteValue = uint8(len(networks[idx].IP))
			case maxByteLength:
				// sort based on cidr size
				cidr, _ := networks[idx].Mask.Size()
				byteValue = uint8(cidr)
			default:
				// sort based on uint value of byte
				byteValue = networks[idx].IP[byteIndex]
			}
			// check if this network wants to be at this index
			if runningPrefixSum[byteValue] != idx {
				if runningPrefixSum[byteValue] < offsetPrefixSum[byteValue] {
					// swap network at current index with provided prefixSum index
					oldNetwork := networks[idx]
					networks[idx] = networks[runningPrefixSum[byteValue]]
					networks[runningPrefixSum[byteValue]] = oldNetwork
				} else {
					idx++
				}
			} else {
				idx++
			}
			runningPrefixSum[byteValue]++
		}

		// break early to avoid the following unnecessary operation
		if byteIndex >= maxByteLength {
			break
		}

		// reset counts back to zero
		for i := 0; i < 256; i++ {
			runningPrefixSum[i] = 0
		}
	}

	return
}
