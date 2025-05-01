package cmd

import (
	"context"
	"errors"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

var (
	ErrInvalidBgType                   = errors.New("invalid bg-type")
	ErrInvalidAndroidFolder            = errors.New("invalid android folder name. possible values (mipmap, drawable)")
	ErrInvalidValueRange               = errors.New("invalid value range")
	ErrPaddingOutOfRange               = errors.New("padding should be between 0..1")
	ErrInvalidColor                    = errors.New("invalid color. e.g of valid colors #0000FF, #FFFFFF")
	ErrColorsAndStopsLengthDidNotMatch = errors.New("The length fo colors should match the length of stops")
)

func androidFolderFlag(folderName *assetsgen.AndroidFolderName) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "folder-name",
		Aliases: []string{"f"},
		Value:   string(*folderName),
		Usage:   "whether to target mipmap or drawable folder",
		Action: func(ctx context.Context, c *cli.Command, s string) error {
			switch s {
			case string(assetsgen.AndroidFolderDrawable):
				*folderName = assetsgen.AndroidFolderDrawable
			case string(assetsgen.AndroidFolderMipmap):
				*folderName = assetsgen.AndroidFolderMipmap
			default:
				return ErrInvalidAndroidFolder
			}
			return nil
		},
	}
}

func imageArg(imagePath *string) *cli.StringArg {
	return &cli.StringArg{
		Name:        "image",
		UsageText:   "The image path",
		Destination: imagePath,
	}
}
