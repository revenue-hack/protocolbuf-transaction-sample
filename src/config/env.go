package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	DBKey = "conn"
)

type ENV struct {
	DBName     string `envconfig:"DB_NAME" required:"true"`
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     string `envconfig:"DB_PORT" default:"3346"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASS" required:"true"`
	GOEnv      string `envconfig:"GO_ENV" default:"local"`
	RDBMaxIdle int    `envconfig:"RDB_MAX_IDLE" default:"10"`
	RDBMaxConn int    `envconfig:"RDB_MAX_CONN" default:"30"`
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
