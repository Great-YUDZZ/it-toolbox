package cisco

import (
	"strings"
	"testing"
)

func TestGetAllCommands(t *testing.T) {
	cmds := GetAllCommands()
	if len(cmds) == 0 {
		t.Fatalf("expected non-empty command catalog, got 0")
	}

	// Verify required fields on each command
	for _, c := range cmds {
		if c.ID == "" {
			t.Errorf("command missing ID")
		}
		if c.Title == "" {
			t.Errorf("command %s missing Title", c.ID)
		}
		if c.Commands == "" {
			t.Errorf("command %s missing Commands", c.ID)
		}
		if c.Verification == "" {
			t.Errorf("command %s missing Verification", c.ID)
		}
		if c.Category == "" {
			t.Errorf("command %s missing Category", c.ID)
		}
		if c.Device == "" {
			t.Errorf("command %s missing Device", c.ID)
		}
	}
}

func TestSearchCommandsFiltering(t *testing.T) {
	// 1. Text Query Test
	resOSPF := SearchCommands(FilterCriteria{Query: "ospf"})
	if len(resOSPF) == 0 {
		t.Errorf("expected to find OSPF commands, got 0")
	}

	// 2. Device Filtering Test
	resRouter := SearchCommands(FilterCriteria{Device: DeviceRouter})
	for _, c := range resRouter {
		if c.Device != DeviceRouter && c.Device != DeviceAll {
			t.Errorf("expected Router or DeviceAll, got %s for cmd %s", c.Device, c.ID)
		}
	}

	// 3. Category Filtering Test
	resVLAN := SearchCommands(FilterCriteria{Category: CategoryVLAN})
	if len(resVLAN) == 0 {
		t.Errorf("expected VLAN commands, got 0")
	}
	for _, c := range resVLAN {
		if c.Category != CategoryVLAN {
			t.Errorf("expected CategoryVLAN, got %s for cmd %s", c.Category, c.ID)
		}
	}

	// 4. Combined Filter
	resCombined := SearchCommands(FilterCriteria{
		Device:   DeviceRouter,
		Category: CategoryRouting,
		Query:    "default",
	})
	if len(resCombined) == 0 {
		t.Errorf("expected Default Route command with combined filter, got 0")
	}
}

func TestRenderCommandsWithParameters(t *testing.T) {
	cmds := GetAllCommands()
	var hostCmd *CiscoCommand
	for _, c := range cmds {
		if c.ID == "basic-hostname-banner" {
			hostCmd = &c
			break
		}
	}
	if hostCmd == nil {
		t.Fatalf("basic-hostname-banner command not found")
	}

	// Test 1: Default fallback when values map is nil
	defaultRendered := hostCmd.RenderCommands(nil)
	if !strings.Contains(defaultRendered, "hostname R1-Pusat") {
		t.Errorf("expected default hostname 'R1-Pusat', got: %s", defaultRendered)
	}

	// Test 2: Custom parameter replacement
	customValues := map[string]string{
		"HOSTNAME": "Router-Kantor-Cabang",
		"BANNER":   "SELAMAT DATANG DI JARINGAN CABANG!",
	}
	customRendered := hostCmd.RenderCommands(customValues)
	if !strings.Contains(customRendered, "hostname Router-Kantor-Cabang") {
		t.Errorf("expected custom hostname 'Router-Kantor-Cabang', got: %s", customRendered)
	}
	if !strings.Contains(customRendered, "SELAMAT DATANG DI JARINGAN CABANG!") {
		t.Errorf("expected custom banner message in commands output")
	}

	// Test 3: IP Example rendering with custom values
	customIP := hostCmd.RenderIPExample(customValues)
	if !strings.Contains(customIP, "Router-Kantor-Cabang") {
		t.Errorf("expected custom hostname in IPExample, got: %s", customIP)
	}
}
