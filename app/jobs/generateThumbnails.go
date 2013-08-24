package jobs

import (
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"os"
  "github.com/paytonrules/image"
)

type GenerateThumbnails struct {
	Images []image.Image
}

func (job GenerateThumbnails) Run() {
	for _, image := range job.Images {
		_, err := os.Lstat(image.Thumbnail)
		if os.IsNotExist(err) {
      jobs.Now(GenerateThumbnail{Image: image})
	  }
  }
}
