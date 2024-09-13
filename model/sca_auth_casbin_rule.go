package model

const ScaAuthCasbinRuleTableName = "sca_auth_casbin_rule"

// ScaAuthCasbinRule 角色/权限/用户关系表
type ScaAuthCasbinRule struct {
	Id    uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Ptype string `gorm:"column:ptype;type:varchar(100)" json:"ptype"`
	V0    string `gorm:"column:v0;type:varchar(100)" json:"v0"`
	V1    string `gorm:"column:v1;type:varchar(100)" json:"v1"`
	V2    string `gorm:"column:v2;type:varchar(100)" json:"v2"`
	V3    string `gorm:"column:v3;type:varchar(100)" json:"v3"`
	V4    string `gorm:"column:v4;type:varchar(100)" json:"v4"`
	V5    string `gorm:"column:v5;type:varchar(100)" json:"v5"`
}

func (m *ScaAuthCasbinRule) TableName() string {
	return ScaAuthCasbinRuleTableName
}
