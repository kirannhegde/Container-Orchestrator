package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//An exported variable that represents the connection to the SQLLite DB
var (
	DBConn *gorm.DB
)
