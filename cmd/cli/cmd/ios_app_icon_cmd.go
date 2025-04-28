package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

func IosAppIcon() *cli.Command {
	var imagePath string

	imageArg := &cli.StringArg{
		Name:        "image",
		UsageText:   "The image path",
		Destination: &imagePath,
	}

	action := func(ctx context.Context, c *cli.Command) error {
		if !assetsgen.IsFileExists(imagePath) {
			return ErrFileNotFound
		}
		return nil
	}
	return &cli.Command{
		Name:    "ios-app-icon",
		Aliases: []string{"iai"},
		Action:  action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{},
	}
}
