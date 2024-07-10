package examples

import (
	"context"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"

	"github.com/jaswdr/faker"
	mysqldb "github.com/si3nloong/sqlgen/examples/db/mysql"
	autopk "github.com/si3nloong/sqlgen/examples/testcase/struct-field/pk/auto-incr"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/pointer"
	"github.com/si3nloong/sqlgen/examples/testcase/struct-field/slice"
)

func TestMain(m *testing.M) {
	// openSqlConn("mysql")
	// ctx := context.Background()
	dbConn = mustValue(openSqlConn("mysql"))
	defer dbConn.Close()

	// m1 := autopk.Model{}
	// sqlutil.FindOne(nil, nil, &m1)

	// if _, err := dbConn.ExecContext(ctx, autopk.Model{}.CreateTableStmt()); err != nil {
	// 	panic(err)
	// }

	// if _, err := dbConn.ExecContext(ctx, pointer.Ptr{}.CreateTableStmt()); err != nil {
	// 	panic(err)
	// }

	// if _, err := dbConn.ExecContext(ctx, array.Array{}.CreateTableStmt()); err != nil {
	// 	panic(err)
	// }

	// mustNot(dbConn.Exec("DROP TABLE `model`;"))
	// mustNot(dbConn.Exec(createTableModel))

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

func TestInsert(t *testing.T) {
	ctx := context.TODO()

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

	t.Run("Insert with array", func(t *testing.T) {
		r1 := slice.Slice{}
		r1.StrList = []string{"a", "b", "c"}
		r1.CustomStrList = append(r1.CustomStrList, "x", "y", "z")
		r1.BoolList = append(r1.BoolList, true, false, true, false, true)
		r1.Int8List = append(r1.Int8List, -88, -13, -1, 6)
		r1.Int32List = append(r1.Int32List, -88, 188, -1)
		r1.Uint8List = append(r1.Uint8List, 10, 5, 1)
		r1.F32List = append(r1.F32List, -88.114, 188.123, -1.0538)
		r1.F64List = append(r1.F64List, -88.114, 188.123, -1.0538)

		inputs := []slice.Slice{r1}
		result, err := mysqldb.Insert(context.TODO(), dbConn, inputs)
		require.NoError(t, err)
		lastID := mustValue(result.LastInsertId())
		require.NotEmpty(t, lastID)

		ptr := slice.Slice{}
		ptr.ID = uint64(lastID)
		mustNoError(mysqldb.FindByPK(ctx, dbConn, &ptr))
	})

	t.Run("Insert with all nil values", func(t *testing.T) {
		inputs := []pointer.Ptr{{}, {}}
		result, err := mysqldb.Insert(ctx, dbConn, inputs)
		require.NoError(t, err)
		lastID := mustValue(result.LastInsertId())
		require.NotEmpty(t, lastID)
		require.Equal(t, int64(2), mustValue(result.RowsAffected()))
	})

	t.Run("Insert with pointer values", func(t *testing.T) {
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
		result, err := mysqldb.Insert(ctx, dbConn, inputs)
		require.NoError(t, err)
		lastID := mustValue(result.LastInsertId())
		require.NoError(t, err)
		require.NotEmpty(t, lastID)
		require.Equal(t, int64(len(inputs)), mustValue(result.RowsAffected()))

		ptr := pointer.Ptr{}
		ptr.ID = lastID
		mustNoError(mysqldb.FindByPK(ctx, dbConn, &ptr))
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
		require.NotZero(t, *ptr.F32)
		require.NotZero(t, *ptr.F64)

		ptrs, err := mysqldb.QueryStmt[pointer.Ptr](ctx, dbConn, mysqldb.SelectStmt{
			Select:    ptr.Columns(),
			FromTable: ptr.TableName(),
			Where:     mysqldb.Equal(ptr.GetInt(), &i),
			Limit:     3,
		})
		_ = ptrs
		require.NoError(t, err)
	})
}

func TestUpdateOne(t *testing.T) {
	var (
		ctx = context.Background()
	)

	data := autopk.Model{}
	result, err := mysqldb.InsertOne(ctx, dbConn, &data)
	if err != nil {
		panic(err)
	}

	i64, _ := result.LastInsertId()
	newData := autopk.Model{}
	newData.ID = uint(i64)
	newData.Name = autopk.LongText(`Updated Text`)

	if _, err := mysqldb.UpdateByPK(ctx, dbConn, newData); err != nil {
		panic(err)
	}
}

func TestDeleteOne(t *testing.T) {
	// ctx := context.TODO()
	// model := newPKModel()
	// _, err := sqlutil.InsertOne(ctx, dbConn, &model)
	// require.NoError(t, err)

	// models, err := sqlutil.SelectFrom[autopk.Model](ctx, dbConn)
	// require.NoError(t, err)

	// log.Println(models)
}
