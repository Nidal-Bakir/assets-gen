package assetsgen

import (
	"encoding/json"
	"path/filepath"

	"github.com/lucasb-eyer/go-colorful"
)

var iosAppIconDpis = []asset{
	iosAppIconDpiAsset{
		Filename: "AppIcon@2x",
		Idiom:    "iphone",
		Scale:    "2x",
		SizeName: "60x60",
		Size:     120,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon@3x",
		Idiom:    "iphone",
		Scale:    "3x",
		SizeName: "60x60",
		Size:     180,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon~ipad",
		Idiom:    "ipad",
		Scale:    "1x",
		SizeName: "76x76",
		Size:     76,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon@2x~ipad",
		Idiom:    "ipad",
		Scale:    "2x",
		SizeName: "76x76",
		Size:     152,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-83.5@2x~ipad",
		Idiom:    "ipad",
		Scale:    "2x",
		SizeName: "83.5x83.5",
		Size:     167,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-40@2x",
		Idiom:    "iphone",
		Scale:    "2x",
		SizeName: "40x40",
		Size:     80,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-40@3x",
		Idiom:    "iphone",
		Scale:    "3x",
		SizeName: "40x40",
		Size:     120,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-40~ipad",
		Idiom:    "ipad",
		Scale:    "1x",
		SizeName: "40x40",
		Size:     40,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-40@2x~ipad",
		Idiom:    "ipad",
		Scale:    "2x",
		SizeName: "40x40",
		Size:     80,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-20@2x",
		Idiom:    "iphone",
		Scale:    "2x",
		SizeName: "20x20",
		Size:     40,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-20@3x",
		Idiom:    "iphone",
		Scale:    "3x",
		SizeName: "20x20",
		Size:     60,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-20~ipad",
		Idiom:    "ipad",
		Scale:    "1x",
		SizeName: "20x20",
		Size:     20,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-20@2x~ipad",
		Idiom:    "ipad",
		Scale:    "2x",
		SizeName: "20x20",
		Size:     40,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-29",
		Idiom:    "iphone",
		Scale:    "1x",
		SizeName: "29x29",
		Size:     29,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-29@2x",
		Idiom:    "iphone",
		Scale:    "2x",
		SizeName: "29x29",
		Size:     58,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-29@3x",
		Idiom:    "iphone",
		Scale:    "3x",
		SizeName: "29x29",
		Size:     87,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-29~ipad",
		Idiom:    "ipad",
		Scale:    "1x",
		SizeName: "29x29",
		Size:     29,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-29@2x~ipad",
		Idiom:    "ipad",
		Scale:    "2x",
		SizeName: "29x29",
		Size:     58,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-60@2x~car",
		Idiom:    "car",
		Scale:    "2x",
		SizeName: "60x60",
		Size:     120,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon-60@3x~car",
		Idiom:    "car",
		Scale:    "3x",
		SizeName: "60x60",
		Size:     180,
	},
	iosAppIconDpiAsset{
		Filename: "AppIcon~ios-marketing",
		Idiom:    "ios-marketing",
		Scale:    "1x",
		SizeName: "1024x1024",
		Size:     1024,
	},
}

type iosAppIconDpiAsset struct {
	Filename string `json:"filename"`
	Idiom    string `json:"idiom"`
	Scale    string `json:"Scale"`
	Size     int    `json:"-"`
	SizeName string `json:"size"`
}

func (a iosAppIconDpiAsset) Name() string {
	return a.Filename
}

func (a iosAppIconDpiAsset) CalcSize(_, _ int) (int, int) {
	return a.Size, a.Size
}

func (a iosAppIconDpiAsset) DirName() string {
	return ""
}

type IosAppIconOptions struct {
	BgIcon BackgroundIcon

	// between [0..1] as percentage of how match the pixel should be transparent to keep its original color. Use -1 to disable
	AlphaThreshold float64

	// between [0..1] as percentage of the maximum axis (w,h) of the image
	Padding float64

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool

	MaskColor *colorful.Color
}

func GenerateAppIconForIos(imagePath string, option IosAppIconOptions) error {
	logoImage, err := genLogoImageForIos(imagePath, option)
	if err != nil {
		return err
	}
	defer logoImage.rootDir.Close()

	bgImage, err := option.BgIcon.generateImgInfo(logoImage)
	if err != nil {
		return err
	}

	err = generateContentsJson(logoImage, iosAppIconDpis)
	if err != nil {
		return err
	}

	err = generateIosAppIcon(logoImage, bgImage, option.AlphaThreshold, iosAppIconDpis)
	if err != nil {
		return err
	}

	return nil

}

func genLogoImageForIos(imagePath string, option IosAppIconOptions) (imageInfo, error) {
	logoImage, err := newImageInfo(
		imagePath,
		filepath.Join(PlatformTypeIos, "Assets.xcassets", "AppIcon.appiconset"),
	)
	if err != nil {
		return logoImage, err
	}

	pad := calPadding(logoImage.img, option.Padding)

	logoImage.
		If(option.TrimWhiteSpace, logoImage.TrimWhiteSpace).
		SquareImageWithEmptyPixels(pad).
		// Padding(pad).
		If(option.MaskColor != nil, func() *imageInfo { return logoImage.ConvertNoneOpaqueToColor(*option.MaskColor) })

	return logoImage, nil
}

func generateIosAppIcon(logoImage imageInfo, bgImage imageInfo, alphaThreshold float64, iosAppIconDpis []asset) error {
	imgs := bgImage.
		IfElse(
			alphaThreshold < 0,
			func() *imageInfo { return bgImage.Stack(logoImage) },
			func() *imageInfo { return bgImage.StackWithNoAlpha(alphaThreshold, logoImage) },
		).
		SplitPerAsset(iosAppIconDpis).
		ResizeForAssets()

	for _, img := range *imgs {
		err := img.SaveWithCustomName(img.asset.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func generateContentsJson(logoImage imageInfo, dpis []asset) error {
	type GenInfo struct {
		Author  string `json:"author"`
		Version int    `json:"version"`
	}
	type output struct {
		Images []asset `json:"images"`
		Info   GenInfo `json:"info"`
	}

	out := output{
		Images: dpis,
		Info: GenInfo{
			Author:  "assets-gen@Nidal-Bakir",
			Version: 1,
		},
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		return err
	}

	file, err := logoImage.rootDir.Create(filepath.Join(logoImage.saveDirPath, "Contents.json"))
	if err != nil {
		return err
	}

	_, err = file.Write(jsonOut)
	if err != nil {
		return err
	}

	return nil
}
