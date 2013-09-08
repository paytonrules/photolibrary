package photolibrary

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemImage struct {
	Thumbnail string
	FullPath  string
}

func (image FileSystemImage) GetThumbnail() string {
  return image.Thumbnail
}

func (image FileSystemImage) GetFullPath() string {
  return image.FullPath;
}

func (image FileSystemImage) Clone() Image {
  return image
}

// Make work with PNG
// Clean up error handling
func (image FileSystemImage) GenerateThumbnail() {
	_, err := os.Lstat(image.Thumbnail)
	if err != nil && os.IsNotExist(err) {
		if strings.Contains(strings.ToUpper(filepath.Ext(image.FullPath)), ".JPG") {
      // touch file at very beginning to reduce duplicate jobs
      // Still a race condition, but as long as there are no errors I can live with it
      // Open the full image file
      file, _ := os.Open(image.FullPath)
      defer file.Close()

      // decode jpeg into image.Image
      img, _ := jpeg.Decode(file)

      // See if there is a thumbnails directory
      _, err := os.Lstat(filepath.Dir(image.Thumbnail))
      if err != nil {
        if os.IsNotExist(err) {
          os.Mkdir(filepath.Dir(image.Thumbnail), os.ModeDir|os.ModePerm)
        }
      }

      // Probably shouldn't continue in the event of error but.....
      // resize to width 200 using NearestNeighbor resampling
      // and preserve aspect ratio
      m := resize.Resize(200, 0, img, resize.NearestNeighbor)

      out, _ := os.Create(image.Thumbnail)
      defer out.Close()

      // write new image to file
      jpeg.Encode(out, m, nil)
		}
	}
}

