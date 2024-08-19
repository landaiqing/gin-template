package model

import (
	"time"
)

const TableNameScaAuthUserDevice = "sca_auth_user_device"

// ScaAuthUserDevice 用户设备信息
type ScaAuthUserDevice struct {
	ID              int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	UserID          *int64     `gorm:"column:user_id;type:bigint(20);comment:用户ID" json:"user_id"`                                   // 用户ID
	IP              *string    `gorm:"column:ip;type:varchar(255);comment:登录IP" json:"ip"`                                           // 登录IP
	Location        *string    `gorm:"column:location;type:varchar(255);comment:地址" json:"location"`                                 // 地址
	Agent           *string    `gorm:"column:agent;type:varchar(255);comment:设备信息" json:"agent"`                                     // 设备信息
	ExtJSON         *string    `gorm:"column:ext_json;type:varchar(255);comment:额外字段" json:"ext_json"`                               // 额外字段
	CreatedBy       *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                             // 创建人
	CreatedTime     *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateBy        *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`                               // 更新人
	UpdateTime      *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted         *int64     `gorm:"column:deleted;type:int(11);default:0;comment:是否删除" json:"deleted"`                            // 是否删除
	Browser         *string    `gorm:"column:browser;type:varchar(255);comment:浏览器" json:"browser"`                                  // 浏览器
	OperatingSystem *string    `gorm:"column:operating_system;type:varchar(255);comment:操作系统" json:"operating_system"`               // 操作系统
	BrowserVersion  *string    `gorm:"column:browser_version;type:varchar(255);comment:浏览器版本" json:"browser_version"`                // 浏览器版本
}

// TableName ScaAuthUserDevice's table name
func (*ScaAuthUserDevice) TableName() string {
	return TableNameScaAuthUserDevice
}
