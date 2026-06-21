package main

import (
	assert "github.com/selyukovn/go-wm-assert"
	"os"
)

type tEnv struct {
	MysqlHostMaster string
	MysqlUser       string
	MysqlPassword   string
	MysqlDb         string
}

func loadEnv() tEnv {
	return tEnv{
		MysqlHostMaster: assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST_MASTER"), "env: MYSQL_HOST_MASTER"),
		MysqlUser:       assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		MysqlPassword:   assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		MysqlDb:         assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	}
}
