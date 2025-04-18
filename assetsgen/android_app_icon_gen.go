package assetsgen

import (
	"image/color"
	"os"
	"path/filepath"
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

func GenerateAppIconForAndroid(imagePath string, folderName androidFolderName, padding int) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, folderName, intentAppIcon)
	if err != nil {
		return err
	}

	imgInfo.squareImageWithPadding(padding)

	err = generateLegacyAppIcon(imgInfo, androidAppIconDpisLegacy)
	if err != nil {
		return err
	}

	err = generateAdaptiveAppIcon(imgInfo, androidAdaptiveAppIconDpisV26)
	if err != nil {
		return err
	}

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

func generateLegacyAppIcon(imgInfo imageInfo, androidAppIconDpisLegacy []Asset) error {
	err := imgInfo.
		convertOpaqueToColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
		clipRRect(80).
		padding(275).
		save(androidAppIconDpisLegacy)

	if err != nil {
		return err
	}

	return nil
}

func generateAdaptiveAppIcon(imgInfo imageInfo, androidAdaptiveAppIconDpisV26 []Asset) error {
	err := generateIcLauncherXml(imgInfo)
	if err != nil {
		return err
	}

	return nil
}

func generateIcLauncherXml(imgInfo imageInfo) error {
	ic_launcher_xml := `<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
  <background android:drawable="@mipmap/ic_launcher_background"/>
  <foreground android:drawable="@mipmap/ic_launcher_foreground"/>
  <monochrome android:drawable="@mipmap/ic_launcher_monochrome"/>
</adaptive-icon>`

	dir, _ := imgInfo.genImageLocation("anydpi-v26")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, "ic_launcher.xml"))
	if err != nil {
		return err
	}

	_, err = file.WriteString(ic_launcher_xml)
	if err != nil {
		return err
	}

	return nil
}
