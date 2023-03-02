package database

import (
	_ "net/http"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

type Database struct{}

var dsn = ""
