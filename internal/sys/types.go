// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"github.com/eja/tibula/sys"
)

type typeConfigSys struct {
	sys.TypeConfig
	MediaPath string `json:"media_path,omitempty"`
	TmpPath   string `json:"tmp_path,omitempty"`
}

var String = sys.String
var Number = sys.Number
var Float = sys.Float
