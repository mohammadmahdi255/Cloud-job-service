package database

import (
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/DBaaS?parseTime=true", global.User, global.Password, global.Url, global.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully connected to Database")
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

func (d *Database) UpdateUpload(id int, isEnable bool) error {
	upload := &models.Upload{Id: id}
	err := d.db.Model(upload).Update("is_enable", isEnable).Error
	return err
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
	fmt.Println(result)
	return err
}

func (d *Database) UpdateResult(id int, output, executeStatus string) error {
	result := &models.Result{Id: id, Output: output, ExecuteStatus: executeStatus}
	err := d.db.Model(result).Updates(result).Error
	return err
}

func (d *Database) GetAllUserResult(email string) ([]map[string]interface{}, error) {

	rows, err := d.db.Table("uploads").Select("uploads.id, email, program_language, execute_status, execute_date").
		Joins("join jobs j join results r on uploads.id = j.upload and j.id = r.job").
		Where("email = ?", email).Rows()

	if err != nil {
		return nil, err
	}

	arr := make([]map[string]interface{}, 0, 10)

	for rows.Next() {
		dic := map[string]interface{}{}
		err := d.db.ScanRows(rows, dic)
		if err != nil {
			return nil, err
		}

		dic["file"] = fmt.Sprintf("%s/%s/%d.%s", global.Endpoint, global.Bucket, dic["id"], dic["program_language"])
		arr = append(arr, dic)
	}

	err = rows.Close()
	return arr, err
}

func (d *Database) Delete(data interface{}) {
	d.db.Delete(data)
}
