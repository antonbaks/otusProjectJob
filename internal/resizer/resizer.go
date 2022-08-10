package resizer

import (
	"errors"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

var (
	ErrResize    = errors.New("resize error")
	ErrSize      = errors.New("incorrect size")
	ErrImageSize = errors.New("incorrect image size")
)

type Resizer struct {
	MaxWidth  uint
	MinWidth  uint
	MaxHeight uint
	MinHeight uint
}

func NewResizer(maxWidth, minWidth, maxHeight, minHeight uint) Resizer {
	return Resizer{
		MaxWidth:  maxWidth,
		MinWidth:  minWidth,
		MaxHeight: maxHeight,
		MinHeight: minHeight,
	}
}

func (r Resizer) ResizeImage(w, h uint, file *os.File) error {
	if w > r.MaxWidth || w < r.MinWidth {
		return ErrSize
	}

	if h > r.MaxHeight || h < r.MinHeight {
		return ErrSize
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Seek err: ", err)
		return ErrResize
	}

	cfg, err := jpeg.DecodeConfig(file)
	if err != nil {
		log.Println("Decode err: ", err)
		return ErrResize
	}

	if uint(cfg.Width) < w {
		log.Println("incorrect image width: ", cfg.Width)
		return ErrImageSize
	}

	if uint(cfg.Height) < h {
		log.Println("incorrect image height: ", cfg.Height)
		return ErrImageSize
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Seek err: ", err)
		return ErrResize
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Decode err: ", err)
		return ErrResize
	}

	m := resize.Resize(w, h, img, resize.Lanczos3)

	if err := file.Truncate(0); err != nil {
		log.Println("Truncate err: ", err)
		return ErrResize
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Seek err: ", err)
		return ErrResize
	}

	if err := jpeg.Encode(file, m, nil); err != nil {
		log.Println(err)
		return ErrResize
	}

	return nil
}
