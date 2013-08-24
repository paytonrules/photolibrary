package main

import (
	"net/http"
  "encoding/json"
  "github.com/paytonrules/image"
  "log"
)

func main() {
  http.HandleFunc("/generateThumbnails", func(w http.ResponseWriter, r *http.Request) {

    decoder := json.NewDecoder(r.Body)
    var images []image.Image
    err := decoder.Decode(&images)
    if err != nil {
      log.Println(err)
    }
    log.Println(images)

    for _, thumbnailImage := range images {
      thumbnailImage.GenerateThumbnail()
    }
	})

	http.ListenAndServe(":8081", nil)
}
