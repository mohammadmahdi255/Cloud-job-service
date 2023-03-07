package services

import (
	"bytes"
	"fmt"
	"github.com/mohammadmahdi255/Cloud-job-service/database"
	"github.com/mohammadmahdi255/Cloud-job-service/database/models"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	apiUrl = "https://api.codex.jaagrav.in"
)

func ExecuteJob() {
	d := database.NewDatabase()

	for {

		job, err := d.GetJob()

		if err != nil {
			if err.Error() == "record not found" {
				time.Sleep(5 * time.Second)
				continue
			} else {
				log.Println(err)
				return
			}
		}

		err = d.UpdateJob(job.Id)
		if err != nil {
			panic(err)
		}
		result := models.NewResult(job.Id)
		err = d.AddResult(result)
		if err != nil {
			panic(err)
		}
		output := executeRequest(job)
		log.Println(output)

		err = d.UpdateResult(result.Id, output, "done")
		if err != nil {
			panic(err)
		}

	}
}

func executeRequest(job *models.Job) string {
	payload := bytes.NewReader([]byte(job.JobQuery))
	client := http.Client{}

	// create http request
	req, _ := http.NewRequest(http.MethodPost, apiUrl, payload)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	// check the request
	reqDump, _ := httputil.DumpRequest(req, true)
	log.Printf("\nrequest:\n%s\n\n", reqDump)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// check the response
	respDump, _ := httputil.DumpResponse(resp, true)
	log.Printf("\nresponse:\n%s\n\n", respDump)

	// create response struct
	type Response struct {
		TimeStamp int    `json:"timeStamp"`
		Status    int    `json:"status"`
		Output    string `json:"output"`
		Error     string `json:"error"`
		Language  string `json:"language"`
		Info      string `json:"info"`
	}

	bytesString, _ := io.ReadAll(resp.Body)
	r := objx.MustFromJSON(string(bytesString))

	fmt.Printf("\noutput:\n")

	if r.Get("error").String() == "" {
		return r.Get("output").String()
	} else {
		return r.Get("error").String()
	}
}
