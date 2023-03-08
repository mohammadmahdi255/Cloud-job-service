package services

import (
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"github.com/mohammadmahdi255/Cloud-job-service/rabbitmq"
	"github.com/mohammadmahdi255/Cloud-job-service/storage"
	"github.com/stretchr/objx"
	"log"
	"net/url"
	"sync"
)

func JobMaker(mutex *sync.RWMutex, bucket, endpoint, URL string) {

	d := database.NewDatabase()
	s := storage.NewStorage("default", endpoint)
	c := rabbitmq.NewConsumer(URL)

	message, err := c.GetMessage()
	if err != nil {
		panic(err)
	}

	mutex.Unlock()

	for data := range message {
		log.Printf("Received Data: %s", string(data.Body))
		dic := objx.MustFromJSON(string(data.Body))

		upload, err := d.GetUpload(dic.Get("id").Int())
		if err != nil {
			log.Println(err)
			return
		}

		//todo: must take object from bucket and replace it with content
		filename := fmt.Sprintf("%d.%s", upload.Id, upload.ProgramLanguage)
		buffer, err := s.Download(bucket, filename)
		if err != nil {
			log.Println(err)
			return
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
			log.Println(err)
			return
		}

	}

}
