package jobs

import (
  "github.com/robfig/revel"
  "photolibrary/app/models"
  "os"
  "image/jpeg"
	"path/filepath"
  "github.com/nfnt/resize"
  "strings"
)

type GenerateThumbnail struct {
  Images []models.Image
}

// Compare no case on the extension duh
// Make work with PNG/JPG
// Turn into simultaneous jobs
// Clean up error handling
func (job GenerateThumbnail) generateJPG(image models.Image) {
  revel.INFO.Println("Generating Thumbnail for " + image.FullPath)
  // touch file at very beginning to reduce duplicate jobs
  // Still a race condition, but as long as there are no errors I can live with it
  // Open the full image file
  file, err := os.Open(image.FullPath)
  if err != nil {
    revel.ERROR.Println(err)
  }

  // decode jpeg into image.Image
  img, err := jpeg.Decode(file)
  if err != nil {
    revel.ERROR.Println(err)
  }
  file.Close()

  // See if there is a thumbnails directory (hmm maybe that path shouldnt
  // be in the controller)
  _, err = os.Lstat(filepath.Dir(image.Thumbnail))
  if err != nil {
    if os.IsNotExist(err) {
      os.Mkdir(filepath.Dir(image.Thumbnail), os.ModeDir | os.ModePerm)
    } else {
      revel.ERROR.Println(err)
    }
  }

  // Probably shouldn't continue in the event of error but.....
  // resize to width 200 using Lanczos resampling
  // and preserve aspect ratio

  revel.INFO.Println("Resizing image")
  m := resize.Resize(200, 0, img, resize.Lanczos3)

  out, err := os.Create(image.Thumbnail)
  if err != nil {
    revel.ERROR.Println(err)
  }
  defer out.Close()

  // write new image to file
  jpeg.Encode(out, m, nil)
}

func (job GenerateThumbnail) Run() {
  for _, image := range job.Images {
    _, err := os.Lstat(image.Thumbnail)
    if (err != nil && os.IsNotExist(err)) {
      if strings.Contains(strings.ToUpper(filepath.Ext(image.FullPath)), ".JPG") {
        job.generateJPG(image)
      }
    }
  }
}
