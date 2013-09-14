package controllers

import (
  "github.com/paytonrules/photolibrary/library"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (c App) renderDirectory(directory string) revel.Result {
  image_url, found := revel.Config.String("image_url")
  if !found {
    return c.RenderText("Could not find image url")
  }
	events := library.MakeFileSystemEvents([]string{".jpg", ".png", ".mov", ".avi"})
	event, err := events.Find(directory)

  rootDir, _ := revel.Config.String("root_dir")
  eventWithThumbnails := event.ReplaceMissingThumbnailsWithTemp()
  relativeEvents := eventWithThumbnails.ImagesRelativeTo(rootDir)

	if err != nil {
		return c.RenderError(err)
	} else {
    c.RenderArgs["root_url"] = image_url
		c.RenderArgs["images"] = relativeEvents.Images
		c.RenderArgs["directories"] = relativeEvents.Events
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
