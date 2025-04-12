package assetsgen

import (
	"errors"
	"image"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type AndroidFolderName string

const (
	Mipmap   AndroidFolderName = "mipmap"
	Drawable AndroidFolderName = "drawable"

	rootFolderName string = "assets_gen_out"
)

type imageInfo struct {
	img              image.Image
	imagePath        string
	imageName        string
	imageExt         string
	encoder          imgio.Encoder
	genImageLocation func(screenType string) (directory string, imageName string)
}

type screenType struct {
	name        string
	scaleFactor float64
}

type screenTypeSlice []screenType

func (s screenTypeSlice) maxScaleFactor() float64 {
	var maxFactor float64
	for _, v := range s {
		maxFactor = max(v.scaleFactor, maxFactor)
	}
	return maxFactor
}

// MDPI - 1.0x
// HDPI - 1.5x
// XHDPI - 2.0x
// XXHDPI - 3x
// XXXHDPI - 4.0x
var androidScreenDpis = screenTypeSlice{
	screenType{
		name:        "mdpi",
		scaleFactor: 1.0,
	},
	screenType{
		name:        "hdpi",
		scaleFactor: 1.5,
	},
	screenType{
		name:        "xhdpi",
		scaleFactor: 2,
	},
	screenType{
		name:        "xxhdpi",
		scaleFactor: 3,
	},
	screenType{
		name:        "xxxhdpi",
		scaleFactor: 4,
	},
}

func imageEncoderFromPath(imagePath string) (imgio.Encoder, error) {
	ext := filepath.Ext(imagePath)
	switch ext {
	case ".png":
		return imgio.PNGEncoder(), nil
	case ".jpeg", ".jpg":
		return imgio.JPEGEncoder(100), nil
	case ".bmp":
		return imgio.BMPEncoder( ), nil

	default:
		return nil, ErrUnsupportedFileType
	}

}
