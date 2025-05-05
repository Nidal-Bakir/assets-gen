package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

// android-asset-gen (aag)
func AndroidAssetGen() *cli.Command {
	var imagePath string
	var trimWhiteSpace bool
	var apply bool

	folderName := assetsgen.AndroidFolderDrawable

	action := func(ctx context.Context, c *cli.Command) error {
		if b := isPathExist(imagePath); !b {
			if len(imagePath) == 0 {
				return ErrPleaseSpecifyImagePath
			}
			return assetsgen.ErrFileNotFound
		}

		err := assetsgen.GenerateImageAssetsForAndroid(
			imagePath, assetsgen.AndroidImageAssetsOptions{
				FolderName:     folderName,
				TrimWhiteSpace: trimWhiteSpace,
			},
		)
		if err != nil {
			return err
		}

		if apply {
			err = applyAndroidAssetImage()
			if err != nil {
				return err
			}
		}

		return nil
	}

	usageText := `android-asset-gen [command [command options]] <image path>

examples:
	aag "./clear_sky.png"
	aag --folder-name drawable --trim "./clear_sky.png"
	aag --apply "./clear_sky.png"`

	return &cli.Command{
		Name:    "android-asset-gen",
		Aliases: []string{"aag"},
		UsageText:   usageText,
		Usage: "Generate Android asset image for all DPIs",
		Action:  action,
		Arguments: []cli.Argument{
			imageArg(&imagePath),
		},
		Flags: []cli.Flag{
			androidFolderFlag(&folderName),
			trimWhiteSpaceFlagFn(&trimWhiteSpace),
			applyFlagFn(&apply),
		},
	}
}

func applyAndroidAssetImage() error {
	err := moveAndroidOutFiles()
	if err != nil {
		return err
	}
	return nil
}
