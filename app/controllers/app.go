package controllers

import (
	"github.com/paytonrules/photoLibrary/app/models"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

// Show a place holder for images that don't have thumbnails yet
// Show an error for invalid files (I don't know how to do this thingy)
// Support .thm and movies
func (c App) renderDirectory(directory string) revel.Result {
  image_url, found := revel.Config.String("image_url")
  if !found {
    return c.RenderText("Could not find image url")
  }
	events := models.FileSystemEvents{}
	event, err := events.Find(directory)
	revel.INFO.Println(event)

	if err != nil {
		return c.RenderError(err)
	} else {
    c.RenderArgs["root_url"] = image_url
		c.RenderArgs["images"] = event.Images
		c.RenderArgs["directories"] = event.Events
		return c.RenderTemplate("App/Show.html")
	}
}

func (c App) Show(directory string) revel.Result {
	return c.renderDirectory(directory)
}

func (c App) Index() revel.Result {
  rootDir, found := revel.Config.String("root_dir")
  if found {
    return c.renderDirectory(rootDir)
  } else {
    return c.RenderText("Could not find root directory")
  }

}
