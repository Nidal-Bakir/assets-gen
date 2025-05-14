package cmd

import (
	"context"
	"fmt"

	"github.com/Nidal-Bakir/assets-gen/assetsgen"
	"github.com/urfave/cli/v3"
)

// android-notification-icon (ani)
func AndroidNotificationIcon() *cli.Command {
	var imagePath string
	var outputName string
	folderName := assetsgen.AndroidFolderMipmap
	var trimWhiteSpace bool
	var alphaThreshold float64
	var apply bool

	action := func(ctx context.Context, c *cli.Command) error {
		if b := isPathExist(imagePath); !b {
			if len(imagePath) == 0 {
				return ErrPleaseSpecifyImagePath
			}
			return assetsgen.ErrFileNotFound
		}

		err := assetsgen.GenerateNotificationIconForAndroid(
			imagePath,
			assetsgen.AndroidNotificationIconOptions{
				FolderName:     folderName,
				TrimWhiteSpace: trimWhiteSpace,
				OutputFileName: outputName,
				AlphaThreshold: alphaThreshold,
			},
		)
		if err != nil {
			return err
		}

		if apply {
			err = applyAndroidNotificationIcon(string(folderName), outputName)
			if err != nil {
				return err
			}
		}

		return nil
	}

	usageText := `android-notification-icon [command [command options]] <image path>

examples:
	aai "./icon.png"
	aai --apply -o "notification_icon" --trim "./icon.png"`

	return &cli.Command{
		Name:      "android-notification-icon",
		Aliases:   []string{"ani"},
		Action:    action,
		UsageText: usageText,
		Usage:     "Generate Android notification asset",
		Arguments: []cli.Argument{
			imageArg(&imagePath),
		},
		Flags: []cli.Flag{
			androidFolderFlag(&folderName),
			outputNameFlagFn(&outputName, "ic_stat_notification_icon"),
			trimWhiteSpaceFlagFn(&trimWhiteSpace),
			alphaThresholdFlagFn(&alphaThreshold),
			applyFlagFn(&apply),
		},
	}
}

func applyAndroidNotificationIcon(folderName, outputName string) error {
	err := moveResAndroidOutFiles()
	if err != nil {
		return err
	}
	err = deleteAssetsGenOutDir()
	if err != nil {
		return err
	}

	fmt.Println("If You are using the Flutter framwork don't forget to add the following " +
		"meta-data schema within the application component: 'android/app/src/main/AndroidManifest.xml'")
	fmt.Println()
	schema := `<meta-data
  android:name="com.google.firebase.messaging.default_notification_icon"
  android:resource="@%s/%s" />`
	fmt.Printf(schema, folderName, outputName)
	fmt.Println()

	return nil
}
