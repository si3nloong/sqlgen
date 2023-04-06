package examples_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jaswdr/faker"
	"github.com/si3nloong/sqlgen/examples/testdata/structfield/autopk"
	sqlutil "github.com/si3nloong/sqlgen/sql"
	"github.com/stretchr/testify/require"
)

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

const createTableModel = `
CREATE TABLE IF NOT EXISTS ` + "`model`" + ` (
	name VARCHAR(100), 
	f TINYINT, 
	id INT PRIMARY KEY AUTO_INCREMENT,
	n INT
)
`

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

func mustNot[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	sqliteDB *sql.DB
)

func TestMain(m *testing.M) {
	// openSqlConn("mysql")
	sqliteDB = mustNot(openSqlConn("mysql"))
	defer sqliteDB.Close()

	// mustNot(sqliteDB.Exec("DROP TABLE `model`;"))
	mustNot(sqliteDB.Exec(createTableModel))

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
	inputs := []autopk.Model{newPKModel(), newPKModel(), newPKModel()}
	result, err := sqlutil.InsertInto(ctx, sqliteDB, inputs)
	require.NoError(t, err)
	// require.Equal(t, int64(0), mustNot(result.LastInsertId()))
	require.Equal(t, int64(3), mustNot(result.RowsAffected()))
}

func TestDeleteOne(t *testing.T) {
	ctx := context.TODO()
	model := newPKModel()
	_, err := sqlutil.InsertOne(ctx, sqliteDB, &model)
	require.NoError(t, err)

	models, err := sqlutil.SelectFrom[autopk.Model](ctx, sqliteDB)
	require.NoError(t, err)

	log.Println(models)
}
