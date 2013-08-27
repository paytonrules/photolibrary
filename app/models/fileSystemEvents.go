package models

import (
	"github.com/paytonrules/image"
	"os"
	"path/filepath"
)

type FileSystemEvents struct {
}

func filterOutHiddenFiles(glob []string) []string {
	filenames := make([]string, 0, len(glob))
	for _, filename := range glob {
		if filepath.Base(filename)[0] != '.' {
			filenames = append(filenames, filename)
		}
	}

	return filenames
}

func separateOutDirectories(glob []string) ([]string, []string) {
	filenames := make([]string, 0, len(glob))
	directories := make([]string, 0, len(glob))

	for _, filename := range glob {
		info, _ := os.Lstat(filename)
		if info.IsDir() {
			directories = append(directories, filename)
		} else {
			filenames = append(filenames, filename)
		}
	}

	return filenames, directories
}

func (events FileSystemEvents) Find(directoryName string) (Event, error) {
	fileNames, err := filepath.Glob(directoryName + "/*")
	if err != nil {
		return Event{}, err
	}
	fileNames = filterOutHiddenFiles(fileNames)
	fileNames, directories := separateOutDirectories(fileNames)

	images := make([]image.Image, 0, len(fileNames))
	for _, file := range fileNames {
		images = append(images, image.Image{FullPath: file,
			Thumbnail: filepath.Dir(file) + "/.thumbnails/" + filepath.Base(file)})
	}

	return Event{Images: images, Events: directories}, nil
}
