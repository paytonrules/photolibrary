package main

import (
	. "launchpad.net/gocheck"
  "github.com/golang/glog"
)

type GoLoggerSuite struct{}

var _ = Suite(&GoLoggerSuite{})

func (s *GoLoggerSuite) TestLogsToTheInfoLogger(c *C) {
  g := &GoLogger{}

  g.Info("Info Message")

  c.Assert(glog.Stats.Info.Lines(), Equals, int64(1))
}
