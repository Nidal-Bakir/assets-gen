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

func GenerateNotificationIconForAndroid(imagePath string, folderName AndroidFolderName, outputFileName string) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, folderName, intentNotificationIcon)
	if err != nil {
		return err
	}

	if len(outputFileName) != 0 {
		outputFileName = fmt.Sprint(outputFileName, imgInfo.imageExt)
	}

	err = imgInfo.
		convertNoneOpaqueToColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		squareImageEmptyPixel().
		splitPerAsset(androidNotificationIconDpis).
		resizeForAssets().
		saveWithCustomName(outputFileName)

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
