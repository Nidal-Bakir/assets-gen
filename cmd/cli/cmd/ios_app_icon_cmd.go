package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

// ios-app-icon (iai)
func IosAppIcon() *cli.Command {
	var imagePath string

	var bgType string
	var bgImagePath string
	var linearGradientDegree int
	var solidColor = colorful.Color{R: 1, G: 1, B: 1}
	var gradientColors = []colorful.Color{{R: 1, G: 1, B: 1}, {R: 0, G: 0, B: 0}}
	var gradientStops = []float64{0.0, 1.0}

	var trimWhiteSpace bool
	var alphaThreshold float64
	var padding float64

	action := func(ctx context.Context, c *cli.Command) error {
		if err := assetsgen.IsFileExistsAndImage(imagePath); err != nil {
			return err
		}

		bgIcon, err := getBgIcon(bgType, gradientColors, gradientStops, solidColor, linearGradientDegree, bgImagePath)
		if err != nil {
			return err
		}

		return assetsgen.GenerateAppIconForIos(
			imagePath,
			assetsgen.IosAppIconOptions{
				BgIcon:         bgIcon,
				Padding:        padding,
				AlphaThreshold: alphaThreshold,
				TrimWhiteSpace: trimWhiteSpace,
			},
		)
	}

	return &cli.Command{
		Name:    "ios-app-icon",
		Aliases: []string{"iai"},
		Action:  action,
		Arguments: []cli.Argument{
			imageArg(&imagePath),
		},
		Flags: []cli.Flag{
			paddingFlagFn(&padding),
			alphaThresholdFlagFn(&alphaThreshold),
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
