package cmd

import (
	"context"
	"errors"
	"sync"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

func GenerateAll() *cli.Command {
	var imagePath string

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

		wg := sync.WaitGroup{}
		wg.Add(4)
		errSlice := make([]error, 4)

		go func() {
			defer wg.Done()
			errSlice[0] = assetsgen.GenerateAppIconForAndroid(
				imagePath,
				assetsgen.AndroidAppIconOptions{
					RoundedCornerPercentRadius: roundedCornerPercentRadius,
					FolderName:                 folderName,
					Padding:                    padding,
					BgIcon:                     bgIcon,
					AlphaThreshold:             alphaThreshold,
					TrimWhiteSpace:             trimWhiteSpace,
					MaskColor:                  maskColor,
					OutputFileName:             "ic_launcher",
				},
			)
		}()

		go func() {
			defer wg.Done()
			errSlice[1] = assetsgen.GenerateAndroidGooglePlayLogo(
				imagePath,
				assetsgen.AndroidGooglePlayLogoOptions{
					Padding:        padding,
					BgIcon:         bgIcon,
					AlphaThreshold: alphaThreshold,
					TrimWhiteSpace: trimWhiteSpace,
					MaskColor:      maskColor,
					OutputFileName: "play_store_logo_512x512",
				},
			)
		}()

		go func() {
			defer wg.Done()
			errSlice[2] = assetsgen.GenerateNotificationIconForAndroid(
				imagePath,
				assetsgen.AndroidNotificationIconOptions{
					FolderName:     folderName,
					TrimWhiteSpace: trimWhiteSpace,
					OutputFileName: "ic_stat_notification_icon",
					AlphaThreshold: alphaThreshold,
				},
			)
		}()

		go func() {
			defer wg.Done()
			errSlice[3] = assetsgen.GenerateAppIconForIos(
				imagePath,
				assetsgen.IosAppIconOptions{
					BgIcon:         bgIcon,
					Padding:        padding,
					AlphaThreshold: alphaThreshold,
					TrimWhiteSpace: trimWhiteSpace,
					MaskColor:      maskColor,
				},
			)
		}()

		wg.Wait()

		err = errors.Join(errSlice...)
		if err != nil {
			return err
		}

		if apply {
			err = applyAll()
			if err != nil {
				return err
			}
		}

		return nil
	}

	return &cli.Command{
		Name:   "all",
		Usage:  "Generate Android app launcher icons, Android notification asset, Google Play logo, and IOS app icon",
		Action: action,
		Arguments: []cli.Argument{
			imageArg,
		},
		Flags: []cli.Flag{
			cornerRadiusFlagFn(&roundedCornerPercentRadius),
			androidFolderFlag(&folderName),
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

func applyAll() error {
	err := moveResAndroidOutFiles()
	if err != nil {
		return err
	}

	err = moveAndroidPlayStoreLogo()
	if err != nil {
		return err
	}

	err = moveIosOutFiles()
	if err != nil {
		return err
	}

	err = deleteAssetsGenOutDir()
	if err != nil {
		return err
	}

	return nil
}
