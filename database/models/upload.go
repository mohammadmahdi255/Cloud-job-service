package models

import (
	"github.com/stretchr/objx"
)

type Upload struct {
	Id              int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Email           string `json:"email"`
	Inputs          string `json:"inputs"`
	ProgramLanguage string `json:"programLanguage"`
	IsEnable        bool   `json:"isEnable"`
}

func NewUpload(j string) (*Upload, error) {
	upload := &Upload{}

	dic, err := objx.FromJSON(j)
	if err != nil {
		return nil, err
	}

	upload.Inputs = dic.Get("inputs").String()
	upload.ProgramLanguage = dic.Get("programLanguage").String()
	upload.IsEnable = dic.Get("isEnable").Bool()

	return upload, nil
}
