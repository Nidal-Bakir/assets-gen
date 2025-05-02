package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Nidal-Bakir/assets-gen/cmd/cli/cmd"
	"github.com/urfave/cli/v3"
)

var cmdErrors = []error{
	cmd.ErrInvalidBgType,
	cmd.ErrInvalidAndroidFolder,
	cmd.ErrInvalidValueRange,
	cmd.ErrPaddingOutOfRange,
	cmd.ErrInvalidColor,
	cmd.ErrColorsAndStopsLengthDidNotMatch,
}

func main() {
	startTime := time.Now()

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

	fmt.Printf("\ntook %.2fsec\n", time.Since(startTime).Seconds())
}

func isCmdDefinedError(err error) bool {
	for _, e := range cmdErrors {
		if strings.Contains(err.Error(), e.Error()) {
			return true
		}
	}
	return false
}
