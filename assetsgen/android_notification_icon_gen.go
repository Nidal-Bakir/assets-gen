package assetsgen

import (
	"image/color"
)

// MDPI    - 24px
// HDPI    - 36px
// XHDPI   - 48px
// XXHDPI  - 72px
// XXXHDPI - 96px
var androidNotificationIconDpis = []Asset{
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

func GenerateNotificationIconForAndroid(imagePath string, folderName androidFolderName) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, folderName, intentNotificationIcon)
	if err != nil {
		return err
	}

	err = imgInfo.
		convertNoneOpaqueToColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		squareImageWithPadding(0).
		save(androidNotificationIconDpis)

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
