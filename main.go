package main

import (
	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	c, err := colorful.Hex("#FFFFFF")
	if err != nil {
		panic(err)
	}
	padding := 0.1
	imagePath := "./test_images/ic_launcher.png"
	BgIcon := assetsgen.NewSolidColorBackground(c)

	err = assetsgen.GenerateAppIconForAndroid(
		imagePath,
		"ic_launcher",
		assetsgen.AndroidAppIconOptions{
			RoundedCornerRadius: 25,
			BgIcon:              BgIcon,
			FolderName:          assetsgen.AndroidFolderMipmap,
			Padding:             padding,
		},
	)
	if err != nil {
		panic(err)
	}
	err = assetsgen.GenerateNotificationIconForAndroid(imagePath, assetsgen.AndroidFolderMipmap, "ic_stat_notification_icon")
	if err != nil {
		panic(err)
	}
	err = assetsgen.GenerateAppIconForIos(imagePath,
		assetsgen.IosAppIconOptions{
			BgIcon:  BgIcon,
			Padding: padding,
		},
	)
	if err != nil {
		panic(err)
	}
}
