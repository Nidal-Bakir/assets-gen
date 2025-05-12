package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nidal-Bakir/assets-gen/cmd/assetsgen/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	startTime := time.Now()

	cmd := &cli.Command{
		Usage:   "A CLI that will help you generate app icons and images for various platforms",
		Version: "v1.0.1",
		Commands: []*cli.Command{
			cmd.AndroidAppIcon(),
			cmd.AndroidNotificationIcon(),
			cmd.AndroidAssetGen(),
			cmd.IosAppIcon(),
			cmd.AndroidGooglePlayLogo(),
			cmd.GenerateAll(),
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ntook %.2fsec\n", time.Since(startTime).Seconds())
}
