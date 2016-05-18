// Package imageEncrypt cut 切图操作
package imageEncrypt

import (
	"image"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sosop/imaging"
)

const (
	// DefaultPatitionX 默认x
	DefaultPatitionX = 4
	// DefaultPatitionY 默认y
	DefaultPatitionY = 4
)

// Cut 切图接口
type Cut interface {
	Cutting(reader io.Reader, filename string, condition ...interface{}) (MetaCuttedImage, error)
}

// RectangleCut 矩形切图
type RectangleCut struct {
	// 横向切几份
	partitionX int
	// 纵向切几份
	partitionY int

	// 切片存储
	storage Storage
	// 切片元信息存储
	meta Meta
}

// NewDefaultRectangleCut 构造默认
func NewDefaultRectangleCut(storage Storage, meta Meta) *RectangleCut {
	return NewRectangleCut(DefaultPatitionX, DefaultPatitionY, storage, meta)
}

// NewRectangleCut 指定值
func NewRectangleCut(partitionX, patitionY int, storage Storage, meta Meta) *RectangleCut {
	return &RectangleCut{partitionX: partitionX, partitionY: patitionY, storage: storage, meta: meta}
}

// Cutting 实现切图接口
func (r RectangleCut) Cutting(reader io.Reader, filename string, condition ...interface{}) (*MetaCuttedImage, error) {
	// 获取图片扩展名判断类型
	ext := strings.ToLower(filepath.Ext(filename))
	src, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	// 计算原始图片大小
	rect := src.Bounds()
	x := rect.Max.X - rect.Min.X
	y := rect.Max.Y - rect.Min.X
	// 横向步长
	stepX := x / r.partitionX
	// 纵向步长
	stepY := y / r.partitionY
	images := make([]CuttedImage, r.partitionX*r.partitionY)
	k := 0
	// 多个gor切割上传
	wg := new(sync.WaitGroup)
	wg.Add(r.partitionX * r.partitionY)
	for row := 0; row < r.partitionY; row++ {
		for col := 0; col < r.partitionX; col++ {
			images[k] = CuttedImage{ID: k}
			p1 := Point{}
			p2 := Point{}
			if col > 0 {
				p1.X = images[k-1].Points[1].X
			}
			if row > 0 {
				p1.Y = images[k-r.partitionX].Points[1].Y
			}
			if col == r.partitionX-1 {
				p2.X = x
			} else {
				p2.X = p1.X + stepX
			}
			if row == r.partitionY-1 {
				p2.Y = y
			} else {
				p2.Y = p1.Y + stepY
			}
			images[k].Points = []Point{p1, p2}
			retangle := image.Rect(p1.X, p1.Y, p2.X, p2.Y)
			subImg := imaging.Crop(src, retangle)
			img := rotate(subImg, &images[k])
			go r.storage.save(&images[k], img, strconv.Itoa(k)+filename, wg, ext)
			k++
		}
	}
	wg.Wait()
	metaImage := MetaCuttedImage{images, x, y, Rectangle, ext}
	r.meta.save(metaImage, condition)
	return &metaImage, nil
}
