package main

import (
	"encoding/json"
	"github.com/paytonrules/photolibrary"
	"github.com/paytonrules/thumbnailRequest"
	"net/http"
)

type GenerateThumbnailsCommand struct {
	Events photolibrary.Events
}

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectory(directory string) {
  event, _ := c.Events.Find(directory)

  for _, img := range event.Images {
    img.GenerateThumbnail()
  }

  for _, childEvent := range event.Events {
    c.generateThumbnailsForDirectory(childEvent.FullName)
  }
}

func (c *GenerateThumbnailsCommand) Execute(r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request thumbnailRequest.Request
	decoder.Decode(&request)

  c.generateThumbnailsForDirectory(request.Directory)
}
// Decode JSON
// Loop through every image
  // - generate image
  // - check clock, see if it's time to stop
  // recurse for each directory
  // those could be go routines

func GenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	obj := GenerateThumbnailsCommand{Events: photolibrary.FileSystemEvents{}}

	obj.Execute(r)
}
