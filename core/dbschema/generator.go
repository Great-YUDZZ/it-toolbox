package dbschema

import (
	"fmt"
	"strings"
)

// GenerateTableDDL generates full DDL SQL statement for a given TableDef and target engine
func GenerateTableDDL(table TableDef) string {
	tableName := strings.TrimSpace(table.TableName)
	if tableName == "" {
		tableName = "contoh_tabel"
	}
	engine := table.Engine
	if engine == "" || engine == EngineAll {
		engine = EngineMySQL
	}

	var sb strings.Builder

	// Header comments
	sb.WriteString(fmt.Sprintf("-- ========================================================\n"))
	sb.WriteString(fmt.Sprintf("-- DDL Skema Tabel: %s\n", tableName))
	sb.WriteString(fmt.Sprintf("-- Target Database: %s\n", engine))
	if table.DatabaseName != "" {
		sb.WriteString(fmt.Sprintf("-- Basis Data: %s\n", table.DatabaseName))
	}
	sb.WriteString(fmt.Sprintf("-- Dihasilkan oleh: IT Toolbox Database Architect\n"))
	sb.WriteString(fmt.Sprintf("-- ========================================================\n\n"))

	// Switch database statement if specified
	if table.DatabaseName != "" {
		switch engine {
		case EngineMySQL, EngineMariaDB:
			sb.WriteString(fmt.Sprintf("USE `%s`;\n\n", table.DatabaseName))
		case EnginePostgres:
			sb.WriteString(fmt.Sprintf("\\c %s;\n\n", table.DatabaseName))
		case EngineSQLServer:
			sb.WriteString(fmt.Sprintf("USE [%s];\nGO\n\n", table.DatabaseName))
		}
	}

	// Drop statement if requested
	if table.IncludeDrop {
		switch engine {
		case EngineSQLServer:
			sb.WriteString(fmt.Sprintf("IF OBJECT_ID(N'dbo.%s', N'U') IS NOT NULL\n    DROP TABLE dbo.%s;\nGO\n\n", tableName, tableName))
		case EngineOracle:
			sb.WriteString(fmt.Sprintf("BEGIN\n    EXECUTE IMMEDIATE 'DROP TABLE %s CASCADE CONSTRAINTS';\nEXCEPTION\n    WHEN OTHERS THEN IF SQLCODE != -942 THEN RAISE; END IF;\nEND;\n/\n\n", strings.ToUpper(tableName)))
		default:
			sb.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\n\n", quoteIdentifier(tableName, engine)))
		}
	}

	// Create table statement
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", quoteIdentifier(tableName, engine)))

	cols := table.Columns
	if len(cols) == 0 {
		cols = PresetColumnsStandard()
	}

	var colLines []string
	var primaryKeys []string

	for _, col := range cols {
		colLine, isInlinePK := formatColumn(col, engine)
		colLines = append(colLines, "    "+colLine)
		if col.IsPrimaryKey && !isInlinePK {
			primaryKeys = append(primaryKeys, quoteIdentifier(col.Name, engine))
		}
	}

	// Add composite/table-level primary key if not defined inline
	if len(primaryKeys) > 0 {
		colLines = append(colLines, fmt.Sprintf("    PRIMARY KEY (%s)", strings.Join(primaryKeys, ", ")))
	}

	sb.WriteString(strings.Join(colLines, ",\n"))
	sb.WriteString("\n)")

	// Engine-specific table suffixes
	switch engine {
	case EngineMySQL, EngineMariaDB:
		sb.WriteString(" ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n")
	case EngineSQLServer:
		sb.WriteString(";\nGO\n")
	case EngineOracle:
		sb.WriteString(";\n")
	default:
		sb.WriteString(";\n")
	}

	return sb.String()
}

func quoteIdentifier(name string, engine DatabaseEngine) string {
	name = strings.TrimSpace(name)
	switch engine {
	case EngineMySQL, EngineMariaDB:
		return "`" + name + "`"
	case EnginePostgres, EngineOracle:
		return "\"" + name + "\""
	case EngineSQLServer:
		return "[" + name + "]"
	default:
		return "\"" + name + "\""
	}
}

func formatColumn(c ColumnDef, engine DatabaseEngine) (string, bool) {
	colName := quoteIdentifier(c.Name, engine)
	var typeStr string
	isInlinePK := false

	length := strings.TrimSpace(c.Length)
	if length == "" {
		if c.Type == TypeVarchar {
			length = "255"
		} else if c.Type == TypeDecimal {
			length = "12,2"
		}
	}

	switch engine {
	case EngineMySQL, EngineMariaDB:
		switch c.Type {
		case TypeAutoID:
			typeStr = "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY"
			isInlinePK = true
		case TypeUUID:
			typeStr = "CHAR(36) NOT NULL PRIMARY KEY"
			isInlinePK = true
		case TypeVarchar:
			typeStr = fmt.Sprintf("VARCHAR(%s)", length)
		case TypeText:
			typeStr = "TEXT"
		case TypeInteger:
			typeStr = "INT"
		case TypeBigInt:
			typeStr = "BIGINT"
		case TypeDecimal:
			typeStr = fmt.Sprintf("DECIMAL(%s)", length)
		case TypeBoolean:
			typeStr = "TINYINT(1)"
		case TypeDateTime:
			typeStr = "DATETIME"
		case TypeDate:
			typeStr = "DATE"
		case TypeJSON:
			typeStr = "JSON"
		default:
			typeStr = "VARCHAR(255)"
		}

	case EnginePostgres:
		switch c.Type {
		case TypeAutoID:
			typeStr = "BIGSERIAL PRIMARY KEY"
			isInlinePK = true
		case TypeUUID:
			typeStr = "UUID PRIMARY KEY DEFAULT gen_random_uuid()"
			isInlinePK = true
		case TypeVarchar:
			typeStr = fmt.Sprintf("VARCHAR(%s)", length)
		case TypeText:
			typeStr = "TEXT"
		case TypeInteger:
			typeStr = "INTEGER"
		case TypeBigInt:
			typeStr = "BIGINT"
		case TypeDecimal:
			typeStr = fmt.Sprintf("NUMERIC(%s)", length)
		case TypeBoolean:
			typeStr = "BOOLEAN"
		case TypeDateTime:
			typeStr = "TIMESTAMPTZ"
		case TypeDate:
			typeStr = "DATE"
		case TypeJSON:
			typeStr = "JSONB"
		default:
			typeStr = "VARCHAR(255)"
		}

	case EngineSQLite:
		switch c.Type {
		case TypeAutoID:
			typeStr = "INTEGER PRIMARY KEY AUTOINCREMENT"
			isInlinePK = true
		case TypeUUID:
			typeStr = "TEXT PRIMARY KEY"
			isInlinePK = true
		case TypeVarchar, TypeText:
			typeStr = "TEXT"
		case TypeInteger, TypeBigInt:
			typeStr = "INTEGER"
		case TypeDecimal:
			typeStr = "REAL"
		case TypeBoolean:
			typeStr = "INTEGER"
		case TypeDateTime, TypeDate:
			typeStr = "TEXT"
		case TypeJSON:
			typeStr = "TEXT"
		default:
			typeStr = "TEXT"
		}

	case EngineSQLServer:
		switch c.Type {
		case TypeAutoID:
			typeStr = "BIGINT IDENTITY(1,1) PRIMARY KEY"
			isInlinePK = true
		case TypeUUID:
			typeStr = "UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID()"
			isInlinePK = true
		case TypeVarchar:
			typeStr = fmt.Sprintf("NVARCHAR(%s)", length)
		case TypeText:
			typeStr = "NVARCHAR(MAX)"
		case TypeInteger:
			typeStr = "INT"
		case TypeBigInt:
			typeStr = "BIGINT"
		case TypeDecimal:
			typeStr = fmt.Sprintf("DECIMAL(%s)", length)
		case TypeBoolean:
			typeStr = "BIT"
		case TypeDateTime:
			typeStr = "DATETIME2"
		case TypeDate:
			typeStr = "DATE"
		case TypeJSON:
			typeStr = "NVARCHAR(MAX)"
		default:
			typeStr = "NVARCHAR(255)"
		}

	case EngineOracle:
		switch c.Type {
		case TypeAutoID:
			typeStr = "NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY"
			isInlinePK = true
		case TypeUUID:
			typeStr = "RAW(16) DEFAULT SYS_GUID() PRIMARY KEY"
			isInlinePK = true
		case TypeVarchar:
			typeStr = fmt.Sprintf("VARCHAR2(%s)", length)
		case TypeText:
			typeStr = "CLOB"
		case TypeInteger:
			typeStr = "NUMBER(10)"
		case TypeBigInt:
			typeStr = "NUMBER(19)"
		case TypeDecimal:
			typeStr = fmt.Sprintf("NUMBER(%s)", length)
		case TypeBoolean:
			typeStr = "NUMBER(1)"
		case TypeDateTime:
			typeStr = "TIMESTAMP WITH TIME ZONE"
		case TypeDate:
			typeStr = "DATE"
		case TypeJSON:
			typeStr = "CLOB CHECK (" + colName + " IS JSON)"
		default:
			typeStr = "VARCHAR2(255)"
		}
	}

	var parts []string
	parts = append(parts, colName, typeStr)

	// If already an inline primary key, skip constraints
	if isInlinePK {
		return strings.Join(parts, " "), true
	}

	if c.IsNotNull {
		parts = append(parts, "NOT NULL")
	}

	if c.IsUnique {
		parts = append(parts, "UNIQUE")
	}

	defVal := strings.TrimSpace(c.DefaultValue)
	if defVal != "" {
		if c.Type == TypeDateTime && strings.ToUpper(defVal) == "CURRENT_TIMESTAMP" {
			switch engine {
			case EngineSQLServer:
				parts = append(parts, "DEFAULT SYSUTCDATETIME()")
			case EngineOracle:
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			default:
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			}
		} else {
			parts = append(parts, "DEFAULT "+defVal)
		}
	}

	return strings.Join(parts, " "), false
}

// PresetColumnsStandard returns standard starter columns
func PresetColumnsStandard() []ColumnDef {
	return []ColumnDef{
		{Name: "id", Type: TypeAutoID, IsPrimaryKey: true, IsNotNull: true},
		{Name: "nama", Type: TypeVarchar, Length: "255", IsNotNull: true},
		{Name: "is_active", Type: TypeBoolean, IsNotNull: true, DefaultValue: "1"},
		{Name: "created_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
	}
}

// PresetColumnsUser returns user authentication table columns
func PresetColumnsUser() []ColumnDef {
	return []ColumnDef{
		{Name: "id", Type: TypeAutoID, IsPrimaryKey: true, IsNotNull: true},
		{Name: "username", Type: TypeVarchar, Length: "50", IsNotNull: true, IsUnique: true},
		{Name: "email", Type: TypeVarchar, Length: "100", IsNotNull: true, IsUnique: true},
		{Name: "password_hash", Type: TypeVarchar, Length: "255", IsNotNull: true},
		{Name: "role", Type: TypeVarchar, Length: "20", IsNotNull: true, DefaultValue: "'user'"},
		{Name: "is_active", Type: TypeBoolean, IsNotNull: true, DefaultValue: "1"},
		{Name: "created_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
	}
}

// PresetColumnsProduct returns e-commerce product table columns
func PresetColumnsProduct() []ColumnDef {
	return []ColumnDef{
		{Name: "id", Type: TypeAutoID, IsPrimaryKey: true, IsNotNull: true},
		{Name: "sku", Type: TypeVarchar, Length: "50", IsNotNull: true, IsUnique: true},
		{Name: "nama_produk", Type: TypeVarchar, Length: "255", IsNotNull: true},
		{Name: "deskripsi", Type: TypeText},
		{Name: "harga", Type: TypeDecimal, Length: "12,2", IsNotNull: true, DefaultValue: "0.00"},
		{Name: "stok", Type: TypeInteger, IsNotNull: true, DefaultValue: "0"},
		{Name: "is_available", Type: TypeBoolean, IsNotNull: true, DefaultValue: "1"},
		{Name: "created_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
		{Name: "updated_at", Type: TypeDateTime, IsNotNull: true, DefaultValue: "CURRENT_TIMESTAMP"},
	}
}
