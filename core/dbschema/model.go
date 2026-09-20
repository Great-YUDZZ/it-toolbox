package dbschema

import "strings"

// DatabaseEngine represents a supported Relational Database Management System (RDBMS)
type DatabaseEngine string

const (
	EngineAll       DatabaseEngine = "Semua Engine"
	EngineMySQL     DatabaseEngine = "MySQL"
	EngineMariaDB   DatabaseEngine = "MariaDB"
	EnginePostgres  DatabaseEngine = "PostgreSQL"
	EngineSQLite    DatabaseEngine = "SQLite"
	EngineSQLServer DatabaseEngine = "SQL Server (T-SQL)"
	EngineOracle    DatabaseEngine = "Oracle Database"
)

// CommandCategory represents the category of database administration or DDL commands
type CommandCategory string

const (
	CatAll            CommandCategory = "Semua Kategori"
	CatDatabaseSetup  CommandCategory = "Pembuatan & Pengaturan Database"
	CatUserPrivileges CommandCategory = "Manajemen Pengguna & Hak Akses"
	CatTableDDL       CommandCategory = "Pembuatan Tabel & Constraint"
	CatBackupRestore  CommandCategory = "Backup, Dump & Restore CLI"
	CatMaintenance    CommandCategory = "Pemeliharaan, Index & Optimasi"
)

// Parameter represents a customizable parameter in a SQL template (e.g. DBNAME, USERNAME, PASSWORD)
type Parameter struct {
	Key          string `json:"key"`           // Placeholder key, e.g. "DB_NAME", "DB_USER"
	Label        string `json:"label"`         // User-friendly label, e.g. "Nama Basis Data"
	DefaultValue string `json:"default_value"` // Default standard value, e.g. "app_production"
	Placeholder  string `json:"placeholder"`   // Hint, e.g. "cth: db_toko_online"
}

// DBCommandItem represents a single database administration or DDL command recipe
type DBCommandItem struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Engine      DatabaseEngine  `json:"engine"`
	Category    CommandCategory `json:"category"`
	Description string          `json:"description"`
	Parameters  []Parameter     `json:"parameters,omitempty"`
	CommandSQL  string          `json:"command_sql"`
	Notes       string          `json:"notes,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
}

// RenderCommand replaces placeholders with custom values or defaults
func (item DBCommandItem) RenderCommand(values map[string]string) string {
	res := item.CommandSQL
	for _, p := range item.Parameters {
		val := ""
		if values != nil {
			val = strings.TrimSpace(values[p.Key])
		}
		if val == "" {
			val = p.DefaultValue
		}
		res = strings.ReplaceAll(res, "{{" + p.Key + "}}", val)
	}
	return res
}

// ColumnType defines common column data types for the interactive builder
type ColumnType string

const (
	TypeAutoID   ColumnType = "ID Auto-Increment (PK)"
	TypeUUID     ColumnType = "UUID / GUID (PK)"
	TypeVarchar  ColumnType = "VARCHAR (Teks Pendek)"
	TypeText     ColumnType = "TEXT (Teks Panjang)"
	TypeInteger  ColumnType = "INTEGER (Bilangan Bulat)"
	TypeBigInt   ColumnType = "BIGINT (Angka Besar)"
	TypeDecimal  ColumnType = "DECIMAL (Uang / Pecahan)"
	TypeBoolean  ColumnType = "BOOLEAN (Benar/Salah)"
	TypeDateTime ColumnType = "DATETIME / TIMESTAMPTZ (Waktu)"
	TypeDate     ColumnType = "DATE (Tanggal)"
	TypeJSON     ColumnType = "JSON (Data Terstruktur)"
)

// ColumnDef represents a column definition in the interactive table builder
type ColumnDef struct {
	Name         string     `json:"name"`
	Type         ColumnType `json:"type"`
	Length       string     `json:"length,omitempty"`        // e.g. "255" or "10,2"
	IsPrimaryKey bool       `json:"is_primary_key"`
	IsNotNull    bool       `json:"is_not_null"`
	IsUnique     bool       `json:"is_unique"`
	DefaultValue string     `json:"default_value,omitempty"` // e.g. "CURRENT_TIMESTAMP" or "'active'"
	Comment      string     `json:"comment,omitempty"`
}

// TableDef represents a table schema created in the interactive builder
type TableDef struct {
	DatabaseName string         `json:"database_name"`
	TableName    string         `json:"table_name"`
	Engine       DatabaseEngine `json:"engine"`
	Columns      []ColumnDef    `json:"columns"`
	IncludeDrop  bool           `json:"include_drop"`
}

// SchemaTemplate represents a complete ready-to-use multi-table schema architecture
type SchemaTemplate struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Category    string                    `json:"category"`
	Tables      []string                  `json:"tables"`
	SQLScripts  map[DatabaseEngine]string `json:"sql_scripts"`
	Notes       string                    `json:"notes"`
}

// RenderScript renders the template SQL for a given engine with custom parameters
func (st SchemaTemplate) RenderScript(engine DatabaseEngine, dbName, tablePrefix string) string {
	sql, ok := st.SQLScripts[engine]
	if !ok {
		sql = st.SQLScripts[EngineMySQL]
	}
	if dbName == "" {
		dbName = "it_toolbox_db"
	}
	res := strings.ReplaceAll(sql, "{{DB_NAME}}", dbName)
	res = strings.ReplaceAll(res, "{{PREFIX}}", tablePrefix)
	return res
}

// DataTypeComparison represents a cross-engine mapping row
type DataTypeComparison struct {
	LogicalType string `json:"logical_type"`
	Description string `json:"description"`
	MySQL       string `json:"mysql"`
	MariaDB     string `json:"mariadb"`
	PostgreSQL  string `json:"postgresql"`
	SQLite      string `json:"sqlite"`
	SQLServer   string `json:"sql_server"`
	Oracle      string `json:"oracle"`
	BestUse     string `json:"best_use"`
}
