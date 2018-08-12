package main

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func parseProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	buildInstructions, err := buildPackInstructions(r.Form)
	log.Println(buildInstructions, err)
	http.Redirect(w, r, "/", 303)
}

// parseUpload will parse and save uploaded images
func parseUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
		myPackWorker.addImage(fh.Filename, &m)
	}
	http.Redirect(w, r, "/", 303)
}

// parseRemove is used to delete an uploaded image that is not longer needed
func parseRemove(w http.ResponseWriter, r *http.Request) {
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
		myPackWorker.removeImage(index)
	}
	http.Redirect(w, r, "/", 303)
}

func buildPackInstructions(f url.Values) (*packInstructions, error) {
	filename, err := formValue(f, "filename")
	if err != nil {
		return nil, err
	}
	var fileType outputFileType
	fileTypeParam, err := formValue(f, "file-type")
	if err == nil {
		fileType, err = parseOutputFileType(fileTypeParam)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	var width int
	widthParam, err := formValue(f, "width")
	if err == nil {
		width, err = strconv.Atoi(widthParam)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	var height int
	heightParam, err := formValue(f, "height")
	if err == nil {
		height, err = strconv.Atoi(heightParam)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	log.Println(filename)
	log.Println(fileType)
	log.Println(width)
	log.Println(height)
	// hasAlpha := r.Form["has-alpha"]
	// redFile := r.Form["red-file"]
	// redChannel := r.Form["red-channel"]
	// greenFile := r.Form["green-file"]
	// greenChannel := r.Form["green-channel"]
	// blueFile := r.Form["blue-file"]
	// blueChannel := r.Form["blue-channel"]
	// alphaFile := r.Form["alpha-file"]
	// alphaChannel := r.Form["alpha-channel"]

	return &packInstructions{}, nil
}

func formValue(f url.Values, param string) (string, error) {
	values := f[param]
	if len(values) == 0 {
		return "", errors.New(param + " was expected and not found")
	}
	if len(values) > 1 {
		return "", errors.New(param + " expected 1 value and found " + strconv.Itoa(len(values)))
	}
	value := values[0]
	if len(value) == 0 {
		return "", errors.New(param + " was expected but found empty string")
	}
	return value, nil
}
