package main

import (
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Labels   []string `yaml:"labels"`
	Tenants  []string `yaml:"tenants"`
	Endpoint string   `yaml:"endpoint"`
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
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomString)
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	client := &http.Client{}

	for {
		epoch := time.Now().UTC().UnixNano()
		var values [][]string
		for i := 0; i < rand.Intn(10)+1; i++ {
			values = append(values, []string{fmt.Sprintf("%d", epoch), generateRandomString(30)})
		}

		orgID := config.Tenants[rand.Intn(len(config.Tenants))]

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

		time.Sleep(3 * time.Second)
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
