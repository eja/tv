// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"github.com/eja/tibula/sys"
)

type typeConfigSys struct {
	sys.TypeConfig
	TvMediaPath     string `json:"tv_media_path,omitempty"`
	TvCheckInterval int    `json:"tv_check_interval,omitempty"`
	TvCheckBatch    int    `json:"tv_check_batch,omitempty"`
}

var String = sys.String
var Number = sys.Number
var Float = sys.Float
