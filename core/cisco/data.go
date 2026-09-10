package cisco

import (
	"strings"
)

// GetAllCommands returns the comprehensive Cisco Packet Tracer command catalog
func GetAllCommands() []CiscoCommand {
	return []CiscoCommand{
		// ====================================================================
		// 1. KONFIGURASI DASAR & KEAMANAN
		// ====================================================================
		{
			ID:          "basic-hostname-banner",
			Title:       "Ganti Nama Perangkat (Hostname) & Pesan Peringatan (Banner MOTD)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryBasic,
			Description: "Mengubah identitas nama router/switch agar mudah dikenali di topologi dan menambahkan banner peringatan saat ada yang mengakses CLI.",
			Parameters: []Parameter{
				{Key: "HOSTNAME", Label: "Nama Hostname", DefaultValue: "R1-Pusat", Placeholder: "cth: R1-Kantor / SW1-Lab"},
				{Key: "BANNER", Label: "Pesan Peringatan (MOTD)", DefaultValue: "PERINGATAN: HANYA PERSONEL RESMI YANG DIIZINKAN!\nAKSES ILEGAL AKAN DITINDAK SECARA HUKUM.", Placeholder: "Teks banner sambutan"},
			},
			IPExample: "Topologi: {{HOSTNAME}} (Router 1941) & SW1-Lantai1 (Switch 2960)",
			Commands: `enable
configure terminal
hostname {{HOSTNAME}}
banner motd #
==================================================
 {{BANNER}}
==================================================
#
exit`,
			Verification:        "Ketik 'exit' hingga kembali ke prompt awal. Saat menekan Enter, pesan Banner MOTD akan muncul dan nama prompt berubah menjadi '{{HOSTNAME}}>'.",
			TroubleshootingTips: "Karakter pembatas (#) di awal dan akhir banner harus sama. Jangan gunakan karakter yang ada di dalam isi pesan.",
		},
		{
			ID:          "basic-passwords-security",
			Title:       "Mengamankan Akses Privileged EXEC & Console Password",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryBasic,
			Description: "Mengunci hak akses 'enable' dengan password terenkripsi MD5 (secret) dan memberi password saat login lewat kabel Console.",
			Parameters: []Parameter{
				{Key: "SECRET_PASS", Label: "Password Enable Secret", DefaultValue: "class", Placeholder: "Password rahasia privilege"},
				{Key: "CONSOLE_PASS", Label: "Password Console Login", DefaultValue: "cisco", Placeholder: "Password login console"},
			},
			IPExample: "Password Privilege: {{SECRET_PASS}} | Password Console: {{CONSOLE_PASS}}",
			Commands: `configure terminal
enable secret {{SECRET_PASS}}
line console 0
password {{CONSOLE_PASS}}
login
logging synchronous
exec-timeout 10 0
exit
service password-encryption
exit`,
			Verification:        "Ketik 'show running-config'. Pastikan 'enable secret' berformat hash acak (bukan teks polos) dan baris 'line con 0' memiliki password terenkripsi.",
			TroubleshootingTips: "Gunakan 'logging synchronous' agar ketikan CLI Anda tidak terpotong saat router memunculkan pesan status/log otomatis.",
		},
		{
			ID:          "basic-ssh-remote",
			Title:       "Konfigurasi Remote Akses Aman (SSHv2 & Local User)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryBasic,
			Description: "Mengaktifkan SSH agar perangkat Cisco dapat diremote secara terenkripsi dari Command Prompt PC di Packet Tracer tanpa Telnet polos.",
			Parameters: []Parameter{
				{Key: "DOMAIN", Label: "Domain Name", DefaultValue: "labcisco.local", Placeholder: "cth: lab.local"},
				{Key: "USER", Label: "Username Admin", DefaultValue: "admin", Placeholder: "cth: admin"},
				{Key: "PASS", Label: "Password Akun", DefaultValue: "cisco123", Placeholder: "cth: admin123"},
				{Key: "IP", Label: "Alamat IP Router", DefaultValue: "192.168.1.1", Placeholder: "cth: 192.168.1.1"},
			},
			IPExample: "IP Router: {{IP}}/24 | Domain: {{DOMAIN}} | User: {{USER}}",
			Commands: `configure terminal
ip domain-name {{DOMAIN}}
crypto key generate rsa
1024
ip ssh version 2
username {{USER}} privilege 15 secret {{PASS}}
line vty 0 4
transport input ssh
login local
exit`,
			Verification:        "Buka PC di Packet Tracer > Command Prompt > ketik: ssh -l {{USER}} {{IP}}, lalu masukkan password. Jika berhasil masuk prompt, SSH sukses.",
			TroubleshootingTips: "RSA Key minimal 1024-bit untuk mendukung SSHv2. Pastikan hostname dan domain-name sudah diset sebelum 'crypto key generate rsa'.",
		},

		// ====================================================================
		// 2. IP ADDRESS & INTERFACE
		// ====================================================================
		{
			ID:          "interface-ipv4-assign",
			Title:       "Setting Alamat IP Interface GigabitEthernet & Mengaktifkan Port",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryInterface,
			Description: "Mengonfigurasi IPv4 dan Subnet Mask pada port LAN Router serta mengaktifkan port yang secara default mati (shutdown).",
			Parameters: []Parameter{
				{Key: "IFACE", Label: "Nama Port Interface", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "IP", Label: "Alamat IPv4", DefaultValue: "192.168.10.1", Placeholder: "cth: 192.168.10.1"},
				{Key: "NETMASK", Label: "Subnet Mask", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "DESC", Label: "Deskripsi Port", DefaultValue: "KONEKSI_KE_LAN_KANTOR", Placeholder: "cth: LINK_LAN_UTAMA"},
			},
			IPExample: "Interface {{IFACE}}: {{IP}} (Subnet Mask: {{NETMASK}})",
			Commands: `configure terminal
interface {{IFACE}}
description {{DESC}}
ip address {{IP}} {{NETMASK}}
no shutdown
exit`,
			Verification:        "Ketik 'show ip interface brief'. Pastikan port berstatus 'Status: up' dan 'Protocol: up' serta lampu link di Packet Tracer berubah hijau.",
			TroubleshootingTips: "Di Router Cisco, SEMUA interface secara default berstatus 'administratively down'. Wajib ketik 'no shutdown' agar port aktif!",
		},
		{
			ID:          "interface-serial-clockrate",
			Title:       "Setting Interface Serial WAN (Kabel Serial DCE + Clock Rate)",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryInterface,
			Description: "Menghubungkan dua Router lewat kabel Serial WAN. Ujung kabel bertanda jam (DCE) memerlukan pengaturan clock rate.",
			Parameters: []Parameter{
				{Key: "IFACE", Label: "Interface Serial", DefaultValue: "Serial 0/0/0", Placeholder: "cth: Serial 0/0/0"},
				{Key: "IP", Label: "Alamat IP (DCE)", DefaultValue: "10.10.10.1", Placeholder: "cth: 10.10.10.1"},
				{Key: "NETMASK", Label: "Subnet Mask WAN", DefaultValue: "255.255.255.252", Placeholder: "cth: 255.255.255.252"},
				{Key: "CLOCK", Label: "Clock Rate (DCE)", DefaultValue: "64000", Placeholder: "cth: 64000 / 128000"},
			},
			IPExample: "R1 ({{IFACE}} DCE): {{IP}}/30 | Serial Subnet Mask: {{NETMASK}}",
			Commands: `configure terminal
interface {{IFACE}}
description LINK_WAN_KE_ROUTER2
ip address {{IP}} {{NETMASK}}
clock rate {{CLOCK}}
no shutdown
exit`,
			Verification:        "Ketik 'show controllers {{IFACE}}'. Pastikan terdeteksi 'DCE' dengan clock rate {{CLOCK}}, lalu ping IP router lawan.",
			TroubleshootingTips: "Gunakan kabel bertanda jam (Serial DCE) di Packet Tracer. Hanya sisi DCE yang butuh 'clock rate', sisi DTE tidak perlu.",
		},
		{
			ID:          "interface-loopback",
			Title:       "Membuat Interface Virtual (Loopback) untuk Uji Coba / Router-ID",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryInterface,
			Description: "Interface virtual yang selalu UP dan tidak pernah down secara fisik. Sangat berguna untuk ID Router OSPF atau target uji ping.",
			Parameters: []Parameter{
				{Key: "LOOP_NUM", Label: "Nomor Loopback", DefaultValue: "0", Placeholder: "cth: 0"},
				{Key: "IP", Label: "Alamat IP Loopback", DefaultValue: "1.1.1.1", Placeholder: "cth: 1.1.1.1"},
				{Key: "NETMASK", Label: "Subnet Mask Host", DefaultValue: "255.255.255.255", Placeholder: "cth: 255.255.255.255"},
			},
			IPExample: "Loopback {{LOOP_NUM}}: {{IP}} (Subnet Mask: {{NETMASK}})",
			Commands: `configure terminal
interface Loopback {{LOOP_NUM}}
description SIMULASI_SERVER_INTERNET
ip address {{IP}} {{NETMASK}}
exit`,
			Verification:        "Ketik 'show ip interface brief'. Interface Loopback{{LOOP_NUM}} langsung berstatus 'up/up' tanpa perlu perintah 'no shutdown'.",
			TroubleshootingTips: "Subnet mask 255.255.255.255 (/32) menghemat IP karena interface loopback hanya memerlukan 1 IP unik.",
		},

		// ====================================================================
		// 3. VLAN & SWITCHING
		// ====================================================================
		{
			ID:          "vlan-create-assign",
			Title:       "Membuat VLAN & Memasukkan Port Access ke VLAN",
			Device:      DeviceSwitchL2,
			Mode:        ModeGlobalConfig,
			Category:    CategoryVLAN,
			Description: "Memisahkan broadcast domain pada Switch ke dalam segmen VLAN berbeda (misal: VLAN Guru & Siswa).",
			Parameters: []Parameter{
				{Key: "VLAN1_ID", Label: "Nomor ID VLAN 1", DefaultValue: "10", Placeholder: "cth: 10"},
				{Key: "VLAN1_NAME", Label: "Nama VLAN 1", DefaultValue: "GURU", Placeholder: "cth: GURU"},
				{Key: "PORT1_RANGE", Label: "Rentang Port VLAN 1", DefaultValue: "FastEthernet 0/1 - 10", Placeholder: "cth: FastEthernet 0/1 - 10"},
				{Key: "VLAN2_ID", Label: "Nomor ID VLAN 2", DefaultValue: "20", Placeholder: "cth: 20"},
				{Key: "VLAN2_NAME", Label: "Nama VLAN 2", DefaultValue: "SISWA", Placeholder: "cth: SISWA"},
				{Key: "PORT2_RANGE", Label: "Rentang Port VLAN 2", DefaultValue: "FastEthernet 0/11 - 20", Placeholder: "cth: FastEthernet 0/11 - 20"},
			},
			IPExample: "VLAN {{VLAN1_ID}} ({{VLAN1_NAME}}): {{PORT1_RANGE}} | VLAN {{VLAN2_ID}} ({{VLAN2_NAME}}): {{PORT2_RANGE}}",
			Commands: `configure terminal
vlan {{VLAN1_ID}}
name {{VLAN1_NAME}}
vlan {{VLAN2_ID}}
name {{VLAN2_NAME}}
exit

interface range {{PORT1_RANGE}}
switchport mode access
switchport access vlan {{VLAN1_ID}}
exit

interface range {{PORT2_RANGE}}
switchport mode access
switchport access vlan {{VLAN2_ID}}
exit`,
			Verification:        "Ketik 'show vlan brief'. Pastikan VLAN {{VLAN1_ID}} ({{VLAN1_NAME}}) dan VLAN {{VLAN2_ID}} ({{VLAN2_NAME}}) aktif dengan port sesuai.",
			TroubleshootingTips: "PC di VLAN {{VLAN1_ID}} TIDAK AKAN BISA saling ping dengan PC di VLAN {{VLAN2_ID}} sebelum dikonfigurasi Router-on-a-Stick (Inter-VLAN Routing).",
		},
		{
			ID:          "vlan-trunk-port",
			Title:       "Mengonfigurasi Port Trunk Antar-Switch atau ke Router",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategoryVLAN,
			Description: "Membuka jalur Trunk pada port penghubung agar dapat melewatkan multi-VLAN sekaligus antar switch atau menuju Router.",
			Parameters: []Parameter{
				{Key: "IFACE", Label: "Port Interface Trunk", DefaultValue: "GigabitEthernet 0/1", Placeholder: "cth: GigabitEthernet 0/1"},
				{Key: "ALLOWED_VLANS", Label: "VLAN yang Diizinkan", DefaultValue: "10,20", Placeholder: "cth: 10,20 / all"},
			},
			IPExample: "Port {{IFACE}} terhubung ke Switch tetangga atau Router (VLAN: {{ALLOWED_VLANS}})",
			Commands: `configure terminal
interface {{IFACE}}
description TRUNK_KE_ROUTER_ATAU_SWITCH2
switchport mode trunk
switchport trunk allowed vlan {{ALLOWED_VLANS}}
exit`,
			Verification:        "Ketik 'show interfaces trunk'. Pastikan port berstatus 'Mode: on / Status: trunking' dengan encapsulation 802.1q.",
			TroubleshootingTips: "Jika menghubungkan Switch ke Switch lain, pastikan kedua ujung port diatur sebagai trunk mode.",
		},
		{
			ID:          "vlan-svi-management",
			Title:       "IP Manajemen Switch (Switch Virtual Interface - SVI VLAN 1)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategoryVLAN,
			Description: "Memberikan alamat IP pada Switch agar Switch dapat diping dan diremote (SSH/Telnet) dari jaringan komputer.",
			Parameters: []Parameter{
				{Key: "VLAN_ID", Label: "Nomor VLAN Management", DefaultValue: "1", Placeholder: "cth: 1 / 99"},
				{Key: "IP", Label: "Alamat IP Switch", DefaultValue: "192.168.1.2", Placeholder: "cth: 192.168.1.2"},
				{Key: "NETMASK", Label: "Subnet Mask", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "GATEWAY", Label: "Default Gateway", DefaultValue: "192.168.1.1", Placeholder: "cth: 192.168.1.1"},
			},
			IPExample: "IP Switch SVI (VLAN {{VLAN_ID}}): {{IP}} | Default Gateway: {{GATEWAY}}",
			Commands: `configure terminal
interface vlan {{VLAN_ID}}
ip address {{IP}} {{NETMASK}}
no shutdown
exit
ip default-gateway {{GATEWAY}}
exit`,
			Verification:        "Buka PC klien di jaringan, buka Command Prompt lalu ketik 'ping {{IP}}'. Pastikan menerima reply.",
			TroubleshootingTips: "'ip default-gateway' sangat penting agar Switch bisa dijangkau dari luar subnet lokalnya.",
		},

		// ====================================================================
		// 4. INTER-VLAN ROUTING (ROAS & MULTILAYER SWITCH)
		// ====================================================================
		{
			ID:          "roas-router-on-a-stick",
			Title:       "Inter-VLAN Routing: Router-on-a-Stick (ROAS Sub-Interface)",
			Device:      DeviceRouter,
			Mode:        ModeSubInterface,
			Category:    CategoryVLAN,
			Description: "Menghubungkan komunikasi antar VLAN yang berbeda melalui satu kabel fisik Router menggunakan sub-interface dan enkapsulasi 802.1Q.",
			Parameters: []Parameter{
				{Key: "PHY_IFACE", Label: "Port Fisik Router", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "VLAN1_ID", Label: "ID VLAN 1", DefaultValue: "10", Placeholder: "cth: 10"},
				{Key: "VLAN1_IP", Label: "IP Gateway VLAN 1", DefaultValue: "192.168.10.1", Placeholder: "cth: 192.168.10.1"},
				{Key: "VLAN1_MASK", Label: "Netmask VLAN 1", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "VLAN2_ID", Label: "ID VLAN 2", DefaultValue: "20", Placeholder: "cth: 20"},
				{Key: "VLAN2_IP", Label: "IP Gateway VLAN 2", DefaultValue: "192.168.20.1", Placeholder: "cth: 192.168.20.1"},
				{Key: "VLAN2_MASK", Label: "Netmask VLAN 2", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
			},
			IPExample: "VLAN {{VLAN1_ID}} Gateway: {{VLAN1_IP}} | VLAN {{VLAN2_ID}} Gateway: {{VLAN2_IP}} (Port Fisik: {{PHY_IFACE}})",
			Commands: `configure terminal
interface {{PHY_IFACE}}
no shutdown
exit

interface {{PHY_IFACE}}.{{VLAN1_ID}}
encapsulation dot1Q {{VLAN1_ID}}
ip address {{VLAN1_IP}} {{VLAN1_MASK}}
exit

interface {{PHY_IFACE}}.{{VLAN2_ID}}
encapsulation dot1Q {{VLAN2_ID}}
ip address {{VLAN2_IP}} {{VLAN2_MASK}}
exit`,
			Verification:        "Ketik 'show ip interface brief' di Router. Pastikan sub-interface berstatus UP, lalu uji ping dari PC VLAN {{VLAN1_ID}} ke PC VLAN {{VLAN2_ID}}.",
			TroubleshootingTips: "Wajib ketik 'encapsulation dot1Q <nomor_vlan>' TERLEBIH DAHULU sebelum memberi IP address di sub-interface, dan port Switch lawan wajib berstatus TRUNK!",
		},
		{
			ID:          "l3-multilayer-switch-routing",
			Title:       "Inter-VLAN Routing pada Switch Layer 3 (Multilayer Switch 3560)",
			Device:      DeviceSwitchL3,
			Mode:        ModeGlobalConfig,
			Category:    CategoryVLAN,
			Description: "Mengaktifkan fungsi routing pada Switch Layer 3 agar switching dan routing antar VLAN berjalan dengan kecepatan kabel tanpa router eksternal.",
			Parameters: []Parameter{
				{Key: "VLAN1_ID", Label: "ID VLAN 1", DefaultValue: "10", Placeholder: "cth: 10"},
				{Key: "VLAN1_NAME", Label: "Nama VLAN 1", DefaultValue: "KANTOR", Placeholder: "cth: KANTOR"},
				{Key: "VLAN1_IP", Label: "IP SVI VLAN 1", DefaultValue: "192.168.10.1", Placeholder: "cth: 192.168.10.1"},
				{Key: "VLAN1_MASK", Label: "Netmask VLAN 1", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "VLAN2_ID", Label: "ID VLAN 2", DefaultValue: "20", Placeholder: "cth: 20"},
				{Key: "VLAN2_NAME", Label: "Nama VLAN 2", DefaultValue: "LAB", Placeholder: "cth: LAB"},
				{Key: "VLAN2_IP", Label: "IP SVI VLAN 2", DefaultValue: "192.168.20.1", Placeholder: "cth: 192.168.20.1"},
				{Key: "VLAN2_MASK", Label: "Netmask VLAN 2", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
			},
			IPExample: "VLAN {{VLAN1_ID}} ({{VLAN1_NAME}}): {{VLAN1_IP}} | VLAN {{VLAN2_ID}} ({{VLAN2_NAME}}): {{VLAN2_IP}}",
			Commands: `configure terminal
ip routing

vlan {{VLAN1_ID}}
name {{VLAN1_NAME}}
vlan {{VLAN2_ID}}
name {{VLAN2_NAME}}

interface vlan {{VLAN1_ID}}
ip address {{VLAN1_IP}} {{VLAN1_MASK}}
no shutdown
exit

interface vlan {{VLAN2_ID}}
ip address {{VLAN2_IP}} {{VLAN2_MASK}}
no shutdown
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan muncul kode 'C {{VLAN1_IP}} is directly connected' dan lakukan ping antar VLAN.",
			TroubleshootingTips: "Perintah 'ip routing' adalah kunci mutlak. Tanpa perintah ini, Switch L3 hanya akan bekerja sebagai Switch L2 biasa.",
		},

		// ====================================================================
		// 5. ROUTING STATIS & DINAMIS
		// ====================================================================
		{
			ID:          "routing-default-route",
			Title:       "Default Route (Gateway of Last Resort ke Internet/ISP)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Meneruskan semua paket data yang tujuannya tidak diketahui di tabel routing ke arah router upstream/ISP.",
			Parameters: []Parameter{
				{Key: "NEXT_HOP", Label: "Alamat IP Next-Hop / Gateway ISP", DefaultValue: "200.100.50.2", Placeholder: "cth: 200.100.50.2"},
			},
			IPExample: "IP Next-Hop ISP: {{NEXT_HOP}}",
			Commands: `configure terminal
ip route 0.0.0.0 0.0.0.0 {{NEXT_HOP}}
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan ada baris 'S* 0.0.0.0/0 [1/0] via {{NEXT_HOP}}' dan 'Gateway of last resort is {{NEXT_HOP}}'.",
			TroubleshootingTips: "0.0.0.0 0.0.0.0 berarti 'segala alamat IP dengan segala subnet mask'. Sangat efisien untuk router cabang yang terhubung ke ISP.",
		},
		{
			ID:          "routing-static-route",
			Title:       "Routing Statis Spesifik (Static Route Menuju Jaringan Tertentu)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Menentukan jalur manual menuju subnet jaringan tertentu di belakang router tetangga.",
			Parameters: []Parameter{
				{Key: "DEST_NET", Label: "Subnet Network Tujuan", DefaultValue: "192.168.20.0", Placeholder: "cth: 192.168.20.0"},
				{Key: "DEST_MASK", Label: "Subnet Mask Tujuan", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "NEXT_HOP", Label: "IP Next-Hop Router Lawan", DefaultValue: "10.10.10.2", Placeholder: "cth: 10.10.10.2"},
			},
			IPExample: "Tujuan: {{DEST_NET}} (Mask: {{DEST_MASK}}) | Lewat Next-Hop: {{NEXT_HOP}}",
			Commands: `configure terminal
ip route {{DEST_NET}} {{DEST_MASK}} {{NEXT_HOP}}
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan ada kode 'S {{DEST_NET}} via {{NEXT_HOP}}', lalu ping IP di subnet tujuan.",
			TroubleshootingTips: "Ingat prinsip routing: 'Koneksi dua arah'. Pastikan router lawan juga memiliki static route balik (return route) ke jaringan asal!",
		},
		{
			ID:          "routing-ospf-single-area",
			Title:       "Routing Dinamis OSPF Single Area (OSPF Area 0)",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Mengaktifkan protokol routing dinamis OSPF open standard yang cepat konvergen menggunakan Wildcard Mask.",
			Parameters: []Parameter{
				{Key: "PROCESS_ID", Label: "OSPF Process ID", DefaultValue: "1", Placeholder: "cth: 1"},
				{Key: "ROUTER_ID", Label: "Router ID OSPF", DefaultValue: "1.1.1.1", Placeholder: "cth: 1.1.1.1"},
				{Key: "NET1", Label: "Network LAN 1", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "WILDCARD1", Label: "Wildcard Mask LAN 1", DefaultValue: "0.0.0.255", Placeholder: "cth: 0.0.0.255"},
				{Key: "NET2", Label: "Network Link WAN 2", DefaultValue: "10.10.10.0", Placeholder: "cth: 10.10.10.0"},
				{Key: "WILDCARD2", Label: "Wildcard Mask WAN 2", DefaultValue: "0.0.0.3", Placeholder: "cth: 0.0.0.3"},
				{Key: "AREA", Label: "Nomor Area OSPF", DefaultValue: "0", Placeholder: "cth: 0 (Backbone Area)"},
			},
			IPExample: "Router-ID: {{ROUTER_ID}} | Net1: {{NET1}} (Wildcard: {{WILDCARD1}}) | Area: {{AREA}}",
			Commands: `configure terminal
router ospf {{PROCESS_ID}}
router-id {{ROUTER_ID}}
network {{NET1}} {{WILDCARD1}} area {{AREA}}
network {{NET2}} {{WILDCARD2}} area {{AREA}}
passive-interface GigabitEthernet 0/0
exit`,
			Verification:        "Ketik 'show ip ospf neighbor'. Pastikan status tetangga 'FULL/BDR' atau 'FULL/DROTHER'. Ketik 'show ip route' dan pastikan rute ditandai kode 'O'.",
			TroubleshootingTips: "Wildcard mask adalah kebalikan dari subnet mask (255.255.255.255 dikurangi Subnet Mask). Contoh /24 = 0.0.0.255, /30 = 0.0.0.3.",
		},
		{
			ID:          "routing-rip-version2",
			Title:       "Routing Dinamis RIP Version 2",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Protokol routing distance vector berbasis Hop Count yang mudah dikonfigurasi untuk jaringan skala kecil.",
			Parameters: []Parameter{
				{Key: "NET1", Label: "Network Lokal 1", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "NET2", Label: "Network Link WAN 2", DefaultValue: "10.0.0.0", Placeholder: "cth: 10.0.0.0"},
			},
			IPExample: "Jaringan lokal: {{NET1}} dan link antar router: {{NET2}}",
			Commands: `configure terminal
router rip
version 2
no auto-summary
network {{NET1}}
network {{NET2}}
exit`,
			Verification:        "Ketik 'show ip route'. Rute yang dipelajari dari RIP akan ditandai dengan huruf 'R'.",
			TroubleshootingTips: "Selalu ketik 'version 2' dan 'no auto-summary' agar RIP mendukung subnetting VLSM/CIDR dan tidak merangkum IP ke classful.",
		},

		// ====================================================================
		// 6. DHCP & LAYANAN JARINGAN
		// ====================================================================
		{
			ID:          "dhcp-server-router",
			Title:       "Konfigurasi DHCP Server di Router Cisco untuk PC Klien",
			Device:      DeviceRouter,
			Mode:        ModeDHCPConfig,
			Category:    CategoryServices,
			Description: "Memberikan alamat IP, Subnet Mask, Gateway, dan DNS secara otomatis kepada komputer klien di jaringan LAN.",
			Parameters: []Parameter{
				{Key: "POOL_NAME", Label: "Nama Pool DHCP", DefaultValue: "POOL_LAN", Placeholder: "cth: POOL_LAN"},
				{Key: "NETWORK", Label: "Network Subnet Klien", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "NETMASK", Label: "Subnet Mask Klien", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "GATEWAY", Label: "Default Gateway Router", DefaultValue: "192.168.1.1", Placeholder: "cth: 192.168.1.1"},
				{Key: "DNS", Label: "DNS Server", DefaultValue: "8.8.8.8", Placeholder: "cth: 8.8.8.8"},
				{Key: "EXCLUDE_START", Label: "Excluded IP Mulai", DefaultValue: "192.168.1.1", Placeholder: "cth: 192.168.1.1"},
				{Key: "EXCLUDE_END", Label: "Excluded IP Selesai", DefaultValue: "192.168.1.10", Placeholder: "cth: 192.168.1.10"},
			},
			IPExample: "Pool: {{POOL_NAME}} | Network: {{NETWORK}} | IP Exclude: {{EXCLUDE_START}} s/d {{EXCLUDE_END}}",
			Commands: `configure terminal
ip dhcp excluded-address {{EXCLUDE_START}} {{EXCLUDE_END}}
ip dhcp pool {{POOL_NAME}}
network {{NETWORK}} {{NETMASK}}
default-router {{GATEWAY}}
dns-server {{DNS}}
exit`,
			Verification:        "Buka PC di Packet Tracer > Desktop > IP Configuration > Pilih tombol 'DHCP'. Pastikan IP terisi otomatis dan muncul pesan 'DHCP request successful'.",
			TroubleshootingTips: "'excluded-address' wajib diset agar IP Gateway router ({{GATEWAY}}) tidak bentrok diberikan ke PC klien!",
		},
		{
			ID:          "dhcp-relay-agent",
			Title:       "DHCP Relay Agent (ip helper-address) Lintas Router",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryServices,
			Description: "Meneruskan pesan broadcast DHCP Discover dari PC di satu LAN menuju Server DHCP terpusat di jaringan lain.",
			Parameters: []Parameter{
				{Key: "CLIENT_IFACE", Label: "Interface Menghadap Klien", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "DHCP_SERVER", Label: "Alamat IP Server DHCP", DefaultValue: "10.10.10.100", Placeholder: "cth: 10.10.10.100"},
			},
			IPExample: "Interface Klien: {{CLIENT_IFACE}} | Alamat IP Server DHCP Terpusat: {{DHCP_SERVER}}",
			Commands: `configure terminal
interface {{CLIENT_IFACE}}
ip helper-address {{DHCP_SERVER}}
exit`,
			Verification:        "Ubah konfigurasi IP PC di sisi klien dari Static ke DHCP. PC akan berhasil memperoleh IP dari Server DHCP meskipun berbeda subnet.",
			TroubleshootingTips: "Perintah 'ip helper-address' diketik pada interface router yang MENGHADAP KE KLIEN (bukan yang menghadap ke server).",
		},

		// ====================================================================
		// 7. KEAMANAN PORT & ACL (ACCESS CONTROL LIST)
		// ====================================================================
		{
			ID:          "security-port-security",
			Title:       "Port Security pada Switch (Mengunci MAC Address PC)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategorySecurity,
			Description: "Mencegah perangkat liar dicolokkan ke port switch. Port akan mati otomatis (shutdown) jika terdeteksi MAC address yang tidak terdaftar.",
			Parameters: []Parameter{
				{Key: "IFACE", Label: "Port FastEthernet", DefaultValue: "FastEthernet 0/1", Placeholder: "cth: FastEthernet 0/1"},
				{Key: "MAX_MAC", Label: "Jumlah Maksimal MAC", DefaultValue: "1", Placeholder: "cth: 1"},
				{Key: "VIOLATION", Label: "Aksi Pelanggaran", DefaultValue: "shutdown", Placeholder: "shutdown / protect / restrict"},
			},
			IPExample: "Port {{IFACE}} hanya boleh dipakai oleh {{MAX_MAC}} PC (Violation: {{VIOLATION}})",
			Commands: `configure terminal
interface {{IFACE}}
switchport mode access
switchport port-security
switchport port-security maximum {{MAX_MAC}}
switchport port-security mac-address sticky
switchport port-security violation {{VIOLATION}}
exit`,
			Verification:        "Ketik 'show port-security interface {{IFACE}}'. Pastikan 'Port Security: Enabled' dan status 'Secure-up'. Coba colokkan PC lain, lampu port akan langsung merah.",
			TroubleshootingTips: "Jika port terkunci merah karena violation, buka port kembali dengan masuk ke interface lalu ketik: 'shutdown' kemudian 'no shutdown'.",
		},
		{
			ID:          "security-standard-acl",
			Title:       "Standard Access Control List (Blokir Berdasarkan IP Sumber)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategorySecurity,
			Description: "Memfilter paket data hanya berdasarkan IP Address sumber (source IP). Nomor rentang Standard ACL: 1 - 99.",
			Parameters: []Parameter{
				{Key: "ACL_NUM", Label: "Nomor ACL (1 - 99)", DefaultValue: "10", Placeholder: "cth: 10"},
				{Key: "BLOCKED_HOST", Label: "IP Host yang Diblokir", DefaultValue: "192.168.1.50", Placeholder: "cth: 192.168.1.50"},
				{Key: "APPLY_IFACE", Label: "Interface Penempatan ACL", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "DIRECTION", Label: "Arah Filter (in / out)", DefaultValue: "in", Placeholder: "in atau out"},
			},
			IPExample: "Blokir PC {{BLOCKED_HOST}} di port {{APPLY_IFACE}} (Arah: {{DIRECTION}})",
			Commands: `configure terminal
access-list {{ACL_NUM}} deny host {{BLOCKED_HOST}}
access-list {{ACL_NUM}} permit any

interface {{APPLY_IFACE}}
ip access-group {{ACL_NUM}} {{DIRECTION}}
exit`,
			Verification:        "Uji ping dari PC {{BLOCKED_HOST}}, pastikan hasilnya 'Destination Host Unreachable'. Uji ping dari PC lain, pastikan berhasil reply.",
			TroubleshootingTips: "Ingat kaidah Cisco: 'Implicit Deny Any' ada di akhir semua ACL. Selalu tambahkan 'permit any' di akhir jika tidak ingin semua traffic terblokir!",
		},
		{
			ID:          "security-extended-acl",
			Title:       "Extended Access Control List (Blokir Port/Protokol: Web/Ping)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategorySecurity,
			Description: "Memfilter paket data secara presisi berdasarkan IP asal, IP tujuan, protokol (TCP/UDP/ICMP), dan nomor port (cth: HTTP port 80).",
			Parameters: []Parameter{
				{Key: "ACL_NUM", Label: "Nomor ACL (100 - 199)", DefaultValue: "100", Placeholder: "cth: 100"},
				{Key: "SRC_NET", Label: "Source Network Asal", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "SRC_WILDCARD", Label: "Wildcard Source", DefaultValue: "0.0.0.255", Placeholder: "cth: 0.0.0.255"},
				{Key: "DEST_HOST", Label: "IP Host Tujuan", DefaultValue: "10.10.10.10", Placeholder: "cth: 10.10.10.10"},
				{Key: "PORT", Label: "Nomor Port Layanan", DefaultValue: "80", Placeholder: "cth: 80 (HTTP) / 443 (HTTPS)"},
				{Key: "APPLY_IFACE", Label: "Interface Penempatan ACL", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
			},
			IPExample: "Izinkan port {{PORT}} ke Server {{DEST_HOST}}, blokir ICMP Ping dari {{SRC_NET}}",
			Commands: `configure terminal
access-list {{ACL_NUM}} permit tcp {{SRC_NET}} {{SRC_WILDCARD}} host {{DEST_HOST}} eq {{PORT}}
access-list {{ACL_NUM}} deny icmp {{SRC_NET}} {{SRC_WILDCARD}} host {{DEST_HOST}}
access-list {{ACL_NUM}} permit ip any any

interface {{APPLY_IFACE}}
ip access-group {{ACL_NUM}} in
exit`,
			Verification:        "Ketik 'show access-lists'. Coba buka Web Browser di PC ke {{DEST_HOST}} (sukses) lalu coba ping {{DEST_HOST}} di CMD (Request timed out). Hit counter ACL akan bertambah.",
			TroubleshootingTips: "Extended ACL ditempatkan sedekat mungkin dengan sumber traffic (traffic source).",
		},

		// ====================================================================
		// 8. NAT & PAT INTERNET SHARING
		// ====================================================================
		{
			ID:          "nat-pat-overload",
			Title:       "PAT (NAT Overload) untuk Sharing Koneksi Internet Seluruh LAN",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryNAT,
			Description: "Menerjemahkan banyak alamat IP Private di LAN ke satu alamat IP Public milik port Serial/Gigabit WAN menuju ISP.",
			Parameters: []Parameter{
				{Key: "INSIDE_IFACE", Label: "Inside Interface (LAN)", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "OUTSIDE_IFACE", Label: "Outside Interface (WAN)", DefaultValue: "Serial 0/0/0", Placeholder: "cth: Serial 0/0/0"},
				{Key: "ACL_NUM", Label: "Nomor Standard ACL", DefaultValue: "1", Placeholder: "cth: 1"},
				{Key: "LAN_NET", Label: "Network Subnet LAN", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "LAN_WILDCARD", Label: "Wildcard Mask LAN", DefaultValue: "0.0.0.255", Placeholder: "cth: 0.0.0.255"},
			},
			IPExample: "LAN Inside: {{LAN_NET}}/24 ({{INSIDE_IFACE}}) | WAN Outside: {{OUTSIDE_IFACE}}",
			Commands: `configure terminal
interface {{INSIDE_IFACE}}
ip nat inside
exit

interface {{OUTSIDE_IFACE}}
ip nat outside
exit

access-list {{ACL_NUM}} permit {{LAN_NET}} {{LAN_WILDCARD}}

ip nat inside source list {{ACL_NUM}} interface {{OUTSIDE_IFACE}} overload
exit`,
			Verification:        "Dari PC LAN, lakukan ping ke Server Internet (misal 8.8.8.8). Di Router, ketik 'show ip nat translations'. Pastikan muncul entri terjemahan IP Inside Global & Local.",
			TroubleshootingTips: "Jangan lupa menandai 'ip nat inside' pada interface LAN dan 'ip nat outside' pada interface WAN. Kata kunci 'overload' adalah penentu fitur PAT.",
		},
		{
			ID:          "nat-static-one-to-one",
			Title:       "Static NAT (1-to-1) untuk Membuka Server Lokal ke Publik",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryNAT,
			Description: "Memetakan satu IP Private Server di LAN secara permanen ke satu IP Public agar dapat diakses oleh pengguna dari internet.",
			Parameters: []Parameter{
				{Key: "INSIDE_IFACE", Label: "Inside Interface (LAN)", DefaultValue: "GigabitEthernet 0/0", Placeholder: "cth: GigabitEthernet 0/0"},
				{Key: "OUTSIDE_IFACE", Label: "Outside Interface (WAN)", DefaultValue: "Serial 0/0/0", Placeholder: "cth: Serial 0/0/0"},
				{Key: "PRIVATE_IP", Label: "IP Private Server Lokal", DefaultValue: "192.168.1.100", Placeholder: "cth: 192.168.1.100"},
				{Key: "PUBLIC_IP", Label: "IP Public Server (ISP)", DefaultValue: "200.100.10.5", Placeholder: "cth: 200.100.10.5"},
			},
			IPExample: "IP Server Private: {{PRIVATE_IP}} <==> IP Public Server: {{PUBLIC_IP}}",
			Commands: `configure terminal
interface {{INSIDE_IFACE}}
ip nat inside
exit

interface {{OUTSIDE_IFACE}}
ip nat outside
exit

ip nat inside source static {{PRIVATE_IP}} {{PUBLIC_IP}}
exit`,
			Verification:        "Dari PC di sisi Internet/ISP, lakukan ping atau browsing ke {{PUBLIC_IP}}. Di Router ketik 'show ip nat translations', akan terlihat entri static translation.",
			TroubleshootingTips: "Static NAT menghabiskan 1 IP Public per 1 Server. Cocok untuk Web Server, Mail Server, atau FTP Server sekolah/perusahaan.",
		},

		// ====================================================================
		// 9. TROUBLESHOOTING & SHOW COMMANDS
		// ====================================================================
		{
			ID:          "show-interface-brief",
			Title:       "Melihat Ringkasan Status Semua Port (show ip interface brief)",
			Device:      DeviceRouter,
			Mode:        ModePrivExec,
			Category:    CategoryShowDiag,
			Description: "Perintah paling penting dalam mendiagnosis masalah jaringan Cisco. Menampilkan nama port, IP address, status Layer 1 (Status) dan Layer 2 (Protocol).",
			IPExample:   "Digunakan di Router dan Switch Layer 3",
			Commands: `enable
show ip interface brief`,
			Verification:        "Status ideal port yang bekerja adalah 'Status: up' dan 'Protocol: up'. Jika 'administratively down', lakukan 'no shutdown'. Jika 'down/down', periksa kabel fisik.",
			TroubleshootingTips: "Bisa disingkat: 'sh ip int br'. Sangat sering dipakai saat ujian praktikum untuk mendeteksi port mana yang lupa dinyalakan.",
		},
		{
			ID:          "show-routing-table",
			Title:       "Melihat Tabel Routing Lengkap (show ip route)",
			Device:      DeviceRouter,
			Mode:        ModePrivExec,
			Category:    CategoryShowDiag,
			Description: "Melihat daftar semua jalur jaringan yang diketahui oleh Router beserta kodenya: C (Connected), L (Local), S (Static), O (OSPF), R (RIP).",
			IPExample:   "Digunakan di Router untuk mengecek apakah subnet tujuan sudah terdaftar",
			Commands: `enable
show ip route`,
			Verification:        "Cari apakah subnet IP tujuan ping Anda tercantum di daftar. Jika tidak ada di tabel routing, paket PASTI dibuang (drop/unreachable).",
			TroubleshootingTips: "Bisa disingkat: 'sh ip ro'. Jika menggunakan routing dinamis (OSPF/RIP), tunggu 10-30 detik konvergensi hingga huruf 'O' atau 'R' muncul.",
		},
		{
			ID:          "show-vlan-and-mac",
			Title:       "Melihat Daftar VLAN & Tabel MAC Address Switch",
			Device:      DeviceSwitchL2,
			Mode:        ModePrivExec,
			Category:    CategoryShowDiag,
			Description: "Memverifikasi apakah VLAN sudah terbuat dan port sudah masuk ke VLAN yang tepat, serta melihat pemetaan MAC address perangkat yang terhubung.",
			IPExample:   "Digunakan pada Switch Cisco 2960",
			Commands: `enable
show vlan brief
show mac address-table`,
			Verification:        "Pastikan nama VLAN muncul dan port aktif tidak salah kamar. Pada 'show mac address-table', pastikan port mengenali MAC address PC yang mencolok.",
			TroubleshootingTips: "Bisa disingkat: 'sh vlan br' dan 'sh mac add'. Jika PC tidak muncul di MAC table, lakukan ping dari PC agar switch belajar MAC-nya.",
		},
		{
			ID:          "show-cdp-neighbors",
			Title:       "Mendeteksi Perangkat Tetangga Cisco (Cisco Discovery Protocol)",
			Device:      DeviceRouter,
			Mode:        ModePrivExec,
			Category:    CategoryShowDiag,
			Description: "Mengetahui perangkat Cisco apa saja yang terhubung ke port kita, nama hostname tetangga, model perangkat, dan port lawan tanpa perlu login ke perangkat lawan.",
			IPExample:   "Berguna saat topologi rumit dengan banyak router & switch",
			Commands: `enable
show cdp neighbors
show cdp neighbors detail`,
			Verification:        "Akan muncul tabel Device ID, Local Intrfce, Capability, Platform, dan Port ID lawan.",
			TroubleshootingTips: "CDP adalah protokol proprietary Cisco. Jika lawan adalah PC atau perangkat non-Cisco, ia tidak akan muncul di tabel CDP.",
		},

		// ====================================================================
		// 10. PANDUAN VERIFIKASI TOPOLOGI & PENGUJIAN KERJA
		// ====================================================================
		{
			ID:          "verify-step-by-step-guide",
			Title:       "Checklist 6 Langkah Memastikan Topologi Cisco Packet Tracer Bekerja",
			Device:      DeviceAll,
			Mode:        ModePrivExec,
			Category:    CategoryVerification,
			Description: "Prosedur standar pengujian menyeluruh dari Layer 1 Fisik hingga Layer 7 Aplikasi untuk memastikan jaringan berfungsi 100%.",
			Parameters: []Parameter{
				{Key: "GW_IP", Label: "IP Default Gateway", DefaultValue: "192.168.1.1", Placeholder: "cth: 192.168.1.1"},
				{Key: "NEXT_HOP_IP", Label: "IP Router Next-Hop", DefaultValue: "10.10.10.2", Placeholder: "cth: 10.10.10.2"},
				{Key: "DEST_IP", Label: "IP Server / PC Tujuan", DefaultValue: "10.10.10.10", Placeholder: "cth: 10.10.10.10"},
			},
			IPExample: "Uji Koneksi PC Klien menuju Gateway {{GW_IP}} dan Server {{DEST_IP}}",
			Commands: `# 1. CEK FISIK (Packet Tracer Link Lights):
#    - Lampu Hijau = Link Up (Normal)
#    - Lampu Oranye = Spanning Tree Protocol (STP) sedang listening (tunggu 30 detik atau klik Fast Forward Time)
#    - Lampu Merah = Port mati (Ketik 'no shutdown') atau tipe kabel salah (Straight vs Cross)

# 2. CEK IP PC DI COMMAND PROMPT:
ipconfig /all

# 3. PING DIRI SENDIRI (Cek TCP/IP Stack PC):
ping 127.0.0.1

# 4. PING DEFAULT GATEWAY (Cek koneksi PC ke Port Router Lokal):
ping {{GW_IP}}

# 5. PING ROUTER NEXT-HOP (Cek link antar router WAN):
ping {{NEXT_HOP_IP}}

# 6. PING END-TO-END KE SERVER / PC TUJUAN:
ping {{DEST_IP}}
tracert {{DEST_IP}}`,
			Verification:        "Jika tracert menampilkan seluruh hop tanpa bintang (* * * Request timed out), maka jalur data pulang-pergi (round-trip) sudah 100% sempurna.",
			TroubleshootingTips: "Ping pertama di Packet Tracer sering 'Request timed out' 1 kali lalu reply (!!!!!). Ini normal karena proses ARP (Address Resolution Protocol). Coba ping kedua kali!",
		},
		{
			ID:          "verify-pdu-simulation-mode",
			Title:       "Cara Melacak Paket Menggunakan Mode Simulasi (Add Simple PDU)",
			Device:      DevicePC,
			Mode:        ModePCTerminal,
			Category:    CategoryVerification,
			Description: "Menggunakan fitur visual amplop surat (Simple PDU) di Cisco Packet Tracer untuk melihat di perangkat mana paket mengalami drop atau salah rute.",
			IPExample:   "Klik ikon Surat Tertutup (Add Simple PDU) atau tekan shortcut keyboard 'P'",
			Commands: `# Langkah Melacak dengan Mode Simulasi Packet Tracer:
# 1. Alihkan mode di kanan bawah dari 'Realtime' ke 'Simulation' (Shift + S).
# 2. Klik ikon Amplop Surat (Add Simple PDU) atau tekan tombol 'P'.
# 3. Klik PC Pengirim, lalu klik PC/Server Penerima.
# 4. Klik tombol 'Capture / Forward' (tanda panah kanan dengan garis).
# 5. Amati pergerakan amplop di topologi:
#    - Jika muncul tanda Centang Hijau di amplop = Sukses (Successful).
#    - Jika muncul tanda Silang Merah di Router/Switch = Paket didrop di perangkat tersebut!
# 6. Klik amplop yang gagal untuk melihat In Layers & Out Layers OSI Model dan alasan drop (cth: 'No route to host').`,
			Verification:        "Tabel 'Simulation Panel' di kanan bawah akan menampilkan status 'Successful' pada daftar User Created PDU List.",
			TroubleshootingTips: "Klik tombol 'Fast Forward Time' (Alt + D) 2-3 kali di mode Realtime sebelum beralih ke Simulasi agar STP dan ARP sudah siap.",
		},

		// ====================================================================
		// 11. MANAJEMEN, SIMPAN & RESET
		// ====================================================================
		{
			ID:          "mgmt-save-and-reset",
			Title:       "Menyimpan Konfigurasi Permanen & Cara Mereset Perangkat",
			Device:      DeviceAll,
			Mode:        ModePrivExec,
			Category:    CategoryBasic,
			Description: "Menyimpan konfigurasi dari RAM ke NVRAM agar tidak hilang saat mati listrik/reload, serta cara menghapus bersih konfigurasi lama.",
			IPExample:   "Wajib dilakukan setelah menyelesaikan praktikum topologi",
			Commands: `# CARA 1: MENYIMPAN KONFIGURASI KE NVRAM (PILIH SALAH SATU)
copy running-config startup-config
# ATAU cukup ketik singkat:
write memory

# CARA 2: RESET TOTAL ROUTER (KEMBALI KE PABRIK)
erase startup-config
reload
# Tekan Enter saat ditanya 'Proceed with reload? [confirm]'

# CARA 3: RESET TOTAL SWITCH (HAPUS VLAN DATABASE & CONFIG)
delete vlan.dat
erase startup-config
reload`,
			Verification:        "Saat router menyala kembali setelah reload, muncul prompt 'Would you like to enter the initial configuration dialog? [yes/no]:' tanda reset sukses.",
			TroubleshootingTips: "Di Switch Cisco, menghapus startup-config SAJA TIDAK MENGHAPUS VLAN! Wajib jalankan 'delete vlan.dat' agar VLAN hilang.",
		},
	}
}

// SearchCommands filters Cisco commands based on criteria
func SearchCommands(criteria FilterCriteria) []CiscoCommand {
	all := GetAllCommands()
	q := strings.ToLower(strings.TrimSpace(criteria.Query))

	var results []CiscoCommand
	for _, cmd := range all {
		// 1. Device filter
		if criteria.Device != "" && criteria.Device != DeviceAll {
			if cmd.Device != criteria.Device && cmd.Device != DeviceAll {
				continue
			}
		}

		// 2. Mode filter
		if criteria.Mode != "" && criteria.Mode != ModeAll {
			if cmd.Mode != criteria.Mode {
				continue
			}
		}

		// 3. Category filter
		if criteria.Category != "" && criteria.Category != CategoryAll {
			if cmd.Category != criteria.Category {
				continue
			}
		}

		// 4. Query text search
		if q != "" {
			matchTitle := strings.Contains(strings.ToLower(cmd.Title), q)
			matchDesc := strings.Contains(strings.ToLower(cmd.Description), q)
			matchCmds := strings.Contains(strings.ToLower(cmd.Commands), q)
			matchIP := strings.Contains(strings.ToLower(cmd.IPExample), q)
			matchVerif := strings.Contains(strings.ToLower(cmd.Verification), q)
			matchTips := strings.Contains(strings.ToLower(cmd.TroubleshootingTips), q)

			if !matchTitle && !matchDesc && !matchCmds && !matchIP && !matchVerif && !matchTips {
				continue
			}
		}

		results = append(results, cmd)
	}

	return results
}
