package cmd

import (
	"context"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

func AndroidAppIcon() *cli.Command {
	var imagePath string
	var outputName string

	var bgType string
	var bgImagePath string
	var linearGradientDegree int
	var solidColor = colorful.Color{R: 1, G: 1, B: 1}
	var gradientColors = []colorful.Color{{R: 1, G: 1, B: 1}, {R: 0, G: 0, B: 0}}
	var gradientStops = []float64{0.0, 1.0}

	var maskColor *colorful.Color
	var trimWhiteSpace bool
	var apply bool
	var roundedCornerPercentRadius float64
	var alphaThreshold float64
	var padding float64
	var folderName = assetsgen.AndroidFolderMipmap

	imageArg := imageArg(&imagePath)

	action := func(ctx context.Context, c *cli.Command) error {
		if b := isPathExist(imagePath); !b {
			if len(imagePath) == 0 {
				return ErrPleaseSpecifyImagePath
			}
			return assetsgen.ErrFileNotFound
		}

		bgIcon, err := getBgIcon(bgType, gradientColors, gradientStops, solidColor, linearGradientDegree, bgImagePath)
		if err != nil {
			return err
		}

		err = assetsgen.GenerateAppIconForAndroid(
			imagePath,
			assetsgen.AndroidAppIconOptions{
				RoundedCornerPercentRadius: roundedCornerPercentRadius,
				FolderName:                 folderName,
				Padding:                    padding,
				BgIcon:                     bgIcon,
				AlphaThreshold:             alphaThreshold,
				TrimWhiteSpace:             trimWhiteSpace,
				MaskColor:                  maskColor,
				OutputFileName:             outputName,
			},
		)
		if err != nil {
			return err
		}

		if apply {
			err = applyAndroidAppIcon(outputName)
			if err != nil {
				return err
			}
		}

		return nil
	}

	usageText := `android-app-icon [command [command options]] <image path>

examples:
	aai "./ic_launcher.png"
	aai -bg linear-gradient --degree 90 --colors "#FF0000, #00FF00, #0000FF" --stops "0.0, 0.5, 1.0" "./ic_launcher.png"
	aai --color "#0000FF" "./ic_launcher.png"
	aai --apply -o "app_icon" -p 0.1 --trim "./ic_launcher.png"`

	return &cli.Command{
		Name:      "android-app-icon",
		Aliases:   []string{"aai"},
		UsageText: usageText,
		Usage:     "Generate Android app launcher icons",
		Action:    action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			cornerRadiusFlagFn(&roundedCornerPercentRadius),
			androidFolderFlag(&folderName),
			paddingFlagFn(&padding),
			alphaThresholdFlagFn(&alphaThreshold),
			outputNameFlagFn(&outputName, "ic_launcher"),
			bgTypeFlagFn(&bgType),
			solidColorFlagFn(&solidColor),
			gradientColorsFlagFn(&gradientColors),
			gradientStopsFlagFn(&gradientStops),
			linearGradientDegreeFlagFn(&linearGradientDegree),
			imageBgFlagFn(&bgImagePath),
			trimWhiteSpaceFlagFn(&trimWhiteSpace),
			maskColorFlagFn(&maskColor),
			applyFlagFn(&apply),
		},
	}
}

func cornerRadiusFlagFn(roundedCornerRadius *float64) *cli.FloatFlag {
	return &cli.FloatFlag{
		Name:        "corner-radius",
		Aliases:     []string{"r"},
		Usage:       "Between [0..1] as percentage of the Radius. For example 1 would make the a full circle clip of the image, and 0 will do nothing, 0.5 will make rounded corners",
		Destination: roundedCornerRadius,
		Value:       1,
		Validator: func(i float64) error {
			if i < 0 || i > 1 {
				return ErrInvalidValueRange
			}
			return nil
		},
	}
}

func applyAndroidAppIcon(outputFileName string) error {
	err := moveResAndroidOutFiles()
	if err != nil {
		return err
	}
	err = deleteAssetsGenOutDir()
	if err != nil {
		return err
	}
	return nil
}
