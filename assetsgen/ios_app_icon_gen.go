package assetsgen

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

func genIosAppIconDpis(ext string) []asset {
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

	for i, v := range iosAppIconDpis {
		iosDpi := v.(iosAppIconDpiAsset)
		iosDpi.Filename = fmt.Sprint(iosDpi.Filename, ext)
		iosAppIconDpis[i] = iosDpi
	}

	return iosAppIconDpis
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
func (a iosAppIconDpiAsset) CalcPadding(_, _ int) int {
	return 0
}

type IosAppIconOptions struct {
	BgIcon BackgroundIcon

	// between [0..1] as percentage of the maximum axis (w,h) of the image
	Padding float64
}

func GenerateAppIconForIos(imagePath string, option IosAppIconOptions) error {
	logoImage, err := genImageInfoForIos(imagePath, intentAppIcon)
	if err != nil {
		return err
	}

	iosAppIconDpis := genIosAppIconDpis(logoImage.imageExt)

	err = generateContentsJson(logoImage, iosAppIconDpis)
	if err != nil {
		return err
	}

	bounds := logoImage.img.Bounds()
	pad := math.Max(float64(bounds.Dx()), float64(bounds.Dy())) * option.Padding
	pad = math.Floor(pad)
	logoImage.squareImageWithPadding(int(pad))

	bgImage, err := option.BgIcon.generateImgInfo(logoImage)
	if err != nil {
		return err
	}

	err = generateIosAppIcon(logoImage, bgImage, iosAppIconDpis)
	if err != nil {
		return err
	}

	return nil
}

func generateIosAppIcon(logoImage imageInfo, bgImage imageInfo, iosAppIconDpis []asset) error {
	imgs := bgImage.
		stack(logoImage).
		splitPerAsset(iosAppIconDpis).
		resizeForAssets()

	for _, img := range *imgs {
		err := img.saveWithCustomName(img.asset.Name())
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

	dir, name := logoImage.genImageLocation("", "Contents.json")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}

	_, err = file.Write(jsonOut)
	if err != nil {
		return err
	}

	return nil
}
