package models

type Upload struct {
	Id              int    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Email           string `json:"email"`
	Inputs          string `json:"inputs"`
	ProgramLanguage string `json:"programLanguage"`
	IsEnable        bool   `json:"isEnable"`
}

func NewUpload(email, inputs, programLanguage string) (*Upload, error) {
	upload := &Upload{}

	upload.Email = email
	upload.Inputs = inputs
	upload.ProgramLanguage = programLanguage
	upload.IsEnable = true

	return upload, nil
}
