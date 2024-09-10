package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"schisandra-cloud-album/global"
	"time"
)

const TableNameScaAuthUser = "sca_auth_user"

// ScaAuthUser 用户表
type ScaAuthUser struct {
	ID          int64      `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true;comment:自增ID" json:"-"`                // 自增ID
	UID         *string    `gorm:"column:uid;type:varchar(20);comment:唯一ID" json:"uid"`                                          // 唯一ID
	Username    *string    `gorm:"column:username;type:varchar(32);comment:用户名" json:"username"`                                 // 用户名
	Nickname    *string    `gorm:"column:nickname;type:varchar(32);comment:昵称" json:"nickname"`                                  // 昵称
	Email       *string    `gorm:"column:email;type:varchar(32);comment:邮箱" json:"email"`                                        // 邮箱
	Phone       *string    `gorm:"column:phone;type:varchar(32);comment:电话" json:"phone"`                                        // 电话
	Password    *string    `gorm:"column:password;type:varchar(64);comment:密码" json:"-"`                                         // 密码
	Gender      *string    `gorm:"column:gender;type:varchar(32);comment:性别" json:"gender"`                                      // 性别
	Avatar      *string    `gorm:"column:avatar;type:longtext;comment:头像" json:"avatar"`                                         // 头像
	Status      *int64     `gorm:"column:status;type:tinyint(4);default:0;comment:状态 0 正常 1 封禁" json:"status"`                   // 状态 0 正常 1 封禁
	Introduce   *string    `gorm:"column:introduce;type:varchar(255);comment:介绍" json:"introduce"`                               // 介绍
	CreatedTime *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdateTime  *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted     *int64     `gorm:"column:deleted;type:int(11);default:0;comment:是否删除 0 未删除 1 已删除" json:"-"`                      // 是否删除 0 未删除 1 已删除
	Blog        *string    `gorm:"column:blog;type:varchar(30);comment:博客" json:"blog"`                                          // 博客
	Location    *string    `gorm:"column:location;type:varchar(50);comment:地址" json:"location"`                                  // 地址
	Company     *string    `gorm:"column:company;type:varchar(50);comment:公司" json:"company"`                                    // 公司
	CreatedBy   *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                             // 创建人
	UpdateBy    *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`                               // 更新人
}

// TableName ScaAuthUser's table name
func (*ScaAuthUser) TableName() string {
	return TableNameScaAuthUser
}

func (user *ScaAuthUser) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (user *ScaAuthUser) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}
func (user *ScaAuthUser) BeforeCreate(tx *gorm.DB) (err error) {
	userId, b := global.DB.Get("user_id")
	if !b {
		creator := "system"
		user.CreatedBy = &creator
		return nil
	}

	userIdStr, ok := userId.(*string)
	if !ok {
		global.LOG.Error("user_id is not of type *string")
		return fmt.Errorf("user_id is not of type *string")
	}

	user.CreatedBy = userIdStr
	return nil
}

func (user *ScaAuthUser) BeforeUpdate(tx *gorm.DB) (err error) {
	userId, b := global.DB.Get("user_id")
	if !b {
		creator := "system"
		user.CreatedBy = &creator
		return nil
	}

	userIdStr, ok := userId.(*string)
	if !ok {
		global.LOG.Error("user_id is not of type *string")
		return fmt.Errorf("user_id is not of type *string")
	}

	user.UpdateBy = userIdStr
	return nil
}
