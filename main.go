package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"time"

	"io/ioutil"

	"strings"

	"github.com/Kagami/go-face"
	"github.com/oliamb/cutter"
	"gocv.io/x/gocv"
)

const (
	imgFolder   = "imgs"
	imgBaseName = "img_"
)

//https://github.com/Kagami/go-face
func main() {
	// set to use a video capture device 0
	deviceID := 0
	reconizePhotos()
	os.Exit(1)
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
			saveCropped(img, r)
			gocv.Rectangle(&img, r, blue, 1)
		}

		// show the image in the window, and wait 1 millisecond

		window.IMShow(img)
		window.WaitKey(1)
	}
}

func saveCropped(img gocv.Mat, r image.Rectangle) {
	pic, err := img.ToImage()
	if err != nil {
		panic(err)
	}

	croppedImg, err := cutter.Crop(pic, cutter.Config{
		Width:  r.Max.X - r.Min.X,
		Height: r.Max.Y - r.Min.Y,
		Anchor: image.Point{
			X: r.Min.X,
			Y: r.Min.Y,
		},
		//Mode:cutter.Centered,
	})
	if err != nil {
		panic(err)
	}

	// Write to file.
	fo, err := os.Create(fmt.Sprintf("imgs/img_%d.jpg", time.Now().UnixNano()))
	if err != nil {
		fmt.Printf("%v", err)
		//panic(err)
	}
	fw := bufio.NewWriter(fo)

	// Encode to jpeg.
	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, croppedImg, nil)

	fw.Write(imageBuf.Bytes())
	fw.Flush()
}

func readFiles() []string {
	dir := "./" + imgFolder + "/"
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, f := range fs {
		if !f.IsDir() && strings.Contains(f.Name(), ".jpg") {
			files = append(files, dir+f.Name())
		}
	}

	return files
}

//
func reconizePhotos() {
	// Init the recognizer.
	rec, err := face.NewRecognizer("models")
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	// Recognize faces on that image.

	var faces []face.Face
	facesStrs := readFiles()
	for _, f := range facesStrs {
		face, err := rec.RecognizeFile(f)
		if err != nil {
			log.Fatalf("Can't recognize: %v", err)
		}
		faces = append(faces, face...)
	}

	// Fill known samples. In the real world you would use a lot of images
	// for each person to get better classification results but in our
	// example we just get them from one big image.
	var samples []face.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		cats = append(cats, int32(i))
	}
	// Name the categories, i.e. people on the image.
	labels := []string{
		"Sungyeon", "Yehana", "Roa", "Eunwoo", "Xiyeon",
		"Kyulkyung", "Nayoung", "Rena", "Kyla", "Yuha",
	}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, cats)

	// Now let's try to classify some not yet known image.
	testImageNayoung := facesStrs[2]
	nayoungFace, err := rec.RecognizeSingleFile(testImageNayoung)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if nayoungFace == nil {
		log.Fatalf("Not a single face on the image")
	}
	catID := rec.Classify(nayoungFace.Descriptor)
	if catID < 0 {
		log.Fatalf("Can't classify")
	}
	// Finally print the classified label. It should be "Nayoung".
	fmt.Println(labels[catID])
}
