package model

import (
	"fmt"
	"gorm.io/gorm"
	"schisandra-cloud-album/global"
	"time"
)

const TableNameScaAuthUserDevice = "sca_auth_user_device"

// ScaAuthUserDevice 用户设备信息
type ScaAuthUserDevice struct {
	ID              int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                  // 主键ID
	UserID          *string    `gorm:"column:user_id;type:varchar(20);comment:用户ID" json:"user_id"`                                  // 用户ID
	IP              *string    `gorm:"column:ip;type:varchar(20);comment:登录IP" json:"ip"`                                            // 登录IP
	Location        *string    `gorm:"column:location;type:varchar(20);comment:地址" json:"location"`                                  // 地址
	Agent           string     `gorm:"column:agent;type:varchar(255);comment:设备信息" json:"agent"`                                     // 设备信息
	CreatedTime     *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateTime      *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted         *int64     `gorm:"column:deleted;type:int(11);default:0;comment:是否删除" json:"deleted"`                            // 是否删除
	Browser         *string    `gorm:"column:browser;type:varchar(20);comment:浏览器" json:"browser"`                                   // 浏览器
	OperatingSystem *string    `gorm:"column:operating_system;type:varchar(20);comment:操作系统" json:"operating_system"`                // 操作系统
	BrowserVersion  *string    `gorm:"column:browser_version;type:varchar(20);comment:浏览器版本" json:"browser_version"`                 // 浏览器版本
	Mobile          *bool      `gorm:"column:mobile;type:int(11);comment:是否为手机" json:"mobile"`                                       // 是否为手机
	Bot             *bool      `gorm:"column:bot;type:int(11);comment:是否为机器人" json:"bot"`                                            // 是否为机器人
	Mozilla         *string    `gorm:"column:mozilla;type:varchar(10);comment:火狐版本" json:"mozilla"`                                  // 火狐版本
	Platform        *string    `gorm:"column:platform;type:varchar(20);comment:平台" json:"platform"`                                  // 平台
	EngineName      *string    `gorm:"column:engine_name;type:varchar(20);comment:引擎名称" json:"engine_name"`                          // 引擎名称
	EngineVersion   *string    `gorm:"column:engine_version;type:varchar(20);comment:引擎版本" json:"engine_version"`                    // 引擎版本
	CreatedBy       *string    `gorm:"column:created_by;type:varchar(32);default:system;comment:创建人" json:"created_by"`              // 创建人
	UpdateBy        *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`
}

// TableName ScaAuthUserDevice's table name
func (*ScaAuthUserDevice) TableName() string {
	return TableNameScaAuthUserDevice
}

func (device *ScaAuthUserDevice) BeforeUpdate(tx *gorm.DB) (err error) {
	userId, b := global.DB.Get("user_id")
	if !b {
		global.LOG.Error("user_id is not found in global.DB")
		return fmt.Errorf("user_id is not found in global.DB")
	}
	userIdStr, ok := userId.(*string)
	if !ok {
		global.LOG.Error("user_id is not of type *string")
		return fmt.Errorf("user_id is not of type *string")
	}

	device.UpdateBy = userIdStr
	return nil
}
