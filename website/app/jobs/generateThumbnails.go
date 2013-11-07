package jobs

import (
	"bytes"
	"github.com/robfig/revel"
	"encoding/json"
	"github.com/paytonrules/photolibrary/thumbnailrequest"
	"net/http"
)

type GenerateThumbnails struct {
	Server    string
	Duration  int
	Directory string
}

func (job GenerateThumbnails) Run() {
	thumbnailRequest := thumbnailrequest.Request{Directory: job.Directory, Duration: job.Duration}
	marshaledThumbnailRequest, _ := json.Marshal(thumbnailRequest)
	body := bytes.NewBuffer(marshaledThumbnailRequest)
  revel.ERROR.Printf("Log this buffer: %s", body)
	resp, _ := http.Post(job.Server+"/generateThumbnails", "text/json", body)
	defer resp.Body.Close()
}
