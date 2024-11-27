package model

import (
	"time"
)

// ScaUserLevel 用户等级表
type ScaUserLevel struct {
	Id               int64      `gorm:"column:id;type:bigint(20);primary_key;comment:主键" json:"id"`
	UserId           string     `gorm:"column:user_id;type:varchar(20);comment:用户Id" json:"user_id"`
	LevelType        int        `gorm:"column:level_type;type:tinyint(1);comment:等级类型" json:"level_type"`
	Level            int        `gorm:"column:level;type:int(11);comment:等级" json:"level"`
	LevelName        string     `gorm:"column:level_name;type:varchar(50);comment:等级名称" json:"level_name"`
	ExpStart         int64      `gorm:"column:exp_start;type:bigint(20);comment:开始经验值" json:"exp_start"`
	ExpEnd           int64      `gorm:"column:exp_end;type:bigint(20);comment:结束经验值" json:"exp_end"`
	LevelDescription string     `gorm:"column:level_description;type:text;comment:等级描述" json:"level_description"`
	CreatedTime      *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"`
	UpdateTime       *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
}

func (m *ScaUserLevel) TableName() string {
	return "sca_user_level"
}
