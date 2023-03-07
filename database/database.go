package database

import (
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	user     = ""
	password = ""
	url      = ""
	port     = 0
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/DBaaS", user, password, url, port)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return &Database{db}
}

func (d *Database) AddUpload(upload *models.Upload) error {
	err := d.db.Create(upload).Error
	return err
}

func (d *Database) GetUpload(id int) (*models.Upload, error) {
	upload := &models.Upload{Id: id}
	err := d.db.First(upload).Error
	return upload, err
}

func (d *Database) AddJob(job *models.Job) error {
	err := d.db.Create(job).Error
	return err
}

func (d *Database) GetJob() (*models.Job, error) {
	job := &models.Job{}
	err := d.db.First(job, "job_status = ?", "none").Error
	return job, err
}

func (d *Database) UpdateJob(id int) error {
	job := &models.Job{Id: id}
	err := d.db.Model(job).Update("job_status", "executed").Error
	return err
}

func (d *Database) AddResult(result *models.Result) error {
	err := d.db.Create(result).Error
	return err
}

func (d *Database) UpdateResult(id int, output, executeStatus string) error {
	result := &models.Result{Id: id, Output: output, ExecuteStatus: executeStatus}
	err := d.db.Model(result).Updates(result).Error
	return err
}
