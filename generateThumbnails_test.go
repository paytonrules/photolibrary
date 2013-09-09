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
	FullPaths   bool
	FindString  string
	FindResults map[string]photolibrary.Event
}

func (evts *PhonyEvents) Find(eventName string) (photolibrary.Event, error) {
	evts.FindString = eventName
	if evts.FindResults != nil {
		return evts.FindResults[eventName], nil
	} else {
		return photolibrary.Event{}, nil
	}
}

func (evts *PhonyEvents) FindResultFor(eventName string, evt photolibrary.Event) {
	if evts.FindResults == nil {
		evts.FindResults = make(map[string]photolibrary.Event)
	}
	evts.FindResults[eventName] = evt
}

type PhonyImage struct {
	FullPath  string
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

func (img PhonyImage) Clone() photolibrary.Image {
	return nil
}

func (s *GenerateThumbnailsSuite) marshalThumbnailRequest(directory string, duration int) *http.Request {
	thumbnailRequest := thumbnailRequest.Request{Directory: directory, Duration: duration}
	marshaledThumbnailRequest, _ := json.Marshal(thumbnailRequest)
	body := bytes.NewBuffer(marshaledThumbnailRequest)
	req, _ := http.NewRequest("dont", "care", body)
	return req
}

func (s *GenerateThumbnailsSuite) TestExecuteFindsTheEventsWithTheRightRoot(c *C) {
	phonyEvent := PhonyEvents{}
	command := GenerateThumbnailsCommand{Events: &phonyEvent}

	req := s.marshalThumbnailRequest("directory", 0)
	command.Execute(req)

	c.Assert(phonyEvent.FindString, Equals, "directory")
}

func (s *GenerateThumbnailsSuite) TestGeneratesThumbnailImages(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", photolibrary.Event{Images: []photolibrary.Image{image}})
	command := GenerateThumbnailsCommand{Events: &phonyEvents}

	req := s.marshalThumbnailRequest("directory", 200)
	command.Execute(req)

	c.Assert(image.Generated, Equals, true)
}

func (s *GenerateThumbnailsSuite) TestGeneratesThumnailImagesForChildEvents(c *C) {
	phonyEvents := PhonyEvents{}
	eventDescription := photolibrary.EventDescription{FullName: "full name"}
	rootEvent := photolibrary.Event{Events: []photolibrary.EventDescription{eventDescription}}
	childImage := &PhonyImage{}
	childEvent := photolibrary.Event{Images: []photolibrary.Image{childImage}}

	phonyEvents.FindResultFor("Root", rootEvent)
	phonyEvents.FindResultFor("full name", childEvent)

	command := GenerateThumbnailsCommand{Events: &phonyEvents}
	req := s.marshalThumbnailRequest("Root", 200)
	command.Execute(req)

	c.Assert(childImage.Generated, Equals, true)
}

func (s *GenerateThumbnailsSuite) TestDoesntGenerateThumbnailsAfterDuration(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", photolibrary.Event{Images: []photolibrary.Image{image}})
	command := GenerateThumbnailsCommand{Events: &phonyEvents}

	req := s.marshalThumbnailRequest("directory", 0)
	command.Execute(req)

	c.Assert(image.Generated, Equals, false)
}

func (s *GenerateThumbnailsSuite) TestItDoesntContinueDownTheEventTreePastTheDuration(c *C) {
	phonyEvents := PhonyEvents{}
	eventDescription := photolibrary.EventDescription{FullName: "full name"}
	rootEvent := photolibrary.Event{Events: []photolibrary.EventDescription{eventDescription}}
	childImage := &PhonyImage{}
	childEvent := photolibrary.Event{Images: []photolibrary.Image{childImage}}

	phonyEvents.FindResultFor("Root", rootEvent)
	phonyEvents.FindResultFor("full name", childEvent)

	command := GenerateThumbnailsCommand{Events: &phonyEvents}
	req := s.marshalThumbnailRequest("Root", 0)
	command.Execute(req)

	c.Assert(childImage.Generated, Equals, false)
}


