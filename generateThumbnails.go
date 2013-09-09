package main

import (
	"encoding/json"
	"github.com/paytonrules/photolibrary"
	"github.com/paytonrules/thumbnailRequest"
	"net/http"
  "fmt"
	"time"
)

type GenerateThumbnailsCommand struct {
	Events    photolibrary.Events
	startTime time.Time
	duration  time.Duration
}

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectory(directory string) {
	event, _ := c.Events.Find(directory)

  for _, img := range event.Images {
    if time.Since(c.startTime) < c.duration {
      img.GenerateThumbnail()
    }
  }

	for _, childEvent := range event.Events {
    if time.Since(c.startTime) < c.duration {
      c.generateThumbnailsForDirectory(childEvent.FullName)
    }
	}
}

func (c *GenerateThumbnailsCommand) Execute(r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request thumbnailRequest.Request
	decoder.Decode(&request)
	c.startTime = time.Now()
  c.duration, _ = time.ParseDuration(fmt.Sprintf("%ds", request.Duration))

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
