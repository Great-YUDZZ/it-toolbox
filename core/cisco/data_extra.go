package cisco

// getAdvancedCommands returns advanced Cisco configurations for the library
func getAdvancedCommands() []CiscoCommand {
	return []CiscoCommand{
		// ====================================================================
		// SPANNING TREE PROTOCOL & ETHERCHANNEL
		// ====================================================================
		{
			ID:          "stp-pvst-root-bridge",
			Title:       "Spanning Tree Protocol (STP): Rapid-PVST+ & Menentukan Root Bridge",
			Device:      DeviceSwitchL2,
			Mode:        ModeGlobalConfig,
			Category:    CategorySpanningTree,
			Description: "Mengaktifkan protokol Rapid Spanning Tree (802.1w) untuk konvergensi link instan dan menentukan Switch mana yang menjadi Root Bridge utama dan cadangan.",
			Parameters: []Parameter{
				{Key: "VLAN_ID", Label: "VLAN ID", DefaultValue: "10", Placeholder: "cth: 1, 10, 20"},
				{Key: "PRIORITY_VAL", Label: "Prioritas STP (Kelipatan 4096)", DefaultValue: "4096", Placeholder: "cth: 0, 4096, 8192"},
			},
			IPExample: "VLAN: {{VLAN_ID}} | Priority: {{PRIORITY_VAL}} (Default: 32768)",
			Commands: `configure terminal
spanning-tree mode rapid-pvst
# Cara 1: Menetapkan Switch sebagai Root Bridge Utama otomatis
spanning-tree vlan {{VLAN_ID}} root primary
# Atau Cara 2: Menetapkan Root Bridge Cadangan (Secondary)
spanning-tree vlan {{VLAN_ID}} root secondary
# Atau Cara 3: Mengatur nilai Priority secara manual (harus kelipatan 4096)
spanning-tree vlan {{VLAN_ID}} priority {{PRIORITY_VAL}}
exit`,
			Explanation: []LineExplanation{
				{Command: "spanning-tree mode rapid-pvst", Explanation: "Mengaktifkan Rapid Per-VLAN Spanning Tree (konvergensi jauh lebih cepat, hanya hitungan 1-2 detik dibanding standard PVST 30-50 detik)."},
				{Command: "spanning-tree vlan X root primary", Explanation: "Secara otomatis menurunkan priority switch menjadi 24576 (atau 4096 di bawah root saat ini) agar terpilih sebagai Root Bridge."},
				{Command: "spanning-tree vlan X root secondary", Explanation: "Mengatur priority switch menjadi 28672 agar otomatis menjadi Root Bridge cadangan jika Root Bridge utama down."},
				{Command: "spanning-tree vlan X priority Y", Explanation: "Menentukan nilai priority secara presisi. Nilai terkecil akan selalu memenangkan pemilihan Root Bridge."},
			},
			Verification:        "Ketik 'show spanning-tree' atau 'show spanning-tree vlan {{VLAN_ID}}'. Pastikan muncul keterangan 'This bridge is the root' jika switch ini terpilih sebagai Root Bridge.",
			TroubleshootingTips: "Nilai STP priority hanya bisa diisi dengan angka kelipatan 4096 (0, 4096, 8192, 12288, 16384, 20480, 24576, 28672, 32768, dst).",
			Tags:                []string{"stp", "rapid-pvst", "root bridge", "spanning tree", "redundansi l2"},
		},
		{
			ID:          "stp-portfast-bpduguard",
			Title:       "Optimasi Port Akses: Spanning Tree PortFast & BPDU Guard",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategorySpanningTree,
			Description: "Mengabaikan penundaan 30 detik STP pada port yang tersambung langsung ke PC (PortFast) dan mematikan port jika ada yang menyolokkan switch ilegal (BPDU Guard).",
			Parameters: []Parameter{
				{Key: "INT_RANGE", Label: "Range Interface Akses", DefaultValue: "fa0/1 - 10", Placeholder: "cth: fa0/1 - 24"},
			},
			IPExample: "Interface: {{INT_RANGE}} (Port langsung ke PC/Laptop/Server)",
			Commands: `configure terminal
interface range {{INT_RANGE}}
switchport mode access
spanning-tree portfast
spanning-tree bpduguard enable
exit`,
			Explanation: []LineExplanation{
				{Command: "switchport mode access", Explanation: "Memastikan interface beroperasi sebagai port akses ke perangkat akhir."},
				{Command: "spanning-tree portfast", Explanation: "Membuat port langsung masuk ke state FORWARDING tanpa melewati tahap Listening & Learning (lampu indikator langsung hijau seketika)."},
				{Command: "spanning-tree bpduguard enable", Explanation: "Fitur keamanan: jika port menerima paket BPDU (tanda ada switch lain terhubung), port akan langsung dimatikan (err-disabled) untuk mencegah loop."},
			},
			Verification:        "Ketik 'show spanning-tree summary'. Pastikan PortFast default dan BPDU Guard tertera aktif. Hubungkan PC baru, lampu link akan langsung hijau seketika tanpa menunggu 30 detik.",
			TroubleshootingTips: "JANGAN PERNAH mengaktifkan 'spanning-tree portfast' pada port trunk atau port yang terhubung ke switch lain karena bisa memicu switching loop fatal (broadcast storm).",
			Tags:                []string{"portfast", "bpdu guard", "stp security", "switch access"},
		},
		{
			ID:          "etherchannel-lacp-l2",
			Title:       "EtherChannel L2 Menggunakan LACP (IEEE 802.3ad Link Aggregation)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategorySpanningTree,
			Description: "Menggabungkan 2 atau lebih kabel fisik menjadi 1 link logikal (Port-Channel) untuk melipatgandakan kapasitas bandwidth dan menyediakan redundansi otomatis menggunakan standar terbuka LACP.",
			Parameters: []Parameter{
				{Key: "INT_RANGE", Label: "Range Interface Fisik", DefaultValue: "fa0/1 - 2", Placeholder: "cth: fa0/1 - 2 / gig0/1 - 2"},
				{Key: "GROUP_ID", Label: "Nomor Port-Channel Group", DefaultValue: "1", Placeholder: "cth: 1, 2, 3"},
			},
			IPExample: "Switch 1 & Switch 2 terhubung kabel Fa0/1 dan Fa0/2 (Bundled ke Port-channel {{GROUP_ID}})",
			Commands: `configure terminal
interface range {{INT_RANGE}}
shutdown
channel-group {{GROUP_ID}} mode active
no shutdown
exit

# Konfigurasi mode Trunk langsung di Port-Channel virtual
interface port-channel {{GROUP_ID}}
switchport mode trunk
exit`,
			Explanation: []LineExplanation{
				{Command: "interface range fa0/1 - 2", Explanation: "Memilih interface fisik yang akan digabungkan ke dalam bundle EtherChannel."},
				{Command: "shutdown", Explanation: "Mematikan sementara interface fisik untuk menghindari desinkronisasi parameter saat pembentukan channel."},
				{Command: "channel-group 1 mode active", Explanation: "Mengaktifkan LACP mode Active (Switch secara proaktif mengirim paket LACP untuk menegosiasikan EtherChannel)."},
				{Command: "interface port-channel 1", Explanation: "Masuk ke interface virtual hasil penggabungan link fisik."},
				{Command: "switchport mode trunk", Explanation: "Menerapkan konfigurasi trunk langsung pada link agregasi port-channel."},
			},
			Verification:        "Ketik 'show etherchannel summary'. Pastikan kode status port-channel adalah 'SU' (Layer 2, In-Use) dan port fisik bertanda 'P' (Bundled in Port-channel).",
			TroubleshootingTips: "Kedua switch harus memiliki mode yang kompatibel: Active-Active atau Active-Passive. Kecepatan (Speed), Duplex, dan status Access/Trunk pada port anggota harus 100% identik.",
			Tags:                []string{"etherchannel", "lacp", "link aggregation", "port-channel", "trunk"},
		},
		{
			ID:          "etherchannel-pagp-l2",
			Title:       "EtherChannel L2 Menggunakan PAgP (Port Aggregation Protocol Cisco)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategorySpanningTree,
			Description: "Menggabungkan beberapa kabel fisik menjadi link agregasi menggunakan protokol kepemilikan Cisco (PAgP) dengan mode Desirable atau Auto.",
			Parameters: []Parameter{
				{Key: "INT_RANGE", Label: "Range Interface Fisik", DefaultValue: "fa0/3 - 4", Placeholder: "cth: fa0/3 - 4"},
				{Key: "GROUP_ID", Label: "Nomor Channel Group", DefaultValue: "2", Placeholder: "cth: 2"},
			},
			IPExample: "PAgP Group: {{GROUP_ID}} | Switch A: desirable | Switch B: desirable atau auto",
			Commands: `configure terminal
interface range {{INT_RANGE}}
channel-group {{GROUP_ID}} mode desirable
exit

interface port-channel {{GROUP_ID}}
switchport mode trunk
exit`,
			Explanation: []LineExplanation{
				{Command: "channel-group X mode desirable", Explanation: "Mengaktifkan PAgP mode Desirable (secara aktif meminta switch lawan untuk membentuk link agregasi EtherChannel)."},
				{Command: "interface port-channel X", Explanation: "Mengonfigurasi interface bundel agregasi secara terpadu."},
			},
			Verification:        "Ketik 'show etherchannel summary'. Protokol akan tertulis 'PAgP' dan status port-channel 'SU'.",
			TroubleshootingTips: "Kombinasi PAgP yang berhasil adalah Desirable-Desirable atau Desirable-Auto. Kombinasi Auto-Auto TIDAK AKAN PERNAH membentuk EtherChannel karena kedua pihak hanya menunggu.",
			Tags:                []string{"pagp", "etherchannel", "cisco proprietary", "trunk"},
		},
		{
			ID:          "etherchannel-l3-routed",
			Title:       "EtherChannel Layer 3 (Routed Port-Channel Antar Switch Multilayer/Router)",
			Device:      DeviceSwitchL3,
			Mode:        ModeInterface,
			Category:    CategorySpanningTree,
			Description: "Menggabungkan port Layer 3 (routed port ber-IP) menjadi satu sambungan logikal berkecepatan tinggi tanpa fungsi switching/STP.",
			Parameters: []Parameter{
				{Key: "INT_RANGE", Label: "Range Interface Fisik", DefaultValue: "gig0/1 - 2", Placeholder: "cth: gig0/1 - 2"},
				{Key: "GROUP_ID", Label: "Nomor Group", DefaultValue: "5", Placeholder: "cth: 5"},
				{Key: "IP_ADDR", Label: "Alamat IP Port-Channel", DefaultValue: "10.0.0.1", Placeholder: "cth: 10.0.0.1"},
				{Key: "NETMASK", Label: "Subnet Mask", DefaultValue: "255.255.255.252", Placeholder: "cth: 255.255.255.252"},
			},
			IPExample: "Port-Channel {{GROUP_ID}} IP: {{IP_ADDR}} /30 (Link antar Switch L3 Utama)",
			Commands: `configure terminal
ip routing
interface range {{INT_RANGE}}
no switchport
channel-group {{GROUP_ID}} mode active
exit

interface port-channel {{GROUP_ID}}
no switchport
ip address {{IP_ADDR}} {{NETMASK}}
no shutdown
exit`,
			Explanation: []LineExplanation{
				{Command: "no switchport", Explanation: "Mengubah port fisik dari switchport Layer 2 menjadi routed port Layer 3 murni."},
				{Command: "channel-group 5 mode active", Explanation: "Menggabungkan interface Layer 3 menggunakan protokol LACP."},
				{Command: "interface port-channel 5", Explanation: "Masuk ke antarmuka port-channel logikal."},
				{Command: "ip address 10.0.0.1 255.255.255.252", Explanation: "Menetapkan alamat IP langsung pada interface port-channel L3."},
			},
			Verification:        "Ketik 'show etherchannel summary'. Status port-channel harus berhuruf 'RU' (Layer 3 Routed, In-Use). Lakukan ping ke IP port-channel switch lawan.",
			TroubleshootingTips: "Perintah 'no switchport' WAJIB diketikkan pada interface fisik TERLEBIH DAHULU sebelum membentuk channel-group, dan diketikkan lagi di interface port-channel.",
			Tags:                []string{"etherchannel l3", "routed port-channel", "switch l3", "ip routing"},
		},

		// ====================================================================
		// VLAN & TRUNKING MANAGEMENT
		// ====================================================================
		{
			ID:          "vlan-vtp-management",
			Title:       "VTP (VLAN Trunking Protocol): Server, Client & Transparent Mode",
			Device:      DeviceSwitchL2,
			Mode:        ModeGlobalConfig,
			Category:    CategoryVLAN,
			Description: "Menyinkronkan pembuatan, pengubahan, dan penghapusan database VLAN secara otomatis dari satu switch Server ke seluruh switch Client di jaringan.",
			Parameters: []Parameter{
				{Key: "DOMAIN_NAME", Label: "Nama Domain VTP", DefaultValue: "LAB-KAMPUS", Placeholder: "Nama domain bersama"},
				{Key: "PASSWORD", Label: "Password VTP", DefaultValue: "cisco123", Placeholder: "Password rahasia VTP"},
			},
			IPExample: "VTP Domain: {{DOMAIN_NAME}} | Pass: {{PASSWORD}} | Versi: 2",
			Commands: `# DI SWITCH UTAMA (VTP SERVER):
configure terminal
vtp mode server
vtp domain {{DOMAIN_NAME}}
vtp password {{PASSWORD}}
vtp version 2
exit

# DI SWITCH AKSES / CABANG (VTP CLIENT):
configure terminal
vtp mode client
vtp domain {{DOMAIN_NAME}}
vtp password {{PASSWORD}}
vtp version 2
exit`,
			Explanation: []LineExplanation{
				{Command: "vtp mode server", Explanation: "Mengatur switch sebagai Server VTP (bisa membuat, mengedit, dan menghapus VLAN serta menyebarkannya ke Client)."},
				{Command: "vtp domain LAB-KAMPUS", Explanation: "Menetapkan nama domain manajemen VTP. Semua switch harus memiliki nama domain yang sama persis."},
				{Command: "vtp password cisco123", Explanation: "Mengamankan pertukaran pesan VTP dengan password terenkripsi MD5."},
				{Command: "vtp mode client", Explanation: "Mengatur switch sebagai Client (tidak bisa membuat VLAN manual, hanya menerima sinkronisasi dari Server)."},
			},
			Verification:        "Di switch Client, ketik 'show vtp status'. Pastikan 'Configuration Revision' bertambah saat Server membuat VLAN baru, dan ketik 'show vlan brief' untuk memastikan VLAN telah terduplikasi.",
			TroubleshootingTips: "VTP HANYA bisa menyebarkan informasi database VLAN melalui kabel yang telah dikonfigurasi sebagai TRUNK. Pastikan sambungan antar-switch sudah bertipe Trunk!",
			Tags:                []string{"vtp", "vlan trunking protocol", "vtp server", "vtp client", "vlan database"},
		},
		{
			ID:          "vlan-native-and-security",
			Title:       "Keamanan Trunking: Mengubah Native VLAN & Mematikan DTP (Dynamic Trunking)",
			Device:      DeviceSwitchL2,
			Mode:        ModeInterface,
			Category:    CategoryVLAN,
			Description: "Praktik terbaik pengamanan switch: memindahkan Native VLAN dari VLAN 1 bawaan ke VLAN khusus yang tidak terpakai dan mematikan negosiasi otomatis DTP untuk mencegah VLAN Hopping attack.",
			Parameters: []Parameter{
				{Key: "INT_TRUNK", Label: "Interface Trunk", DefaultValue: "gig0/1", Placeholder: "cth: gig0/1 / fa0/24"},
				{Key: "NATIVE_VLAN", Label: "ID Native VLAN Baru", DefaultValue: "999", Placeholder: "cth: 99, 999"},
				{Key: "ALLOWED_LIST", Label: "Daftar VLAN yang Diizinkan", DefaultValue: "10,20,30,999", Placeholder: "cth: 10,20,999"},
			},
			IPExample: "Port Trunk: {{INT_TRUNK}} | Native VLAN: {{NATIVE_VLAN}} | Allowed: {{ALLOWED_LIST}}",
			Commands: `configure terminal
vlan {{NATIVE_VLAN}}
name NATIVE_UNUSED
exit

interface {{INT_TRUNK}}
switchport mode trunk
switchport trunk native vlan {{NATIVE_VLAN}}
switchport trunk allowed vlan {{ALLOWED_LIST}}
switchport nonegotiate
exit`,
			Explanation: []LineExplanation{
				{Command: "vlan 999; name NATIVE_UNUSED", Explanation: "Membuat VLAN khusus kosong untuk menampung lalu lintas untagged (Native) agar terpisah dari VLAN data."},
				{Command: "switchport trunk native vlan 999", Explanation: "Memindahkan Native VLAN dari default VLAN 1 ke VLAN 999 pada interface trunk ini."},
				{Command: "switchport trunk allowed vlan 10,20,30,999", Explanation: "Membatasi hanya VLAN yang terdaftar saja yang boleh melewati kabel trunk (pruning manual)."},
				{Command: "switchport nonegotiate", Explanation: "Mematikan protokol DTP (Dynamic Trunking Protocol) agar port tidak bisa dimanipulasi oleh attacker menjadi trunk ilegal."},
			},
			Verification:        "Ketik 'show interfaces {{INT_TRUNK}} trunk'. Periksa baris 'Native VLAN' (harus {{NATIVE_VLAN}}) dan pastikan tidak ada peringatan 'Native VLAN mismatch' di log terminal.",
			TroubleshootingTips: "Native VLAN HARUS bernilai SAMA pada kedua ujung kabel trunk! Jika switch sebelah menggunakan Native VLAN 1 dan switch ini 999, akan terjadi 'Native VLAN Mismatch' dan paket data bocor.",
			Tags:                []string{"native vlan", "dtp", "switchport nonegotiate", "vlan security", "vlan hopping"},
		},

		// ====================================================================
		// ROUTING ADVANCED: OSPF MULTI-AREA, EIGRP, BGP, FLOATING STATIC
		// ====================================================================
		{
			ID:          "routing-floating-static-route",
			Title:       "Floating Static Route: Jalur Cadangan Otomatis (Failover Backup Route)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Membuat rute statis cadangan dengan Administrative Distance (AD) lebih tinggi yang otomatis aktif menggantikan rute utama hanya saat jalur primer mengalami gangguan/down.",
			Parameters: []Parameter{
				{Key: "DEST_NET", Label: "Network Tujuan", DefaultValue: "192.168.20.0", Placeholder: "cth: 192.168.20.0"},
				{Key: "NETMASK", Label: "Subnet Mask", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
				{Key: "PRIMARY_HOP", Label: "IP Next-Hop Jalur Utama (AD 1)", DefaultValue: "10.10.10.2", Placeholder: "cth: 10.10.10.2"},
				{Key: "BACKUP_HOP", Label: "IP Next-Hop Jalur Cadangan", DefaultValue: "172.16.1.2", Placeholder: "cth: 172.16.1.2"},
				{Key: "BACKUP_AD", Label: "Nilai AD Cadangan (cth: 10/50)", DefaultValue: "10", Placeholder: "cth: 10, 50, 200"},
			},
			IPExample: "Tujuan: {{DEST_NET}} | Jalur Primer: {{PRIMARY_HOP}} (AD: 1) | Jalur Backup: {{BACKUP_HOP}} (AD: {{BACKUP_AD}})",
			Commands: `configure terminal
# 1. Rute Statis Utama (Default AD = 1)
ip route {{DEST_NET}} {{NETMASK}} {{PRIMARY_HOP}}

# 2. Rute Statis Floating / Cadangan (Diberi AD {{BACKUP_AD}})
ip route {{DEST_NET}} {{NETMASK}} {{BACKUP_HOP}} {{BACKUP_AD}}
exit`,
			Explanation: []LineExplanation{
				{Command: "ip route 192.168.20.0 255.255.255.0 10.10.10.2", Explanation: "Rute utama beroperasi normal dengan AD default 1. Rute inilah yang masuk ke Routing Table."},
				{Command: "ip route 192.168.20.0 255.255.255.0 172.16.1.2 10", Explanation: "Rute cadangan disimpan di konfigurasi tetapi 'mengapung' (tidak aktif di routing table) karena AD 10 kalah bersaing dengan AD 1."},
			},
			Verification:        "Ketik 'show ip route'. Saat jalur utama normal, hanya {{PRIMARY_HOP}} yang tertera. Matikan interface jalur utama ('shutdown'), lalu ketik 'show ip route' lagi: rute cadangan {{BACKUP_HOP}} akan langsung muncul menggantikan rute utama.",
			TroubleshootingTips: "Pastikan nilai AD jalur cadangan selalu lebih besar dari rute utama. Jika rute utama adalah OSPF (AD 110), maka floating static route harus diberi AD minimal 115 atau 200.",
			Tags:                []string{"floating static", "backup route", "failover", "administrative distance", "routing"},
		},
		{
			ID:          "routing-ipv6-static-and-default",
			Title:       "IPv6 Routing: Default Route (::/0) & Rute Statis IPv6 Global",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Mengaktifkan kemampuan forwarding paket IPv6 pada router dan membuat rute statis serta default gateway IPv6.",
			Parameters: []Parameter{
				{Key: "NEXT_HOP_IPV6", Label: "IP Next-Hop Router Lawan", DefaultValue: "2001:db8:acad:1::2", Placeholder: "Alamat IPv6 tetangga"},
				{Key: "REMOTE_NET_IPV6", Label: "Subnet IPv6 Remote", DefaultValue: "2001:db8:acad:2::/64", Placeholder: "Subnet tujuan"},
			},
			IPExample: "Default Route: ::/0 | Remote Subnet: {{REMOTE_NET_IPV6}} | Next-Hop: {{NEXT_HOP_IPV6}}",
			Commands: `configure terminal
# WAJIB: Aktifkan proses routing paket IPv6 di router
ipv6 unicast-routing

# 1. Rute Statis Tertuju ke Subnet Spesifik
ipv6 route {{REMOTE_NET_IPV6}} {{NEXT_HOP_IPV6}}

# 2. Default Route IPv6 (Setara 0.0.0.0 0.0.0.0 di IPv4)
ipv6 route ::/0 {{NEXT_HOP_IPV6}}
exit`,
			Explanation: []LineExplanation{
				{Command: "ipv6 unicast-routing", Explanation: "Menghidupkan mesin pengalihan paket IPv6 pada router. Tanpa ini, router Cisco hanya bertindak sebagai end-device IPv6 biasa."},
				{Command: "ipv6 route 2001:db8:acad:2::/64 <next-hop>", Explanation: "Menambahkan entri rute statis menuju prefiks subnet IPv6 tujuan melalui IP hop berikutnya."},
				{Command: "ipv6 route ::/0 <next-hop>", Explanation: "Default route IPv6 (::/0 melambangkan semua alamat IPv6 di internet)."},
			},
			Verification:        "Ketik 'show ipv6 route'. Pastikan rute statis bertanda 'S' dan rute default tertera 'S ::/0'. Coba lakukan 'ping ipv6 <alamat_tujuan>'.",
			TroubleshootingTips: "Jika menggunakan Link-Local Address (fe80::) sebagai next-hop, Anda WAJIB menyertakan nama interface keluar (contoh: 'ipv6 route ::/0 GigabitEthernet0/0 fe80::2').",
			Tags:                []string{"ipv6", "ipv6 route", "unicast-routing", "default route", "dual stack"},
		},
		{
			ID:          "routing-ospf-multi-area",
			Title:       "Routing Dinamis OSPFv2 Multi-Area (Backbone Area 0 & ABR Router)",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Mengonfigurasi OSPF dengan pembagian banyak area untuk mempercepat konvergensi dan membatasi penyebaran banjir LSA, di mana Area Border Router (ABR) menghubungkan Area 1/2 ke Backbone Area 0.",
			Parameters: []Parameter{
				{Key: "PROCESS_ID", Label: "Process ID OSPF", DefaultValue: "1", Placeholder: "cth: 1, 10"},
				{Key: "ROUTER_ID", Label: "Router ID Unik", DefaultValue: "2.2.2.2", Placeholder: "cth: 2.2.2.2"},
				{Key: "NET_AREA0", Label: "Network Jalur Backbone", DefaultValue: "10.10.10.0 0.0.0.3", Placeholder: "cth: 10.10.10.0 0.0.0.3"},
				{Key: "NET_AREA1", Label: "Network LAN Cabang", DefaultValue: "192.168.10.0 0.0.0.255", Placeholder: "cth: 192.168.10.0 0.0.0.255"},
				{Key: "AREA_CABANG", Label: "Nomor Area Cabang", DefaultValue: "1", Placeholder: "cth: 1, 2, 10"},
			},
			IPExample: "ABR Router: Area 0 (Backbone WAN) <---> Area {{AREA_CABANG}} (LAN Departemen)",
			Commands: `configure terminal
router ospf {{PROCESS_ID}}
router-id {{ROUTER_ID}}
# Hubungkan ke Jalur Backbone (Wajib Area 0)
network {{NET_AREA0}} area 0
# Hubungkan ke Area Regional / Cabang
network {{NET_AREA1}} area {{AREA_CABANG}}
# Optimasi: Matikan broadcast hello OSPF ke arah port PC
passive-interface GigabitEthernet0/0
exit`,
			Explanation: []LineExplanation{
				{Command: "router ospf 1", Explanation: "Mengaktifkan proses perutean OSPF nomor 1."},
				{Command: "router-id 2.2.2.2", Explanation: "Menetapkan identitas pengenal router dalam topologi OSPF."},
				{Command: "network 10.10.10.0 0.0.0.3 area 0", Explanation: "Mendaftarkan link serial/gigabit ke dalam Area 0 (Backbone OSPF yang wajib ada)."},
				{Command: "network 192.168.10.0 0.0.0.255 area 1", Explanation: "Mendaftarkan link LAN cabang ke Area 1. Karena router ini memiliki kaki di Area 0 dan Area 1, router ini resmi menjadi ABR (Area Border Router)."},
				{Command: "passive-interface GigabitEthernet0/0", Explanation: "Mencegah pengiriman paket Hello OSPF ke arah switch/PC LAN agar menghemat bandwidth dan meningkatkan keamanan."},
			},
			Verification:        "Ketik 'show ip ospf neighbor' untuk melihat hubungan ketetanggaan (status FULL). Ketik 'show ip route ospf' di router Area 1: rute dari area lain akan berkode 'O IA' (OSPF Inter-Area).",
			TroubleshootingTips: "Aturan utama OSPF Multi-Area: SEMUA Area non-backbone (Area 1, 2, 3...) HARUS terhubung langsung secara fisik atau logikal ke Area 0 (Backbone).",
			Tags:                []string{"ospf", "multi-area", "abr", "area 0", "inter-area", "o ia"},
		},
		{
			ID:          "routing-ospf-tuning-bandwidth",
			Title:       "Tuning Metric OSPF: Auto-Cost Reference-Bandwidth & Manual Interface Cost",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryRouting,
			Description: "Menyesuaikan rumus kalkulasi metric OSPF agar mampu membedakan link Gigabit (1 Gbps) dan 10 Gigabit dari FastEthernet (100 Mbps), serta cara mengatur cost link secara manual.",
			Parameters: []Parameter{
				{Key: "INT_NAME", Label: "Interface Spesifik", DefaultValue: "GigabitEthernet0/0", Placeholder: "cth: Gig0/0 / Serial0/0/0"},
				{Key: "MANUAL_COST", Label: "Nilai Cost Manual", DefaultValue: "10", Placeholder: "cth: 10, 50, 100"},
			},
			IPExample: "Ref Bandwidth: 1000 Mbps (1 Gbps) | Interface: {{INT_NAME}} Cost: {{MANUAL_COST}}",
			Commands: `configure terminal
# Langkah 1: Ubah Reference Bandwidth agar link Gigabit terhitung akurat
router ospf 1
auto-cost reference-bandwidth 1000
exit

# Langkah 2: Mengatur Cost OSPF secara manual pada interface spesifik
interface {{INT_NAME}}
ip ospf cost {{MANUAL_COST}}
exit`,
			Explanation: []LineExplanation{
				{Command: "auto-cost reference-bandwidth 1000", Explanation: "Mengubah basis pembagi metric default (100 Mbps) menjadi 1000 Mbps (1 Gbps), sehingga GigabitEthernet mendapat cost 1 dan FastEthernet mendapat cost 10."},
				{Command: "ip ospf cost 10", Explanation: "Memaksa interface memiliki cost spesifik secara manual tanpa memedulikan bandwidth bawaan kabel."},
			},
			Verification:        "Ketik 'show ip ospf interface {{INT_NAME}}'. Periksa nilai 'Cost: {{MANUAL_COST}}'.",
			TroubleshootingTips: "Perintah 'auto-cost reference-bandwidth' HARUS diterapkan dengan nilai yang SAMA pada SELURUH router OSPF di jaringan agar kalkulasi metric konsisten.",
			Tags:                []string{"ospf cost", "reference-bandwidth", "ospf metric", "tuning"},
		},
		{
			ID:          "routing-ospf-default-route-seeding",
			Title:       "Menyebarkan Rute Internet ke Seluruh Jaringan: OSPF Default-Information Originate",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Membuat rute default internet di router gateway border dan menyebarkannya secara otomatis ke seluruh router internal tanpa perlu mengetik rute manual di tiap router.",
			Parameters: []Parameter{
				{Key: "ISP_IP", Label: "Alamat IP Router ISP / Internet", DefaultValue: "203.0.113.2", Placeholder: "cth: 203.0.113.2"},
			},
			IPExample: "Border Gateway Router terhubung ke ISP ({{ISP_IP}}) menyebarkan rute 0.0.0.0/0 via OSPF",
			Commands: `configure terminal
# 1. Buat rute default statis mengarah ke ISP luar
ip route 0.0.0.0 0.0.0.0 {{ISP_IP}}

# 2. Sebar rute default ke seluruh tetangga OSPF
router ospf 1
default-information originate
exit`,
			Explanation: []LineExplanation{
				{Command: "ip route 0.0.0.0 0.0.0.0 203.0.113.2", Explanation: "Menentukan arah lalu lintas internet ke alamat IP router penyedia layanan (ISP)."},
				{Command: "default-information originate", Explanation: "Memerintahkan OSPF untuk mengiklankan rute default (0.0.0.0/0) sebagai LSA Tipe 5 (External) ke seluruh area OSPF."},
			},
			Verification:        "Buka router tetangga di dalam jaringan, ketik 'show ip route'. Akan muncul baris rute 'O*E2 0.0.0.0/0' yang mengarah kembali ke router border gateway ini.",
			TroubleshootingTips: "Jika rute default statis (ip route 0.0.0.0 0.0.0.0) belum ada di router, OSPF tidak akan menyebarkannya. Jika ingin OSPF menyebarkannya walau belum ada rute statis, tambahkan kata 'always': 'default-information originate always'.",
			Tags:                []string{"ospf default route", "default-information originate", "isp sharing", "o*e2"},
		},
		{
			ID:          "routing-eigrp-as",
			Title:       "Routing Dinamis EIGRP (Enhanced Interior Gateway Routing Protocol)",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Mengonfigurasi protokol routing canggih EIGRP dengan konvergensi instan DUAL algorithm, metric gabungan bandwidth/delay, dan penonaktifan auto-summary.",
			Parameters: []Parameter{
				{Key: "AS_NUMBER", Label: "Nomor Autonomous System (AS)", DefaultValue: "100", Placeholder: "cth: 10, 100"},
				{Key: "NET1", Label: "Network 1 & Wildcard", DefaultValue: "192.168.1.0 0.0.0.255", Placeholder: "cth: 192.168.1.0 0.0.0.255"},
				{Key: "NET2", Label: "Network 2 & Wildcard", DefaultValue: "10.10.10.0 0.0.0.3", Placeholder: "cth: 10.10.10.0 0.0.0.3"},
			},
			IPExample: "EIGRP AS: {{AS_NUMBER}} | Subnet 1: {{NET1}} | Subnet 2: {{NET2}}",
			Commands: `configure terminal
router eigrp {{AS_NUMBER}}
no auto-summary
network {{NET1}}
network {{NET2}}
passive-interface GigabitEthernet0/0
exit`,
			Explanation: []LineExplanation{
				{Command: "router eigrp 100", Explanation: "Mengaktifkan proses routing EIGRP dengan nomor Autonomous System (AS) 100."},
				{Command: "no auto-summary", Explanation: "Mencegah EIGRP meringkas subnet ke kelas alaminya (classful), sehingga mendukung VLSM dan subnetting modern."},
				{Command: "network 192.168.1.0 0.0.0.255", Explanation: "Mendaftarkan subnet interface LAN beserta wildcard mask-nya ke EIGRP."},
				{Command: "passive-interface GigabitEthernet0/0", Explanation: "Menghentikan pengiriman paket Hello EIGRP ke interface yang mengarah ke PC user."},
			},
			Verification:        "Ketik 'show ip eigrp neighbors' untuk melihat daftar router tetangga yang terbentuk. Ketik 'show ip route eigrp': rute yang dipelajari akan memiliki kode huruf 'D' (DUAL).",
			TroubleshootingTips: "Nomor Autonomous System (AS) HARUS IDENTIK pada semua router yang ingin bertukar informasi EIGRP. Jika router 1 memakai AS 100 dan router 2 memakai AS 200, keduanya tidak akan pernah bertetangga.",
			Tags:                []string{"eigrp", "dual", "autonomous system", "as number", "routing dinamis"},
		},
		{
			ID:          "routing-bgp-basic",
			Title:       "Routing Eksterior BGP (Border Gateway Protocol) Dasar: Peering eBGP Antar AS",
			Device:      DeviceRouter,
			Mode:        ModeRouterConfig,
			Category:    CategoryRouting,
			Description: "Mengonfigurasi peering External BGP (eBGP) antar dua Autonomous System (AS) yang berbeda untuk menyimulasikan interkoneksi ISP atau kantor multi-site berskala besar.",
			Parameters: []Parameter{
				{Key: "MY_AS", Label: "Nomor AS Router Lokal", DefaultValue: "65001", Placeholder: "cth: 65001"},
				{Key: "PEER_IP", Label: "Alamat IP Neighbor Lawan", DefaultValue: "198.51.100.2", Placeholder: "cth: 198.51.100.2"},
				{Key: "PEER_AS", Label: "Nomor AS Neighbor Lawan", DefaultValue: "65002", Placeholder: "cth: 65002"},
				{Key: "LOCAL_NET", Label: "Network Lokal yang Diiklankan", DefaultValue: "192.168.1.0", Placeholder: "cth: 192.168.1.0"},
				{Key: "LOCAL_MASK", Label: "Subnet Mask Subnet Lokal", DefaultValue: "255.255.255.0", Placeholder: "cth: 255.255.255.0"},
			},
			IPExample: "Router Lokal (AS {{MY_AS}}) <--- eBGP Peering ---> Router Lawan (AS {{PEER_AS}} @ {{PEER_IP}})",
			Commands: `configure terminal
router bgp {{MY_AS}}
neighbor {{PEER_IP}} remote-as {{PEER_AS}}
network {{LOCAL_NET}} mask {{LOCAL_MASK}}
exit`,
			Explanation: []LineExplanation{
				{Command: "router bgp 65001", Explanation: "Mengaktifkan proses perutean BGP dengan nomor Autonomous System lokal."},
				{Command: "neighbor 198.51.100.2 remote-as 65002", Explanation: "Mendaftarkan tetangga peering eBGP (karena remote-as 65002 berbeda dengan AS lokal 65001, mode yang aktif otomatis adalah eBGP)."},
				{Command: "network 192.168.1.0 mask 255.255.255.0", Explanation: "Mengiklankan awalan subnet lokal ke dalam tabel BGP dunia."},
			},
			Verification:        "Ketik 'show ip bgp summary'. Kolom 'State/PfxRcd' harus menunjukkan angka (bukan kata 'Active' atau 'Idle'). Ketik 'show ip route bgp' untuk melihat rute berkode huruf 'B'.",
			TroubleshootingTips: "Dalam perintah 'network' di BGP, subnet mask HARUS cocok persis dengan entri yang ada di routing table lokal router. Jika tidak cocok, rute tidak akan pernah diiklankan.",
			Tags:                []string{"bgp", "ebgp", "autonomous system", "isp peering", "wan"},
		},

		// ====================================================================
		// GATEWAY REDUNDANCY (FHRP: HSRP & VRRP)
		// ====================================================================
		{
			ID:          "fhrp-hsrp-gateway-redundancy",
			Title:       "Redundansi Default Gateway: HSRP (Hot Standby Router Protocol Cisco)",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryRedundancy,
			Description: "Menggabungkan dua router menjadi 1 Virtual IP Default Gateway bersama. Jika router Active mati, router Standby otomatis mengambil alih lalu lintas dalam hitungan detik tanpa putus koneksi.",
			Parameters: []Parameter{
				{Key: "INT_LAN", Label: "Interface LAN", DefaultValue: "GigabitEthernet0/0", Placeholder: "cth: Gig0/0 / Fa0/0"},
				{Key: "GROUP_ID", Label: "Nomor Group HSRP", DefaultValue: "1", Placeholder: "cth: 1"},
				{Key: "VIRTUAL_IP", Label: "Alamat Virtual IP Gateway Bersama", DefaultValue: "192.168.1.254", Placeholder: "IP Gateway untuk PC user"},
				{Key: "PRIORITY_VAL", Label: "Nilai Priority (Router Utama)", DefaultValue: "110", Placeholder: "Default: 100, Utama: 110"},
			},
			IPExample: "Router 1 (Active, Priority: {{PRIORITY_VAL}}) | Router 2 (Standby, Priority: 100) | Virtual IP: {{VIRTUAL_IP}}",
			Commands: `# KONFIGURASI DI ROUTER 1 (ACTIVE / UTAMA):
configure terminal
interface {{INT_LAN}}
standby version 2
standby {{GROUP_ID}} ip {{VIRTUAL_IP}}
standby {{GROUP_ID}} priority {{PRIORITY_VAL}}
standby {{GROUP_ID}} preempt
exit

# KONFIGURASI DI ROUTER 2 (STANDBY / CADANGAN):
configure terminal
interface {{INT_LAN}}
standby version 2
standby {{GROUP_ID}} ip {{VIRTUAL_IP}}
# Priority default 100, tidak perlu diubah agar Router 1 menjadi pemenang
exit`,
			Explanation: []LineExplanation{
				{Command: "standby version 2", Explanation: "Mengaktifkan HSRP Versi 2 yang mendukung IPv6 dan nomor group lebih besar (0-4095)."},
				{Command: "standby 1 ip 192.168.1.254", Explanation: "Menentukan Alamat Virtual IP yang akan dipasang sebagai 'Default Gateway' pada semua PC di jaringan LAN."},
				{Command: "standby 1 priority 110", Explanation: "Menaikkan prioritas Router 1 di atas default (100), sehingga Router 1 terpilih menjadi ACTIVE Router."},
				{Command: "standby 1 preempt", Explanation: "Fitur Preempt: Jika Router 1 sempat mati lalu menyala kembali, Router 1 akan langsung merebut kembali status Active-nya secara terhormat."},
			},
			Verification:        "Ketik 'show standby brief'. Pastikan Router 1 berstatus 'Active' dan Router 2 berstatus 'Standby'. Pada PC, pasang gateway {{VIRTUAL_IP}}, jalankan ping berkelanjutan ('ping -t'), lalu matikan port Router 1; ping hanya akan drop 1 paket lalu kembali reply!",
			TroubleshootingTips: "Semua PC/klien di jaringan LAN HARUS memasang IP Virtual Gateway ({{VIRTUAL_IP}}) sebagai Default Gateway mereka, BUKAN IP fisik dari Router 1 maupun Router 2.",
			Tags:                []string{"hsrp", "fhrp", "gateway redundancy", "virtual ip", "standby preempt"},
		},
		{
			ID:          "fhrp-vrrp-gateway",
			Title:       "Redundansi Default Gateway Terbuka: VRRP (Virtual Router Redundancy Protocol)",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryRedundancy,
			Description: "Standar terbuka industri (RFC 3768) untuk redundansi gateway antar-vendor yang mirip dengan HSRP, membagi peran Master dan Backup Router.",
			Parameters: []Parameter{
				{Key: "INT_LAN", Label: "Interface LAN", DefaultValue: "GigabitEthernet0/0", Placeholder: "cth: Gig0/0"},
				{Key: "GROUP_ID", Label: "Nomor Group VRRP", DefaultValue: "10", Placeholder: "cth: 10"},
				{Key: "VIRTUAL_IP", Label: "Alamat Virtual IP Bersama", DefaultValue: "192.168.10.254", Placeholder: "cth: 192.168.10.254"},
			},
			IPExample: "Master Router (Priority 120) <---> Backup Router (Priority 100) | Virtual IP: {{VIRTUAL_IP}}",
			Commands: `# DI ROUTER MASTER:
configure terminal
interface {{INT_LAN}}
vrrp {{GROUP_ID}} ip {{VIRTUAL_IP}}
vrrp {{GROUP_ID}} priority 120
exit

# DI ROUTER BACKUP:
configure terminal
interface {{INT_LAN}}
vrrp {{GROUP_ID}} ip {{VIRTUAL_IP}}
exit`,
			Explanation: []LineExplanation{
				{Command: "vrrp 10 ip 192.168.10.254", Explanation: "Mendaftarkan Virtual IP gateway untuk group VRRP 10."},
				{Command: "vrrp 10 priority 120", Explanation: "Menaikkan prioritas Master Router di atas 100. Pada VRRP, preempt aktif secara default."},
			},
			Verification:        "Ketik 'show vrrp brief'. Router Master akan berstatus 'Master' dan router kedua berstatus 'Backup'.",
			TroubleshootingTips: "Berbeda dengan HSRP, pada VRRP fitur 'preempt' sudah otomatis aktif secara default tanpa perlu diketik.",
			Tags:                []string{"vrrp", "fhrp", "open standard", "gateway backup"},
		},

		// ====================================================================
		// KEAMANAN SISTEM, HARDENING & LOGGING
		// ====================================================================
		{
			ID:          "hardening-aaa-local-security",
			Title:       "Keamanan Akses Perangkat: AAA Model Baru & Akun Administrator Privilege 15",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryHardening,
			Description: "Menerapkan framework AAA (Authentication, Authorization, Accounting) lokal untuk login Console dan SSH menggunakan database username dan password terenkripsi.",
			Parameters: []Parameter{
				{Key: "ADMIN_USER", Label: "Nama User Admin", DefaultValue: "netadmin", Placeholder: "cth: netadmin / sysadmin"},
				{Key: "ADMIN_PASS", Label: "Password Admin (Secret)", DefaultValue: "Admin#StrongPass99!", Placeholder: "Password kuat"},
			},
			IPExample: "Akun: {{ADMIN_USER}} (Hak Akses Penuh Privilege 15)",
			Commands: `configure terminal
# 1. Buat akun administrator dengan hak akses tertinggi (Privilege 15)
username {{ADMIN_USER}} privilege 15 secret {{ADMIN_PASS}}

# 2. Aktifkan framework AAA
aaa new-model

# 3. Tetapkan metode autentikasi login lokal sebagai standar
aaa authentication login default local
aaa authorization exec default local
exit`,
			Explanation: []LineExplanation{
				{Command: "username netadmin privilege 15 secret ...", Explanation: "Membuat akun lokal dengan hak akses tertinggi (Privilege 15 otomatis langsung masuk ke mode Privileged EXEC '#')."},
				{Command: "aaa new-model", Explanation: "Mengaktifkan arsitektur keamanan Authentication, Authorization, dan Accounting Cisco."},
				{Command: "aaa authentication login default local", Explanation: "Memerintahkan semua jalur login (Console, AUX, Telnet, SSH) untuk memverifikasi username & password ke database lokal router."},
			},
			Verification:        "Ketik 'exit' hingga keluar ke prompt login. Saat diminta login, masukkan username '{{ADMIN_USER}}' dan password '{{ADMIN_PASS}}'.",
			TroubleshootingTips: "JANGAN PERNAH mengaktifkan 'aaa new-model' sebelum membuat akun username lokal! Jika lupa membuat user, Anda bisa terkunci dari perangkat sama sekali.",
			Tags:                []string{"aaa", "privilege 15", "local user", "hardening", "keamanan"},
		},
		{
			ID:          "hardening-login-bruteforce",
			Title:       "Mitigasi Serangan Password: Login Block-For (Anti Brute-Force Lockout)",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryHardening,
			Description: "Memblokir secara otomatis upaya login jika terjadi kegagalan memasukkan password berulang kali dalam durasi waktu tertentu.",
			Parameters: []Parameter{
				{Key: "BLOCK_SEC", Label: "Durasi Blokir (Detik)", DefaultValue: "120", Placeholder: "cth: 120 (2 menit)"},
				{Key: "ATTEMPTS", Label: "Maksimal Percobaan Gagal", DefaultValue: "3", Placeholder: "cth: 3 atau 5"},
				{Key: "WINDOW_SEC", Label: "Jendela Waktu Pantau (Detik)", DefaultValue: "60", Placeholder: "cth: 60 detik"},
			},
			IPExample: "Blokir {{BLOCK_SEC}} detik jika salah {{ATTEMPTS}}x dalam kurun {{WINDOW_SEC}} detik",
			Commands: `configure terminal
login block-for {{BLOCK_SEC}} attempts {{ATTEMPTS}} within {{WINDOW_SEC}}
login delay 2
login on-failure log
login on-success log
exit`,
			Explanation: []LineExplanation{
				{Command: "login block-for 120 attempts 3 within 60", Explanation: "Jika ada yang salah memasukkan password 3 kali dalam waktu 60 detik, router otomatis mengunci akses selama 120 detik (2 menit)."},
				{Command: "login delay 2", Explanation: "Menambahkan jeda paksa 2 detik setiap kali percobaan login gagal untuk memperlambat serangan otomatis/script brute force."},
				{Command: "login on-failure log", Explanation: "Mencatat log peringatan setiap kali ada upaya login yang gagal."},
			},
			Verification:        "Ketik 'show login'. Status 'Login quiet mode' akan aktif jika ada yang salah password berturut-turut melebihi batas percobaan.",
			TroubleshootingTips: "Simpan selalu backup konfigurasi di luar perangkat sebelum menguji brute force lockout agar tidak terhambat saat praktikum.",
			Tags:                []string{"login block-for", "brute force", "lockout", "hardening", "keamanan router"},
		},
		{
			ID:          "hardening-ntp-clock",
			Title:       "Sinkronisasi Waktu Akurat: NTP Client & Konfigurasi Timezone WIB",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryHardening,
			Description: "Menyelaraskan jam internal router/switch dengan server waktu NTP (Network Time Protocol) dan mengatur zona waktu Indonesia Barat (WIB +7) agar rekaman log kejadian akurat.",
			Parameters: []Parameter{
				{Key: "NTP_SERVER_IP", Label: "Alamat IP Server NTP", DefaultValue: "192.168.1.250", Placeholder: "cth: 192.168.1.250"},
			},
			IPExample: "Zona: WIB (UTC+7) | NTP Server: {{NTP_SERVER_IP}}",
			Commands: `configure terminal
# 1. Atur Zona Waktu Indonesia Barat (WIB UTC+7)
clock timezone WIB 7 0

# 2. Hubungkan ke Server NTP
ntp server {{NTP_SERVER_IP}}
service timestamps log datetime msec
service timestamps debug datetime msec
exit`,
			Explanation: []LineExplanation{
				{Command: "clock timezone WIB 7 0", Explanation: "Menentukan nama zona waktu lokal 'WIB' dengan offset +7 jam dari UTC."},
				{Command: "ntp server 192.168.1.250", Explanation: "Menunjuk router atau server NTP yang menjadi acuan jam waktu bersama."},
				{Command: "service timestamps log datetime msec", Explanation: "Memastikan pesan syslog terminal mencantumkan tanggal, jam, menit, detik, dan milidetik yang presisi."},
			},
			Verification:        "Ketik 'show clock'. Jam akan tertera sesuai zona WIB. Ketik 'show ntp status' untuk melihat status sinkronisasi ('Clock is synchronized').",
			TroubleshootingTips: "Di Cisco Packet Tracer, sinkronisasi NTP membutuhkan waktu beberapa siklus. Klik tombol 'Fast Forward Time' (Alt + D) beberapa kali untuk mempercepat sinkronisasi.",
			Tags:                []string{"ntp", "clock", "timezone", "wib", "logging timestamps"},
		},
		{
			ID:          "hardening-syslog-server",
			Title:       "Pusat Log Keamanan: Mengirim Pesan Event ke Server Syslog Eksternal",
			Device:      DeviceRouter,
			Mode:        ModeGlobalConfig,
			Category:    CategoryHardening,
			Description: "Mengirimkan seluruh pesan log kejadian (link down, login gagal, perubahan config) secara real-time ke Server Syslog untuk analisis forensik jaringan.",
			Parameters: []Parameter{
				{Key: "SYSLOG_IP", Label: "Alamat IP Server Syslog", DefaultValue: "192.168.1.251", Placeholder: "IP Server Syslog di Packet Tracer"},
			},
			IPExample: "Syslog Server: {{SYSLOG_IP}} | Trap Level: Informational (6)",
			Commands: `configure terminal
logging on
logging host {{SYSLOG_IP}}
logging trap informational
logging source-interface GigabitEthernet0/0
exit`,
			Explanation: []LineExplanation{
				{Command: "logging on", Explanation: "Mengaktifkan fasilitas logging perangkat."},
				{Command: "logging host 192.168.1.251", Explanation: "Menentukan IP server penerima paket Syslog UDP port 514."},
				{Command: "logging trap informational", Explanation: "Mengirimkan pesan dari tingkat keparahan 0 (Emergency) hingga level 6 (Informational)."},
				{Command: "logging source-interface Gig0/0", Explanation: "Memastikan paket syslog dikirim menggunakan alamat IP antarmuka LAN lokal."},
			},
			Verification:        "Di Cisco Packet Tracer, buka perangkat Server ({{SYSLOG_IP}}), buka tab 'Services' -> 'SYSLOG'. Lakukan 'shutdown' pada salah satu port router, periksa tabel log di server bertambah secara otomatis.",
			TroubleshootingTips: "Pastikan koneksi routing antara router dengan Server Syslog berjalan lancar (uji ping terlebih dahulu).",
			Tags:                []string{"syslog", "logging", "monitoring", "server log"},
		},

		// ====================================================================
		// WAN & TUNNELING
		// ====================================================================
		{
			ID:          "wan-hdlc-ppp-encapsulation",
			Title:       "Koneksi WAN Serial: Enkapsulasi PPP dengan Autentikasi CHAP / PAP",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryWAN,
			Description: "Mengonfigurasi link WAN serial point-to-point menggunakan protokol terbuka Point-to-Point Protocol (PPP) yang diamankan oleh sistem tantangan sandi CHAP (Challenge Handshake Authentication Protocol).",
			Parameters: []Parameter{
				{Key: "INT_SERIAL", Label: "Interface Serial", DefaultValue: "Serial0/0/0", Placeholder: "cth: Serial0/0/0"},
				{Key: "REMOTE_ROUTER_NAME", Label: "Hostname Router Lawan", DefaultValue: "R2-Cabang", Placeholder: "cth: R2-Cabang"},
				{Key: "SHARED_SECRET", Label: "Password Bersama (Shared Secret)", DefaultValue: "ciscoWanPass", Placeholder: "Password yang sama persis"},
			},
			IPExample: "R1-Pusat <=== PPP CHAP ===> {{REMOTE_ROUTER_NAME}} (Password: {{SHARED_SECRET}})",
			Commands: `# DI ROUTER 1 (R1-Pusat):
configure terminal
# 1. Daftarkan username router lawan beserta password bersama
username {{REMOTE_ROUTER_NAME}} password {{SHARED_SECRET}}

# 2. Pasang enkapsulasi PPP dan autentikasi CHAP di port serial
interface {{INT_SERIAL}}
encapsulation ppp
ppp authentication chap
no shutdown
exit

# CATATAN: DI ROUTER LAWAN ({{REMOTE_ROUTER_NAME}}):
# Anda wajib mendaftarkan hostname router pertama dan password yang sama:
# username R1-Pusat password {{SHARED_SECRET}}
# interface {{INT_SERIAL}}
# encapsulation ppp
# ppp authentication chap`,
			Explanation: []LineExplanation{
				{Command: "username R2-Cabang password ...", Explanation: "Mendaftarkan identitas router lawan pada database lokal (kunci nama harus sama persis dengan hostname lawan)."},
				{Command: "encapsulation ppp", Explanation: "Mengganti enkapsulasi default HDLC menjadi protokol standar terbuka PPP."},
				{Command: "ppp authentication chap", Explanation: "Mengaktifkan three-way handshake autentikasi CHAP yang aman karena password tidak pernah dikirim dalam teks polos."},
			},
			Verification:        "Ketik 'show interfaces {{INT_SERIAL}}'. Baris protokol harus menyatakan: 'Encapsulation PPP, LCP Open, loopback not set'.",
			TroubleshootingTips: "Pada PPP CHAP, nama username HARUS SAMA PERSIS (case-sensitive) dengan 'hostname' router sebelah, dan password 'secret' di kedua router harus 100% identik.",
			Tags:                []string{"ppp", "chap", "wan", "serial interface", "encapsulation"},
		},
		{
			ID:          "wan-gre-tunnel-site-to-site",
			Title:       "Site-to-Site Virtual Tunnel: GRE (Generic Routing Encapsulation) Tunnel",
			Device:      DeviceRouter,
			Mode:        ModeInterface,
			Category:    CategoryWAN,
			Description: "Membangun terowongan virtual point-to-point melintasi jaringan internet/WAN publik sehingga dua kantor cabang dapat saling berkomunikasi dan menjalankan routing protocol privat.",
			Parameters: []Parameter{
				{Key: "TUNNEL_ID", Label: "Nomor Tunnel", DefaultValue: "0", Placeholder: "cth: 0, 1"},
				{Key: "TUNNEL_IP", Label: "IP Privat Interface Tunnel", DefaultValue: "10.255.255.1", Placeholder: "cth: 10.255.255.1"},
				{Key: "TUNNEL_NETMASK", Label: "Subnet Mask Tunnel", DefaultValue: "255.255.255.252", Placeholder: "cth: 255.255.255.252 (/30)"},
				{Key: "LOCAL_WAN_INT", Label: "Interface Fisik WAN Lokal", DefaultValue: "GigabitEthernet0/1", Placeholder: "cth: Gig0/1 / Serial0/0/0"},
				{Key: "REMOTE_WAN_IP", Label: "IP Publik WAN Router Lawan", DefaultValue: "203.0.113.2", Placeholder: "cth: 203.0.113.2"},
			},
			IPExample: "Tunnel {{TUNNEL_ID}} IP: {{TUNNEL_IP}}/30 | Sumber: {{LOCAL_WAN_INT}} | Tujuan: {{REMOTE_WAN_IP}}",
			Commands: `configure terminal
interface Tunnel{{TUNNEL_ID}}
ip address {{TUNNEL_IP}} {{TUNNEL_NETMASK}}
tunnel source {{LOCAL_WAN_INT}}
tunnel destination {{REMOTE_WAN_IP}}
tunnel mode gre ip
no shutdown
exit`,
			Explanation: []LineExplanation{
				{Command: "interface Tunnel0", Explanation: "Membuat antarmuka logikal terowongan virtual."},
				{Command: "ip address 10.255.255.1 255.255.255.252", Explanation: "Memberikan alamat IP privat point-to-point (/30) pada link virtual tunnel."},
				{Command: "tunnel source Gig0/1", Explanation: "Menentukan antarmuka fisik lokal yang menjadi pintu keluar terowongan paket."},
				{Command: "tunnel destination 203.0.113.2", Explanation: "Menentukan alamat IP publik router tujuan di sisi seberang WAN."},
				{Command: "tunnel mode gre ip", Explanation: "Menggunakan enkapsulasi GRE standar Cisco/RFC."},
			},
			Verification:        "Ketik 'show ip interface brief'. Pastikan 'Tunnel{{TUNNEL_ID}}' berstatus Status UP dan Protocol UP. Coba ping alamat IP tunnel lawan (misal 10.255.255.2).",
			TroubleshootingTips: "Sebelum membuat tunnel, pastikan IP sumber dan IP tujuan ({{REMOTE_WAN_IP}}) sudah bisa saling ping melalui jaringan publik internet.",
			Tags:                []string{"gre", "tunnel", "vpn", "wan site-to-site", "tunneling"},
		},

		// ====================================================================
		// BACKUP, RESET & PASSWORD RECOVERY
		// ====================================================================
		{
			ID:          "recovery-password-rommon",
			Title:       "Prosedur Pemulihan Password Router Cisco (Rommon Mode & Config-Register 0x2142)",
			Device:      DeviceRouter,
			Mode:        ModePrivExec,
			Category:    CategoryRecovery,
			Description: "Panduan resmi mengatasi lupa password enable secret pada Router Cisco dengan mem-bypass pembacaan file startup-config di NVRAM via Rommon mode tanpa menghapus konfigurasi lama.",
			Commands: `# LANGKAH 1: REBOOT ROUTER & MASUK KE ROMMON
# Matikan saklar power router di Packet Tracer lalu nyalakan kembali.
# Segera tekan tombol 'Ctrl + Break' (atau Ctrl + C) di terminal.
# Prompt akan berubah menjadi: 'rommon 1 >'

# LANGKAH 2: UBAH REGISTER KE 0x2142 (BYPASS STARTUP-CONFIG)
rommon 1 > confreg 0x2142
rommon 2 > reset
# Router akan booting ulang tanpa membaca konfigurasi lama (seperti router baru).
# Saat ditanya 'Would you like to enter initial configuration? [yes/no]:' ketik: no

# LANGKAH 3: PULIHKAN CONFIG LAMA KE RAM & GANTI PASSWORD
Router> enable
Router# copy startup-config running-config
# Nama router dan konfigurasi lama Anda kini kembali muncul!
Router# configure terminal
Router(config)# enable secret passwordBaru123

# LANGKAH 4: KEMBALIKAN REGISTER KE NORMAL (0x2102) & SIMPAN
Router(config)# config-register 0x2102
Router(config)# exit
Router# copy running-config startup-config
Router# reload`,
			Explanation: []LineExplanation{
				{Command: "confreg 0x2142", Explanation: "Mengubah nilai register prosesor agar mengabaikan file startup-config di NVRAM saat booting."},
				{Command: "reset", Explanation: "Me-reboot router langsung dari lingkungan darurat Rommon."},
				{Command: "copy startup-config running-config", Explanation: "Menyalin kembali konfigurasi asli yang tersimpan di NVRAM ke memori kerja RAM tanpa menghapusnya."},
				{Command: "config-register 0x2102", Explanation: "WAJIB: Mengembalikan nilai register ke standar normal (0x2102) agar router kembali membaca startup-config saat dinyalakan berikutnya."},
			},
			Verification:        "Ketik 'show version'. Pada baris paling akhir, pastikan tertulis: 'Configuration register is 0x2102'. Coba uji login enable menggunakan password baru.",
			TroubleshootingTips: "JANGAN LUPA mengembalikan config-register ke '0x2102'! Jika lupa, setiap kali router dimatikan dan dinyalakan, konfigurasi akan selalu tampak hilang kembali ke setelan pabrik.",
			Tags:                []string{"password recovery", "rommon", "0x2142", "0x2102", "lupa password router"},
		},
		{
			ID:          "recovery-backup-tftp",
			Title:       "Pencadangan & Pemulihan: Backup / Restore Konfigurasi Cisco ke Server TFTP",
			Device:      DeviceRouter,
			Mode:        ModePrivExec,
			Category:    CategoryRecovery,
			Description: "Mencadangkan file konfigurasi dan file citra IOS (.bin) dari router/switch ke server TFTP di jaringan LAN serta cara memulihkannya kembali saat perangkat rusak.",
			Parameters: []Parameter{
				{Key: "TFTP_SERVER_IP", Label: "Alamat IP Server TFTP", DefaultValue: "192.168.1.252", Placeholder: "cth: 192.168.1.252"},
				{Key: "BACKUP_FILENAME", Label: "Nama Berkas Cadangan", DefaultValue: "R1-Pusat-Config-Backup.cfg", Placeholder: "cth: R1-Backup.cfg"},
			},
			IPExample: "TFTP Server: {{TFTP_SERVER_IP}} | File Backup: {{BACKUP_FILENAME}}",
			Commands: `# 1. CARA BACKUP RUNNING-CONFIG KE SERVER TFTP:
Router# copy running-config tftp:
Address or name of remote host []? {{TFTP_SERVER_IP}}
Destination filename [Router-confg]? {{BACKUP_FILENAME}}
# Muncul: 'Writing running-config...!! [OK - 1250 bytes]'

# 2. CARA RESTORE / MEMULIHKAN CONFIG DARI TFTP KE ROUTER:
Router# copy tftp: running-config
Address or name of remote host []? {{TFTP_SERVER_IP}}
Source filename []? {{BACKUP_FILENAME}}
Destination filename [running-config]? 
# Muncul: 'Loading {{BACKUP_FILENAME}} from {{TFTP_SERVER_IP}}... [OK]'

# 3. CARA BACKUP CITRA SISTEM OPERASI CISCO IOS (FLASH KE TFTP):
Router# copy flash: tftp:`,
			Explanation: []LineExplanation{
				{Command: "copy running-config tftp:", Explanation: "Menyalin konfigurasi aktif di RAM langsung ke file di server TFTP."},
				{Command: "copy tftp: running-config", Explanation: "Mengunduh file konfigurasi dari server TFTP dan langsung menerapkannya ke perangkat."},
				{Command: "copy flash: tftp:", Explanation: "Menyimpan file binari IOS Cisco (.bin) ke server untuk cadangan jika sistem operasi corrupt."},
			},
			Verification:        "Buka Server TFTP di Packet Tracer, buka tab 'Services' -> 'TFTP'. Pastikan nama file '{{BACKUP_FILENAME}}' terdaftar di daftar berkas server.",
			TroubleshootingTips: "Uji koneksi ping ke IP {{TFTP_SERVER_IP}} sebelum menjalankan perintah 'copy'. Di Packet Tracer, pastikan Service TFTP pada perangkat Server berstatus 'On'.",
			Tags:                []string{"tftp", "backup config", "restore", "disaster recovery", "ios backup"},
		},

		// ====================================================================
		// PC & END DEVICE COMMAND PROMPT SUITE
		// ====================================================================
		{
			ID:          "diag-pc-terminal-suite",
			Title:       "Perintah Command Prompt PC di Cisco Packet Tracer (ipconfig, ping, tracert, arp)",
			Device:      DevicePC,
			Mode:        ModePCTerminal,
			Category:    CategoryShowDiag,
			Description: "Koleksi lengkap perintah terminal command prompt pada PC/Laptop di Cisco Packet Tracer untuk memeriksa alamat IP, uji koneksi, melacak hop rute, dan melihat tabel ARP.",
			IPExample:   "Dijalankan di Desktop -> Command Prompt (PC> / Command Prompt)",
			Commands: `# 1. LIHAT ALAMAT IP LENGKAP PC, SUBNET MASK, GATEWAY, & MAC:
ipconfig /all

# 2. MINTA ALAMAT IP BARU DARI SERVER DHCP:
ipconfig /release
ipconfig /renew

# 3. UJI SAMBUNGAN KONEKSI JARINGAN KE TUJUAN:
ping 192.168.1.1
# Ping tanpa henti (hanya berhenti jika ditekan Ctrl + C):
ping -t 192.168.1.1

# 4. LACAK JALUR LOMPATAN ROUTER (TRACEROUTE):
tracert 192.168.20.10

# 5. LIHAT TABEL PEMETAAN IP KE MAC ADDRESS (ARP CACHE):
arp -a
# Hapus cache ARP lokal:
arp -d

# 6. UJI RESOLUSI NAMA DOMAIN KE SERVER DNS:
nslookup www.cisco.com

# 7. REMOT CLI KE ROUTER/SWITCH LEWAT PC:
telnet 192.168.1.1
ssh -l admin 192.168.1.1`,
			Explanation: []LineExplanation{
				{Command: "ipconfig /all", Explanation: "Menampilkan konfigurasi antarmuka FastEthernet PC secara terperinci termasuk IP, Mask, Gateway, DNS, dan Alamat Fisik MAC."},
				{Command: "ipconfig /renew", Explanation: "Mengirim paket DHCP Request untuk meminta konfigurasi IP baru dari DHCP Server di jaringan."},
				{Command: "ping 192.168.1.1", Explanation: "Mengirim 4 paket ICMP Echo Request untuk menguji konektivitas dua arah dengan komputer/router tujuan."},
				{Command: "tracert 192.168.20.10", Explanation: "Melacak setiap hop IP router yang dilewati paket data hingga mencapai alamat tujuan akhir."},
				{Command: "arp -a", Explanation: "Menampilkan tabel korespondensi antara alamat IP logikal dengan alamat fisik kartu jaringan (MAC Address)."},
				{Command: "nslookup www.cisco.com", Explanation: "Mengirim kueri nama domain ke DNS server yang dikonfigurasi pada PC."},
				{Command: "ssh -l admin 192.168.1.1", Explanation: "Membuka sesi terminal jarak jauh terenkripsi SSH ke alamat IP router manajemen."},
			},
			Verification:        "Semua output muncul langsung di layar Command Prompt PC Packet Tracer.",
			TroubleshootingTips: "Di Command Prompt Packet Tracer, jika ping terhenti atau tidak merespons, tekan kombinasi 'Ctrl + C' untuk membatalkan proses.",
			Tags:                []string{"pc commands", "ipconfig", "ping", "tracert", "arp -a", "nslookup", "end device"},
		},
	}
}
