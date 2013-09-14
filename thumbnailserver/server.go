package main

import (
	"github.com/gorilla/mux"
	"github.com/paytonrules/photolibrary/library"
	"html/template"
	"net/http"
	"strconv"
)

var index = template.Must(template.ParseFiles(
	"templates/index.html",
))

func GenerateThumbnailsPost(w http.ResponseWriter, r *http.Request) {
	obj := GenerateThumbnailsCommand{Events: library.FileSystemEvents{}}

	duration, _ := strconv.Atoi(r.FormValue("duration"))
	obj.generateThumbnailsForDirectoryAndDuration(r.FormValue("directory"), duration)
}

func GenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	obj := GenerateThumbnailsCommand{Events: library.FileSystemEvents{}}

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
