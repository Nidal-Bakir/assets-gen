package assetsgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

// make sure to update it if you update the dpis slices below
const MAX_DPI_SIZE_FOR_ANDROID_APP_ICON = 432

// MDPI    - 108px
// HDPI    - 162px
// XHDPI   - 216px
// XXHDPI  - 324px
// XXXHDPI - 432px
var androidAdaptiveAppIconLayerDpisV26 = []asset{
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
		size:    MAX_DPI_SIZE_FOR_ANDROID_APP_ICON,
	},
}

// MDPI    - 66px
// HDPI    - 99px
// XHDPI   - 132px
// XXHDPI  - 198px
// XXXHDPI - 264px
var androidAdaptiveAppIconLogoDpisV26 = []asset{
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
var androidAppIconDpisLegacy = []asset{
	androidAppIconDpiAsset{
		dpiName: "mdpi",
		size:    48,
		padding: 4,
	},
	androidAppIconDpiAsset{
		dpiName: "hdpi",
		size:    72,
		padding: 6,
	},
	androidAppIconDpiAsset{
		dpiName: "xhdpi",
		size:    96,
		padding: 8,
	},
	androidAppIconDpiAsset{
		dpiName: "xxhdpi",
		size:    144,
		padding: 12,
	},
	androidAppIconDpiAsset{
		dpiName: "xxxhdpi",
		size:    192,
		padding: 24,
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

type AndroidAppIconOptions struct {
	// between [0..1] as percentage of the Radius. For example 1 would make the a full circle clip of the image, and 0 will do nothing, 0.5 will make rounded corners
	RoundedCornerPercentRadius float64

	// between [0..1] as percentage of how match the pixel should be transparent to keep its original color. Use -1 to disable
	AlphaThreshold float64

	BgIcon     BackgroundIcon
	FolderName AndroidFolderName

	// between [0..1] as percentage of the maximum axis (w,h) of the image
	Padding float64

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool

	MaskColor *colorful.Color

	OutputFileName string
}

func GenerateAppIconForAndroid(imagePath string, option AndroidAppIconOptions) error {
	logoImage, err := genImageInfoForAndroid(imagePath, option.FolderName, intentAppIcon)
	if err != nil {
		return err
	}

	pad := calPadding(logoImage.img, option.Padding)

	logoImage.
		If(option.TrimWhiteSpace, logoImage.TrimWhiteSpace).
		SquareImageWithEmptyPixels().
		ResizeSquare(MAX_DPI_SIZE_FOR_ANDROID_APP_ICON). // for performance optimization
		Padding(pad).
		ResizeSquare(MAX_DPI_SIZE_FOR_ANDROID_APP_ICON). // for performance optimization
		If(option.AlphaThreshold >= 0, func() *imageInfo { return logoImage.RemoveAlphaOnThreshold(option.AlphaThreshold) }).
		If(option.MaskColor != nil, func() *imageInfo { return logoImage.ConvertNoneOpaqueToColor(*option.MaskColor) })

	bgImage, err := option.BgIcon.generateImgInfo(logoImage)
	if err != nil {
		return err
	}

	w := sync.WaitGroup{}
	w.Add(2)

	var legacyAppIconError error
	var adaptiveAppIconError error

	go func() {
		defer w.Done()
		legacyAppIconError = generateLegacyAppIcon(logoImage, bgImage, option.RoundedCornerPercentRadius, option.AlphaThreshold, androidAppIconDpisLegacy, option.OutputFileName)
	}()

	go func() {
		defer w.Done()
		adaptiveAppIconError = generateAdaptiveAppIcon(logoImage, bgImage, androidAdaptiveAppIconLayerDpisV26, androidAdaptiveAppIconLogoDpisV26, option.OutputFileName)
	}()

	w.Wait()

	if legacyAppIconError != nil {
		return legacyAppIconError
	}
	if adaptiveAppIconError != nil {
		return adaptiveAppIconError
	}

	return nil
}

func generateLegacyAppIcon(logoImage imageInfo, bgImage imageInfo, roundedCornerPercentRadius float64, AlphaThreshold float64, androidAppIconDpisLegacy []asset, outputFileName string) error {
	if len(outputFileName) != 0 {
		outputFileName = fmt.Sprint(outputFileName, logoImage.imageExt)
	}

	err := bgImage.
		StackWithNoAlpha(AlphaThreshold, logoImage).
		ClipRRect(roundedCornerPercentRadius).
		SplitPerAsset(androidAppIconDpisLegacy).
		ResizeForAssets().
		PadForAsset().
		ResizeForAssets().
		SaveWithCustomName(outputFileName)

	if err != nil {
		return err
	}

	return nil
}

func generateAdaptiveAppIcon(logoImage imageInfo, bgImage imageInfo, androidAdaptiveAppIconLayerDpisV26 []asset, androidAdaptiveAppIconLogoDpisV26 []asset, outputFileName string) error {
	err := generateIcLauncherXml(logoImage, outputFileName)
	if err != nil {
		return err
	}

	logos := logoImage.
		SplitPerAsset(androidAdaptiveAppIconLogoDpisV26).
		ResizeForAssets().
		PadForAsset().
		SetAssets(androidAdaptiveAppIconLayerDpisV26).
		ResizeForAssets()

	shouldUseAssetName := len(outputFileName) == 0
	foregroundName := fmt.Sprint(outputFileName, "_foreground", logoImage.imageExt)
	monochromeName := fmt.Sprint(outputFileName, "_monochrome", logoImage.imageExt)
	backgroundName := fmt.Sprint(outputFileName, "_background", logoImage.imageExt)

	for _, logo := range *logos {
		if shouldUseAssetName {
			foregroundName = fmt.Sprint(logo.imgNameWithoutExt, "_foreground", logoImage.imageExt)
			monochromeName = fmt.Sprint(logo.imgNameWithoutExt, "_monochrome", logoImage.imageExt)
		}

		err := logo.SaveWithCustomName(foregroundName)
		if err != nil {
			return err
		}

		err = logo.SaveWithCustomName(monochromeName)
		if err != nil {
			return err
		}
	}

	bgs := bgImage.
		SplitPerAsset(androidAdaptiveAppIconLayerDpisV26).
		ResizeForAssets()

	for _, bg := range *bgs {
		if shouldUseAssetName {
			backgroundName = fmt.Sprint(bg.imgNameWithoutExt, "_background", bg.imageExt)
		}
		err := bg.SaveWithCustomName(backgroundName)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateIcLauncherXml(logoImage imageInfo, outputFileName string) error {
	sb := strings.Builder{}

	var name string
	if len(outputFileName) == 0 {
		name = logoImage.imgNameWithoutExt
	} else {
		name = outputFileName
	}

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
