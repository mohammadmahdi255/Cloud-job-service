package models

import (
	"gorm.io/datatypes"
)

type Result struct {
	Id            int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Job           int `gorm:"references:id"`
	Output        string
	ExecuteStatus string
	ExecuteDate   datatypes.Date
}

func NewResult(jobId int) *Result {
	result := &Result{}

	result.Job = jobId
	result.ExecuteStatus = "in progress"

	return result
}
