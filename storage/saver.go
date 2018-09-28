package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"

	"github.com/oliamb/cutter"
	"gocv.io/x/gocv"
)

func SaveCropped(tmpl string, img gocv.Mat, r image.Rectangle) {
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
	fo, err := os.Create(fmt.Sprintf(tmpl, time.Now().UnixNano()))
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
