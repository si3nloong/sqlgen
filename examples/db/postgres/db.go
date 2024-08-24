package postgresdb

import (
	"context"
	"database/sql"
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

func InsertOne[T sequel.TableColumnValuer, Ptr interface {
	sequel.TableColumnValuer
	sequel.PtrScanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.SingleInserter:
		query, args := v.InsertOneStmt()
		return sqlConn.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	default:
		columns, values := model.Columns(), model.Values()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		cols := strings.Join(columns, ",")
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES (")
		for i := range values {
			if i > 0 {
				stmt.WriteString("," + wrapVar(i+1))
			} else {
				// argument always started from 1
				stmt.WriteString(wrapVar(i + 1))
			}
		}
		stmt.WriteString(") RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		return sqlConn.QueryRowContext(ctx, stmt.String(), values...).Scan(model.Addrs()...)
	}
}

// Insert is a helper function to insert multiple records.
func Insert[T sequel.Inserter, Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model   = data[0]
		columns = model.Columns()
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
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
		stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var i int64
		for rows.Next() {
			if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
				return nil, err
			}
			i++
		}
		return sequel.NewRowsAffectedResult(i), rows.Close()
	default:
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
		stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var i int64
		for rows.Next() {
			if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
				return nil, err
			}
			i++
		}
		return sequel.NewRowsAffectedResult(i), rows.Close()
	}
}

type duplicateKeyer interface {
	DuplicateKeys() []string
}
type upsertOpts struct {
	doNothing     bool
	omitFields    []string
	duplicateKeys []string
}
type UpsertOption func(*upsertOpts)

func WithDoNothing(doNothing bool) UpsertOption {
	return func(opt *upsertOpts) {
		opt.doNothing = doNothing
	}
}

func WithOmitFields(fields ...string) UpsertOption {
	return func(opt *upsertOpts) {
		opt.omitFields = fields
	}
}

func WithDuplicateKeys(keys []string) UpsertOption {
	return func(opt *upsertOpts) {
		opt.duplicateKeys = keys
	}
}

func UpsertOne[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, opts ...UpsertOption) error {
	var opt upsertOpts
	for i := range opts {
		opts[i](&opt)
	}

	var (
		stmt     = strpool.AcquireString()
		columns  = model.Columns()
		noOfCols = len(columns)
		values   = model.Values()
	)
	defer strpool.ReleaseString(stmt)

	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		columns = append(columns[:idx], columns[idx+1:]...)
		noOfCols = len(columns)
		values = append(values[:idx], values[idx+1:]...)
		// Don't include auto increment primary key on INSERT
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	default:
		panic("unreachable")
	}
	for i := 1; i <= noOfCols; i++ {
		if i > 1 {
			stmt.WriteString("," + wrapVar(i))
		} else {
			stmt.WriteString(wrapVar(i))
		}
	}
	if len(opt.duplicateKeys) > 0 {
		stmt.WriteString(") ON CONFLICT(" + strings.Join(opt.duplicateKeys, ",") + ")")
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			stmt.WriteString(") ON CONFLICT(" + strings.Join(keys, ",") + ")")
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(") ON CONFLICT(" + pkName + ")")
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			stmt.WriteString(") ON CONFLICT(" + strings.Join(names, ",") + ")")
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(" DO NOTHING")
	} else {
		stmt.WriteString(" DO UPDATE SET ")
		omitDict := make(map[string]struct{})
		noOfCols = len(columns)
		for i := range opt.omitFields {
			omitDict[opt.omitFields[i]] = struct{}{}
		}
		for i := 0; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i] + ",")
			} else {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			}
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
	return sqlConn.QueryRowContext(ctx, stmt.String(), values...).Scan(model.Addrs()...)
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T interface {
	sequel.Keyer
	sequel.Inserter
}, Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T, opts ...UpsertOption) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.Columns()
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*noOfData)
		opt      upsertOpts
	)

	for i := range opts {
		opts[i](&opt)
	}

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		columns = append(columns[:idx], columns[idx+1:]...)
		// Don't include auto increment primary key on INSERT
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
	default:
		panic("unreachable")
	}
	if len(opt.duplicateKeys) > 0 {
		stmt.WriteString(" ON CONFLICT(" + strings.Join(opt.duplicateKeys, ",") + ")")
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			stmt.WriteString(" ON CONFLICT(" + strings.Join(keys, ",") + ")")
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(" ON CONFLICT(" + pkName + ")")
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			// // Exclude primary key, don't update it
			// for i := len(idxs) - 1; i >= 0; i-- {
			// 	columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
			// }
			stmt.WriteString(" ON CONFLICT(" + strings.Join(names, ",") + ")")
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(" DO NOTHING")
	} else {
		stmt.WriteString(" DO UPDATE SET ")
		omitDict := make(map[string]struct{})
		columns = model.Columns()
		for i := range opt.omitFields {
			omitDict[opt.omitFields[i]] = struct{}{}
		}
		for i := 0; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i] + ",")
			} else {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			}
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
	rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var i int64
	for rows.Next() {
		if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
			return nil, err
		}
		i++
	}
	return sequel.NewRowsAffectedResult(i), rows.Close()
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.KeyFinder:
		query, args := v.FindOneByPKStmt()
		return sqlConn.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	case sequel.PrimaryKeyer:
		columns := TableColumns(model)
		pkName, _, pk := v.PK()
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+DbTable(model)+" WHERE "+pkName+" = $1 LIMIT 1;", pk).Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		columns := TableColumns(model)
		names, _, keys := v.CompositeKey()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("SELECT " + strings.Join(columns, ",") + " FROM " + DbTable(model) + " WHERE ")
		noOfKey := len(names)
		for i := 0; i < noOfKey; i++ {
			if i > 0 {
				stmt.WriteString(" AND " + names[i] + " = " + wrapVar(i+1))
			} else {
				stmt.WriteString(names[i] + " = " + wrapVar(i+1))
			}
		}
		stmt.WriteString(" LIMIT 1;")
		return sqlConn.QueryRowContext(ctx, stmt.String(), keys...).Scan(model.Addrs()...)
	default:
		panic("unreachable")
	}
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyUpdater:
		query, args := v.UpdateOneByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, pkIdx, pk := v.PK()
		columns := model.Columns()
		columns = append(columns[:pkIdx], columns[pkIdx+1:]...)
		values := model.Values()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("UPDATE " + DbTable(model) + " SET ")
		for idx := range columns {
			if idx > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(columns[idx] + " = " + wrapVar(idx+1))
		}
		stmt.WriteString(" WHERE " + pkName + " = " + wrapVar(len(columns)+2) + ";")
		return sqlConn.ExecContext(ctx, stmt.String(), append(values, pk)...)
	default:
		panic("unreachable")
	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyDeleter:
		query, args := v.DeleteOneByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		return sqlConn.ExecContext(ctx, "DELETE FROM "+DbTable(model)+" WHERE "+pkName+" = $1;", pk)
	case sequel.CompositeKeyer:
		names, _, keys := v.CompositeKey()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("DELETE FROM " + DbTable(model) + " WHERE ")
		noOfKey := len(names)
		for i := 0; i < noOfKey; i++ {
			if i == 1 {
				stmt.WriteString(names[i] + " = " + wrapVar(i+1))
			} else {
				stmt.WriteString(" AND " + names[i] + " = " + wrapVar(i+1))
			}
		}
		stmt.WriteByte(';')
		return sqlConn.ExecContext(ctx, stmt.String(), keys...)
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
				blr.WriteString(strings.Join(TableColumns(vj), ","))
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

var (
	_ sequel.Stmt = (*sqlStmt)(nil)
)

func (s *sqlStmt) Var(value any) string {
	s.pos++
	s.args = append(s.args, value)
	return wrapVar(s.pos)
}

func (s *sqlStmt) Vars(values []any) string {
	noOfLen := len(values)
	s.args = append(s.args, values...)
	buf := new(strings.Builder)
	buf.WriteByte('(')
	i := s.pos
	s.pos += noOfLen
	for ; i < s.pos; i++ {
		if i < s.pos-1 {
			buf.WriteString(wrapVar(i+1) + ",")
		} else {
			buf.WriteString(wrapVar(i + 1))
		}
	}
	buf.WriteByte(')')
	return buf.String()
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

func TableColumns[T sequel.Columner](model T) []string {
	if v, ok := any(model).(sequel.SQLColumner); ok {
		return v.SQLColumns()
	}
	return model.Columns()
}

func dbName(model any) string {
	if v, ok := model.(sequel.Databaser); ok {
		return v.DatabaseName() + "."
	}
	return ""
}

func wrapVar(i int) string {
	return `$` + strconv.Itoa(i)
}
