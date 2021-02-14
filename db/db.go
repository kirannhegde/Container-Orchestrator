package db

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DBConn *gorm.db
)

type clusternode struct {
	id           int
	nodeName     string
	nodeIpaddr   string
	nodeCapacity int64
}

type container struct {
	containerRegistry string
	imageName         string
	imageVer          string
}
