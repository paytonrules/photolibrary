package jobs

import (
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"os"
	"photolibrary/app/models"
)

type GenerateThumbnails struct {
	Images []models.Image
}

func (job GenerateThumbnails) Run() {
	for _, image := range job.Images {
		_, err := os.Lstat(image.Thumbnail)
		if os.IsNotExist(err) {
      jobs.Now(GenerateThumbnail{Image: image})
	  }
  }
}
