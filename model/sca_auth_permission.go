package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"schisandra-cloud-album/global"
	"time"
)

const TableNameScaAuthPermission = "sca_auth_permission"

// ScaAuthPermission 权限表
type ScaAuthPermission struct {
	ID             int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	PermissionName *string    `gorm:"column:permission_name;type:varchar(64);comment:权限名称" json:"permission_name"`                  // 权限名称
	ParentID       *int64     `gorm:"column:parent_id;type:bigint(20);comment:父ID" json:"parent_id"`                                // 父ID
	Type           *int64     `gorm:"column:type;type:tinyint(4);comment:类型 0 菜单 1 接口" json:"type"`                                 // 类型 0 菜单 1 目录 2 按钮 -1其他
	Path           *string    `gorm:"column:path;type:varchar(30);comment:路径" json:"path"`                                          // 路径
	Status         *int64     `gorm:"column:status;type:tinyint(4);comment:状态 0 启用 1 停用" json:"status"`                             // 状态 0 启用 1 停用
	Method         *string    `gorm:"column:method;type:varchar(20);comment:请求方式" json:"method"`                                    // 请求方式
	Icon           *string    `gorm:"column:icon;type:varchar(128);comment:图标" json:"icon"`                                         // 图标
	PermissionKey  *string    `gorm:"column:permission_key;type:varchar(64);comment:权限关键字" json:"permission_key"`                   // 权限关键字
	Order          *int64     `gorm:"column:order;type:int(11);comment:排序" json:"order"`                                            // 排序
	CreatedTime    *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateTime     *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted        *int64     `gorm:"column:deleted;type:int(11);comment:是否删除" json:"deleted"`                                      // 是否删除
	Remark         *string    `gorm:"column:remark;type:varchar(255);comment:备注 描述" json:"remark"`                                  // 备注 描述
	CreatedBy      *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                             // 创建人
	UpdateBy       *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`
}

// TableName ScaAuthPermission's table name
func (*ScaAuthPermission) TableName() string {
	return TableNameScaAuthPermission
}

func (permission *ScaAuthPermission) MarshalBinary() ([]byte, error) {
	return json.Marshal(permission)
}

func (permission *ScaAuthPermission) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, permission)
}
func (permission *ScaAuthPermission) BeforeCreate(tx *gorm.DB) (err error) {
	userId, b := global.DB.Get("user_id")
	if !b {
		global.LOG.Error("user_id is not in global.DB")
		return fmt.Errorf("user_id is not in global.DB")
	}

	userIdStr, ok := userId.(*string)
	if !ok {
		global.LOG.Error("user_id is not of type *string")
		return fmt.Errorf("user_id is not of type *string")
	}

	permission.CreatedBy = userIdStr
	return nil
}

func (permission *ScaAuthPermission) BeforeUpdate(tx *gorm.DB) (err error) {
	userId, b := global.DB.Get("user_id")
	if !b {
		global.LOG.Error("user_id is not in global.DB")
		return fmt.Errorf("user_id is not in global.DB")
	}

	userIdStr, ok := userId.(*string)
	if !ok {
		global.LOG.Error("user_id is not of type *string")
		return fmt.Errorf("user_id is not of type *string")
	}

	permission.UpdateBy = userIdStr
	return nil
}
