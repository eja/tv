// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"embed"

	tibula "github.com/eja/tibula/db"
	"github.com/eja/tibula/sys"
)

//go:embed all:assets
var dbAssets embed.FS

func Wizard() error {
	configFile := sys.Options.ConfigFile
	if err := sys.ConfigRead(configFile, &Options); err != nil {
		return err
	}

	Options.TvMediaPath = sys.WizardPrompt("Media folder path")
	checkInterval := sys.WizardPrompt("Channels check interval")
	if checkInterval != "" {
		Options.TvCheckInterval = int(Number(checkInterval))
	}
	checkBatch := sys.WizardPrompt("Channels check batch size")
	if checkBatch != "" {
		Options.TvCheckBatch = int(Number(checkBatch))
	}

	tibula.Assets = dbAssets

	db := tibula.Session()
	if err := db.Open(Options.DbType, Options.DbName, Options.DbUser, Options.DbPass, Options.DbHost, Options.DbPort); err != nil {
		return err
	}
	if err := db.Setup(""); err != nil {
		return err
	}

	Options.ConfigFile = ""
	return sys.ConfigWrite(configFile, &Options)
}
