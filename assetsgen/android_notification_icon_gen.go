package assetsgen

import (
	"image/color"

	"github.com/anthonynsimon/bild/adjust"
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

	// convert any non opaque pixel to the color white
	imgInfo.img = adjust.Apply(imgInfo.img, func(pxColor color.RGBA) color.RGBA {
		if pxColor.A == 0 {
			return pxColor
		}
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	})

	imgInfo.img = squareImageWithPadding(imgInfo.img)

	err = generateImageAsstes(imgInfo, androidNotificationIconDpis)
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
