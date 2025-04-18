package main

import "github.com/Nidal-Bakir/assets-gen/assetsgen"

func main() {
	padding := 0

	err := assetsgen.GenerateAppIconForAndroid("./test_images/ic_launcher.png", assetsgen.AndroidFolderMipmap, padding)
	if err != nil {
		panic(err)
	}
}
