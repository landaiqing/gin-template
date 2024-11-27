package model

import (
	"time"
)

// ScaUserFollows 用户关注表
type ScaUserFollows struct {
	FollowerId  string     `gorm:"column:follower_id;type:varchar(20);primary_key;comment:关注者" json:"follower_id"`
	FolloweeId  string     `gorm:"column:followee_id;type:varchar(20);comment:被关注者;NOT NULL" json:"followee_id"`
	Status      int        `gorm:"column:status;type:tinyint(1);default:0;comment:关注状态（0 未互关 1 互关）" json:"status"`
	CreatedTime *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"`
	UpdateTime  *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
}

// TableName ScaUserFollows 表名
func (ScaUserFollows) TableName() string {
	return "sca_user_follows"
}
