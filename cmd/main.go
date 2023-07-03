package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	SavedImages []string `yaml:"savedImages"`
}

func main() {
	// Create the "scan-reports" folder if it doesn't exist
	err := createScanReportsFolder()
	if err != nil {
		log.Fatalf("Failed to create scan reports folder: %v", err)
	}

	// Read the manifest.yaml file
	data, err := os.ReadFile("manifest.yaml")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the YAML content
	var manifest Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Iterate over savedImages and process each Docker image
	for _, image := range manifest.SavedImages {
		// Extract image name and tag
		imageParts := strings.Split(image, ":")
		imageName := imageParts[0]
		imageTag := imageParts[1]

		// Check if the Docker image already exists locally
		existsLocally, err := isDockerImageExistsLocally(imageName, imageTag)
		if err != nil {
			log.Printf("Failed to check if image %s:%s exists locally: %v", imageName, imageTag, err)
			continue
		}

		if existsLocally {
			fmt.Printf("Image %s:%s already exists locally\n", imageName, imageTag)
		} else {
			// Pull the Docker image
			err := pullDockerImage(imageName, imageTag)
			if err != nil {
				log.Printf("Failed to pull image %s:%s: %v", imageName, imageTag, err)
				continue
			}
			fmt.Printf("Image %s:%s pulled successfully\n", imageName, imageTag)
		}

		// Perform Trivy scan
		err = trivyScan(imageName, imageTag)
		if err != nil {
			log.Printf("Failed to perform Trivy scan on image %s:%s: %v", imageName, imageTag, err)
			continue
		}
	}
}

func createScanReportsFolder() error {
	// Get the absolute path to the project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get project root directory: %v", err)
	}

	// Create the "scan-reports" folder if it doesn't exist
	scanReportsFolder := filepath.Join(projectRoot, "scan-reports")
	if _, err := os.Stat(scanReportsFolder); os.IsNotExist(err) {
		err := os.Mkdir(scanReportsFolder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create scan reports folder: %v", err)
		}
	}
	return nil
}

func isDockerImageExistsLocally(imageName, imageTag string) (bool, error) {
	// Execute "docker images" command to list locally available images
	cmd := exec.Command("docker", "images", fmt.Sprintf("%s:%s", imageName, imageTag), "--format", "{{.ID}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to list Docker images: %v, output: %s", err, string(output))
	}

	// If the output contains any data, the image exists locally
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func pullDockerImage(imageName, imageTag string) error {
	// Execute "docker pull" command
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", imageName, imageTag))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull Docker image: %v, output: %s", err, string(output))
	}
	return nil
}

func trivyScan(imageName, imageTag string) error {
	// Combine the image name and tag together
	imageFullName := fmt.Sprintf("%s:%s", imageName, imageTag)
	// Execute "trivy" command to perform the scan
	cmd := exec.Command("trivy", "image", imageFullName)

	// Create a file to store the Trivy report
	reportFileName := getReportFileName(imageName, imageTag)

	reportFile, err := os.Create(reportFileName)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer reportFile.Close()

	// Pipe the Trivy output to the report file and also to os.Stdout
	cmd.Stdout = io.MultiWriter(reportFile, os.Stdout)
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run Trivy: %v", err)
	}

	fmt.Printf("Trivy scan report for image %s:%s saved to %s\n", imageName, imageTag, reportFileName)
	return nil
}

func getReportFileName(imageName, imageTag string) string {
	// Extract the image name without the repository prefix
	imageNameParts := strings.Split(imageName, "/")
	imageName = imageNameParts[len(imageNameParts)-1]

	// Extract the tag name without any special characters
	tag := strings.ReplaceAll(imageTag, ":", "")

	// Construct the report file name
	reportFileName := fmt.Sprintf("%s-%s-trivy-report.txt", imageName, tag)

	// Get the absolute path of the report file
	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", "scan-reports", reportFileName))
	if err != nil {
		log.Fatalf("Failed to get absolute path for report file: %v", err)
	}
	return absPath
}
