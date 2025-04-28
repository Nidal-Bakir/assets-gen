package assetsgen

import (
	"fmt"
)

type AndroidFolderName string

const (
	AndroidFolderMipmap   AndroidFolderName = "mipmap"
	AndroidFolderDrawable AndroidFolderName = "drawable"
)

func genImageInfoForAndroid(imagePath string, folderName AndroidFolderName, intent intention) (imageInfo, error) {
	return newImageInfo(
		imagePath,
		platformTypeAndroid,
		intent,
		func(screenType string) string {
			return fmt.Sprint(string(folderName), "-", screenType)
		},
	)
}
