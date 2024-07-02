// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"flag"
	"os"

	"github.com/eja/tibula/sys"
)

var Options typeConfigSys

func Configure() error {
	flag.StringVar(&Options.TvMediaPath, "tv-media-path", "/opt/eja/tv/media", "Media folder path")
	flag.IntVar(&Options.TvCheckInterval, "tv-check-interval", 3600, "Channels check interval")
	flag.IntVar(&Options.TvCheckBatch, "tv-check-batch", 10, "Channels check batch")

	if err := sys.Configure(); err != nil {
		return err
	}
	Options.TypeConfig = sys.Options

	if sys.Commands.Start && sys.Options.ConfigFile == "" {
		sys.Options.ConfigFile = Name + ".json"
		if _, err := os.Stat(sys.Options.ConfigFile); err != nil {
			sys.Options.ConfigFile = ""
		}
	}

	if sys.Options.ConfigFile != "" {
		if err := sys.ConfigRead(sys.Options.ConfigFile, &Options); err != nil {
			return err
		}
	}

	return nil
}
