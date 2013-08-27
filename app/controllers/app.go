package controllers

import (
	"github.com/paytonrules/image"
	"github.com/paytonrules/photoLibrary/app/models"
	photoJobs "github.com/paytonrules/photolibrary/app/jobs"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"os"
	"path/filepath"
)

type Directory struct {
	FullName  string
	ShortName string
}

type App struct {
	*revel.Controller
	Events models.Events
}

func thumbnailPathFor(imagePath string) string {
	return filepath.Dir(imagePath) + "/.thumbnails/" + filepath.Base(imagePath)
}

func makeImage(imagePath string) image.Image {
	return image.Image{Thumbnail: thumbnailPathFor(imagePath), FullPath: imagePath}
}

func removeHiddenFilesFrom(allFiles []string) []string {
	filesWithoutHiddenFiles := make([]string, 0, len(allFiles))
	for _, filename := range allFiles {
		if filepath.Base(filename)[0] != '.' {
			filesWithoutHiddenFiles = append(filesWithoutHiddenFiles, filename)
		}
	}

	return filesWithoutHiddenFiles
}

func findEligibleFiles(glob string) ([]string, error) {
	eligibleFiles, error := filepath.Glob(glob)
	return removeHiddenFilesFrom(eligibleFiles), error
}

// Show a place holder for images that don't have thumbnails yet
// Show an error for invalid files (I don't know how to do this thingy)
// Support .thm and movies
func (c App) renderDirectory(directory string) revel.Result {
	filesWithoutHiddenFiles, error := findEligibleFiles(directory + "/*")

	if error != nil {
		return c.RenderError(error)
	}

	directories := make([]Directory, 0, len(filesWithoutHiddenFiles))
	images := make([]image.Image, 0, len(filesWithoutHiddenFiles))

	for _, currentFile := range filesWithoutHiddenFiles {
		lstat, error := os.Lstat(currentFile)
		if error != nil {
			return c.RenderError(error)
		}

		if lstat.IsDir() {
			directories = append(directories, Directory{FullName: currentFile, ShortName: filepath.Base(currentFile)})
		} else {
			images = append(images, makeImage(currentFile))
		}
	}

	jobs.Now(photoJobs.GenerateThumbnails{Images: images})

	imagesWithTemporaryThumbnails := make([]image.Image, 0, len(images))
	for _, currentImage := range images {
		_, error := os.Lstat(currentImage.Thumbnail)

		if os.IsNotExist(error) {
			// You should use 'images url' or whatever the revel syntax is
			imagesWithTemporaryThumbnails = append(imagesWithTemporaryThumbnails,
				*image.NewImage("public/img/thumbnail_being_generated.png", currentImage.FullPath))
		} else {
			imagesWithTemporaryThumbnails = append(imagesWithTemporaryThumbnails,
				*image.CopyImage(currentImage))
		}
	}

	c.RenderArgs["images"] = imagesWithTemporaryThumbnails
	c.RenderArgs["directories"] = directories
	return c.RenderTemplate("App/Show.html")
	/*
  // It's not events - events only does the find.
     A more complicated interaction does the job launching and the thumbnail detection
	   event, err := c.events.Find(directory)

	   if err != nil {
	     return c.RenderError(err)
	   } else {
	     c.RenderArgs["images"] = event.images
	     c.RenderArgs["directories"] = event.directories
	     return c.RenderTemplate("App/Show.html")
	   }*/
}

func (c App) init() revel.Result {
	c.Events = models.FileSystemEvents{}
	return nil
}

func (c App) Show(directory string) revel.Result {
	return c.renderDirectory(directory)
}

func (c App) Index() revel.Result {
	return c.renderDirectory("public/events")
}
