package photolibrary

import (
	"github.com/nfnt/resize"
	"log"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemImage struct {
	Thumbnail string
	FullPath  string
}
// Make work with PNG
// Clean up error handling
func (image FileSystemImage) GenerateThumbnail() {
	_, err := os.Lstat(image.Thumbnail)
  log.Println("generating " + image.FullPath)
	if err != nil && os.IsNotExist(err) {
		if strings.Contains(strings.ToUpper(filepath.Ext(image.FullPath)), ".JPG") {
      // touch file at very beginning to reduce duplicate jobs
      // Still a race condition, but as long as there are no errors I can live with it
      // Open the full image file
      file, err := os.Open(image.FullPath)
      if err != nil {
        log.Println(err)
      }

      // decode jpeg into image.Image
      img, err := jpeg.Decode(file)
      if err != nil {
        log.Println(err)
      }
      file.Close()

      // See if there is a thumbnails directory
      _, err = os.Lstat(filepath.Dir(image.Thumbnail))
      if err != nil {
        if os.IsNotExist(err) {
          os.Mkdir(filepath.Dir(image.Thumbnail), os.ModeDir|os.ModePerm)
        } else {
          log.Println(err)
        }
      }

      // Probably shouldn't continue in the event of error but.....
      // resize to width 200 using NearestNeighbor resampling
      // and preserve aspect ratio
      m := resize.Resize(200, 0, img, resize.NearestNeighbor)

      out, err := os.Create(image.Thumbnail)
      if err != nil {
        log.Println(err)
      }
      defer out.Close()

      // write new image to file
      jpeg.Encode(out, m, nil)
		}
	}
}

