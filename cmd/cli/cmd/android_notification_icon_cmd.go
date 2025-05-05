package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

// android-notification-icon (ani)
func AndroidNotificationIcon() *cli.Command {
	var imagePath string
	var outputName string
	folderName := assetsgen.AndroidFolderMipmap
	var trimWhiteSpace bool
	var alphaThreshold float64

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}
		return assetsgen.GenerateNotificationIconForAndroid(
			imagePath,
			assetsgen.AndroidNotificationIconOptions{
				FolderName:     folderName,
				TrimWhiteSpace: trimWhiteSpace,
				OutputFileName: outputName,
				AlphaThreshold: alphaThreshold,
			},
		)
	}

	return &cli.Command{
		Name:    "android-notification-icon",
		Aliases: []string{"ani"},
		Action:  action,
		Arguments: []cli.Argument{
			imageArg(&imagePath),
		},
		Flags: []cli.Flag{
			androidFolderFlag(&folderName),
			outputNameFlagFn(&outputName),
			trimWhiteSpaceFlagFn(&trimWhiteSpace),
			alphaThresholdFlagFn(&alphaThreshold),
		},
	}
}
