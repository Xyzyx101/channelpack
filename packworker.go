package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/nfnt/resize"
)

const progressChunks = 100.0

// packWorker holds the uploaded images and creates the final product for the user
type packWorker struct {
	UploadImages   []*uploadImage
	DownloadImage  *downloadImage
	Output         chan interface{}
	OutputProgress float32
}

// newPackWorker creates and initializes a new pack worker
func newPackWorker() *packWorker {
	uploadImages := make([]*uploadImage, 0, 4)
	output := make(chan interface{})
	close(output) // output starts closed because before processing '/output' return no content
	return &packWorker{uploadImages, nil, output, 0.0}
}

// addImage adds a newly uploaded image to the pack worker
func (p *packWorker) addImage(name string, image image.Image) {
	u := newUploadImage(name, image)
	p.UploadImages = append(p.UploadImages, u)
}

func (p packWorker) image(name string) (image.Image, error) {
	for _, image := range p.UploadImages {
		if name == image.name {
			return image.image, nil
		}
	}
	return nil, errors.New("Packworker has no image : " + name)
}

// imageNames returns a list of image names that can be used in the html template
func (p packWorker) imageNames() []string {
	var imageData = make([]string, 0, len(p.UploadImages))
	for _, image := range p.UploadImages {
		imageData = append(imageData, image.name)
	}
	return imageData
}

// imageChannels returns a list of input channel options that are valid for that image
func (p packWorker) imageChannels() []string {
	var ic = make([]string, 0, len(p.UploadImages))
	for _, image := range p.UploadImages {
		var channels channelConst
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
			channels = red | green | blue | alpha | grey
		case color.YCbCrModel:
			channels = red | green | blue | grey
		case color.AlphaModel:
			fallthrough
		case color.Alpha16Model:
			channels = alpha
		case color.GrayModel:
			fallthrough
		case color.Gray16Model:
			channels = grey
		default:
			channels = none
		}
		ic = append(ic, channels.String())
	}
	return ic
}

// removeImage an image from the worker
func (p *packWorker) removeImage(index int) error {
	if index < 0 || index >= len(p.UploadImages) {
		return errors.New("Tried to remove image with bad index")
	}
	p.UploadImages = append(p.UploadImages[:index], p.UploadImages[index+1:]...)
	return nil
}

// serveThumbnail converts and sends thumbnail versions of the uploaded images
func (p *packWorker) serveThumbnail(w http.ResponseWriter, filepath string) error {
	var image *uploadImage
	for _, needle := range p.UploadImages {
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

// serveOutput serves the processed image
func (p *packWorker) serveOutput(w http.ResponseWriter) error {
	output := <-p.Output
	switch output.(type) {
	case float32:
		newProgress, _ := output.(float32)
		p.OutputProgress += newProgress
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "progress:%f:%f", p.OutputProgress, progressChunks)
		break
	case string:
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", output)
		break
	default:
		w.WriteHeader(http.StatusNoContent)
		p.OutputProgress = 0.0
	}
	return nil
}

func (p *packWorker) createImage(i packInstructions) error {
	p.Output = make(chan interface{}, progressChunks+1)
	for i := 0; i < progressChunks; i++ {
		p.Output <- float32(1.0) //float32(i)
	}

	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return err
	}
	p.DownloadImage = newDownloadImage(buf.Bytes(), jpgContent)
	var href = "./download"
	var filename = "test.jpg"
	p.Output <- "download:" + href + ":" + filename
	return nil
}

func resizeImage(i *image.Image) *image.Image {
	return nil
}

func (p *packWorker) serveDownload(w http.ResponseWriter) error {
	r := bytes.NewReader(p.DownloadImage.data)
	w.Header().Set("Content-Type", p.DownloadImage.contentType.String())
	io.Copy(w, r)
	return nil
}
