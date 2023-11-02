package examples_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"

	"github.com/jaswdr/faker"
	autopk "github.com/si3nloong/sqlgen/examples/testcase/struct-field/pk/auto-incr"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/pointer"
	"github.com/si3nloong/sqlgen/sequel/db"
)

var (
	sqliteDB *sql.DB
)

func openSqlConn(driver string) (*sql.DB, error) {
	switch driver {
	case "mysql":
		return sql.Open("mysql", "root:abcd1234@/sqlbench")
	case "sqlite":
		os.Remove("./sqlite.db")
		return sql.Open("sqlite3", "./sqlite.db")
	default:
		return nil, errors.New("unsupported sql driver")
	}
}

func mustValue[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	// openSqlConn("mysql")
	sqliteDB = mustValue(openSqlConn("mysql"))
	defer sqliteDB.Close()

	// m1 := autopk.Model{}
	// sqlutil.FindOne(nil, nil, &m1)

	if err := db.Migrate[pointer.Ptr](context.TODO(), sqliteDB); err != nil {
		panic(err)
	}
	// mustNot(sqliteDB.Exec("DROP TABLE `model`;"))
	// mustNot(sqliteDB.Exec(createTableModel))

	m.Run()
}

func newPKModel() autopk.Model {
	fake := faker.New()
	return autopk.Model{
		F:    true,
		Name: autopk.LongText(fake.Person().Name()),
		N:    fake.Int64Between(0, 100),
	}
}

func TestInsertInto(t *testing.T) {
	ctx := context.TODO()

	t.Run("with all nil values", func(t *testing.T) {
		inputs := []pointer.Ptr{{}, {}}
		result, err := db.InsertInto(ctx, sqliteDB, inputs)
		lastID := mustValue(result.LastInsertId())
		require.NoError(t, err)
		require.NotEmpty(t, lastID)
		require.Equal(t, int64(2), mustValue(result.RowsAffected()))
	})

	t.Run("with pointer values", func(t *testing.T) {
		str := "hello world"
		flag := true
		dt := time.Now().UTC()
		u8 := uint8(100)
		u16 := uint16(1203)
		u32 := uint32(5784182)
		u64 := uint64(11829290203)
		u := uint(67284)
		i8 := int8(-100)
		i16 := int16(-1203)
		i32 := int32(-5784182)
		i64 := int64(-11829290203)
		i := int(-67284)
		f32 := float32(16263.8888)
		f64 := float64(-16263.8888)
		inputs := []pointer.Ptr{
			{Str: &str, Bool: &flag, Time: &dt, F32: &f32, F64: &f64, Uint: &u, Uint8: &u8, Uint16: &u16, Uint32: &u32, Uint64: &u64, Int: &i, Int8: &i8, Int16: &i16, Int32: &i32, Int64: &i64},
			{Str: &str, Bool: &flag, Time: &dt, F32: &f32, F64: &f64, Uint: &u, Uint8: &u8, Uint16: &u16, Uint32: &u32, Uint64: &u64, Int: &i, Int8: &i8, Int16: &i16, Int32: &i32, Int64: &i64},
		}
		result, err := db.InsertInto(ctx, sqliteDB, inputs)
		require.NoError(t, err)
		lastID := mustValue(result.LastInsertId())
		require.NoError(t, err)
		require.NotEmpty(t, lastID)
		require.Equal(t, int64(len(inputs)), mustValue(result.RowsAffected()))

		ptr := pointer.Ptr{}
		ptr.ID = lastID
		mustNoError(db.FindOne(ctx, sqliteDB, &ptr))
		require.Equal(t, str, *ptr.Str)
		require.Equal(t, dt.Format(time.DateOnly), (*ptr.Time).Format(time.DateOnly))
		require.True(t, *ptr.Bool)
		require.Equal(t, u, *ptr.Uint)
		require.Equal(t, u8, *ptr.Uint8)
		require.Equal(t, u16, *ptr.Uint16)
		require.Equal(t, u32, *ptr.Uint32)
		require.Equal(t, u64, *ptr.Uint64)
		require.Equal(t, i, *ptr.Int)
		require.Equal(t, i8, *ptr.Int8)
		require.Equal(t, i16, *ptr.Int16)
		require.Equal(t, i32, *ptr.Int32)
		require.Equal(t, i64, *ptr.Int64)
		require.Equal(t, f32, *ptr.F32)
		require.Equal(t, f64, *ptr.F64)
	})
}

func TestDeleteOne(t *testing.T) {
	// ctx := context.TODO()
	// model := newPKModel()
	// _, err := sqlutil.InsertOne(ctx, sqliteDB, &model)
	// require.NoError(t, err)

	// models, err := sqlutil.SelectFrom[autopk.Model](ctx, sqliteDB)
	// require.NoError(t, err)

	// log.Println(models)
}
