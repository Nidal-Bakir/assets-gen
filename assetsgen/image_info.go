package assetsgen

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

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

func (s *imageInfoSlice) ForEach(fn func(imageInfo) imageInfo) *imageInfoSlice {
	for i, v := range *s {
		(*s)[i] = fn(v)
	}
	return s
}

func (s *imageInfoSlice) ResizeForAssets() *imageInfoSlice {
	return s.ForEach(
		func(imgInfo imageInfo) imageInfo {
			return *imgInfo.ResizeForAsset()
		},
	)
}

func (s *imageInfoSlice) PadForAsset() *imageInfoSlice {
	return s.ForEach(
		func(imgInfo imageInfo) imageInfo {
			return *imgInfo.PadForAsset()
		},
	)
}

func (s imageInfoSlice) Save() error {
	for _, v := range s {
		v.Save()
	}
	return nil
}

func (s imageInfoSlice) SaveWithCustomName(customImageName string) error {
	for _, v := range s {
		v.SaveWithCustomName(customImageName)
	}
	return nil
}

func (s *imageInfoSlice) SetAssets(assets []asset) *imageInfoSlice {
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
		genImageLocation: func(screenType, customImageName string) (string, string) {
			dir := filepath.Join(
				rootFolderName,
				string(platform),
				string(intent),
				imgNameWithoutExt,
				lastFolderName(screenType),
			)
			name := imgName
			if len(customImageName) != 0 {
				name = customImageName
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

func (imgInfo *imageInfo) ResizeSquare(x int) *imageInfo {
	return imgInfo.Resize(x, x)
}

func (imgInfo *imageInfo) Resize(w, h int) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	imgW := imgBounds.Dx()
	imgH := imgBounds.Dy()
	if imgW == w && imgH == h { // it's already a resized to w,h
		return imgInfo
	}

	imgInfo.img = transform.Resize(imgInfo.img, w, h, transform.Linear)
	return imgInfo
}

func (imgInfo *imageInfo) SquareImageWithEmptyPixels() *imageInfo {
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

func (imgInfo *imageInfo) Padding(padding int) *imageInfo {
	if padding == 0 {
		return imgInfo
	}

	imgInfo.img = clone.Pad(imgInfo.img, padding, padding, clone.NoFill)
	return imgInfo
}

func (imgInfo *imageInfo) ConvertColors(fn func(color.Color) color.Color) *imageInfo {

	imgInfo.img = adjust.Apply(
		imgInfo.img,
		func(pxColor color.RGBA) color.RGBA {
			return color.RGBAModel.Convert(fn(pxColor)).(color.RGBA)
		},
	)

	return imgInfo
}

func (imgInfo *imageInfo) ConvertNoneOpaqueToColor(newColor color.Color) *imageInfo {
	return imgInfo.ConvertColors(func(pxColor color.Color) color.Color {
		c := color.RGBAModel.Convert(pxColor).(color.RGBA)
		if c.A == 0 {
			return pxColor
		}
		return newColor
	})
}

func (imgInfo *imageInfo) UpdatePixels(updater func(x, y int, c color.Color) color.Color) *imageInfo {
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

func (imgInfo *imageInfo) ClipRRect(percentRadius float64) *imageInfo {
	if percentRadius == 0 {
		return imgInfo
	}
	if percentRadius == 1 {
		return imgInfo.ClipToCircle()
	}

	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	r := math.Max(float64(imgBounds.Dx()), float64(imgBounds.Dy())) / 2
	roundedCornerRadius := int(math.Floor(r) * percentRadius)

	return imgInfo.UpdatePixels(
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

func (imgInfo *imageInfo) ClipToCircle() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := float64(imgBounds.Dx())
	h := float64(imgBounds.Dy())
	r := math.Max(w, h) / 2
	cx := w / 2
	cy := h / 2

	return imgInfo.UpdatePixels(
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

func (imgInfo imageInfo) SplitPerAsset(assets []asset) *imageInfoSlice {
	s := make(imageInfoSlice, len(assets))
	for i, a := range assets {
		s[i] = imgInfo
		s[i].asset = a
	}
	return &s
}

func (imgInfo *imageInfo) ResizeForAsset() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w, h := imgInfo.asset.CalcSize(imgBounds.Dx(), imgBounds.Dy())
	return imgInfo.Resize(w, h)
}

func (imgInfo *imageInfo) PadForAsset() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w, h := imgInfo.asset.CalcSize(imgBounds.Dx(), imgBounds.Dy())
	return imgInfo.Padding(imgInfo.asset.CalcPadding(w, h))
}

func (imgInfo imageInfo) Save() error {
	return imgInfo.SaveWithCustomName("")
}

func (imgInfo imageInfo) SaveWithCustomName(customImageName string) error {
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

func (imgInfo *imageInfo) LinearGradient(colorsTable GradientTable, degree int) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	imgInfo.img = createLinearGradient(colorsTable, degree, imgBounds.Dx(), imgBounds.Dy())
	return imgInfo
}

func (imgInfo *imageInfo) RadialGradient(colorsTable GradientTable) *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	imgInfo.img = createRadialGradient(colorsTable, imgBounds.Dx(), imgBounds.Dy())
	return imgInfo
}

func (imgInfo imageInfo) Copy() *imageInfo {
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

// All the images should be the same width and height
//
// layout the images on top of each other the last image will be laid out at last
func (imgInfo *imageInfo) Stack(images ...imageInfo) *imageInfo {
	images = append([]imageInfo{*imgInfo}, images...)
	slices.Reverse(images)

	return imgInfo.UpdatePixels(
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

func (imgInfo *imageInfo) StackWithNoAlpha(threshold float64, images ...imageInfo) *imageInfo {
	images = append([]imageInfo{*imgInfo}, images...)
	slices.Reverse(images)

	l := len(images)
	return imgInfo.UpdatePixels(
		func(x, y int, _ color.Color) color.Color {

			for i, img := range images {
				isLastImage := i == l-1
				rgba := color.RGBAModel.Convert(img.img.At(x, y)).(color.RGBA)

				if isColorNotTransparent(rgba) {
					return rgba
				}
				if isColorFullTransparent(rgba) || isLastImage {
					if isLastImage {
						rgba.A = 255
						return rgba
					}
					continue
				}

				// not last image and the alpha value is  0 < A < 255

				if float64(rgba.A) > 255*threshold { // if the alpha is grater then threshold% then remove it from the color, or use the background color otherwise
					rgba.A = 255
					return rgba
				}

				continue
			}

			return color.RGBA{R: 255, G: 0, B: 0, A: 255} // red to catch any errors
		},
	)
}

// can not be 0
func (imgInfo *imageInfo) RemoveAlphaOnThreshold(threshold float64) *imageInfo {
	return imgInfo.UpdatePixels(
		func(x, y int, c color.Color) color.Color {
			rgba := color.RGBAModel.Convert(c).(color.RGBA)

			if isColorFullTransparent(rgba) || isColorNotTransparent(rgba) {
				return rgba
			}

			// the alpha value is  0 < A < 255
			if float64(rgba.A) > 255*threshold {
				rgba.A = 255
				return rgba
			}

			rgba.A = 0
			return rgba
		},
	)
}

func isColorFullTransparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

func isColorNotTransparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 255

}

func (imgInfo *imageInfo) RemoveAlpha() *imageInfo {
	return imgInfo.UpdatePixels(
		func(x, y int, c color.Color) color.Color {
			rgba := color.RGBAModel.Convert(c).(color.RGBA)
			rgba.A = 255
			return rgba
		},
	)
}

func (imgInfo *imageInfo) TrimWhiteSpace() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var topTrimCount, bottomTrimCount, leftTrimCount, rightTrimCount int
	go func() {
		defer wg.Done()
		topTrimCount, bottomTrimCount = reportX(imgInfo.img)
	}()
	go func() {
		defer wg.Done()
		leftTrimCount, rightTrimCount = reportY(imgInfo.img)
	}()

	wg.Wait()

	dst := image.NewRGBA(image.Rect(0, 0, w-rightTrimCount-leftTrimCount, h-bottomTrimCount-topTrimCount))
	draw.Draw(dst, dst.Rect, imgInfo.img, image.Point{leftTrimCount, topTrimCount}, draw.Src)

	imgInfo.img = dst
	return imgInfo
}

func reportX(img image.Image) (top, bottom int) {
	imgBounds := img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	canTrimX := func(y int) bool {
		for x := range w {
			if isColorFullTransparent(img.At(x, y)) {
				continue
			}
			return false
		}
		return true
	}

	var didTopStoped, didBottomStoped bool

	half := int(math.Floor(float64(h) / 2))
	y := 0
	for {
		if (didTopStoped && didBottomStoped) || (y >= half) {
			break
		}

		if canTrimX(y) {
			top++
		} else {
			didTopStoped = true
		}

		if canTrimX(h - y - 1) {
			bottom++
		} else {
			didBottomStoped = true
		}

		y++
	}

	return top, bottom
}

func reportY(img image.Image) (left, right int) {
	imgBounds := img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()

	canTrimY := func(x int) bool {
		for y := range h {
			if isColorFullTransparent(img.At(x, y)) {
				continue
			}
			return false
		}
		return true
	}

	var didLeftStoped, didRightStoped bool

	half := int(math.Floor(float64(w) / 2))
	x := 0
	for {
		if (didLeftStoped && didRightStoped) || (x >= half) {
			break
		}

		if canTrimY(x) {
			left++
		} else {
			didLeftStoped = true
		}

		if canTrimY(w - x - 1) {
			right++
		} else {
			didRightStoped = true
		}

		x++
	}

	return left, right
}

func (imgInfo *imageInfo) CropToSquare() *imageInfo {
	imgBounds := imgInfo.img.Bounds()
	w := imgBounds.Dx()
	h := imgBounds.Dy()
	if w == h {
		return imgInfo
	}

	cx := w / 2
	cy := h / 2

	size := int(math.Min(float64(w), float64(h)))

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(dst, dst.Rect, imgInfo.img, image.Point{cx - size/2, cy - size/2}, draw.Src)

	imgInfo.img = dst
	return imgInfo
}
