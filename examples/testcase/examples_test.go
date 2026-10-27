package testcase_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ory/dockertest/v3"
	mysqldb "github.com/si3nloong/sqlgen/examples/db/mysql"
	sqlitedb "github.com/si3nloong/sqlgen/examples/db/sqlite"
	"github.com/si3nloong/sqlgen/examples/testcase/core"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/pk"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/pointer"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/stretchr/testify/require"
)

// The checked-in Car does not emit SQLInsertColumns or SQLInsertValues, which Upsert's
// sequel.Inserter constraint requires. Keep this compatibility adapter local
// to the tests; all SQL, values and scanners still come from the generated Car.
type car struct{ pk.Car }

func (car) SQLInsertColumns() string { return "id,no,color,year,manuc_date" }
func (v car) SQLInsertValues() []any { return v.Values() }

// SQLGEN_MYSQL_DSN may point to a dedicated test database. Otherwise a disposable
// MySQL 8 container is used. Each test uses connection-local temporary tables,
// so no existing application tables are dropped or modified.
func mysqlConn(t *testing.T) *sql.Conn {
	t.Helper()
	dsn := os.Getenv("SQLGEN_MYSQL_DSN")
	if dsn == "" {
		pool, err := dockertest.NewPool(dockerEndpoint())
		require.NoError(t, err)
		if err := pool.Client.Ping(); err != nil {
			if os.Getenv("SQLGEN_REQUIRE_MYSQL") == "1" {
				t.Fatalf("MySQL integration tests require Docker or SQLGEN_MYSQL_DSN: %v", err)
			}
			t.Skipf("MySQL integration tests require Docker or SQLGEN_MYSQL_DSN: %v", err)
		}
		resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_DATABASE=sqlgen", "MYSQL_ROOT_PASSWORD=secret"})
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, pool.Purge(resource)) })
		dsn = fmt.Sprintf("root:secret@tcp(%s)/sqlgen?parseTime=true&loc=UTC", resource.GetHostPort("3306/tcp"))
		pool.MaxWait = 90 * time.Second
		db, err := sql.Open("mysql", dsn)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })
		require.NoError(t, pool.Retry(db.Ping))
		conn, err := db.Conn(t.Context())
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, conn.Close()) })
		return conn
	}
	config, err := mysql.ParseDSN(dsn)
	require.NoError(t, err)
	config.ParseTime = true
	config.Loc = time.UTC
	db, err := sql.Open("mysql", config.FormatDSN())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })
	require.NoError(t, conn.PingContext(t.Context()))
	return conn
}

func affected(t *testing.T, result sql.Result, err error, want int64) {
	t.Helper()
	require.NoError(t, err)
	require.NotNil(t, result)
	n, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, want, n)
}

func sampleCar(id int) car {
	year := 2020 + id
	return car{pk.Car{ID: pk.PK(id), No: fmt.Sprintf("CAR-%02d", id), Color: pk.White, Year: &year, ManucDate: time.Date(2024, 2, 3, 4, 5, 6, 123456000, time.UTC)}}
}

func TestMySQLGeneratedAPI(t *testing.T) {
	conn := mysqlConn(t)
	for _, bound := range []bool{false, true} {
		name := "functions"
		if bound {
			name = "NewDB methods"
		}
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			_, err := conn.ExecContext(ctx, "CREATE TEMPORARY TABLE car (id BIGINT PRIMARY KEY, no VARCHAR(100) NOT NULL, color INT NOT NULL, year INT NULL, manuc_date DATETIME(6) NOT NULL)")
			require.NoError(t, err)
			defer func() { _, err := conn.ExecContext(ctx, "DROP TEMPORARY TABLE car"); require.NoError(t, err) }()
			insertOne := func(v *car) (sql.Result, error) { return mysqldb.InsertOne[car](ctx, conn, v) }
			insert := func(v []car) (sql.Result, error) { return mysqldb.Insert[car](ctx, conn, v) }
			upsertOne := func(v *car, override bool, omit ...string) (sql.Result, error) {
				return mysqldb.UpsertOne[car](ctx, conn, v, override, omit...)
			}
			upsert := func(v []car, override bool, omit ...string) (sql.Result, error) {
				return mysqldb.Upsert[car](ctx, conn, v, override, omit...)
			}
			find := func(v *car) error { return mysqldb.FindByPK[car](ctx, conn, v) }
			update := func(v car) (sql.Result, error) { return mysqldb.UpdateByPK[car](ctx, conn, v) }
			deletePK := func(v car) (sql.Result, error) { return mysqldb.DeleteByPK[car](ctx, conn, v) }
			query := func(fn func(car) mysqldb.SelectStmt) ([]car, error) { return mysqldb.QueryStmt[car](ctx, conn, fn) }
			queryOne := func(fn func(car) mysqldb.SelectOneStmt) (*car, error) {
				return mysqldb.QueryOneStmt[car](ctx, conn, fn)
			}
			execUpdate := func(fn func(car) mysqldb.UpdateStmt) (sql.Result, error) { return mysqldb.ExecStmt[car](ctx, conn, fn) }
			execDelete := func(fn func(car) mysqldb.DeleteStmt) (sql.Result, error) { return mysqldb.ExecStmt[car](ctx, conn, fn) }
			paginate := mysqldb.Paginate[car]
			if bound {
				db := mysqldb.NewDB(conn)
				insertOne = func(v *car) (sql.Result, error) { return db.InsertOne[car](ctx, v) }
				insert = func(v []car) (sql.Result, error) { return db.Insert[car](ctx, v) }
				upsertOne = func(v *car, override bool, omit ...string) (sql.Result, error) {
					return db.UpsertOne[car](ctx, v, override, omit...)
				}
				upsert = func(v []car, override bool, omit ...string) (sql.Result, error) {
					return db.Upsert[car](ctx, v, override, omit...)
				}
				find = func(v *car) error { return db.FindByPK[car](ctx, v) }
				update = func(v car) (sql.Result, error) { return db.UpdateByPK[car](ctx, v) }
				deletePK = func(v car) (sql.Result, error) { return db.DeleteByPK[car](ctx, v) }
				query = func(fn func(car) mysqldb.SelectStmt) ([]car, error) { return db.QueryStmt[car](ctx, fn) }
				queryOne = func(fn func(car) mysqldb.SelectOneStmt) (*car, error) { return db.QueryOneStmt[car](ctx, fn) }
				execUpdate = func(fn func(car) mysqldb.UpdateStmt) (sql.Result, error) { return db.ExecStmt[car](ctx, fn) }
				execDelete = func(fn func(car) mysqldb.DeleteStmt) (sql.Result, error) { return db.ExecStmt[car](ctx, fn) }
				paginate = db.Paginate[car]
			}

			check := func(want car) {
				t.Helper()
				got := car{Car: pk.Car{ID: want.ID}}
				require.NoError(t, find(&got))
				require.Equal(t, want, got)
			}
			first := sampleCar(1)
			result, err := insertOne(&first)
			affected(t, result, err, 1)
			check(first)
			others := []car{sampleCar(2), sampleCar(3)}
			others[0].Year = nil
			result, err = insert(others)
			affected(t, result, err, 2)
			for _, v := range others {
				check(v)
			}
			result, err = insert(nil)
			affected(t, result, err, 0)
			result, err = upsert(nil, true)
			affected(t, result, err, 0)

			// A duplicate insert must fail without changing the stored row.
			duplicate := sampleCar(1)
			duplicate.No = "duplicate"
			_, err = insertOne(&duplicate)
			require.Error(t, err)
			check(first)

			inserted := sampleCar(4)
			result, err = upsertOne(&inserted, false)
			affected(t, result, err, 1)
			check(inserted)
			ignored := inserted
			ignored.No = "ignored"
			result, err = upsertOne(&ignored, false)
			affected(t, result, err, 0)
			check(inserted)
			changed := inserted
			changed.No = "updated"
			changed.Color = pk.Red
			result, err = upsertOne(&changed, true, "no")
			affected(t, result, err, 2)
			changed.No = inserted.No
			check(changed)
			changed.No = "overwritten"
			result, err = upsertOne(&changed, true)
			affected(t, result, err, 2)
			check(changed)

			batch := []car{sampleCar(5), sampleCar(6)}
			result, err = upsert(batch, false)
			affected(t, result, err, 2)
			for _, v := range batch {
				check(v)
			}
			batch[0].No = "ignored batch"
			result, err = upsert(batch, false)
			affected(t, result, err, 0)
			check(sampleCar(5))
			batch[0].Color = pk.Red
			batch[1].Color = pk.Red
			result, err = upsert(batch, true, "no")
			affected(t, result, err, 4)
			batch[0].No = sampleCar(5).No
			for _, v := range batch {
				check(v)
			}
			batch[0].No = "batch updated"
			batch[1].No = "batch updated too"
			result, err = upsert(batch, true)
			affected(t, result, err, 4)
			for _, v := range batch {
				check(v)
			}

			first.No = "update by key"
			first.Year = nil
			result, err = update(first)
			affected(t, result, err, 1)
			check(first)
			check(others[0]) // Updating one key must not change its neighbours.

			got, err := queryOne(func(v car) mysqldb.SelectOneStmt {
				return mysqldb.SelectOneStmt{Where: mysqldb.Equal(v.ColumnID(), first.ID)}
			})
			require.NoError(t, err)
			require.Equal(t, first, *got)
			rows, err := query(func(v car) mysqldb.SelectStmt {
				return mysqldb.SelectStmt{Where: mysqldb.In(v.ColumnID(), []pk.PK{1, 2, 3}), OrderBy: []sequel.OrderByClause{mysqldb.Asc(v.ColumnID())}, Limit: 2, Offset: 1}
			})
			require.NoError(t, err)
			require.Equal(t, others, rows)
			rows, err = query(func(v car) mysqldb.SelectStmt {
				return mysqldb.SelectStmt{Where: mysqldb.Equal(v.ColumnID(), pk.PK(999))}
			})
			require.NoError(t, err)
			require.Empty(t, rows)
			_, err = queryOne(func(v car) mysqldb.SelectOneStmt {
				return mysqldb.SelectOneStmt{Where: mysqldb.Equal(v.ColumnID(), pk.PK(999))}
			})
			require.ErrorIs(t, err, sql.ErrNoRows)

			result, err = execUpdate(func(v car) mysqldb.UpdateStmt {
				return mysqldb.UpdateStmt{Set: []sequel.SetClause{mysqldb.Set(v.ColumnNo(), "statement updated")}, Where: mysqldb.Equal(v.ColumnID(), first.ID)}
			})
			affected(t, result, err, 1)
			first.No = "statement updated"
			check(first)
			check(others[0])

			for _, mode := range []mysqldb.LockMode{mysqldb.LockForUpdate, mysqldb.LockForShare} {
				locked := car{Car: pk.Car{ID: first.ID}}
				if bound {
					require.NoError(t, mysqldb.NewDB(conn).LockByPK(ctx, &locked, mode))
				} else {
					tx, err := conn.BeginTx(ctx, nil)
					require.NoError(t, err)
					err = mysqldb.LockByPK(ctx, tx, &locked, mode)
					rollbackErr := tx.Rollback()
					require.NoError(t, err)
					require.NoError(t, rollbackErr)
				}
				require.Equal(t, first, locked)
			}

			want := []car{first, others[0], others[1], changed, batch[0], batch[1]}
			for _, reverse := range []bool{false, true} {
				pager := paginate(mysqldb.PaginateStmt{Limit: 2})
				seq := pager.Next(ctx, conn)
				expected := append([]car(nil), want...)
				if reverse {
					seq = pager.Prev(ctx, conn)
					for i, j := 0, len(expected)-1; i < j; i, j = i+1, j-1 {
						expected[i], expected[j] = expected[j], expected[i]
					}
				}
				var all []car
				pages := 0
				for page, err := range seq {
					require.NoError(t, err)
					require.LessOrEqual(t, len(page), 2)
					pages++
					require.LessOrEqual(t, pages, 4, "pagination must terminate")
					all = append(all, page...)
				}
				require.Equal(t, expected, all)
			}
			// Explicit cursors are inclusive, matching the pager's look-ahead cursor.
			for page, err := range paginate(mysqldb.PaginateStmt{Limit: 2}).Next(ctx, conn, others[0]) {
				require.NoError(t, err)
				require.Equal(t, others, page)
				break // Also exercise early iterator termination and connection cleanup.
			}
			check(first)

			result, err = deletePK(first)
			affected(t, result, err, 1)
			require.ErrorIs(t, find(&first), sql.ErrNoRows)
			result, err = deletePK(first)
			affected(t, result, err, 0)
			result, err = update(first)
			affected(t, result, err, 0)
			result, err = execDelete(func(v car) mysqldb.DeleteStmt {
				return mysqldb.DeleteStmt{Where: mysqldb.In(v.ColumnID(), []pk.PK{2, 3})}
			})
			affected(t, result, err, 2)
			for _, v := range others {
				require.ErrorIs(t, find(&v), sql.ErrNoRows)
			}
			check(changed)
			for _, v := range batch {
				check(v)
			}
		})
	}
}

func ptrOf[T any](v T) *T { return &v }

func sampleUser() core.User {
	v := core.User{No: 42, JoinedTime: time.Date(2024, 2, 3, 4, 5, 6, 123456000, time.UTC), Address: core.Address{Line1: "12 Main Street", Line2: "Unit 3", CountryCode: "MY"}, Kind: reflect.String, Type: core.HouseUnitTypeB, Chan: reflect.BothDir, PostalCode: ptrOf("12345"), Nicknames: [2]string{"first", "second"}, Slice: []float64{0.125, -3.5, 100}, Map: map[string]float64{"a": 1.25, "b": -2.5}}
	v.Name = "Ada"
	v.T = v.JoinedTime
	v.ExtraInfo.Flag = true
	v.Nested = &struct{ Deep struct{ Bool bool } }{}
	v.Nested.Deep.Bool = true
	return v
}

// Run the serialization round trips without Docker too, against a real SQLite
// database and its generated helpers. MySQL uses the same model fixtures.
func TestGeneratedRoundTrips(t *testing.T) {
	for _, dialect := range []string{"sqlite", "mysql"} {
		t.Run(dialect, func(t *testing.T) {
			var conn *sql.Conn
			if dialect == "mysql" {
				conn = mysqlConn(t)
			} else {
				db, err := sql.Open("sqlite3", ":memory:")
				require.NoError(t, err)
				t.Cleanup(func() { require.NoError(t, db.Close()) })
				conn, err = db.Conn(t.Context())
				require.NoError(t, err)
				t.Cleanup(func() { require.NoError(t, conn.Close()) })
			}
			pkType := "INTEGER PRIMARY KEY AUTOINCREMENT"
			if dialect == "mysql" {
				pkType = "BIGINT PRIMARY KEY AUTO_INCREMENT"
			}
			// Migrate before invoking any generated API. DATETIME(6) preserves the
			// fixture's microseconds, so comparisons do not discard time precision.
			err := migrate(t, conn, dialect, "CREATE TEMPORARY TABLE user (id "+pkType+", no BIGINT, joined_time DATETIME(6), address JSON, kind INT, type INT, chan INT, postal_code TEXT, extra_info JSON, nicknames JSON, slice JSON, map JSON, nested JSON, t DATETIME(6), name TEXT)")
			require.NoError(t, err)
			insertOne := func(v *core.User) (sql.Result, error) {
				return sqlitedb.InsertOne(t.Context(), conn, v)
			}
			insert := func(v []core.User) (sql.Result, error) {
				return sqlitedb.Insert(t.Context(), conn, v)
			}
			find := func(v *core.User) error {
				return sqlitedb.FindByPK(t.Context(), conn, v)
			}
			if dialect == "mysql" {
				insertOne = func(v *core.User) (sql.Result, error) { return mysqldb.InsertOne(t.Context(), conn, v) }
				insert = func(v []core.User) (sql.Result, error) { return mysqldb.Insert(t.Context(), conn, v) }
				find = func(v *core.User) error { return mysqldb.FindByPK(t.Context(), conn, v) }
			}
			want := sampleUser()
			result, err := insertOne(&want)
			affected(t, result, err, 1)
			id, err := result.LastInsertId()
			require.NoError(t, err)
			require.Positive(t, id)
			require.Equal(t, id, want.ID, "InsertOne populates the generated primary key")
			got := core.User{ID: id}
			require.NoError(t, find(&got))
			require.Equal(t, want, got)

			batch := []core.User{sampleUser(), sampleUser()}
			batch[0].Name = "batch one"
			batch[1].Name = "batch two"
			batch[1].PostalCode = nil
			batch[1].Nested = nil
			batch[1].Slice = []float64{}
			batch[1].Map = map[string]float64{}
			result, err = insert(batch)
			affected(t, result, err, 2)
			// Discover generated IDs by a unique fixture value instead of assuming
			// auto_increment_increment=1 or identical LastInsertId semantics.
			for _, v := range batch {
				require.NoError(t, conn.QueryRowContext(t.Context(), "SELECT id FROM user WHERE name = ?", v.Name).Scan(&v.ID))
				got := core.User{ID: v.ID}
				require.NoError(t, find(&got))
				require.Equal(t, v, got)
			}
			testPointerRoundTrip(t, conn, dialect, pkType)
		})
	}
}

func testPointerRoundTrip(t *testing.T, conn *sql.Conn, dialect, pkType string) {
	t.Helper()
	err := migrate(t, conn, dialect, "CREATE TEMPORARY TABLE ptr (id "+pkType+", str TEXT, bytes BLOB, bool BOOLEAN, `int` BIGINT, int_8 SMALLINT, int_16 SMALLINT, int_32 INT, int_64 BIGINT, uint BIGINT, uint_8 SMALLINT, uint_16 INT, uint_32 BIGINT, uint_64 BIGINT, f_32 FLOAT, f_64 DOUBLE, time DATETIME(6), nested JSON, embedded_time DATETIME(6), any_time DATETIME(6), ptr_str TEXT)")
	require.NoError(t, err)
	values := []pointer.Ptr{{}, {Str: ptrOf("hello"), Bytes: ptrOf([]byte{0, 1, 127, 255}), Bool: ptrOf(true), Int: ptrOf(-123), Int8: ptrOf(int8(-88)), Int16: ptrOf(int16(-1234)), Int32: ptrOf(int32(-123456)), Int64: ptrOf(int64(-12345678901)), Uint: ptrOf(uint(123)), Uint8: ptrOf(uint8(88)), Uint16: ptrOf(uint16(1234)), Uint32: ptrOf(uint32(123456)), Uint64: ptrOf(uint64(12345678901)), F32: ptrOf(float32(1.25)), F64: ptrOf(-123.125), Time: ptrOf(time.Date(2024, 2, 3, 4, 5, 6, 123456000, time.UTC))}}
	for i := range values {
		want := &values[i]
		// Addrs initializes the unexported embedded structs used by the fixture.
		// Restore optional fields afterwards to cover actual SQL NULL values.
		want.Addrs()
		want.AnyTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		if i == 0 {
			want.Str = nil
			want.Bytes = nil
			want.Bool = nil
			want.Int = nil
			want.Int8 = nil
			want.Int16 = nil
			want.Int32 = nil
			want.Int64 = nil
			want.Uint = nil
			want.Uint8 = nil
			want.Uint16 = nil
			want.Uint32 = nil
			want.Uint64 = nil
			want.F32 = nil
			want.F64 = nil
			want.Time = nil
		}
		want.Nested = nil
		want.EmbeddedTime = nil
		want.PtrStr = nil
		var result sql.Result
		if dialect == "mysql" {
			result, err = mysqldb.InsertOne(t.Context(), conn, want)
		} else {
			result, err = sqlitedb.InsertOne(t.Context(), conn, want)
		}
		affected(t, result, err, 1)
		require.Positive(t, want.ID)
		got := pointer.Ptr{ID: want.ID}
		if dialect == "mysql" {
			err = mysqldb.FindByPK(t.Context(), conn, &got)
		} else {
			err = sqlitedb.FindByPK(t.Context(), conn, &got)
		}
		require.NoError(t, err)
		require.Equal(t, *want, got)
	}
	var result sql.Result
	if dialect == "mysql" {
		result, err = mysqldb.Insert(t.Context(), conn, values)
	} else {
		result, err = sqlitedb.Insert(t.Context(), conn, values)
	}
	affected(t, result, err, int64(len(values)))
	rows, err := conn.QueryContext(t.Context(), "SELECT id FROM ptr ORDER BY id")
	require.NoError(t, err)
	var ids []int64
	for rows.Next() {
		var id int64
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.Len(t, ids, 2*len(values))
	for i, want := range values {
		want.ID = ids[len(values)+i]
		got := pointer.Ptr{ID: want.ID}
		if dialect == "mysql" {
			err = mysqldb.FindByPK(t.Context(), conn, &got)
		} else {
			err = sqlitedb.FindByPK(t.Context(), conn, &got)
		}
		require.NoError(t, err)
		require.Equal(t, want, got)
	}

}

func migrate(t *testing.T, conn *sql.Conn, dialect, ddl string) error {
	t.Helper()
	if dialect == "sqlite" {
		ddl = strings.ReplaceAll(ddl, "DATETIME(6)", "DATETIME")
	}
	_, err := conn.ExecContext(t.Context(), ddl)
	return err
}

// Docker Desktop on macOS can expose only a per-user socket, without the
// /var/run/docker.sock symlink. Always honor an explicitly configured host.
func dockerEndpoint() string {
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		return host
	}
	if _, err := os.Stat("/var/run/docker.sock"); err == nil {
		return ""
	}
	if home, err := os.UserHomeDir(); err == nil {
		socket := filepath.Join(home, ".docker", "run", "docker.sock")
		if info, err := os.Stat(socket); err == nil && info.Mode()&os.ModeSocket != 0 {
			return "unix://" + socket
		}
	}
	return ""
}
