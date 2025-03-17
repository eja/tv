// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func install() error {
	if _, err := os.Stat("/usr/bin/apt"); os.IsNotExist(err) {
		return fmt.Errorf("this installation does only work on Debian or a derivate")
	}

	if _, err := os.Stat("/sys/kernel/debug"); os.IsNotExist(err) {
		return fmt.Errorf("you need root/sudo privileges to execute this procedure")
	}

	commands := []string{
		"apt update",
		"apt install -y psmisc",
		"apt install -y wget",
	}

	for _, cmd := range commands {
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return fmt.Errorf("failed to execute command %s: %w\n", cmd, err)
		}
	}

	var tsduckURL string
	switch runtime.GOARCH {
	case "arm64":
		tsduckURL = "https://github.com/tsduck/tsduck/releases/download/v3.39-3956/tsduck_3.39-3956.debian12_arm64.deb"
	case "amd64":
		tsduckURL = "https://github.com/tsduck/tsduck/releases/download/v3.39-3956/tsduck_3.39-3956.debian12_amd64.deb"
	default:
		return fmt.Errorf("Unsupported architecture")
	}

	if err := exec.Command("wget", "-O", "/tmp/tsduck.deb", tsduckURL).Run(); err != nil {
		return fmt.Errorf("failed to download tsduck: %w\n", err)
	}

	if err := exec.Command("apt", "install", "-y", "/tmp/tsduck.deb").Run(); err != nil {
		return fmt.Errorf("failed to install tsduck: %w\n", err)
	}

	return nil
}
