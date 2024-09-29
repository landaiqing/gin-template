package impl

import (
	"encoding/json"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"time"
)

var userDao = impl.UserDaoImpl{}

type UserServiceImpl struct{}

// ResponseData 返回数据
type ResponseData struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresAt    int64   `json:"expires_at"`
	UID          *string `json:"uid"`
}

func (res ResponseData) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res ResponseData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &res)
}

// GetUserListService 返回用户列表
func (UserServiceImpl) GetUserListService() []*model.ScaAuthUser {
	return userDao.GetUserList()
}

// QueryUserByUsernameService 根据用户名查询用户
func (UserServiceImpl) QueryUserByUsernameService(username string) model.ScaAuthUser {
	return userDao.QueryUserByUsername(username)
}

// QueryUserByUuidService 根据uid查询用户
func (UserServiceImpl) QueryUserByUuidService(uid *string) model.ScaAuthUser {
	user, err := userDao.QueryUserByUuid(uid)
	if err != nil {
		return model.ScaAuthUser{}
	}
	return user
}

// DeleteUserService 根据uid删除用户
func (UserServiceImpl) DeleteUserService(uid string) error {
	return userDao.DeleteUser(uid)
}

// QueryUserByPhoneService 根据手机号查询用户
func (UserServiceImpl) QueryUserByPhoneService(phone string) model.ScaAuthUser {
	return userDao.QueryUserByPhone(phone)
}

// QueryUserByEmailService 根据邮箱查询用户
func (UserServiceImpl) QueryUserByEmailService(email string) model.ScaAuthUser {
	return userDao.QueryUserByEmail(email)
}

// AddUserService 新增用户
func (UserServiceImpl) AddUserService(user model.ScaAuthUser) (*model.ScaAuthUser, error) {
	return userDao.AddUser(user)
}

// UpdateUserService 更新用户信息
func (UserServiceImpl) UpdateUserService(phone, encrypt string) error {
	return userDao.UpdateUser(phone, encrypt)
}

// RefreshTokenService 刷新用户token
func (UserServiceImpl) RefreshTokenService(refreshToken string) (*ResponseData, bool) {
	parseRefreshToken, isUpd, err := utils.ParseRefreshToken(refreshToken)
	if err != nil || !isUpd {
		global.LOG.Errorln(err)
		return nil, false
	}
	accessTokenString, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: parseRefreshToken.UserID})
	if err != nil {
		return nil, false
	}
	tokenKey := constant.UserLoginTokenRedisKey + *parseRefreshToken.UserID
	token, err := redis.Get(tokenKey).Result()
	if err != nil || token == "" {
		global.LOG.Errorln(err)
		return nil, false
	}
	data := ResponseData{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		UID:          parseRefreshToken.UserID,
	}
	if err = redis.Set(tokenKey, data, time.Hour*24*7).Err(); err != nil {
		global.LOG.Errorln(err)
		return nil, false
	}
	return &data, true
}
