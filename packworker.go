package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/nfnt/resize"
)

// UploadImage represents the uploaded image from the user.
// name is the filename.UploadImage
// image is the actual image data
// thumb is a shrunk down version to display over the internet
type UploadImage struct {
	name  string
	image *image.Image
	thumb []byte
}

// ColorModel is required for UploadImage to satisfy the Image interface
func (u *UploadImage) ColorModel() color.Model { return (*u.image).ColorModel() }

// Bounds is required for UploadImage to satisfy the Image interface
func (u *UploadImage) Bounds() image.Rectangle { return (*u.image).Bounds() }

// At is required for UploadImage to satisfy the Image interface
func (u *UploadImage) At(x, y int) color.Color { return (*u.image).At(x, y) }

// NewUploadImage creates a new UploadImage
func NewUploadImage(name string, image *image.Image) *UploadImage {
	return &UploadImage{name, image, nil}
}

// PackWorker holds the uploaded images and creates the final product for the user
type PackWorker struct {
	Images []*UploadImage
}

// NewPackWorker creates and initializes a new pack worker
func NewPackWorker() *PackWorker {
	images := make([]*UploadImage, 0, 4)
	return &PackWorker{images}
}

// AddImage adds a newly uploaded image to the pack worker
func (p *PackWorker) AddImage(name string, image *image.Image) {
	u := NewUploadImage(name, image)
	p.Images = append(p.Images, u)
}

// ImageNames returns a list of image names that can be used in the html template
func (p PackWorker) ImageNames() []string {
	var imageData = make([]string, 0, len(p.Images))
	for _, image := range p.Images {
		imageData = append(imageData, image.name)
	}
	return imageData
}

// RemoveImage an image from the worker
func (p *PackWorker) RemoveImage(index int) error {
	if index < 0 || index >= len(p.Images) {
		return errors.New("Tried to remove image with bad index")
	}
	p.Images = append(p.Images[:index], p.Images[index+1:]...)
	return nil
}

// ServeThumbnail converts and sends thumbnail versions of the uploaded images
func (p *PackWorker) ServeThumbnail(w http.ResponseWriter, filepath string) error {
	var image *UploadImage
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
