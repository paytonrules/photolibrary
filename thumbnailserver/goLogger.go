package main

import (
  "github.com/golang/glog"
)

type GoLogger struct{}


func (g *GoLogger) Info(message string) {
  glog.MaxSize = 1024 * 1024 * 10 // 10 megabytes max
  glog.Info(message)
}
