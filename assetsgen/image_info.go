package assetsgen

import (
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type imageInfo struct {
	img               image.Image
	imagePath         string
	imageName         string
	imgNameWithoutExt string
	imageExt          string
	encoder           imgio.Encoder
	asset             asset
	genImageLocation  func(screenType, customImageName string) (directory string, imageName string)
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
			return *imgInfo.resizeForAsset()
		},
	)
}

func (s *imageInfoSlice) padForAsset() *imageInfoSlice {
	return s.forEeach(
		func(imgInfo imageInfo) imageInfo {
			return *imgInfo.padForAsset()
		},
	)
}

func (s imageInfoSlice) save() error {
	for _, v := range s {
		v.save()
	}
	return nil
}

func (s imageInfoSlice) saveWithCustomName(customImageName string) error {
	for _, v := range s {
		v.saveWithCustomName(customImageName)
	}
	return nil
}

func (s *imageInfoSlice) setAssets(assets []asset) *imageInfoSlice {
	for i, v := range *s {
		v.asset = assets[i]
		(*s)[i] = v
	}
	return s
}

func newImageInfo(imagePath string, platform platformType, intent intention, lastFolderName func(screenType string) string) (imageInfo, error) {
	if err := IsFileExistsAndImage(imagePath); err != nil {
		return imageInfo{}, err
	}

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
		img:               img,
		encoder:           enc,
		imagePath:         imagePath,
		imageName:         imgName,
		imageExt:          imageExt,
		imgNameWithoutExt: imgNameWithoutExt,
		genImageLocation: func(screenType, cutomImageName string) (string, string) {
			dir := filepath.Join(
				rootFolderName,
				string(platform),
				string(intent),
				imgNameWithoutExt,
				lastFolderName(screenType),
			)
			name := imgName
			if len(cutomImageName) != 0 {
				name = cutomImageName
			}
			return dir, name
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

func (imgInfo *imageInfo) squareImageEmptyPixel() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	if w == h { // it's already a square
		return imgInfo
	}

	var padX, padY int
	if w < h {
		padX = (h - w) / 2
	} else {
		padY = (w - h) / 2
	}

	imgInfo.img = clone.Pad(imgInfo.img, padX, padY, clone.NoFill)
	return imgInfo
}

func (imgInfo *imageInfo) padding(padding int) *imageInfo {
	if padding == 0 {
		return imgInfo
	}

	imgInfo.img = clone.Pad(imgInfo.img, padding, padding, clone.NoFill)
	return imgInfo
}

func (imgInfo *imageInfo) convertColors(fn func(color.RGBA) color.RGBA) *imageInfo {
	imgInfo.img = adjust.Apply(imgInfo.img, func(pxColor color.RGBA) color.RGBA {
		return fn(pxColor)
	})
	return imgInfo
}

func (imgInfo *imageInfo) convertNoneOpaqueToColor(newColor color.RGBA) *imageInfo {
	return imgInfo.convertColors(func(pxColor color.RGBA) color.RGBA {
		if pxColor.A == 0 {
			return pxColor
		}
		return newColor
	})
}

func (imgInfo *imageInfo) updatePixels(updater func(x, y int, c color.Color) color.Color) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			dst.Set(x, y, updater(x, y, imgInfo.img.At(x, y)))
		}
	}

	imgInfo.img = dst
	return imgInfo
}

func (imgInfo *imageInfo) clipRRect(percentRadius float64) *imageInfo {
	if percentRadius == 1 {
		return imgInfo.clipToCircle()
	}

	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	r := math.Max(float64(imgBounds.Dx()), float64(imgBounds.Dy())) / 2
	roundedCornerRadius := int(math.Floor(r) * percentRadius)

	return imgInfo.updatePixels(
		func(x, y int, c color.Color) color.Color {
			if isOnRoundedCorner(x, y, w, h, roundedCornerRadius) {
				return color.RGBA{}
			}
			return c
		},
	)
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

func (imgInfo *imageInfo) clipToCircle() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := float64(imgBounds.Dx())
	h := float64(imgBounds.Dy())
	r := math.Max(w, h) / 2
	cx := w / 2
	cy := h / 2

	return imgInfo.updatePixels(
		func(x, y int, c color.Color) color.Color {
			if isPixelInsideCircle(x, y, cx, cy, r) {
				return c
			}
			return color.RGBA{}
		},
	)
}

func isPixelInsideCircle(x, y int, cx, cy, radius float64) bool {
	distance := math.Sqrt(math.Pow(float64(x)-cx, 2) + math.Pow(float64(y)-cy, 2))
	return distance <= radius
}

func (imgInfo imageInfo) splitPerAsset(assets []asset) *imageInfoSlice {
	s := make(imageInfoSlice, len(assets))
	for i, a := range assets {
		s[i] = imgInfo
		s[i].asset = a
	}
	return &s
}

func (imgInfo *imageInfo) resizeForAsset() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w, h := imgInfo.asset.CalcSize(imgBounds.Dx(), imgBounds.Dy())
	return imgInfo.resize(w, h)
}

func (imgInfo *imageInfo) padForAsset() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w, h := imgInfo.asset.CalcSize(imgBounds.Dx(), imgBounds.Dy())
	return imgInfo.padding(imgInfo.asset.CalcPadding(w, h))
}

func (imgInfo imageInfo) save() error {
	return imgInfo.saveWithCustomName("")
}

func (imgInfo imageInfo) saveWithCustomName(customImageName string) error {
	dir, name := imgInfo.genImageLocation(imgInfo.asset.Name(), customImageName)
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

func (imgInfo *imageInfo) linearGradient(colorsTable GradientTable, degree int) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	imgInfo.img = createLinearGradient(colorsTable, degree, imgBounds.Dx(), imgBounds.Dy())
	return imgInfo
}

func (imgInfo *imageInfo) radialGradient(colorsTable GradientTable) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	imgInfo.img = createRadialGradient(colorsTable, imgBounds.Dx(), imgBounds.Dy())
	return imgInfo
}

func (imgInfo imageInfo) copy() *imageInfo {
	return &imageInfo{
		img:               clone.AsRGBA(imgInfo.img),
		imagePath:         imgInfo.imagePath,
		imageName:         imgInfo.imageName,
		imageExt:          imgInfo.imageExt,
		imgNameWithoutExt: imgInfo.imgNameWithoutExt,
		encoder:           imgInfo.encoder,
		asset:             imgInfo.asset,
		genImageLocation:  imgInfo.genImageLocation,
	}
}

// All the imgs should be the same width and height
//
// layout the imgs on top of each other the last image will be laid out at last
func (imgInfo *imageInfo) stack(images ...imageInfo) *imageInfo {
	images = append([]imageInfo{*imgInfo}, images...)
	slices.Reverse(images)

	return imgInfo.updatePixels(
		func(x, y int, _ color.Color) color.Color {
			for _, img := range images {
				c := img.img.At(x, y)
				_, _, _, a := c.RGBA()
				if a != 0 {
					return c
				}
			}
			return color.RGBA{}
		},
	)
}
