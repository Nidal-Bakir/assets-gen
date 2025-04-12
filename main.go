package main

import "github.com/Nidal-Bakir/assets-gen/assetsgen"

func main() {
	err := assetsgen.GenerateImageAsstesForAndroid("./launch_image.png", assetsgen.Mipmap)
	if err != nil {
		panic(err)
	}
}
