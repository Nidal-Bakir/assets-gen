package assetsgen

import (
	"path/filepath"

	"github.com/lucasb-eyer/go-colorful"
)

type AndroidGooglePlayLogoOptions struct {
	// between [0..1] as percentage of how match the pixel should be transparent to keep its original color. Use -1 to disable
	AlphaThreshold float64

	BgIcon BackgroundIcon

	// between [0..1] as percentage of the maximum axis (w,h) of the image
	Padding float64

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool

	MaskColor *colorful.Color

	OutputFileName string
}

func GenerateAndroidGooglePlayLogo(imagePath string, option AndroidGooglePlayLogoOptions) error {
	logoImage, err := newImageInfo(
		imagePath,
		filepath.Join(PlatformTypeAndroid, "main"),
	)
	if err != nil {
		return err
	}
	defer logoImage.rootDir.Close()

	pad := calPadding(logoImage.img, option.Padding)

	logoImage.
		If(option.TrimWhiteSpace, logoImage.TrimWhiteSpace).
		SquareImageWithEmptyPixels(pad).
		If(option.MaskColor != nil, func() *imageInfo { return logoImage.ConvertNoneOpaqueToColor(*option.MaskColor) })

	bgImage, err := option.BgIcon.generateImgInfo(logoImage)
	if err != nil {
		return err
	}

	bgImage.asset = androidGooglePlayLogoDpiAsset{
		dpiName: "main",
		Size:    512,
	}

	err = bgImage.
		StackWithNoAlpha(option.AlphaThreshold, logoImage).
		ResizeForAsset().
		SaveWithCustomName(option.OutputFileName)

	if err != nil {
		return err
	}

	return nil
}

type androidGooglePlayLogoDpiAsset struct {
	dpiName string
	Size    int
}

func (a androidGooglePlayLogoDpiAsset) Name() string {
	return a.dpiName
}

func (a androidGooglePlayLogoDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.Size, a.Size
}

func (a androidGooglePlayLogoDpiAsset) DirName() string {
	return ""
}
