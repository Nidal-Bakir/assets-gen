package assetsgen

func genImageInfoForIos(imagePath string, intent intention) (imageInfo, error) {
	return newImageInfo(
		imagePath,
		platformTypeIos,
		intent,
		func(screenType string) string {
			return ""
		},
	)
}
