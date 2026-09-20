package dbschema

import "strings"

// GetDataTypesMatrix returns the comprehensive cross-database data types comparison table
func GetDataTypesMatrix() []DataTypeComparison {
	return []DataTypeComparison{
		{
			LogicalType: "ID Auto-Increment (PK)",
			Description: "Kunci utama bertambah otomatis untuk baris tabel",
			MySQL:       "BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			MariaDB:     "BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY",
			PostgreSQL:  "BIGSERIAL PRIMARY KEY",
			SQLite:      "INTEGER PRIMARY KEY AUTOINCREMENT",
			SQLServer:   "BIGINT IDENTITY(1,1) PRIMARY KEY",
			Oracle:      "NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY",
			BestUse:     "Primary key default yang cepat diindeks dan hemat memori",
		},
		{
			LogicalType: "UUID / GUID",
			Description: "Pengidentifikasi unik global 128-bit (aman untuk multi-server & distributed)",
			MySQL:       "CHAR(36) atau BINARY(16)",
			MariaDB:     "UUID",
			PostgreSQL:  "UUID DEFAULT gen_random_uuid()",
			SQLite:      "TEXT (36 karakter)",
			SQLServer:   "UNIQUEIDENTIFIER DEFAULT NEWID()",
			Oracle:      "RAW(16) DEFAULT SYS_GUID()",
			BestUse:     "Microservices, sistem terdistribusi, dan API publik tanpa mengekspos ID berurutan",
		},
		{
			LogicalType: "Teks Pendek (Nama, Kode, Judul)",
			Description: "Untaian karakter bervariasi dengan batas panjang maksimal",
			MySQL:       "VARCHAR(n) [misal: VARCHAR(255)]",
			MariaDB:     "VARCHAR(n)",
			PostgreSQL:  "VARCHAR(n) atau TEXT",
			SQLite:      "TEXT",
			SQLServer:   "NVARCHAR(n)",
			Oracle:      "VARCHAR2(n CHAR)",
			BestUse:     "Nama pengguna, email, judul artikel, slug, kode pos",
		},
		{
			LogicalType: "Teks Panjang (Deskripsi, Artikel)",
			Description: "Teks berukuran besar tanpa batasan karakter kaku",
			MySQL:       "TEXT / MEDIUMTEXT / LONGTEXT",
			MariaDB:     "TEXT / LONGTEXT",
			PostgreSQL:  "TEXT",
			SQLite:      "TEXT",
			SQLServer:   "NVARCHAR(MAX)",
			Oracle:      "CLOB",
			BestUse:     "Isi blog, konten markdown, deskripsi produk, log sistem",
		},
		{
			LogicalType: "Bilangan Bulat Standar",
			Description: "Integer 32-bit (rentang -2.1 miliar s/d +2.1 miliar)",
			MySQL:       "INT atau INT SIGNED",
			MariaDB:     "INT",
			PostgreSQL:  "INTEGER",
			SQLite:      "INTEGER",
			SQLServer:   "INT",
			Oracle:      "NUMBER(10)",
			BestUse:     "Jumlah stok, urutan antrean, hitungan kuantitas",
		},
		{
			LogicalType: "Bilangan Bulat Besar",
			Description: "Integer 64-bit untuk data volume raksasa",
			MySQL:       "BIGINT",
			MariaDB:     "BIGINT",
			PostgreSQL:  "BIGINT",
			SQLite:      "INTEGER",
			SQLServer:   "BIGINT",
			Oracle:      "NUMBER(19)",
			BestUse:     "Nomor telepon, view count, saldo poin, id transaksi berskala besar",
		},
		{
			LogicalType: "Mata Uang & Angka Presisi",
			Description: "Desimal eksak dengan digit tetap di belakang koma (anti floating point error)",
			MySQL:       "DECIMAL(12,2) / DECIMAL(18,4)",
			MariaDB:     "DECIMAL(12,2)",
			PostgreSQL:  "NUMERIC(12,2)",
			SQLite:      "REAL atau INTEGER (simpan dalam sen)",
			SQLServer:   "DECIMAL(12,2) atau MONEY",
			Oracle:      "NUMBER(12,2)",
			BestUse:     "Harga barang, total pembayaran, persentase diskon, akuntansi",
		},
		{
			LogicalType: "Nilai Boolean (Ya/Tidak)",
			Description: "Status biner benar/salah atau aktif/non-aktif",
			MySQL:       "TINYINT(1) [0 = false, 1 = true]",
			MariaDB:     "BOOLEAN / TINYINT(1)",
			PostgreSQL:  "BOOLEAN (TRUE / FALSE)",
			SQLite:      "INTEGER (0 atau 1)",
			SQLServer:   "BIT (0 atau 1)",
			Oracle:      "NUMBER(1) [atau BOOLEAN di Oracle 23c+]",
			BestUse:     "is_active, is_verified, is_admin, has_paid",
		},
		{
			LogicalType: "Tanggal Saja (Tanpa Jam)",
			Description: "Informasi kalender tahun-bulan-hari (YYYY-MM-DD)",
			MySQL:       "DATE",
			MariaDB:     "DATE",
			PostgreSQL:  "DATE",
			SQLite:      "TEXT (format ISO 8601 YYYY-MM-DD)",
			SQLServer:   "DATE",
			Oracle:      "DATE",
			BestUse:     "Tanggal lahir, tanggal jatuh tempo, tanggal libur",
		},
		{
			LogicalType: "Waktu & Zona Waktu (Timestamp)",
			Description: "Tanggal dan jam presisi tinggi disertai dukungan zona waktu UTC",
			MySQL:       "DATETIME / TIMESTAMP",
			MariaDB:     "DATETIME / TIMESTAMP",
			PostgreSQL:  "TIMESTAMPTZ (TIMESTAMP WITH TIME ZONE)",
			SQLite:      "TEXT DEFAULT CURRENT_TIMESTAMP",
			SQLServer:   "DATETIMEOFFSET atau DATETIME2",
			Oracle:      "TIMESTAMP WITH TIME ZONE",
			BestUse:     "Audit log, created_at, updated_at, jadwal booking",
		},
		{
			LogicalType: "Dokumen JSON / Data Fleksibel",
			Description: "Struktur data semi-terstruktur JSON dengan validasi sintaks",
			MySQL:       "JSON",
			MariaDB:     "JSON (alias untuk LONGTEXT dengan cek JSON)",
			PostgreSQL:  "JSONB (Binary JSON, terindeks cepat)",
			SQLite:      "TEXT (didukung fungsi JSON)",
			SQLServer:   "NVARCHAR(MAX) (dengan ISJSON check)",
			Oracle:      "JSON / CLOB (CHECK IS JSON)",
			BestUse:     "Metadata konfigurasi, atribut dinamis produk, payload webhook",
		},
		{
			LogicalType: "Berkas Biner (BLOB / Gambar / File)",
			Description: "Byte mentah untuk menyimpan berkas atau tanda tangan digital",
			MySQL:       "BLOB / LONGBLOB",
			MariaDB:     "LONGBLOB",
			PostgreSQL:  "BYTEA",
			SQLite:      "BLOB",
			SQLServer:   "VARBINARY(MAX)",
			Oracle:      "BLOB",
			BestUse:     "Kunci kriptografi, sertifikat, berkas kecil (disarankan URL untuk file besar)",
		},
		{
			LogicalType: "Pilihan Terbatas (Enum / Status)",
			Description: "Kumpulan nilai statis yang sudah ditentukan sebelumnya",
			MySQL:       "ENUM('pending', 'proses', 'selesai')",
			MariaDB:     "ENUM('pending', 'proses', 'selesai')",
			PostgreSQL:  "CREATE TYPE status_enum AS ENUM(...) atau CHECK",
			SQLite:      "TEXT CHECK(col IN ('pending', 'proses', 'selesai'))",
			SQLServer:   "VARCHAR(20) CHECK(col IN (...))",
			Oracle:      "VARCHAR2(20) CHECK(col IN (...))",
			BestUse:     "Status pesanan, peran hak akses, metode pengiriman",
		},
	}
}

// SearchDataTypes searches comparative data types by query
func SearchDataTypes(query string) []DataTypeComparison {
	all := GetDataTypesMatrix()
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return all
	}

	var results []DataTypeComparison
	for _, item := range all {
		if strings.Contains(strings.ToLower(item.LogicalType), q) ||
			strings.Contains(strings.ToLower(item.Description), q) ||
			strings.Contains(strings.ToLower(item.MySQL), q) ||
			strings.Contains(strings.ToLower(item.PostgreSQL), q) ||
			strings.Contains(strings.ToLower(item.BestUse), q) {
			results = append(results, item)
		}
	}
	return results
}
