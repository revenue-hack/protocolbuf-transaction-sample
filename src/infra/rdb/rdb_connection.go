package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/revenue-hack/protobuf-transaction-sample/src/config"
	"github.com/revenue-hack/protobuf-transaction-sample/src/db"
	_ "github.com/walf443/go-sql-tracer"
	"golang.org/x/xerrors"
)

var (
	once   sync.Once
	driver = "mysql"
)

func NewConnection() (*sql.DB, error) {
	var conn *sql.DB
	var err error
	once.Do(func() {
		protocol := fmt.Sprintf("tcp(%s:%s)", config.Env.DBHost, config.Env.DBPort)
		connInfo := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
			config.Env.DBUser, config.Env.DBPassword, protocol, config.Env.DBName)

		if config.IsLocal() {
			driver = "mysql:trace"
		}
		conn, err = sql.Open(driver, connInfo)
		if err != nil {
			return
		}

		if e := conn.Ping(); e != nil {
			err = e
			return
		}

		// DBへの全体のコネクション総数
		conn.SetMaxOpenConns(config.Env.RDBMaxConn)
		// DB接続を待機させておくコネクション総数
		conn.SetMaxIdleConns(config.Env.RDBMaxIdle)
	})

	return conn, err
}

func ExecFromCtx(ctx context.Context) (db.RDBHandler, error) {
	val := ctx.Value(config.DBKey)
	if val == nil {
		return nil, xerrors.New("fail to get connection from context")
	}

	conn, ok := val.(db.RDBHandler)
	if !ok {
		return nil, xerrors.New("can't get context executor")
	}
	return conn, nil
}
