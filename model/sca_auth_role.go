package model

import (
	"encoding/json"
	"time"
)

const TableNameScaAuthRole = "sca_auth_role"

// ScaAuthRole 角色表
type ScaAuthRole struct {
	ID          int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	RoleName    string     `gorm:"column:role_name;type:varchar(32);not null;comment:角色名称" json:"role_name"`                     // 角色名称
	RoleKey     string     `gorm:"column:role_key;type:varchar(64);not null;comment:角色关键字" json:"role_key"`                      // 角色关键字
	CreatedTime *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateTime  *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted     *int64     `gorm:"column:deleted;type:int(11);comment:是否删除 0 未删除 1已删除" json:"deleted"`                           // 是否删除 0 未删除 1已删除
	CreatedBy   *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                             // 创建人
	UpdateBy    *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`
}

// TableName ScaAuthRole's table name
func (*ScaAuthRole) TableName() string {
	return TableNameScaAuthRole
}

func (role *ScaAuthRole) MarshalBinary() ([]byte, error) {
	return json.Marshal(role)
}

func (role *ScaAuthRole) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, role)
}
