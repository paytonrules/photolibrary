package main

import (
	"github.com/paytonrules/photolibrary/library"
	. "launchpad.net/gocheck"
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


func (s *GenerateThumbnailsSuite) TestExecuteFindsTheEventsWithTheRightRoot(c *C) {
	phonyEvent := PhonyEvents{}
	command := MakeGenerateThumbnailCommand(&phonyEvent)

  command.generateThumbnailsForDirectory("directory")

	c.Assert(phonyEvent.FindString, Equals, "directory")
}

func (s *GenerateThumbnailsSuite) TestGeneratesThumbnailImages(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", library.Event{Images: []library.Image{image}})
	command := MakeGenerateThumbnailCommand(&phonyEvents)

	command.generateThumbnailsForDirectoryAndDuration("directory", 200)

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
	command.generateThumbnailsForDirectoryAndDuration("Root", 200)

	c.Assert(childImage.Generated, Equals, true)
}

func (s *GenerateThumbnailsSuite) TestDoesntGenerateThumbnailsAfterDuration(c *C) {
	phonyEvents := PhonyEvents{}
	image := &PhonyImage{}
	phonyEvents.FindResultFor("directory", library.Event{Images: []library.Image{image}})
	command := MakeGenerateThumbnailCommand(&phonyEvents)

	command.generateThumbnailsForDirectory("directory")

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
	command.generateThumbnailsForDirectory("Root")

	c.Assert(childImage.Generated, Equals, false)
}

func (s *GenerateThumbnailsSuite) TestLoggingEvents(c *C) {
  logger := new(TestLogger)
	phonyEvents := &PhonyEvents{}
  generator := MakeGenerateThumbnailCommandWithLogger(phonyEvents, logger)

  generator.generateThumbnailsForDirectory("Root")

  c.Assert(logger.info[0], Equals, "Generating Images in Root")
}
