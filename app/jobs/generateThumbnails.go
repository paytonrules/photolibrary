package jobs

import (
	"bytes"
	"encoding/json"
	"github.com/paytonrules/image"
	"github.com/robfig/revel"
	"io"
	"net/http"
	"path/filepath"
)

type GenerateThumbnails struct {
	Images []image.Image
  Duration int
}

func (job GenerateThumbnails) Run() {
	thumbnailsWithFullPaths := make([]image.Image, 0, len(job.Images))

	for _, currentImage := range job.Images {
		imageWithFullPaths := image.Image{}
		imageWithFullPaths.FullPath, _ = filepath.Abs(currentImage.FullPath)
		imageWithFullPaths.Thumbnail, _ = filepath.Abs(currentImage.Thumbnail)

		thumbnailsWithFullPaths = append(thumbnailsWithFullPaths, imageWithFullPaths)
	}

	var body io.Reader
	imagesJson, err := json.Marshal(thumbnailsWithFullPaths)
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	body = bytes.NewBuffer(imagesJson)
	http.Post("http://localhost:8081/generateThumbnails",
		"text/json",
		body)
}
