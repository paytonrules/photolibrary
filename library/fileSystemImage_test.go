package library

import (
	. "launchpad.net/gocheck"
)

type FileSystemImageSuite struct {}
var _ = Suite(&FileSystemImageSuite{})

func (s *FileSystemImageSuite) TestCloneMakesACopy(c *C) {
  originalImage := FileSystemImage{Thumbnail: "thumb", FullPath: "path"}
  copiedImage := originalImage.Clone()

  c.Assert(copiedImage.GetThumbnail(), Equals, "thumb")
  c.Assert(copiedImage.GetFullPath(), Equals, "path")
}

func (s *FileSystemImageSuite) TestCloneDoesNotReturnTheSameObject(c *C) {
  originalImage := FileSystemImage{Thumbnail: "thumb", FullPath: "path"}
  copiedImage := originalImage.Clone()

  originalImage.Thumbnail = "not thumb"

  c.Assert(copiedImage.GetThumbnail(), Equals, "thumb")
}
