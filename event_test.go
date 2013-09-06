package photolibrary

import (
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
	images := make([]Image, 0, 2)
	images = append(images, FileSystemImage{Thumbnail: "mything.jpg"})
	images = append(images, FileSystemImage{Thumbnail: "anotherThing.jpg"})

	event := Event{Images: images}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].GetThumbnail(), Equals, "thumbnail_being_generated.jpg")
	c.Assert(eventWithTemp.Images[1].GetThumbnail(), Equals, "thumbnail_being_generated.jpg")
}

func (s *EventSuite) TestWeKeepTheFullPath(c *C) {
	images := make([]Image, 0, 2)
	images = append(images, FileSystemImage{FullPath: "mything.jpg"})
	images = append(images, FileSystemImage{FullPath: "anotherThing.jpg"})

	event := Event{Images: images}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].GetFullPath(), Equals, "mything.jpg")
	c.Assert(eventWithTemp.Images[1].GetFullPath(), Equals, "anotherThing.jpg")
}

func (s *EventSuite) TestWeKeepTheEventDescriptions(c *C) {
	event := Event{Images: []Image{},
		Events: []EventDescription{{FullName: "directory", ShortName: "d"}}}
	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Events, HasLen, 1)
	c.Assert(eventWithTemp.Events[0].FullName, Equals, "directory")
}

func (s *EventSuite) TestWeDontReplaceThumbnailsThatExist(c *C) {
	s.CreateFile("test.jpg", c)
	images := []Image{FileSystemImage{Thumbnail: s.directory + "/test.jpg"}}
	event := Event{Images: images}

	eventWithTemp := event.ReplaceMissingThumbnailsWithTemp()

	c.Assert(eventWithTemp.Images[0].GetThumbnail(), Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertRelativePathToFullPathForFullPath(c *C) {
	s.CreateFile("test.jpg", c)
	relativePath := s.RelativePathTo("test.jpg", c)
	images := []Image{FileSystemImage{FullPath: relativePath}}
	event := Event{Images: images}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

	c.Assert(eventWithFullPaths.Images[0].GetFullPath(), Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertRelativePathToFullPathForThumbnail(c *C) {
	s.CreateFile("test.jpg", c)
	relativePath := s.RelativePathTo("test.jpg", c)
	images := []Image{FileSystemImage{Thumbnail: relativePath}}
	event := Event{Images: images}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

	c.Assert(eventWithFullPaths.Images[0].GetThumbnail(), Equals, s.directory+"/test.jpg")
}

func (s *EventSuite) TestConvertEventFullPathsAsWell(c *C) {
	s.directory = c.MkDir()
	relativePath := s.RelativePathTo("please/", c)
	event := Event{Events: []EventDescription{{FullName: relativePath}}}

	eventWithFullPaths := event.ReplaceRelativePathsWithFullPaths()

	c.Assert(eventWithFullPaths.Events[0].FullName, Equals, s.directory+"/please")
}

// Generating jobs
// Past this - URL's.  You don't configure the img url or use it yet
//  - that happens in the view.  So the two things you needed to configure, you don't test
// Generating the two types of images
// Damn
