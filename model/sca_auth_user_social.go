package model

import (
	"time"
)

const TableNameScaAuthUserSocial = "sca_auth_user_social"

// ScaAuthUserSocial 社会用户信息表
type ScaAuthUserSocial struct {
	ID          int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	UserID      *string    `gorm:"column:user_id;type:varchar(20);not null;comment:用户ID" json:"user_id"`                         // 用户ID
	Source      *string    `gorm:"column:source;type:varchar(10);comment:第三方用户来源" json:"source"`                                 // 第三方用户来源
	OpenID      *string    `gorm:"column:open_id;type:varchar(50);comment:第三方用户的 open id" json:"open_id"`                        // 第三方用户的 open id
	Status      *int64     `gorm:"column:status;type:int(11);default:0;comment:状态 0正常 1 封禁" json:"status"`                       // 状态 0正常 1 封禁
	CreatedTime *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateTime  *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted     *int64     `gorm:"column:deleted;type:int(11);default:0;comment:是否删除" json:"deleted"`                            // 是否删除
}

// TableName ScaAuthUserSocial's table name
func (*ScaAuthUserSocial) TableName() string {
	return TableNameScaAuthUserSocial
}
