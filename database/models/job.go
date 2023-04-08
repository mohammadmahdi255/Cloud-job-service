package models

type Job struct {
	Id        int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Upload    int `gorm:"references:id"`
	JobQuery  string
	JobStatus string
}

func NewJob(uploadId int, query string) *Job {
	job := &Job{}

	job.Upload = uploadId
	job.JobQuery = query
	job.JobStatus = "suspend"

	return job
}
