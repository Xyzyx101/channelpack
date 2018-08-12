package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var templates = make(map[string]*template.Template)

// initTemplates parses and verifies and caches of the templates channel packer will use
func initTemplates() {
	layout := template.Must(template.ParseFiles(filepath.Join("templates", "layout.gohtml")))

	indexFiles := []string{
		filepath.Join("templates", "layout.gohtml"),
		filepath.Join("templates", "index.gohtml"),
		filepath.Join("templates", "input.gohtml"),
		filepath.Join("templates", "output.gohtml"),
		filepath.Join("templates", "input-images.gohtml"),
	}
	index, err := layout.Clone()
	if err != nil {
		log.Fatal("cloning layout: ", err)
	}
	_, err = index.ParseFiles(indexFiles...)
	if err != nil {
		log.Fatal("parsing index: ", err)
	}
	templates["index"] = index

	foo, err := layout.Clone()
	if err != nil {
		log.Fatal("cloning layout: ", err)
	}
	_, err = foo.ParseFiles(filepath.Join("templates", "foo.gohtml"))
	if err != nil {
		log.Fatal("parsing foo: ", err)
	}
	templates["foo"] = foo
}

// serveStatic parses templates as serves the static css and js resources
func serveStatic(w http.ResponseWriter, r *http.Request) {
	var page *template.Template
	switch path := strings.TrimPrefix(r.URL.Path, "/"); path {
	case "favicon.ico":
		w.WriteHeader(http.StatusNoContent)
		return
	case "":
		page = templates["index"]
	default:
		if t, ok := templates[path]; ok {
			page = t
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}
	log.Println(maskPack)
	log.Println(allPackTypes)

	imageNames := myPackWorker.imageNames()
	imageChannels := myPackWorker.imageChannels()
	configData := struct {
		ImageNames    []string
		ImageChannels []string
		AllPackTypes  []packType
		AllChannels   []imageChannel
	}{imageNames, imageChannels, allPackTypes, allChannels}
	page.ExecuteTemplate(w, "layout", configData)
}

// serveThumbnail serves the images stored in PackWorker as web thumbnails
func serveThumbnail(w http.ResponseWriter, r *http.Request) {
	_, file := filepath.Split(r.URL.Path)
	if err := myPackWorker.serveThumbnail(w, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
