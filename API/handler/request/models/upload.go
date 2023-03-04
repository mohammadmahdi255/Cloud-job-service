package models

import (
	"encoding/json"
	"log"
)

type Upload struct {
	Email           string `json:"email"`
	Inputs          string `json:"inputs"`
	ProgramLanguage string `json:"programLanguage"`
	IsEnable        bool   `json:"isEnable"`
}

func NewUpload(j []string) *Upload {
	upload := &Upload{}

	err := json.Unmarshal([]byte(j[0]), upload)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return upload
}
