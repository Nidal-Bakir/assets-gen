package cmd

import (
	"context"
	"path/filepath"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

func AndroidGooglePlayLogo() *cli.Command {
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

	var alphaThreshold float64
	var padding float64

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

		err = assetsgen.GenerateAndroidGooglePlayLogo(
			imagePath,
			assetsgen.AndroidGooglePlayLogoOptions{
				Padding:        padding,
				BgIcon:         bgIcon,
				AlphaThreshold: alphaThreshold,
				TrimWhiteSpace: trimWhiteSpace,
				MaskColor:      maskColor,
				OutputFileName: outputName,
			},
		)
		if err != nil {
			return err
		}

		if apply {
			err = applyAndroidPlayStoreLogo()
			if err != nil {
				return err
			}
		}

		return nil
	}

	usageText := `android-app-icon [command [command options]] <image path>

examples:
	apsl "./logo.png"
	apsl -bg linear-gradient --degree 90 --colors "#FF0000, #00FF00, #0000FF" --stops "0.0, 0.5, 1.0" "./logo.png"
	apsl --color "#0000FF" "./logo.png"
	apsl --apply -o "play_store" -p 0.1 --trim "./logo.png"`

	return &cli.Command{
		Name:      "android-google-play-logo",
		Aliases:   []string{"agpl"},
		UsageText: usageText,
		Usage:     "Generate Android Google Play logo 512x512",
		Action:    action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			paddingFlagFn(&padding),
			alphaThresholdFlagFn(&alphaThreshold),
			outputNameFlagFn(&outputName, "play_store_logo_512x512"),
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

func applyAndroidPlayStoreLogo() error {
	err := moveAndroidPlayStoreLogo()
	if err != nil {
		return err
	}
	err = deleteAssetsGenOutDir()
	if err != nil {
		return err
	}
	return nil
}

func moveAndroidPlayStoreLogo() error {
	adroidMainRootDir, err := getAndroidMainDirAsRoot()
	if err != nil {
		return err
	}
	adroidMainRootDir.Close()
	assetsOutRootDir, err := assetsgen.GetRootDir()
	if err != nil {
		return err
	}
	assetsOutRootDir.Close()

	src := filepath.Join(assetsOutRootDir.Name(), assetsgen.PlatformTypeAndroid, "main")
	dst := filepath.Join(adroidMainRootDir.Name())

	err = moveFilesR(src, dst)
	if err != nil {
		return err
	}

	return nil
}
