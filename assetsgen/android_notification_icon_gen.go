package assetsgen

import (
	"fmt"
	"image/color"
	"path/filepath"
)

func androidNotificationIconDpis(androidFolderName string) []asset {
	// MDPI    - 24px
	// HDPI    - 36px
	// XHDPI   - 48px
	// XXHDPI  - 72px
	// XXXHDPI - 96px
	var dpis = []asset{
		androidNotificationIconDpiAsset{
			dpiName: "mdpi",
			width:   24,
			height:  24,
		},
		androidNotificationIconDpiAsset{
			dpiName: "hdpi",
			width:   36,
			height:  36,
		},
		androidNotificationIconDpiAsset{
			dpiName: "xhdpi",
			width:   48,
			height:  48,
		},
		androidNotificationIconDpiAsset{
			dpiName: "xxhdpi",
			width:   72,
			height:  72,
		},
		androidNotificationIconDpiAsset{
			dpiName: "xxxhdpi",
			width:   96,
			height:  96,
		},
	}

	for i, v := range dpis {
		dpi := v.(androidNotificationIconDpiAsset)
		dpi.dirName = fmt.Sprint(dpi.dirName, androidFolderName)
		dpis[i] = dpi
	}
	return dpis
}

type AndroidNotificationIconOptions struct {
	// between [0..1] as percentage of how match the pixel should be transparent to keep its original color. Use -1 to disable
	AlphaThreshold float64

	FolderName AndroidFolderName

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool

	OutputFileName string
}

func GenerateNotificationIconForAndroid(imagePath string, option AndroidNotificationIconOptions) error {
	logoImage, err := newImageInfo(
		imagePath,
		filepath.Join(PlatformTypeAndroid, "res"),
	)
	if err != nil {
		return err
	}
	defer logoImage.rootDir.Close()

	err = logoImage.
		If(option.TrimWhiteSpace, logoImage.TrimWhiteSpace).
		If(option.AlphaThreshold >= 0, func() *imageInfo { return logoImage.RemoveAlphaOnThreshold(option.AlphaThreshold) }).
		ConvertNoneOpaqueToColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		SquareImageWithEmptyPixels(0).
		SplitPerAsset(androidNotificationIconDpis(string(option.FolderName))).
		ResizeForAssets().
		SaveWithCustomName(option.OutputFileName)

	if err != nil {
		return err
	}

	return nil
}

type androidNotificationIconDpiAsset struct {
	dpiName string
	dirName string
	width   int
	height  int
}

func (a androidNotificationIconDpiAsset) Name() string {
	return a.dpiName
}

func (a androidNotificationIconDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.width, a.height
}

func (a androidNotificationIconDpiAsset) DirName() string {
	return fmt.Sprint(a.dirName, "-", a.dpiName)
}
