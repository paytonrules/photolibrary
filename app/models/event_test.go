package models

import (
	"github.com/paytonrules/image"
	. "launchpad.net/gocheck"
	"os"
)

type EventSuite struct {
	directory string
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) TearDownTest(c *C) {
	if s.directory != "" {
		os.RemoveAll(s.directory)
	}
}

func (s *EventSuite) TestReplacingThumbnailsWithPlaceHolders(c *C) {
	images := make([]image.Image, 0, 2)
	images = append(images, image.Image{Thumbnail: "mything.jpg"})
	images = append(images, image.Image{Thumbnail: "anotherThing.jpg"})

	event := Event{Images: images}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].Thumbnail, Equals, "thumbnail_being_generated.jpg")
	c.Assert(eventWithTemp.Images[1].Thumbnail, Equals, "thumbnail_being_generated.jpg")
}

func (s *EventSuite) TestWeKeepTheFullPath(c *C) {
	images := make([]image.Image, 0, 2)
	images = append(images, image.Image{FullPath: "mything.jpg"})
	images = append(images, image.Image{FullPath: "anotherThing.jpg"})

	event := Event{Images: images}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].FullPath, Equals, "mything.jpg")
	c.Assert(eventWithTemp.Images[1].FullPath, Equals, "anotherThing.jpg")
}

func (s *EventSuite) TestWeKeepTheEvents(c *C) {
	event := Event{Images: []image.Image{},
		Events: []string{"directory"}}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Events, HasLen, 1)
	c.Assert(eventWithTemp.Events[0], Equals, "directory")
}

func (s *EventSuite) TestweDontReplaceThumbnailsThatExist(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/test.jpg")
	c.Assert(err, IsNil)
	defer file.Close()

  images := make([]image.Image, 0, 1)
  images = append(images, image.Image{Thumbnail: s.directory + "/test.jpg"}) 
	event := Event{Images: images}
  eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].Thumbnail, Equals, s.directory + "/test.jpg")
}

// It doesnt replace if the thumbnail exists
// It doesnt lose the Event

// Past this - URL's.  You don't configure the img url or use it yet
// Full paths and generating jobs
