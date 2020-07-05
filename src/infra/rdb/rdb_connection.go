package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/revenue-hack/protobuf-transaction-sample/src/config"
	_ "github.com/walf443/go-sql-tracer"
	"golang.org/x/xerrors"
)

var (
	once   sync.Once
	driver = "postgres"
)

func NewConnection() (*sql.DB, error) {
	var conn *sql.DB
	var err error
	once.Do(func() {
		connInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			config.Env.DBHost, config.Env.DBPort, config.Env.DBName, config.Env.DBUser, config.Env.DBPassword)
		if config.IsLocal() || config.IsTest() {
			connInfo += " sslmode=disable"
		}

		if config.IsLocal() {
			driver = "postgres:trace"
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

		loc, e := time.LoadLocation("Asia/Tokyo")
		if e != nil {
			err = e
			return
		}
		boil.SetLocation(loc)
		//boil.SetDB(conn)
	})

	return conn, err
}

func ExecFromCtx(ctx context.Context) (boil.ContextExecutor, error) {
	val := ctx.Value(config.DBKey)
	if val == nil {
		return nil, dserr.Internal(xerrors.New("fail to get connection from context"))
	}

	conn, ok := val.(boil.ContextExecutor)
	if !ok {
		return nil, dserr.Internal(xerrors.New("can't get context executor"))
	}
	return conn, nil
}
