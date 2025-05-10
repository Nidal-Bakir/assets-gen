package cmd

import (
	"context"
	"os"
	"path/filepath"

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

	var maskColor *colorful.Color
	var trimWhiteSpace bool
	var alphaThreshold float64
	var padding float64
	var apply bool

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

		err = assetsgen.GenerateAppIconForIos(
			imagePath,
			assetsgen.IosAppIconOptions{
				BgIcon:         bgIcon,
				Padding:        padding,
				AlphaThreshold: alphaThreshold,
				TrimWhiteSpace: trimWhiteSpace,
				MaskColor:      maskColor,
			},
		)
		if err != nil {
			return err
		}

		if apply {
			err = applyIosAppIcon()
			if err != nil {
				return err
			}
		}

		return nil
	}

	usageText := `ios-app-icon [command [command options]] <image path>

examples:
	iai "./app_icon.png"
	iai -bg linear-gradient --degree 90 --colors "#FF0000, #00FF00, #0000FF" --stops "0.0, 0.5, 1.0" "./app_icon.png"
	iai --color "#0000FF" "./app_icon.png"
	iai --apply -p 0.1 --trim "./app_icon.png"`

	return &cli.Command{
		Name:      "ios-app-icon",
		Aliases:   []string{"iai"},
		UsageText: usageText,
		Usage:     "Generate IOS app icon",
		Action:    action,
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
			maskColorFlagFn(&maskColor),
			applyFlagFn(&apply),
		},
	}
}

func applyIosAppIcon() error {
	err := moveIosOutFiles()
	if err != nil {
		return err
	}
	return nil
}

func moveIosOutFiles() error {
	xcassetsRootDir, err := getIosXcassetsAsRoot()
	if err != nil {
		return err
	}
	xcassetsRootDir.Close()
	assetsOutRootDir, err := assetsgen.GetRootDir()
	if err != nil {
		return err
	}
	assetsOutRootDir.Close()

	src := filepath.Join(assetsOutRootDir.Name(), assetsgen.PlatformTypeIos, "Assets.xcassets", "AppIcon.appiconset")
	dst := filepath.Join(xcassetsRootDir.Name(), "AppIcon.appiconset")
	err = os.RemoveAll(dst)
	if err != nil {
		return err
	}

	err = moveFilesR(src, dst)
	if err != nil {
		return err
	}

	err = os.RemoveAll(assetsOutRootDir.Name())
	if err != nil {
		return err
	}

	return nil
}
