package main

import (
  "net/http"
  . "launchpad.net/gocheck"
)

type ServerSuite struct {
  directory string
  duration int
}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) generateThumbnailsForDirectoryAndDuration(directory string, duration int) {
  s.directory = directory;
  s.duration = duration;
}

func (s *ServerSuite) TestWeCreateThumbnailCommandWithRequestBodyOnPost(c *C) {
  r, _ := http.NewRequest("POST", "http://a.com?directory=TheBody&duration=300", nil)

  makeCommand = func() ThumbnailGenerator {
    return s
  }

  GenerateThumbnailsPost(nil, r)

  c.Assert(s.directory, Equals, "TheBody")
  c.Assert(s.duration, Equals, 300)
}

func (s *ServerSuite) TestWeHandleSpacesOnDirectory(c *C) {
  r, _ := http.NewRequest("POST", "http://a.com?directory=The+Body&duration=300", nil)

  makeCommand = func() ThumbnailGenerator {
    return s
  }

  GenerateThumbnailsPost(nil, r)

  c.Assert(s.directory, Equals, "The Body")
  c.Assert(s.duration, Equals, 300)
}
