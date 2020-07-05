package repoimpl_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/rdb"
)

var (
	conn *sql.DB
)

func setup() {
	db, err := rdb.NewConnection()
	if err != nil {
		panic(err)
	}
	conn = db

}

func tearDown() {
	err := conn.Close()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	defer tearDown()
	os.Exit(m.Run())
}
