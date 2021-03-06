package library

import (
	"os"
	"path/filepath"
	"strings"
)

type FileSystemEvents struct {
	supportedFiles map[string]bool
}

func MakeFileSystemEvents(SupportedFiles []string) FileSystemEvents {
	events := FileSystemEvents{supportedFiles: map[string]bool{}}
	for _, supportedFileExt := range SupportedFiles {
		events.supportedFiles[strings.ToLower(supportedFileExt)] = true
	}
	return events
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

func (evts FileSystemEvents) filterOutUnsupportedFiles(glob []string) []string {
	filenames := []string{}
	for _, filename := range glob {
		if evts.supportedFiles[strings.ToLower(filepath.Ext(filename))] {
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
	fileNames = events.filterOutUnsupportedFiles(fileNames)

	images := make([]Image, 0, len(fileNames))
	for _, file := range fileNames {
		images = append(images, &FileSystemImage{FullPath: file,
			Thumbnail: filepath.Dir(file) + "/.thumbnails/" + filepath.Base(file)})
	}

	eventDescriptions := make([]EventDescription, 0, len(directories))
	for _, directory := range directories {
		eventDescriptions = append(eventDescriptions, EventDescription{
			FullName:  directory,
			ShortName: filepath.Base(directory)})
	}

	return Event{Images: images, Events: eventDescriptions}, nil
}
