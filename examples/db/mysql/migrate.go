package mysqldb

import (
	"context"
	"strings"
	"sync"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

var (
	mu sync.RWMutex
)

type Column struct {
	Name       string
	Type       string
	DataType   string
	Size       int
	Pos        int
	IsNullable bool
	Default    any
	Comment    string
}

type Index struct {
	Name     string
	Type     string
	Nullable bool
}

func TableExists(ctx context.Context, sqlConn sequel.DB, dbName, tableName string) (bool, error) {
	var count int64
	if err := sqlConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`, dbName, tableName).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func UnsafeMigrate[T interface {
	sequel.Tabler
	sequel.Migrator
}](ctx context.Context, sqlConn sequel.DB, dbName string) error {
	mu.Lock()
	defer mu.Unlock()

	var v T
	tableName := v.TableName()
	exists, err := TableExists(ctx, sqlConn, dbName, tableName)
	if err != nil {
		return err
	}
	def := v.Schemas()
	stmt := strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	// If the table exists, we will use alter table
	if exists {
		// Alter table need to check primary key, foreign key and indexes
		stmt.WriteString("ALTER TABLE " + tableName + " (")
		colDict := make(map[string]Column)
		if err := tableColumns(ctx, sqlConn, dbName, tableName, func(c Column, _ int) {
			colDict[c.Name] = c
		}); err != nil {
			return err
		}
		for i, col := range def.Columns {
			if _, ok := colDict[col.Name]; !ok {
				stmt.WriteString(",ADD COLUMN " + col.Definition)
			} else {
				stmt.WriteString(",MODIFY COLUMN " + col.Definition)
			}
			if i > 0 {
				stmt.WriteString(" FIRST")
			} else {
				stmt.WriteString(" AFTER " + def.Columns[i-1].Name)
			}
		}
		clear(colDict)
		idxDict := make(map[string]Index)
		if err := tableIndexes(ctx, sqlConn, dbName, tableName, func(idx Index, _ int) {
			idxDict[idx.Name] = idx
		}); err != nil {
			return err
		}
		for _, idx := range def.Indexes {
			if _, ok := idxDict[idx.Name]; !ok {
				stmt.WriteString(",DROP INDEX " + idx.Name)
			} else {
				stmt.WriteString(",ADD " + idx.Definition)
			}
		}
		clear(idxDict)
		stmt.WriteString(") ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	} else {
		stmt.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (")
		for i, col := range def.Columns {
			if i > 0 {
				stmt.WriteString("," + col.Definition + " AFTER " + def.Columns[i-1].Name)
			} else {
				stmt.WriteString(col.Definition + " FIRST")
			}
		}
		switch vi := any(v).(type) {
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(",PRIMARY KEY(" + pkName + ")")
		case sequel.CompositeKeyer:
			keys, _, _ := vi.CompositeKey()
			stmt.WriteString(",PRIMARY KEY(" + strings.Join(keys, ",") + ")")
		}
		for _, idx := range def.Indexes {
			stmt.WriteString(",ADD " + idx.Definition)
		}
		stmt.WriteString(") ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	}
	if _, err := sqlConn.ExecContext(ctx, stmt.String()); err != nil {
		return err
	}
	return nil
}

func tableColumns(ctx context.Context, sqlConn sequel.DB, dbName, tableName string, reduceFunc func(Column, int)) error {
	rows, err := sqlConn.QueryContext(ctx, `SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION;`, dbName, tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var column Column
		if err := rows.Scan(&column.Pos, &column.Name, &column.Type, &column.Default, &column.IsNullable, &column.DataType, &column.Comment); err != nil {
			return err
		}
		reduceFunc(column, i)
		i++
	}
	return rows.Close()
}

func tableIndexes(ctx context.Context, sqlConn sequel.DB, dbName, tableName string, reduceFunc func(Index, int)) error {
	rows, err := sqlConn.QueryContext(ctx, `SELECT DISTINCT INDEX_NAME, INDEX_TYPE, NON_UNIQUE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`, dbName, tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var index Index
		if err := rows.Scan(&index.Name, &index.Type, &index.Nullable); err != nil {
			return err
		}
		reduceFunc(index, i)
		i++
	}
	return rows.Close()
}
