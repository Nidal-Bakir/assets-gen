package main

import "github.com/Nidal-Bakir/assets-gen/assetsgen"

func main() {
	err := assetsgen.GenerateAppIconForAndroid("./test_image.png", assetsgen.AndroidFolderMipmap)
	if err != nil {
		panic(err)
	}
}
