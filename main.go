package main

import (
	"fmt"
	"time"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
)

func MustParsHex(hex string) colorful.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}

func main() {
	startTime := time.Now()

	table := assetsgen.GradientTable{
		{
			MustParsHex("#262d4d"),
			0.0,
		},
		{
			MustParsHex("#7BD9EF"),
			1.0,
		},
	}

	err := assetsgen.GenerateAppIconForAndroid(
		"./test_images/ic_launcher.png",
		assetsgen.AppIconOptions{
			FolderName:          assetsgen.AndroidFolderMipmap,
			Padding:             0,
			BgIcon:              assetsgen.NewLinearGradientBackground(table, 0),
			RoundedCornerRadius: 100,
		})

	if err != nil {
		panic(err)
	}

	fmt.Printf("took %.1f seconds\n", time.Now().Sub(startTime).Seconds())
}
