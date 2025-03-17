// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"tv/internal/core"
	"tv/internal/sys"
	"tv/internal/web"

	"github.com/eja/tibula/log"
	tibulaSys "github.com/eja/tibula/sys"
	tibulaWeb "github.com/eja/tibula/web"
)

func main() {
	if err := sys.Configure(); err != nil {
		log.Fatal(err)
	}

	if tibulaSys.Commands.DbSetup {
		if err := tibulaSys.Setup(); err != nil {
			log.Fatal(err)
		}
	} else if tibulaSys.Commands.Wizard {
		if err := tibulaSys.WizardSetup(); err != nil {
			log.Fatal(err)
		}
		if err := sys.Wizard(); err != nil {
			log.Fatal(err)
		}

	} else if tibulaSys.Commands.Start {
		if sys.Options.DbName == "" {
			log.Fatal("Database name/file is mandatory.")
		}
		if err := web.Router(); err != nil {
			log.Fatal(err)
		}
		if sys.Options.TvCheck {
			if err := core.Start(); err != nil {
				log.Fatal(err)
			}
		}
		if err := tibulaWeb.Start(); err != nil {
			log.Fatal("Cannot start the web service: ", err)
		}
	} else {
		sys.Help()
	}
}
