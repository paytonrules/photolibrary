package main

import (
	"encoding/json"
	"fmt"
	"github.com/paytonrules/photolibrary/library"
	"github.com/paytonrules/photolibrary/thumbnailrequest"
	"net/http"
	"time"
)

type GenerateThumbnailsCommand struct {
	Events    library.Events
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

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectoryAndDuration(directory string, duration int) {
	c.startTime = time.Now()
	c.duration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))
	c.generateThumbnailsForDirectory(directory)
}

func (c *GenerateThumbnailsCommand) Execute(r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request thumbnailRequest.Request
	decoder.Decode(&request)

	c.generateThumbnailsForDirectoryAndDuration(request.Directory, request.Duration)
}
