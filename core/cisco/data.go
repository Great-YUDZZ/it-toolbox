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
			IPExample:   "Topologi: R1-Kantor-Pusat (Router 1941) & SW1-Lantai1 (Switch 2960)",
			Commands: `enable
configure terminal
hostname R1-Pusat
banner motd #
==================================================
 PERINGATAN: HANYA PERSONEL RESMI YANG DIIZINKAN!
 AKSES ILEGAL AKAN DITINDAK SECARA HUKUM.
==================================================
#
exit`,
			Verification:        "Ketik 'exit' hingga kembali ke prompt awal. Saat menekan Enter, pesan Banner MOTD akan muncul dan nama prompt berubah menjadi 'R1-Pusat>'.",
			TroubleshootingTips: "Karakter pembatas (#) di awal dan akhir banner harus sama. Jangan gunakan karakter yang ada di dalam isi pesan.",
		},
		{
			ID:          "basic-passwords-security",
			Title:       "Mengamankan Akses Privileged EXEC & Console Password",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryBasic,
			Description: "Mengunci hak akses 'enable' dengan password terenkripsi MD5 (secret) dan memberi password saat login lewat kabel Console.",
			IPExample:   "Password Privilege: class | Password Console: cisco",
			Commands: `configure terminal
enable secret class
line console 0
password cisco
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
			IPExample:   "IP Router: 192.168.1.1/24 | Domain: labcisco.local | User: admin",
			Commands: `configure terminal
ip domain-name labcisco.local
crypto key generate rsa
1024
ip ssh version 2
username admin privilege 15 secret cisco123
line vty 0 4
transport input ssh
login local
exit`,
			Verification:        "Buka PC di Packet Tracer > Command Prompt > ketik: ssh -l admin 192.168.1.1, lalu masukkan password 'cisco123'. Jika berhasil masuk prompt 'R1#', SSH sukses.",
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
			IPExample:   "Interface G0/0: 192.168.10.1/24 (Subnet Mask: 255.255.255.0)",
			Commands: `configure terminal
interface GigabitEthernet 0/0
description KONEKSI_KE_LAN_KANTOR
ip address 192.168.10.1 255.255.255.0
no shutdown
exit`,
			Verification:        "Ketik 'show ip interface brief'. Pastikan GigabitEthernet0/0 berstatus 'Status: up' dan 'Protocol: up' serta lampu link di Packet Tracer berubah hijau.",
			TroubleshootingTips: "Di Router Cisco, SEMUA interface secara default berstatus 'administratively down'. Wajib ketik 'no shutdown' agar port aktif!",
		},
		{
			ID:          "interface-serial-clockrate",
			Title:       "Setting Interface Serial WAN (Kabel Serial DCE + Clock Rate)",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryInterface,
			Description: "Menghubungkan dua Router lewat kabel Serial WAN. Ujung kabel bertanda jam (DCE) memerlukan pengaturan clock rate.",
			IPExample:   "R1 (Serial0/0/0 DCE): 10.10.10.1/30 | R2 (Serial0/0/0 DTE): 10.10.10.2/30",
			Commands: `configure terminal
interface Serial 0/0/0
description LINK_WAN_KE_ROUTER2
ip address 10.10.10.1 255.255.255.252
clock rate 64000
no shutdown
exit`,
			Verification:        "Ketik 'show controllers Serial 0/0/0'. Pastikan terdeteksi 'DCE V.35' dengan clock rate 64000, lalu ping 10.10.10.2 pastikan Success rate 100% (!!!!!).",
			TroubleshootingTips: "Gunakan kabel bertanda jam (Serial DCE) di Packet Tracer. Hanya sisi DCE yang butuh 'clock rate', sisi DTE tidak perlu.",
		},
		{
			ID:          "interface-loopback",
			Title:       "Membuat Interface Virtual (Loopback) untuk Uji Coba / Router-ID",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryInterface,
			Description: "Interface virtual yang selalu UP dan tidak pernah down secara fisik. Sangat berguna untuk ID Router OSPF atau target uji ping.",
			IPExample:   "Loopback 0: 1.1.1.1/32 (Subnet Mask: 255.255.255.255)",
			Commands: `configure terminal
interface Loopback 0
description SIMULASI_SERVER_INTERNET
ip address 1.1.1.1 255.255.255.255
exit`,
			Verification:        "Ketik 'show ip interface brief'. Interface Loopback0 langsung berstatus 'up/up' tanpa perlu perintah 'no shutdown'.",
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
			Description: "Memisahkan broadcast domain pada Switch ke dalam segmen VLAN berbeda (misal: VLAN 10 untuk Guru, VLAN 20 untuk Siswa).",
			IPExample:   "VLAN 10: 192.168.10.0/24 (Port Fa0/1 - Fa0/10) | VLAN 20: 192.168.20.0/24 (Port Fa0/11 - Fa0/20)",
			Commands: `configure terminal
vlan 10
name GURU
vlan 20
name SISWA
exit

interface range FastEthernet 0/1 - 10
switchport mode access
switchport access vlan 10
exit

interface range FastEthernet 0/11 - 20
switchport mode access
switchport access vlan 20
exit`,
			Verification:        "Ketik 'show vlan brief'. Pastikan VLAN 10 (GURU) aktif dan memuat port Fa0/1 s/d Fa0/10, serta VLAN 20 (SISWA) memuat Fa0/11 s/d Fa0/20.",
			TroubleshootingTips: "PC di VLAN 10 TIDAK AKAN BISA saling ping dengan PC di VLAN 20 sebelum dikonfigurasi Router-on-a-Stick (Inter-VLAN Routing).",
		},
		{
			ID:          "vlan-trunk-port",
			Title:       "Mengonfigurasi Port Trunk Antar-Switch atau ke Router",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategoryVLAN,
			Description: "Membuka jalur Trunk pada port penghubung agar dapat melewatkan multi-VLAN sekaligus antar switch atau menuju Router.",
			IPExample:   "Port GigabitEthernet0/1 terhubung ke Switch tetangga atau port GigabitEthernet Router",
			Commands: `configure terminal
interface GigabitEthernet 0/1
description TRUNK_KE_ROUTER_ATAU_SWITCH2
switchport mode trunk
switchport trunk allowed vlan 10,20
exit`,
			Verification:        "Ketik 'show interfaces trunk'. Pastikan port Gi0/1 berstatus 'Mode: on / Status: trunking' dengan encapsulation 802.1q.",
			TroubleshootingTips: "Jika menghubungkan Switch ke Switch lain, pastikan kedua ujung port diatur sebagai trunk mode.",
		},
		{
			ID:          "vlan-svi-management",
			Title:       "IP Manajemen Switch (Switch Virtual Interface - SVI VLAN 1)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategoryVLAN,
			Description: "Memberikan alamat IP pada Switch agar Switch dapat diping dan diremote (SSH/Telnet) dari jaringan komputer.",
			IPExample:   "IP Switch: 192.168.1.2/24 | Default Gateway: 192.168.1.1",
			Commands: `configure terminal
interface vlan 1
ip address 192.168.1.2 255.255.255.0
no shutdown
exit
ip default-gateway 192.168.1.1
exit`,
			Verification:        "Buka PC ber-IP 192.168.1.10, buka Command Prompt lalu ketik 'ping 192.168.1.2'. Pastikan menerima reply.",
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
			IPExample:   "VLAN 10 Gateway: 192.168.10.1/24 | VLAN 20 Gateway: 192.168.20.1/24 (Port Fisik: G0/0)",
			Commands: `configure terminal
interface GigabitEthernet 0/0
no shutdown
exit

interface GigabitEthernet 0/0.10
encapsulation dot1Q 10
ip address 192.168.10.1 255.255.255.0
exit

interface GigabitEthernet 0/0.20
encapsulation dot1Q 20
ip address 192.168.20.1 255.255.255.0
exit`,
			Verification:        "Ketik 'show ip interface brief' di Router. Pastikan G0/0.10 dan G0/0.20 berstatus UP. Lalu uji ping dari PC VLAN 10 ke PC VLAN 20.",
			TroubleshootingTips: "Wajib ketik 'encapsulation dot1Q <nomor_vlan>' TERLEBIH DAHULU sebelum memberi IP address di sub-interface, dan port Switch lawan wajib berstatus TRUNK!",
		},
		{
			ID:          "l3-multilayer-switch-routing",
			Title:       "Inter-VLAN Routing pada Switch Layer 3 (Multilayer Switch 3560)",
			Device:      DeviceSwitchL3,
			Mode:        ModeGlobalConfig,
			Category:    CategoryVLAN,
			Description: "Mengaktifkan fungsi routing pada Switch Layer 3 agar switching dan routing antar VLAN berjalan dengan kecepatan kabel tanpa router eksternal.",
			IPExample:   "VLAN 10: 192.168.10.1/24 | VLAN 20: 192.168.20.1/24",
			Commands: `configure terminal
ip routing

vlan 10
name KANTOR
vlan 20
name LAB

interface vlan 10
ip address 192.168.10.1 255.255.255.0
no shutdown
exit

interface vlan 20
ip address 192.168.20.1 255.255.255.0
no shutdown
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan muncul kode 'C 192.168.10.0/24 is directly connected, Vlan10' dan lakukan ping antar VLAN.",
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
			IPExample:   "IP Next-Hop ISP: 200.100.50.2 (atau keluar lewat Serial0/0/0)",
			Commands: `configure terminal
ip route 0.0.0.0 0.0.0.0 200.100.50.2
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan ada baris 'S* 0.0.0.0/0 [1/0] via 200.100.50.2' dan 'Gateway of last resort is 200.100.50.2 to network 0.0.0.0'.",
			TroubleshootingTips: "0.0.0.0 0.0.0.0 berarti 'segala alamat IP dengan segala subnet mask'. Sangat efisien untuk router cabang yang terhubung ke ISP.",
		},
		{
			ID:          "routing-static-route",
			Title:       "Routing Statis Spesifik (Static Route Menuju Jaringan Tertentu)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Menentukan jalur manual menuju subnet jaringan tertentu di belakang router tetangga.",
			IPExample:   "Tujuan LAN Cabang: 192.168.20.0/24 | Lewat Next-Hop IP: 10.10.10.2",
			Commands: `configure terminal
ip route 192.168.20.0 255.255.255.0 10.10.10.2
exit`,
			Verification:        "Ketik 'show ip route'. Pastikan ada kode 'S 192.168.20.0/24 [1/0] via 10.10.10.2', lalu ping dari Router 1 ke 192.168.20.1.",
			TroubleshootingTips: "Ingat prinsip routing: 'Koneksi dua arah'. Pastikan router lawan juga memiliki static route balik (return route) ke jaringan asal!",
		},
		{
			ID:          "routing-ospf-single-area",
			Title:       "Routing Dinamis OSPF Single Area (OSPF Area 0)",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Mengaktifkan protokol routing dinamis OSPF open standard yang cepat konvergen menggunakan Wildcard Mask.",
			IPExample:   "Router-ID: 1.1.1.1 | Net1: 192.168.1.0/24 | Net2: 10.10.10.0/30 (Wildcard: 0.0.0.3)",
			Commands: `configure terminal
router ospf 1
router-id 1.1.1.1
network 192.168.1.0 0.0.0.255 area 0
network 10.10.10.0 0.0.0.3 area 0
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
			IPExample:   "Jaringan lokal: 192.168.1.0/24 dan link antar router: 10.0.0.0/8",
			Commands: `configure terminal
router rip
version 2
no auto-summary
network 192.168.1.0
network 10.0.0.0
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
			IPExample:   "Pool: POOL_LAN | Network: 192.168.1.0/24 | IP Dikecualikan: 192.168.1.1 s/d 192.168.1.10",
			Commands: `configure terminal
ip dhcp excluded-address 192.168.1.1 192.168.1.10
ip dhcp pool POOL_LAN
network 192.168.1.0 255.255.255.0
default-router 192.168.1.1
dns-server 8.8.8.8
exit`,
			Verification:        "Buka PC di Packet Tracer > Desktop > IP Configuration > Pilih tombol 'DHCP'. Pastikan IP terisi otomatis (cth: 192.168.1.11) dan muncul pesan 'DHCP request successful'.",
			TroubleshootingTips: "'excluded-address' wajib diset agar IP Gateway router (192.168.1.1) tidak bentrok diberikan ke PC klien!",
		},
		{
			ID:          "dhcp-relay-agent",
			Title:       "DHCP Relay Agent (ip helper-address) Lintas Router",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryServices,
			Description: "Meneruskan pesan broadcast DHCP Discover dari PC di satu LAN menuju Server DHCP terpusat di jaringan lain.",
			IPExample:   "Interface Klien: G0/0 | Alamat IP Server DHCP Terpusat: 10.10.10.100",
			Commands: `configure terminal
interface GigabitEthernet 0/0
ip helper-address 10.10.10.100
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
			IPExample:   "Port Fa0/1 hanya boleh dipakai oleh 1 PC yang pertama kali terhubung",
			Commands: `configure terminal
interface FastEthernet 0/1
switchport mode access
switchport port-security
switchport port-security maximum 1
switchport port-security mac-address sticky
switchport port-security violation shutdown
exit`,
			Verification:        "Ketik 'show port-security interface fa0/1'. Pastikan 'Port Security: Enabled' dan status 'Secure-up'. Coba colokkan PC lain, lampu port akan langsung merah (err-disabled).",
			TroubleshootingTips: "Jika port terkunci merah karena violation, buka port kembali dengan masuk ke interface lalu ketik: 'shutdown' kemudian 'no shutdown'.",
		},
		{
			ID:          "security-standard-acl",
			Title:       "Standard Access Control List (Blokir Berdasarkan IP Sumber)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategorySecurity,
			Description: "Memfilter paket data hanya berdasarkan IP Address sumber (source IP). Nomor rentang Standard ACL: 1 - 99.",
			IPExample:   "Blokir PC 192.168.1.50 dari mengakses jaringan lain, izinkan PC lainnya",
			Commands: `configure terminal
access-list 10 deny host 192.168.1.50
access-list 10 permit any

interface GigabitEthernet 0/0
ip access-group 10 in
exit`,
			Verification:        "Uji ping dari PC 192.168.1.50, pastikan hasilnya 'Destination Host Unreachable'. Uji ping dari PC lain (cth: 192.168.1.51), pastikan berhasil reply.",
			TroubleshootingTips: "Ingat kaidah Cisco: 'Implicit Deny Any' ada di akhir semua ACL. Selalu tambahkan 'permit any' di akhir jika tidak ingin semua traffic terblokir!",
		},
		{
			ID:          "security-extended-acl",
			Title:       "Extended Access Control List (Blokir Port/Protokol: Web/Ping)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategorySecurity,
			Description: "Memfilter paket data secara presisi berdasarkan IP asal, IP tujuan, protokol (TCP/UDP/ICMP), dan nomor port (cth: HTTP port 80).",
			IPExample:   "Blokir seluruh LAN 192.168.1.0/24 dari ping (ICMP) ke Server 10.10.10.10, tapi izinkan akses Web HTTP (port 80)",
			Commands: `configure terminal
access-list 100 permit tcp 192.168.1.0 0.0.0.255 host 10.10.10.10 eq 80
access-list 100 deny icmp 192.168.1.0 0.0.0.255 host 10.10.10.10
access-list 100 permit ip any any

interface GigabitEthernet 0/0
ip access-group 100 in
exit`,
			Verification:        "Ketik 'show access-lists'. Coba buka Web Browser di PC ke 10.10.10.10 (sukses) lalu coba ping 10.10.10.10 di CMD (Request timed out). Hit counter ACL akan bertambah.",
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
			IPExample:   "LAN Inside: 192.168.1.0/24 (G0/0) | WAN Outside: Serial0/0/0 (IP Public dari ISP)",
			Commands: `configure terminal
interface GigabitEthernet 0/0
ip nat inside
exit

interface Serial 0/0/0
ip nat outside
exit

access-list 1 permit 192.168.1.0 0.0.0.255

ip nat inside source list 1 interface Serial 0/0/0 overload
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
			IPExample:   "IP Server Private: 192.168.1.100 <==> IP Public Server: 200.100.10.5",
			Commands: `configure terminal
interface GigabitEthernet 0/0
ip nat inside
exit

interface Serial 0/0/0
ip nat outside
exit

ip nat inside source static 192.168.1.100 200.100.10.5
exit`,
			Verification:        "Dari PC di sisi Internet/ISP, lakukan ping atau browsing ke 200.100.10.5. Di Router ketik 'show ip nat translations', akan terlihat entri static translation.",
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
			IPExample:   "Uji Koneksi PC Klien (192.168.1.10) menuju Server Tujuan (10.10.10.10)",
			Commands: `# 1. CEK FISIK (Packet Tracer Link Lights):
#    - Lampu Hijau = Link Up (Normal)
#    - Lampu Oranye = Spanning Tree Protocol (STP) sedang listening (tunggu 30 detik atau klik Fast Forward Time)
#    - Lampu Merah = Port mati (Ketik 'no shutdown') atau tipe kabel salah (Straight vs Cross)

# 2. CEK IP PC DI COMMAND PROMPT:
ipconfig /all

# 3. PING DIRI SENDIRI (Cek TCP/IP Stack PC):
ping 127.0.0.1

# 4. PING DEFAULT GATEWAY (Cek koneksi PC ke Port Router Lokal):
ping 192.168.1.1

# 5. PING ROUTER NEXT-HOP (Cek link antar router WAN):
ping 10.10.10.2

# 6. PING END-TO-END KE SERVER / PC TUJUAN:
ping 10.10.10.10
tracert 10.10.10.10`,
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
