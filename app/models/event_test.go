package models

import (
	"github.com/paytonrules/image"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
)

type EventSuite struct {
	directory string
}

func (s *EventSuite) CreateFile(filename string, c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/" + filename)
	c.Assert(err, IsNil)
	defer file.Close()
}

func (s *EventSuite) RelativePathTo(filename string, c *C) string {
	workingDir, err := os.Getwd()
	c.Assert(err, IsNil)
	relativePath, err := filepath.Rel(workingDir, s.directory+"/"+filename)
	c.Assert(err, IsNil)

	return relativePath
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

func (s *EventSuite) TestWeDontReplaceThumbnailsThatExist(c *C) {
	s.CreateFile("test.jpg", c)
	images := []image.Image{image.Image{Thumbnail: s.directory + "/test.jpg"}}
	event := Event{Images: images}

	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].Thumbnail, Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertRelativePathToFullPathForFullPath(c *C) {
	s.CreateFile("test.jpg", c)
	relativePath := s.RelativePathTo("test.jpg", c)
	images := []image.Image{image.Image{FullPath: relativePath}}
	event := Event{Images: images}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

	c.Assert(eventWithFullPaths.Images[0].FullPath, Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertRelativePathToFullPathForThumbnail(c *C) {
	s.CreateFile("test.jpg", c)
	relativePath := s.RelativePathTo("test.jpg", c)
	images := []image.Image{image.Image{Thumbnail: relativePath}}
	event := Event{Images: images}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

	c.Assert(eventWithFullPaths.Images[0].Thumbnail, Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertDirectoriesAsWell(c *C) {
	s.directory = c.MkDir()
	relativePath := s.RelativePathTo("please/", c)
	event := Event{Events: []string{relativePath}}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

  c.Assert(eventWithFullPaths.Events[0], Equals, s.directory + "/please")
}

// Generating jobs
// Past this - URL's.  You don't configure the img url or use it yet
//  - that happens in the view.  So the two things you needed to configure, you don't test
// Damn
