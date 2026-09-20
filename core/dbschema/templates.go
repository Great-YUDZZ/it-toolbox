package dbschema

// GetSchemaTemplates returns the curated collection of production-ready schema templates
func GetSchemaTemplates() []SchemaTemplate {
	return []SchemaTemplate{
		// ====================================================================
		// 1. AUTH & RBAC (Role-Based Access Control)
		// ====================================================================
		{
			ID:          "schema_auth_rbac",
			Title:       "Sistem Autentikasi Pengguna & RBAC",
			Description: "Arsitektur keamanan enterprise untuk autentikasi user, role hierarkis, perizinan modular (permissions), session token, dan audit log pelacakan aktivitas.",
			Category:    "Keamanan & Akses",
			Tables:      []string{"roles", "permissions", "role_permissions", "users", "user_roles", "user_sessions", "audit_logs"},
			Notes:       "Dilengkapi foreign key cascading dan indexing pada email, username, dan session_token untuk response login instan.",
			SQLScripts: map[DatabaseEngine]string{
				EngineMySQL: `-- ========================================================
-- Schema: Sistem Autentikasi & RBAC
-- Target: MySQL / MariaDB (InnoDB, UTF8mb4)
-- Basis Data: {{DB_NAME}}
-- ========================================================

CREATE DATABASE IF NOT EXISTS ` + "`{{DB_NAME}}`" + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ` + "`{{DB_NAME}}`" + `;

-- 1. Tabel Peran (Roles)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}roles`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nama_role`" + ` VARCHAR(50) NOT NULL UNIQUE,
    ` + "`deskripsi`" + ` VARCHAR(255) NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Tabel Hak Izin (Permissions)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}permissions`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`kode_izin`" + ` VARCHAR(80) NOT NULL UNIQUE,
    ` + "`nama_izin`" + ` VARCHAR(150) NOT NULL,
    ` + "`modul`" + ` VARCHAR(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Pivot Tabel Role - Permissions
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}role_permissions`" + ` (
    ` + "`role_id`" + ` INT UNSIGNED NOT NULL,
    ` + "`permission_id`" + ` INT UNSIGNED NOT NULL,
    PRIMARY KEY (` + "`role_id`, `permission_id`" + `),
    CONSTRAINT ` + "`fk_rp_role`" + ` FOREIGN KEY (` + "`role_id`" + `) REFERENCES ` + "`{{PREFIX}}roles` (`id`)" + ` ON DELETE CASCADE,
    CONSTRAINT ` + "`fk_rp_perm`" + ` FOREIGN KEY (` + "`permission_id`" + `) REFERENCES ` + "`{{PREFIX}}permissions` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. Tabel Akun Pengguna (Users)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}users`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`username`" + ` VARCHAR(50) NOT NULL UNIQUE,
    ` + "`email`" + ` VARCHAR(120) NOT NULL UNIQUE,
    ` + "`password_hash`" + ` VARCHAR(255) NOT NULL,
    ` + "`nama_lengkap`" + ` VARCHAR(150) NOT NULL,
    ` + "`is_active`" + ` TINYINT(1) NOT NULL DEFAULT 1,
    ` + "`last_login_at`" + ` DATETIME NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ` + "`updated_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX ` + "`idx_users_email` (`email`)" + `
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. Pivot Tabel User - Roles (Many to Many)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}user_roles`" + ` (
    ` + "`user_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`role_id`" + ` INT UNSIGNED NOT NULL,
    PRIMARY KEY (` + "`user_id`, `role_id`" + `),
    CONSTRAINT ` + "`fk_ur_user`" + ` FOREIGN KEY (` + "`user_id`" + `) REFERENCES ` + "`{{PREFIX}}users` (`id`)" + ` ON DELETE CASCADE,
    CONSTRAINT ` + "`fk_ur_role`" + ` FOREIGN KEY (` + "`role_id`" + `) REFERENCES ` + "`{{PREFIX}}roles` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 6. Tabel Sesi Login (User Sessions)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}user_sessions`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`user_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`session_token`" + ` VARCHAR(255) NOT NULL UNIQUE,
    ` + "`ip_address`" + ` VARCHAR(45) NULL,
    ` + "`user_agent`" + ` TEXT NULL,
    ` + "`expires_at`" + ` DATETIME NOT NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX ` + "`idx_session_token` (`session_token`)" + `,
    CONSTRAINT ` + "`fk_sess_user`" + ` FOREIGN KEY (` + "`user_id`" + `) REFERENCES ` + "`{{PREFIX}}users` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 7. Tabel Log Audit Aktivitas (Audit Logs)
CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}audit_logs`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`user_id`" + ` BIGINT UNSIGNED NULL,
    ` + "`aksi`" + ` VARCHAR(100) NOT NULL,
    ` + "`tabel_target`" + ` VARCHAR(80) NULL,
    ` + "`rekaman_id`" + ` VARCHAR(100) NULL,
    ` + "`perubahan_json`" + ` JSON NULL,
    ` + "`ip_address`" + ` VARCHAR(45) NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX ` + "`idx_audit_user` (`user_id`)" + `,
    INDEX ` + "`idx_audit_created` (`created_at`)" + `
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

				EnginePostgres: `-- ========================================================
-- Schema: Sistem Autentikasi & RBAC
-- Target: PostgreSQL
-- ========================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS {{PREFIX}}roles (
    id SERIAL PRIMARY KEY,
    nama_role VARCHAR(50) NOT NULL UNIQUE,
    deskripsi VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}permissions (
    id SERIAL PRIMARY KEY,
    kode_izin VARCHAR(80) NOT NULL UNIQUE,
    nama_izin VARCHAR(150) NOT NULL,
    modul VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}role_permissions (
    role_id INT NOT NULL REFERENCES {{PREFIX}}roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES {{PREFIX}}permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(120) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nama_lengkap VARCHAR(150) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_{{PREFIX}}users_email ON {{PREFIX}}users(email);

CREATE TABLE IF NOT EXISTS {{PREFIX}}user_roles (
    user_id BIGINT NOT NULL REFERENCES {{PREFIX}}users(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES {{PREFIX}}roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES {{PREFIX}}users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_{{PREFIX}}sess_token ON {{PREFIX}}user_sessions(session_token);

CREATE TABLE IF NOT EXISTS {{PREFIX}}audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES {{PREFIX}}users(id) ON DELETE SET NULL,
    aksi VARCHAR(100) NOT NULL,
    tabel_target VARCHAR(80),
    rekaman_id VARCHAR(100),
    perubahan_json JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_{{PREFIX}}audit_created ON {{PREFIX}}audit_logs(created_at);`,

				EngineSQLite: `-- ========================================================
-- Schema: Sistem Autentikasi & RBAC
-- Target: SQLite (Mode WAL & Foreign Keys)
-- ========================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS {{PREFIX}}roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_role TEXT NOT NULL UNIQUE,
    deskripsi TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kode_izin TEXT NOT NULL UNIQUE,
    nama_izin TEXT NOT NULL,
    modul TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES {{PREFIX}}roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES {{PREFIX}}permissions (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    nama_lengkap TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    last_login_at TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES {{PREFIX}}users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES {{PREFIX}}roles (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    ip_address TEXT,
    user_agent TEXT,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES {{PREFIX}}users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    aksi TEXT NOT NULL,
    tabel_target TEXT,
    rekaman_id TEXT,
    perubahan_json TEXT,
    ip_address TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES {{PREFIX}}users (id) ON DELETE SET NULL
);`,
			},
		},

		// ====================================================================
		// 2. E-COMMERCE & TOKO ONLINE
		// ====================================================================
		{
			ID:          "schema_ecommerce",
			Title:       "E-Commerce & Toko Online",
			Description: "Struktur data transaksi komersial lengkap: Pelanggan, Kategori Berjenjang, Katalog Produk, Stok, Keranjang, Pesanan (Orders), Rincian Item, dan Riwayat Pembayaran.",
			Category:    "Bisnis & Transaksi",
			Tables:      []string{"customers", "categories", "products", "orders", "order_items", "payments", "shippings"},
			Notes:       "Menggunakan tipe data DECIMAL(12,2) untuk akurasi nominal mata uang anti-pembulatan desimal.",
			SQLScripts: map[DatabaseEngine]string{
				EngineMySQL: `-- ========================================================
-- Schema: E-Commerce & Toko Online
-- Target: MySQL / MariaDB (InnoDB, UTF8mb4)
-- ========================================================

CREATE DATABASE IF NOT EXISTS ` + "`{{DB_NAME}}`" + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ` + "`{{DB_NAME}}`" + `;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}customers`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nama`" + ` VARCHAR(150) NOT NULL,
    ` + "`email`" + ` VARCHAR(120) NOT NULL UNIQUE,
    ` + "`telepon`" + ` VARCHAR(25) NULL,
    ` + "`alamat_utama`" + ` TEXT NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}categories`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`parent_id`" + ` INT UNSIGNED NULL,
    ` + "`nama_kategori`" + ` VARCHAR(100) NOT NULL,
    ` + "`slug`" + ` VARCHAR(120) NOT NULL UNIQUE,
    CONSTRAINT ` + "`fk_cat_parent`" + ` FOREIGN KEY (` + "`parent_id`" + `) REFERENCES ` + "`{{PREFIX}}categories` (`id`)" + ` ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}products`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`category_id`" + ` INT UNSIGNED NULL,
    ` + "`sku`" + ` VARCHAR(50) NOT NULL UNIQUE,
    ` + "`nama_produk`" + ` VARCHAR(255) NOT NULL,
    ` + "`slug`" + ` VARCHAR(255) NOT NULL UNIQUE,
    ` + "`harga`" + ` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    ` + "`stok`" + ` INT NOT NULL DEFAULT 0,
    ` + "`berat_gram`" + ` INT NOT NULL DEFAULT 100,
    ` + "`is_active`" + ` TINYINT(1) NOT NULL DEFAULT 1,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX ` + "`idx_prod_cat` (`category_id`)" + `,
    CONSTRAINT ` + "`fk_prod_cat`" + ` FOREIGN KEY (` + "`category_id`" + `) REFERENCES ` + "`{{PREFIX}}categories` (`id`)" + ` ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}orders`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nomor_invoice`" + ` VARCHAR(60) NOT NULL UNIQUE,
    ` + "`customer_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`total_belanja`" + ` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    ` + "`ongkos_kirim`" + ` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    ` + "`grand_total`" + ` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    ` + "`status_pesanan`" + ` ENUM('menunggu_pembayaran', 'diproses', 'dikirim', 'selesai', 'dibatalkan') NOT NULL DEFAULT 'menunggu_pembayaran',
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX ` + "`idx_order_cust` (`customer_id`)" + `,
    CONSTRAINT ` + "`fk_order_cust`" + ` FOREIGN KEY (` + "`customer_id`" + `) REFERENCES ` + "`{{PREFIX}}customers` (`id`)" + ` ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}order_items`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`order_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`product_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`kuantitas`" + ` INT NOT NULL DEFAULT 1,
    ` + "`harga_satuan`" + ` DECIMAL(12,2) NOT NULL,
    ` + "`subtotal`" + ` DECIMAL(12,2) NOT NULL,
    CONSTRAINT ` + "`fk_item_order`" + ` FOREIGN KEY (` + "`order_id`" + `) REFERENCES ` + "`{{PREFIX}}orders` (`id`)" + ` ON DELETE CASCADE,
    CONSTRAINT ` + "`fk_item_prod`" + ` FOREIGN KEY (` + "`product_id`" + `) REFERENCES ` + "`{{PREFIX}}products` (`id`)" + ` ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}payments`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`order_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`metode_bayar`" + ` VARCHAR(50) NOT NULL,
    ` + "`jumlah_bayar`" + ` DECIMAL(12,2) NOT NULL,
    ` + "`status_bayar`" + ` ENUM('pending', 'berhasil', 'gagal', 'kadaluarsa') NOT NULL DEFAULT 'pending',
    ` + "`waktu_bayar`" + ` DATETIME NULL,
    CONSTRAINT ` + "`fk_pay_order`" + ` FOREIGN KEY (` + "`order_id`" + `) REFERENCES ` + "`{{PREFIX}}orders` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

				EnginePostgres: `-- ========================================================
-- Schema: E-Commerce & Toko Online
-- Target: PostgreSQL
-- ========================================================

CREATE TABLE IF NOT EXISTS {{PREFIX}}customers (
    id BIGSERIAL PRIMARY KEY,
    nama VARCHAR(150) NOT NULL,
    email VARCHAR(120) NOT NULL UNIQUE,
    telepon VARCHAR(25),
    alamat_utama TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}categories (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES {{PREFIX}}categories(id) ON DELETE SET NULL,
    nama_kategori VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}products (
    id BIGSERIAL PRIMARY KEY,
    category_id INT REFERENCES {{PREFIX}}categories(id) ON DELETE SET NULL,
    sku VARCHAR(50) NOT NULL UNIQUE,
    nama_produk VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    harga NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    stok INT NOT NULL DEFAULT 0,
    berat_gram INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}orders (
    id BIGSERIAL PRIMARY KEY,
    nomor_invoice VARCHAR(60) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL REFERENCES {{PREFIX}}customers(id) ON DELETE RESTRICT,
    total_belanja NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    ongkos_kirim NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    grand_total NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    status_pesanan VARCHAR(30) NOT NULL DEFAULT 'menunggu_pembayaran',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES {{PREFIX}}orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES {{PREFIX}}products(id) ON DELETE RESTRICT,
    kuantitas INT NOT NULL DEFAULT 1,
    harga_satuan NUMERIC(12,2) NOT NULL,
    subtotal NUMERIC(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES {{PREFIX}}orders(id) ON DELETE CASCADE,
    metode_bayar VARCHAR(50) NOT NULL,
    jumlah_bayar NUMERIC(12,2) NOT NULL,
    status_bayar VARCHAR(20) NOT NULL DEFAULT 'pending',
    waktu_bayar TIMESTAMPTZ
);`,

				EngineSQLite: `-- ========================================================
-- Schema: E-Commerce & Toko Online
-- Target: SQLite
-- ========================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS {{PREFIX}}customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    telepon TEXT,
    alamat_utama TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    nama_kategori TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    FOREIGN KEY (parent_id) REFERENCES {{PREFIX}}categories (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER,
    sku TEXT NOT NULL UNIQUE,
    nama_produk TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    harga REAL NOT NULL DEFAULT 0.00,
    stok INTEGER NOT NULL DEFAULT 0,
    berat_gram INTEGER NOT NULL DEFAULT 100,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES {{PREFIX}}categories (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor_invoice TEXT NOT NULL UNIQUE,
    customer_id INTEGER NOT NULL,
    total_belanja REAL NOT NULL DEFAULT 0.00,
    ongkos_kirim REAL NOT NULL DEFAULT 0.00,
    grand_total REAL NOT NULL DEFAULT 0.00,
    status_pesanan TEXT NOT NULL DEFAULT 'menunggu_pembayaran',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES {{PREFIX}}customers (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    kuantitas INTEGER NOT NULL DEFAULT 1,
    harga_satuan REAL NOT NULL,
    subtotal REAL NOT NULL,
    FOREIGN KEY (order_id) REFERENCES {{PREFIX}}orders (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES {{PREFIX}}products (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    metode_bayar TEXT NOT NULL,
    jumlah_bayar REAL NOT NULL,
    status_bayar TEXT NOT NULL DEFAULT 'pending',
    waktu_bayar TEXT,
    FOREIGN KEY (order_id) REFERENCES {{PREFIX}}orders (id) ON DELETE CASCADE
);`,
			},
		},

		// ====================================================================
		// 3. SISTEM INFORMASI AKADEMIK / LAB TKJ
		// ====================================================================
		{
			ID:          "schema_akademik_tkj",
			Title:       "Sistem Akademik Sekolah & Lab TKJ",
			Description: "Skema relasional lengkap untuk institusi pendidikan SMK/TKJ: Data Siswa (NISN), Guru (NIP), Rombel/Kelas, Mata Pelajaran, Penilaian, Presensi, dan Inventaris Aset Lab Komputer.",
			Category:    "Pendidikan & Sekolah",
			Tables:      []string{"guru", "kelas", "siswa", "mata_pelajaran", "jadwal_pelajaran", "nilai_siswa", "aset_lab_komputer"},
			Notes:       "Dilengkapi tabel aset_lab_komputer untuk mencatat PC Lab, nomor switch, port LAN, dan status hardware.",
			SQLScripts: map[DatabaseEngine]string{
				EngineMySQL: `-- ========================================================
-- Schema: Sistem Akademik & Lab TKJ
-- Target: MySQL / MariaDB (InnoDB, UTF8mb4)
-- ========================================================

CREATE DATABASE IF NOT EXISTS ` + "`{{DB_NAME}}`" + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ` + "`{{DB_NAME}}`" + `;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}guru`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nip`" + ` VARCHAR(30) NOT NULL UNIQUE,
    ` + "`nama_guru`" + ` VARCHAR(150) NOT NULL,
    ` + "`keahlian`" + ` VARCHAR(100) NULL DEFAULT 'Teknik Komputer & Jaringan',
    ` + "`telepon`" + ` VARCHAR(25) NULL,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}kelas`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nama_kelas`" + ` VARCHAR(50) NOT NULL UNIQUE,
    ` + "`tingkat`" + ` ENUM('X', 'XI', 'XII') NOT NULL,
    ` + "`wali_guru_id`" + ` INT UNSIGNED NULL,
    CONSTRAINT ` + "`fk_kelas_wali`" + ` FOREIGN KEY (` + "`wali_guru_id`" + `) REFERENCES ` + "`{{PREFIX}}guru` (`id`)" + ` ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}siswa`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`nisn`" + ` VARCHAR(20) NOT NULL UNIQUE,
    ` + "`nama_siswa`" + ` VARCHAR(150) NOT NULL,
    ` + "`kelas_id`" + ` INT UNSIGNED NOT NULL,
    ` + "`jenis_kelamin`" + ` ENUM('L', 'P') NOT NULL,
    ` + "`status_aktif`" + ` TINYINT(1) NOT NULL DEFAULT 1,
    ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ` + "`fk_siswa_kelas`" + ` FOREIGN KEY (` + "`kelas_id`" + `) REFERENCES ` + "`{{PREFIX}}kelas` (`id`)" + ` ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}mata_pelajaran`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`kode_mapel`" + ` VARCHAR(20) NOT NULL UNIQUE,
    ` + "`nama_mapel`" + ` VARCHAR(100) NOT NULL,
    ` + "`kkm`" + ` INT NOT NULL DEFAULT 75
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}nilai_siswa`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`siswa_id`" + ` BIGINT UNSIGNED NOT NULL,
    ` + "`mapel_id`" + ` INT UNSIGNED NOT NULL,
    ` + "`nilai_tugas`" + ` DECIMAL(5,2) DEFAULT 0,
    ` + "`nilai_uts`" + ` DECIMAL(5,2) DEFAULT 0,
    ` + "`nilai_uas`" + ` DECIMAL(5,2) DEFAULT 0,
    ` + "`nilai_akhir`" + ` DECIMAL(5,2) DEFAULT 0,
    CONSTRAINT ` + "`fk_nilai_siswa`" + ` FOREIGN KEY (` + "`siswa_id`" + `) REFERENCES ` + "`{{PREFIX}}siswa` (`id`)" + ` ON DELETE CASCADE,
    CONSTRAINT ` + "`fk_nilai_mapel`" + ` FOREIGN KEY (` + "`mapel_id`" + `) REFERENCES ` + "`{{PREFIX}}mata_pelajaran` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}aset_lab_komputer`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`kode_aset`" + ` VARCHAR(50) NOT NULL UNIQUE,
    ` + "`lab_ruang`" + ` VARCHAR(50) NOT NULL DEFAULT 'Lab TKJ 1',
    ` + "`nama_perangkat`" + ` VARCHAR(100) NOT NULL,
    ` + "`ip_address`" + ` VARCHAR(45) NULL,
    ` + "`mac_address`" + ` VARCHAR(17) NULL,
    ` + "`kondisi`" + ` ENUM('baik', 'perlu_perbaikan', 'rusak_total') NOT NULL DEFAULT 'baik',
    ` + "`keterangan`" + ` TEXT NULL,
    ` + "`updated_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

				EnginePostgres: `-- ========================================================
-- Schema: Sistem Akademik & Lab TKJ
-- Target: PostgreSQL
-- ========================================================

CREATE TABLE IF NOT EXISTS {{PREFIX}}guru (
    id SERIAL PRIMARY KEY,
    nip VARCHAR(30) NOT NULL UNIQUE,
    nama_guru VARCHAR(150) NOT NULL,
    keahlian VARCHAR(100) DEFAULT 'Teknik Komputer & Jaringan',
    telepon VARCHAR(25),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}kelas (
    id SERIAL PRIMARY KEY,
    nama_kelas VARCHAR(50) NOT NULL UNIQUE,
    tingkat VARCHAR(10) NOT NULL,
    wali_guru_id INT REFERENCES {{PREFIX}}guru(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}siswa (
    id BIGSERIAL PRIMARY KEY,
    nisn VARCHAR(20) NOT NULL UNIQUE,
    nama_siswa VARCHAR(150) NOT NULL,
    kelas_id INT NOT NULL REFERENCES {{PREFIX}}kelas(id) ON DELETE RESTRICT,
    jenis_kelamin CHAR(1) NOT NULL,
    status_aktif BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}mata_pelajaran (
    id SERIAL PRIMARY KEY,
    kode_mapel VARCHAR(20) NOT NULL UNIQUE,
    nama_mapel VARCHAR(100) NOT NULL,
    kkm INT NOT NULL DEFAULT 75
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}nilai_siswa (
    id BIGSERIAL PRIMARY KEY,
    siswa_id BIGINT NOT NULL REFERENCES {{PREFIX}}siswa(id) ON DELETE CASCADE,
    mapel_id INT NOT NULL REFERENCES {{PREFIX}}mata_pelajaran(id) ON DELETE CASCADE,
    nilai_tugas NUMERIC(5,2) DEFAULT 0,
    nilai_uts NUMERIC(5,2) DEFAULT 0,
    nilai_uas NUMERIC(5,2) DEFAULT 0,
    nilai_akhir NUMERIC(5,2) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}aset_lab_komputer (
    id SERIAL PRIMARY KEY,
    kode_aset VARCHAR(50) NOT NULL UNIQUE,
    lab_ruang VARCHAR(50) NOT NULL DEFAULT 'Lab TKJ 1',
    nama_perangkat VARCHAR(100) NOT NULL,
    ip_address VARCHAR(45),
    mac_address VARCHAR(17),
    kondisi VARCHAR(30) NOT NULL DEFAULT 'baik',
    keterangan TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,

				EngineSQLite: `-- ========================================================
-- Schema: Sistem Akademik & Lab TKJ
-- Target: SQLite
-- ========================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS {{PREFIX}}guru (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nip TEXT NOT NULL UNIQUE,
    nama_guru TEXT NOT NULL,
    keahlian TEXT DEFAULT 'Teknik Komputer & Jaringan',
    telepon TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}kelas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_kelas TEXT NOT NULL UNIQUE,
    tingkat TEXT NOT NULL,
    wali_guru_id INTEGER,
    FOREIGN KEY (wali_guru_id) REFERENCES {{PREFIX}}guru (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}siswa (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nisn TEXT NOT NULL UNIQUE,
    nama_siswa TEXT NOT NULL,
    kelas_id INTEGER NOT NULL,
    jenis_kelamin TEXT NOT NULL,
    status_aktif INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (kelas_id) REFERENCES {{PREFIX}}kelas (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}mata_pelajaran (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kode_mapel TEXT NOT NULL UNIQUE,
    nama_mapel TEXT NOT NULL,
    kkm INTEGER NOT NULL DEFAULT 75
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}nilai_siswa (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    siswa_id INTEGER NOT NULL,
    mapel_id INTEGER NOT NULL,
    nilai_tugas REAL DEFAULT 0,
    nilai_uts REAL DEFAULT 0,
    nilai_uas REAL DEFAULT 0,
    nilai_akhir REAL DEFAULT 0,
    FOREIGN KEY (siswa_id) REFERENCES {{PREFIX}}siswa (id) ON DELETE CASCADE,
    FOREIGN KEY (mapel_id) REFERENCES {{PREFIX}}mata_pelajaran (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}aset_lab_komputer (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kode_aset TEXT NOT NULL UNIQUE,
    lab_ruang TEXT NOT NULL DEFAULT 'Lab TKJ 1',
    nama_perangkat TEXT NOT NULL,
    ip_address TEXT,
    mac_address TEXT,
    kondisi TEXT NOT NULL DEFAULT 'baik',
    keterangan TEXT,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
			},
		},

		// ====================================================================
		// 4. IT NETWORK INFRASTRUCTURE & IPAM
		// ====================================================================
		{
			ID:          "schema_network_ipam",
			Title:       "Manajemen Infrastruktur Jaringan & IPAM",
			Description: "Struktur basis data untuk Network Engineer & Sysadmin: Manajemen Gedung/Ruang, Switch/Router, VLAN, Alokasi Subnet CIDR, Mapping IP Address & Port Switch.",
			Category:    "Jaringan & Server",
			Tables:      []string{"lokasi_gedung", "vlans", "subnets", "perangkat_jaringan", "ports_switch", "ip_allocations"},
			Notes:       "Sangat cocok untuk topologi jaringan kampus, kantor bertingkat, dan data center lokal.",
			SQLScripts: map[DatabaseEngine]string{
				EngineMySQL: `-- ========================================================
-- Schema: Manajemen Infrastruktur Jaringan & IPAM
-- Target: MySQL / MariaDB (InnoDB, UTF8mb4)
-- ========================================================

CREATE DATABASE IF NOT EXISTS ` + "`{{DB_NAME}}`" + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ` + "`{{DB_NAME}}`" + `;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}lokasi_gedung`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`kode_lokasi`" + ` VARCHAR(20) NOT NULL UNIQUE,
    ` + "`nama_gedung`" + ` VARCHAR(100) NOT NULL,
    ` + "`lantai`" + ` INT NOT NULL DEFAULT 1,
    ` + "`deskripsi`" + ` VARCHAR(255) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}vlans`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`vlan_id`" + ` INT NOT NULL UNIQUE,
    ` + "`nama_vlan`" + ` VARCHAR(50) NOT NULL,
    ` + "`deskripsi`" + ` VARCHAR(150) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}subnets`" + ` (
    ` + "`id`" + ` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`vlan_id`" + ` INT UNSIGNED NULL,
    ` + "`network_cidr`" + ` VARCHAR(45) NOT NULL UNIQUE,
    ` + "`gateway_ip`" + ` VARCHAR(45) NOT NULL,
    ` + "`netmask`" + ` VARCHAR(45) NOT NULL,
    ` + "`dns_server`" + ` VARCHAR(100) NULL DEFAULT '8.8.8.8, 1.1.1.1',
    CONSTRAINT ` + "`fk_subnet_vlan`" + ` FOREIGN KEY (` + "`vlan_id`" + `) REFERENCES ` + "`{{PREFIX}}vlans` (`id`)" + ` ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}perangkat_jaringan`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`lokasi_id`" + ` INT UNSIGNED NULL,
    ` + "`hostname`" + ` VARCHAR(80) NOT NULL UNIQUE,
    ` + "`tipe_perangkat`" + ` ENUM('Router', 'Switch L2', 'Switch L3', 'Access Point', 'Server', 'Firewall') NOT NULL,
    ` + "`merek_model`" + ` VARCHAR(100) NULL,
    ` + "`ip_manajemen`" + ` VARCHAR(45) NOT NULL UNIQUE,
    ` + "`status`" + ` ENUM('online', 'offline', 'maintenance') NOT NULL DEFAULT 'online',
    CONSTRAINT ` + "`fk_netdev_loc`" + ` FOREIGN KEY (` + "`lokasi_id`" + `) REFERENCES ` + "`{{PREFIX}}lokasi_gedung` (`id`)" + ` ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ` + "`{{PREFIX}}ip_allocations`" + ` (
    ` + "`id`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ` + "`subnet_id`" + ` INT UNSIGNED NOT NULL,
    ` + "`ip_address`" + ` VARCHAR(45) NOT NULL UNIQUE,
    ` + "`mac_address`" + ` VARCHAR(17) NULL,
    ` + "`hostname_terkait`" + ` VARCHAR(100) NULL,
    ` + "`status_alokasi`" + ` ENUM('dhcp_pool', 'statis_server', 'statis_network', 'reservasi') NOT NULL DEFAULT 'statis_server',
    ` + "`keterangan`" + ` VARCHAR(255) NULL,
    CONSTRAINT ` + "`fk_ip_sub`" + ` FOREIGN KEY (` + "`subnet_id`" + `) REFERENCES ` + "`{{PREFIX}}subnets` (`id`)" + ` ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

				EnginePostgres: `-- ========================================================
-- Schema: Manajemen Jaringan & IPAM
-- Target: PostgreSQL
-- ========================================================

CREATE TABLE IF NOT EXISTS {{PREFIX}}lokasi_gedung (
    id SERIAL PRIMARY KEY,
    kode_lokasi VARCHAR(20) NOT NULL UNIQUE,
    nama_gedung VARCHAR(100) NOT NULL,
    lantai INT NOT NULL DEFAULT 1,
    deskripsi VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}vlans (
    id SERIAL PRIMARY KEY,
    vlan_id INT NOT NULL UNIQUE,
    nama_vlan VARCHAR(50) NOT NULL,
    deskripsi VARCHAR(150)
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}subnets (
    id SERIAL PRIMARY KEY,
    vlan_id INT REFERENCES {{PREFIX}}vlans(id) ON DELETE SET NULL,
    network_cidr INET NOT NULL UNIQUE,
    gateway_ip INET NOT NULL,
    netmask VARCHAR(45) NOT NULL,
    dns_server VARCHAR(100) DEFAULT '8.8.8.8, 1.1.1.1'
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}perangkat_jaringan (
    id BIGSERIAL PRIMARY KEY,
    lokasi_id INT REFERENCES {{PREFIX}}lokasi_gedung(id) ON DELETE SET NULL,
    hostname VARCHAR(80) NOT NULL UNIQUE,
    tipe_perangkat VARCHAR(50) NOT NULL,
    merek_model VARCHAR(100),
    ip_manajemen INET NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'online'
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}ip_allocations (
    id BIGSERIAL PRIMARY KEY,
    subnet_id INT NOT NULL REFERENCES {{PREFIX}}subnets(id) ON DELETE CASCADE,
    ip_address INET NOT NULL UNIQUE,
    mac_address MACADDR,
    hostname_terkait VARCHAR(100),
    status_alokasi VARCHAR(40) NOT NULL DEFAULT 'statis_server',
    keterangan VARCHAR(255)
);`,

				EngineSQLite: `-- ========================================================
-- Schema: Manajemen Jaringan & IPAM
-- Target: SQLite
-- ========================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS {{PREFIX}}lokasi_gedung (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kode_lokasi TEXT NOT NULL UNIQUE,
    nama_gedung TEXT NOT NULL,
    lantai INTEGER NOT NULL DEFAULT 1,
    deskripsi TEXT
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}vlans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vlan_id INTEGER NOT NULL UNIQUE,
    nama_vlan TEXT NOT NULL,
    deskripsi TEXT
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}subnets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vlan_id INTEGER,
    network_cidr TEXT NOT NULL UNIQUE,
    gateway_ip TEXT NOT NULL,
    netmask TEXT NOT NULL,
    dns_server TEXT DEFAULT '8.8.8.8, 1.1.1.1',
    FOREIGN KEY (vlan_id) REFERENCES {{PREFIX}}vlans (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}perangkat_jaringan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lokasi_id INTEGER,
    hostname TEXT NOT NULL UNIQUE,
    tipe_perangkat TEXT NOT NULL,
    merek_model TEXT,
    ip_manajemen TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'online',
    FOREIGN KEY (lokasi_id) REFERENCES {{PREFIX}}lokasi_gedung (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS {{PREFIX}}ip_allocations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subnet_id INTEGER NOT NULL,
    ip_address TEXT NOT NULL UNIQUE,
    mac_address TEXT,
    hostname_terkait TEXT,
    status_alokasi TEXT NOT NULL DEFAULT 'statis_server',
    keterangan TEXT,
    FOREIGN KEY (subnet_id) REFERENCES {{PREFIX}}subnets (id) ON DELETE CASCADE
);`,
			},
		},
	}
}
