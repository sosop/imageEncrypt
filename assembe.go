// Package imageEncrypt assembe 拼接图片
package imageEncrypt

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"sync"

	"github.com/sosop/imaging"
)

// Assembe 接口
type Assembe interface {
	assembing(condition ...interface{}) ([]byte, error)
}

// FileSystemAssembe 从文件系统中获取图片进行封装
type FileSystemAssembe struct {
	s Storage
	m Meta
}

// NewFileSystemAssembe 构造
func NewFileSystemAssembe(s Storage, m Meta) *FileSystemAssembe {
	return &FileSystemAssembe{s, m}
}

func (a *FileSystemAssembe) assembing(condition ...interface{}) ([]byte, error) {
	metaImage, err := a.m.get(condition)
	if err != nil {
		return nil, err
	}
	// 创建图片
	full := image.NewNRGBA(image.Rect(0, 0, metaImage.MaxX, metaImage.MaxY))
	n := len(metaImage.Images)
	wg := new(sync.WaitGroup)
	wg.Add(n)
	flag := true
	for _, cuttedImage := range metaImage.Images {
		go drawIt(a.s, cuttedImage, full, &flag, wg)
	}
	wg.Wait()
	if !flag {
		return nil, errors.New("加载失败")
	}
	imaging.Save(full, fmt.Sprint("test", metaImage.Ext))
	/*
		buf := bytes.NewBuffer(nil)
		f, _ := formats[metaImage.Ext]
		err := imaging.Encode(buf, full, f)
		if err != nil {
			return nil, err
		}
		data := base64.StdEncoding.EncodeToString(buf.Bytes())
	*/
	return nil, nil
}

func drawIt(s Storage, cuttedImage CuttedImage, bg *image.NRGBA, flag *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	rc, err := s.get(cuttedImage.Location)
	if err != nil {
		*flag = false
		return
	}
	defer rc.Close()
	img, err := imaging.Decode(rc)
	if err != nil {
		*flag = false
		return
	}
	invImg := inverseRotate(img, cuttedImage.Rotate)
	draw.Draw(bg, image.Rect(cuttedImage.Points[0].X, cuttedImage.Points[0].Y, cuttedImage.Points[1].X, cuttedImage.Points[1].Y), invImg, image.Pt(0, 0), draw.Src)
}
