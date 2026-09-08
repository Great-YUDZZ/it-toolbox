package calculators

import (
	"testing"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		cidr         string
		wantNetwork  string
		wantBcast    string
		wantFirst    string
		wantLast     string
		wantHosts    int
		wantNetmask  string
		wantErr      bool
	}{
		{
			cidr:        "192.168.1.10/24",
			wantNetwork: "192.168.1.0",
			wantBcast:   "192.168.1.255",
			wantFirst:   "192.168.1.1",
			wantLast:    "192.168.1.254",
			wantHosts:   254,
			wantNetmask: "255.255.255.0",
			wantErr:     false,
		},
		{
			cidr:        "10.0.0.0/30",
			wantNetwork: "10.0.0.0",
			wantBcast:   "10.0.0.3",
			wantFirst:   "10.0.0.1",
			wantLast:    "10.0.0.2",
			wantHosts:   2,
			wantNetmask: "255.255.255.252",
			wantErr:     false,
		},
		{
			cidr:    "invalid-cidr",
			wantErr: true,
		},
		{
			cidr:    "2001:db8::/32",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		got, err := ParseCIDR(tc.cidr)
		if (err != nil) != tc.wantErr {
			t.Errorf("ParseCIDR(%q) error = %v; wantErr %v", tc.cidr, err, tc.wantErr)
			continue
		}
		if !tc.wantErr {
			if got.NetworkAddress != tc.wantNetwork {
				t.Errorf("ParseCIDR(%q).NetworkAddress = %v; want %v", tc.cidr, got.NetworkAddress, tc.wantNetwork)
			}
			if got.BroadcastAddress != tc.wantBcast {
				t.Errorf("ParseCIDR(%q).BroadcastAddress = %v; want %v", tc.cidr, got.BroadcastAddress, tc.wantBcast)
			}
			if got.FirstHost != tc.wantFirst {
				t.Errorf("ParseCIDR(%q).FirstHost = %v; want %v", tc.cidr, got.FirstHost, tc.wantFirst)
			}
			if got.LastHost != tc.wantLast {
				t.Errorf("ParseCIDR(%q).LastHost = %v; want %v", tc.cidr, got.LastHost, tc.wantLast)
			}
			if got.TotalHosts != tc.wantHosts {
				t.Errorf("ParseCIDR(%q).TotalHosts = %v; want %v", tc.cidr, got.TotalHosts, tc.wantHosts)
			}
			if got.Netmask != tc.wantNetmask {
				t.Errorf("ParseCIDR(%q).Netmask = %v; want %v", tc.cidr, got.Netmask, tc.wantNetmask)
			}
		}
	}
}

func TestCalculateVLSM(t *testing.T) {
	// 192.168.1.0/24 with subnets for 100 hosts, 50 hosts, and 20 hosts
	subnets := []int{100, 50, 20}
	results, err := CalculateVLSM("192.168.1.0/24", subnets)
	if err != nil {
		t.Fatalf("CalculateVLSM unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("CalculateVLSM expected 3 subnets, got %d", len(results))
	}

	// 100 hosts -> requires /25 (126 usable hosts)
	if results[0].CIDR != "192.168.1.0/25" {
		t.Errorf("Subnet 0 CIDR = %s, want 192.168.1.0/25", results[0].CIDR)
	}

	// 50 hosts -> requires /26 (62 usable hosts)
	if results[1].CIDR != "192.168.1.128/26" {
		t.Errorf("Subnet 1 CIDR = %s, want 192.168.1.128/26", results[1].CIDR)
	}

	// 20 hosts -> requires /27 (30 usable hosts)
	if results[2].CIDR != "192.168.1.192/27" {
		t.Errorf("Subnet 2 CIDR = %s, want 192.168.1.192/27", results[2].CIDR)
	}

	// Overflow test: trying to allocate 500 hosts in /24
	_, errOverflow := CalculateVLSM("192.168.1.0/24", []int{500})
	if errOverflow == nil {
		t.Errorf("Expected overflow error, got nil")
	}
}

func TestFindSubnetForHosts(t *testing.T) {
	// 50 hosts -> should recommend /26 (62 usable hosts)
	rec, err := FindSubnetForHosts(50, "192.168.1.0")
	if err != nil {
		t.Fatalf("FindSubnetForHosts(50) error: %v", err)
	}

	if rec.Prefix != 26 {
		t.Errorf("FindSubnetForHosts(50).Prefix = %d; want 26", rec.Prefix)
	}
	if rec.AllocatedHosts != 62 {
		t.Errorf("FindSubnetForHosts(50).AllocatedHosts = %d; want 62", rec.AllocatedHosts)
	}
	if rec.Netmask != "255.255.255.192" {
		t.Errorf("FindSubnetForHosts(50).Netmask = %s; want 255.255.255.192", rec.Netmask)
	}
	if rec.WastedHosts != 12 {
		t.Errorf("FindSubnetForHosts(50).WastedHosts = %d; want 12", rec.WastedHosts)
	}

	// 1000 hosts -> should recommend /22 (1022 usable hosts)
	rec1000, err := FindSubnetForHosts(1000, "10.0.0.0")
	if err != nil {
		t.Fatalf("FindSubnetForHosts(1000) error: %v", err)
	}
	if rec1000.Prefix != 22 {
		t.Errorf("FindSubnetForHosts(1000).Prefix = %d; want 22", rec1000.Prefix)
	}

	// Invalid hosts
	_, errInv := FindSubnetForHosts(0, "192.168.1.0")
	if errInv == nil {
		t.Errorf("Expected error for 0 hosts, got nil")
	}
}

