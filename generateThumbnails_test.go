package main

import (
	"bytes"
	"encoding/json"
	"github.com/paytonrules/photolibrary"
	"github.com/paytonrules/thumbnailRequest"
	. "launchpad.net/gocheck"
	"net/http"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GenerateThumbnailsSuite struct{}

var _ = Suite(&GenerateThumbnailsSuite{})

type PhonyEvents struct {
	FindString string
	FindResult photolibrary.Event
}

type PhonyImage struct {
  FullPath string
  Thumbnail string
  Generated bool
}

func (img *PhonyImage) GenerateThumbnail() {
  img.Generated = true
}

func (img PhonyImage) GetFullPath() string {
  return img.FullPath
}

func (img PhonyImage) GetThumbnail() string {
  return img.Thumbnail
}

func (evts *PhonyEvents) Find(eventName string) (photolibrary.Event, error) {
	evts.FindString = eventName
	return evts.FindResult, nil
}

func (s *GenerateThumbnailsSuite) TestExecuteFindsTheEventsWithTheRightRoot(c *C) {
	phonyEvent := PhonyEvents{}
	command := GenerateThumbnailsCommand{Events: &phonyEvent}

	thumbnailRequest := thumbnailRequest.Request{Directory: "directory", Duration: 0}

	marshaledThumbnailRequest, _ := json.Marshal(thumbnailRequest)
	body := bytes.NewBuffer(marshaledThumbnailRequest)
	req, _ := http.NewRequest("dont", "care", body)

	command.Execute(req)

	c.Assert(phonyEvent.FindString, Equals, "directory")
}

func (s *GenerateThumbnailsSuite) TestGenerateThumbnailOnEachImage(c *C) {
  image := &PhonyImage{}
  images := []photolibrary.Image{image}
  phonyEvents := &PhonyEvents{FindResult: photolibrary.Event{Images:images}}
	command := GenerateThumbnailsCommand{Events: phonyEvents}

	thumbnailRequest := thumbnailRequest.Request{Directory: "directory", Duration: 0}
	marshaledThumbnailRequest, _ := json.Marshal(thumbnailRequest)
	body := bytes.NewBuffer(marshaledThumbnailRequest)
	req, _ := http.NewRequest("dont", "care", body)

	command.Execute(req)

  c.Assert(image.Generated, Equals, true)
}
