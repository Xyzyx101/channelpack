package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/nfnt/resize"
)

// packWorker holds the uploaded images and creates the final product for the user
type packWorker struct {
	Images []*uploadImage
}

// newPackWorker creates and initializes a new pack worker
func newPackWorker() *packWorker {
	images := make([]*uploadImage, 0, 4)
	return &packWorker{images}
}

// addImage adds a newly uploaded image to the pack worker
func (p *packWorker) addImage(name string, image *image.Image) {
	u := newUploadImage(name, image)
	p.Images = append(p.Images, u)
}

// imageNames returns a list of image names that can be used in the html template
func (p packWorker) imageNames() []string {
	var imageData = make([]string, 0, len(p.Images))
	for _, image := range p.Images {
		imageData = append(imageData, image.name)
	}
	return imageData
}

// imageChannels returns a list of input channel options that are valid for that image
func (p packWorker) imageChannels() []string {
	var ic = make([]string, 0, len(p.Images))
	for _, image := range p.Images {
		var validChannels string
		cm := image.ColorModel()
		switch cm {
		case color.RGBAModel:
			fallthrough
		case color.RGBA64Model:
			fallthrough
		case color.NRGBAModel:
			fallthrough
		case color.NRGBA64Model:
			fallthrough
		case color.NYCbCrAModel:
			validChannels = "R|G|B|A|Grey"
		case color.YCbCrModel:
			validChannels = "R|G|B|Grey"
		case color.AlphaModel:
			fallthrough
		case color.Alpha16Model:
			validChannels = "A"
		case color.GrayModel:
			fallthrough
		case color.Gray16Model:
			validChannels = "Grey"
		default:
			log.Println("Unknown colour model")
			validChannels = "XXX"
		}
		ic = append(ic, validChannels)
	}
	return ic
}

// removeImage an image from the worker
func (p *packWorker) removeImage(index int) error {
	if index < 0 || index >= len(p.Images) {
		return errors.New("Tried to remove image with bad index")
	}
	p.Images = append(p.Images[:index], p.Images[index+1:]...)
	return nil
}

// serveThumbnail converts and sends thumbnail versions of the uploaded images
func (p *packWorker) serveThumbnail(w http.ResponseWriter, filepath string) error {
	var image *uploadImage
	for _, needle := range p.Images {
		if needle.name == filepath {
			image = needle
		}
	}
	if image == nil {
		return errors.New("No image found for filepath")
	}

	if len(image.thumb) == 0 {
		thumb := resize.Resize(240, 240, image, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, thumb, nil)
		if err != nil {
			return err
		}
		image.thumb = buf.Bytes()
	}
	r := bytes.NewReader(image.thumb)
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, r)
	return nil
}
