package client_controller

import (
	"sync"
	"time"

	"github.com/ccpwcn/kgo"
	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
)

type ClientController struct{}

var mu sync.Mutex

// GenerateClientId 生成客户端ID
// @Summary 生成客户端ID
// @Description 生成客户端ID
// @Tags 微信公众号
// @Produce json
// @Router /controller/oauth/generate_client_id [get]
func (ClientController) GenerateClientId(c *gin.Context) {
	// 获取客户端IP
	ip := utils.GetClientIP(c)
	// 加锁
	mu.Lock()
	defer mu.Unlock()

	// 从Redis获取客户端ID
	clientId := redis.Get(constant.UserLoginClientRedisKey + ip).Val()
	if clientId != "" {
		result.OkWithData(clientId, c)
		return
	}
	// 生成新的客户端ID
	simpleUuid := kgo.SimpleUuid()
	err := redis.Set(constant.UserLoginClientRedisKey+ip, simpleUuid, time.Hour*24*7).Err()
	if err != nil {
		global.LOG.Error(err)
		return
	}
	result.OkWithData(simpleUuid, c)
	return
}
