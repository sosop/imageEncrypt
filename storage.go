// Package imageEncrypt storage 保存被切图片
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

// Storage 存储接口
type Storage interface {
	save(image *CuttedImage, subImage image.Image, filename string, wg *sync.WaitGroup, exts ...string)
	get(path ...string) (io.ReadCloser, error)
}

// FileStorage 文件存储
type FileStorage struct {
	dir string
}

// NewFileStorage 构造文件存储
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

// 字节缓冲区
func data(img image.Image, ext string) (*bytes.Buffer, error) {
	f, ok := formats[ext]
	if !ok {
		return nil, imaging.ErrUnsupportedFormat
	}
	buf := bytes.NewBuffer(nil)
	err := imaging.Encode(buf, img, f)
	return buf, err
}
