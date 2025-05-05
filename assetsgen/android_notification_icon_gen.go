package assetsgen

import (
	"fmt"
	"image/color"
)

// MDPI    - 24px
// HDPI    - 36px
// XHDPI   - 48px
// XXHDPI  - 72px
// XXXHDPI - 96px
var androidNotificationIconDpis = []asset{
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

type AndroidNotificationIconOptions struct {
	// between [0..1] as percentage of how match the pixel should be transparent to keep its original color. Use -1 to disable
	AlphaThreshold float64

	FolderName AndroidFolderName

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool

	OutputFileName string
}

func GenerateNotificationIconForAndroid(imagePath string, option AndroidNotificationIconOptions) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, option.FolderName, intentNotificationIcon)
	if err != nil {
		return err
	}

	if len(option.OutputFileName) != 0 {
		option.OutputFileName = fmt.Sprint(option.OutputFileName, imgInfo.imageExt)
	}

	err = imgInfo.
		If(option.TrimWhiteSpace, imgInfo.TrimWhiteSpace).
		If(option.AlphaThreshold >= 0, func() *imageInfo { return imgInfo.RemoveAlphaOnThreshold(option.AlphaThreshold) }).
		ConvertNoneOpaqueToColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		SquareImageWithEmptyPixels().
		SplitPerAsset(androidNotificationIconDpis).
		ResizeForAssets().
		SaveWithCustomName(option.OutputFileName)

	if err != nil {
		return err
	}

	return nil
}

type androidNotificationIconDpiAsset struct {
	dpiName string
	width   int
	height  int
}

func (a androidNotificationIconDpiAsset) Name() string {
	return a.dpiName
}

func (a androidNotificationIconDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.width, a.height
}

func (a androidNotificationIconDpiAsset) CalcPadding(_, _ int) int {
	return 0
}
