package controllers

import (
  "github.com/robfig/revel"
  "path/filepath"
  "os"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
  allFiles, error := filepath.Glob("public/events/*")
  if error != nil {
    return c.Render(error)
  }

  directories := make([]string, 0, len(allFiles))
  files := make([]string, 0, len(allFiles))

  for i := range allFiles {
    lstat, error := os.Lstat(allFiles[i])
    if error != nil {
      return c.Render(error)
    }

    if lstat.IsDir() {
      directories = append(directories, allFiles[i])
    } else {
      files = append(files, allFiles[i])
    }
  }
  return c.Render(files, directories)
}
