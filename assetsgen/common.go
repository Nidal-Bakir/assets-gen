package assetsgen

import (
	"errors"
	"math"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type platformType string

const (
	platformTypeAndroid platformType = "android"
	platformTypeIos     platformType = "ios"

	rootFolderName string = "assets_gen_out"
)

type intention string

const (
	intentAppIcon          intention = "app_icon"
	intentNotificationIcon intention = "notification_icon"
	intentAsset            intention = "asset"
)

type Asset interface {
	Name() string
	CalcSize(w, h int) (int, int)
	CalcPadding(w, h int) int
}

type backgroundIcon interface {
	generateImgInfo(logo imageInfo) (imageInfo, error)
}

type gradientBackground struct {
	table        GradientTable
	degree       float64
	gradientType GradientType
}

func (g gradientBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.copy()

	switch g.gradientType {
	case LinearGradient:
		bgImage.linearGradient(g.table, g.degree)
	case RadialGradient:
		bgImage.radialGradient(g.table)
	}

	return *bgImage, nil
}

func NewLinearGradientBackground(table GradientTable, degree float64) backgroundIcon {
	return gradientBackground{table: table, degree: degree, gradientType: LinearGradient}
}

func NewRadialGradientBackground(table GradientTable) backgroundIcon {
	return gradientBackground{table: table, gradientType: RadialGradient}
}

type imageBackground struct {
	imagePath string

	// between [0..1] as percentage of the maximum axis (w,h) of the image;
	padding float64
}

func (i imageBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.copy()
	img, err := imgio.Open(i.imagePath)
	if err != nil {
		return imageInfo{}, err
	}
	bgImage.img = img

	bounds := bgImage.img.Bounds()

	pad := math.Max(float64(bounds.Dx()), float64(bounds.Dy())) * i.padding
	pad = math.Floor(pad)
	bgImage.squareImageWithPadding(int(pad))

	logoBounds := logo.img.Bounds()
	bgImage.resize(logoBounds.Dx(), logoBounds.Dy())

	return *bgImage, nil
}

// [padding] between [0..1] as percentage of the maximum axis (w,h) of the image
func NewImageBackground(imagePath string, padding float64) backgroundIcon {
	return imageBackground{imagePath: imagePath, padding: padding}
}

type solidColorBackground struct {
	color colorful.Color
}

func (s solidColorBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.copy()
	solidColorGradient := GradientTable{{Col: s.color, Pos: 1.0}}
	bgImage.linearGradient(solidColorGradient, 0)
	return *bgImage, nil
}

func NewSolidColorBackground(c colorful.Color) backgroundIcon {
	return solidColorBackground{c}
}
