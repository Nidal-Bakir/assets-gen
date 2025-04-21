package assetsgen

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// MDPI    - 108px
// HDPI    - 162px
// XHDPI   - 216px
// XXHDPI  - 324px
// XXXHDPI - 432px
var androidAdaptiveAppIconLayerDpisV26 = []Asset{
	androidAppIconDpiAsset{
		dpiName: "mdpi",
		size:    108,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		size:    162,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		size:    216,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		size:    324,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		size:    432,
	},
}

// MDPI    - 66px
// HDPI    - 99px
// XHDPI   - 132px
// XXHDPI  - 198px
// XXXHDPI - 264px
var androidAdaptiveAppIconLogoDpisV26 = []Asset{
	androidAppIconDpiAsset{
		dpiName: "mdpi",
		size:    66,
		padding: 24,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		size:    99,
		padding: 36,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		size:    132,
		padding: 48,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		size:    198,
		padding: 72,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		size:    264,
		padding: 96,
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
		size:    48,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		size:    72,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		size:    96,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		size:    144,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		size:    192,
	},
}

type androidAppIconDpiAsset struct {
	dpiName string
	size    int
	padding int
}

func (a androidAppIconDpiAsset) Name() string {
	return a.dpiName
}

func (a androidAppIconDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.size, a.size
}
func (a androidAppIconDpiAsset) CalcPadding(_, _ int) int {
	return a.padding
}

type AppIconOptions struct {
	RoundedCornerRadius int
	BgColor             GradientTable
	Degree              float64
	GradientType        GradientType
	FolderName          androidFolderName
	// between [0..1] as percentage of the maximum axis (w,h) of the image
	Padding float64
}

func GenerateAppIconForAndroid(imagePath string, option AppIconOptions) error {
	logoImage, err := genImageInfoForAndroid(imagePath, option.FolderName, intentAppIcon)
	if err != nil {
		return err
	}

	bounds := logoImage.img.Bounds()
	pad := math.Max(float64(bounds.Dx()), float64(bounds.Dy())) * option.Padding
	pad = math.Floor(pad)
	logoImage.squareImageWithPadding(int(pad))

	bgImage := generateBackgroudAppImage(logoImage, option)

	err = generateLegacyAppIcon(logoImage, bgImage, option.RoundedCornerRadius, androidAppIconDpisLegacy)
	if err != nil {
		return err
	}

	err = generateAdaptiveAppIcon(logoImage, bgImage, androidAdaptiveAppIconLayerDpisV26, androidAdaptiveAppIconLogoDpisV26)
	if err != nil {
		return err
	}

	return nil
}

func generateBackgroudAppImage(logoImage imageInfo, option AppIconOptions) imageInfo {
	bgImage := logoImage.copy()

	switch option.GradientType {
	case LinearGradient:
		bgImage.linearGradient(option.BgColor, option.Degree)
	case RadialGradient:
		bgImage.radialGradient(option.BgColor)
	}

	return *bgImage
}

func generateLegacyAppIcon(logoImage imageInfo, bgImage imageInfo, roundedCornerRadius int, androidAppIconDpisLegacy []Asset) error {
	err := bgImage.
		stack(logoImage).
		clipRRect(roundedCornerRadius).
		padding(150). // TODO: check this value
		splitPerAsset(androidAppIconDpisLegacy).
		resizeForAssets().
		save()

	if err != nil {
		return err
	}

	return nil
}

func generateAdaptiveAppIcon(logoImage imageInfo, bgImage imageInfo, androidAdaptiveAppIconLayerDpisV26 []Asset, androidAdaptiveAppIconLogoDpisV26 []Asset) error {
	err := generateIcLauncherXml(logoImage)
	if err != nil {
		return err
	}

	logos := logoImage.
		splitPerAsset(androidAdaptiveAppIconLogoDpisV26).
		resizeForAssets().
		setAssets(androidAdaptiveAppIconLayerDpisV26).
		resizeForAssets()

	for _, logo := range *logos {
		err := logo.saveWithCustomName(fmt.Sprint(logo.imgNameWithoutExt, "_foreground", logo.imageExt))
		if err != nil {
			return err
		}

		err = logo.saveWithCustomName(fmt.Sprint(logo.imgNameWithoutExt, "_monochrome", logo.imageExt))
		if err != nil {
			return err
		}
	}

	bgs := bgImage.
		splitPerAsset(androidAdaptiveAppIconLayerDpisV26).
		resizeForAssets()

	for _, bg := range *bgs {
		err := bg.saveWithCustomName(fmt.Sprint(bg.imgNameWithoutExt, "_background", bg.imageExt))
		if err != nil {
			return err
		}
	}

	return nil
}

func generateIcLauncherXml(logoImage imageInfo) error {
	sb := strings.Builder{}

	name := logoImage.imgNameWithoutExt

	sb.WriteString(`<?xml version="1.0" encoding="utf-8" ?>`)
	sb.WriteRune('\n')

	sb.WriteString(`<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">`)
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <background android:drawable="@mipmap/`, name, `_background" />`))
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <foreground android:drawable="@mipmap/`, name, `_foreground" />`))
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <monochrome android:drawable="@mipmap/`, name, `_monochrome" />`))
	sb.WriteRune('\n')

	sb.WriteString(`</adaptive-icon>`)
	sb.WriteRune('\n')

	ic_launcher_xml := sb.String()

	dir, name := logoImage.genImageLocation("anydpi-v26", "ic_launcher.xml")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}

	_, err = file.WriteString(ic_launcher_xml)
	if err != nil {
		return err
	}

	return nil
}
