package assetsgen

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func GenerateImageAsstesForAndroid(imagePath string, folderName AndroidFolderName) error {
	img, err := imgio.Open(imagePath)
	if err != nil {
		return err
	}

	enc, err := imageEncoderFromPath(imagePath)
	if err != nil {
		return err
	}

	imgName := filepath.Base(imagePath)
	imageExt := filepath.Ext(imagePath)
	imgNameWithoutExt := strings.ReplaceAll(imgName, imageExt, "")

	imgInfo := imageInfo{
		img:       img,
		encoder:   enc,
		imagePath: imagePath,
		imageName: imgName,
		imageExt:  imageExt,
		genImageLocation: func(screenType string) (directory string, imageName string) {
			dir := filepath.Join(
				rootFolderName,
				"android",
				fmt.Sprint("asset", "-", imgNameWithoutExt),
				fmt.Sprint(string(folderName), "-", screenType),
			)
			return dir, imgName
		},
	}

	err = generateImageAsstesForAndroid(imgInfo, androidScreenDpis)
	if err != nil {
		return err
	}

	return nil
}

func generateImageAsstesForAndroid(imgInfo imageInfo, dpis screenTypeSlice) error {
	maxScaleFactor := dpis.maxScaleFactor()

	imageBounds := imgInfo.img.Bounds().Max
	baseW := float64(imageBounds.X) / maxScaleFactor
	baseH := float64(imageBounds.Y) / maxScaleFactor

	for _, dpi := range dpis {
		w := math.Floor(baseW * dpi.scaleFactor)
		h := math.Floor(baseH * dpi.scaleFactor)
		resizedImg := transform.Resize(imgInfo.img, int(w), int(h), transform.Linear)

		dir, name := imgInfo.genImageLocation(dpi.name)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		err = imgio.Save(filepath.Join(dir, name), resizedImg, imgInfo.encoder)
		if err != nil {
			return err
		}
	}

	return nil
}
