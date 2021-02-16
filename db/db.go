package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//An exported variable that represents the connection to the SQLLite DB
var (
	DBConn *gorm.DB
)
