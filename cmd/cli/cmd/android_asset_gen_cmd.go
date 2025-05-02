package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

// android-asset-gen (aag)
func AndroidAssetGen() *cli.Command {
	var imagePath string
	folderName := assetsgen.AndroidFolderDrawable

	imageArg := imageArg(&imagePath)

	folderNameFlag := androidFolderFlag(&folderName)

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}
		return assetsgen.GenerateImageAssetsForAndroid(imagePath, folderName)
	}

	return &cli.Command{
		Name:      "android-asset-gen",
		Aliases:   []string{"aag"},
		Usage:     "generate asset image for all the DPIs: MDPI, HDPI, XHDPI, XXHDPI, XXXHDPI",
		ArgsUsage: "the image path",
		Action:    action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			folderNameFlag,
		},
	}
}
