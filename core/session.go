package core

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"net/http"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/global"
)

func InitSession(client *redis.Client) {
	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		global.LOG.Fatal("failed to create redis store: ", err)
	}

	// Example changing configuration for sessions
	store.KeyPrefix(constant.UserSessionRedisKey)
	store.Options(sessions.Options{
		Path: "/",
		//Domain: global.CONFIG.System.Web,
		MaxAge:      86400 * 7,
		HttpOnly:    true,
		Secure:      true,
		Partitioned: true,
		SameSite:    http.SameSiteLaxMode,
	})
	global.Session = store
}
