package utils

import (
	"encoding/json"
	"fmt"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/global"
	"strconv"
)

// CheckSlideData 校验滑动验证码
func CheckSlideData(point []int64, key string) bool {
	if point == nil || key == "" {
		return false
	}
	cacheDataByte, err := redis.Get(constant.CommentSubmitCaptchaRedisKey + key).Bytes()
	if len(cacheDataByte) == 0 || err != nil {
		return false
	}
	var dct *slide.Block
	if err = json.Unmarshal(cacheDataByte, &dct); err != nil {
		return false
	}

	chkRet := false
	if 2 == len(point) {
		sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", point[0]), 64)
		sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", point[1]), 64)
		chkRet = slide.CheckPoint(int64(sx), int64(sy), int64(dct.X), int64(dct.Y), 4)
	}
	if chkRet {
		return true
	}
	return false
}

// CheckRotateData 校验旋转验证码
func CheckRotateData(angle string, key string) bool {
	if angle == "" || key == "" {
		return false
	}
	cacheDataByte, err := redis.Get(constant.UserLoginCaptchaRedisKey + key).Bytes()
	if err != nil || len(cacheDataByte) == 0 {
		global.LOG.Error(err)
		return false
	}
	var dct *rotate.Block
	if err = json.Unmarshal(cacheDataByte, &dct); err != nil {
		global.LOG.Error(err)
		return false
	}
	sAngle, err := strconv.ParseFloat(fmt.Sprintf("%v", angle), 64)
	if err != nil {
		global.LOG.Error(err)
		return false
	}
	chkRet := rotate.CheckAngle(int64(sAngle), int64(dct.Angle), 2)
	if chkRet {
		return true
	}
	return false

}
