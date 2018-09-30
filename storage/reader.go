package storage

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"
)

const FileExt = ".jpg"

type FaceReader struct {
}

// Read read all pics in folder, return array of persons
// path base path of faces with slash on end
func (fr *FaceReader) Read(dir string) ([]Person, error) {

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var prsns []Person
	for _, f := range fs {
		// dirs in root of dir is a names of persons, its contains person photos
		if f.IsDir() {
			p := Person{Name: f.Name()}
			imgs, err := fr.read(dir + f.Name())
			if err != nil {
				return nil, err
			}
			p.Imgs = imgs
			prsns = append(prsns, p)

		}

	}

	return prsns, nil
}

func (fr *FaceReader) ReadAsJpg(out chan image.Image, path string) {

	imgs, err := fr.read(path)
	if err != nil {
		panic(err)
	}

	for _, v := range imgs {
		fmt.Println(v)
		h, err := os.Open(v)
		if err != nil {
			panic(err)
		}
		pic, _, err := image.Decode(h)
		if err != nil {
			panic(err)
		}
		out <- pic
	}
}

func (fr *FaceReader) read(path string) ([]string, error) {
	var imgs []string

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		if !f.IsDir() {
			if strings.Contains(f.Name(), FileExt) {
				imgs = append(imgs, path+"/"+f.Name())
			}
		} else {
			inDir, err := fr.read(path + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			imgs = append(imgs, inDir...)
		}
	}

	return imgs, nil
}
