package models

import (
	"github.com/paytonrules/image"
	"os"
	"path/filepath"
)

type Event struct {
	Images []image.Image
	Events []string
}

type mappingFunction func(image.Image) image.Image

func (e Event) mapImages(f mappingFunction) []image.Image {
	images := make([]image.Image, 0, len(e.Images))

	for _, originalImage := range e.Images {
		images = append(images, f(originalImage))
	}

	return images
}

func (e Event) ReplaceMissingThumbnailsWithTemp() (newEvent *Event) {
	imagesWithTempThumbnail := e.mapImages(func(original image.Image) image.Image {
		_, err := os.Lstat(original.Thumbnail)

		var newImage image.Image
		if os.IsNotExist(err) {
			newImage = image.Image{Thumbnail: "thumbnail_being_generated.jpg",
				FullPath: original.FullPath}
		} else {
			newImage = image.Image{Thumbnail: original.Thumbnail,
				FullPath: original.FullPath}
		}
		return newImage
	})

	return &Event{Images: imagesWithTempThumbnail, Events: e.Events}
}

func (e Event) ReplaceRelativePathsWithFullPaths() (newEvent *Event) {
	imagesWithFullPaths := e.mapImages(func(original image.Image) image.Image {
		absPath, _ := filepath.Abs(original.FullPath)
		absPathThumbnail, _ := filepath.Abs(original.Thumbnail)
		return image.Image{FullPath: absPath, Thumbnail: absPathThumbnail}
	})

  directories := []string{}
  for _, relativeDirectory := range e.Events {
    absEvent, _ := filepath.Abs(relativeDirectory)
    directories = append(directories, absEvent)
  }

  return &Event{Images: imagesWithFullPaths, Events: directories}
}
