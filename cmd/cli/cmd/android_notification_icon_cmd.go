package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

func AndroidNotificationIcon() *cli.Command {
	var imagePath string
	folderName := assetsgen.AndroidFolderMipmap

	imageArg := imageArg(&imagePath)

	folderNameFlag := androidFolderFlag(&folderName)

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}
		return nil
	}
	return &cli.Command{
		Name:    "android-notification-icon",
		Aliases: []string{"ani"},
		Action:  action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			folderNameFlag,
		},
	}
}
