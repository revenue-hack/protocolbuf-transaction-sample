package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	DBKey = "conn"
)

type ENV struct {
	DBName     string `envconfig:"PSQL_DBNAME" required:"true"`
	DBHost     string `envconfig:"PSQL_HOST" required:"true"`
	DBPort     string `envconfig:"PSQL_PORT" default:"5432"`
	DBUser     string `envconfig:"PSQL_USER" required:"true"`
	DBPassword string `envconfig:"PSQL_PASS" required:"true"`
	GOEnv      string `envconfig:"GO_ENV" default:"local"`
	CorsHost   string `envconfig:"CORS_HOST" required:"true"`
	GRPCPort   int64  `envconfig:"GRPC_PORT" default:"9090"`
	RDBMaxIdle int    `envconfig:"RDB_MAX_IDLE" default:"10"` // DBへの待機させておく最大のコネクション総数
	RDBMaxConn int    `envconfig:"RDB_MAX_CONN" default:"30"` // DBへこのアプリケーションが接続できる最大コネクション数
}

var Env ENV

func init() {
	env := os.Getenv("GO_ENV")

	if err := envconfig.Process(env, &Env); err != nil {
		panic(err)
	}
}

func IsLocal() bool {
	return Env.GOEnv == "local"
}
