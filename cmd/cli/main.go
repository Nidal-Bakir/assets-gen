package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Nidal-Bakir/assets-gen/cmd/cli/cmd"
	"github.com/urfave/cli/v3"
)

var cmdErrors = []error{
	cmd.ErrInvalidBgType,
	cmd.ErrInvalidAndroidFolder,
	cmd.ErrNigativeValueCorners,
	cmd.ErrPaddingOutOfRange,
	cmd.ErrInvalidColor,
}

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			cmd.AndroidAppIcon(),
			cmd.AndroidNotificationIcon(),
			cmd.AndroidAssetGen(),
			cmd.IosAppIcon(),
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil && !isCmdDefinedError(err) {
		log.Fatal(err)
	}
}

func isCmdDefinedError(err error) bool {
	for _, e := range cmdErrors {
		if strings.Contains(err.Error(), e.Error()) {
			return true
		}
	}
	return false
}
