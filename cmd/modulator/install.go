// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"io"
	"net/http"
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
	}

	for _, cmd := range commands {
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return fmt.Errorf(`failed to execute command "%s": %w\n`, cmd, err)
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

	tsduckFilePath := "/tmp/tsduck.deb"
	if err := downloadFile(tsduckURL, tsduckFilePath); err != nil {
		return fmt.Errorf("failed to download tsduck: %w", err)
	}

	if err := exec.Command("apt", "install", "-y", tsduckFilePath).Run(); err != nil {
		return fmt.Errorf("failed to install tsduck: %w", err)
	}

	return nil
}

func downloadFile(url, filePath string) error {
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %w", err)
	}

	return nil
}
