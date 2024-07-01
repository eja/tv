// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"github.com/eja/tibula/sys"
)

type typeConfigSys struct {
	sys.TypeConfig
	MediaPath     string `json:"media_path,omitempty"`
	CheckInterval int    `json:"check_interval,omitempty"`
}

var String = sys.String
var Number = sys.Number
var Float = sys.Float
