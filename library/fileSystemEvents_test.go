package library

import (
	. "launchpad.net/gocheck"
	"os"
)

type FileSystemEventsSuite struct {
	directory string
}

var _ = Suite(&FileSystemEventsSuite{})

func (s *FileSystemEventsSuite) TearDownTest(c *C) {
	if s.directory != "" {
		os.RemoveAll(s.directory)
	}
}

func (s *FileSystemEventsSuite) TestAnEmptyDirectory(c *C) {
	s.directory = c.MkDir()
	events := MakeFileSystemEvents([]string{})

	event, err := events.Find(s.directory)

	c.Assert(event.Images, HasLen, 0)
	c.Assert(err, IsNil)
}

func (s *FileSystemEventsSuite) TestADirectoryWithOneImage(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/test.jpg")
	c.Assert(err, IsNil)
	defer file.Close()

  events := MakeFileSystemEvents([]string {".jpg"})
	event, _ := events.Find(s.directory)

	c.Assert(event.Images, HasLen, 1)
	c.Assert(event.Images[0].GetFullPath(), Equals, s.directory+"/test.jpg")
}

func (s *FileSystemEventsSuite) TestDoesNotIncludeHiddenFiles(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/.DS_Store")
	c.Assert(err, IsNil)
	defer file.Close()

	events := MakeFileSystemEvents([]string {})
	event, _ := events.Find(s.directory)

	c.Assert(event.Images, HasLen, 0)
}

func (s *FileSystemEventsSuite) TestReturnsAnErrorFromABadDirectory(c *C) {
	events := MakeFileSystemEvents([]string{})
	_, err := events.Find("[]")

	// syntax error in pattern
	c.Assert(err, ErrorMatches, "syntax error in pattern")
}

func (s *FileSystemEventsSuite) TestIncludingDirectoriesAsEvents(c *C) {
	s.directory = c.MkDir()
	err := os.Mkdir(s.directory+"/Events", 0755)
	c.Assert(err, IsNil)

	events := MakeFileSystemEvents([]string{})
	event, err := events.Find(s.directory)

	c.Assert(event.Images, HasLen, 0)
	c.Assert(event.Events, HasLen, 1)
	c.Assert(event.Events[0].FullName, Equals, s.directory+"/Events")
	c.Assert(event.Events[0].ShortName, Equals, "Events")
}

func (s *FileSystemEventsSuite) TestMakingAPathToEachThumbnail(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/silly.jpg")
	c.Assert(err, IsNil)
	defer file.Close()

	events := MakeFileSystemEvents([]string{".jpg"})
	event, err := events.Find(s.directory)

	c.Assert(event.Images[0].GetThumbnail(), Equals, s.directory+"/.thumbnails/silly.jpg")
}

func (s *FileSystemEventsSuite) TestExcludingUnsupportedFiles(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/silly.mov")
	c.Assert(err, IsNil)
	defer file.Close()

	events := MakeFileSystemEvents([]string{".jpg"})
	event, err := events.Find(s.directory)

  c.Assert(event.Images, HasLen, 0)
}

func (s *FileSystemEventsSuite) TestExcludingUnsupportedFilesIsCaseInsensitive(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/silly.JPG")
	c.Assert(err, IsNil)
	defer file.Close()

	events := MakeFileSystemEvents([]string{".jpg"})
	event, err := events.Find(s.directory)

	c.Assert(event.Images[0].GetThumbnail(), Equals, s.directory+"/.thumbnails/silly.JPG")
}

func (s *FileSystemEventsSuite) TestExcludingUnsupportedFilesIsCaseInsensitiveOnExtension(c *C) {
	s.directory = c.MkDir()
	file, err := os.Create(s.directory + "/silly.jpg")
	c.Assert(err, IsNil)
	defer file.Close()

	events := MakeFileSystemEvents([]string{".JPG"})
	event, err := events.Find(s.directory)

	c.Assert(event.Images[0].GetThumbnail(), Equals, s.directory+"/.thumbnails/silly.jpg")
}

