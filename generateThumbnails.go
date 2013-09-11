package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paytonrules/photolibrary"
	"github.com/paytonrules/thumbnailRequest"
	"net/http"
  "html/template"
	"time"
  "strconv"
)

var index = template.Must(template.ParseFiles(
  "templates/index.html",
))

type GenerateThumbnailsCommand struct {
	Events    photolibrary.Events
	startTime time.Time
	duration  time.Duration
}

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectory(directory string) {
	event, _ := c.Events.Find(directory)

	for _, img := range event.Images {
		if time.Since(c.startTime) < c.duration {
			img.GenerateThumbnail()
		}
	}

	for _, childEvent := range event.Events {
		if time.Since(c.startTime) < c.duration {
			c.generateThumbnailsForDirectory(childEvent.FullName)
		}
	}
}

func (c *GenerateThumbnailsCommand) generateThumbnailsForDirectoryAndDuration(directory string, duration int) {
  c.startTime = time.Now()
	c.duration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))
	c.generateThumbnailsForDirectory(directory)
}

func (c *GenerateThumbnailsCommand) Execute(r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request thumbnailRequest.Request
	decoder.Decode(&request)

	c.generateThumbnailsForDirectoryAndDuration(request.Directory, request.Duration)
}

func GenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	obj := GenerateThumbnailsCommand{Events: photolibrary.FileSystemEvents{}}

	obj.Execute(r)
}

func GenerateThumbnailsPost(w http.ResponseWriter, r *http.Request) {
	obj := GenerateThumbnailsCommand{Events: photolibrary.FileSystemEvents{}}

  duration, _ := strconv.Atoi(r.FormValue("duration"))
  obj.generateThumbnailsForDirectoryAndDuration(r.FormValue("directory"), duration)
}

func RenderTestPage(w http.ResponseWriter, r *http.Request) {
  index.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generateThumbnails", GenerateThumbnails)
	r.HandleFunc("/generateThumbnailsPost", GenerateThumbnailsPost)
  r.HandleFunc("/", RenderTestPage)
	http.Handle("/", r)
	http.ListenAndServe(":9001", nil)
}
