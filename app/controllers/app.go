package controllers

import (
  "github.com/robfig/revel"
  "path/filepath"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
  files, error := filepath.Glob("*")
  if (error == nil) {
    return c.Render(files)
  }
  return c.Render(error)
}
