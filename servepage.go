package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// ServeStatic parses templates as serves the static resources and form
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.gohtml")
	var page string
	if r.URL.Path == "/" {
		page = filepath.Join("templates", "index.gohtml")
	} else if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		page = filepath.Join("templates", filepath.Clean(r.URL.Path))
	}
	inputImages := filepath.Join("templates", "input-images.gohtml")
	imageNames := MyPackWorker.ImageNames()
	imageData := struct {
		Names []string
	}{imageNames}

	tmpl, err := template.ParseFiles(layout, page, inputImages)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", imageData)
}

// ServeThumbnail serves the images stored in PackWorker as web thumbnails
func ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	_, file := filepath.Split(r.URL.Path)
	if err := MyPackWorker.ServeThumbnail(w, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
