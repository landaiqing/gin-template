package core

import (
	sensitive "github.com/zmexing/go-sensitive-word"
	"log"
	"schisandra-cloud-album/global"
)

func InitSensitive() {
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		log.Fatalf("init sensitive filter failed, err:%v", err)
		return
	}
	// 加载敏感词库
	err = filter.Store.LoadDictPath("resource/sensitive/反动词库.txt",
		"resource/sensitive/暴恐词库.txt",
		"resource/sensitive/色情词库.txt",
		"resource/sensitive/贪腐词库.txt",
		"resource/sensitive/民生词库.txt")
	if err != nil {
		log.Fatalf("load sensitive dict failed, err:%v", err)
		return
	}
	global.SensitiveManager = filter
}
