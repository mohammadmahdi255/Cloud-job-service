package main

import (
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/stretchr/objx"
	"net/url"
)

func main() {

	for {

		d := database.NewDatabase()
		s := storage.NewStorage("default", global.Endpoint)
		c := rabbitmq.NewConsumer(global.CloudamqpUrl)

		message, err := c.GetMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for data := range message {
			fmt.Printf("Received Data: %s\n", string(data.Body))
			dic := objx.MustFromJSON(string(data.Body))

			upload, err := d.GetUpload(dic.Get("id").Int())
			if err != nil {
				fmt.Println(err)
				break
			}

			//todo: must take object from bucket and replace it with content
			filename := fmt.Sprintf("%d.%s", upload.Id, upload.ProgramLanguage)
			buffer, err := s.Download(global.Bucket, filename)
			if err != nil {
				fmt.Println(err)
				break
			}
			content := string(buffer.Bytes())

			queryParams := url.Values{
				"code":     {content},
				"language": {upload.ProgramLanguage},
				"input":    {upload.Inputs},
			}

			job := models.NewJob(upload.Id, queryParams.Encode())

			err = d.AddJob(job)
			if err != nil {
				fmt.Println(err)
				break
			}

		}
	}

}
