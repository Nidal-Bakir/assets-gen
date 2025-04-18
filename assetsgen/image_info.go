package assetsgen

import (
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

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

func newImageInfo(imagePath string, platform platformType, intent intention, lastFolderName func(screenType string) string) (imageInfo, error) {
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

type imageInfo struct {
	img              image.Image
	imagePath        string
	imageName        string
	imageExt         string
	encoder          imgio.Encoder
	genImageLocation func(screenType string) (directory string, imageName string)
}

func (imgInfo *imageInfo) squareImageWithPadding(padding int) *imageInfo {
	imageBounds := imgInfo.img.Bounds()
	w := imageBounds.Dx() + padding
	h := imageBounds.Dy() + padding

	if w == h { // it's already a square
		return imgInfo
	}

	// to center the image in the Square
	wOffset := padding / 2
	hOffset := padding / 2
	if w < h {
		wOffset += (h - w) / 2
	} else {
		hOffset += (w - h) / 2
	}

	w = int(math.Max(float64(w), float64(h)))
	h = w

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h - hOffset {
		for x := range w - wOffset {
			dst.Set(x+wOffset, y+hOffset, imgInfo.img.At(x, y))
		}
	}

	imgInfo.img = dst
	return imgInfo
}

func (imgInfo *imageInfo) padding(padding int) *imageInfo {
	imageBounds := imgInfo.img.Bounds()
	w := imageBounds.Dx() + padding
	h := imageBounds.Dy() + padding

	wOffset := padding / 2
	hOffset := padding / 2

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h - hOffset {
		for x := range w - wOffset {
			dst.Set(x+wOffset, y+hOffset, imgInfo.img.At(x, y))
		}
	}

	imgInfo.img = dst
	return imgInfo
}

func (imgInfo *imageInfo) convertColors(fn func(color.RGBA) color.RGBA) *imageInfo {
	imgInfo.img = adjust.Apply(imgInfo.img, func(pxColor color.RGBA) color.RGBA {
		return fn(pxColor)
	})
	return imgInfo
}

func (imgInfo *imageInfo) convertOpaqueToColor(newColor color.RGBA) *imageInfo {
	return imgInfo.convertColors(func(pxColor color.RGBA) color.RGBA {
		if pxColor.A == 0 {
			return newColor
		}
		return pxColor
	})
}

func (imgInfo *imageInfo) convertNoneOpaqueToColor(newColor color.RGBA) *imageInfo {
	return imgInfo.convertColors(func(pxColor color.RGBA) color.RGBA {
		if pxColor.A == 0 {
			return pxColor
		}
		return newColor
	})
}

func (imgInfo *imageInfo) clipRRect(r int) *imageInfo {
	imageBounds := imgInfo.img.Bounds()
	w := imageBounds.Dx()
	h := imageBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			if isOnRoundedCorner(x, y, w, h, r) {
				continue
			}
			dst.Set(x, y, imgInfo.img.At(x, y))
		}
	}

	imgInfo.img = dst
	return imgInfo
}

func isOnRoundedCorner(x, y, w, h, r int) bool {
	// Top-left
	if x < r && y < r {
		dx := x - r
		dy := y - r
		return dx*dx+dy*dy > r*r
	}

	// Bottom-left
	if x < r && y >= h-r {
		dx := x - r
		dy := y - (h - r)
		return dx*dx+dy*dy > r*r
	}

	// Top-right
	if x >= w-r && y < r {
		dx := x - (w - r)
		dy := y - r
		return dx*dx+dy*dy > r*r
	}

	// Bottom-right
	if x >= w-r && y >= h-r {
		dx := x - (w - r)
		dy := y - (h - r)
		return dx*dx+dy*dy > r*r
	}

	return false
}

func (imgInfo *imageInfo) save(assets []Asset) error {
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
