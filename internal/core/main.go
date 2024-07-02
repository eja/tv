// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package core

import (
	"os"

	"tv/internal/sys"
)

const tag = "[tv] [core]"

func Start() (err error) {
	if _, err = os.Stat(sys.Options.TvMediaPath); err != nil {
		err = os.MkdirAll(sys.Options.TvMediaPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return checkChannels()
}
