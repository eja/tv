// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package av

import (
	"os/exec"

	"github.com/eja/tibula/log"
)

const tag = "[FF]"

func FFmpeg(args []string) error {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner"}
	cmd := exec.Command("ffmpeg", append(baseArgs, args...)...)
	log.Trace(tag, "ffmpeg", args)
	return cmd.Run()
}

func FFprobe(args []string) ([]byte, error) {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner", "-v", "error"}
	cmd := exec.Command("ffprobe", append(baseArgs, args...)...)
	log.Trace(tag, "ffprobe", args)
	return cmd.Output()
}
