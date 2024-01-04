package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	checkFileName   = "check.txt"
	validFileName   = "valid.txt"
	invalidFileName = "invalid.txt"
)

func main() {
	startTime := time.Now()

	color.Red("\n\n----------------------")
	color.Cyan("Starting...")

	links, err := readLinksFromFile(checkFileName)
	if err != nil {
		color.Red("Failed, cannot load check file, error: %s", err)
		return
	}

	var count = 0
	for _, link := range links {
		if isValidLink(link) {
			color.Green("[+] %s", link)
			appendToValidFile(link)
		} else {
			color.Red("[-] %s", link)
			appendToInvalidFile(link)
		}
		count++
	}

	if count < 1 {
		color.Red("No results")
	}

	color.Cyan("Program finished..")
	color.Red("----------------------")
	color.Cyan("\n\nBy Zile42O")

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Program took: %d:%02d\n", int(elapsedTime.Minutes()), int(elapsedTime.Seconds())%60)
}

func readLinksFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func isValidLink(link string) bool {
	parts := strings.Split(link, "/")
	if len(parts) < 5 {
		return false
	}
	code := parts[4]
	apiURL := fmt.Sprintf("https://discordapp.com/api/invite/%s", code)
	response, err := http.Get(apiURL)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	if message, ok := result["message"].(string); ok && strings.Contains(strings.ToLower(message), "unknown invite") {
		return false
	}
	return true
}

func appendToFile(filename, link string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("Error while append to file, error: %s", err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(link + "\n"); err != nil {
		color.Red("Error while writing to file, error: %s", err)
	}
}

func appendToValidFile(link string) {
	appendToFile(validFileName, link)
}

func appendToInvalidFile(link string) {
	appendToFile(invalidFileName, link)
}
