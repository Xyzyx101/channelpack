package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"
)

// ParseForm parses the form input and starts the conversion process
func ParseForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test worked...")
	fmt.Println("method:", r.Method)
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//r.ParseForm()
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

// ParseUpload will parse and save uploaded images
func ParseUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	r.ParseMultipartForm(1 << 10)
	fhs := r.MultipartForm.File["image-file"]
	for _, fh := range fhs {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		m, _, err := image.Decode(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		MyPackWorker.AddImage(fh.Filename, &m)
	}
	http.Redirect(w, r, "/", 303)
}

// ParseRemove is used to delete an uploaded image that is not longer needed
func ParseRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	indexProp, ok := r.Form["file-index"]
	if ok {
		index, err := strconv.Atoi(indexProp[0])
		if err != nil {
			http.Error(w, "Form contains no file-index property", http.StatusInternalServerError)
		}
		MyPackWorker.RemoveImage(index)
	}
	http.Redirect(w, r, "/", 303)
}
