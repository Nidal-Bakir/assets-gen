package cmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v3"
)

var bgTypes = []string{"solid-color", "linear-gradien", "radial-gradient", "image"}

func AndroidAppIcon() *cli.Command {
	var imagePath string
	var outputName string
	var roundedCornerPercentRadius float64

	var bgType string
	var bgImagePath string
	var linearGradientDegree int
	var solidColor = colorful.Color{R: 1, G: 1, B: 1}
	var gradientColors = []colorful.Color{colorful.Color{R: 1, G: 1, B: 1}, colorful.Color{R: 0, G: 0, B: 0}}
	var gradientStops = []float64{0.0, 1.0}

	var padding float64
	var folderName = assetsgen.AndroidFolderMipmap

	imageArg := imageArg(&imagePath)

	cornerRadiusFlag := cornerRadiusFlagFn(&roundedCornerPercentRadius)
	folderNameFlag := androidFolderFlag(&folderName)
	paddingFlag := paddingFlagFn(&padding)
	outputNameFlag := outputNameFlagFn(&outputName)
	bgTypeFlag := bgTypeFlagFn(&bgType)
	solidColorFlag := solidColorFlagFn(&solidColor)
	gradientColorsFlag := gradientColorsFlagFn(&gradientColors)
	gradientStopsFlag := gradientStopsFlagFn(&gradientStops)
	linearGradientDegreeFlag := linearGradientDegreeFlagFn(&linearGradientDegree)
	imageBgFlag := imageBgFlagFn(&bgImagePath)

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
			cornerRadiusFlag,
			folderNameFlag,
			outputNameFlag,
			paddingFlag,
			bgTypeFlag,
			solidColorFlag,
			gradientColorsFlag,
			gradientStopsFlag,
			linearGradientDegreeFlag,
			imageBgFlag,
		},
	}
}

func getBgIcon(
	bgType string,
	gradientColors []colorful.Color,
	gradientStops []float64,
	solidColor colorful.Color,
	linearGradientDegree int,
	bgImagePath string,
) (assetsgen.BackgroundIcon, error) {
	var BgIcon assetsgen.BackgroundIcon

	switch bgType {
	case "solid-color":
		BgIcon = assetsgen.NewSolidColorBackground(solidColor)

	case "linear-gradien":
		table, err := generateGradientTable(gradientColors, gradientStops)
		if err != nil {
			return nil, err
		}
		BgIcon = assetsgen.NewLinearGradientBackground(table, linearGradientDegree)

	case "radial-gradient":
		table, err := generateGradientTable(gradientColors, gradientStops)
		if err != nil {
			return nil, err
		}
		BgIcon = assetsgen.NewRadialGradientBackground(table)

	case "image":
		BgIcon = assetsgen.NewImageBackground(bgImagePath)

	default:
		panic("we should not be here")
	}

	return BgIcon, nil
}

func generateGradientTable(colors []colorful.Color, stops []float64) (assetsgen.GradientTable, error) {
	if len(colors) != len(stops) {
		return nil, ErrColorsAndStopsLengthDidNotMatch
	}

	table := make(assetsgen.GradientTable, len(colors))

	for i, c := range colors {
		table[i] = assetsgen.GradientTableItem{
			Col: c,
			Pos: stops[i],
		}
	}

	return table, nil
}

func paddingFlagFn(padding *float64) *cli.FloatFlag {
	return &cli.FloatFlag{
		Name:        "padding",
		Aliases:     []string{"p"},
		Destination: padding,
		Value:       0.1,
		Validator: func(i float64) error {
			if i < 0 || i > 1 {
				return ErrPaddingOutOfRange
			}
			return nil
		},
	}
}

func outputNameFlagFn(outputName *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Value:       "",
		Usage:       "Set a cutom output name for the generated files. Only name without extension.",
		DefaultText: "The default is to use the image name as output name",
		Destination: outputName,
	}
}

func cornerRadiusFlagFn(roundedCornerRadius *float64) *cli.FloatFlag {
	return &cli.FloatFlag{
		Name:        "corner-radius",
		Aliases:     []string{"r"},
		Destination: roundedCornerRadius,
		Value:       0,
		Validator: func(i float64) error {
			if i < 0 || i > 1 {
				return ErrInvalidValueRange
			}
			return nil
		},
	}
}

func bgTypeFlagFn(bgType *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "bg-type",
		Aliases: []string{"bg"},
		Value:   "solid-color",
		Usage:   fmt.Sprint("Set the backgorund type: ", strings.Join(bgTypes, ",")),
		Validator: func(s string) error {

			if slices.Contains(bgTypes, s) {
				return nil
			}
			return ErrInvalidBgType
		},
		Destination: bgType,
	}
}

func imageBgFlagFn(bgImagePath *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "bg-path",
		Value:       "",
		Usage:       "Path to the backgound image",
		Destination: bgImagePath,
		Validator: func(imagePath string) error {
			return assetsgen.IsFileExistsAndImage(imagePath)
		},
	}
}

func solidColorFlagFn(solidBgColor *colorful.Color) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "color",
		Value: "#FFFFFF",
		Usage: "The solid backgorund color default to white",
		Action: func(ctx context.Context, c *cli.Command, s string) error {
			color, err := colorful.Hex(s)
			*solidBgColor = color
			return err
		},
	}
}

func linearGradientDegreeFlagFn(degree *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "degree",
		Value:       0,
		Usage:       "The angle of rotation for the linear gradient background",
		Destination: degree,
	}
}

func gradientColorsFlagFn(colors *[]colorful.Color) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "colors",
		Value: "#FFFFFF,#FFFFFF",
		Usage: "The gradient backgorund colors, comman saperated e.g: #0000FF,#FF0000. You should supply the stops also. The colors count should match the stops",
		Action: func(ctx context.Context, c *cli.Command, s string) error {
			if len(s) == 0 {
				return nil
			}

			colorsFromUser := strings.Split(s, ",")
			*colors = make([]colorful.Color, len(colorsFromUser))
			for i, colorStr := range colorsFromUser {
				c, err := colorful.Hex(colorStr)
				if err != nil {
					return ErrInvalidColor
				}
				(*colors)[i] = c
			}

			return nil
		},
	}
}

func gradientStopsFlagFn(stops *[]float64) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "stops",
		Value: "0.0,1.0",
		Usage: "The gradient backgorund colors stops, comman saperated e.g: 0.0,1.0. The stops count should match the colors",
		Action: func(ctx context.Context, c *cli.Command, s string) error {
			if len(s) == 0 {
				return nil
			}

			stopsFromUser := strings.Split(s, ",")
			*stops = make([]float64, len(stopsFromUser))
			for i, stopStr := range stopsFromUser {
				stop, err := strconv.ParseFloat(stopStr, 64)
				if err != nil {
					return err
				}
				(*stops)[i] = stop
			}

			return nil
		},
	}
}
