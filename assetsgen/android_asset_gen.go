package assetsgen

import (
	"fmt"
	"math"
	"path/filepath"
)

type AndroidImageAssetsOptions struct {
	FolderName AndroidFolderName

	// removes the white spaces from the edges of the logo
	TrimWhiteSpace bool
}

func GenerateImageAssetsForAndroid(imagePath string, option AndroidImageAssetsOptions) error {
	imgInfo, err := newImageInfo(
		imagePath,
		filepath.Join(PlatformTypeAndroid, "res"),
	)
	if err != nil {
		return err
	}
	defer imgInfo.rootDir.Close()

	if option.TrimWhiteSpace {
		imgInfo.TrimWhiteSpace()
	}

	imgBounds := imgInfo.img.Bounds()
	androidScreenDpis := generateAndroidScreenDpis(imgBounds.Dx(), imgBounds.Dy(), string(option.FolderName))

	err = imgInfo.
		SplitPerAsset(androidScreenDpis).
		ResizeForAssets().
		Save()
	if err != nil {
		return err
	}

	return nil
}

// MDPI    - 1.0x
// HDPI    - 1.5x
// XHDPI   - 2.0x
// XXHDPI  - 3.0x
// XXXHDPI - 4.0x
func generateAndroidScreenDpis(w, h int, androidFolderName string) []asset {
	androidScreenDpis := []asset{
		androidScreenDpiAsset{
			dpiName:     "mdpi",
			scaleFactor: 1.0,
		},
		androidScreenDpiAsset{
			dpiName:     "hdpi",
			scaleFactor: 1.5,
		},
		androidScreenDpiAsset{
			dpiName:     "xhdpi",
			scaleFactor: 2,
		},
		androidScreenDpiAsset{
			dpiName:     "xxhdpi",
			scaleFactor: 3,
		},
		androidScreenDpiAsset{
			dpiName:     "xxxhdpi",
			scaleFactor: 4,
		},
	}

	var maxScaleFactor float64
	for _, v := range androidScreenDpis {
		screenDpi := v.(androidScreenDpiAsset)
		maxScaleFactor = max(screenDpi.scaleFactor, maxScaleFactor)
	}

	baseW := float64(w) / maxScaleFactor
	baseH := float64(h) / maxScaleFactor
	for i, v := range androidScreenDpis {
		screenDpi := v.(androidScreenDpiAsset)
		screenDpi.baseW = int(math.Floor(baseW))
		screenDpi.baseH = int(math.Floor(baseH))
		screenDpi.dirName = androidFolderName
		androidScreenDpis[i] = screenDpi
	}

	return androidScreenDpis
}

type androidScreenDpiAsset struct {
	dpiName     string
	dirName     string
	scaleFactor float64
	baseW       int
	baseH       int
}

func (a androidScreenDpiAsset) Name() string {
	return a.dpiName
}

func (a androidScreenDpiAsset) CalcSize(_, _ int) (int, int) {
	w := int(math.Floor(float64(a.baseW) * a.scaleFactor))
	h := int(math.Floor(float64(a.baseH) * a.scaleFactor))
	return w, h
}

func (a androidScreenDpiAsset) DirName() string {
	return fmt.Sprint(a.dirName, "-", a.dpiName)
}
