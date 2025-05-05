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

func androidAdaptiveAppIconLayerDpisV26(androidFolderName string) []asset {
	// MDPI    - 108px
	// HDPI    - 162px
	// XHDPI   - 216px
	// XXHDPI  - 324px
	// XXXHDPI - 432px
	var dpis = []asset{
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
	for i, v := range dpis {
		dpi := v.(androidAppIconDpiAsset)
		dpi.dirName = fmt.Sprint(dpi.dirName, androidFolderName)
		dpis[i] = dpi
	}
	return dpis
}

func androidAdaptiveAppIconLogoDpisV26(androidFolderName string) []asset {
	// MDPI    - 66px
	// HDPI    - 99px
	// XHDPI   - 132px
	// XXHDPI  - 198px
	// XXXHDPI - 264px
	var dpis = []asset{
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
	for i, v := range dpis {
		dpi := v.(androidAppIconDpiAsset)
		dpi.dirName = fmt.Sprint(dpi.dirName, androidFolderName)
		dpis[i] = dpi
	}
	return dpis
}

func androidAppIconDpisLegacy(androidFolderName string) []asset {
	// MDPI    - 48px
	// HDPI    - 72px
	// XHDPI   - 96px
	// XXHDPI  - 144px
	// XXXHDPI - 192px
	var dpis = []asset{
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

	for i, v := range dpis {
		dpi := v.(androidAppIconDpiAsset)
		dpi.dirName = fmt.Sprint(dpi.dirName, androidFolderName)
		dpis[i] = dpi
	}
	return dpis
}

type androidAppIconDpiAsset struct {
	dpiName string
	dirName string
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

func (a androidAppIconDpiAsset) DirName() string {
	return fmt.Sprint(a.dirName, "-", a.dpiName)
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
	logoImage, err := newImageInfo(
		imagePath,
		filepath.Join(PlatformTypeAndroid, "res"),
	)
	if err != nil {
		return err
	}
	defer logoImage.rootDir.Close()

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
		legacyAppIconError = generateLegacyAppIcon(
			logoImage,
			bgImage,
			option.RoundedCornerPercentRadius,
			option.AlphaThreshold,
			androidAppIconDpisLegacy(string(option.FolderName)),
			option.OutputFileName,
		)
	}()

	go func() {
		defer w.Done()
		adaptiveAppIconError = generateAdaptiveAppIcon(
			logoImage,
			bgImage,
			androidAdaptiveAppIconLayerDpisV26(string(option.FolderName)),
			androidAdaptiveAppIconLogoDpisV26(string(option.FolderName)),
			option.OutputFileName,
		)
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
	foregroundName := fmt.Sprint(outputFileName, "_foreground")
	monochromeName := fmt.Sprint(outputFileName, "_monochrome")
	backgroundName := fmt.Sprint(outputFileName, "_background")

	for _, logo := range *logos {
		if shouldUseAssetName {
			foregroundName = fmt.Sprint(logo.imgNameWithoutExt, "_foreground")
			monochromeName = fmt.Sprint(logo.imgNameWithoutExt, "_monochrome")
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
			backgroundName = fmt.Sprint(bg.imgNameWithoutExt, "_background")
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

	sb.WriteString(`<?xml version="1.0" encoding="utf-8" ?>`)
	sb.WriteRune('\n')

	sb.WriteString(`<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">`)
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <background android:drawable="@mipmap/`, outputFileName, `_background" />`))
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <foreground android:drawable="@mipmap/`, outputFileName, `_foreground" />`))
	sb.WriteRune('\n')

	sb.WriteString(fmt.Sprint(`  <monochrome android:drawable="@mipmap/`, outputFileName, `_monochrome" />`))
	sb.WriteRune('\n')

	sb.WriteString(`</adaptive-icon>`)
	sb.WriteRune('\n')

	ic_launcher_xml := sb.String()

	dir := filepath.Join(logoImage.saveDirPath, "anydpi-v26")
	err := logoImage.rootDir.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := logoImage.rootDir.Create(filepath.Join(dir, "ic_launcher.xml"))
	if err != nil {
		return err
	}

	_, err = file.WriteString(ic_launcher_xml)
	if err != nil {
		return err
	}

	return nil
}
