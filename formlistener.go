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
	if err != nil {
		log.Println(err)
	} else {
		log.Println(buildInstructions)
	}
	myPackWorker.createImage(*buildInstructions)
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
		myPackWorker.addImage(fh.Filename, m)
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
	var fileType outputContentType
	fileTypeParam, err := formValue(f, "file-type")
	if err == nil {
		fileType, err = parseOutputContentType(fileTypeParam)
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
	var packType *packType
	packTypeParam, err := formValue(f, "pack-type")
	if err == nil {
		packType, err = parsePackType(packTypeParam)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	var redInputChannel, greenInputChannel, blueInputChannel, alphaInputChannel, greyInputChannel *inputChannel
	switch {
	case packType.Equals(maskPack):
		fallthrough
	case packType.Equals(rgbaPack):
		alphaInputChannel, err = channelParams(f, "alpha-file", "alpha-channel")
		if err != nil {
			return nil, err
		}
		fallthrough
	case packType.Equals(rgbPack):
		redInputChannel, err = channelParams(f, "red-file", "red-channel")
		if err != nil {
			return nil, err
		}
		greenInputChannel, err = channelParams(f, "green-file", "green-channel")
		if err != nil {
			return nil, err
		}
		blueInputChannel, err = channelParams(f, "blue-file", "blue-channel")
		if err != nil {
			return nil, err
		}
	case packType.Equals(greyPack):
		greyInputChannel, err = channelParams(f, "grey-file", "grey-channel")
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unhandled pack type")
	}

	return &packInstructions{
		filename,
		fileType,
		width,
		height,
		redInputChannel,
		greenInputChannel,
		blueInputChannel,
		alphaInputChannel,
		greyInputChannel,
	}, nil
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

func channelParams(form url.Values, filename string, channelName string) (*inputChannel, error) {
	fileParam, err := formValue(form, filename)
	if err != nil {
		return nil, err
	}
	image, err := myPackWorker.image(fileParam)
	if err != nil {
		return nil, err
	}
	channelParam, err := formValue(form, channelName)
	if err != nil {
		return nil, err
	}
	channel := parseChannel(channelParam)
	return &inputChannel{image, channel}, nil
}
