package assetsgen

// MDPI    - 24px
// HDPI    - 36px
// XHDPI   - 48px
// XXHDPI  - 72px
// XXXHDPI - 96px
var androidAppIconDpis = []Asset{
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

func GenerateAppIconForAndroid(imagePath string, folderName androidFolderName) error {
	imgInfo, err := genImageInfoForAndroid(imagePath, folderName, intentAppIcon)
	if err != nil {
		return err
	}

	imgInfo.img = squareImageWithPadding(imgInfo.img)

	err = generateImageAsstes(imgInfo, androidAppIconDpis)
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
