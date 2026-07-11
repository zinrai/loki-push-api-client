package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type Config struct {
	Labels        []string `yaml:"labels"`
	Tenants       []string `yaml:"tenants"`
	Endpoint      string   `yaml:"endpoint"`
	SleepInterval int      `yaml:"sleep_interval"`
}

type PushRequest struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Stream StreamInfo `json:"stream"`
	Values [][]string `json:"values"`
}

type StreamInfo struct {
	Label string `json:"label"`
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s, commit: %s, date: %s\n", version, commit, date)
		return
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	client := &http.Client{}

	seq := 0
	for {
		orgID := config.Tenants[rand.Intn(len(config.Tenants))]

		var values [][]string
		for i := 0; i < rand.Intn(10)+1; i++ {
			seq++
			epoch := time.Now().UTC().UnixNano()
			line := fmt.Sprintf("tenant=%s seq=%d %s", orgID, seq, generateRandomString(30))
			values = append(values, []string{fmt.Sprintf("%d", epoch), line})
		}

		reqBody, err := json.Marshal(PushRequest{
			Streams: []Stream{
				{
					Stream: StreamInfo{
						Label: config.Labels[rand.Intn(len(config.Labels))],
					},
					Values: values,
				},
			},
		})
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", config.Endpoint, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Scope-OrgID", orgID)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Response Status: %s, X-Scope-OrgID: %s\n", resp.Status, orgID)

		time.Sleep(time.Duration(config.SleepInterval) * time.Second)
	}
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
