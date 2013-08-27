package models

import (
	"github.com/paytonrules/image"
	"os"
)

type Event struct {
	Images []image.Image
	Events []string
}

func (e Event) ReplaceMissingThumbnailsWithTemp() (newEvent *Event) {
	imagesWithoutThumbnail := make([]image.Image, 0, len(e.Images))
	for _, oldImage := range e.Images {
		_, err := os.Lstat(oldImage.Thumbnail)

		var newImage image.Image
		if os.IsNotExist(err) {
			newImage = image.Image{Thumbnail: "thumbnail_being_generated.jpg",
				FullPath: oldImage.FullPath}
		} else {
			newImage = image.Image{Thumbnail: oldImage.Thumbnail,
				FullPath: oldImage.FullPath}
		}

		imagesWithoutThumbnail = append(imagesWithoutThumbnail, newImage)
	}

	return &Event{Images: imagesWithoutThumbnail, Events: e.Events}
}
