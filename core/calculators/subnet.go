package calculators

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sort"
	"strings"
)

var (
	ErrInvalidCIDR   = errors.New("invalid CIDR format, expected format like 192.168.1.0/24")
	ErrIPv6NotSupported = errors.New("IPv6 is not supported, please provide an IPv4 CIDR")
	ErrVLSMOverflow  = errors.New("requested subnets exceed available network address space")
	ErrInvalidSubnet = errors.New("subnet host requirement must be greater than 0")
)

type SubnetInfo struct {
	NetworkAddress   string
	BroadcastAddress string
	FirstHost        string
	LastHost         string
	TotalHosts       int
	Netmask          string
	CIDR             string
}

// ipToUint32 converts an IPv4 net.IP to uint32
func ipToUint32(ip net.IP) uint32 {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ipv4)
}

// uint32ToIP converts uint32 back to IPv4 string
func uint32ToIP(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

// ParseCIDR parses an IPv4 CIDR string and calculates network boundaries and usable host range
func ParseCIDR(cidr string) (SubnetInfo, error) {
	cidr = strings.TrimSpace(cidr)
	if !strings.Contains(cidr, "/") {
		return SubnetInfo{}, ErrInvalidCIDR
	}

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return SubnetInfo{}, fmt.Errorf("%w: %v", ErrInvalidCIDR, err)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return SubnetInfo{}, ErrIPv6NotSupported
	}

	ones, bits := ipNet.Mask.Size()
	if bits != 32 {
		return SubnetInfo{}, ErrIPv6NotSupported
	}

	netIPUint := ipToUint32(ipNet.IP)
	maskUint := ipToUint32(net.IP(ipNet.Mask))
	wildcard := ^maskUint
	bcastUint := netIPUint | wildcard

	var firstHostUint, lastHostUint uint32
	var totalHosts int

	if ones == 32 {
		totalHosts = 1
		firstHostUint = netIPUint
		lastHostUint = netIPUint
	} else if ones == 31 {
		totalHosts = 2
		firstHostUint = netIPUint
		lastHostUint = bcastUint
	} else {
		totalHosts = int(math.Pow(2, float64(32-ones))) - 2
		firstHostUint = netIPUint + 1
		lastHostUint = bcastUint - 1
	}

	maskIP := net.IP(ipNet.Mask)

	return SubnetInfo{
		NetworkAddress:   uint32ToIP(netIPUint),
		BroadcastAddress: uint32ToIP(bcastUint),
		FirstHost:        uint32ToIP(firstHostUint),
		LastHost:         uint32ToIP(lastHostUint),
		TotalHosts:       totalHosts,
		Netmask:          maskIP.String(),
		CIDR:             fmt.Sprintf("%s/%d", uint32ToIP(netIPUint), ones),
	}, nil
}

type vlsReq struct {
	index     int
	hostCount int
}

// CalculateVLSM calculates optimal subnet divisions based on requested host counts
func CalculateVLSM(network string, subnets []int) ([]SubnetInfo, error) {
	if len(subnets) == 0 {
		return nil, nil
	}

	baseInfo, err := ParseCIDR(network)
	if err != nil {
		return nil, err
	}

	// Prepare requests with original index preserved
	reqs := make([]vlsReq, len(subnets))
	for i, count := range subnets {
		if count <= 0 {
			return nil, ErrInvalidSubnet
		}
		reqs[i] = vlsReq{index: i, hostCount: count}
	}

	// Sort descending by host count (VLSM standard rule)
	sort.Slice(reqs, func(i, j int) bool {
		return reqs[i].hostCount > reqs[j].hostCount
	})

	_, baseNet, _ := net.ParseCIDR(baseInfo.CIDR)
	baseStart := ipToUint32(baseNet.IP)
	ones, _ := baseNet.Mask.Size()
	baseCapacity := uint32(1 << (32 - ones))
	baseEnd := baseStart + baseCapacity

	currentStart := baseStart
	results := make([]SubnetInfo, len(subnets))

	for _, req := range reqs {
		// Minimum power of 2 needed to accommodate req.hostCount + 2 (network + broadcast)
		neededHosts := req.hostCount + 2
		hostBits := 0
		for (1 << hostBits) < neededHosts {
			hostBits++
		}
		if hostBits < 2 {
			hostBits = 2
		}

		prefix := 32 - hostBits
		blockSize := uint32(1 << hostBits)

		if currentStart+blockSize > baseEnd {
			return nil, ErrVLSMOverflow
		}

		subCIDR := fmt.Sprintf("%s/%d", uint32ToIP(currentStart), prefix)
		info, err := ParseCIDR(subCIDR)
		if err != nil {
			return nil, err
		}

		results[req.index] = info
		currentStart += blockSize
	}

	return results, nil
}

// SubnetRecommendation holds the suggested subnet parameters for a requested host count
type SubnetRecommendation struct {
	NeededHosts    int
	AllocatedHosts int
	Prefix         int
	CIDR           string
	Netmask        string
	WildcardMask   string
	FirstHost      string
	LastHost       string
	Broadcast      string
	WastedHosts    int
	Efficiency     float64
	ClassHint      string
}

// FindSubnetForHosts calculates the optimal IPv4 subnet and mask to accommodate a target host count
func FindSubnetForHosts(neededHosts int, baseIP string) (SubnetRecommendation, error) {
	if neededHosts <= 0 {
		return SubnetRecommendation{}, ErrInvalidSubnet
	}
	if neededHosts > 1073741822 { // Max usable hosts in IPv4 (/2)
		return SubnetRecommendation{}, errors.New("jumlah host melebihi kapasitas pengalamatan IPv4")
	}

	baseIP = strings.TrimSpace(baseIP)
	if baseIP == "" {
		baseIP = "192.168.1.0"
	} else if strings.Contains(baseIP, "/") {
		parts := strings.Split(baseIP, "/")
		baseIP = strings.TrimSpace(parts[0])
	}

	parsedIP := net.ParseIP(baseIP)
	if parsedIP == nil || parsedIP.To4() == nil {
		return SubnetRecommendation{}, errors.New("format IP dasar tidak valid, gunakan format seperti 192.168.1.0")
	}

	// Calculate needed host bits h such that (2^h - 2) >= neededHosts
	hostBits := 2
	for (1<<hostBits)-2 < neededHosts {
		hostBits++
		if hostBits > 30 {
			return SubnetRecommendation{}, errors.New("jumlah host melebihi batas maksimum subnet IPv4")
		}
	}

	prefix := 32 - hostBits
	allocatedHosts := (1 << hostBits) - 2
	wastedHosts := allocatedHosts - neededHosts
	efficiency := (float64(neededHosts) / float64(allocatedHosts)) * 100.0

	// Align base IP to subnet boundary
	ipUint := ipToUint32(parsedIP)
	maskUint := ^uint32((1 << hostBits) - 1)
	netUint := ipUint & maskUint

	cidr := fmt.Sprintf("%s/%d", uint32ToIP(netUint), prefix)
	info, err := ParseCIDR(cidr)
	if err != nil {
		return SubnetRecommendation{}, err
	}

	wildcardUint := ^maskUint
	wildcardIP := uint32ToIP(wildcardUint)

	var classHint string
	switch {
	case prefix >= 24:
		classHint = fmt.Sprintf("Prefix /%d (Setara Class C) — Sangat ideal untuk ruang kerja, lab komputer, atau sub-divisi kantor.", prefix)
	case prefix >= 16:
		classHint = fmt.Sprintf("Prefix /%d (Setara Class B) — Cocok untuk gedung kampus, jaringan enterprise, atau datacenter.", prefix)
	default:
		classHint = fmt.Sprintf("Prefix /%d (Setara Class A) — Alokasi skala besar untuk infrastruktur backbone atau provider.", prefix)
	}

	return SubnetRecommendation{
		NeededHosts:    neededHosts,
		AllocatedHosts: allocatedHosts,
		Prefix:         prefix,
		CIDR:           info.CIDR,
		Netmask:        info.Netmask,
		WildcardMask:   wildcardIP,
		FirstHost:      info.FirstHost,
		LastHost:       info.LastHost,
		Broadcast:      info.BroadcastAddress,
		WastedHosts:    wastedHosts,
		Efficiency:     efficiency,
		ClassHint:      classHint,
	}, nil
}

