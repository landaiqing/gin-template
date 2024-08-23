package core

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"schisandra-cloud-album/global"
	"sync"
)

var (
	once sync.Once
)

func InitCasbin() {
	once.Do(func() {
		adapter, _ := gormadapter.NewAdapterByDBUseTableName(global.DB, "sca_sys_", "casbin_rule")
		m, err := model.NewModelFromFile("config/rbac_model.pml")
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		e, _ := casbin.NewCachedEnforcer(m, adapter)
		e.SetExpireTime(60 * 60)
		err = e.LoadPolicy()
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		global.Casbin = e
	})
}
