package main

import (
	"bytes"
	"encoding/json"
	"github.com/paytonrules/photolibrary/library"
	"github.com/paytonrules/photolibrary/thumbnailrequest"
	. "launchpad.net/gocheck"
	"net/http"
)

type GenerateThumbnailsSuite struct{}

var _ = Suite(&GenerateThumbnailsSuite{})

type TestLogger struct {
  info []string
}

func (l *TestLogger) Info(message string) {
  if l.info == nil {
    l.info = []string{}
  }

  l.info = append(l.info, message)
}

type PhonyEvents struct {
	FullPaths   bool
	FindString  string
	FindResults map[string]library.Event
}

func (evts *PhonyEvents) Find(eventName string) (library.Event, error) {
	evts.FindString = eventName
	if evts.FindResults != nil {
		return evts.FindResults[eventName], nil
	} else {
		return library.Event{}, nil
	}
}

func (evts *PhonyEvents) FindResultFor(eventName string, evt library.Event) {
	if evts.FindResults == nil {
		evts.FindResults = make(map[string]library.Event)
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

func (img PhonyImage) Clone() library.Image {
	return nil
}

func (s *GenerateThumbnailsSuite) marshalThumbnailRequest(directory string, duration int) *http.Request {
	thumbnailRequest := thumbnailrequest.Request{Directory: directory, Duration: duration}
	marshaledThumbnailRequest, _ := json.Marshal(thumbnailRequest)
	body := bytes.NewBuffer(marshaledThumbnailRequest)
	req, _ := http.NewRequest("dont", "care", body)
	return req
}

func (s *GenerateThumbnailsSuite) TestExecuteFindsTheEventsWithTheRightRoot(c *C) {
	phonyEvent := PhonyEvents{}
	command := MakeGenerateThumbnailCommand(&phonyEvent)

	req := s.marshalThumbnailRequest("directory", 0)
	command.Execute(req)

	c.Assert(phonyEvent.FindString, Equals, "directory")
}

func (s *GenerateThumbnailsSuite) TestGeneratesThumbnailImages(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", library.Event{Images: []library.Image{image}})
	command := MakeGenerateThumbnailCommand(&phonyEvents)

	req := s.marshalThumbnailRequest("directory", 200)
	command.Execute(req)

	c.Assert(image.Generated, Equals, true)
}

func (s *GenerateThumbnailsSuite) TestGeneratesThumnailImagesForChildEvents(c *C) {
	phonyEvents := PhonyEvents{}
	eventDescription := library.EventDescription{FullName: "full name"}
	rootEvent := library.Event{Events: []library.EventDescription{eventDescription}}
	childImage := &PhonyImage{}
	childEvent := library.Event{Images: []library.Image{childImage}}

	phonyEvents.FindResultFor("Root", rootEvent)
	phonyEvents.FindResultFor("full name", childEvent)

	command := MakeGenerateThumbnailCommand(&phonyEvents)
	req := s.marshalThumbnailRequest("Root", 200)
	command.Execute(req)

	c.Assert(childImage.Generated, Equals, true)
}

func (s *GenerateThumbnailsSuite) TestDoesntGenerateThumbnailsAfterDuration(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", library.Event{Images: []library.Image{image}})
	command := MakeGenerateThumbnailCommand(&phonyEvents)

	req := s.marshalThumbnailRequest("directory", 0)
	command.Execute(req)

	c.Assert(image.Generated, Equals, false)
}

func (s *GenerateThumbnailsSuite) TestItDoesntContinueDownTheEventTreePastTheDuration(c *C) {
	phonyEvents := PhonyEvents{}
	eventDescription := library.EventDescription{FullName: "full name"}
	rootEvent := library.Event{Events: []library.EventDescription{eventDescription}}
	childImage := &PhonyImage{}
	childEvent := library.Event{Images: []library.Image{childImage}}

	phonyEvents.FindResultFor("Root", rootEvent)
	phonyEvents.FindResultFor("full name", childEvent)

	command := MakeGenerateThumbnailCommand(&phonyEvents)
	req := s.marshalThumbnailRequest("Root", 0)
	command.Execute(req)

	c.Assert(childImage.Generated, Equals, false)
}

func (s *GenerateThumbnailsSuite) TestLoggingEvents(c *C) {
  logger := new(TestLogger)
	phonyEvents := &PhonyEvents{}
  generator := MakeGenerateThumbnailCommandWithLogger(phonyEvents, logger)

  generator.generateThumbnailsForDirectory("Root")

  c.Assert(logger.info[0], Equals, "Generating Images in Root")
}
