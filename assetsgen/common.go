package assetsgen

import (
	"errors"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileNotFound        = errors.New("file not found")
)

type platformType = string

const (
	PlatformTypeAndroid platformType = "android"
	PlatformTypeIos     platformType = "ios"

	RootFolderName string = "assets_gen_out"
)

type asset interface {
	Name() string
	CalcSize(w, h int) (int, int)
	DirName() string
}

type BackgroundIcon interface {
	generateImgInfo(logo *imageInfo) (*imageInfo, error)
}

type gradientBackground struct {
	table        GradientTable
	degree       int
	gradientType GradientType
}

func (g gradientBackground) generateImgInfo(logo *imageInfo) (*imageInfo, error) {
	bgImage := logo.Copy()

	switch g.gradientType {
	case LinearGradient:
		bgImage.LinearGradient(g.table, g.degree)
	case RadialGradient:
		bgImage.RadialGradient(g.table)
	}

	return bgImage, nil
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

func (i imageBackground) generateImgInfo(logo *imageInfo) (*imageInfo, error) {
	bgImage := logo.Copy()
	img, err := imgio.Open(i.imagePath)
	if err != nil {
		return &imageInfo{}, err
	}
	bgImage.img = img

	logoBounds := logo.img.Bounds()

	bgImage.CropToSquare().
		Resize(logoBounds.Dx(), logoBounds.Dy()).
		RemoveAlpha()

	return bgImage, nil
}

// [padding] between [0..1] as percentage of the maximum axis (w,h) of the image
func NewImageBackground(imagePath string) BackgroundIcon {
	return imageBackground{imagePath: imagePath}
}

type solidColorBackground struct {
	color colorful.Color
}

func (s solidColorBackground) generateImgInfo(logo *imageInfo) (*imageInfo, error) {
	return logo.Copy().SoldiColor(s.color), nil
}

func NewSolidColorBackground(c colorful.Color) BackgroundIcon {
	return solidColorBackground{c}
}

func IsFileExistsAndImage(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
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

func calPadding(img image.Image, padding float64) int {
	bounds := img.Bounds()
	pad := math.Max(float64(bounds.Dx()), float64(bounds.Dy())) * padding
	pad = math.Floor(pad)
	return int(pad)
}

func GetRootDir() (*os.Root, error) {
	p := filepath.Join("./", RootFolderName)
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return nil, err
	}
	rootDir, err := os.OpenRoot(p)
	if err != nil {
		return nil, err
	}

	return rootDir, nil
}

func saveImage(root *os.Root, filename string, img image.Image, encoder imgio.Encoder) error {
	f, err := root.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return encoder(f, img)
}

func splitPath(path string) []string {
	dir, last := filepath.Split(path)
	if dir == "" {
		return []string{last}
	}
	return append(splitPath(filepath.Clean(dir)), last)
}
