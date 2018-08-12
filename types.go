package main

import (
	"errors"
	"image"
	"image/color"
	"strings"
)

type imageChannel struct {
	Name, PrettyName string
}

var (
	redChannel   = imageChannel{"red", "Red Channel"}
	greenChannel = imageChannel{"green", "Green Channel"}
	blueChannel  = imageChannel{"blue", "Blue Channel"}
	alphaChannel = imageChannel{"alpha", "Alpha Channel"}
	greyChannel  = imageChannel{"grey", "Greyscale Channel"}
)

type packType struct {
	Name          string
	ImageChannels []imageChannel
}

var (
	maskPack = packType{"Mask", []imageChannel{redChannel, greenChannel, blueChannel, alphaChannel}}
	rgbPack  = packType{"RGB", []imageChannel{redChannel, greenChannel, blueChannel}}
	rgbaPack = packType{"RGBA", []imageChannel{redChannel, greenChannel, blueChannel, alphaChannel}}
	greyPack = packType{"Greyscale", []imageChannel{greyChannel}}
)

func parsePackType(s string) (*packType, error) {
	for _, p := range allPackTypes {
		if s == p.Name {
			return &p, nil
		}
	}
	return nil, errors.New("Unable to parse packType : " + s)
}

var allPackTypes = []packType{maskPack, rgbPack, rgbaPack, greyPack}

type outputFileType string

const (
	png outputFileType = "png"
	jpg outputFileType = "jpg"
	tga outputFileType = "tga"
)

func parseOutputFileType(s string) (outputFileType, error) {
	s = strings.ToLower(s)
	switch s {
	case "png":
		return png, nil
	case "jpg":
		return jpg, nil
	case "tga":
		return tga, nil
	default:
		return "", errors.New("Unable to parse output file type : " + s)
	}
}

type packInstructions struct {
	outputName                    string
	outputType                    outputFileType
	width, height                 int
	red, green, blue, alpha, grey inputChannel
}

type inputChannel struct {
	filename string
	channel  imageChannel
}

// uploadImage represents the uploaded image from the user.
// name is the filename.uploadImage
// image is the actual image data
// thumb is a shrunk down version to display over the internet
type uploadImage struct {
	name  string
	image *image.Image
	thumb []byte
}

// ColorModel is required for UploadImage to satisfy the Image interface
func (u *uploadImage) ColorModel() color.Model { return (*u.image).ColorModel() }

// Bounds is required for UploadImage to satisfy the Image interface
func (u *uploadImage) Bounds() image.Rectangle { return (*u.image).Bounds() }

// At is required for UploadImage to satisfy the Image interface
func (u *uploadImage) At(x, y int) color.Color { return (*u.image).At(x, y) }

// newUploadImage creates a new UploadImage
func newUploadImage(name string, image *image.Image) *uploadImage {
	return &uploadImage{name, image, nil}
}
