// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"tv/internal/av"
	"tv/internal/sys"

	tibula "github.com/eja/tibula/db"
	"github.com/eja/tibula/log"
)

const channelsBatch = 100

func checkChannels() (err error) {
	db := tibula.Session()
	if err = db.Open(sys.Options.DbType, sys.Options.DbName, sys.Options.DbUser, sys.Options.DbPass, sys.Options.DbHost, sys.Options.DbPort); err != nil {
		return
	}
	for {
		timeLastCheck := time.Now().Add(-time.Duration(sys.Options.CheckInterval) * time.Second).Format("2006-01-02 15:04:05")
		timeLastWorking := time.Now().Add(-30 * 24 * time.Hour).Format("2006-01-02 15:04:05")
		rows, err := db.Rows(`SELECT * FROM tvChannels WHERE 
		  (checkLast < ? OR checkLast IS NULL OR checkLast = "") AND 
			(
			 (checkLastWorking > ?) OR 
			 (ejaLog > ? AND (checkLastWorking IS NULL OR checkLastWorking = ""))
			) AND
			power > 0 AND 
			name != "" 
			ORDER BY power DESC, checkLast ASC 
			LIMIT ?
		`, timeLastCheck, timeLastWorking, timeLastWorking, channelsBatch)
		if err != nil {
			return err
		}
		for _, row := range rows {
			var videoWidth int64
			var videoHeight int64
			var videoSize string
			status := 0
			ABR := 0
			framePath := filepath.Join(sys.Options.MediaPath, row["name"]+".png")

			cors, subtitles, err := checkPlaylist(row["sourceUrl"])
			if err != nil {
				log.Warn(tag, "playlist check error", row["name"], err)
			} else {
				status += 1
				log.Trace(tag, "playlist ok", row["label"], cors, subtitles)
				if strings.HasPrefix(strings.ToLower(row["sourceUrl"]), "https") {
					status += 100
				}
				probeData, err := av.FFprobe([]string{"-timeout", "10000000", "-print_format", "json", "-show_format", "-show_streams", row["sourceUrl"]})
				if err != nil {
					log.Warn(tag, "probe check error", row["name"], err)
				} else {
					status += 10
					var probeJson map[string]interface{}
					err = json.Unmarshal(probeData, &probeJson)
					if err != nil {
						log.Warn(tag, "error unmarshalling ffprobe json", err)
					} else {
						av.FFmpeg([]string{"-timeout", "10000000", "-i", row["sourceUrl"], "-vframes", "1", framePath})
						if streams, ok := probeJson["streams"].([]interface{}); !ok {
							log.Warn(tag, "json streams not found")
						} else {
							for _, streamMap := range streams {
								if stream, ok := streamMap.(map[string]interface{}); ok {
									if stream["codec_type"] == "video" {
										w := sys.Number(stream["width"])
										h := sys.Number(stream["height"])
										if w > videoWidth {
											videoWidth = w
										}
										if h > videoHeight {
											videoHeight = h
										}
										ABR++
										log.Trace(stream["width"], stream["height"])
									}
								}
							}
							if videoWidth > 0 && videoHeight > 0 {
								videoSize = fmt.Sprintf("%dx%d", videoWidth, videoHeight)
							}
						}
					}
				}
			}
			db.Run("UPDATE tvChannels SET status=?,checkLast=?,size=?,subtitle=?,abr=? WHERE name=?", status, db.Now(), videoSize, subtitles, ABR, row["name"])
			log.Trace(tag, "status:", row["country"], row["label"], status, videoSize)
			if status > 0 {
				db.Run("UPDATE tvChannels SET checkLastWorking=? WHERE name=?", db.Now(), row["name"])
			}
		}

		log.Trace(tag, "sleep")
		time.Sleep(10 * time.Second)
	}

	return
}
