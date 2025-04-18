package assetsgen

import (
	"fmt"
)

type androidFolderName string

const (
	AndroidFolderMipmap   androidFolderName = "mipmap"
	AndroidFolderDrawable androidFolderName = "drawable"

	rootFolderName string = "assets_gen_out"
)

func genImageInfoForAndroid(imagePath string, folderName androidFolderName, intent intention) (imageInfo, error) {
	return newImageInfo(
		imagePath,
		platformTypeAndroid,
		intent,
		func(screenType string) string {
			return fmt.Sprint(string(folderName), "-", screenType)
		},
	)
}
