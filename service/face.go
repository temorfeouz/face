package service

import (
	"github.com/temorfeouz/face/storage"

	"github.com/Kagami/go-face"
	"github.com/pkg/errors"
)

type FaceReconizer struct {
	persons  []storage.Person
	modelDir string
}

type reader interface {
	Read(dir string) ([]storage.Person, error)
}

func (fr *FaceReconizer) Init(r reader, dir, modelDir string) error {
	prsns, err := r.Read(dir)
	fr.persons = prsns

	fr.modelDir = modelDir

	return err
}

func (fr *FaceReconizer) processFaces() error {
	rec, err := face.NewRecognizer(fr.modelDir)
	if err != nil {
		return errors.Errorf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	// Recognize faces on that image.

	var faces []face.Face
	facesStrs := fr.persons
	for _, f := range facesStrs {

		var (
			face []face.Face
			err  error
		)
		// load all faces
		for _, userPic := range f.Imgs {
			face, err = rec.RecognizeFile(userPic)
			if err != nil {
				return errors.Errorf("Can't recognize: %v", err)
			}
			faces = append(faces, face...)
		}

		// categorize

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
}

func (fr *FaceReconizer) Reconize() (string, error) {
	return "", nil
}
