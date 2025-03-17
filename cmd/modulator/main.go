// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	Name    = "modulator"
	Version = "9.3.17"
)

var (
	Install bool
	Start   bool
	ApiUrl  string
	Retries int
)

func main() {
	flag.BoolVar(&Install, "install", false, "Install the modulator")
	flag.BoolVar(&Start, "start", false, "Start injecting")
	flag.StringVar(&ApiUrl, "api", "https://api.eja.tv", "API server URL")
	flag.IntVar(&Retries, "retry", 0, "Max retries before rebooting (default don't reboot)")
	flag.Parse()

	if Install {
		if err := install(); err != nil {
			log.Fatalf("Cannot install: %v", err)
		}
	} else if Start {
		if err := start(ApiUrl, Retries); err != nil {
			log.Fatalf("Cannot start: %v", err)
		}
	} else {
		fmt.Println("Copyright:", "Ubaldo Porcheddu <ubaldo@eja.it>")
		fmt.Printf("Version: %s\nUsage: %s [options]\n\nOptions:\n", Version, Name)
		flag.PrintDefaults()
		fmt.Println()
	}
}
