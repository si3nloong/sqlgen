package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

type autoIncrKeyInserter interface {
	sequel.AutoIncrKeyer
	sequel.SingleInserter
}

func InsertOne[T sequel.TableColumnValuer[T], Ptr interface {
	sequel.TableColumnValuer[T]
	sequel.PtrScanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) (sql.Result, error) {
	switch v := any(model).(type) {
	case autoIncrKeyInserter:
		_, idx, _ := v.PK()
		query, args := v.InsertOneStmt()
		result, err := sqlConn.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		i64, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		switch v := model.Addrs()[idx].(type) {
		case *int64:
			*v = i64
		case sql.Scanner:
			if err := v.Scan(i64); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(`sqlgen: invalid auto increment data type`)
		}
		return result, nil
	case sequel.SingleInserter:
		query, args := v.InsertOneStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	default:
		columns, values := model.ColumnNames(), model.Values()
		return sqlConn.ExecContext(ctx, "INSERT INTO "+DbTable(model)+" ("+strings.Join(columns, ",")+") VALUES ("+strings.Repeat(",?", len(columns))[1:]+");", values...)
	}
}

// Insert is a helper function to insert multiple records.
func Insert[T sequel.TableColumnValuer[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model   = data[0]
		columns = model.ColumnNames()
		stmt    = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)

	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		_, idx, _ := v.PK()
		columns = append(columns[:idx], columns[idx+1:]...)
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		placeholder := "(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + placeholder)
			} else {
				stmt.WriteString(placeholder)
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
		stmt.WriteByte(';')
		return sqlConn.ExecContext(ctx, stmt.String(), args...)
	default:
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + cols + ") VALUES ")
		placeholder := "(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + placeholder)
			} else {
				stmt.WriteString(placeholder)
			}
			args = append(args, data[i].Values()...)
		}
		stmt.WriteByte(';')
		return sqlConn.ExecContext(ctx, stmt.String(), args...)
	}
}

func UpsertOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, override bool, omittedFields ...string) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.PrimaryKeyer:
		pkName, idx, _ := v.PK()
		columns := model.ColumnNames()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			omitDict := map[string]struct{}{pkName: {}}
			for i := range omittedFields {
				omitDict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			for i := range columns {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if i < noOfCols-1 {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
				} else {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
				}
			}
			clear(omitDict)
		}
		if _, ok := any(model).(sequel.AutoIncrKeyer); ok {
			result, err := sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
			if err != nil {
				return nil, err
			}
			i64, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			switch v := model.Addrs()[idx].(type) {
			case *int64:
				*v = i64
			case sql.Scanner:
				if err := v.Scan(i64); err != nil {
					return nil, err
				}
			default:
				return nil, errors.New(`sqlgen: invalid auto increment data type`)
			}
			return result, nil
		}
		return sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
	case sequel.CompositeKeyer:
		names, idxs, _ := v.CompositeKey()
		columns := model.ColumnNames()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			dict := make(map[string]struct{})
			for i := range append(names, omittedFields...) {
				dict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			// Exclude composite key, don't update it
			for i := len(idxs) - 1; i >= 0; i-- {
				columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
			}
			for i := range columns {
				if _, ok := dict[columns[i]]; ok {
					continue
				}
				if i < noOfCols-1 {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
				} else {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
				}
			}
			clear(dict)
		}
		return sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
	default:
		panic("unreachable")
	}
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T sequel.KeyValuer[T], Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T, override bool, omittedFields ...string) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.ColumnNames()
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*noOfData)
	)

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		// Don't include auto increment primary key on INSERT
		columns = append(columns[:idx], columns[idx+1:]...)
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
			args = append(args, data[i].Values()...)
		}
		_, idxs, _ := v.CompositeKey()
		// Exclude primary key, don't update it
		for i := len(idxs) - 1; i >= 0; i-- {
			columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
		}
	}
	if override {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		// Don't update primary key when we do upsert
		omitDict := map[string]struct{}{}
		for i := range omittedFields {
			omitDict[omittedFields[i]] = struct{}{}
		}
		noOfCols = len(columns)
		for i := range columns {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + "=VALUES(" + columns[i] + "),")
			} else {
				stmt.WriteString(columns[i] + "=VALUES(" + columns[i] + ")")
			}
		}
		clear(omitDict)
	}
	stmt.WriteByte(';')
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
}

type primaryKeyFinder interface {
	sequel.PrimaryKeyer
	sequel.KeyFinder
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case primaryKeyFinder:
		query, args := v.FindByPKStmt()
		return sqlConn.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	case sequel.PrimaryKeyer:
		columns := Columns(model)
		pkName, _, pk := v.PK()
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+DbTable(model)+" WHERE "+pkName+" = ? LIMIT 1;", pk).Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		columns := Columns(model)
		names, _, keys := v.CompositeKey()
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+DbTable(model)+" WHERE "+strings.Join(names, " = ? AND ")+" = ? LIMIT 1;", keys...).Scan(model.Addrs()...)
	default:
		panic("unreachable")
	}
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyUpdater:
		query, args := v.UpdateByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, pkIdx, pk := v.PK()
		columns := model.ColumnNames()
		columns = append(columns[:pkIdx], columns[pkIdx+1:]...)
		values := model.Values()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
		return sqlConn.ExecContext(ctx, "UPDATE "+DbTable(model)+" SET "+strings.Join(columns, " = ?,")+" = ? WHERE "+pkName+" = ?;", append(values, pk)...)
	default:
		panic("unreachable")
	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyDeleter:
		query, args := v.DeleteByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		return sqlConn.ExecContext(ctx, "DELETE FROM "+DbTable(model)+" WHERE "+pkName+" = ?;", pk)
	case sequel.CompositeKeyer:
		names, _, keys := v.CompositeKey()
		return sqlConn.ExecContext(ctx, "DELETE FROM "+DbTable(model)+" WHERE "+strings.Join(names, " = ? AND ")+" = ?;", keys...)
	default:
		panic("unreachable")
	}
}

type SelectStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	GroupBy   []string
	Offset    uint64
	Limit     uint16
}

type SQLStatement struct {
	Query     string
	Arguments []any
}

func QueryStmt[T any, Ptr sequel.PtrScanner[T], Stmt interface {
	SelectStmt | SQLStatement
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) ([]T, error) {
	var (
		rows *sql.Rows
		err  error
	)

	switch vi := any(stmt).(type) {
	case SelectStmt:
		var (
			blr = AcquireStmt()
			v   T
		)
		defer ReleaseStmt(blr)
		blr.WriteString("SELECT ")
		if len(vi.Select) > 0 {
			blr.WriteString(strings.Join(vi.Select, ","))
		} else {
			switch vj := any(v).(type) {
			case sequel.Columner:
				blr.WriteString(strings.Join(Columns(vj), ","))
			default:
				blr.WriteByte('*')
			}
		}
		if vi.FromTable != "" {
			blr.WriteString(" FROM " + dbName(v) + vi.FromTable)
		} else {
			switch vj := any(v).(type) {
			case sequel.Tabler:
				blr.WriteString(" FROM " + DbTable(vj))
			default:
				return nil, fmt.Errorf("missing table name for model %T", v)
			}
		}
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.GroupBy) > 0 {
			blr.WriteString(" GROUP BY ")
			for i := range vi.GroupBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				blr.WriteString(vi.GroupBy[i])
			}
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		if vi.Offset > 0 {
			blr.WriteString(" OFFSET " + strconv.FormatUint(vi.Offset, 10))
		}
		blr.WriteByte(';')
		rows, err = sqlConn.QueryContext(ctx, blr.String(), blr.Args()...)
		ReleaseStmt(blr)

	case SQLStatement:
		rows, err = sqlConn.QueryContext(ctx, vi.Query, vi.Arguments...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []T
	for rows.Next() {
		var v T
		if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateStmt struct {
	Table   string
	Set     []sequel.SetClause
	Where   sequel.WhereClause
	OrderBy []sequel.OrderByClause
	Limit   uint16
}

type DeleteStmt struct {
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

func ExecStmt[T any, Stmt interface {
	UpdateStmt | DeleteStmt
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) (sql.Result, error) {
	blr := AcquireStmt()
	defer ReleaseStmt(blr)

	var v T
	switch vi := any(stmt).(type) {
	case UpdateStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("UPDATE " + DbTable(vt))
		} else {
			blr.WriteString("UPDATE " + dbName(v) + vi.Table)
		}
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.Set) > 0 {
			blr.WriteString(" SET ")
			for i := range vi.Set {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.Set[i](blr)
			}
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}

	case DeleteStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("DELETE FROM " + DbTable(vt))
		} else {
			blr.WriteString("DELETE FROM " + dbName(v) + vi.FromTable)
		}
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
	}
	blr.WriteByte(';')
	return sqlConn.ExecContext(ctx, blr.String(), blr.Args()...)
}

var (
	pool = sync.Pool{
		New: func() any {
			return new(sqlStmt)
		},
	}
)

func AcquireStmt() sequel.Stmt {
	return pool.Get().(*sqlStmt)
}

func ReleaseStmt(stmt sequel.Stmt) {
	if stmt != nil {
		stmt.Reset()
		pool.Put(stmt)
	}
}

type sqlStmt struct {
	strings.Builder
	pos  int
	args []any
}

func (s *sqlStmt) Var(query string, value any) {
	s.pos++
	s.WriteString(query + "?")
	s.args = append(s.args, value)
}

func (s *sqlStmt) Vars(query string, values []any) {
	s.WriteString(query)
	noOfLen := len(values)
	s.WriteString("(" + strings.Repeat(",?", noOfLen)[1:] + ")")
	s.args = append(s.args, values...)
}

func (s sqlStmt) Args() []any {
	return s.args
}

func (s *sqlStmt) Reset() {
	s.args = nil
	s.pos = 0
	s.Builder.Reset()
}

func DbTable[T sequel.Tabler](model T) string {
	if v, ok := any(model).(sequel.Databaser); ok {
		return v.DatabaseName() + "." + model.TableName()
	}
	return model.TableName()
}

func dbName(model any) string {
	if v, ok := model.(sequel.Databaser); ok {
		return v.DatabaseName() + "."
	}
	return ""
}

func Columns[T sequel.Columner](model T) []string {
	if v, ok := any(model).(sequel.SQLColumner); ok {
		return v.SQLColumns()
	}
	return model.ColumnNames()
}
