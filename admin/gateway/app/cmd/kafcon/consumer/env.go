package main

import (
	assert "github.com/selyukovn/go-wm-assert"
	"os"
	"strconv"
	"strings"
)

type tEnv struct {
	MysqlHost             string
	MysqlUser             string
	MysqlPassword         string
	MysqlDb               string
	RedisCacheHost        string
	RedisCacheUser        string
	RedisCachePassword    string
	RedisCacheDb          uint
	KafkaBrokersHostPorts []string
}

func loadEnv() tEnv {
	redisDbUint64, redisDbErr := strconv.ParseUint(os.Getenv("REDIS_CACHE_DB"), 10, 64)
	assert.TrueMust(redisDbErr == nil, "env: REDIS_CACHE_DB")

	kafkaBrokersHostPorts := strings.Split(os.Getenv("KAFKA_BROKERS_HOSTPORTS"), ",")
	assert.SliceCmp[[]string, string]().LenMin(1).Uniques().CustomElementEach("each", func(s string) bool {
		return nil == assert.Str().NotEmpty(). /* todo : UrlHostPort(). */ Check(s)
	}).Must(kafkaBrokersHostPorts, "env: KAFKA_BROKERS_HOSTPORTS")

	return tEnv{
		MysqlHost:             assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST"), "env: MYSQL_HOST"),
		MysqlUser:             assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		MysqlPassword:         assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		MysqlDb:               assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
		RedisCacheHost:        assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_HOST"), "env: REDIS_CACHE_HOST"),
		RedisCacheUser:        assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_USER"), "env: REDIS_CACHE_USER"),
		RedisCachePassword:    assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_PASSWORD"), "env: REDIS_CACHE_PASSWORD"),
		RedisCacheDb:          uint(redisDbUint64),
		KafkaBrokersHostPorts: kafkaBrokersHostPorts,
	}
}
