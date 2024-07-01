// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package web

import (
	"github.com/eja/tibula/api"
)

func Router() (err error) {

	api.Plugins["test"] = func(eja api.TypeApi, db api.TypeDbSession) api.TypeApi {
		return eja
	}

	return
}
