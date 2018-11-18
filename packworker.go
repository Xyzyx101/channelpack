package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"golang.org/x/image/tiff"

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
	for output := range p.Output {
		switch output.(type) {
		case float64:
			p.OutputProgress += float32(output.(float64))
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "progress:%f:%f", p.OutputProgress, progressChunks)
			return nil
		case string:
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "%s", output)
			return nil
		}
	}
	w.WriteHeader(http.StatusNoContent)
	p.OutputProgress = 0.0
	return nil
	// output := <-p.Output
	// switch output.(type) {
	// case float64:
	// 	p.OutputProgress += float32(output.(float64))
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Header().Set("Content-Type", "text/plain")
	// 	fmt.Fprintf(w, "progress:%f:%f", p.OutputProgress, progressChunks)
	// 	break
	// case string:
	// 	w.WriteHeader(http.StatusCreated)
	// 	fmt.Fprintf(w, "%s", output)
	// 	break
	// default:
	// 	w.WriteHeader(http.StatusNoContent)
	// 	p.OutputProgress = 0.0
	// }
	// return nil
}

func (p *packWorker) createImage(pInst packInstructions) error {
	p.Output = make(chan interface{}, progressChunks+1)
	resizeInputImages(&pInst, p.Output)
	var img image.Image
	var err error
	switch *pInst.packType {
	case maskPack:
		img, err = createMaskPack(pInst)
	case rgbPack:
		img, err = createRGBPack(pInst)
	case rgbaPack:
		img, err = createRGBAPack(pInst)
	case greyPack:
		img, err = createGreyPack(pInst)
	}
	if err != nil {
		return err
	}
	p.Output <- 20.0
	buf := new(bytes.Buffer)
	switch pInst.outputType {
	case jpgContent:
		err = jpeg.Encode(buf, img, nil)
	case pngContent:
		err = png.Encode(buf, img)
	case tiffContent:
		err = tiff.Encode(buf, img, nil)
	}
	if err != nil {
		return err
	}
	p.Output <- 30.0
	p.DownloadImage = newDownloadImage(buf.Bytes(), pInst.outputType)
	var href = "./download"
	var filename = pInst.outputName
	p.Output <- "download:" + href + ":" + filename
	close(p.Output)
	return nil
}

func createMaskPack(pInst packInstructions) (image.Image, error) {
	width := pInst.width
	height := pInst.height
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := getColorFromChannel(x, y, pInst.red)
			g := getColorFromChannel(x, y, pInst.green)
			b := getColorFromChannel(x, y, pInst.blue)
			a := getColorFromChannel(x, y, pInst.alpha)
			pixelCol := color.NRGBA{r, g, b, a}
			img.Set(x, y, pixelCol)
		}
	}
	return img, nil
}

func createRGBPack(pInst packInstructions) (image.Image, error) {
	width := pInst.width
	height := pInst.height
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := getColorFromChannel(x, y, pInst.red)
			g := getColorFromChannel(x, y, pInst.green)
			b := getColorFromChannel(x, y, pInst.blue)
			a := uint8(255)
			pixelCol := color.RGBA{r, g, b, a}
			img.Set(x, y, pixelCol)
		}
	}
	return img, nil
}
func createRGBAPack(pInst packInstructions) (image.Image, error) {
	width := pInst.width
	height := pInst.height
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := getColorFromChannel(x, y, pInst.red)
			g := getColorFromChannel(x, y, pInst.green)
			b := getColorFromChannel(x, y, pInst.blue)
			a := getColorFromChannel(x, y, pInst.alpha)
			pixelCol := color.RGBA{r, g, b, a}
			img.Set(x, y, pixelCol)
		}
	}
	return img, nil
}
func createGreyPack(pInst packInstructions) (image.Image, error) {
	width := pInst.width
	height := pInst.height
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewGray(image.Rectangle{upLeft, lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grey := getColorFromChannel(x, y, pInst.grey)
			pixelCol := color.Gray{grey}
			img.Set(x, y, pixelCol)
		}
	}
	return img, nil
}

func getColorFromChannel(x, y int, iChan *inputChannel) uint8 {
	rgbaCol := iChan.image.At(x, y)
	if iChan.channel == grey {
		greyCol := color.GrayModel.Convert(rgbaCol)
		grey, _ := greyCol.(color.Gray)
		return grey.Y
	}
	nrgbaCol := color.NRGBAModel.Convert(rgbaCol)
	nrgba, _ := nrgbaCol.(color.NRGBA)
	switch iChan.channel {
	case red:
		return uint8(nrgba.R)
	case green:
		return uint8(nrgba.G)
	case blue:
		return uint8(nrgba.B)
	case alpha:
		return uint8(nrgba.A)
	}
	return 0
}

func resizeInputImages(pInst *packInstructions, progress chan<- interface{}) {
	inputChannels := [5]*inputChannel{pInst.red, pInst.green, pInst.blue, pInst.alpha, pInst.grey}
	for _, iC := range inputChannels {
		func(pInst *packInstructions) {
			defer func() {
				progress <- 10.0
			}()
			if iC == nil {
				return
			}
			if iC.image.Bounds().Dx() == pInst.width && iC.image.Bounds().Dy() == pInst.height {
				return
			}
			resized := resize.Resize(uint(pInst.width), uint(pInst.height), iC.image, resize.Lanczos3)
			*iC = inputChannel{resized, iC.channel}
		}(pInst)
	}
}

func (p *packWorker) serveDownload(w http.ResponseWriter) error {
	r := bytes.NewReader(p.DownloadImage.data)
	w.Header().Set("Content-Type", p.DownloadImage.contentType.String())
	io.Copy(w, r)
	return nil
}
