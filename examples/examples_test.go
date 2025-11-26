package examples

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect/postgres"

	"log"

	"github.com/golang-migrate/migrate/v4"

	"github.com/jaswdr/faker"
	mysqldb "github.com/si3nloong/sqlgen/examples/db/mysql"
	"github.com/si3nloong/sqlgen/examples/testcase/core"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/pointer"
)

var (
	//go:embed migrate/*.sql
	migrationFiles embed.FS
	fake           = faker.New()
	conn           *sql.DB
)

func openSqlConn(driver, username, password string, addr string, dbname string) (*sql.DB, error) {
	switch driver {
	case "mysql":
		return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/sqlgen?parseTime=true", username, password, addr))
	case "sqlite3":
		// os.Remove("./sqlite.db")
		return sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	default:
		return nil, errors.New("unsupported sql driver")
	}
}

func TestMain(m *testing.M) {
	exitCode, err := initContainer(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func initContainer(m *testing.M) (int, error) {
	const (
		driver   = "mysql"
		dbname   = "sqlgen"
		password = "secret"
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		return 0, err
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return 0, fmt.Errorf("Could not connect to Docker: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run(driver, "8.0", []string{
		fmt.Sprintf("MYSQL_DATABASE=%s", dbname),
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password),
	})
	if err != nil {
		return 0, fmt.Errorf("Could not start resource: %w", err)
	}
	defer pool.Purge(resource)
	addr := resource.GetHostPort("3306/tcp")

	if err := pool.Retry(func() error {
		conn, err = openSqlConn(driver, "root", password, addr, dbname)
		if err != nil {
			return err
		}
		return conn.Ping()
	}); err != nil {
		return 0, err
	}
	defer conn.Close()

	d, err := iofs.New(migrationFiles, "migrate")
	if err != nil {
		return 0, err
	}
	instance, err := mysql.WithInstance(conn, &mysql.Config{})
	if err != nil {
		return 0, err
	}
	mg, err := migrate.NewWithInstance("iofs", d, "sqlgen", instance)
	if err != nil {
		return 0, fmt.Errorf("unable to create a sql instance: %w", err)
	}
	if err := mg.Up(); err != nil {
		return 0, fmt.Errorf("migration up failed: %w", err)
	}
	exitCode := m.Run()
	if err := mg.Down(); err != nil {
		return 0, fmt.Errorf("migration down failed: %w", err)
	}
	return exitCode, nil
}

func TestInsert(t *testing.T) {
	utcNow := time.Now().UTC()

	t.Run("InsertOne with AutoIncrement PK", func(t *testing.T) {
		u1 := core.User{}
		u1.Name = fake.Company().Name()
		u1.No = fake.UIntBetween(1, 100)
		u1.Address.Line1 = fake.Address().Address()
		u1.Address.Line2 = fake.Address().SecondaryAddress()
		u1.Address.CountryCode = fake.Address().CountryCode()
		postalCode := fake.Address().PostCode()
		u1.PostalCode = &postalCode
		u1.ExtraInfo.Flag = fake.Bool()
		u1.Slice = []float64{0.15, 0.3, 0.88}
		u1.Nicknames = [2]string{"John Pinto", "JP"}
		u1.Kind = reflect.String
		u1.T = utcNow
		u1.Map = map[string]float64{"a": 1, "b": 2}
		u1.JoinedTime = utcNow
		result, err := mysqldb.InsertOne(t.Context(), conn, &u1)
		require.NoError(t, err)
		affected, err := result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		u2 := core.User{}
		u2.ID, err = result.LastInsertId()
		require.NoError(t, err)

		require.NoError(t, mysqldb.FindByPK(t.Context(), conn, &u2))
		require.Equal(t, u1.ID, u2.ID)
		require.Equal(t, u1.No, u2.No)
		require.Equal(t, u1.ExtraInfo, u2.ExtraInfo)
		require.Equal(t, u1.Address, u2.Address)
		require.Equal(t, u1.PostalCode, u2.PostalCode)
		require.Equal(t, u1.Kind, u2.Kind)
		require.Equal(t, u1.Name, u2.Name)
		require.Equal(t, u1.Map, u2.Map)
		require.Equal(t, u1.T.Format(time.RFC822), u2.T.Format(time.RFC822))
		require.NotEmpty(t, u2.JoinedTime)
		require.ElementsMatch(t, u1.Nicknames, u2.Nicknames)
		require.ElementsMatch(t, u1.Slice, u2.Slice)
	})

	t.Run("InsertOne with pointer", func(t *testing.T) {
		ptr1 := pointer.Ptr{}
		ptr1.Str = ptrOf(fake.App().Name())
		ptr1.Bool = ptrOf(fake.Bool())
		ptr1.Int8 = ptrOf(fake.Int8Between(-88, 88))
		ptr1.Int16 = ptrOf(fake.Int16Between(-888, 888))
		ptr1.Int32 = ptrOf(fake.Int32Between(-888_888_888, 888_888_888))
		ptr1.Int64 = ptrOf(fake.Int64Between(-888_888_888_888, 888_888_888_888))
		ptr1.Int = ptrOf(fake.IntBetween(-888_888, 888_888))
		ptr1.Uint8 = ptrOf(fake.UInt8Between(0, 88))
		ptr1.Uint16 = ptrOf(fake.UInt16Between(0, 888))
		ptr1.Uint32 = ptrOf(fake.UInt32Between(0, 888_888_888))
		ptr1.Uint64 = ptrOf(fake.UInt64Between(0, 888_888_888_888_888))
		ptr1.Uint = ptrOf(fake.UIntBetween(0, 888_888))
		ptr1.Time = ptrOf(time.Now().UTC())
		ptr1.F32 = ptrOf(fake.Float32(0, 1, 100))
		ptr1.F64 = ptrOf(fake.Float64(0, 1, 88888))
		result, err := mysqldb.InsertOne(t.Context(), conn, &ptr1)
		require.NoError(t, err)
		lastID, err := result.LastInsertId()
		require.NoError(t, err)
		require.NotEmpty(t, lastID)

		ptr2 := pointer.Ptr{}
		ptr2.ID = lastID
		require.NoError(t, mysqldb.FindByPK(t.Context(), conn, &ptr2))
		require.Equal(t, ptr1.Str, ptr2.Str)
		require.Equal(t, ptr1.Bool, ptr2.Bool)
		require.Equal(t, ptr1.Int, ptr2.Int)
		require.Equal(t, ptr1.Int8, ptr2.Int8)
		require.Equal(t, ptr1.Int16, ptr2.Int16)
		require.Equal(t, ptr1.Int32, ptr2.Int32)
		require.Equal(t, ptr1.Int64, ptr2.Int64)
		require.Equal(t, ptr1.Uint, ptr2.Uint)
		require.Equal(t, ptr1.Uint8, ptr2.Uint8)
		require.Equal(t, ptr1.Uint16, ptr2.Uint16)
		require.Equal(t, ptr1.Uint32, ptr2.Uint32)
		require.Equal(t, ptr1.Uint64, ptr2.Uint64)
		require.Equal(t, ptr1.Time.Format(time.RFC822), ptr2.Time.Format(time.RFC822))
		require.Nil(t, ptr2.Nested)
	})

	// t.Run("Insert with array", func(t *testing.T) {
	// 	r1 := slice.Slice{}
	// 	r1.StrList = []string{"a", "b", "c"}
	// 	r1.CustomStrList = append(r1.CustomStrList, "x", "y", "z")
	// 	r1.BoolList = append(r1.BoolList, true, false, true, false, true)
	// 	r1.Int8List = append(r1.Int8List, -88, -13, -1, 6)
	// 	r1.Int32List = append(r1.Int32List, -88, 188, -1)
	// 	r1.Uint8List = append(r1.Uint8List, 10, 5, 1)
	// 	r1.F32List = append(r1.F32List, -88.114, 188.123, -1.0538)
	// 	r1.F64List = append(r1.F64List, -88.114, 188.123, -1.0538)

	// 	inputs := []slice.Slice{r1}
	// 	result, err := mysqldb.Insert(context.TODO(), dbConn, inputs)
	// 	require.NoError(t, err)
	// 	lastID := mustValue(result.LastInsertId())
	// 	require.NotEmpty(t, lastID)

	// 	ptr := slice.Slice{}
	// 	ptr.ID = uint64(lastID)
	// 	mustNoError(mysqldb.FindByPK(t.Context(), dbConn, &ptr))
	// })

	// t.Run("Insert with all nil values", func(t *testing.T) {
	// 	inputs := []pointer.Ptr{{}, {}}
	// 	result, err := mysqldb.Insert(t.Context(), dbConn, inputs)
	// 	require.NoError(t, err)
	// 	lastID := mustValue(result.LastInsertId())
	// 	require.NotEmpty(t, lastID)
	// 	require.Equal(t, int64(2), mustValue(result.RowsAffected()))
	// })

	// t.Run("Insert with pointer values", func(t *testing.T) {
	// 	str := "hello world"
	// 	flag := true
	// 	dt := time.Now().UTC()
	// 	u8 := uint8(100)
	// 	u16 := uint16(1203)
	// 	u32 := uint32(5784182)
	// 	u64 := uint64(11829290203)
	// 	u := uint(67284)
	// 	i8 := int8(-100)
	// 	i16 := int16(-1203)
	// 	i32 := int32(-5784182)
	// 	i64 := int64(-11829290203)
	// 	i := int(-67284)
	// 	f32 := float32(16263.8888)
	// 	f64 := float64(-16263.8888)
	// 	inputs := []pointer.Ptr{
	// 		{Str: &str, Bool: &flag, Time: &dt, F32: &f32, F64: &f64, Uint: &u, Uint8: &u8, Uint16: &u16, Uint32: &u32, Uint64: &u64, Int: &i, Int8: &i8, Int16: &i16, Int32: &i32, Int64: &i64},
	// 		{Str: &str, Bool: &flag, Time: &dt, F32: &f32, F64: &f64, Uint: &u, Uint8: &u8, Uint16: &u16, Uint32: &u32, Uint64: &u64, Int: &i, Int8: &i8, Int16: &i16, Int32: &i32, Int64: &i64},
	// 	}
	// 	result, err := mysqldb.Insert(t.Context(), dbConn, inputs)
	// 	require.NoError(t, err)
	// 	lastID := mustValue(result.LastInsertId())
	// 	require.NoError(t, err)
	// 	require.NotEmpty(t, lastID)
	// 	require.Equal(t, int64(len(inputs)), mustValue(result.RowsAffected()))

	// 	ptrs, err := mysqldb.QueryStmt(t.Context(), dbConn, func(p pointer.Ptr) mysqldb.SelectStmt {
	// 		return mysqldb.SelectStmt{
	// 			Select:    p.Columns(),
	// 			FromTable: p.TableName(),
	// 			Where:     mysqldb.Equal(p.ColumnInt(), &i),
	// 			Limit:     3,
	// 		}
	// 	})
	// 	require.NotEmpty(t, ptrs)
	// 	require.NoError(t, err)
	// })
}

func TestUpdateOne(t *testing.T) {
	// data := autopk.Model{}
	// result, err := mysqldb.InsertOne(t.Context(), dbConn, &data)
	// if err != nil {
	// 	panic(err)
	// }

	// i64, _ := result.LastInsertId()
	// newData := autopk.Model{}
	// newData.ID = uint(i64)
	// newData.Name = autopk.LongText(`Updated Text`)

	// if _, err := mysqldb.UpdateByPK(t.Context(), dbConn, newData); err != nil {
	// 	panic(err)
	// }
}

func TestDeleteOne(t *testing.T) {
	// ctx := context.TODO()
	// model := newPKModel()
	// _, err := sqlutil.InsertOne(ctx, dbConn, &model)
	// require.NoError(t, err)

	// models, err := sqlutil.SelectFrom[autopk.Model](ctx, dbConn)
	// require.NoError(t, err)
}

func TestPaginate(t *testing.T) {
	// t.Run("Without cursor", func(t *testing.T) {
	// 	p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
	// 	p.Next(t.Context(), conn)
	// })

	// t.Run("With cursor", func(t *testing.T) {
	// 	p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
	// 	p.Next(t.Context(), conn)
	// })

	// t.Run(`With "WHERE" clause`, func(t *testing.T) {
	// 	p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
	// 	p.Next(t.Context(), conn)
	// })

	// t.Run(`With "ORDER BY" clause`, func(t *testing.T) {
	// 	p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
	// 	p.Next(t.Context(), conn)
	// })

	// for v, err := range p.Next(t.Context(), dbConn) {
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	_ = v
	// }
}

func ptrOf[T any](v T) *T {
	return &v
}
