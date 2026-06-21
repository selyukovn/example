package main

import (
	assert "github.com/selyukovn/go-wm-assert"
	"os"
	"strconv"
)

type tEnv = struct {
	RedisCacheHost            string
	RedisCacheUser            string
	RedisCachePassword        string
	RedisCacheDb              uint
	ServiceAuthApiGrpcBaseUrl string
	ServiceAuthApiGrpcApiKey  string
}

func loadEnv() tEnv {
	redisDbUint64, redisDbErr := strconv.ParseUint(os.Getenv("REDIS_CACHE_DB"), 10, 64)
	assert.TrueMust(redisDbErr == nil, "env: REDIS_CACHE_DB")

	return tEnv{
		RedisCacheHost:            assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_HOST"), "env: REDIS_CACHE_HOST"),
		RedisCacheUser:            assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_USER"), "env: REDIS_CACHE_USER"),
		RedisCachePassword:        assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_PASSWORD"), "env: REDIS_CACHE_PASSWORD"),
		RedisCacheDb:              uint(redisDbUint64),
		ServiceAuthApiGrpcBaseUrl: assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_AUTH_API_GRPC_BASEURL"), "env: SERVICE_AUTH_API_GRPC_BASEURL"),
		ServiceAuthApiGrpcApiKey:  assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_AUTH_API_GRPC_APIKEY"), "env: SERVICE_AUTH_API_GRPC_APIKEY"),
	}
}
