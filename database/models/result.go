package models

import (
	"time"
)

type Result struct {
	Id            int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Job           int `gorm:"references:id"`
	Output        string
	ExecuteStatus string
	ExecuteDate   time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}

func NewResult(jobId int) *Result {
	result := &Result{}

	result.Job = jobId
	result.ExecuteStatus = "in progress"

	return result
}
