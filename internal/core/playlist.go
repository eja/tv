// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func checkPlaylist(url string) (cors bool, subtitles bool, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, false, fmt.Errorf("error getting the playlist: %v", err)
	}
	defer resp.Body.Close()

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin == "*" {
		cors = true
	} else {
		cors = false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, fmt.Errorf("error reading response body: %v", err)
	}

	lines := strings.Split(string(body), "\n")
	for number, line := range lines {
		if number == 0 && !strings.HasPrefix(strings.ToLower(lines[0]), "#extm3u") {
			return false, false, fmt.Errorf("not a playlist file")
		}
		if strings.HasPrefix(strings.ToLower(line), "#ext-x-media:type=subtitles") {
			subtitles = true
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "#ext-x-media:type=closed-captions") {
			subtitles = true
			break
		}
	}

	return
}
