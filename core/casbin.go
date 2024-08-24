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
		adapter, err := gormadapter.NewAdapterByDBUseTableName(global.DB, global.CONFIG.Casbin.TablePrefix, global.CONFIG.Casbin.TableName)
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		m, err := model.NewModelFromFile(global.CONFIG.Casbin.ModelPath)
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		e, err := casbin.NewCachedEnforcer(m, adapter)
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		e.EnableCache(true)
		e.SetExpireTime(60 * 60)
		err = e.LoadPolicy()
		if err != nil {
			global.LOG.Error(err.Error())
			panic(err)
		}
		global.Casbin = e
	})
}
