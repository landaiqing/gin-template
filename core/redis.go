package core

import (
	"github.com/redis/go-redis/v9"

	"schisandra-cloud-album/global"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:           global.CONFIG.Redis.Addr(),
		Password:       global.CONFIG.Redis.Password,
		DB:             global.CONFIG.Redis.Db,
		MaxActiveConns: global.CONFIG.Redis.MaxActive,
		MaxIdleConns:   global.CONFIG.Redis.MaxIdle,
		PoolSize:       global.CONFIG.Redis.PoolSize,
		MinIdleConns:   global.CONFIG.Redis.MinIdle,
		PoolTimeout:    global.CONFIG.Redis.PoolTimeout,
		ReadTimeout:    -1,
	})
	InitSession(rdb)
	global.REDIS = rdb
}
