// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameScaAuthUserSocial = "sca_auth_user_social"

// ScaAuthUserSocial 社会用户信息表
type ScaAuthUserSocial struct {
	ID               int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键ID" json:"id"`                                                 // 主键ID
	UserID           int64      `gorm:"column:user_id;type:bigint(20);not null;comment:用户ID" json:"user_id"`                                         // 用户ID
	UUID             *string    `gorm:"column:uuid;type:varchar(255);comment:第三方系统的唯一ID" json:"uuid"`                                                // 第三方系统的唯一ID
	Source           *string    `gorm:"column:source;type:varchar(255);comment:第三方用户来源" json:"source"`                                               // 第三方用户来源
	AccessToken      *string    `gorm:"column:access_token;type:varchar(255);comment:用户的授权令牌" json:"access_token"`                                   // 用户的授权令牌
	ExpireIn         *int64     `gorm:"column:expire_in;type:int(11);comment:第三方用户的授权令牌的有效期" json:"expire_in"`                                       // 第三方用户的授权令牌的有效期
	RefreshToken     *string    `gorm:"column:refresh_token;type:varchar(255);comment:刷新令牌" json:"refresh_token"`                                    // 刷新令牌
	OpenID           *string    `gorm:"column:open_id;type:varchar(255);comment:第三方用户的 open id" json:"open_id"`                                      // 第三方用户的 open id
	UID              *string    `gorm:"column:uid;type:varchar(255);comment:第三方用户的 ID" json:"uid"`                                                   // 第三方用户的 ID
	AccessCode       *string    `gorm:"column:access_code;type:varchar(255);comment:个别平台的授权信息" json:"access_code"`                                   // 个别平台的授权信息
	UnionID          *string    `gorm:"column:union_id;type:varchar(255);comment:第三方用户的 union id" json:"union_id"`                                   // 第三方用户的 union id
	Scope            *string    `gorm:"column:scope;type:varchar(255);comment:第三方用户授予的权限" json:"scope"`                                              // 第三方用户授予的权限
	TokenType        *string    `gorm:"column:token_type;type:varchar(255);comment:个别平台的授权信息" json:"token_type"`                                     // 个别平台的授权信息
	IDToken          *string    `gorm:"column:id_token;type:varchar(255);comment:id token" json:"id_token"`                                          // id token
	MacAlgorithm     *string    `gorm:"column:mac_algorithm;type:varchar(255);comment:小米平台用户的附带属性" json:"mac_algorithm"`                             // 小米平台用户的附带属性
	MacKey           *string    `gorm:"column:mac_key;type:varchar(255);comment:小米平台用户的附带属性" json:"mac_key"`                                         // 小米平台用户的附带属性
	Code             *string    `gorm:"column:code;type:varchar(255);comment:用户的授权code" json:"code"`                                                 // 用户的授权code
	OauthToken       *string    `gorm:"column:oauth_token;type:varchar(255);comment:Twitter平台用户的附带属性" json:"oauth_token"`                            // Twitter平台用户的附带属性
	OauthTokenSecret *string    `gorm:"column:oauth_token_secret;type:varchar(255);comment:Twitter平台用户的附带属性" json:"oauth_token_secret"`              // Twitter平台用户的附带属性
	Status           *string    `gorm:"column:status;type:varchar(255);comment:状态 0正常 1 封禁" json:"status"`                                           // 状态 0正常 1 封禁
	ExtJSON          *string    `gorm:"column:ext_json;type:varchar(255);comment:额外字段" json:"ext_json"`                                              // 额外字段
	CreatedBy        *string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`                                            // 创建人
	CreatedTime      *time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime;comment:创建时间" json:"created_time"` // 创建时间
	UpdateBy         *string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`                                              // 更新人
	UpdateTime       *time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间" json:"update_time"`   // 更新时间
	Deleted          *int64     `gorm:"column:deleted;type:int(11);comment:是否删除" json:"deleted"`                                                     // 是否删除
}

// TableName ScaAuthUserSocial's table name
func (*ScaAuthUserSocial) TableName() string {
	return TableNameScaAuthUserSocial
}
