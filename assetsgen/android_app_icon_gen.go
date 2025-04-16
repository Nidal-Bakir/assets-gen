package assetsgen

import (
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
)

// MDPI    - 108px
// HDPI    - 162px
// XHDPI   - 216px
// XXHDPI  - 324px
// XXXHDPI - 432px
var androidAdaptiveAppIconDpisV26 = []Asset{
	androidAppIconDpiAsset{
		dpiName: "mdpi",
		width:   108,
		height:  108,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		width:   162,
		height:  162,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		width:   216,
		height:  216,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		width:   324,
		height:  324,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		width:   432,
		height:  432,
	},
}

// MDPI    - 48px
// HDPI    - 72px
// XHDPI   - 96px
// XXHDPI  - 144px
// XXXHDPI - 192px
var androidAppIconDpisLegacy = []Asset{
	androidAppIconDpiAsset{
		dpiName: "mdpi",
		width:   48,
		height:  48,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		width:   72,
		height:  72,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		width:   96,
		height:  96,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		width:   144,
		height:  144,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		width:   192,
		height:  192,
	},
}

func GenerateAppIconForAndroid(imagePath string, folderName androidFolderName) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, folderName, intentAppIcon)
	if err != nil {
		return err
	}

	dir, name := imgInfo.genImageLocation("anydpi-v26")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	imgInfo.img = squareImageWithPadding(imgInfo.img, 0)

	err = imgio.Save(filepath.Join(dir, name), imgInfo.img, imgInfo.encoder)
	if err != nil {
		return err
	}

	// err = generateImageAsstes(imgInfo, androidAppIconDpisLegacy)
	// if err != nil {
	// 	return err
	// }

	return nil
}

type androidAppIconDpiAsset struct {
	dpiName string
	width   int
	height  int
}

func (a androidAppIconDpiAsset) Name() string {
	return a.dpiName
}

func (a androidAppIconDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.width, a.height
}
