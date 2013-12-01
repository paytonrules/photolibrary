package main

import (
	"fmt"
	"github.com/paytonrules/photolibrary/library"
	"time"
)

type NullLogger struct {
}

func (l *NullLogger) Info(message string) {
}

type GenerateThumbnailsCommand struct {
	Events    library.Events
	startTime time.Time
	duration  time.Duration
  logger    Logger
}

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectory(directory string) {
  c.logger.Info("Generating Images in " + directory)
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
  c.logger.Info("Try generating Images in " + directory)
	c.startTime = time.Now()
	c.duration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))
	c.generateThumbnailsForDirectory(directory)
}

func MakeGenerateThumbnailCommand(events library.Events) *GenerateThumbnailsCommand {
  nullLogger := new(NullLogger)
  return MakeGenerateThumbnailCommandWithLogger(events, nullLogger)
}

func MakeGenerateThumbnailCommandWithLogger(events library.Events, logger Logger) *GenerateThumbnailsCommand {
  g := GenerateThumbnailsCommand{Events: events, logger: logger}

  return &g
}
