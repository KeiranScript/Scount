package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define command-line arguments
	filePath := flag.String("file", "", "Path to the JSON file")
	userFilter := flag.String("u", "", "Username to filter messages by")
	contentFilter := flag.String("s", "", "Substring to search for in message content")
	flag.Parse()

	// Check if the file path was provided
	if *filePath == "" {
		fmt.Println("Please provide the path to the JSON file using the -file argument.")
		return
	}

	// Open the JSON file
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Define a map to hold the parsed JSON data
	data := make(map[string][]map[string]interface{})

	// Decode the entire JSON file into the map
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Store message counts for each user
	messageCounts := make(map[string]int)

	// Iterate over the map to count the messages
	for user, messages := range data {
		// Skip users if we're filtering by user and the user doesn't match the filter
		if *userFilter != "" && user != *userFilter {
			continue
		}

		count := 0
		for _, message := range messages {
			// Check if the "Media Type" is "TEXT" and if the "From" field matches the filter (if provided)
			if mediaType, ok := message["Media Type"].(string); ok && mediaType == "TEXT" {
				if *userFilter == "" || message["From"] == *userFilter {
					// Check if "Content" exists and is a string
					if content, ok := message["Content"].(string); ok {
						// If a content filter is provided, check if it is in the "Content" field
						if *contentFilter == "" || strings.Contains(strings.ToLower(content), strings.ToLower(*contentFilter)) {
							count++
						}
					}
				}
			}
		}
		messageCounts[user] = count
	}

	// Output the message counts
	for user, count := range messageCounts {
		fmt.Printf("User: %s, Message Count: %d\n", user, count)
	}
}

