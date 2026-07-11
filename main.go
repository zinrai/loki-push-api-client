package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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

func generateNames(prefix string, n int) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("%s%d", prefix, i+1)
	}
	return names
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
	endpoint := flag.String("endpoint", "http://localhost:3100/loki/api/v1/push", "Loki push API endpoint")
	tenantCount := flag.Int("tenants", 5, "number of tenants to generate")
	labelCount := flag.Int("labels", 5, "number of labels to generate")
	interval := flag.Duration("interval", 30*time.Second, "interval between requests")
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s, commit: %s, date: %s\n", version, commit, date)
		return
	}

	if *tenantCount < 1 || *labelCount < 1 {
		log.Fatal("tenants and labels must be at least 1")
	}

	tenants := generateNames("tenant", *tenantCount)
	labels := generateNames("label", *labelCount)

	fmt.Println("tenants:", strings.Join(tenants, " "))
	fmt.Println("labels:", strings.Join(labels, " "))

	client := &http.Client{}

	seq := 0
	for {
		orgID := tenants[rand.Intn(len(tenants))]

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
						Label: labels[rand.Intn(len(labels))],
					},
					Values: values,
				},
			},
		})
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(reqBody))
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
		resp.Body.Close()

		log.Printf("Response Status: %s, X-Scope-OrgID: %s\n", resp.Status, orgID)

		time.Sleep(*interval)
	}
}
