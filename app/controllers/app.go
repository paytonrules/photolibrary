package controllers

import (
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"os"
	"path/filepath"
	photoJobs "photolibrary/app/jobs"
	"photolibrary/app/models"
)

type Directory struct {
	FullName  string
	ShortName string
}

type App struct {
	*revel.Controller
}

func thumbnailPathFor(imagePath string) string {
	return filepath.Dir(imagePath) + "/.thumbnails/" + filepath.Base(imagePath)
}

func makeImage(imagePath string) models.Image {
	return models.Image{Thumbnail: thumbnailPathFor(imagePath), FullPath: imagePath}
}

// .thumbnails is not one of the directories
// Actually skip all hiddens
// Show a place holder for images that don't have thumbnails yet
// Show an error for invalid files (I don't know how to do this thingy)
// Support .thm and movies
func (c App) renderDirectory(directory string) revel.Result {
	allFiles, error := filepath.Glob(directory + "/*")
	if error != nil {
		return c.Render(error)
	}

	directories := make([]Directory, 0, len(allFiles))
	images := make([]models.Image, 0, len(allFiles))

	for i := range allFiles {
		lstat, error := os.Lstat(allFiles[i])
		if error != nil {
			return c.Render(error)
		}

		if lstat.IsDir() {
			directories = append(directories, Directory{FullName: allFiles[i], ShortName: filepath.Base(allFiles[i])})
		} else {
			images = append(images, makeImage(allFiles[i]))
		}
	}

	revel.INFO.Println("GenerateThumbnail")
	jobs.Now(photoJobs.GenerateThumbnail{Images: images})

	c.RenderArgs["images"] = images
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
