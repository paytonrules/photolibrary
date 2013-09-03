package photolibrary

import (
	"os"
	"path/filepath"
)

type EventDescription struct {
	FullName  string
	ShortName string
}

type Event struct {
	Images []Image
	Events []EventDescription
}

type mappingFunction func(Image) Image

func (e Event) mapImages(f mappingFunction) []Image {
	images := make([]Image, 0, len(e.Images))

	for _, originalImage := range e.Images {
		images = append(images, f(originalImage))
	}

	return images
}

func (e Event) ReplaceMissingThumbnailsWithTemp() (newEvent *Event) {
	imagesWithTempThumbnail := e.mapImages(func(original Image) Image {
		_, err := os.Lstat(original.Thumbnail)

		var newImage Image
		if os.IsNotExist(err) {
			newImage = Image{Thumbnail: "thumbnail_being_generated.jpg",
				FullPath: original.FullPath}
		} else {
			newImage = Image{Thumbnail: original.Thumbnail,
				FullPath: original.FullPath}
		}
		return newImage
	})

	return &Event{Images: imagesWithTempThumbnail, Events: e.Events}
}

func (e Event) ReplaceRelativePathsWithFullPaths() (newEvent *Event) {
	imagesWithFullPaths := e.mapImages(func(original Image) Image {
		absPath, _ := filepath.Abs(original.FullPath)
		absPathThumbnail, _ := filepath.Abs(original.Thumbnail)
		return Image{FullPath: absPath, Thumbnail: absPathThumbnail}
	})

	events := []EventDescription{}
	for _, event := range e.Events {
		fullDirectory, _ := filepath.Abs(event.FullName)
		absEvent := EventDescription{FullName: fullDirectory}
		events = append(events, absEvent)
	}

	return &Event{Images: imagesWithFullPaths, Events: events}
}
