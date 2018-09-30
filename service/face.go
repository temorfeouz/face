package service

//
//import (
//	"github.com/temorfeouz/face/storage"
//
//	"github.com/Kagami/go-face"
//	"github.com/pkg/errors"
//)
//
//type FaceReconizer struct {
//	persons  []storage.Person
//	modelDir string
//	rec *face.Recognizer
//}
//
//type reader interface {
//	Read(dir string) ([]storage.Person, error)
//}
//
//func (fr *FaceReconizer) Init(r reader, dir, modelDir string) error {
//	prsns, err := r.Read(dir)
//	fr.persons = prsns
//
//	fr.modelDir = modelDir
//
//	return err
//}
//
//func (fr *FaceReconizer) processFaces() error {
//	recTmp, err := face.NewRecognizer(fr.modelDir)
//
//	if err != nil {
//		return errors.Errorf("Can't init face recognizer: %v", err)
//	}
//
//	fr.rec=recTmp
//
//	// Recognize faces on that image.
//
//	var (
//		//faces []face.Face
//		samples []face.Descriptor
//		cats []int32
//	)
//
//	facesStrs := fr.persons
//	for i, f := range facesStrs {
//
//		var (
//			face []face.Face
//			err  error
//		)
//		// load all faces
//		for _, userPic := range f.Imgs {
//			face, err = fr.rec.RecognizeFile(userPic)
//			if err != nil {
//				return errors.Errorf("Can't recognize: %v", err)
//			}
//
//			// add samples
//			for _, faceSample:=range face{
//				samples=append(samples, faceSample.Descriptor)
//				cats=append(cats, int32(i))
//			}
//		}
//	}
//
//	fr.rec.SetSamples(samples, cats)
//
//	return nil
//}
//
//func (fr *FaceReconizer) Reconize() (string, error) {
//	return "", nil
//}
