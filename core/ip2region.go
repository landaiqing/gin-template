package core

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"schisandra-cloud-album/global"
)

func InitIP2Region() {
	var dbPath = "ip2region/ip2region.xdb"
	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		global.LOG.Errorf("failed to load vector index from `%s`: %s\n", dbPath, err)
		return
	}
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		global.LOG.Errorf("failed to create searcher with vector index: %s\n", err)
		return
	}
	global.IP2Location = searcher
	return
}
