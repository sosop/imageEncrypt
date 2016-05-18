// Package imageEncrypt storage is storing the splice image
package imageEncrypt

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"sync"

	"github.com/sosop/imaging"
)

// Storage interface
type Storage interface {
	save(image *CuttedImage, subImage image.Image, filename string, wg *sync.WaitGroup, exts ...string)
	get(path ...string) (io.ReadCloser, error)
}

// FileStorage Use file system to store splice image
type FileStorage struct {
	dir string
}

// NewFileStorage constructor
func NewFileStorage(dir string) *FileStorage {
	return &FileStorage{dir}
}

func (s *FileStorage) save(image *CuttedImage, subImage image.Image, filename string, wg *sync.WaitGroup, exts ...string) {
	defer wg.Done()
	fullname := fmt.Sprint(s.dir, fmt.Sprintf("%x", md5.Sum([]byte(filename))), exts[0])
	err := imaging.Save(subImage, fullname)
	if err != nil {
		return
	}
	image.Location = fullname
}

func (s *FileStorage) get(paths ...string) (io.ReadCloser, error) {
	if len(paths) == 0 {
		return nil, errors.New("paths is empty")
	}
	f, err := os.Open(paths[0])
	if err != nil {
		return nil, err
	}
	return f, nil
}

// byte buffer
func data(img image.Image, ext string) (*bytes.Buffer, error) {
	f, ok := formats[ext]
	if !ok {
		return nil, imaging.ErrUnsupportedFormat
	}
	buf := bytes.NewBuffer(nil)
	err := imaging.Encode(buf, img, f)
	return buf, err
}
