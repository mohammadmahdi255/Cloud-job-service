package database

import (
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/handler/request/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "net/http"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

const user = "root"
const password = "jUlfRTMMkd9GMUaVJLUn10JB"
const url = "may.iran.liara.ir"
const port = 33761

type Database struct {
	db *gorm.DB
}

func NewDatabase() *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/DBaaS", user, password, url, port)
	fmt.Println(dsn)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return &Database{db}
}

func (d *Database) AddUpload(upload *models.Upload) error {

	err := d.db.Create(upload).Error
	return err
}
