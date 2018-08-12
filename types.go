package main

import (
	"errors"
	"image"
	"image/color"
	"strings"
)

type imageChannel int

const (
	channelError imageChannel = -1
	r            imageChannel = iota
	g
	b
	a
	grey
)

func (c imageChannel) String() string {
	switch c {
	case r:
		return "r"
	case g:
		return "g"
	case b:
		return "b"
	case a:
		return "a"
	case grey:
		return "grey"
	default:
		return "error channel"
	}
}

func parseImageChannel(s string) (imageChannel, error) {
	s = strings.ToLower(s)
	switch s {
	case "r":
		return r, nil
	case "g":
		return g, nil
	case "b":
		return b, nil
	case "a":
		return a, nil
	default:
		return channelError, errors.New("Unable to parse image channel : " + s)
	}
}

type outputFileType int

const (
	outputFileTypeError outputFileType = -1
	png                 outputFileType = iota
	jpg
	tga
)

func (f outputFileType) String() string {
	switch f {
	case png:
		return "png"
	case jpg:
		return "jpg"
	case tga:
		return "tga"
	default:
		return "output file type error"
	}
}

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
		return outputFileTypeError, errors.New("Unable to parse output file type : " + s)
	}
}

type packInstructions struct {
	filename                string
	fileType                outputFileType
	width, height           int
	red, green, blue, alpha packChannel
}

type packChannel struct {
	name    string
	channel imageChannel
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
