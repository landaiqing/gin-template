package client_api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
)

// GenerateClientId 生成客户端ID
// @Summary 生成客户端ID
// @Description 生成客户端ID
// @Tags 微信公众号
// @Produce json
// @Router /api/oauth/generate_client_id [get]
func (ClientAPI) GenerateClientId(c *gin.Context) {
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
	v1 := uuid.NewV1()
	err := redis.Set(constant.UserLoginClientRedisKey+ip, v1.String(), 0).Err()
	if err != nil {
		global.LOG.Error(err)
		return
	}
	result.OkWithData(v1.String(), c)
	return
}
