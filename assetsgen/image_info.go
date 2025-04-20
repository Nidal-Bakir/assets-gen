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

type imageInfo struct {
	img              image.Image
	imagePath        string
	imageName        string
	imageExt         string
	encoder          imgio.Encoder
	asset            Asset
	genImageLocation func(screenType string) (directory string, imageName string)
}

type imageInfoSlice []imageInfo

func (s *imageInfoSlice) forEeach(fn func(imageInfo) imageInfo) *imageInfoSlice {
	for i, v := range *s {
		(*s)[i] = fn(v)
	}
	return s
}

func (s *imageInfoSlice) resizeForAssets() *imageInfoSlice {
	return s.forEeach(
		func(imgInfo imageInfo) imageInfo {
			return *imgInfo.resizeFroAsset()
		},
	)
}

func (s imageInfoSlice) save() error {
	for _, v := range s {
		v.save()
	}
	return nil
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

func (imgInfo *imageInfo) resize(w, h int) *imageInfo {
	imgInfo.img = transform.Resize(imgInfo.img, w, h, transform.Linear)
	return imgInfo
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

func (imgInfo imageInfo) splitPerAsset(assets []Asset) *imageInfoSlice {
	s := make(imageInfoSlice, len(assets))
	for i, a := range assets {
		s[i] = imgInfo
		s[i].asset = a
	}
	return &s
}

// make sure to set the asset before calling this function,
// it will panic if the asset is nil when it try to dereference
// the asset object functions
func (imgInfo *imageInfo) resizeFroAsset() *imageInfo {
	imageBounds := imgInfo.img.Bounds()
	w, h := imgInfo.asset.CalcSize(imageBounds.Dx(), imageBounds.Dy())
	return imgInfo.resize(w, h)
}

// make sure to set the asset before calling this function,
// it will panic if the asset is nil when it try to dereference
// the asset object functions
func (imgInfo imageInfo) save() error {
	dir, name := imgInfo.genImageLocation(imgInfo.asset.Name())
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = imgio.Save(filepath.Join(dir, name), imgInfo.img, imgInfo.encoder)
	if err != nil {
		return err
	}

	return nil
}
