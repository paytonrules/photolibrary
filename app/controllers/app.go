package controllers

import (
	//	"github.com/nfnt/resize"
	"github.com/robfig/revel"
	//  "image/jpeg"
	//  "image/png"
	"os"
	"path/filepath"
)

type Directory struct {
	FullName  string
	ShortName string
}

type Image struct {
	Thumbnail string
	FullPath  string
}

type App struct {
	*revel.Controller
}

func thumbnailPathFor(imagePath string) string {
  return filepath.Dir(imagePath) + "/thumbnails/" + filepath.Base(imagePath)
}

func makeImage(imagePath string) Image {
  return Image{Thumbnail: thumbnailPathFor(imagePath), FullPath: imagePath}
}

func (c App) renderDirectory(directory string) revel.Result {
	allFiles, error := filepath.Glob(directory + "/*")
	if error != nil {
		return c.Render(error)
	}

	directories := make([]Directory, 0, len(allFiles))
	files := make([]Image, 0, len(allFiles))

	for i := range allFiles {
		lstat, error := os.Lstat(allFiles[i])
		if error != nil {
			return c.Render(error)
		}

		if lstat.IsDir() {
			directories = append(directories, Directory{FullName: allFiles[i], ShortName: filepath.Base(allFiles[i])})
		} else {
			files = append(files, makeImage(allFiles[i]))
		}
	}
	c.RenderArgs["files"] = files
	c.RenderArgs["directories"] = directories
	return c.RenderTemplate("App/Show.html")
}

func (c App) Show(directory string) revel.Result {
	return c.renderDirectory(directory)
}

func (c App) Index() revel.Result {
	// Convert to a url
	return c.renderDirectory("public/events")
}
