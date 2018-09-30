package service

import (
	"fmt"
	"image"
	"image/color"

	"github.com/temorfeouz/face/storage"

	"gocv.io/x/gocv"
)

type FaceFinder struct {
	classifier gocv.CascadeClassifier
	blue       color.RGBA
}

func (ff *FaceFinder) Dispose() {
	ff.classifier.Close()
}

func (ff *FaceFinder) Init() {
	ff.blue = color.RGBA{0, 0, 255, 0}

	ff.classifier = gocv.NewCascadeClassifier()

	if !ff.classifier.Load("data/haarcascade_frontalface_default.xml") {
		panic("Error reading cascade file: data/haarcascade_frontalface_default.xml")
	}
}
func (ff *FaceFinder) FindFace(pic image.Image) ([]gocv.Mat, error) {

	img, err := gocv.ImageToMatRGB(pic)
	if err != nil {
		return nil, err
	}

	// detect faces
	rects := ff.classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		storage.SaveCropped("imgs/img_%d.jpg", img, r)
		gocv.Rectangle(&img, r, ff.blue, 1)
	}

	return nil, nil
}
