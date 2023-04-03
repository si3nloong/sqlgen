package main_test

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"os"
// 	"testing"

// 	sqlutil "github.com/si3nloong/sqlgen/sql"
// 	"github.com/si3nloong/sqlgen/testdata/structfield/primitive"
// 	"github.com/stretchr/testify/require"
// )

// func openSqlConn(driver string) (*sql.DB, error) {
// 	switch driver {
// 	case "mysql":
// 		return sql.Open("mysql", "user:password@/dbname")
// 	// case "postgres":
// 	case "sqlite":
// 		os.Remove("./sqlite.db")
// 		return sql.Open("sqlite3", "./sqlite.db")
// 	default:
// 		return nil, errors.New("unsupported sql driver")
// 	}
// }

// func mustNot[T any](v T, err error) T {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return v
// }

// func TestMain(m *testing.M) {
// 	// openSqlConn("mysql")
// 	conn := mustNot(openSqlConn("sqlite"))
// 	defer conn.Close()

// 	m.Run()

// 	testing.Benchmark(func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {

// 		}
// 	})
// }

// func TestDatabases(t *testing.T) {
// 	conn := mustNot(openSqlConn("sqlite"))
// 	defer conn.Close()
// 	testSqlite(t, conn)
// }

// func testSqlite(t *testing.T, conn *sql.DB) {
// 	ctx := context.TODO()
// 	inputs := []primitive.Primitive{}
// 	result, err := sqlutil.InsertInto(ctx, conn, inputs)
// 	require.NoError(t, err)
// 	require.Equal(t, int64(0), mustNot(result.LastInsertId()))
// 	require.Equal(t, int64(0), mustNot(result.RowsAffected()))
// }

// func TestDeleteOne(t *testing.T) {
// 	//type args[T sqlutil.KeyValuer[T]] struct {
// 	//	ctx context.Context
// 	//	db  sqlutil.DB
// 	//	v   T
// 	//}
// 	//type testCase[T sqlutil.KeyValuer[T]] struct {
// 	//	name    string
// 	//	args    args[T]
// 	//	want    sql.Result
// 	//	wantErr bool
// 	//}
// 	//tests := []testCase[ /* TODO: Insert concrete types here */ ]{
// 	//	// TODO: Add test cases.
// 	//}
// 	//for _, tt := range tests {
// 	//	t.Run(tt.name, func(t *testing.T) {
// 	//		got, err := sqlutil.DeleteOne(tt.args.ctx, tt.args.db, tt.args.v)
// 	//		if (err != nil) != tt.wantErr {
// 	//			t.Errorf("DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
// 	//			return
// 	//		}
// 	//		if !reflect.DeepEqual(got, tt.want) {
// 	//			t.Errorf("DeleteOne() got = %v, want %v", got, tt.want)
// 	//		}
// 	//	})
// 	//}
// }
