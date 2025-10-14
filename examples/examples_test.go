package examples

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	_ "github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect/postgres"

	"log"

	"github.com/jaswdr/faker"
	mysqldb "github.com/si3nloong/sqlgen/examples/db/mysql"
	"github.com/si3nloong/sqlgen/examples/testcase/core"
	autopk "github.com/si3nloong/sqlgen/examples/testcase/struct-field/pk/auto-incr"
)

var (
	// migrationFiles embed.FS
	sqlConn *sql.DB
)

func openSqlConn(driver string) (*sql.DB, error) {
	switch driver {
	case "mysql":
		return sql.Open("mysql", "root:abcd1234@/sqlbench?parseTime=true")
	case "sqlite3":
		// os.Remove("./sqlite.db")
		return sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	default:
		return nil, errors.New("unsupported sql driver")
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		sqlConn, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return sqlConn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// b, err := migrationFiles.ReadFile("migrations/1_create_tables.up.sql")
	// if err != nil {
	// 	log.Fatalf("Unable to find migration file: %s", err)
	// }
	// sqlConn.Exec("DROP TABLE IF EXISTS `user`;")
	// if _, err := sqlConn.Exec(string(b)); err != nil {
	// 	log.Fatalf("Cannot do migration: %s", err)
	// }

	os.Exit(m.Run())
}

func newPKModel() autopk.Model {
	fake := faker.New()
	return autopk.Model{
		F:    true,
		Name: autopk.LongText(fake.Person().Name()),
		N:    fake.Int64Between(0, 100),
	}
}

func TestInsert(t *testing.T) {
	fake := faker.New()

	t.Run("Insert", func(t *testing.T) {
		utcNow := time.Now()
		u1 := core.User{}
		u1.No = fake.UIntBetween(1, 100)
		u1.Address.Line1 = fake.Address().Address()
		u1.Address.Line2 = fake.Address().SecondaryAddress()
		u1.Address.CountryCode = fake.Address().CountryCode()
		postalCode := fake.Address().PostCode()
		u1.PostalCode = &postalCode
		u1.ExtraInfo.Flag = fake.Bool()
		u1.Nicknames = [2]string{"John Pinto", "JP"}
		u1.Kind = reflect.String
		u1.JoinedTime = utcNow
		result, err := mysqldb.InsertOne(t.Context(), sqlConn, &u1)
		require.NoError(t, err)
		affected, err := result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		u2 := core.User{}
		u2.ID, err = result.LastInsertId()
		require.NoError(t, err)

		require.NoError(t, mysqldb.FindByPK(t.Context(), sqlConn, &u2))
		require.Equal(t, u1.ID, u2.ID)
		require.Equal(t, u1.No, u2.No)
		require.Equal(t, u1.ExtraInfo, u2.ExtraInfo)
		require.Equal(t, u1.Address, u2.Address)
		require.Equal(t, u1.PostalCode, u2.PostalCode)
		require.Equal(t, u1.Kind, u2.Kind)
		require.NotEmpty(t, u2.JoinedTime)
		require.ElementsMatch(t, u1.Nicknames, u2.Nicknames)
	})

	// t.Run("Insert with double ptr", func(t *testing.T) {
	// 	u8 := uint(188)
	// 	str := "Hello, james!"
	// 	cStr := doubleptr.LongStr(`Hi, bye`)
	// 	data := doubleptr.DoublePtr{}
	// 	data.L3PtrUint = ptrOf(ptrOf(ptrOf(u8)))
	// 	data.L3PtrCustomStr = ptrOf(ptrOf(ptrOf(cStr)))
	// 	data.L7PtrStr = ptrOf(ptrOf(ptrOf(ptrOf(ptrOf(ptrOf(ptrOf(str)))))))
	// 	inputs := []doubleptr.DoublePtr{data}
	// 	result, err := mysqldb.Insert(context.TODO(), dbConn, inputs)
	// 	require.NoError(t, err)
	// 	lastID := mustValue(result.LastInsertId())
	// 	require.NotEmpty(t, lastID)
	// })

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

	// 	ptr := pointer.Ptr{}
	// 	ptr.ID = lastID
	// 	mustNoError(mysqldb.FindByPK(t.Context(), dbConn, &ptr))
	// 	require.Equal(t, str, *ptr.Str)
	// 	require.Equal(t, dt.Format(time.DateOnly), (*ptr.Time).Format(time.DateOnly))
	// 	require.True(t, *ptr.Bool)
	// 	require.Equal(t, u, *ptr.Uint)
	// 	require.Equal(t, u8, *ptr.Uint8)
	// 	require.Equal(t, u16, *ptr.Uint16)
	// 	require.Equal(t, u32, *ptr.Uint32)
	// 	require.Equal(t, u64, *ptr.Uint64)
	// 	require.Equal(t, i, *ptr.Int)
	// 	require.Equal(t, i8, *ptr.Int8)
	// 	require.Equal(t, i16, *ptr.Int16)
	// 	require.Equal(t, i32, *ptr.Int32)
	// 	require.Equal(t, i64, *ptr.Int64)
	// 	require.NotZero(t, *ptr.F32)
	// 	require.NotZero(t, *ptr.F64)

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
	t.Run("Without cursor", func(t *testing.T) {
		p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
		p.Next(t.Context(), sqlConn)
	})

	t.Run("With cursor", func(t *testing.T) {
		p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
		p.Next(t.Context(), sqlConn)
	})

	t.Run(`With "WHERE" clause`, func(t *testing.T) {
		p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
		p.Next(t.Context(), sqlConn)
	})

	t.Run(`With "ORDER BY" clause`, func(t *testing.T) {
		p := mysqldb.Paginate[core.User](mysqldb.PaginateStmt{})
		p.Next(t.Context(), sqlConn)
	})

	// for v, err := range p.Next(t.Context(), dbConn) {
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	_ = v
	// }
}
