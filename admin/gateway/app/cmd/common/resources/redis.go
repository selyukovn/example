package resources

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	assert "github.com/selyukovn/go-wm-assert"
)

func OpenRedis(host string, username string, password string, dbNumber uint) *redis.Client {
	assert.Str().NotEmpty().Must(host)
	assert.Str().NotEmpty().Must(username)
	assert.Str().NotEmpty().Must(password)

	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:6379?db=%d", username, password, host, dbNumber))
	if err != nil {
		panic(err)
	}

	r := redis.NewClient(opt)

	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return r
}
