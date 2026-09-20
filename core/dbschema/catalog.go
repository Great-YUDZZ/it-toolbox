package dbschema

import "strings"

// GetCommandCatalog returns the complete catalog of database administration & DDL commands
func GetCommandCatalog() []DBCommandItem {
	return []DBCommandItem{
		// ====================================================================
		// 1. DATABASE SETUP
		// ====================================================================
		{
			ID:          "mysql_create_db",
			Title:       "Membuat Database Standar UTF8mb4 (MySQL / MariaDB)",
			Engine:      EngineMySQL,
			Category:    CatDatabaseSetup,
			Description: "Membuat database baru dengan charset utf8mb4 dan collation unicode agar mendukung seluruh karakter internasional dan emoji secara sempurna.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: db_portal_sekolah"},
			},
			CommandSQL: "CREATE DATABASE IF NOT EXISTS `{{DB_NAME}}`\n" +
				"    CHARACTER SET utf8mb4\n" +
				"    COLLATE utf8mb4_unicode_ci;\n\n" +
				"-- Verifikasi pembuatan basis data\n" +
				"SHOW CREATE DATABASE `{{DB_NAME}}`;",
			Notes: "utf8mb4 adalah standar wajib modern untuk MySQL/MariaDB menggantikan utf8 lama yang hanya mendukung 3-byte.",
			Tags:  []string{"create", "database", "utf8mb4", "charset", "collation"},
		},
		{
			ID:          "postgres_create_db",
			Title:       "Membuat Database dengan Encoding UTF8 & Owner (PostgreSQL)",
			Engine:      EnginePostgres,
			Category:    CatDatabaseSetup,
			Description: "Membuat database baru di PostgreSQL dengan kepemilikan role aplikasi, encoding UTF8, dan template standar.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: db_portal"},
				{Key: "DB_OWNER", Label: "Pemilik (Owner Role)", DefaultValue: "app_user", Placeholder: "cth: postgres"},
			},
			CommandSQL: "CREATE DATABASE {{DB_NAME}}\n" +
				"    WITH\n" +
				"    OWNER = {{DB_OWNER}}\n" +
				"    ENCODING = 'UTF8'\n" +
				"    LC_COLLATE = 'en_US.UTF-8'\n" +
				"    LC_CTYPE = 'en_US.UTF-8'\n" +
				"    TABLESPACE = pg_default\n" +
				"    CONNECTION LIMIT = -1;\n\n" +
				"-- Menghubungkan sesi CLI ke database baru\n" +
				"\\c {{DB_NAME}};",
			Notes: "Pastikan role pemilik (owner) telah dibuat terlebih dahulu sebelum mengeksekusi CREATE DATABASE.",
			Tags:  []string{"create", "database", "postgres", "utf8", "owner"},
		},
		{
			ID:          "sqlite_setup_wal",
			Title:       "Inisialisasi Database & Mode WAL Performa Tinggi (SQLite)",
			Engine:      EngineSQLite,
			Category:    CatDatabaseSetup,
			Description: "Mengaktifkan Write-Ahead Logging (WAL) dan Foreign Key Constraints pada SQLite untuk konkurensi baca/tulis tinggi.",
			Parameters: []Parameter{
				{Key: "DB_FILE", Label: "Nama Berkas SQLite", DefaultValue: "aplikasi.db", Placeholder: "cth: production.sqlite3"},
			},
			CommandSQL: "-- Buka atau inisialisasi berkas via sqlite3 terminal:\n" +
				"-- sqlite3 {{DB_FILE}}\n\n" +
				"-- Aktifkan mode WAL untuk konkurensi baca-tulis tinggi\n" +
				"PRAGMA journal_mode = WAL;\n\n" +
				"-- Aktifkan penegakan integritas Foreign Key (default SQLite adalah OFF)\n" +
				"PRAGMA foreign_keys = ON;\n\n" +
				"-- Optimalkan sinkronisasi disk\n" +
				"PRAGMA synchronous = NORMAL;\n\n" +
				"-- Periksa integritas berkas basis data\n" +
				"PRAGMA integrity_check;",
			Notes: "SQLite memerlukan PRAGMA foreign_keys = ON; di setiap awal koneksi aplikasi agar relasi antar tabel ditegakkan.",
			Tags:  []string{"sqlite", "wal", "foreign_keys", "pragma", "performance"},
		},
		{
			ID:          "mssql_create_db",
			Title:       "Membuat Database Enterprise dengan Collation UTF-8 (SQL Server)",
			Engine:      EngineSQLServer,
			Category:    CatDatabaseSetup,
			Description: "Membuat database di Microsoft SQL Server dengan pengaturan file data, file log, dan collation UTF-8 modern.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "AplikasiDB", Placeholder: "cth: SalesDB"},
			},
			CommandSQL: "CREATE DATABASE [{{DB_NAME}}]\n" +
				"COLLATE Latin1_General_100_CI_AS_SC_UTF8;\n" +
				"GO\n\n" +
				"-- Mengatur Recovery Model ke SIMPLE (untuk dev) atau FULL (untuk prod)\n" +
				"ALTER DATABASE [{{DB_NAME}}] SET RECOVERY FULL;\n" +
				"GO\n\n" +
				"-- Mengaktifkan READ COMMITTED SNAPSHOT untuk konkurensi tinggi\n" +
				"ALTER DATABASE [{{DB_NAME}}] SET READ_COMMITTED_SNAPSHOT ON;\n" +
				"GO",
			Notes: "Collation _SC_UTF8 tersedia pada SQL Server 2019 ke atas untuk penyimpanan UTF-8 hemat ruang.",
			Tags:  []string{"mssql", "sqlserver", "create database", "collation", "snapshot"},
		},
		{
			ID:          "db_drop_safety",
			Title:       "Menghapus Database dengan Aman (Drop Database)",
			Engine:      EngineMySQL,
			Category:    CatDatabaseSetup,
			Description: "Perintah drop database dengan pengecekan keberadaan (IF EXISTS) untuk skrip otomasi/migrasi.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_lama_tidak_terpakai", Placeholder: "cth: db_staging_old"},
			},
			CommandSQL: "-- PERINGATAN: Perintah ini menghapus database beserta seluruh tabel & data di dalamnya!\n" +
				"DROP DATABASE IF EXISTS `{{DB_NAME}}`;",
			Notes: "Selalu buat cadangan (dump) terlebih dahulu sebelum menjalankan perintah DROP di lingkungan produksi.",
			Tags:  []string{"drop", "delete", "cleanup", "database"},
		},

		// ====================================================================
		// 2. USER PRIVILEGES & SECURITY
		// ====================================================================
		{
			ID:          "mysql_create_app_user",
			Title:       "Membuat User Aplikasi & Memberikan Hak Akses Penuh (MySQL/MariaDB)",
			Engine:      EngineMySQL,
			Category:    CatUserPrivileges,
			Description: "Membuat user baru untuk backend aplikasi dan memberikan hak akses penuh hanya pada database yang ditargetkan.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: db_toko"},
				{Key: "DB_USER", Label: "Nama Pengguna", DefaultValue: "app_user", Placeholder: "cth: backend_svc"},
				{Key: "DB_HOST", Label: "Host Akses (% untuk remote / localhost)", DefaultValue: "%", Placeholder: "localhost"},
				{Key: "DB_PASS", Label: "Kata Sandi User", DefaultValue: "KatasandiKuat_123!", Placeholder: "password_rahasia"},
			},
			CommandSQL: "-- 1. Buat pengguna baru dengan autentikasi kata sandi yang kuat\n" +
				"CREATE USER IF NOT EXISTS '{{DB_USER}}'@'{{DB_HOST}}' IDENTIFIED BY '{{DB_PASS}}';\n\n" +
				"-- 2. Berikan seluruh hak akses operasi DDL & DML pada database target\n" +
				"GRANT ALL PRIVILEGES ON `{{DB_NAME}}`.* TO '{{DB_USER}}'@'{{DB_HOST}}';\n\n" +
				"-- 3. Terapkan perubahan hak akses ke memori\n" +
				"FLUSH PRIVILEGES;\n\n" +
				"-- 4. Verifikasi daftar hak akses yang baru diberikan\n" +
				"SHOW GRANTS FOR '{{DB_USER}}'@'{{DB_HOST}}';",
			Notes: "Gunakan host 'localhost' jika database dan aplikasi backend berada pada satu VPS/server fisik yang sama.",
			Tags:  []string{"user", "grant", "privileges", "mysql", "security"},
		},
		{
			ID:          "mysql_create_readonly_user",
			Title:       "Membuat User Read-Only / Analitik (MySQL/MariaDB)",
			Engine:      EngineMySQL,
			Category:    CatUserPrivileges,
			Description: "Membuat user khusus pelaporan (Reporting / BI / Read-Only) yang hanya diizinkan mengeksekusi perintah SELECT.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: db_toko"},
				{Key: "RO_USER", Label: "Nama Pengguna Read-Only", DefaultValue: "bi_readonly", Placeholder: "cth: analyst_user"},
				{Key: "RO_PASS", Label: "Kata Sandi", DefaultValue: "AnalystPassword_2026!", Placeholder: "password_rahasia"},
			},
			CommandSQL: "CREATE USER IF NOT EXISTS '{{RO_USER}}'@'%' IDENTIFIED BY '{{RO_PASS}}';\n\n" +
				"-- Berikan hanya izin SELECT (tidak bisa INSERT, UPDATE, DELETE, atau DROP)\n" +
				"GRANT SELECT ON `{{DB_NAME}}`.* TO '{{RO_USER}}'@'%';\n\n" +
				"FLUSH PRIVILEGES;",
			Notes: "Sangat direkomendasikan untuk koneksi tools seperti Metabase, Tableau, Grafana, atau PowerBI.",
			Tags:  []string{"readonly", "select only", "bi", "reporting", "security"},
		},
		{
			ID:          "postgres_user_permissions",
			Title:       "Membuat Role & Hak Akses Schema Public Lengkap (PostgreSQL)",
			Engine:      EnginePostgres,
			Category:    CatUserPrivileges,
			Description: "Membuat role baru di PostgreSQL serta mengatur izin akses database dan skema public hingga tabel masa depan.",
			Parameters: []Parameter{
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: my_postgres_db"},
				{Key: "APP_USER", Label: "Nama Pengguna (Role)", DefaultValue: "app_backend", Placeholder: "cth: webapp"},
				{Key: "APP_PASS", Label: "Kata Sandi", DefaultValue: "SecretPostgresPass_123!", Placeholder: "password_rahasia"},
			},
			CommandSQL: "-- 1. Buat role pengguna dengan izin login\n" +
				"CREATE ROLE {{APP_USER}} WITH LOGIN PASSWORD '{{APP_PASS}}';\n\n" +
				"-- 2. Izinkan koneksi ke database target\n" +
				"GRANT CONNECT ON DATABASE {{DB_NAME}} TO {{APP_USER}};\n\n" +
				"-- 3. Hubungkan ke database target\n" +
				"\\c {{DB_NAME}};\n\n" +
				"-- 4. Berikan akses ke schema public\n" +
				"GRANT USAGE, CREATE ON SCHEMA public TO {{APP_USER}};\n\n" +
				"-- 5. Berikan akses penuh ke seluruh tabel & sequence yang sudah ada\n" +
				"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO {{APP_USER}};\n" +
				"GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO {{APP_USER}};\n\n" +
				"-- 6. Atur default privilege agar tabel yang dibuat di masa depan otomatis bisa diakses\n" +
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO {{APP_USER}};\n" +
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO {{APP_USER}};",
			Notes: "ALTER DEFAULT PRIVILEGES wajib diatur di PostgreSQL agar user tidak mengalami permission denied saat migrasi tabel baru.",
			Tags:  []string{"postgres", "role", "permissions", "schema public", "sequences"},
		},

		// ====================================================================
		// 3. TABLE DDL & CONSTRAINTS
		// ====================================================================
		{
			ID:          "mysql_foreign_key_cascade",
			Title:       "Pola Relasi Foreign Key dengan ON DELETE CASCADE (MySQL)",
			Engine:      EngineMySQL,
			Category:    CatTableDDL,
			Description: "Contoh pembuatan tabel induk dan tabel anak berelasi 1-to-Many dengan integritas referensial CASCADE.",
			Parameters: []Parameter{
				{Key: "PARENT_TABLE", Label: "Tabel Induk", DefaultValue: "kategori", Placeholder: "cth: departemen"},
				{Key: "CHILD_TABLE", Label: "Tabel Anak", DefaultValue: "produk", Placeholder: "cth: karyawan"},
			},
			CommandSQL: "-- Tabel Induk\n" +
				"CREATE TABLE IF NOT EXISTS `{{PARENT_TABLE}}` (\n" +
				"    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,\n" +
				"    `nama_kategori` VARCHAR(100) NOT NULL,\n" +
				"    `slug` VARCHAR(120) NOT NULL UNIQUE,\n" +
				"    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n\n" +
				"-- Tabel Anak yang merujuk ke Tabel Induk\n" +
				"CREATE TABLE IF NOT EXISTS `{{CHILD_TABLE}}` (\n" +
				"    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,\n" +
				"    `kategori_id` BIGINT UNSIGNED NOT NULL,\n" +
				"    `nama_item` VARCHAR(255) NOT NULL,\n" +
				"    `harga` DECIMAL(12,2) NOT NULL DEFAULT 0.00,\n" +
				"    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
				"    \n" +
				"    -- Index untuk mempercepat query JOIN\n" +
				"    INDEX `idx_{{CHILD_TABLE}}_kategori` (`kategori_id`),\n" +
				"    \n" +
				"    -- Foreign Key Constraint\n" +
				"    CONSTRAINT `fk_{{CHILD_TABLE}}_{{PARENT_TABLE}}`\n" +
				"        FOREIGN KEY (`kategori_id`)\n" +
				"        REFERENCES `{{PARENT_TABLE}}` (`id`)\n" +
				"        ON DELETE CASCADE\n" +
				"        ON UPDATE CASCADE\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
			Notes: "ON DELETE CASCADE menghapus seluruh data anak otomatis saat baris induk dihapus. Gunakan RESTRICT atau SET NULL jika data anak tidak boleh ikut terhapus.",
			Tags:  []string{"foreign key", "cascade", "join index", "one to many", "ddl"},
		},
		{
			ID:          "postgres_uuid_pk",
			Title:       "Tabel dengan UUID Primary Key & Auto Timestamp Trigger (PostgreSQL)",
			Engine:      EnginePostgres,
			Category:    CatTableDDL,
			Description: "Tabel standar industri dengan UUID v4 otomatis, JSONB kolom dinamis, dan trigger update_timestamp otomatis.",
			Parameters: []Parameter{
				{Key: "TABLE_NAME", Label: "Nama Tabel", DefaultValue: "akun_pengguna", Placeholder: "cth: master_pelanggan"},
			},
			CommandSQL: "-- Aktifkan modul pgcrypto / gen_random_uuid bawaan\n" +
				"CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";\n\n" +
				"-- Fungsi Trigger untuk memperbarui kolom updated_at otomatis\n" +
				"CREATE OR REPLACE FUNCTION trigger_set_timestamp()\n" +
				"RETURNS TRIGGER AS $$\n" +
				"BEGIN\n" +
				"    NEW.updated_at = NOW();\n" +
				"    RETURN NEW;\n" +
				"END;\n" +
				"$$ LANGUAGE plpgsql;\n\n" +
				"-- Definisi Tabel\n" +
				"CREATE TABLE IF NOT EXISTS {{TABLE_NAME}} (\n" +
				"    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
				"    email VARCHAR(255) NOT NULL UNIQUE,\n" +
				"    nama_lengkap VARCHAR(150) NOT NULL,\n" +
				"    metadata JSONB DEFAULT '{}'::jsonb,\n" +
				"    is_active BOOLEAN NOT NULL DEFAULT TRUE,\n" +
				"    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
				"    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP\n" +
				");\n\n" +
				"-- Index GIN untuk pencarian isi data JSONB super cepat\n" +
				"CREATE INDEX IF NOT EXISTS idx_{{TABLE_NAME}}_metadata ON {{TABLE_NAME}} USING GIN (metadata);\n\n" +
				"-- Pasang trigger pada tabel\n" +
				"DROP TRIGGER IF EXISTS set_timestamp_{{TABLE_NAME}} ON {{TABLE_NAME}};\n" +
				"CREATE TRIGGER set_timestamp_{{TABLE_NAME}}\n" +
				"    BEFORE UPDATE ON {{TABLE_NAME}}\n" +
				"    FOR EACH ROW\n" +
				"    EXECUTE PROCEDURE trigger_set_timestamp();",
			Notes: "Trigger otomatis diperlukan di PostgreSQL karena tidak ada fitur bawaan 'ON UPDATE CURRENT_TIMESTAMP' seperti di MySQL.",
			Tags:  []string{"uuid", "postgres", "trigger", "jsonb", "gin index", "timestamp"},
		},

		// ====================================================================
		// 4. BACKUP & RESTORE CLI
		// ====================================================================
		{
			ID:          "mysqldump_backup_commands",
			Title:       "Perintah Backup Database Lengkap via mysqldump (CLI)",
			Engine:      EngineMySQL,
			Category:    CatBackupRestore,
			Description: "Kumpulan perintah terminal Linux / Windows Command Prompt untuk mencadangkan database MySQL/MariaDB.",
			Parameters: []Parameter{
				{Key: "DB_USER", Label: "User Database", DefaultValue: "root", Placeholder: "root"},
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "cth: db_toko"},
				{Key: "BACKUP_FILE", Label: "Nama Berkas Cadangan", DefaultValue: "cadangan_db.sql", Placeholder: "backup.sql"},
			},
			CommandSQL: "# 1. Backup penuh database (Struktur + Data) dengan kompresi gzip:\n" +
				"mysqldump -u {{DB_USER}} -p --single-transaction --quick --routines --triggers {{DB_NAME}} | gzip > {{BACKUP_FILE}}.gz\n\n" +
				"# 2. Backup hanya skema/struktur tabel saja (tanpa data):\n" +
				"mysqldump -u {{DB_USER}} -p --no-data --routines --triggers {{DB_NAME}} > skema_{{DB_NAME}}.sql\n\n" +
				"# 3. Backup hanya data saja (tanpa DDL create table):\n" +
				"mysqldump -u {{DB_USER}} -p --no-create-info {{DB_NAME}} > data_{{DB_NAME}}.sql\n\n" +
				"# 4. Backup seluruh database sekaligus (All Databases):\n" +
				"mysqldump -u {{DB_USER}} -p --all-databases --single-transaction --quick > full_server_backup.sql",
			Notes: "--single-transaction mencegah tabel terkunci saat proses dump berjalan pada tabel ber-engine InnoDB.",
			Tags:  []string{"mysqldump", "backup", "restore", "cli", "gzip", "mariadb-dump"},
		},
		{
			ID:          "mysql_restore_command",
			Title:       "Perintah Restore Database SQL via Terminal (CLI)",
			Engine:      EngineMySQL,
			Category:    CatBackupRestore,
			Description: "Mengembalikan atau mengimpor berkas SQL cadangan kembali ke dalam database MySQL/MariaDB.",
			Parameters: []Parameter{
				{Key: "DB_USER", Label: "User Database", DefaultValue: "root", Placeholder: "root"},
				{Key: "DB_NAME", Label: "Nama Database Tujuan", DefaultValue: "db_aplikasi_utama", Placeholder: "db_target"},
				{Key: "SQL_FILE", Label: "Berkas SQL", DefaultValue: "cadangan_db.sql", Placeholder: "backup.sql"},
			},
			CommandSQL: "# 1. Mengimpor file SQL biasa ke database tujuan:\n" +
				"mysql -u {{DB_USER}} -p {{DB_NAME}} < {{SQL_FILE}}\n\n" +
				"# 2. Jika file backup terkompresi .gz, ekstrak langsung ke pipe mysql:\n" +
				"gunzip < {{SQL_FILE}}.gz | mysql -u {{DB_USER}} -p {{DB_NAME}}\n\n" +
				"# 3. Mode verbose untuk memantau baris yang dieksekusi:\n" +
				"mysql -u {{DB_USER}} -p --verbose {{DB_NAME}} < {{SQL_FILE}}",
			Notes: "Pastikan database tujuan sudah dibuat terlebih dahulu jika di dalam berkas .sql tidak terdapat perintah CREATE DATABASE.",
			Tags:  []string{"restore", "import", "mysql", "gunzip", "terminal"},
		},
		{
			ID:          "pg_dump_backup_restore",
			Title:       "Perintah Backup & Restore Modern pg_dump (PostgreSQL CLI)",
			Engine:      EnginePostgres,
			Category:    CatBackupRestore,
			Description: "Sintaks pg_dump format kustom binary (.dump) berkecepatan tinggi dan proses restore menggunakan pg_restore.",
			Parameters: []Parameter{
				{Key: "DB_USER", Label: "User PostgreSQL", DefaultValue: "postgres", Placeholder: "postgres"},
				{Key: "DB_NAME", Label: "Nama Database", DefaultValue: "db_aplikasi_utama", Placeholder: "mydb"},
				{Key: "DUMP_FILE", Label: "Nama File Dump", DefaultValue: "backup_pg.dump", Placeholder: "backup.dump"},
			},
			CommandSQL: "# 1. Backup format kustom terkompresi (Format -Fc paling direkomendasikan):\n" +
				"pg_dump -U {{DB_USER}} -d {{DB_NAME}} -F c -b -v -f {{DUMP_FILE}}\n\n" +
				"# 2. Restore file dump format kustom ke database baru:\n" +
				"pg_restore -U {{DB_USER}} -d {{DB_NAME}} -v --clean --if-exists {{DUMP_FILE}}\n\n" +
				"# 3. Backup skrip SQL teks biasa (plain text .sql):\n" +
				"pg_dump -U {{DB_USER}} -d {{DB_NAME}} -F p -f {{DB_NAME}}_plain.sql\n\n" +
				"# 4. Restore file SQL plain text melalui psql:\n" +
				"psql -U {{DB_USER}} -d {{DB_NAME}} -f {{DB_NAME}}_plain.sql",
			Notes: "Format binary (-F c) memungkinkan pg_restore berjalan secara multi-thread dengan parameter -j 4 untuk restore sangat cepat.",
			Tags:  []string{"pg_dump", "pg_restore", "psql", "postgres", "backup"},
		},

		// ====================================================================
		// 5. MAINTENANCE & DIAGNOSTICS
		// ====================================================================
		{
			ID:          "mysql_maintenance_diag",
			Title:       "Pemeliharaan, Optimasi Tabel & Analisis Query (MySQL)",
			Engine:      EngineMySQL,
			Category:    CatMaintenance,
			Description: "Kumpulan perintah diagnostik status tabel, defragmentasi ruang kosong, serta analisis rencana eksekusi query (EXPLAIN).",
			Parameters: []Parameter{
				{Key: "TABLE_NAME", Label: "Nama Tabel", DefaultValue: "transaksi", Placeholder: "cth: pesanan"},
			},
			CommandSQL: "-- 1. Menganalisis distribusi statistik indeks tabel\n" +
				"ANALYZE TABLE `{{TABLE_NAME}}`;\n\n" +
				"-- 2. Mengoptimalkan tabel dan mereklamasi ruang penyimpanan yang terbuang\n" +
				"OPTIMIZE TABLE `{{TABLE_NAME}}`;\n\n" +
				"-- 3. Memeriksa adanya kerusakan pada struktur tabel\n" +
				"CHECK TABLE `{{TABLE_NAME}}` EXTENDED;\n\n" +
				"-- 4. Analisis performa query dan indeks yang digunakan (Execution Plan)\n" +
				"EXPLAIN FORMAT=JSON\n" +
				"SELECT * FROM `{{TABLE_NAME}}` WHERE `id` = 1;",
			Notes: "OPTIMIZE TABLE merekonstruksi tabel dan indeks di latar belakang. Lakukan pada jam lalu lintas rendah.",
			Tags:  []string{"optimize", "analyze", "explain", "check table", "maintenance"},
		},
		{
			ID:          "postgres_vacuum_analyze",
			Title:       "Vacuuming, Reindexing & Buffer Analysis (PostgreSQL)",
			Engine:      EnginePostgres,
			Category:    CatMaintenance,
			Description: "Pembersihan dead tuples di PostgreSQL (VACUUM), regenerasi indeks, dan inspeksi buffer eksekusi query.",
			Parameters: []Parameter{
				{Key: "TABLE_NAME", Label: "Nama Tabel", DefaultValue: "transaksi", Placeholder: "cth: order_items"},
			},
			CommandSQL: "-- 1. Bersihkan dead tuples dan perbarui statistik perencana query\n" +
				"VACUUM (VERBOSE, ANALYZE) {{TABLE_NAME}};\n\n" +
				"-- 2. Buat ulang indeks tabel yang mengalami fragmentasi atau pembengkakan\n" +
				"REINDEX TABLE CONCURRENTLY {{TABLE_NAME}};\n\n" +
				"-- 3. Analisis mendalam query dengan laporan penggunaan memori buffer & waktu aktual\n" +
				"EXPLAIN (ANALYZE, BUFFERS, VERBOSE)\n" +
				"SELECT * FROM {{TABLE_NAME}} ORDER BY id DESC LIMIT 20;",
			Notes: "Gunakan REINDEX TABLE CONCURRENTLY di lingkungan produksi agar pembacaan dan penulisan tabel tidak terblokir.",
			Tags:  []string{"vacuum", "reindex", "explain analyze", "buffers", "postgres"},
		},
	}
}

// SearchCommands filters command items by engine, category, and query string
func SearchCommands(engine DatabaseEngine, category CommandCategory, query string) []DBCommandItem {
	all := GetCommandCatalog()
	q := strings.TrimSpace(strings.ToLower(query))

	var results []DBCommandItem
	for _, item := range all {
		if engine != "" && engine != EngineAll {
			if item.Engine != engine && item.Engine != EngineAll {
				if !(engine == EngineMariaDB && item.Engine == EngineMySQL) {
					continue
				}
			}
		}

		if category != "" && category != CatAll {
			if item.Category != category {
				continue
			}
		}

		if q != "" {
			match := strings.Contains(strings.ToLower(item.Title), q) ||
				strings.Contains(strings.ToLower(item.Description), q) ||
				strings.Contains(strings.ToLower(item.CommandSQL), q)
			if !match {
				for _, t := range item.Tags {
					if strings.Contains(strings.ToLower(t), q) {
						match = true
						break
					}
				}
			}
			if !match {
				continue
			}
		}

		results = append(results, item)
	}
	return results
}
