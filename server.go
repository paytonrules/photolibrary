package main

import (
	"net/http"
  "github.com/paytonrules/image"
  "fmt"
  "log"
  "html"
)

func main() {
  http.HandleFunc("/generateThumbnails", func(w http.ResponseWriter, r *http.Request) {
    log.Print(r.Form)
    image := image.NewImage("bleh", "blue")
    image.GenerateThumbnail()

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":8080", nil)
}
