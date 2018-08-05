package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
)

// ListenForm waits for the form to be submitted from the static stite and parses it
// func ListenForm() {
// 	http.HandleFunc("/process", Test)
// }

// ParseForm parses the form input and starts the conversion process
func ParseForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test worked...")
	fmt.Println("method:", r.Method)
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r.ParseForm()

	// r.ParseMultipartForm(32 << 20)
	// fmt.Println("test:", r.Form["test"])
	// fmt.Println("file:", r.Form["myFile"])

	r.ParseMultipartForm(1 << 10)
	fhs := r.MultipartForm.File["image-file"]
	for _, fh := range fhs {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		m, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		bounds := m.Bounds()
		fmt.Println(bounds)
	}
	// f is one of the files

	//fhs := r.MultipartForm.File["myfiles"]
	// for _, fh := range fhs {
	// 	f, err := fh.Open()
	// 	// f is one of the files
	// 	//f.Read
	// }

}
