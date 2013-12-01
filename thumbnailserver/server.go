package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/paytonrules/photolibrary/library"
	"github.com/paytonrules/photolibrary/thumbnailrequest"
	"html/template"
	"net/http"
	"strconv"
)

var index = template.Must(template.ParseFiles(
	"templates/index.html",
))

type ThumbnailGenerator interface {
  generateThumbnailsForDirectoryAndDuration(directory string, duration int)
}

var makeCommand = func() ThumbnailGenerator {
	logger := new(GoLogger)
  return MakeGenerateThumbnailCommandWithLogger(library.MakeFileSystemEvents([]string{".jpg", ".png"}),
		logger)
}

func GenerateThumbnailsPost(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  requestBody := r.FormValue("directory")
	duration, _ := strconv.Atoi(r.FormValue("duration"))

  obj := makeCommand()
	obj.generateThumbnailsForDirectoryAndDuration(requestBody, duration)
}

func GenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request thumbnailrequest.Request
  decoder.Decode(&request)

  obj := makeCommand()
  obj.generateThumbnailsForDirectoryAndDuration(request.Directory, request.Duration)
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
