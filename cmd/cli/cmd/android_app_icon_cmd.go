package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

var bgTypes = []string{"solid-color", "linear-gradient", "radial-gradient", "image"}

func AndroidAppIcon() *cli.Command {
	var imagePath string
	var outputName string

	var bgType string
	var bgImagePath string
	var linearGradientDegree int
	var solidColor = colorful.Color{R: 1, G: 1, B: 1}
	var gradientColors = []colorful.Color{{R: 1, G: 1, B: 1}, {R: 0, G: 0, B: 0}}
	var gradientStops = []float64{0.0, 1.0}

	var trimWhiteSpace bool
	var roundedCornerPercentRadius float64
	var alphaThreshold float64
	var padding float64
	var folderName = assetsgen.AndroidFolderMipmap

	imageArg := imageArg(&imagePath)

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}

		bgIcon, err := getBgIcon(bgType, gradientColors, gradientStops, solidColor, linearGradientDegree, bgImagePath)
		if err != nil {
			return err
		}

		return assetsgen.GenerateAppIconForAndroid(
			imagePath,
			outputName,
			assetsgen.AndroidAppIconOptions{
				RoundedCornerPercentRadius: roundedCornerPercentRadius,
				FolderName:                 folderName,
				Padding:                    padding,
				BgIcon:                     bgIcon,
				AlphaThreshold:             alphaThreshold,
				TrimWhiteSpace:             trimWhiteSpace,
			},
		)
	}

	return &cli.Command{
		Name:    "android-app-icon",
		Aliases: []string{"aai"},
		Action:  action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			cornerRadiusFlagFn(&roundedCornerPercentRadius),
			androidFolderFlag(&folderName),
			paddingFlagFn(&padding),
			alphaThresholdFlagFn(&alphaThreshold),
			outputNameFlagFn(&outputName),
			bgTypeFlagFn(&bgType),
			solidColorFlagFn(&solidColor),
			gradientColorsFlagFn(&gradientColors),
			gradientStopsFlagFn(&gradientStops),
			linearGradientDegreeFlagFn(&linearGradientDegree),
			imageBgFlagFn(&bgImagePath),
			trimWhiteSpaceFlagFn(&trimWhiteSpace),
		},
	}
}

func cornerRadiusFlagFn(roundedCornerRadius *float64) *cli.FloatFlag {
	return &cli.FloatFlag{
		Name:        "corner-radius",
		Aliases:     []string{"r"},
		Usage:       "Between [0..1] as percentage of the Radius. For example 1 would make the a full circle clip of the image, and 0 will do nothing, 0.5 will make rounded corners",
		Destination: roundedCornerRadius,
		Value:       0.2,
		Validator: func(i float64) error {
			if i < 0 || i > 1 {
				return ErrInvalidValueRange
			}
			return nil
		},
	}
}
