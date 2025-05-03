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

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}
		return assetsgen.GenerateNotificationIconForAndroid(imagePath, folderName, outputName, trimWhiteSpace)
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
		},
	}
}
