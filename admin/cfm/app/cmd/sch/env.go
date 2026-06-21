package main

import (
	assert "github.com/selyukovn/go-wm-assert"
	"os"
)

type tEnv struct {
	MysqlHost     string
	MysqlUser     string
	MysqlPassword string
	MysqlDb       string
}

func loadEnv() tEnv {
	return tEnv{
		MysqlHost:     assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST"), "env: MYSQL_HOST"),
		MysqlUser:     assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		MysqlPassword: assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		MysqlDb:       assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	}
}
