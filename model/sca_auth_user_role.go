package model

import (
	"time"
)

const TableNameScaAuthUserRole = "sca_auth_user_role"

// ScaAuthUserRole 用户-角色映射表
type ScaAuthUserRole struct {
	ID          int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	UserID      string     `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                        // 用户ID
	RoleID      int64      `gorm:"column:role_id;type:bigint(20);not null;comment:角色ID" json:"role_id"`                          // 角色ID
	CreatedBy   *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                             // 创建人
	CreatedTime *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateBy    *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`                               // 更新人
	UpdateTime  *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted     *int64     `gorm:"column:deleted;type:int(11);comment:是否删除" json:"deleted"`                                      // 是否删除
}

// TableName ScaAuthUserRole's table name
func (*ScaAuthUserRole) TableName() string {
	return TableNameScaAuthUserRole
}
