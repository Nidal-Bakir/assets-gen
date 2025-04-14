package assetsgen

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type platformType string

const (
	platformTypeAndroid platformType = "android"
	platformTypeIos     platformType = "ios"
)

type intention string

const (
	intentAppIcon          intention = "app_icon"
	intentNotificationIcon intention = "notification_icon"
	intentAsset            intention = "asset"
)

type imageInfo struct {
	img              image.Image
	imagePath        string
	imageName        string
	imageExt         string
	encoder          imgio.Encoder
	genImageLocation func(screenType string) (directory string, imageName string)
}

type Asset interface {
	Name() string
	CalcSize(w, h int) (int, int)
}

func imageEncoderFromPath(imagePath string) (imgio.Encoder, error) {
	ext := filepath.Ext(imagePath)
	switch ext {
	case ".png":
		return imgio.PNGEncoder(), nil
	case ".jpeg", ".jpg":
		return imgio.JPEGEncoder(100), nil
	case ".bmp":
		return imgio.BMPEncoder(), nil

	default:
		return nil, ErrUnsupportedFileType
	}
}

func generateImageInfo(imagePath string, platform platformType, intent intention, lastFolderName func(screenType string) string) (imageInfo, error) {
	img, err := imgio.Open(imagePath)
	if err != nil {
		return imageInfo{}, err
	}

	enc, err := imageEncoderFromPath(imagePath)
	if err != nil {
		return imageInfo{}, err
	}

	imgName := filepath.Base(imagePath)
	imageExt := filepath.Ext(imagePath)
	imgNameWithoutExt := strings.ReplaceAll(imgName, imageExt, "")

	imgInfo := imageInfo{
		img:       img,
		encoder:   enc,
		imagePath: imagePath,
		imageName: imgName,
		imageExt:  imageExt,
		genImageLocation: func(screenType string) (directory string, imageName string) {
			dir := filepath.Join(
				rootFolderName,
				string(platform),
				string(intent),
				imgNameWithoutExt,
				lastFolderName(screenType),
			)
			return dir, imgName
		},
	}

	return imgInfo, nil
}

func squareImageWithPadding(img image.Image) image.Image {
	imageBounds := img.Bounds()
	w := imageBounds.Dx()
	h := imageBounds.Dy()

	if w == h { // it's already a square
		return img
	}

	// to center the image in the Square
	var wOffset, hOffset int

	if w < h {
		w += (h - w)
		wOffset = w / 3
	} else {
		h += (w - h)
		hOffset = h / 3
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h - hOffset {
		for x := range w - wOffset {
			dst.Set(x+wOffset, y+hOffset, img.At(x, y))
		}
	}

	return dst
}

func generateImageAsstes(imgInfo imageInfo, assets []Asset) error {
	imageBounds := imgInfo.img.Bounds()

	for _, asset := range assets {
		w, h := asset.CalcSize(imageBounds.Dx(), imageBounds.Dy())

		resizedImg := transform.Resize(imgInfo.img, w, h, transform.Linear)

		dir, name := imgInfo.genImageLocation(asset.Name())
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		err = imgio.Save(filepath.Join(dir, name), resizedImg, imgInfo.encoder)
		if err != nil {
			return err
		}
	}

	return nil
}
