package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paytonrules/photolibrary/library"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

var index = template.Must(template.ParseFiles(
	"templates/index.html",
))

func GenerateThumbnailsPost(w http.ResponseWriter, r *http.Request) {
	logger := new(GoLogger)
	requestBody, _ := ioutil.ReadAll(r.Body)
  requestAsString := fmt.Sprintf("Request body %s", requestBody)
	logger.Info("Recieved Generate Thumbnails Request as Post" + requestAsString)
  logger.Info("FormValue for directory is: " + r.FormValue("directory"))
	obj := MakeGenerateThumbnailCommandWithLogger(library.MakeFileSystemEvents([]string{".jpg", ".png"}),
		new(GoLogger))

	duration, _ := strconv.Atoi(r.FormValue("duration"))
	obj.generateThumbnailsForDirectoryAndDuration(r.FormValue("directory"), duration)
}

func GenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	logger := new(GoLogger)
	requestBody, _ := ioutil.ReadAll(r.Body)
  requestAsString := fmt.Sprintf("Request body %s", requestBody)
	logger.Info("Recieved Generate Thumbnails Request " + requestAsString)
	obj := MakeGenerateThumbnailCommandWithLogger(library.MakeFileSystemEvents([]string{".jpg", ".png"}),
		new(GoLogger))

	obj.Execute(r)
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
