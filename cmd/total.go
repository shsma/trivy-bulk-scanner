package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Specify the folder containing the Trivy report files
	folderPath := "scan-reports"

	// Initialize counters
	totalUnknown := 0
	totalLow := 0
	totalMedium := 0
	totalHigh := 0
	totalCritical := 0

	// Get the list of files in the folder
	fileList, err := getFileList(folderPath)
	if err != nil {
		log.Fatalf("Failed to get file list: %v", err)
	}

	// Process each Trivy report file
	for _, fileName := range fileList {
		// Read the Trivy report file
		file, err := os.Open(fileName)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			continue
		}
		defer file.Close()

		// Read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			// Check if the line contains severity information
			if strings.HasPrefix(line, "Total: ") {
				// Extract the severity counts
				severityCounts := strings.TrimPrefix(line, "Total: ")
				severityCounts = strings.TrimSuffix(severityCounts, ")")

				// Split the severity counts
				counts := strings.Split(severityCounts, ", ")

				// Iterate over the severity counts and update the corresponding total
				for _, count := range counts {
					parts := strings.Split(count, ": ")
					if len(parts) == 2 {
						severity := strings.TrimSpace(parts[0])
						count := strings.TrimSpace(parts[1])

						switch severity {
						case "UNKNOWN":
							totalUnknown += parseInt(count)
						case "LOW":
							totalLow += parseInt(count)
						case "MEDIUM":
							totalMedium += parseInt(count)
						case "HIGH":
							totalHigh += parseInt(count)
						case "CRITICAL":
							totalCritical += parseInt(count)
						}
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read file: %v", err)
		}
	}

	// Print the totals
	fmt.Println("Totals:")
	fmt.Printf("UNKNOWN: %d\n", totalUnknown)
	fmt.Printf("LOW: %d\n", totalLow)
	fmt.Printf("MEDIUM: %d\n", totalMedium)
	fmt.Printf("HIGH: %d\n", totalHigh)
	fmt.Printf("CRITICAL: %d\n", totalCritical)
}

func getFileList(folderPath string) ([]string, error) {
	var fileList []string

	// Walk through the folder and get the list of files
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileList = append(fileList, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func parseInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		log.Fatalf("Failed to parse integer: %v", err)
	}
	return result
}
