// Copyright by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ApiRequest struct {
	UUID    string `json:"uuid"`
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

type ApiResponse struct {
	Url              string  `json:"url"`
	Frequency        float64 `json:"frequency"`
	Bandwidth        int     `json:"bandwidth"`
	Constellation    int     `json:"constellation"`
	CodeRate         int     `json:"codeRate"`
	GuardInterval    int     `json:"guardInterval"`
	TransmissionMode int     `json:"transmissionMode"`
}

func start(apiUrl string, maxRetries int) error {
	rebootCount := 0

	uuid, err := macGet()
	if err != nil {
		return err
	}

	for {
		config, err := apiPost(apiUrl, uuid, "init", "modulator")
		if err != nil {
			log.Printf("api retrieval problem: %v", err)
		} else {
			if config.Url != "" {
				log.Println("injecting")
				codeRate := []string{"1/2", "2/3", "3/4", "5/6", "7/8"}
				constellation := []string{"QPSK", "16-QAM", "64-QAM"}
				transmissionMode := []string{"2K", "8K", "4K"}
				guardInterval := []string{"1/32", "1/16", "1/8", "1/4"}
				cmd := exec.Command("tsp",
					"-P", "bitrate_monitor", "--min", "1000000", "--alarm-command", "killall -9 tsp #",
					"-I", "http", "--compressed", config.Url,
					"-O", "vatek", "--modulation", "DVB-T", "--remux", "passthrough",
					"--frequency", fmt.Sprintf("%.0f", config.Frequency*1000),
					"--bandwidth", fmt.Sprintf("%d", config.Bandwidth/1000),
					"--constellation", constellation[config.Constellation],
					"--convolutional-rate", codeRate[config.CodeRate],
					"--guard-interval", guardInterval[config.GuardInterval],
					"--transmission-mode", transmissionMode[config.TransmissionMode],
				)
				if output, err := cmd.CombinedOutput(); err != nil {
					log.Printf("error executing command: %v\n%s\n", err, string(output))
				}

			} else {
				log.Printf("api configuration empty")
			}
		}
		rebootCount++
		if maxRetries > 0 && rebootCount > maxRetries {
			log.Println("too many restarts, rebooting")
			if err := exec.Command("/sbin/reboot", "--force").Run(); err != nil {
				return err
			}

		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func apiPost(url, uuid, action, payload string) (*ApiResponse, error) {
	requestBody := ApiRequest{
		UUID:    uuid,
		Action:  action,
		Payload: payload,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func macGet() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}
		if fields[1] == "00000000" && fields[2] != "00000000" && fields[7] == "00000000" {
			device := fields[0]
			filePath := fmt.Sprintf("/sys/class/net/%s/address", device)
			macFile, err := os.Open(filePath)
			if err != nil {
				return "", err
			}
			defer macFile.Close()

			macScanner := bufio.NewScanner(macFile)
			if macScanner.Scan() {
				return strings.ToUpper(strings.ReplaceAll(macScanner.Text(), ":", "")), nil
			}

			if err := macScanner.Err(); err != nil {
				return "", err
			}

			return "", fmt.Errorf("MAC address not found for device %s", device)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("default gateway device not found")
}
