package main

import (
	"fmt"
	"image/color"
_"github.com/oliamb/cutter"
	"gocv.io/x/gocv"
	"github.com/oliamb/cutter"
	"os"
	"bufio"
	"image"
	"bytes"
	"image/jpeg"
	"time"
)
//https://github.com/Kagami/go-face
func main() {
	// set to use a video capture device 0
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	
	defer classifier.Close()

	if !classifier.Load("data/haarcascade_frontalface_default.xml") {
		fmt.Println("Error reading cascade file: data/haarcascade_frontalface_default.xml")
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image
		for _, r := range rects {
			//os.Exit(1)
saveCropped(img,r)
			
			gocv.Rectangle(&img, r, blue, 1)
		}

		// show the image in the window, and wait 1 millisecond
		
		
		window.IMShow(img)
		window.WaitKey(1)
	}
}

func saveCropped(img gocv.Mat,  r image.Rectangle){
	pic , err:=img.ToImage()
	if err!=nil{
		panic(err)
	}

	croppedImg, err := cutter.Crop(pic, cutter.Config{
		Width:  r.Max.X-r.Min.X,
		Height: r.Max.Y-r.Min.Y,
		Anchor:image.Point{
			X:r.Min.X,
			Y:r.Min.Y,
		},
		//Mode:cutter.Centered,
	})
	if err!=nil{
		panic(err)
	}
	
	// Write to file.
	fo, err := os.Create(fmt.Sprintf("imgs/img_%d.jpg", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)

	// Encode to jpeg.
	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, croppedImg, nil)
	
	fw.Write(imageBuf.Bytes())
}