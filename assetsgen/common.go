package assetsgen

import (
	"errors"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileNotFound        = errors.New("file not found")
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

type asset interface {
	Name() string
	CalcSize(w, h int) (int, int)
	CalcPadding(w, h int) int
}

type BackgroundIcon interface {
	generateImgInfo(logo imageInfo) (imageInfo, error)
}

type gradientBackground struct {
	table        GradientTable
	degree       int
	gradientType GradientType
}

func (g gradientBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.Copy()

	switch g.gradientType {
	case LinearGradient:
		bgImage.LinearGradient(g.table, g.degree)
	case RadialGradient:
		bgImage.RadialGradient(g.table)
	}

	return *bgImage, nil
}

func NewLinearGradientBackground(table GradientTable, degree int) BackgroundIcon {
	return gradientBackground{table: table, degree: degree, gradientType: LinearGradient}
}

func NewRadialGradientBackground(table GradientTable) BackgroundIcon {
	return gradientBackground{table: table, gradientType: RadialGradient}
}

type imageBackground struct {
	imagePath string
}

func (i imageBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.Copy()
	img, err := imgio.Open(i.imagePath)
	if err != nil {
		return imageInfo{}, err
	}
	bgImage.img = img

	logoBounds := logo.img.Bounds()

	bgImage.SquareImageEmptyPixel().
		Resize(logoBounds.Dx(), logoBounds.Dy()).
		RemoveAlpha()

	return *bgImage, nil
}

// [padding] between [0..1] as percentage of the maximum axis (w,h) of the image
func NewImageBackground(imagePath string) BackgroundIcon {
	return imageBackground{imagePath: imagePath}
}

type solidColorBackground struct {
	color colorful.Color
}

func (s solidColorBackground) generateImgInfo(logo imageInfo) (imageInfo, error) {
	bgImage := logo.Copy()
	solidColorGradient := GradientTable{{Col: s.color, Pos: 1.0}}
	bgImage.LinearGradient(solidColorGradient, 0)
	return *bgImage, nil
}

func NewSolidColorBackground(c colorful.Color) BackgroundIcon {
	return solidColorBackground{c}
}

func IsFileExistsAndImage(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileNotFound
		}
		return err
	}
	name := info.Name()

	if strings.HasSuffix(name, ".png") ||
		strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg") {
		return nil
	}

	return ErrUnsupportedFileType
}
