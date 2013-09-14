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
		_, err := os.Lstat(original.GetThumbnail())

		var newImage Image
		if os.IsNotExist(err) {
			newImage = &FileSystemImage{Thumbnail: "thumbnail_being_generated.png",
				FullPath: original.GetFullPath()}
		} else {
			newImage = original.Clone()
		}
		return newImage
	})

	return &Event{Images: imagesWithTempThumbnail, Events: e.Events}
}

func (e Event) ReplaceRelativePathsWithFullPaths() *Event {
	imagesWithFullPaths := e.mapImages(func(original Image) Image {
		absPath, _ := filepath.Abs(original.GetFullPath())
		absPathThumbnail, _ := filepath.Abs(original.GetThumbnail())
		return &FileSystemImage{FullPath: absPath, Thumbnail: absPathThumbnail}
	})

	events := []EventDescription{}
	for _, event := range e.Events {
		fullDirectory, _ := filepath.Abs(event.FullName)
		absEvent := EventDescription{FullName: fullDirectory}
		events = append(events, absEvent)
	}

	return &Event{Images: imagesWithFullPaths, Events: events}
}

func (e Event) ImagesRelativeTo(root string) *Event {
  imagesWithRelativePaths := e.mapImages(func(original Image) Image {
		relPath, _ := filepath.Rel(root, original.GetFullPath())
		relPathThumbnail, _ := filepath.Rel(root, original.GetThumbnail())
		return &FileSystemImage{FullPath: relPath, Thumbnail: relPathThumbnail}
  });

  return &Event{Images: imagesWithRelativePaths, Events: e.Events}
}
