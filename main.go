package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide your chat_history.json file")
		return
	}

	filePath := os.Args[1]

	userFilter := ""
	contentFilter := ""

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-u":
			if i+1 < len(os.Args) {
				userFilter = os.Args[i+1]
				i++
			}
		case "-s":
			if i+1 < len(os.Args) {
				contentFilter = os.Args[i+1]
				i++
			}
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data := make(map[string][]map[string]interface{})

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	messageCounts := make(map[string]int)

	for user, messages := range data {
		if userFilter != "" && user != userFilter {
			continue
		}

		count := 0
		for _, message := range messages {
			if mediaType, ok := message["Media Type"].(string); ok && mediaType == "TEXT" {
				if userFilter == "" || message["From"] == userFilter {
					if content, ok := message["Content"].(string); ok {
						if contentFilter == "" || strings.Contains(strings.ToLower(content), strings.ToLower(contentFilter)) {
							count++
						}
					}
				}
			}
		}
		messageCounts[user] = count
	}

	for user, count := range messageCounts {
		fmt.Printf("User: %s, Message Count: %d\n", user, count)
	}
}
