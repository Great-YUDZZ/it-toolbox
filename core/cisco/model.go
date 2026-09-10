package cisco

// DeviceType represents the target hardware in Cisco Packet Tracer
type DeviceType string

const (
	DeviceAll      DeviceType = "Semua Perangkat"
	DeviceRouter   DeviceType = "Router Cisco"
	DeviceSwitchL2 DeviceType = "Switch L2 (2960)"
	DeviceSwitchL3 DeviceType = "Switch L3 Multilayer (3560)"
	DevicePC       DeviceType = "PC / End Device"
)

// CLIMode represents the Cisco IOS prompt/mode level
type CLIMode string

const (
	ModeAll           CLIMode = "Semua Mode"
	ModeUserExec      CLIMode = "User EXEC (Router>)"
	ModePrivExec      CLIMode = "Privileged EXEC (Router#)"
	ModeGlobalConfig  CLIMode = "Global Config (config)#"
	ModeInterface     CLIMode = "Interface Config (config-if)#"
	ModeSubInterface  CLIMode = "Sub-Interface (config-subif)#"
	ModeVLAN          CLIMode = "VLAN Config (config-vlan)#"
	ModeLine          CLIMode = "Line Config (config-line)#"
	ModeRouterConfig  CLIMode = "Router Config (config-router)#"
	ModeDHCPConfig    CLIMode = "DHCP Config (config-dhcp)#"
	ModePCTerminal    CLIMode = "PC Command Prompt (PC>)"
)

// Category represents the configuration objective
type Category string

const (
	CategoryAll          Category = "Semua Kategori"
	CategoryBasic        Category = "Konfigurasi Dasar & Keamanan"
	CategoryInterface    Category = "IP Address & Interface"
	CategoryVLAN         Category = "VLAN & Switching"
	CategoryRouting      Category = "Routing Statis & Dinamis"
	CategoryServices     Category = "DHCP & Layanan Jaringan"
	CategorySecurity     Category = "Keamanan Port & ACL"
	CategoryNAT          Category = "NAT & PAT Internet Sharing"
	CategoryShowDiag     Category = "Troubleshooting (Show Commands)"
	CategoryVerification Category = "Panduan Verifikasi Topologi"
)

// CiscoCommand represents a comprehensive Cisco Packet Tracer configuration recipe
type CiscoCommand struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Device              DeviceType `json:"device"`
	Mode                CLIMode    `json:"mode"`
	Category            Category   `json:"category"`
	Description         string     `json:"description"`
	IPExample           string     `json:"ip_example,omitempty"`
	Commands            string     `json:"commands"`
	Verification        string     `json:"verification"`
	TroubleshootingTips string     `json:"troubleshooting_tips,omitempty"`
}

// FilterCriteria defines search and filter options
type FilterCriteria struct {
	Query    string
	Device   DeviceType
	Mode     CLIMode
	Category Category
}
