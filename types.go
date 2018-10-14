package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"
)

type channelConst uint8

const (
	none  channelConst = 0x00
	all   channelConst = 0xff
	red   channelConst = 0x01
	green channelConst = 0x02
	blue  channelConst = 0x04
	alpha channelConst = 0x08
	grey  channelConst = 0x10
)

func (c channelConst) String() string {
	var names []string
	if c&red > 0 {
		names = append(names, "red")
	}
	if c&green > 0 {
		names = append(names, "green")
	}
	if c&blue > 0 {
		names = append(names, "blue")
	}
	if c&alpha > 0 {
		names = append(names, "alpha")
	}
	if c&grey > 0 {
		names = append(names, "grey")
	}
	return strings.Join(names, "|")
}

func (c channelConst) PrettyName() string {
	var names []string
	if c&red > 0 {
		names = append(names, "Red Channel")
	}
	if c&green > 0 {
		names = append(names, "Green Channel")
	}
	if c&blue > 0 {
		names = append(names, "Blue Channel")
	}
	if c&alpha > 0 {
		names = append(names, "Alpha Channel")
	}
	if c&grey > 0 {
		names = append(names, "Greyscale")
	}
	return strings.Join(names, "|")
}

type channel interface {
	fmt.Stringer
	PrettyName() string
}

// var (
// 	redChannel   = imageChannel{"red", "Red Channel"}
// 	greenChannel = imageChannel{"green", "Green Channel"}
// 	blueChannel  = imageChannel{"blue", "Blue Channel"}
// 	alphaChannel = imageChannel{"alpha", "Alpha Channel"}
// 	greyChannel  = imageChannel{"grey", "Greyscale Channel"}
// )

//var allChannels = []channel{red, green, blue, alpha, grey}
var allChannelsForJS = []struct{ Name, PrettyName string }{
	{red.String(), red.PrettyName()},
	{green.String(), green.PrettyName()},
	{blue.String(), blue.PrettyName()},
	{alpha.String(), alpha.PrettyName()},
	{grey.String(), grey.PrettyName()},
}

func parseChannel(s string) channel {
	channel := none
	s = strings.Replace(s, "|", "", -1)
	if strings.Contains(s, red.String()) {
		channel |= red
		s = strings.Replace(s, "red", "", -1)
	}
	if strings.Contains(s, green.String()) {
		channel |= green
		s = strings.Replace(s, "green", "", -1)
	}
	if strings.Contains(s, blue.String()) {
		channel |= blue
		s = strings.Replace(s, "blue", "", -1)
	}
	if strings.Contains(s, alpha.String()) {
		channel |= alpha
		s = strings.Replace(s, "alpha", "", -1)
	}
	if strings.Contains(s, grey.String()) {
		channel |= grey
		s = strings.Replace(s, "grey", "", -1)
	}
	return channel
}

type packType struct {
	Name    string
	Channel channel
}

var (
	maskPack = packType{"Mask", red | green | blue | alpha}
	rgbPack  = packType{"RGB", red | green | blue}
	rgbaPack = packType{"RGBA", red | green | blue | alpha}
	greyPack = packType{"Greyscale", grey}
)

func (p packType) Equals(other packType) bool {
	if p.Name != other.Name {
		return false
	}
	if p.Channel != other.Channel {
		return false
	}
	return true
}

func parsePackType(s string) (*packType, error) {
	switch s {
	case maskPack.Name:
		return &maskPack, nil
	case rgbaPack.Name:
		return &rgbaPack, nil
	case rgbPack.Name:
		return &rgbPack, nil
	case greyPack.Name:
		return &greyPack, nil
	default:
		return nil, errors.New("Unable to parse packType : " + s)
	}
}

var allPackTypesForJS = []struct{ Name, ImageChannels string }{
	{maskPack.Name, maskPack.Channel.String()},
	{rgbPack.Name, rgbPack.Channel.String()},
	{rgbaPack.Name, rgbaPack.Channel.String()},
	{greyPack.Name, greyPack.Channel.String()},
}

type outputContentType string

const (
	pngContent outputContentType = "image/png"
	jpgContent outputContentType = "image/jpeg"
	tgaContent outputContentType = "image/x-tga"
)

func (o outputContentType) String() string {
	return string(o)
}
func parseOutputContentType(s string) (outputContentType, error) {
	s = strings.ToLower(s)
	switch s {
	case "png":
		return pngContent, nil
	case "jpg":
		return jpgContent, nil
	case "tga":
		return tgaContent, nil
	default:
		return "", errors.New("Unable to parse output file type : " + s)
	}
}

type packInstructions struct {
	outputName                    string
	outputType                    outputContentType
	width, height                 int
	red, green, blue, alpha, grey *inputChannel
}

type inputChannel struct {
	image   image.Image
	channel channel
}

// uploadImage represents the uploaded image from the user.
// name is the filename.uploadImage
// image is the actual image data
// thumb is a shrunk down version to display over the internet
type uploadImage struct {
	name  string
	image image.Image
	thumb []byte
}

// ColorModel is required for UploadImage to satisfy the Image interface
func (u *uploadImage) ColorModel() color.Model { return u.image.ColorModel() }

// Bounds is required for UploadImage to satisfy the Image interface
func (u *uploadImage) Bounds() image.Rectangle { return u.image.Bounds() }

// At is required for UploadImage to satisfy the Image interface
func (u *uploadImage) At(x, y int) color.Color { return u.image.At(x, y) }

// newUploadImage creates a new UploadImage
func newUploadImage(name string, image image.Image) *uploadImage {
	return &uploadImage{name, image, nil}
}

type downloadImage struct {
	data        []byte
	contentType outputContentType
}

func newDownloadImage(data []byte, contentType outputContentType) *downloadImage {
	return &downloadImage{data, contentType}
}
